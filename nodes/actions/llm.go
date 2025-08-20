package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LLMActionNode{
		BaseNode: *node.NewBaseNode("teatime.action.llm", node.NodeTypeAction, "LLM", "LLM을 통해 텍스트를 생성하는 액션 노드입니다.", "Bot"),
	})
}

// LLM을 통해 텍스트를 생성하는 액션 노드
type LLMActionNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *LLMActionNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		node.OutputProp(node.String, "response", "Response",
			node.WithDescription("LLM의 응답 텍스트입니다."),
		),
	}
}

func (r *LLMActionNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		node.SelectProp("model", "Model", []string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo", "claude-3-opus", "claude-3-sonnet", "claude-3-haiku"},
			node.WithDescription("사용할 LLM 모델을 선택하세요"),
			node.RequiredWithDefault("gpt-3.5-turbo"),
		),
		node.StringProp("apiKey", "API Key",
			node.WithDescription("LLM API 키를 입력하세요"),
			node.Required(),
		),
		node.StringProp("systemPrompt", "System Prompt",
			node.WithDescription("시스템 프롬프트를 입력하세요"),
			node.WithPlaceholder("You are a helpful assistant..."),
			node.TextArea(3),
			node.Optional(),
		),
		node.StringProp("userPrompt", "User Prompt",
			node.WithDescription("사용자 프롬프트를 입력하세요"),
			node.WithPlaceholder("Enter your prompt here..."),
			node.TextArea(5),
			node.Required(),
		),
		node.FloatProp("temperature", "Temperature",
			node.WithDescription("응답의 창의성 정도 (0.0 ~ 2.0)"),
			node.RangeSlider(0, 2),
			node.WithStep(0.1),
			node.OptionalWithDefault(0.7),
		),
		node.IntProp("maxTokens", "Max Tokens",
			node.WithDescription("최대 토큰 수"),
			node.WithRange(1, 8192, 1),
			node.OptionalWithDefault(int64(2048)),
		),
		node.FloatProp("topP", "Top P",
			node.WithDescription("Top-p 샘플링 (0.0 ~ 1.0)"),
			node.OptionalWithDefault(1.0),
		),
		node.SelectProp("responseFormat", "Response Format", []string{"text", "json"},
			node.WithDescription("응답 형식"),
			node.OptionalWithDefault("text"),
		),
	}
}

func (r *LLMActionNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *LLMActionNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
