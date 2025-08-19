package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&CommandActionNode{
		BaseNode: *node.NewBaseNode("teatime.action.command", node.NodeTypeAction, "Command", "시스템 명령어를 실행하는 액션 노드입니다.", "Terminal"),
	})
}

// 명령어를 실행하는 액션 노드
type CommandActionNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *CommandActionNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Output",
			Description: "명령어 실행 결과입니다.",
			Key:         "output",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Exit Code",
			Description: "명령어 종료 코드입니다.",
			Key:         "exitCode",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Error",
			Description: "오류 메시지입니다.",
			Key:         "error",
			Value:       "",
			Type:        node.String,
			Optional:    true,
		},
	}
}

func (r *CommandActionNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Command",
			Description: "실행할 명령어를 입력하세요",
			Optional:    false,
			Key:         "command",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Working Directory",
			Description: "작업 디렉토리를 입력하세요",
			Optional:    true,
			Key:         "workdir",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Timeout",
			Description: "타임아웃 시간(초)",
			Optional:    true,
			Key:         "timeout",
			Value:       "30",
			Type:        node.Int64,
		},
	}
}

func (r *CommandActionNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *CommandActionNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
