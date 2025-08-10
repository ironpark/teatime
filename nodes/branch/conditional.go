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
		{
			Name:        "Branch Taken",
			Description: "선택된 분기 (true/false)",
			Key:         "branchTaken",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Condition Result",
			Description: "조건 평가 결과",
			Key:         "conditionResult",
			Value:       "",
			Type:        node.Bool,
		},
		{
			Name:        "Output Value",
			Description: "분기에서 출력된 값",
			Key:         "outputValue",
			Value:       "",
			Type:        node.JSON,
			Optional:    true,
		},
	}
}

func (r *ConditionalBranchNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Left Value",
			Description: "비교할 왼쪽 값",
			Optional:    false,
			Key:         "leftValue",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Operator",
			Description: "비교 연산자",
			Optional:    false,
			Key:         "operator",
			Value:       "==",
			Type:        node.String,
			Options:     []string{"==", "!=", ">", "<", ">=", "<=", "contains", "startsWith", "endsWith", "regex", "isEmpty", "isNotEmpty"},
		},
		{
			Name:        "Right Value",
			Description: "비교할 오른쪽 값",
			Optional:    false,
			Key:         "rightValue",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Data Type",
			Description: "데이터 타입",
			Optional:    true,
			Key:         "dataType",
			Value:       "string",
			Type:        node.String,
			Options:     []string{"string", "number", "boolean", "json"},
		},
		{
			Name:        "Case Sensitive",
			Description: "대소문자 구분 (문자열 비교 시)",
			Optional:    true,
			Key:         "caseSensitive",
			Value:       "true",
			Type:        node.Bool,
		},
		{
			Name:        "Default Branch",
			Description: "조건이 실패할 때 기본 분기",
			Optional:    true,
			Key:         "defaultBranch",
			Value:       "false",
			Type:        node.String,
			Options:     []string{"true", "false", "error"},
		},
	}
}

func (r *ConditionalBranchNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *ConditionalBranchNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
