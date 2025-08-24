package trigger

import (
	"time"

	"github.com/go-viper/mapstructure/v2"
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
