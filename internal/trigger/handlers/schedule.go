package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/trigger"
	"github.com/robfig/cron/v3"
)

// ScheduleHandler handles cron-based schedule triggers
type ScheduleHandler struct {
	manager   *trigger.Manager
	scheduler *cron.Cron
}

// ScheduleConfig represents schedule configuration
type ScheduleConfig struct {
	Cron string `mapstructure:"cron"`
}

// Validate validates the schedule configuration
func (c *ScheduleConfig) Validate() error {
	if c.Cron == "" {
		return fmt.Errorf("cron expression is required")
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(c.Cron)
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", c.Cron, err)
	}

	return nil
}

func (h *ScheduleHandler) Initialize(manager *trigger.Manager) error {
	h.manager = manager
	h.scheduler = cron.New(cron.WithSeconds())
	return nil
}

func (h *ScheduleHandler) Run(ctx context.Context) error {
	if h.scheduler != nil {
		h.scheduler.Start()
		<-ctx.Done()
		stopCtx := h.scheduler.Stop()
		<-stopCtx.Done()
	}
	return nil
}

func (h *ScheduleHandler) Type() trigger.TriggerType {
	return trigger.TypeSchedule
}

func (h *ScheduleHandler) Validate(configMap map[string]any) error {
	// Use mapstructure directly for validation
	var config ScheduleConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &config,
		TagName: "mapstructure",
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(configMap); err != nil {
		return fmt.Errorf("invalid schedule config: %w", err)
	}

	// Validate using ScheduleConfig's Validate method
	return config.Validate()
}

func (h *ScheduleHandler) Register(instance *trigger.Instance) error {
	if h.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	var config ScheduleConfig
	if err := instance.Bind(&config); err != nil {
		return fmt.Errorf("failed to bind schedule config: %w", err)
	}

	entryID, err := h.scheduler.AddFunc(config.Cron, func() {
		if h.manager != nil {
			data := map[string]any{
				"timestamp": time.Now(),
				"cron":      config.Cron,
				"scheduled": true,
			}
			h.manager.ExecuteTrigger(context.Background(), instance.ID, data)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %w", err)
	}

	instance.SetCleanup(func() error {
		if h.scheduler != nil {
			h.scheduler.Remove(entryID)
		}
		return nil
	})

	instance.Config["entryID"] = int(entryID)

	fmt.Printf("Scheduled cron job: %s (ID: %d)\n", config.Cron, entryID)
	return nil
}

func (h *ScheduleHandler) Unregister(instance *trigger.Instance) error {
	return nil
}
