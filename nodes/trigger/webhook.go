package trigger

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&WebhookTriggerNode{
		BaseNode: node.NewBaseNode(
			"teatime.trigger.webhook",
			node.NodeTypeTrigger,
			"Webhook",
			"Webhook을 통해 워크플로우를 실행하는 트리거 노드입니다.",
			"Webhook",
			[]node.NodeProperty{
				node.StringProp("url", "URL",
					node.WithDescription("Webhook URL을 입력하세요"),
					node.Optional(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.Date, "timestamp", "Timestamp",
					node.WithDescription("호출시점의 날짜와 시간입니다."),
				),
				node.OutputProp(node.JSON, "body", "Request Body",
					node.WithDescription("요청 바디입니다."),
				),
			},
			nil, // Use default output handle
		),
	})
}

// Webhook을 통해 워크플로우를 실행하는 트리거 노드
type WebhookTriggerNode struct {
	node.BaseNode
}
