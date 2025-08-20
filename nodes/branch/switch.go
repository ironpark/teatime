package branch

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&SwitchBranchNode{
		BaseNode: *node.NewBaseNode("teatime.branch.switch", node.NodeTypeBranch, "Switch", "여러 조건을 평가하여 워크플로우를 분기하는 스위치 브랜치 노드입니다.", "Shuffle"),
	})
}

// 여러 조건을 평가하여 워크플로우를 분기하는 스위치 브랜치 노드
type SwitchBranchNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *SwitchBranchNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		node.OutputProp(node.Int64, "selectedCase", "Selected Case",
			node.WithDescription("선택된 케이스 인덱스"),
		),
		node.OutputProp(node.String, "caseLabel", "Case Label",
			node.WithDescription("선택된 케이스 레이블"),
		),
		node.OutputProp(node.String, "matchValue", "Match Value",
			node.WithDescription("매칭된 값"),
		),
		node.OutputProp(node.JSON, "outputValue", "Output Value",
			node.WithDescription("선택된 케이스의 출력 값"),
		),
	}
}

func (r *SwitchBranchNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		node.StringProp("switchValue", "Switch Value",
			node.WithDescription("평가할 값"),
			node.Required(),
		),
		node.JSONProp("cases", "Cases",
			node.WithDescription("케이스 목록 (JSON 배열: [{\"label\": \"case1\", \"value\": \"value1\", \"output\": \"output1\"}, ...])"),
			node.RequiredWithDefault("[]"),
		),
		node.SelectProp("matchType", "Match Type", []string{"exact", "contains", "regex", "range", "custom"},
			node.WithDescription("매칭 타입"),
			node.OptionalWithDefault("exact"),
		),
		node.BoolProp("caseSensitive", "Case Sensitive",
			node.WithDescription("대소문자 구분"),
			node.OptionalWithDefault(true),
		),
		node.StringProp("defaultCase", "Default Case",
			node.WithDescription("기본 케이스 출력 값"),
			node.Optional(),
		),
		node.BoolProp("stopOnFirstMatch", "Stop on First Match",
			node.WithDescription("첫 번째 매치에서 중단"),
			node.OptionalWithDefault(true),
		),
		node.SelectProp("dataType", "Data Type", []string{"string", "number", "boolean", "json"},
			node.WithDescription("데이터 타입"),
			node.OptionalWithDefault("string"),
		),
	}
}

func (r *SwitchBranchNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *SwitchBranchNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
