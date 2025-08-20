package trigger

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&WebhookTriggerNode{
		BaseNode: *node.NewBaseNode("teatime.trigger.webhook", node.NodeTypeTrigger, "Webhook", "Webhook을 통해 워크플로우를 실행하는 트리거 노드입니다.", "Webhook"),
	})
}

// Webhook을 통해 워크플로우를 실행하는 트리거 노드
type WebhookTriggerNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *WebhookTriggerNode) Output() []node.NodeProperty {
	output := []node.NodeProperty{
		node.OutputProp(node.Date, "timestamp", "Timestamp",
			node.WithDescription("호출시점의 날짜와 시간입니다."),
		),
		node.OutputProp(node.JSON, "body", "Request Body",
			node.WithDescription("요청 바디입니다."),
		),
	}
	return append(output, r.customParams...)
}

func (r *WebhookTriggerNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		node.StringProp("url", "URL",
			node.WithDescription("Webhook URL을 입력하세요"),
			node.Optional(),
		),
	}
}

func (r *WebhookTriggerNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *WebhookTriggerNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
