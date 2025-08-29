package agent

import "context"

// Provider defines the interface for AI language model providers
type Provider interface {
	// Complete generates a completion for the given messages
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	
	// GetModel returns the model identifier
	GetModel() string
	
	// GetProvider returns the provider name
	GetProvider() string
	
	// SupportsTools returns whether the agent supports tool calling
	SupportsTools() bool
	
	// GetMaxTokens returns the maximum context length
	GetMaxTokens() int
}