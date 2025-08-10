package trigger

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&CommandTriggerNode{
		BaseNode: *node.NewBaseNode("teatime.trigger.command", node.NodeTypeTrigger, "Command", "명령어를 통해 워크플로우를 실행하는 트리거 노드입니다.", "Zap"),
	})
}

// 명령어를 통해 워크플로우를 실행하는 트리거 노드
type CommandTriggerNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *CommandTriggerNode) Output() []node.NodeProperty {
	output := []node.NodeProperty{
		{
			Name:        "Working Directory",
			Description: "명령어가 호출된 디렉토리입니다.",
			Key:         "workdir",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Timestamp",
			Description: "호출시점의 날짜와 시간입니다.",
			Key:         "timestamp",
			Value:       "",
			Type:        node.Date,
		},
	}
	return append(output, r.customParams...)
}

func (r *CommandTriggerNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "명령어",
			Description: "명령어를 입력하세요",
			Optional:    true,
			Key:         "cmd",
			Value:       "",
			Type:        node.String,
		},
	}
}

func (r *CommandTriggerNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *CommandTriggerNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
