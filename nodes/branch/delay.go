package branch

import (
	"context"
	"fmt"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&DelayBranchNode{
		BaseNode: node.NewBaseNode(
			"teatime.branch.delay",
			node.NodeTypeBranch,
			"Delay",
			"지정된 시간만큼 실행을 지연합니다.",
			"Clock",
			[]node.NodeProperty{
				node.IntProp("duration", "Duration",
					node.WithDescription("지연 시간"),
					node.WithDefault(1000),
					node.WithMin(0),
					node.Required(),
				),
				node.SelectProp("unit", "Time Unit", []string{"ms", "s", "m", "h"},
					node.WithDescription("시간 단위 선택"),
					node.RequiredWithDefault("ms"),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.Int64, "actualDelay", "Actual Delay (ms)",
					node.WithDescription("실제 지연된 시간 (밀리초)"),
				),
				node.OutputProp(node.String, "startTime", "Start Time",
					node.WithDescription("지연 시작 시간 (ISO 8601)"),
				),
				node.OutputProp(node.String, "endTime", "End Time",
					node.WithDescription("지연 종료 시간 (ISO 8601)"),
				),
			},
			[]node.OutputHandle{
				{ID: "delayed", Label: "Delayed", Description: "지연 완료 후 실행"},
			},
		),
	})
}

type delayBranchProps struct {
	Duration int    `mapstructure:"duration"`
	Unit     string `mapstructure:"unit"`
}

// DelayBranchNode delays workflow execution for a specified duration.
type DelayBranchNode struct {
	node.BaseNode
}

// Run executes the delay logic.
func (d *DelayBranchNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var props delayBranchProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{},
		}
	}

	// Convert duration to milliseconds based on unit
	var durationMs int64
	switch props.Unit {
	case "ms":
		durationMs = int64(props.Duration)
	case "s":
		durationMs = int64(props.Duration) * 1000
	case "m":
		durationMs = int64(props.Duration) * 1000 * 60
	case "h":
		durationMs = int64(props.Duration) * 1000 * 60 * 60
	default:
		return node.NodeResult{
			Error:         fmt.Errorf("unsupported time unit: %s", props.Unit),
			Continue:      false,
			OutputHandles: []string{},
		}
	}

	// Record start time
	startTime := time.Now()
	startTimeISO := startTime.Format(time.RFC3339)

	// Create duration and sleep
	duration := time.Duration(durationMs) * time.Millisecond

	// Check for context cancellation during delay
	select {
	case <-time.After(duration):
		// Delay completed normally
	case <-ctx.Done():
		// Context was cancelled
		return node.NodeResult{
			Error:         fmt.Errorf("delay interrupted by context cancellation"),
			Continue:      false,
			OutputHandles: []string{},
		}
	}

	// Record end time and calculate actual delay
	endTime := time.Now()
	endTimeISO := endTime.Format(time.RFC3339)
	actualDelayMs := endTime.Sub(startTime).Milliseconds()

	return node.NodeResult{
		Output: map[string]any{
			"actualDelay": actualDelayMs,
			"startTime":   startTimeISO,
			"endTime":     endTimeISO,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"delayed"},
	}
}
