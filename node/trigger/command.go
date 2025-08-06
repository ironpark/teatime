package trigger

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&CommandTriggerNode{})
}

// 명령어를 통해 워크플로우를 실행하는 트리거 노드
type CommandTriggerNode struct {
	customParams []types.NodeProperty
}

func (r *CommandTriggerNode) Name() string {
	return "Command"
}

func (r *CommandTriggerNode) Description() string {
	return "명령어를 통해 워크플로우를 실행하는 트리거 노드입니다."
}

func (r *CommandTriggerNode) Type() types.NodeType {
	return types.NodeTypeTrigger
}

func (r *CommandTriggerNode) ID() string {
	return "teatime.trigger.command"
}

func (r *CommandTriggerNode) Output() []types.NodeProperty {
	output := []types.NodeProperty{
		{
			Name:        "Working Directory",
			Description: "명령어가 호출된 디렉토리입니다.",
			Key:         "workdir",
			Value:       "",
			Type:        types.String,
		},
		{
			Name:        "Timestamp",
			Description: "호출시점의 날짜와 시간입니다.",
			Key:         "timestamp",
			Value:       "",
			Type:        types.Date,
		},
	}
	return append(output, r.customParams...)
}

func (r *CommandTriggerNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "명령어",
			Description: "명령어를 입력하세요",
			Optional:    true,
			Key:         "cmd",
			Value:       "",
			Type:        types.String,
		},
	}
}

func (r *CommandTriggerNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *CommandTriggerNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
