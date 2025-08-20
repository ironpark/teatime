package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&CommandActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.command",
			node.NodeTypeAction,
			"Command",
			"시스템 명령어를 실행하는 액션 노드입니다.",
			"Terminal",
			[]node.NodeProperty{
				node.StringProp("command", "Command",
					node.WithDescription("실행할 명령어를 입력하세요"),
					node.Required(),
				),
				node.StringProp("workdir", "Working Directory",
					node.WithDescription("작업 디렉토리를 입력하세요"),
					node.Optional(),
				),
				node.IntProp("timeout", "Timeout",
					node.WithDescription("타임아웃 시간(초)"),
					node.OptionalWithDefault(int64(30)),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "output", "Output",
					node.WithDescription("명령어 실행 결과입니다."),
				),
				node.OutputProp(node.Int64, "exitCode", "Exit Code",
					node.WithDescription("명령어 종료 코드입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			nil, // Use default output handle
		),
	})
}

// 명령어를 실행하는 액션 노드
type CommandActionNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

