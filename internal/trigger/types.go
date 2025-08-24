package trigger

import (
	"context"
	"time"

	"github.com/ironpark/teatime/internal/recipe"
)

type TriggerType string

const (
	TypeWebhook   TriggerType = "webhook"
	TypeSchedule  TriggerType = "schedule"
	TypeCommand   TriggerType = "command"
	TypeFileWatch TriggerType = "filewatch"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusError    Status = "error"
)

// Instance represents an active trigger instance
type Instance struct {
	ID        string                 `json:"id"`
	RecipeID  string                 `json:"recipeId"`
	NodeID    string                 `json:"nodeId"`
	Type      TriggerType            `json:"type"`
	Config    map[string]interface{} `json:"config"`
	Status    Status                 `json:"status"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`

	// Runtime statistics (memory only)
	TriggerCount  int64      `json:"triggerCount"`
	LastTriggered *time.Time `json:"lastTriggered,omitempty"`
	LastError     string     `json:"lastError,omitempty"`

	// Cleanup function (optional)
	cleanup func() error `json:"-"`
}

// SetCleanup sets the cleanup function for this trigger instance
func (i *Instance) SetCleanup(fn func() error) {
	i.cleanup = fn
}

// Event represents a trigger event
type Event struct {
	TriggerID   string                 `json:"triggerId"`
	RecipeID    string                 `json:"recipeId"`
	NodeID      string                 `json:"nodeId"`
	Data        map[string]interface{} `json:"data"`
	TriggeredAt time.Time              `json:"triggeredAt"`
}

// Handler interface for different trigger types
type Handler interface {
	Type() TriggerType
	Register(ctx context.Context, instance *Instance) error
	Unregister(instance *Instance) error
	Validate(config map[string]interface{}) error
}

// External dependencies interfaces
type RecipeLoader interface {
	LoadRecipe(recipeID string) (*recipe.Recipe, error)
}

type RecipeRunner interface {
	Execute(ctx context.Context, recipe *recipe.Recipe, startNodeID string, data map[string]interface{}) error
}

// TriggerStat represents trigger execution statistics
type TriggerStat struct {
	TriggerID     string     `json:"triggerId"`
	RecipeID      string     `json:"recipeId"`
	TriggerCount  int64      `json:"triggerCount"`
	LastTriggered *time.Time `json:"lastTriggered,omitempty"`
	LastError     string     `json:"lastError,omitempty"`
}

// Store interface for optional persistence
type Store interface {
	SaveStats(stats []TriggerStat) error
}