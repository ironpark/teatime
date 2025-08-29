package agent

import (
	"errors"
	"fmt"
)

// Common error variables
var (
	ErrToolNotFound       = errors.New("tool not found")
	ErrInvalidConfig      = errors.New("invalid configuration")
	ErrConversationNotFound = errors.New("conversation not found")
	ErrToolHandlerNil     = errors.New("tool handler is nil")
	ErrEmptyToolName      = errors.New("tool name cannot be empty")
	ErrInvalidParameters  = errors.New("invalid parameters")
	ErrAPIError          = errors.New("API error")
	ErrContextCanceled   = errors.New("context canceled")
	ErrRateLimited       = errors.New("rate limited")
)

// AgentError represents an error from an AI agent operation
type AgentError struct {
	Type     ErrorType
	Provider string
	Message  string
	Cause    error
}

// ErrorType categorizes different types of errors
type ErrorType string

const (
	ErrorTypeAPI          ErrorType = "api"
	ErrorTypeAuth         ErrorType = "auth"
	ErrorTypeRateLimit    ErrorType = "rate_limit"
	ErrorTypeInvalidInput ErrorType = "invalid_input"
	ErrorTypeTimeout      ErrorType = "timeout"
	ErrorTypeInternal     ErrorType = "internal"
	ErrorTypeTool         ErrorType = "tool"
	ErrorTypeMemory       ErrorType = "memory"
)

// Error implements the error interface
func (e *AgentError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s error from %s: %s (caused by: %v)", e.Type, e.Provider, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s error from %s: %s", e.Type, e.Provider, e.Message)
}

// Unwrap returns the underlying cause error
func (e *AgentError) Unwrap() error {
	return e.Cause
}

// Is implements error matching for errors.Is
func (e *AgentError) Is(target error) bool {
	if target == nil {
		return false
	}
	
	var agentErr *AgentError
	if errors.As(target, &agentErr) {
		return e.Type == agentErr.Type
	}
	
	return e.Cause != nil && errors.Is(e.Cause, target)
}

// NewAgentError creates a new agent error
func NewAgentError(errorType ErrorType, provider, message string, cause error) *AgentError {
	return &AgentError{
		Type:     errorType,
		Provider: provider,
		Message:  message,
		Cause:    cause,
	}
}

// WrapError wraps an existing error with agent error context
func WrapError(errorType ErrorType, provider string, cause error) *AgentError {
	return &AgentError{
		Type:     errorType,
		Provider: provider,
		Message:  cause.Error(),
		Cause:    cause,
	}
}

// ToolError represents an error from tool execution
type ToolError struct {
	ToolName string
	Message  string
	Cause    error
}

// Error implements the error interface
func (e *ToolError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("tool '%s' error: %s (caused by: %v)", e.ToolName, e.Message, e.Cause)
	}
	return fmt.Sprintf("tool '%s' error: %s", e.ToolName, e.Message)
}

// Unwrap returns the underlying cause error
func (e *ToolError) Unwrap() error {
	return e.Cause
}

// NewToolError creates a new tool error
func NewToolError(toolName, message string, cause error) *ToolError {
	return &ToolError{
		ToolName: toolName,
		Message:  message,
		Cause:    cause,
	}
}

// ValidationError represents configuration or parameter validation errors
type ValidationError struct {
	Field   string
	Value   any
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("validation error for field '%s' with value '%v': %s", e.Field, e.Value, e.Message)
	}
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value any, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}