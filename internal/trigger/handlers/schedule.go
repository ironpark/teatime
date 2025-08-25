package handlers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ironpark/teatime/internal/trigger"
	"github.com/robfig/cron/v3"
)

// ScheduleHandler handles cron-based schedule triggers
type ScheduleHandler struct {
	scheduler *cron.Cron
	eventCh   chan<- trigger.Event
	entries   map[string]cron.EntryID // triggerID -> entryID mapping
	mu        sync.RWMutex
}

// ScheduleContext represents the schedule event context for schedule triggers.
type ScheduleContext struct {
	Timestamp time.Time `mapstructure:"timestamp"`
	Cron      string    `mapstructure:"cron"`
	Scheduled bool      `mapstructure:"scheduled"`
}

// ScheduleConfig represents schedule configuration
type ScheduleConfig struct {
	Cron        string `mapstructure:"cron"`
	Description string `mapstructure:"description"`
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

func (h *ScheduleHandler) Initialize(ctx context.Context, eventCh chan<- trigger.Event) error {
	h.scheduler = cron.New(cron.WithSeconds())
	h.eventCh = eventCh
	h.entries = make(map[string]cron.EntryID)
	return nil
}

func (h *ScheduleHandler) Start(ctx context.Context) error {
	if h.scheduler != nil {
		h.scheduler.Start()
		<-ctx.Done()
		stopCtx := h.scheduler.Stop()
		<-stopCtx.Done()
	}
	return nil
}

func (h *ScheduleHandler) NodeRef() string {
	return "teatime.trigger.schedule"
}

func (h *ScheduleHandler) Name() string {
	return "Scheduled"
}

func (h *ScheduleHandler) Description() string {
	return "Triggers workflows based on cron schedule"
}

func (h *ScheduleHandler) Register(ctx context.Context, id string, configMap map[string]any) error {
	if h.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	var config ScheduleConfig
	if err := trigger.BindAndValidate(&config, configMap); err != nil {
		return fmt.Errorf("failed to validate schedule config: %w", err)
	}

	entryID, err := h.scheduler.AddFunc(config.Cron, func() {
		if h.eventCh != nil {
			data := map[string]any{
				"timestamp": time.Now(),
				"cron":      config.Cron,
				"scheduled": true,
			}
			event := trigger.Event{
				TriggerID:   id,
				Data:        data,
				TriggeredAt: time.Now(),
			}
			select {
			case h.eventCh <- event:
			default:
				fmt.Printf("Warning: event channel full for trigger %s\n", id)
			}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %w", err)
	}

	// Store entry ID for later removal
	h.mu.Lock()
	h.entries[id] = entryID
	h.mu.Unlock()

	fmt.Printf("Scheduled cron job: %s (ID: %d)\n", config.Cron, entryID)
	return nil
}

func (h *ScheduleHandler) Unregister(ctx context.Context, id string) error {
	h.mu.Lock()
	entryID, exists := h.entries[id]
	if exists {
		delete(h.entries, id)
	}
	h.mu.Unlock()

	if exists && h.scheduler != nil {
		h.scheduler.Remove(entryID)
		fmt.Printf("Removed scheduled job for trigger %s (entry ID: %d)\n", id, entryID)
	}

	return nil
}
