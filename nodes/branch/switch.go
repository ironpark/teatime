package branch

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&SwitchBranchNode{
		BaseNode: *node.NewBaseNode("teatime.branch.switch", node.NodeTypeBranch, "Switch", "여러 조건을 평가하여 워크플로우를 분기하는 스위치 브랜치 노드입니다."),
	})
}

// 여러 조건을 평가하여 워크플로우를 분기하는 스위치 브랜치 노드
type SwitchBranchNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *SwitchBranchNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Selected Case",
			Description: "선택된 케이스 인덱스",
			Key:         "selectedCase",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Case Label",
			Description: "선택된 케이스 레이블",
			Key:         "caseLabel",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Match Value",
			Description: "매칭된 값",
			Key:         "matchValue",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Output Value",
			Description: "선택된 케이스의 출력 값",
			Key:         "outputValue",
			Value:       "",
			Type:        node.JSON,
			Optional:    true,
		},
	}
}

func (r *SwitchBranchNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Switch Value",
			Description: "평가할 값",
			Optional:    false,
			Key:         "switchValue",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Cases",
			Description: "케이스 목록 (JSON 배열: [{\"label\": \"case1\", \"value\": \"value1\", \"output\": \"output1\"}, ...])",
			Optional:    false,
			Key:         "cases",
			Value:       "[]",
			Type:        node.JSON,
		},
		{
			Name:        "Match Type",
			Description: "매칭 타입",
			Optional:    true,
			Key:         "matchType",
			Value:       "exact",
			Type:        node.String,
			Options:     []string{"exact", "contains", "regex", "range", "custom"},
		},
		{
			Name:        "Case Sensitive",
			Description: "대소문자 구분",
			Optional:    true,
			Key:         "caseSensitive",
			Value:       "true",
			Type:        node.Bool,
		},
		{
			Name:        "Default Case",
			Description: "기본 케이스 출력 값",
			Optional:    true,
			Key:         "defaultCase",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Stop on First Match",
			Description: "첫 번째 매치에서 중단",
			Optional:    true,
			Key:         "stopOnFirstMatch",
			Value:       "true",
			Type:        node.Bool,
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
	}
}

func (r *SwitchBranchNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *SwitchBranchNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
