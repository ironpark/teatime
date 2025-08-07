package trigger

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&WebhookTriggerNode{
		BaseNode: *node.NewBaseNode("teatime.trigger.webhook", node.NodeTypeTrigger, "Webhook", "Webhook을 통해 워크플로우를 실행하는 트리거 노드입니다."),
	})
}

// Webhook을 통해 워크플로우를 실행하는 트리거 노드
type WebhookTriggerNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *WebhookTriggerNode) Output() []node.NodeProperty {
	output := []node.NodeProperty{
		{
			Name:        "Timestamp",
			Description: "호출시점의 날짜와 시간입니다.",
			Key:         "timestamp",
			Value:       "",
			Type:        node.Date,
		},
		{
			Name:        "Request Body",
			Description: "요청 바디입니다.",
			Key:         "body",
			Value:       "",
			Type:        node.JSON,
		},
	}
	return append(output, r.customParams...)
}

func (r *WebhookTriggerNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "URL",
			Description: "Webhook URL을 입력하세요",
			Optional:    true,
			Key:         "url",
			Value:       "",
			Type:        node.String,
		},
	}
}

func (r *WebhookTriggerNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *WebhookTriggerNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
