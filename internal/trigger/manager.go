package trigger

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ironpark/teatime/internal/recipe"
)

type Manager struct {
	// Dependencies
	recipeLoader RecipeLoader
	recipeRunner RecipeRunner

	// Trigger handlers
	handlers map[TriggerType]Handler

	// Active triggers management
	activeTriggers map[string]*Instance
	recipeTriggers map[string][]*Instance

	// Optional state management
	store Store
	mu    sync.RWMutex

	// Execution state
	running bool
}

// NewManager creates a new trigger manager
func NewManager(loader RecipeLoader, runner RecipeRunner, handlers []Handler) *Manager {
	// Convert slice to map
	handlersMap := make(map[TriggerType]Handler)
	for _, handler := range handlers {
		handlersMap[handler.Type()] = handler
	}

	m := &Manager{
		recipeLoader:   loader,
		recipeRunner:   runner,
		handlers:       handlersMap,
		activeTriggers: make(map[string]*Instance),
		recipeTriggers: make(map[string][]*Instance),
	}

	// Initialize handlers
	for _, handler := range handlers {
		if err := handler.Initialize(m); err != nil {
			log.Printf("Failed to initialize handler %s: %v", handler.Type(), err)
		}
	}

	return m
}

// SetStore sets the optional store for persistence
func (m *Manager) SetStore(store Store) {
	m.store = store
}

// Start starts the trigger manager
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("trigger manager already running")
	}
	wg := sync.WaitGroup{}
	// Start all handlers
	for _, handler := range m.handlers {
		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()
			if err := h.Run(ctx); err != nil {
				log.Printf("Handler %s run error: %v", h.Type(), err)
			}
		}(handler)
	}
	m.mu.Unlock()
	wg.Wait()
	// Unregister all triggers
	m.mu.Lock()
	for recipeID := range m.recipeTriggers {
		m.unregisterRecipeInternal(recipeID)
	}
	m.running = true
	m.mu.Unlock()
	log.Println("Trigger manager started")
	return nil
}

// RegisterRecipe registers all triggers for a recipe
func (m *Manager) RegisterRecipe(recipe *recipe.Recipe) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	recipeID := recipe.Path // Use path as recipe identifier

	// Cleanup existing triggers for this recipe
	m.unregisterRecipeInternal(recipeID)

	// Find trigger nodes
	triggerNodes := m.findTriggerNodes(recipe.Nodes)
	if len(triggerNodes) == 0 {
		return nil // No triggers to register
	}

	var instances []*Instance

	// Register each trigger node
	for _, node := range triggerNodes {
		instance, err := m.registerTriggerNode(recipeID, node)
		if err != nil {
			// Rollback on partial failure
			for _, inst := range instances {
				m.unregisterTriggerInternal(inst)
			}
			return fmt.Errorf("failed to register trigger %s: %w", node.Id, err)
		}
		instances = append(instances, instance)
	}

	// Store mapping on success
	m.recipeTriggers[recipeID] = instances

	log.Printf("Registered %d triggers for recipe %s", len(instances), recipeID)
	return nil
}

// UnregisterRecipe unregisters all triggers for a recipe
func (m *Manager) UnregisterRecipe(recipeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.unregisterRecipeInternal(recipeID)
}

// registerTriggerNode registers a single trigger node
func (m *Manager) registerTriggerNode(recipeID string, node recipe.Node) (*Instance, error) {
	// Determine trigger type
	triggerType := m.determineTriggerType(node.Ref)
	handler, exists := m.handlers[triggerType]
	if !exists {
		return nil, fmt.Errorf("unsupported trigger type: %s", triggerType)
	}

	// Convert node properties to config map
	config := make(map[string]any)
	for _, prop := range node.Properties {
		config[prop.Key] = prop.Value
	}

	// Create instance
	instance := &Instance{
		ID:        generateTriggerID(recipeID, node.Id),
		RecipeID:  recipeID,
		NodeID:    node.Id,
		Type:      triggerType,
		Config:    config,
		Status:    StatusInactive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate configuration
	if err := handler.Validate(instance.Config); err != nil {
		return nil, fmt.Errorf("invalid trigger config: %w", err)
	}

	// Register with handler
	if err := handler.Register(instance); err != nil {
		return nil, err
	}

	// Mark as active and store in memory
	instance.Status = StatusActive
	instance.UpdatedAt = time.Now()
	m.activeTriggers[instance.ID] = instance

	return instance, nil
}

// unregisterRecipeInternal unregisters all triggers for a recipe (internal, assumes lock held)
func (m *Manager) unregisterRecipeInternal(recipeID string) error {
	instances, exists := m.recipeTriggers[recipeID]
	if !exists {
		return nil
	}

	// Unregister each trigger
	for _, instance := range instances {
		if err := m.unregisterTriggerInternal(instance); err != nil {
			log.Printf("Failed to unregister trigger %s: %v", instance.ID, err)
		}
	}

	// Remove mapping
	delete(m.recipeTriggers, recipeID)

	log.Printf("Unregistered triggers for recipe %s", recipeID)
	return nil
}

// unregisterTriggerInternal unregisters a single trigger (internal, assumes lock held)
func (m *Manager) unregisterTriggerInternal(instance *Instance) error {
	// Unregister from handler
	if handler, exists := m.handlers[instance.Type]; exists {
		if err := handler.Unregister(instance); err != nil {
			return err
		}
	}

	// Call cleanup function if exists
	if instance.cleanup != nil {
		if err := instance.cleanup(); err != nil {
			return err
		}
	}

	// Remove from memory
	delete(m.activeTriggers, instance.ID)
	instance.Status = StatusInactive

	return nil
}

// ExecuteTrigger executes a trigger (called by handlers)
func (m *Manager) ExecuteTrigger(ctx context.Context, triggerID string, data map[string]any) {
	// Execute asynchronously to prevent blocking
	go func() {
		if err := m.executeTriggerSync(ctx, triggerID, data); err != nil {
			log.Printf("Trigger execution failed [%s]: %v", triggerID, err)
		}
	}()
}

// executeTriggerSync synchronously executes a trigger
func (m *Manager) executeTriggerSync(ctx context.Context, triggerID string, data map[string]any) error {
	m.mu.RLock()
	instance, exists := m.activeTriggers[triggerID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("trigger not found: %s", triggerID)
	}

	// Update execution statistics
	now := time.Now()
	instance.TriggerCount++
	instance.LastTriggered = &now
	instance.LastError = ""

	// Load recipe
	recipe, err := m.recipeLoader.LoadRecipe(instance.RecipeID)
	if err != nil {
		instance.LastError = err.Error()
		return fmt.Errorf("failed to load recipe %s: %w", instance.RecipeID, err)
	}

	// Execute recipe
	ctx, cancel := context.WithTimeout(ctx, 60*time.Minute)
	defer cancel()

	err = m.recipeRunner.Execute(ctx, recipe, instance.NodeID, data)
	if err != nil {
		instance.LastError = err.Error()
		return fmt.Errorf("recipe execution failed: %w", err)
	}

	// Update success timestamp
	instance.UpdatedAt = now
	log.Printf("Trigger %s executed successfully", triggerID)

	return nil
}

// GetActiveTriggers returns all active triggers
func (m *Manager) GetActiveTriggers() []*Instance {
	m.mu.RLock()
	defer m.mu.RUnlock()

	triggers := make([]*Instance, 0, len(m.activeTriggers))
	for _, trigger := range m.activeTriggers {
		triggers = append(triggers, trigger)
	}

	return triggers
}

// GetTriggersByRecipe returns triggers for a specific recipe
func (m *Manager) GetTriggersByRecipe(recipeID string) []*Instance {
	m.mu.RLock()
	defer m.mu.RUnlock()

	instances, exists := m.recipeTriggers[recipeID]
	if !exists {
		return []*Instance{}
	}

	// Return a copy to prevent external modification
	result := make([]*Instance, len(instances))
	copy(result, instances)
	return result
}

// GetHandler returns a handler by type
func (m *Manager) GetHandler(triggerType TriggerType) Handler {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.handlers[triggerType]
}

// SaveTriggerStats saves trigger statistics (optional)
func (m *Manager) SaveTriggerStats() error {
	if m.store == nil {
		return nil // No store configured
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var stats []TriggerStat
	for _, instance := range m.activeTriggers {
		if instance.TriggerCount > 0 { // Only save triggers that have been executed
			stats = append(stats, TriggerStat{
				TriggerID:     instance.ID,
				RecipeID:      instance.RecipeID,
				TriggerCount:  instance.TriggerCount,
				LastTriggered: instance.LastTriggered,
				LastError:     instance.LastError,
			})
		}
	}

	return m.store.SaveStats(stats)
}

// Helper functions

// findTriggerNodes finds all trigger nodes in a recipe
func (m *Manager) findTriggerNodes(nodes []recipe.Node) []recipe.Node {
	var triggerNodes []recipe.Node
	for _, node := range nodes {
		if m.isTriggerNode(node.Ref) {
			triggerNodes = append(triggerNodes, node)
		}
	}
	return triggerNodes
}

// isTriggerNode checks if a node reference is a trigger node
func (m *Manager) isTriggerNode(nodeRef string) bool {
	return strings.HasPrefix(nodeRef, "teatime.trigger.")
}

// determineTriggerType determines the trigger type from node reference
func (m *Manager) determineTriggerType(nodeRef string) TriggerType {
	switch {
	case strings.Contains(nodeRef, ".webhook"):
		return TypeWebhook
	case strings.Contains(nodeRef, ".schedule"):
		return TypeSchedule
	case strings.Contains(nodeRef, ".command"):
		return TypeCommand
	case strings.Contains(nodeRef, ".filewatch"):
		return TypeFileWatch
	default:
		return TriggerType(nodeRef) // fallback
	}
}

// generateTriggerID generates a deterministic trigger ID
func generateTriggerID(recipeID, nodeID string) string {
	return fmt.Sprintf("trigger_%s_%s", recipeID, nodeID)
}
