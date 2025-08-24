package services

import (
	"context"
	"fmt"
	"log"

	"github.com/ironpark/teatime/internal/recipe"
	"github.com/ironpark/teatime/internal/runner"
	"github.com/ironpark/teatime/internal/trigger"
	"github.com/ironpark/teatime/internal/trigger/handlers"
	"github.com/ironpark/teatime/stores"
)

// TriggersService manages trigger registration and execution
type TriggersService struct {
	store           *stores.Store
	recipesService  *RecipesService
	triggerManager  *trigger.Manager
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

func (r *recipeRunner) Execute(ctx context.Context, rec *recipe.Recipe, startNodeID string, data map[string]interface{}) error {
	// Convert trigger data to properties map
	properties := make(map[string]any)
	for key, value := range data {
		properties[key] = value
	}
	
	// Execute the recipe starting from the trigger node
	return runner.Run(ctx, rec, startNodeID, properties, func(rec *recipe.Recipe, state runner.NodeExecutionStatus, node recipe.Node, output map[string]any, err error) {
		if err != nil {
			log.Printf("Recipe execution error - Recipe: %s, Node: %s, State: %s, Error: %v", rec.Name, node.Id, state, err)
		} else {
			log.Printf("Recipe execution - Recipe: %s, Node: %s, State: %s", rec.Name, node.Id, state)
		}
	})
}

// NewTriggersService creates a new triggers service
func NewTriggersService(store *stores.Store, recipesService *RecipesService) *TriggersService {
	loader := &recipeLoader{store: store}
	runner := &recipeRunner{recipesService: recipesService}
	
	// Create handlers slice
	handlersList := []trigger.Handler{
		&handlers.WebhookHandler{},
		&handlers.ScheduleHandler{},
		&handlers.CommandHandler{},
		&handlers.FileWatchHandler{},
	}
	
	manager := trigger.NewManager(loader, runner, handlersList)
	
	return &TriggersService{
		store:          store,
		recipesService: recipesService,
		triggerManager: manager,
	}
}

// Start starts the triggers service
func (s *TriggersService) Start() error {
	if err := s.triggerManager.Start(); err != nil {
		return fmt.Errorf("failed to start trigger manager: %w", err)
	}
	
	// Register triggers for all existing recipes
	if err := s.registerAllRecipes(); err != nil {
		log.Printf("Warning: failed to register some triggers: %v", err)
	}
	
	log.Println("Triggers service started")
	return nil
}

// Shutdown gracefully shuts down the triggers service
func (s *TriggersService) Shutdown() error {
	return s.triggerManager.Shutdown()
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
func (s *TriggersService) ExecuteCommand(command string, args map[string]interface{}) error {
	if cmdHandler, ok := s.triggerManager.GetHandler(trigger.TypeCommand).(*handlers.CommandHandler); ok {
		return cmdHandler.ExecuteCommand(command, args)
	}
	return fmt.Errorf("command handler not found")
}

// GetRegisteredCommands returns all registered commands
func (s *TriggersService) GetRegisteredCommands() []string {
	if cmdHandler, ok := s.triggerManager.GetHandler(trigger.TypeCommand).(*handlers.CommandHandler); ok {
		return cmdHandler.GetRegisteredCommands()
	}
	return []string{}
}

// RefreshRecipeTriggers refreshes triggers for a specific recipe (call after recipe updates)
func (s *TriggersService) RefreshRecipeTriggers(recipeID string) error {
	// First unregister existing triggers for this recipe
	recipe, err := s.store.GetRecipe(recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}
	
	if err := s.triggerManager.UnregisterRecipe(recipe.Path); err != nil {
		log.Printf("Warning: failed to unregister triggers for recipe %s: %v", recipeID, err)
	}
	
	// Re-register triggers
	return s.triggerManager.RegisterRecipe(recipe)
}

// RefreshAllTriggers refreshes all triggers (useful for bulk updates)
func (s *TriggersService) RefreshAllTriggers() error {
	return s.registerAllRecipes()
}

// GetTriggerStats returns trigger statistics
func (s *TriggersService) GetTriggerStats() map[string]interface{} {
	triggers := s.triggerManager.GetActiveTriggers()
	stats := make(map[string]interface{})
	
	stats["total_triggers"] = len(triggers)
	
	typeCount := make(map[string]int)
	totalExecutions := int64(0)
	
	for _, t := range triggers {
		typeCount[string(t.Type)]++
		totalExecutions += t.TriggerCount
	}
	
	stats["by_type"] = typeCount
	stats["total_executions"] = totalExecutions
	
	return stats
}

// registerAllRecipes registers triggers for all existing recipes
func (s *TriggersService) registerAllRecipes() error {
	recipes, err := s.store.ListRecipes()
	if err != nil {
		return fmt.Errorf("failed to list recipes: %w", err)
	}
	
	var errors []error
	for _, dbRecipe := range recipes {
		recipe, err := s.store.GetRecipe(dbRecipe.ID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get recipe %s: %w", dbRecipe.ID, err))
			continue
		}
		
		if err := s.triggerManager.RegisterRecipe(recipe); err != nil {
			errors = append(errors, fmt.Errorf("failed to register triggers for recipe %s: %w", dbRecipe.ID, err))
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("multiple registration errors occurred: %v", errors)
	}
	
	return nil
}