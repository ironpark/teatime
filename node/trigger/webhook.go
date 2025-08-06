package trigger

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&WebhookTriggerNode{})
}

// 명령어를 통해 워크플로우를 실행하는 트리거 노드
type WebhookTriggerNode struct {
	customParams []types.NodeProperty
}

func (r *WebhookTriggerNode) Name() string {
	return "Webhook"
}

func (r *WebhookTriggerNode) Description() string {
	return "명령어를 통해 워크플로우를 실행하는 트리거 노드입니다."
}

func (r *WebhookTriggerNode) Type() types.NodeType {
	return types.NodeTypeTrigger
}

func (r *WebhookTriggerNode) ID() string {
	return "teatime.trigger.webhook"
}

func (r *WebhookTriggerNode) Output() []types.NodeProperty {
	output := []types.NodeProperty{
		{
			Name:        "Timestamp",
			Description: "호출시점의 날짜와 시간입니다.",
			Key:         "timestamp",
			Value:       "",
			Type:        types.Date,
		},
		{
			Name:        "Request Body",
			Description: "요청 바디입니다.",
			Key:         "body",
			Value:       "",
			Type:        types.JSON,
		},
	}
	return append(output, r.customParams...)
}

func (r *WebhookTriggerNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "URL",
			Description: "Webhook URL을 입력하세요",
			Optional:    true,
			Key:         "url",
			Value:       "",
			Type:        types.String,
		},
	}
}

func (r *WebhookTriggerNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *WebhookTriggerNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
