package actions

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicOption "github.com/anthropics/anthropic-sdk-go/option"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var anthropicModels = []string{
	"claude-3-7-sonnet-latest",
	"claude-3-7-sonnet-20250219",
	"claude-3-5-haiku-latest",
	"claude-3-5-haiku-20241022",
	"claude-sonnet-4-20250514",
	"claude-sonnet-4-0",
	"claude-4-sonnet-20250514",
	"claude-3-5-sonnet-latest",
	"claude-3-5-sonnet-20241022",
	"claude-3-5-sonnet-20240620",
	"claude-opus-4-0",
	"claude-opus-4-20250514",
	"claude-4-opus-20250514",
	"claude-opus-4-1-20250805",
}

var openaiModels = []string{
	"gpt-5", "gpt-5-mini", "gpt-5-nano", "gpt-5-2025-08-07", "gpt-5-mini-2025-08-07", "gpt-5-nano-2025-08-07", "gpt-5-chat-latest", "gpt-4.1", "gpt-4.1-mini", "gpt-4.1-nano", "o4-mini", "o3", "o3-mini", "o1", "o1-mini", "gpt-4o", "chatgpt-4o-latest", "gpt-4o-mini", "gpt-4-turbo", "gpt-4",
}

var (
	providers       = []string{"openai", "anthropic", "openrouter", "custom"}
	providerToModel = map[string][]string{
		"openai":     openaiModels,
		"anthropic":  anthropicModels,
		"openrouter": {"openai/gpt-4", "openai/gpt-3.5-turbo", "anthropic/claude-3-opus", "anthropic/claude-3-sonnet", "anthropic/claude-3-haiku", "meta-llama/llama-2-70b-chat"},
		"custom":     {},
	}
	providerDefaultModel = map[string]string{
		"openai":     "gpt-4o-mini",
		"anthropic":  "claude-3-haiku-20240307",
		"openrouter": "openai/gpt-4o-mini",
		"custom":     "",
	}
)

func init() {
	node.RegisterNode(&LLMActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.llm",
			node.NodeTypeAction,
			"LLM",
			"LLM을 통해 텍스트를 생성합니다.",
			"Bot",
			[]node.NodeProperty{
				node.SelectProp("provider", "Provider", providers,
					node.WithDescription("LLM 제공업체를 선택하세요"),
					node.RequiredWithDefault("openai"),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "response", "Response",
					node.WithDescription("LLM의 응답 텍스트입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "LLM response generated successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "LLM request failed",
				},
			},
		),
	})
}

// LLMActionNode generates text using Large Language Models via multiple providers.
type LLMActionNode struct {
	node.BaseNode
}

// GetProperties returns dynamic properties based on the selected provider.
func (l *LLMActionNode) GetProperties(ctx node.PropertyContext) []node.NodeProperty {
	// Get base properties
	baseProps := l.BaseNode.GetProperties(ctx)

	// Get current provider selection
	provider, _ := ctx["provider"].(string)
	if provider == "" {
		provider = "openai"
	}

	// Start with base properties (just provider)
	props := make([]node.NodeProperty, 0, len(baseProps)+8)
	props = append(props, baseProps...)

	// Add model property based on provider
	props = append(props, node.SelectProp("model", "Model", providerToModel[provider],
		node.WithDescription("모델을 선택하세요"),
		node.RequiredWithDefault(providerDefaultModel[provider]),
	))

	// Add API Key property
	var apiKeyProp node.NodeProperty
	switch provider {
	case "custom":
		apiKeyProp = node.StringProp("apiKey", "API Key",
			node.WithDescription("커스텀 API 키를 입력하세요"),
			node.Optional(),
		)
	default:
		apiKeyProp = node.StringProp("apiKey", "API Key",
			node.WithDescription("API 키를 입력하세요"),
			node.Required(),
		)
	}
	props = append(props, apiKeyProp)

	// Add Base URL property (only for custom provider)
	if provider == "custom" {
		baseURLProp := node.StringProp("baseURL", "Base URL",
			node.WithDescription("커스텀 API 기본 URL을 입력하세요"),
			node.WithPlaceholder("https://your-api.com/v1"),
			node.Required(),
		)
		props = append(props, baseURLProp)
	}

	// Add remaining properties
	props = append(props,
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
	)

	return props
}

type llmActionProps struct {
	Provider       string  `mapstructure:"provider"`
	Model          string  `mapstructure:"model"`
	APIKey         string  `mapstructure:"apiKey"`
	BaseURL        string  `mapstructure:"baseURL"`
	SystemPrompt   string  `mapstructure:"systemPrompt"`
	UserPrompt     string  `mapstructure:"userPrompt"`
	Temperature    float64 `mapstructure:"temperature"`
	MaxTokens      int64   `mapstructure:"maxTokens"`
	TopP           float64 `mapstructure:"topP"`
	ResponseFormat string  `mapstructure:"responseFormat"`
}

// Run executes the LLM request and returns the generated response.
func (l *LLMActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters
	var props llmActionProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	fmt.Println("props", props)
	fmt.Println("resolvedProps", resolvedProps)

	if props.APIKey == "" {
		return node.NodeResult{
			Output: map[string]any{
				"response": "",
			},
			Error:         fmt.Errorf("API key is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	if props.UserPrompt == "" {
		return node.NodeResult{
			Output: map[string]any{
				"response": "",
			},
			Error:         fmt.Errorf("user prompt is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Set defaults
	if props.Temperature == 0 {
		props.Temperature = 0.7
	}
	if props.MaxTokens == 0 {
		props.MaxTokens = 2048
	}
	if props.TopP == 0 {
		props.TopP = 1.0
	}
	if props.ResponseFormat == "" {
		props.ResponseFormat = "text"
	}
	var response string
	var err error
	switch props.Provider {
	case "anthropic":
		response, err = anthropicApi(ctx, props)
	case "openai":
		response, err = openaiApi(ctx, props)
	}
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"response": "",
			},
			Error:         err,
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"response": response,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

func anthropicApi(ctx context.Context, props llmActionProps) (response string, err error) {

	client := anthropic.NewClient(
		anthropicOption.WithAPIKey(props.APIKey),
	)
	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: props.MaxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewAssistantMessage(anthropic.NewTextBlock(props.SystemPrompt)),
			anthropic.NewUserMessage(anthropic.NewTextBlock(props.UserPrompt)),
		},
		System: []anthropic.TextBlockParam{
			{Text: props.SystemPrompt},
		},
		Model:       anthropic.Model(props.Model),
		Temperature: anthropic.Float(props.Temperature),
		TopP:        anthropic.Float(props.TopP),
	})
	if err != nil {
		return "", err
	}
	response = message.Content[0].Text
	return
}

func openaiApi(ctx context.Context, props llmActionProps) (response string, err error) {
	client := openai.NewClient(
		option.WithAPIKey(props.APIKey),
	)
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(props.SystemPrompt),
			openai.UserMessage(props.UserPrompt),
		},
		Model:       openai.ChatModel(props.Model),
		Temperature: openai.Float(props.Temperature),
		TopP:        openai.Float(props.TopP),
	})
	if err != nil {
		return "", err
	}

	response = chatCompletion.Choices[0].Message.Content
	return
}
