package trigger

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&LLMActionNode{})
}

// LLM을 통해 텍스트를 생성하는 액션 노드
type LLMActionNode struct {
	customParams []types.NodeProperty
}

func (r *LLMActionNode) Name() string {
	return "LLM"
}

func (r *LLMActionNode) Description() string {
	return "LLM을 통해 텍스트를 생성하는 액션 노드입니다."
}

func (r *LLMActionNode) Type() types.NodeType {
	return types.NodeTypeAction
}

func (r *LLMActionNode) ID() string {
	return "teatime.action.llm"
}

func (r *LLMActionNode) Output() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Response",
			Description: "LLM의 응답 텍스트입니다.",
			Key:         "response",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Token Usage",
			Description: "토큰 사용량 정보입니다.",
			Key:         "tokenUsage",
			Value:       "",
			Type:        types.JSON,
			Optional:    true,
		},
	}
}

func (r *LLMActionNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Model",
			Description: "사용할 LLM 모델을 선택하세요",
			Optional:    false,
			Key:         "model",
			Value:       "gpt-3.5-turbo",
			Type:        types.String,
			Options:     []string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo", "claude-3-opus", "claude-3-sonnet", "claude-3-haiku"},
		},
		{
			Name:        "API Key",
			Description: "LLM API 키를 입력하세요",
			Optional:    false,
			Key:         "apiKey",
			Value:       "",
			Type:        types.String,
		},
		{
			Name:        "System Prompt",
			Description: "시스템 프롬프트를 입력하세요",
			Optional:    true,
			Key:         "systemPrompt",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "User Prompt",
			Description: "사용자 프롬프트를 입력하세요",
			Optional:    false,
			Key:         "userPrompt",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Temperature",
			Description: "응답의 창의성 정도 (0.0 ~ 2.0)",
			Optional:    true,
			Key:         "temperature",
			Value:       "0.7",
			Type:        types.Float64,
		},
		{
			Name:        "Max Tokens",
			Description: "최대 토큰 수",
			Optional:    true,
			Key:         "maxTokens",
			Value:       "2048",
			Type:        types.Float64,
		},
		{
			Name:        "Top P",
			Description: "Top-p 샘플링 (0.0 ~ 1.0)",
			Optional:    true,
			Key:         "topP",
			Value:       "1.0",
			Type:        types.Float64,
		},
		{
			Name:        "Response Format",
			Description: "응답 형식",
			Optional:    true,
			Key:         "responseFormat",
			Value:       "text",
			Type:        types.String,
			Options:     []string{"text", "json"},
		},
	}
}

func (r *LLMActionNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *LLMActionNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
