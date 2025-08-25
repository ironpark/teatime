package trigger

import (
	"time"

	"github.com/go-viper/mapstructure/v2"
)

// Instance represents an active trigger instance with configuration and runtime data.
type Instance struct {
	RecipeRunner `json:"-"` // Embedded recipe runner for execution

	ID        string                 `json:"id"`        // Unique instance identifier
	RecipeID  string                 `json:"recipeId"`  // Associated recipe ID  
	NodeID    string                 `json:"nodeId"`    // Source node ID in recipe
	NodeRef   string                 `json:"nodeRef"`   // Node reference
	Config    map[string]any `json:"config"`    // Trigger configuration
	Status    Status                 `json:"status"`    // Current status
	CreatedAt time.Time              `json:"createdAt"` // Creation timestamp
	UpdatedAt time.Time              `json:"updatedAt"` // Last update timestamp

	// Runtime statistics (memory only)
	TriggerCount  int64      `json:"triggerCount"`              // Total trigger executions
	LastTriggered *time.Time `json:"lastTriggered,omitempty"`   // Last execution time
	LastError     string     `json:"lastError,omitempty"`       // Last error message
}

// SetRunner sets the recipe runner for this instance.
func (i *Instance) SetRunner(runner RecipeRunner) {
	i.RecipeRunner = runner
}

// Bind decodes the trigger configuration into the provided struct using mapstructure tags.
// This is a convenience method for handlers to parse their specific configuration format.
func (i *Instance) Bind(v any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  v,
		TagName: "mapstructure",
	})
	if err != nil {
		return err
	}
	return decoder.Decode(i.Config)
}
