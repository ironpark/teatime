package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/ironpark/teatime/internal/trigger"
)

// ScheduleHandler handles cron-based schedule triggers
type ScheduleHandler struct {
	manager   *trigger.Manager
	scheduler *cron.Cron
}

func (h *ScheduleHandler) Initialize(ctx context.Context, manager *trigger.Manager) error {
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

func (h *ScheduleHandler) Validate(config map[string]interface{}) error {
	cronExpr, ok := config["cron"].(string)
	if !ok || cronExpr == "" {
		return fmt.Errorf("cron expression is required")
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", cronExpr, err)
	}

	return nil
}

func (h *ScheduleHandler) Register(ctx context.Context, instance *trigger.Instance) error {
	if h.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	cronExpr := instance.Config["cron"].(string)
	
	entryID, err := h.scheduler.AddFunc(cronExpr, func() {
		if h.manager != nil {
			data := map[string]interface{}{
				"timestamp": time.Now(),
				"cron":      cronExpr,
				"scheduled": true,
			}
			h.manager.ExecuteTrigger(instance.ID, data)
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

	fmt.Printf("Scheduled cron job: %s (ID: %d)\n", cronExpr, entryID)
	return nil
}

func (h *ScheduleHandler) Unregister(instance *trigger.Instance) error {
	return nil
}

