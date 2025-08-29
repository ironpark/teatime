package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Message represents a single message in a conversation
type Message struct {
	ID          string       `json:"id"`
	Role        Role         `json:"role"`
	Content     string       `json:"content"`
	Timestamp   time.Time    `json:"timestamp"`
	ToolCalls   []ToolCall   `json:"tool_calls,omitempty"`
	ToolResults []ToolResult `json:"tool_results,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// Role defines the type of message sender
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
	RoleTool      Role = "tool"
)

// ToolCall represents a function call made by the agent
type ToolCall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
	IsError    bool   `json:"is_error"`
	Data       any    `json:"data,omitempty"`
}

// Conversation represents a conversation with metadata
type Conversation struct {
	ID           string            `json:"id"`
	Title        string            `json:"title,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	MessageCount int               `json:"message_count"`
	Metadata     map[string]any    `json:"metadata,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
}

// SearchQuery represents a query for searching conversations or messages
type SearchQuery struct {
	Query     string            `json:"query,omitempty"`
	Roles     []Role            `json:"roles,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
	StartTime *time.Time        `json:"start_time,omitempty"`
	EndTime   *time.Time        `json:"end_time,omitempty"`
	Limit     int               `json:"limit,omitempty"`
	Offset    int               `json:"offset,omitempty"`
	Metadata  map[string]any    `json:"metadata,omitempty"`
}

// SearchResult represents search results
type SearchResult struct {
	Messages      []Message      `json:"messages,omitempty"`
	Conversations []Conversation `json:"conversations,omitempty"`
	Total         int            `json:"total"`
	HasMore       bool           `json:"has_more"`
}

// CompletionRequest contains parameters for generating completions
type CompletionRequest struct {
	Messages      []Message `json:"messages"`
	MaxTokens     int       `json:"max_tokens,omitempty"`
	Temperature   float64   `json:"temperature,omitempty"`
	StopSequences []string  `json:"stop_sequences,omitempty"`
	Tools         []Tool    `json:"tools,omitempty"`
}

// CompletionResponse represents the response from completion
type CompletionResponse struct {
	Message       Message `json:"message"`
	StopReason    string  `json:"stop_reason"`
	Usage         Usage   `json:"usage"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// Memory defines the interface for conversation storage
type Memory interface {
	// Conversation management
	CreateConversation(ctx context.Context, title string) (*Conversation, error)
	GetConversation(ctx context.Context, conversationID string) (*Conversation, error)
	UpdateConversation(ctx context.Context, conversationID string, updates map[string]any) error
	DeleteConversation(ctx context.Context, conversationID string) error
	ListConversations(ctx context.Context, limit, offset int) ([]Conversation, error)

	// Message management
	AddMessage(ctx context.Context, conversationID string, message Message) error
	GetMessages(ctx context.Context, conversationID string, limit int) ([]Message, error)
	GetMessagesByRange(ctx context.Context, conversationID string, startTime, endTime time.Time) ([]Message, error)
	UpdateMessage(ctx context.Context, messageID string, updates map[string]any) error
	DeleteMessage(ctx context.Context, messageID string) error

	// Bulk operations
	AddMessages(ctx context.Context, conversationID string, messages []Message) error
	ClearConversation(ctx context.Context, conversationID string) error

	// Search and filtering
	SearchMessages(ctx context.Context, query SearchQuery) (*SearchResult, error)
	SearchConversations(ctx context.Context, query SearchQuery) (*SearchResult, error)
}

// ExecuteToolCalls executes all tool calls in the message concurrently and returns results
func ExecuteToolCalls(ctx context.Context, message Message, toolRegistry *ToolRegistry) []ToolResult {
	if len(message.ToolCalls) == 0 {
		return nil
	}
	
	// Create channels for results and errors
	resultChan := make(chan ToolResult, len(message.ToolCalls))
	var wg sync.WaitGroup
	
	// Execute each tool call concurrently
	for _, toolCall := range message.ToolCalls {
		wg.Add(1)
		go func(tc ToolCall) {
			defer wg.Done()
			
			toolResult, err := toolRegistry.Execute(ctx, tc.Name, tc.Arguments)
			if err != nil {
				// Create error result
				errorResult := ToolResult{
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf("Tool %s failed: %v", tc.Name, err),
					IsError:    true,
				}
				resultChan <- errorResult
				return
			}
			
			// Create success result
			successResult := ToolResult{
				ToolCallID: tc.ID,
				Content:    toolResult.Content,
				IsError:    toolResult.IsError,
				Data:       toolResult.Data,
			}
			resultChan <- successResult
		}(toolCall)
	}
	
	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect all results
	var toolResultsData []ToolResult
	for result := range resultChan {
		toolResultsData = append(toolResultsData, result)
	}
	
	return toolResultsData
}