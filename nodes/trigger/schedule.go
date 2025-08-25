package trigger

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
	"github.com/ironpark/teatime/internal/trigger/handlers"
	"github.com/robfig/cron/v3"
)

func init() {
	node.RegisterNode(&ScheduleTriggerNode{
		BaseNode: node.NewBaseNode(
			"teatime.trigger.schedule",
			node.NodeTypeTrigger,
			"Scheduled",
			"지정된 스케줄에 따라 워크플로우를 실행합니다.",
			"Clock",
			[]node.NodeProperty{
				node.StringProp("cron", "Cron Expression",
					node.WithDescription("스케줄을 정의하는 cron 표현식 (예: 0 */5 * * * * - 5분마다)"),
					node.Required(),
				),
				node.StringProp("description", "Description",
					node.WithDescription("스케줄에 대한 설명"),
					node.Optional(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.Date, "timestamp", "Timestamp",
					node.WithDescription("스케줄 실행 시점의 날짜와 시간입니다."),
				),
				node.OutputProp(node.String, "cron", "Cron Expression",
					node.WithDescription("실행된 cron 표현식입니다."),
				),
				node.OutputProp(node.Bool, "scheduled", "Scheduled",
					node.WithDescription("스케줄에 의해 실행되었는지 여부입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Scheduled",
					Description: "Scheduled execution triggered",
				},
			},
		),
	})
}


// ScheduleTriggerNode triggers workflow execution based on cron schedule.
type ScheduleTriggerNode struct {
	node.BaseNode
}

// Run executes the schedule trigger logic.
// This is called when the cron schedule fires.
func (s *ScheduleTriggerNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var config handlers.ScheduleConfig
	if err := mapstructure.Decode(resolvedProps, &config); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"success"},
		}
	}

	// Extract schedule event information from execution context
	var scheduleContext handlers.ScheduleContext
	if err := states.BindExecContext(&scheduleContext); err != nil {
		return node.NodeResult{
			Error:    fmt.Errorf("failed to bind execution context: %w", err),
			Continue: false,
		}
	}

	// Extract individual fields for convenience
	timestamp := scheduleContext.Timestamp
	cronExpr := scheduleContext.Cron
	scheduled := scheduleContext.Scheduled

	// Build output data
	output := map[string]any{
		"timestamp": timestamp,
		"cron":      cronExpr,
		"scheduled": scheduled,
	}

	return node.NodeResult{
		Output:        output,
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// GetCronExpression returns the configured cron expression.
func (s *ScheduleTriggerNode) GetCronExpression() string {
	props := s.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "cron" {
			if cronExpr, ok := prop.Value.(string); ok {
				return cronExpr
			}
		}
	}
	return ""
}

// ValidateCronExpression validates the cron expression format.
func (s *ScheduleTriggerNode) ValidateCronExpression(cronExpr string) error {
	if cronExpr == "" {
		return fmt.Errorf("cron expression is required")
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", cronExpr, err)
	}

	return nil
}

// GetDescription returns the schedule description.
func (s *ScheduleTriggerNode) GetDescription() string {
	props := s.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "description" {
			if desc, ok := prop.Value.(string); ok {
				return desc
			}
		}
	}
	return "Scheduled trigger"
}