package branch

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LoopBranchNode{
		BaseNode: node.NewBaseNode(
			"teatime.branch.loop",
			node.NodeTypeBranch,
			"Loop",
			"지정된 횟수만큼 워크플로우를 반복 실행합니다.",
			"RotateCcw",
			[]node.NodeProperty{
				node.SelectProp("loopType", "Loop Type", []string{"count", "while", "foreach"},
					node.WithDescription("반복 유형을 선택하세요"),
					node.RequiredWithDefault("count"),
				),
				node.IntProp("maxIterations", "Max Iterations",
					node.WithDescription("최대 반복 횟수 (무한 루프 방지)"),
					node.WithDefault(100),
					node.WithMin(1),
					node.WithMax(10000),
					node.Required(),
				),
				node.IntProp("count", "Count",
					node.WithDescription("반복 횟수 (count 타입일 때)"),
					node.WithDefault(3),
					node.WithMin(1),
					node.WithMax(1000),
				),
				node.BoolProp("condition", "Condition",
					node.WithDescription("반복 조건 (while 타입일 때) 예: {{counter < 10}}"),
					node.Expression(),
				),
				node.JSONProp("array", "Array",
					node.WithDescription("반복할 배열 (foreach 타입일 때)"),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.Int64, "currentIteration", "Current Iteration",
					node.WithDescription("현재 반복 횟수 (1부터 시작)"),
				),
				node.OutputProp(node.Int64, "totalIterations", "Total Iterations",
					node.WithDescription("총 반복 횟수"),
				),
				node.OutputProp(node.JSON, "currentItem", "Current Item",
					node.WithDescription("현재 배열 아이템 (foreach 타입일 때)"),
				),
				node.OutputProp(node.Int64, "currentIndex", "Current Index",
					node.WithDescription("현재 배열 인덱스 (foreach 타입일 때)"),
				),
				node.OutputProp(node.Bool, "isFirst", "Is First",
					node.WithDescription("첫 번째 반복인지 여부"),
				),
				node.OutputProp(node.Bool, "isLast", "Is Last",
					node.WithDescription("마지막 반복인지 여부"),
				),
			},
			[]node.OutputHandle{
				{ID: "loop", Label: "Loop", Description: "각 반복마다 실행"},
				{ID: "complete", Label: "Complete", Description: "모든 반복 완료 후 실행"},
				{ID: "break", Label: "Break", Description: "반복 중단 시 실행"},
			},
		),
	})
}

type loopBranchProps struct {
	LoopType       string      `mapstructure:"loopType"`
	MaxIterations  int         `mapstructure:"maxIterations"`
	Count          int         `mapstructure:"count"`
	Condition      bool        `mapstructure:"condition"`
	Array          any `mapstructure:"array"`
}

// LoopBranchNode executes workflow iterations based on different loop types.
type LoopBranchNode struct {
	node.BaseNode
}

// Run executes the loop logic.
func (l *LoopBranchNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var props loopBranchProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"break"},
		}
	}

	// Get current iteration from workflow state
	currentIteration := int64(1)
	if iter := states.Get("_loop_iteration"); iter != nil {
		if iterInt, ok := iter.(int64); ok {
			currentIteration = iterInt + 1
		}
	}

	var shouldContinueLoop bool
	var totalIterations int64
	var currentItem any
	var currentIndex int64

	// Determine if loop should continue based on type
	switch props.LoopType {
	case "count":
		totalIterations = int64(props.Count)
		shouldContinueLoop = currentIteration <= totalIterations

	case "while":
		totalIterations = int64(props.MaxIterations) // Use max as estimate
		shouldContinueLoop = props.Condition && currentIteration <= int64(props.MaxIterations)

	case "foreach":
		// Handle array iteration
		if arraySlice, ok := props.Array.([]any); ok {
			totalIterations = int64(len(arraySlice))
			currentIndex = currentIteration - 1
			
			if currentIndex < totalIterations {
				currentItem = arraySlice[currentIndex]
				shouldContinueLoop = true
			} else {
				shouldContinueLoop = false
			}
		} else {
			return node.NodeResult{
				Error:         fmt.Errorf("array property is not a valid array for foreach loop"),
				Continue:      false,
				OutputHandles: []string{"break"},
			}
		}

	default:
		return node.NodeResult{
			Error:         fmt.Errorf("unsupported loop type: %s", props.LoopType),
			Continue:      false,
			OutputHandles: []string{"break"},
		}
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return node.NodeResult{
			Error:         fmt.Errorf("loop interrupted by context cancellation"),
			Continue:      false,
			OutputHandles: []string{"break"},
		}
	default:
	}

	// Determine output
	isFirst := currentIteration == 1
	var isLast bool
	var outputHandle string

	if shouldContinueLoop {
		isLast = (currentIteration == totalIterations)
		outputHandle = "loop"
	} else {
		isLast = true
		outputHandle = "complete"
	}

	return node.NodeResult{
		Output: map[string]any{
			"currentIteration": currentIteration,
			"totalIterations":  totalIterations,
			"currentItem":      currentItem,
			"currentIndex":     currentIndex,
			"isFirst":          isFirst,
			"isLast":           isLast,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{outputHandle},
	}
}