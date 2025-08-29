package agent

import (
	"context"
	"fmt"
	"time"
)

// AiAgent represents an AI agent with a provider and configuration
type AiAgent struct {
	Name     string
	Provider Provider
	Config   map[string]any
	Tools    *ToolRegistry
	Memory
}

// NewAgent creates a new AI agent with the given provider
func NewAgent(name string, provider Provider) *AiAgent {
	return &AiAgent{
		Name:     name,
		Provider: provider,
		Config:   make(map[string]any),
		Tools:    NewToolRegistry(),
		Memory:   NewInMemoryStore(),
	}
}

// AddTool adds a tool to the agent's tool registry
func (a *AiAgent) AddTool(tool Tool) error {
	return a.Tools.Register(tool)
}

// GetTool retrieves a tool by name from the agent's registry
func (a *AiAgent) GetTool(name string) (*Tool, bool) {
	return a.Tools.Get(name)
}

// ListTools returns all available tools for this agent
func (a *AiAgent) ListTools() []Tool {
	return a.Tools.List()
}

// SendMessage sends a message and handles the complete interaction cycle
func (agent *AiAgent) SendMessage(ctx context.Context, conversationID string, content string) (*CompletionResponse, error) {
	// Add user message to memory
	userMsg := Message{
		Role:      RoleUser,
		Content:   content,
		Timestamp: time.Now(),
	}

	if err := agent.AddMessage(ctx, conversationID, userMsg); err != nil {
		return nil, fmt.Errorf("failed to add user message: %w", err)
	}

	// Get conversation history
	messages, err := agent.GetMessages(ctx, conversationID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Prepare completion request
	req := CompletionRequest{
		Messages:    messages,
		MaxTokens:   2000,
		Temperature: 0.7,
	}

	// Add tools if available
	if agent.Tools != nil {
		req.Tools = agent.Tools.List()
	}

	// Generate completion
	response, err := agent.Provider.Complete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate completion: %w", err)
	}
	agent.AddMessage(ctx, conversationID, response.Message)

	for len(response.Message.ToolCalls) > 0 {
		// Execute all tool calls using ExecuteToolCalls function
		toolResultsData := ExecuteToolCalls(ctx, response.Message, agent.Tools)
		if err := agent.AddMessage(ctx, conversationID, Message{
			Role:        RoleTool,
			ToolResults: toolResultsData,
			Timestamp:   time.Now(),
		}); err != nil {
			return nil, fmt.Errorf("failed to add tool result message: %w", err)
		}
		// Get updated conversation for follow-up response
		updatedMessages, err := agent.GetMessages(ctx, conversationID, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get updated conversation: %w", err)
		}

		// Get follow-up response from agent
		followUpReq := CompletionRequest{
			Messages:    updatedMessages,
			MaxTokens:   2000,
			Temperature: 0.7,
		}

		// Include tools for potential additional calls
		if agent.Tools != nil {
			followUpReq.Tools = agent.Tools.List()
		}

		followUpResponse, err := agent.Provider.Complete(ctx, followUpReq)
		if err != nil {
			return nil, fmt.Errorf("failed to generate follow-up response: %w", err)
		}

		// Add follow-up response to memory
		if err := agent.AddMessage(ctx, conversationID, followUpResponse.Message); err != nil {
			return nil, fmt.Errorf("failed to add follow-up message: %w", err)
		}

		response = followUpResponse
	}

	return response, nil
}
