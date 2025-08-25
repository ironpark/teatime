package services

import (
	"context"
	"fmt"
	"log"
	"maps"

	"github.com/ironpark/teatime/internal/recipe"
	"github.com/ironpark/teatime/internal/runner"
	"github.com/ironpark/teatime/internal/trigger"
	"github.com/ironpark/teatime/internal/trigger/handlers"
	"github.com/ironpark/teatime/stores"
)

// TriggersService manages trigger registration and execution
type TriggersService struct {
	store          *stores.Store
	recipesService *RecipesService
	triggerManager *trigger.Manager
}

// RecipeLoader implementation for trigger manager
type recipeLoader struct {
	store *stores.Store
}

func (r *recipeLoader) LoadRecipe(recipeID string) (*recipe.Recipe, error) {
	return r.store.GetRecipe(recipeID)
}

// RecipeRunner implementation for trigger manager
type recipeRunner struct {
	recipesService *RecipesService
}

func (r *recipeRunner) Execute(ctx context.Context, rec *recipe.Recipe, startNodeID string, data map[string]any) error {
	// Convert trigger data to properties map
	properties := make(map[string]any, len(data))
	maps.Copy(properties, data)

	// Execute the recipe starting from the trigger node
	return runner.Run(ctx, rec, startNodeID, properties, func(rec *recipe.Recipe, state runner.NodeExecutionStatus, node recipe.Node, output map[string]any, err error) {
		if err != nil {
			log.Printf("Recipe execution error - Recipe: %s, Node: %s, State: %s, Error: %v", rec.Name, node.Id, state, err)
		} else {
			log.Printf("Recipe execution - Recipe: %s, Node: %s, State: %s", rec.Name, node.Id, state)
		}
	})
}

// NewTriggersService creates a new triggers service with internal registry management.
func NewTriggersService(store *stores.Store, recipesService *RecipesService) *TriggersService {
	loader := &recipeLoader{store: store}
	runner := &recipeRunner{recipesService: recipesService}

	// Create manager with internal registry
	manager := trigger.NewManager(loader, runner)

	// Register handler instances directly with the manager
	handlerInstances := []trigger.Handler{
		&handlers.WebhookHandler{},
		&handlers.ScheduleHandler{},
		&handlers.CommandHandler{},
		&handlers.FileWatchHandler{},
	}

	for _, handler := range handlerInstances {
		if err := manager.RegisterHandler(handler); err != nil {
			log.Printf("Failed to register handler %s: %v", handler.NodeRef(), err)
		}
	}

	return &TriggersService{
		store:          store,
		recipesService: recipesService,
		triggerManager: manager,
	}
}

// Start starts the triggers service and registers all existing recipe triggers.
func (s *TriggersService) Start(ctx context.Context) error {
	if err := s.triggerManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start trigger manager: %w", err)
	}

	// Register triggers for all existing recipes
	if err := s.registerAllRecipes(); err != nil {
		// Log but don't fail - service can start even if some triggers fail to register
		log.Printf("Trigger registration warning: %v", err)
	}

	log.Println("Triggers service started successfully")
	return nil
}

// RegisterRecipe registers triggers for a recipe
func (s *TriggersService) RegisterRecipe(recipeID string) error {
	recipe, err := s.store.GetRecipe(recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	return s.triggerManager.RegisterRecipe(recipe)
}

// UnregisterRecipe unregisters triggers for a recipe
func (s *TriggersService) UnregisterRecipe(recipeID string) error {
	recipe, err := s.store.GetRecipe(recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	return s.triggerManager.UnregisterRecipe(recipe.Path)
}

// GetActiveTriggers returns all active triggers
func (s *TriggersService) GetActiveTriggers() []*trigger.Instance {
	return s.triggerManager.GetActiveTriggers()
}

// GetTriggersByRecipe returns triggers for a specific recipe
func (s *TriggersService) GetTriggersByRecipe(recipeID string) ([]*trigger.Instance, error) {
	recipe, err := s.store.GetRecipe(recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	return s.triggerManager.GetTriggersByRecipe(recipe.Path), nil
}

// ExecuteCommand executes a command trigger (for command-based triggers)
func (s *TriggersService) ExecuteCommand(ctx context.Context, command string, args map[string]any) error {
	if cmdHandler, ok := s.triggerManager.GetHandler("teatime.trigger.command").(*handlers.CommandHandler); ok {
		return cmdHandler.ExecuteCommand(ctx, command, args)
	}
	return fmt.Errorf("command handler not found")
}

// GetRegisteredCommands returns all registered commands
func (s *TriggersService) GetRegisteredCommands() []string {
	if cmdHandler, ok := s.triggerManager.GetHandler("teatime.trigger.command").(*handlers.CommandHandler); ok {
		return cmdHandler.GetRegisteredCommands()
	}
	return []string{}
}

// GetSupportedNodeRefs returns all supported node references from the registry
func (s *TriggersService) GetSupportedNodeRefs() []string {
	return s.triggerManager.GetSupportedNodeRefs()
}

// RegisterHandler dynamically registers a handler instance
func (s *TriggersService) RegisterHandler(handler trigger.Handler) error {
	return s.triggerManager.RegisterHandler(handler)
}

// UnregisterHandler dynamically unregisters a trigger handler
func (s *TriggersService) UnregisterHandler(nodeRef string) error {
	return s.triggerManager.UnregisterHandler(nodeRef)
}

// RefreshRecipeTriggers refreshes triggers for a specific recipe (call after recipe updates).
func (s *TriggersService) RefreshRecipeTriggers(recipeID string) error {
	recipe, err := s.store.GetRecipe(recipeID)
	if err != nil {
		return fmt.Errorf("failed to load recipe %s: %w", recipeID, err)
	}

	// First unregister existing triggers
	if err := s.triggerManager.UnregisterRecipe(recipe.Path); err != nil {
		log.Printf("Warning: failed to unregister triggers for recipe %s: %v", recipeID, err)
	}

	// Re-register triggers
	if err := s.triggerManager.RegisterRecipe(recipe); err != nil {
		return fmt.Errorf("failed to register triggers for recipe %s: %w", recipeID, err)
	}

	log.Printf("Successfully refreshed triggers for recipe: %s", recipeID)
	return nil
}

// RefreshAllTriggers refreshes all triggers (useful for bulk updates)
func (s *TriggersService) RefreshAllTriggers() error {
	return s.registerAllRecipes()
}

// GetTriggerStats returns trigger statistics
func (s *TriggersService) GetTriggerStats() map[string]any {
	triggers := s.triggerManager.GetActiveTriggers()
	stats := make(map[string]any)

	stats["total_triggers"] = len(triggers)

	typeCount := make(map[string]int)
	totalExecutions := int64(0)

	for _, t := range triggers {
		typeCount[t.NodeRef]++
		totalExecutions += t.TriggerCount
	}

	stats["by_type"] = typeCount
	stats["total_executions"] = totalExecutions

	return stats
}

// registerAllRecipes registers triggers for all existing recipes.
func (s *TriggersService) registerAllRecipes() error {
	recipes, err := s.store.ListRecipes()
	if err != nil {
		return fmt.Errorf("failed to list recipes: %w", err)
	}

	var errors []error
	successCount := 0

	for _, dbRecipe := range recipes {
		recipe, err := s.store.GetRecipe(dbRecipe.ID)
		if err != nil {
			errors = append(errors, fmt.Errorf("recipe %s: failed to load: %w", dbRecipe.ID, err))
			continue
		}

		if err := s.triggerManager.RegisterRecipe(recipe); err != nil {
			errors = append(errors, fmt.Errorf("recipe %s: failed to register triggers: %w", dbRecipe.ID, err))
			continue
		}

		successCount++
	}

	// Log results
	log.Printf("Registered triggers for %d/%d recipes", successCount, len(recipes))

	if len(errors) > 0 {
		log.Printf("Registration errors: %v", errors)
		if successCount == 0 {
			return fmt.Errorf("failed to register any triggers: %d errors occurred", len(errors))
		}
		// Return warning if some succeeded
		return fmt.Errorf("partial registration failure: %d/%d failed", len(errors), len(recipes))
	}

	return nil
}
