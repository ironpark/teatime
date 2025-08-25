package branch

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&MatchBranchNode{
		BaseNode: node.NewBaseNode(
			"teatime.branch.match",
			node.NodeTypeBranch,
			"Match",
			"여러 조건을 평가하여 매치된 모든 분기를 실행합니다.",
			"GitMerge",
			[]node.NodeProperty{
				node.IntProp("expressionCount", "Expression Count",
					node.WithDescription("조건 표현식 개수"),
					node.WithDefault(2),
					node.WithMin(1),
					node.WithMax(10),
					node.Required(),
				),
				node.BoolProp("executeAll", "Execute All Matches",
					node.WithDescription("모든 매치된 조건을 실행할지, 첫 번째만 실행할지 선택"),
					node.WithDefault(true),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.StringArray, "matchedBranches", "Matched Branches",
					node.WithDescription("매치된 분기 ID 목록"),
				),
				node.OutputProp(node.Int64, "matchCount", "Match Count",
					node.WithDescription("매치된 조건 개수"),
				),
			},
			[]node.OutputHandle{}, // Dynamic output handles based on expressionCount
		),
	})
}

type matchBranchProps struct {
	ExpressionCount int                    `mapstructure:"expressionCount"`
	ExecuteAll      bool                   `mapstructure:"executeAll"`
	Expressions     map[string]interface{} `mapstructure:",remain"`
}

// MatchBranchNode evaluates multiple conditions and routes to all matching branches.
type MatchBranchNode struct {
	node.BaseNode
}

// GetProperties returns dynamic properties based on expression count
func (m *MatchBranchNode) GetProperties(ctx node.PropertyContext) []node.NodeProperty {
	baseProps := m.BaseNode.GetProperties(ctx)

	// Get expression count from context
	expressionCount := 2 // default
	if count, ok := ctx["expressionCount"].(float64); ok {
		expressionCount = int(count)
	} else if count, ok := ctx["expressionCount"].(int); ok {
		expressionCount = count
	}

	// Add dynamic expression properties
	for i := 1; i <= expressionCount; i++ {
		key := fmt.Sprintf("expression%d", i)
		label := fmt.Sprintf("Expression %d", i)
		desc := fmt.Sprintf("조건 표현식 %d 예: {{node1.output.value == '%s'}}", i, fmt.Sprintf("value%d", i))

		baseProps = append(baseProps, node.BoolProp(key, label,
			node.WithDescription(desc),
			node.Expression(),
		))

		// Add label property for each expression
		labelKey := fmt.Sprintf("label%d", i)
		labelLabel := fmt.Sprintf("Branch %d Label", i)
		labelDesc := fmt.Sprintf("분기 %d의 라벨", i)

		baseProps = append(baseProps, node.StringProp(labelKey, labelLabel,
			node.WithDescription(labelDesc),
			node.WithDefault(fmt.Sprintf("Branch %d", i)),
		))
	}

	return baseProps
}

// GetOutputHandles returns dynamic output handles based on expression count
func (m *MatchBranchNode) GetOutputHandles(ctx node.PropertyContext) []node.OutputHandle {
	baseHandles := m.BaseNode.GetOutputHandles(ctx)

	// Get expression count from context
	expressionCount := 2 // default
	if count, ok := ctx["expressionCount"].(float64); ok {
		expressionCount = int(count)
	} else if count, ok := ctx["expressionCount"].(int); ok {
		expressionCount = count
	}

	// Add dynamic output handles
	for i := 1; i <= expressionCount; i++ {
		handleID := fmt.Sprintf("match%d", i)

		// Get label from context if available
		label := fmt.Sprintf("Branch %d", i)
		labelKey := fmt.Sprintf("label%d", i)
		if labelValue, ok := ctx[labelKey].(string); ok && labelValue != "" {
			label = labelValue
		}

		desc := fmt.Sprintf("조건 %d이 매치될 때 실행", i)

		baseHandles = append(baseHandles, node.OutputHandle{
			ID:          handleID,
			Label:       label,
			Description: desc,
		})
	}

	// Add default handle for no matches
	baseHandles = append(baseHandles, node.OutputHandle{
		ID:          "nomatch",
		Label:       "No Match",
		Description: "어떤 조건도 매치되지 않을 때 실행",
	})

	return baseHandles
}

// Run executes the match branch logic.
func (m *MatchBranchNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract base parameters
	var props matchBranchProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"nomatch"},
		}
	}

	// Evaluate each expression
	var matchedBranches []string
	var outputHandles []string

	for i := 1; i <= props.ExpressionCount; i++ {
		exprKey := fmt.Sprintf("expression%d", i)

		// Get the resolved boolean value
		if exprValue, exists := resolvedProps[exprKey]; exists {
			if result, ok := exprValue.(bool); ok && result {
				branchID := fmt.Sprintf("match%d", i)
				matchedBranches = append(matchedBranches, branchID)
				outputHandles = append(outputHandles, branchID)

				// If executeAll is false, stop at first match
				if !props.ExecuteAll {
					break
				}
			}
		}
	}

	// If no matches found, use nomatch handle
	if len(matchedBranches) == 0 {
		outputHandles = []string{"nomatch"}
		matchedBranches = []string{"nomatch"}
	}

	return node.NodeResult{
		Output: map[string]any{
			"matchedBranches": matchedBranches,
			"matchCount":      len(matchedBranches),
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: outputHandles,
	}
}
