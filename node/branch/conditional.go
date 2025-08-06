package branch

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&ConditionalBranchNode{})
}

// 조건에 따라 워크플로우를 분기하는 브랜치 노드
type ConditionalBranchNode struct {
	customParams []types.NodeProperty
}

func (r *ConditionalBranchNode) Name() string {
	return "Conditional"
}

func (r *ConditionalBranchNode) Description() string {
	return "조건에 따라 워크플로우를 분기하는 브랜치 노드입니다."
}

func (r *ConditionalBranchNode) Type() types.NodeType {
	return types.NodeTypeBranch
}

func (r *ConditionalBranchNode) ID() string {
	return "teatime.branch.conditional"
}

func (r *ConditionalBranchNode) Output() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Branch Taken",
			Description: "선택된 분기 (true/false)",
			Key:         "branchTaken",
			Value:       "",
			Type:        types.String,
		},
		{
			Name:        "Condition Result",
			Description: "조건 평가 결과",
			Key:         "conditionResult",
			Value:       "",
			Type:        types.Bool,
		},
		{
			Name:        "Output Value",
			Description: "분기에서 출력된 값",
			Key:         "outputValue",
			Value:       "",
			Type:        types.JSON,
			Optional:    true,
		},
	}
}

func (r *ConditionalBranchNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Left Value",
			Description: "비교할 왼쪽 값",
			Optional:    false,
			Key:         "leftValue",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Operator",
			Description: "비교 연산자",
			Optional:    false,
			Key:         "operator",
			Value:       "==",
			Type:        types.String,
			Options:     []string{"==", "!=", ">", "<", ">=", "<=", "contains", "startsWith", "endsWith", "regex", "isEmpty", "isNotEmpty"},
		},
		{
			Name:        "Right Value",
			Description: "비교할 오른쪽 값",
			Optional:    false,
			Key:         "rightValue",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Data Type",
			Description: "데이터 타입",
			Optional:    true,
			Key:         "dataType",
			Value:       "string",
			Type:        types.String,
			Options:     []string{"string", "number", "boolean", "json"},
		},
		{
			Name:        "Case Sensitive",
			Description: "대소문자 구분 (문자열 비교 시)",
			Optional:    true,
			Key:         "caseSensitive",
			Value:       "true",
			Type:        types.Bool,
		},
		{
			Name:        "Default Branch",
			Description: "조건이 실패할 때 기본 분기",
			Optional:    true,
			Key:         "defaultBranch",
			Value:       "false",
			Type:        types.String,
			Options:     []string{"true", "false", "error"},
		},
	}
}

func (r *ConditionalBranchNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *ConditionalBranchNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
