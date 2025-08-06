package branch

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&SwitchBranchNode{})
}

// 여러 조건을 평가하여 워크플로우를 분기하는 스위치 브랜치 노드
type SwitchBranchNode struct {
	customParams []types.NodeProperty
}

func (r *SwitchBranchNode) Name() string {
	return "Switch"
}

func (r *SwitchBranchNode) Description() string {
	return "여러 조건을 평가하여 워크플로우를 분기하는 스위치 브랜치 노드입니다."
}

func (r *SwitchBranchNode) Type() types.NodeType {
	return types.NodeTypeBranch
}

func (r *SwitchBranchNode) ID() string {
	return "teatime.branch.switch"
}

func (r *SwitchBranchNode) Output() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Selected Case",
			Description: "선택된 케이스 인덱스",
			Key:         "selectedCase",
			Value:       "",
			Type:        types.Int64,
		},
		{
			Name:        "Case Label",
			Description: "선택된 케이스 레이블",
			Key:         "caseLabel",
			Value:       "",
			Type:        types.String,
		},
		{
			Name:        "Match Value",
			Description: "매칭된 값",
			Key:         "matchValue",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Output Value",
			Description: "선택된 케이스의 출력 값",
			Key:         "outputValue",
			Value:       "",
			Type:        types.JSON,
			Optional:    true,
		},
	}
}

func (r *SwitchBranchNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Switch Value",
			Description: "평가할 값",
			Optional:    false,
			Key:         "switchValue",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Cases",
			Description: "케이스 목록 (JSON 배열: [{\"label\": \"case1\", \"value\": \"value1\", \"output\": \"output1\"}, ...])",
			Optional:    false,
			Key:         "cases",
			Value:       "[]",
			Type:        types.JSON,
		},
		{
			Name:        "Match Type",
			Description: "매칭 타입",
			Optional:    true,
			Key:         "matchType",
			Value:       "exact",
			Type:        types.String,
			Options:     []string{"exact", "contains", "regex", "range", "custom"},
		},
		{
			Name:        "Case Sensitive",
			Description: "대소문자 구분",
			Optional:    true,
			Key:         "caseSensitive",
			Value:       "true",
			Type:        types.Bool,
		},
		{
			Name:        "Default Case",
			Description: "기본 케이스 출력 값",
			Optional:    true,
			Key:         "defaultCase",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Stop on First Match",
			Description: "첫 번째 매치에서 중단",
			Optional:    true,
			Key:         "stopOnFirstMatch",
			Value:       "true",
			Type:        types.Bool,
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
	}
}

func (r *SwitchBranchNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *SwitchBranchNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
