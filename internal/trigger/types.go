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
