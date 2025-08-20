package trigger

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&CommandTriggerNode{
		BaseNode: node.NewBaseNode(
			"teatime.trigger.command",
			node.NodeTypeTrigger,
			"Command",
			"명령어를 통해 워크플로우를 실행하는 트리거 노드입니다.",
			"Zap",
			[]node.NodeProperty{
				node.StringProp("cmd", "명령어",
					node.WithDescription("명령어를 입력하세요"),
					node.Optional(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "workdir", "Working Directory",
					node.WithDescription("명령어가 호출된 디렉토리입니다."),
				),
				node.OutputProp(node.Date, "timestamp", "Timestamp",
					node.WithDescription("호출시점의 날짜와 시간입니다."),
				),
			},
			nil, // Use default output handle
		),
	})
}

// 명령어를 통해 워크플로우를 실행하는 트리거 노드
type CommandTriggerNode struct {
	node.BaseNode
}
