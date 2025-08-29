package anthropic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/ironpark/teatime/internal/agent"
)

// Config holds the configuration for the Anthropic client
type Config struct {
	APIKey      string
	Model       anthropic.Model
	MaxTokens   int
	Temperature float64
}

// Client implements the agent.AiAgent interface for Anthropic's Claude
type Client struct {
	config Config
	client *anthropic.Client
}

// NewClient creates a new Anthropic client
func NewClient(config Config) *Client {
	if config.Model == "" {
		config.Model = anthropic.ModelClaude3_5SonnetLatest
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4096
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}

	client := anthropic.NewClient(option.WithAPIKey(config.APIKey))

	return &Client{
		config: config,
		client: &client,
	}
}

// Complete generates a completion using Anthropic's API
func (c *Client) Complete(ctx context.Context, req agent.CompletionRequest) (*agent.CompletionResponse, error) {
	// Convert agent messages to Anthropic format
	var messages []anthropic.MessageParam
	var systemMessage string
	// fmt.Printf("🔍 %+v\n", len(req.Messages))
	// for _, msg := range req.Messages {
	// 	fmt.Printf("🔍 Message: %+v\n", msg)
	// }
	for _, msg := range req.Messages {
		contentBlocks := []anthropic.ContentBlockParamUnion{}
		// Handle assistant messages
		if msg.Content != "" {
			contentBlocks = append(contentBlocks, anthropic.NewTextBlock(msg.Content))
		}
		switch msg.Role {
		case agent.RoleSystem:
			systemMessage = msg.Content
		case agent.RoleUser, agent.RoleTool:
			for _, toolResult := range msg.ToolResults {
				contentBlocks = append(contentBlocks, anthropic.NewToolResultBlock(
					toolResult.ToolCallID,
					toolResult.Content,
					toolResult.IsError,
				))
			}
			messages = append(messages, anthropic.NewUserMessage(contentBlocks...))
		case agent.RoleAssistant:
			for _, call := range msg.ToolCalls {
				contentBlocks = append(contentBlocks, anthropic.NewToolUseBlock(
					call.ID,
					call.Arguments,
					call.Name,
				))
			}
			messages = append(messages, anthropic.NewAssistantMessage(contentBlocks...))
		}
	}

	// Set max tokens
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = c.config.MaxTokens
	}

	// Set temperature
	temperature := req.Temperature
	if temperature == 0 {
		temperature = c.config.Temperature
	}

	// Build message parameters
	params := anthropic.MessageNewParams{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   int64(maxTokens),
		Temperature: anthropic.Float(temperature),
	}

	// Set system message
	if systemMessage != "" {
		params.System = []anthropic.TextBlockParam{
			{
				Text: systemMessage,
			},
		}
	}

	// Set stop sequences
	if len(req.StopSequences) > 0 {
		params.StopSequences = req.StopSequences
	}

	// Add tools if provided and supported
	if len(req.Tools) > 0 {
		var tools []anthropic.ToolUnionParam
		for _, tool := range req.Tools {
			// Create tool parameter properly
			toolParam := anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type:       "object",
					Properties: convertParameters(tool.Parameters),
					Required:   getRequiredParameters(tool.Parameters),
				},
			}
			tools = append(tools, anthropic.ToolUnionParam{OfTool: &toolParam})
		}
		params.Tools = tools
	}

	// Call the API
	response, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert response to our format with tool support
	var content string
	var toolCalls []agent.ToolCall

	// Process all content blocks
	for _, block := range response.Content {
		switch block.Type {
		case "text":
			content += block.Text

		case "tool_use":
			// Parse Input from JSON RawMessage
			var input map[string]any
			if err := json.Unmarshal(block.Input, &input); err != nil {
				continue // Skip malformed input
			}

			toolCall := agent.ToolCall{
				ID:        block.ID,
				Name:      block.Name,
				Arguments: input,
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}

	return &agent.CompletionResponse{
		Message: agent.Message{
			ID:        response.ID,
			Role:      agent.RoleAssistant,
			Content:   content,
			Timestamp: time.Now(),
			ToolCalls: toolCalls,
		},
		StopReason: string(response.StopReason),
		Usage: agent.Usage{
			InputTokens:  int(response.Usage.InputTokens),
			OutputTokens: int(response.Usage.OutputTokens),
			TotalTokens:  int(response.Usage.InputTokens + response.Usage.OutputTokens),
		},
	}, nil
}

// GetModel returns the model identifier
func (c *Client) GetModel() string {
	return string(c.config.Model)
}

// GetProvider returns the provider name
func (c *Client) GetProvider() string {
	return "anthropic"
}

// SupportsTools returns whether the agent supports tool calling
func (c *Client) SupportsTools() bool {
	model := string(c.config.Model)
	return strings.Contains(model, "claude-3") || strings.Contains(model, "sonnet") || strings.Contains(model, "haiku") || strings.Contains(model, "opus")
}

// GetMaxTokens returns the maximum context length
func (c *Client) GetMaxTokens() int {
	model := string(c.config.Model)
	switch {
	case strings.Contains(model, "claude-3-5-sonnet") || strings.Contains(model, "claude-3-5-haiku"):
		return 200000
	case strings.Contains(model, "claude-3"):
		return 200000
	default:
		return 200000 // Default for newer models
	}
}

// convertParameters converts agent.Parameter map to Anthropic's schema format
func convertParameters(params map[string]agent.Parameter) map[string]any {
	properties := make(map[string]any)

	for name, param := range params {
		prop := map[string]any{
			"type":        string(param.Type),
			"description": param.Description,
		}

		if len(param.Enum) > 0 {
			prop["enum"] = param.Enum
		}

		if param.Items != nil {
			prop["items"] = map[string]any{
				"type":        string(param.Items.Type),
				"description": param.Items.Description,
			}
		}

		if param.Properties != nil {
			prop["properties"] = convertParameterPointers(param.Properties)
		}

		properties[name] = prop
	}

	return properties
}

// convertParameterPointers converts parameter pointers to interface map
func convertParameterPointers(params map[string]*agent.Parameter) map[string]any {
	properties := make(map[string]any)

	for name, param := range params {
		if param != nil {
			prop := map[string]any{
				"type":        string(param.Type),
				"description": param.Description,
			}
			properties[name] = prop
		}
	}

	return properties
}

// getRequiredParameters extracts required parameter names
func getRequiredParameters(params map[string]agent.Parameter) []string {
	var required []string
	for name, param := range params {
		if param.Required {
			required = append(required, name)
		}
	}
	return required
}
