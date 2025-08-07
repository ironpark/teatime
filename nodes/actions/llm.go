package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LLMActionNode{
		BaseNode: *node.NewBaseNode("teatime.action.llm", node.NodeTypeAction, "LLM", "LLM을 통해 텍스트를 생성하는 액션 노드입니다."),
	})
}

// LLM을 통해 텍스트를 생성하는 액션 노드
type LLMActionNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *LLMActionNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Response",
			Description: "LLM의 응답 텍스트입니다.",
			Key:         "response",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Token Usage",
			Description: "토큰 사용량 정보입니다.",
			Key:         "tokenUsage",
			Value:       "",
			Type:        node.JSON,
			Optional:    true,
		},
	}
}

func (r *LLMActionNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Model",
			Description: "사용할 LLM 모델을 선택하세요",
			Optional:    false,
			Key:         "model",
			Value:       "gpt-3.5-turbo",
			Type:        node.String,
			Options:     []string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo", "claude-3-opus", "claude-3-sonnet", "claude-3-haiku"},
		},
		{
			Name:        "API Key",
			Description: "LLM API 키를 입력하세요",
			Optional:    false,
			Key:         "apiKey",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "System Prompt",
			Description: "시스템 프롬프트를 입력하세요",
			Optional:    true,
			Key:         "systemPrompt",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "User Prompt",
			Description: "사용자 프롬프트를 입력하세요",
			Optional:    false,
			Key:         "userPrompt",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Temperature",
			Description: "응답의 창의성 정도 (0.0 ~ 2.0)",
			Optional:    true,
			Key:         "temperature",
			Value:       "0.7",
			Type:        node.Float64,
		},
		{
			Name:        "Max Tokens",
			Description: "최대 토큰 수",
			Optional:    true,
			Key:         "maxTokens",
			Value:       "2048",
			Type:        node.Float64,
		},
		{
			Name:        "Top P",
			Description: "Top-p 샘플링 (0.0 ~ 1.0)",
			Optional:    true,
			Key:         "topP",
			Value:       "1.0",
			Type:        node.Float64,
		},
		{
			Name:        "Response Format",
			Description: "응답 형식",
			Optional:    true,
			Key:         "responseFormat",
			Value:       "text",
			Type:        node.String,
			Options:     []string{"text", "json"},
		},
	}
}

func (r *LLMActionNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *LLMActionNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
