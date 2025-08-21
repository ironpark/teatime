package branch

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&ConditionalBranchNode{
		BaseNode: node.NewBaseNode(
			"teatime.branch.conditional",
			node.NodeTypeBranch,
			"Conditional",
			"조건에 따라 워크플로우를 분기하는 브랜치 노드입니다.",
			"GitBranch",
			[]node.NodeProperty{
				node.BoolProp("expression", "Expression",
					node.WithDescription("조건 표현식 예: {{node1.output.value > 10}}"),
					node.Required(),
					node.Expression(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "branchTaken", "Branch Taken",
					node.WithDescription("선택된 분기 (true/false)"),
				),
				node.OutputProp(node.Bool, "conditionResult", "Condition Result",
					node.WithDescription("조건 평가 결과"),
				),
			},
			[]node.OutputHandle{
				{ID: "true", Label: "Then", Description: "조건이 참일 때 실행"},
				{ID: "false", Label: "Else", Description: "조건이 거짓일 때 실행"},
			},
		),
	})
}

type conditionalBranchProps struct {
	Expression bool `mapstructure:"expression"`
}

// ConditionalBranchNode evaluates a condition and routes workflow execution based on the result.
type ConditionalBranchNode struct {
	node.BaseNode
}

// Run executes the conditional branch logic.
func (c *ConditionalBranchNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var props conditionalBranchProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"false"},
		}
	}

	// Use the resolved boolean value directly
	conditionResult := props.Expression

	// Determine output handle and branch taken
	var outputHandle string
	var branchTaken string
	if conditionResult {
		outputHandle = "true"
		branchTaken = "true"
	} else {
		outputHandle = "false"
		branchTaken = "false"
	}

	return node.NodeResult{
		Output: map[string]any{
			"branchTaken":     branchTaken,
			"conditionResult": conditionResult,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{outputHandle},
	}
}
