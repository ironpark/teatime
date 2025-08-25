// Package trigger provides a flexible trigger system for workflow automation.
// It supports multiple trigger types including webhooks, schedules, commands, and file watching.
package trigger

import (
	"context"
	"time"

	"github.com/ironpark/teatime/internal/recipe"
)

// Status represents the current state of a trigger instance.
type Status string

// Trigger instance states.
const (
	StatusActive   Status = "active"   // Trigger is active and monitoring
	StatusInactive Status = "inactive" // Trigger is disabled
	StatusError    Status = "error"    // Trigger encountered an error
)

// Event represents a trigger event with metadata and data payload.
type Event struct {
	TriggerID   string         `json:"triggerId"`   // Unique trigger identifier
	RecipeID    string         `json:"recipeId"`    // Associated recipe ID
	NodeID      string         `json:"nodeId"`      // Trigger node ID
	Data        map[string]any `json:"data"`        // Event payload data
	TriggeredAt time.Time      `json:"triggeredAt"` // Event timestamp
}

// Handler defines the interface that all trigger handlers must implement.
// Handlers are responsible for managing specific types of triggers.
type Handler interface {
	// NodeRef returns the node reference this handler manages.
	NodeRef() string

	// Name returns the human-readable name of this handler.
	Name() string

	// Description returns a description of this handler's functionality.
	Description() string

	// Register registers a new trigger instance with the handler.
	Register(ctx context.Context, id string, config map[string]any) error

	// Unregister removes a trigger instance from the handler.
	Unregister(ctx context.Context, id string) error

	// Start starts the handler's background processing (if needed).
	Start(ctx context.Context) error

	// Initialize prepares the handler with event channel.
	Initialize(ctx context.Context, eventCh chan<- Event) error
}

// RecipeLoader defines the interface for loading recipe definitions.
type RecipeLoader interface {
	// LoadRecipe loads a recipe by its identifier.
	LoadRecipe(recipeID string) (*recipe.Recipe, error)
}

// RecipeRunner defines the interface for executing recipes.
type RecipeRunner interface {
	// Execute runs a recipe starting from the specified node with given data.
	Execute(ctx context.Context, recipe *recipe.Recipe, startNodeID string, data map[string]any) error
}

// TriggerStat represents execution statistics for a trigger.
type TriggerStat struct {
	TriggerID     string     `json:"triggerId"`               // Trigger identifier
	RecipeID      string     `json:"recipeId"`                // Associated recipe ID
	NodeRef       string     `json:"nodeRef"`                 // Node reference
	TriggerCount  int64      `json:"triggerCount"`            // Number of times triggered
	LastTriggered *time.Time `json:"lastTriggered,omitempty"` // Last trigger timestamp
	LastError     string     `json:"lastError,omitempty"`     // Last error message
}

// Store defines the interface for persisting trigger statistics.
type Store interface {
	// SaveStats persists trigger execution statistics.
	SaveStats(stats []TriggerStat) error
}
