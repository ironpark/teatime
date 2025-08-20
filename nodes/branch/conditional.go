package branch

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&ConditionalBranchNode{
		BaseNode: *node.NewBaseNode("teatime.branch.conditional", node.NodeTypeBranch, "Conditional", "조건에 따라 워크플로우를 분기하는 브랜치 노드입니다.", "GitBranch"),
	})
}

// 조건에 따라 워크플로우를 분기하는 브랜치 노드
type ConditionalBranchNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *ConditionalBranchNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		node.OutputProp(node.String, "branchTaken", "Branch Taken",
			node.WithDescription("선택된 분기 (true/false)"),
		),
		node.OutputProp(node.Bool, "conditionResult", "Condition Result",
			node.WithDescription("조건 평가 결과"),
		),
		node.OutputProp(node.JSON, "outputValue", "Output Value",
			node.WithDescription("분기에서 출력된 값"),
		),
	}
}

func (r *ConditionalBranchNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		node.StringProp("leftValue", "Left Value",
			node.WithDescription("비교할 왼쪽 값"),
			node.Required(),
		),
		node.SelectProp("operator", "Operator", []string{"==", "!=", ">", "<", ">=", "<=", "contains", "startsWith", "endsWith", "regex", "isEmpty", "isNotEmpty"},
			node.WithDescription("비교 연산자"),
			node.RequiredWithDefault("=="),
		),
		node.StringProp("rightValue", "Right Value",
			node.WithDescription("비교할 오른쪽 값"),
			node.Required(),
		),
		node.SelectProp("dataType", "Data Type", []string{"string", "number", "boolean", "json"},
			node.WithDescription("데이터 타입"),
			node.OptionalWithDefault("string"),
		),
		node.BoolProp("caseSensitive", "Case Sensitive",
			node.WithDescription("대소문자 구분 (문자열 비교 시)"),
			node.OptionalWithDefault(true),
		),
		node.SelectProp("defaultBranch", "Default Branch", []string{"true", "false", "error"},
			node.WithDescription("조건이 실패할 때 기본 분기"),
			node.OptionalWithDefault("false"),
		),
	}
}

func (r *ConditionalBranchNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *ConditionalBranchNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
