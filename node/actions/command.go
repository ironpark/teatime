package trigger

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&CommandActionNode{})
}

// 명령어를 실행하는 액션 노드
type CommandActionNode struct {
	customParams []types.NodeProperty
}

func (r *CommandActionNode) Name() string {
	return "Command"
}

func (r *CommandActionNode) Description() string {
	return "시스템 명령어를 실행하는 액션 노드입니다."
}

func (r *CommandActionNode) Type() types.NodeType {
	return types.NodeTypeAction
}

func (r *CommandActionNode) ID() string {
	return "teatime.action.command"
}

func (r *CommandActionNode) Output() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Output",
			Description: "명령어 실행 결과입니다.",
			Key:         "output",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Exit Code",
			Description: "명령어 종료 코드입니다.",
			Key:         "exitCode",
			Value:       "",
			Type:        types.Int64,
		},
		{
			Name:        "Error",
			Description: "오류 메시지입니다.",
			Key:         "error",
			Value:       "",
			Type:        types.Text,
			Optional:    true,
		},
	}
}

func (r *CommandActionNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Command",
			Description: "실행할 명령어를 입력하세요",
			Optional:    false,
			Key:         "command",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Working Directory",
			Description: "작업 디렉토리를 입력하세요",
			Optional:    true,
			Key:         "workdir",
			Value:       "",
			Type:        types.String,
		},
		{
			Name:        "Environment Variables",
			Description: "환경 변수를 JSON 형식으로 입력하세요",
			Optional:    true,
			Key:         "env",
			Value:       "{}",
			Type:        types.JSON,
		},
		{
			Name:        "Timeout",
			Description: "타임아웃 시간(초)",
			Optional:    true,
			Key:         "timeout",
			Value:       "30",
			Type:        types.Int64,
		},
	}
}

func (r *CommandActionNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *CommandActionNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
