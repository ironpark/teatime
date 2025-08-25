package trigger

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	rc "github.com/ironpark/teatime/internal/recipe"
	"github.com/ironpark/teatime/internal/runner"
)

// Manager coordinates trigger handlers and manages trigger instances for recipes.
// It provides registration, lifecycle management, and execution orchestration for triggers.
type Manager struct {
	// Dependencies
	recipeLoader RecipeLoader
	registry     *Registry

	// Active triggers management
	activeTriggers map[string]*Instance
	recipeTriggers map[string][]*Instance

	// Event channel for trigger events
	eventCh chan Event

	// Optional state management
	store Store
	mu    sync.RWMutex

	// Execution state
	running bool
}

// NewManager creates a new trigger manager with internal registry.
func NewManager(loader RecipeLoader) *Manager {
	m := &Manager{
		recipeLoader:   loader,
		registry:       NewRegistry(),
		activeTriggers: make(map[string]*Instance),
		recipeTriggers: make(map[string][]*Instance),
		eventCh:        make(chan Event, 100), // Buffered channel for trigger events
	}

	return m
}

// NewManagerWithHandlers creates a new trigger manager with specific handlers (legacy).
// This function is provided for backward compatibility. New code should use NewManager and RegisterHandler.
func NewManagerWithHandlers(loader RecipeLoader, handlers []Handler) *Manager {
	m := NewManager(loader)

	// Register handlers (initialization handled by registry)
	for _, handler := range handlers {
		if err := m.RegisterHandler(handler); err != nil {
			log.Printf("Failed to register handler %s: %v", handler.NodeRef(), err)
			continue
		}
	}

	return m
}

// SetStore sets the optional store for persistence
func (m *Manager) SetStore(store Store) {
	m.store = store
}

// Start starts the trigger manager and all registered handlers.
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("trigger manager already running")
	}

	// Start all handlers in background
	handlers := m.registry.List()
	for _, handlerInfo := range handlers {
		handler, err := m.registry.GetHandler(handlerInfo.NodeRef)
		if err != nil {
			log.Printf("Failed to get handler %s: %v", handlerInfo.NodeRef, err)
			continue
		}
		go func(h Handler) {
			if err := h.Start(ctx); err != nil {
				log.Printf("Handler %s start error: %v", h.NodeRef(), err)
			}
		}(handler)
	}

	m.running = true
	log.Println("Trigger manager started")

	// Start event processing goroutine
	go m.processEvents(ctx)

	// Wait for context cancellation to cleanup
	go m.waitForShutdown(ctx)

	return nil
}

// processEvents processes trigger events from the event channel
func (m *Manager) processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-m.eventCh:
			m.handleTriggerEvent(ctx, event)
		}
	}
}

// handleTriggerEvent handles a single trigger event
func (m *Manager) handleTriggerEvent(ctx context.Context, event Event) {
	// Execute asynchronously to prevent blocking
	go func() {
		if err := m.executeTriggerSync(ctx, event.TriggerID, event.Data); err != nil {
			log.Printf("Trigger execution failed [%s]: %v", event.TriggerID, err)
		}
	}()
}

// waitForShutdown waits for context cancellation and cleans up resources.
func (m *Manager) waitForShutdown(ctx context.Context) {
	<-ctx.Done()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Unregister all triggers
	for recipeID := range m.recipeTriggers {
		m.unregisterRecipeInternal(recipeID)
	}

	m.running = false
	log.Println("Trigger manager stopped")
}

// RegisterRecipe registers all triggers for a recipe
func (m *Manager) RegisterRecipe(recipe *rc.Recipe) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	recipeID := recipe.Path // Use path as recipe identifier

	// Cleanup existing triggers for this recipe
	m.unregisterRecipeInternal(recipeID)

	// Find trigger nodes
	triggerNodes := recipe.GetTriggerNodes()
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
func (m *Manager) registerTriggerNode(recipeID string, node rc.Node) (*Instance, error) {
	// Get handler from registry
	handler, err := m.registry.GetHandler(node.Ref)
	if err != nil {
		return nil, fmt.Errorf("unsupported trigger node: %s", node.Ref)
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
		NodeRef:   node.Ref,
		Config:    config,
		Status:    StatusInactive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Register with handler (validation is performed inside Register)
	if err := handler.Register(context.Background(), instance.ID, instance.Config); err != nil {
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
	if handler, err := m.registry.GetHandler(instance.NodeRef); err == nil {
		if err := handler.Unregister(context.Background(), instance.ID); err != nil {
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

	workflowState := runner.NewWorkflowState()
	workflowState.SetExecContext(data)
	err = runner.Run(ctx, recipe, instance.NodeID, workflowState, data, func(rec *rc.Recipe, state runner.NodeExecutionStatus, node rc.Node, output map[string]any, err error) {
		if err != nil {
			log.Printf("Recipe execution error - Recipe: %s, Node: %s, State: %s, Error: %v", rec.Name, node.Id, state, err)
		} else {
			log.Printf("Recipe execution - Recipe: %s, Node: %s, State: %s", rec.Name, node.Id, state)
		}
	})
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

// GetHandler returns a handler by node reference
func (m *Manager) GetHandler(nodeRef string) Handler {
	handler, err := m.registry.GetHandler(nodeRef)
	if err != nil {
		return nil
	}
	return handler
}

// RegisterHandler dynamically registers a handler instance
func (m *Manager) RegisterHandler(handler Handler) error {
	return m.registry.Register(context.Background(), handler, m.eventCh)
}

// UnregisterHandler dynamically unregisters a handler
func (m *Manager) UnregisterHandler(nodeRef string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if there are active triggers of this type
	for _, instances := range m.recipeTriggers {
		for _, instance := range instances {
			if instance.NodeRef == nodeRef {
				return fmt.Errorf("cannot unregister handler '%s': active triggers exist", nodeRef)
			}
		}
	}

	// Remove from registry
	m.registry.UnregisterHandler(nodeRef)

	return nil
}

// GetSupportedNodeRefs returns all supported node references.
func (m *Manager) GetSupportedNodeRefs() []string {
	handlers := m.registry.List()
	nodeRefs := make([]string, len(handlers))
	for i, handler := range handlers {
		nodeRefs[i] = handler.NodeRef
	}
	return nodeRefs
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
				NodeRef:       instance.NodeRef,
				TriggerCount:  instance.TriggerCount,
				LastTriggered: instance.LastTriggered,
				LastError:     instance.LastError,
			})
		}
	}

	return m.store.SaveStats(stats)
}

// Helper functions

// generateTriggerID generates a deterministic trigger ID.
func generateTriggerID(recipeID, nodeID string) string {
	return "trigger_" + recipeID + "_" + nodeID
}
