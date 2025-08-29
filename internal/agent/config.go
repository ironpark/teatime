package agent

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AgentConfig represents common configuration for AI agents
type AgentConfig struct {
	// Required fields
	APIKey   string `json:"api_key" yaml:"api_key"`
	Provider string `json:"provider" yaml:"provider"`
	
	// Optional fields with defaults
	Model        string        `json:"model,omitempty" yaml:"model,omitempty"`
	MaxTokens    int          `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
	Temperature  float64      `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	TopP         float64      `json:"top_p,omitempty" yaml:"top_p,omitempty"`
	Timeout      time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RetryAttempts int         `json:"retry_attempts,omitempty" yaml:"retry_attempts,omitempty"`
	BaseURL      string       `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	
	// Advanced options
	Tools          []string               `json:"tools,omitempty" yaml:"tools,omitempty"`
	SystemPrompt   string                `json:"system_prompt,omitempty" yaml:"system_prompt,omitempty"`
	StopSequences  []string              `json:"stop_sequences,omitempty" yaml:"stop_sequences,omitempty"`
	Metadata       map[string]any        `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Validator defines interface for configuration validation
type Validator interface {
	Validate() error
}

// Validate validates the agent configuration
func (c *AgentConfig) Validate() error {
	var errs []error
	
	// Required fields
	if strings.TrimSpace(c.APIKey) == "" {
		errs = append(errs, NewValidationError("api_key", c.APIKey, "API key is required"))
	}
	
	if strings.TrimSpace(c.Provider) == "" {
		errs = append(errs, NewValidationError("provider", c.Provider, "provider is required"))
	}
	
	// Optional field validation
	if c.MaxTokens < 0 || c.MaxTokens > 200000 {
		errs = append(errs, NewValidationError("max_tokens", c.MaxTokens, "max_tokens must be between 0 and 200000"))
	}
	
	if c.Temperature < 0.0 || c.Temperature > 2.0 {
		errs = append(errs, NewValidationError("temperature", c.Temperature, "temperature must be between 0.0 and 2.0"))
	}
	
	if c.TopP < 0.0 || c.TopP > 1.0 {
		errs = append(errs, NewValidationError("top_p", c.TopP, "top_p must be between 0.0 and 1.0"))
	}
	
	if c.RetryAttempts < 0 || c.RetryAttempts > 10 {
		errs = append(errs, NewValidationError("retry_attempts", c.RetryAttempts, "retry_attempts must be between 0 and 10"))
	}
	
	if c.Timeout > 0 && c.Timeout < time.Second {
		errs = append(errs, NewValidationError("timeout", c.Timeout, "timeout must be at least 1 second"))
	}
	
	// Provider-specific validation
	if err := c.validateProviderSpecific(); err != nil {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		return &MultiValidationError{Errors: errs}
	}
	
	return nil
}

// validateProviderSpecific performs provider-specific validation
func (c *AgentConfig) validateProviderSpecific() error {
	switch strings.ToLower(c.Provider) {
	case "anthropic":
		return c.validateAnthropicConfig()
	case "openai":
		return c.validateOpenAIConfig()
	default:
		return NewValidationError("provider", c.Provider, "unsupported provider")
	}
}

// validateAnthropicConfig validates Anthropic-specific configuration
func (c *AgentConfig) validateAnthropicConfig() error {
	validModels := []string{
		"claude-3-5-sonnet-20241022", "claude-3-5-sonnet-latest",
		"claude-3-5-haiku-20241022", "claude-3-5-haiku-latest",
		"claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307",
	}
	
	if c.Model != "" && !contains(validModels, c.Model) {
		return NewValidationError("model", c.Model, "invalid model for Anthropic provider")
	}
	
	return nil
}

// validateOpenAIConfig validates OpenAI-specific configuration
func (c *AgentConfig) validateOpenAIConfig() error {
	validModels := []string{
		"gpt-4", "gpt-4-turbo", "gpt-4o", "gpt-4o-mini",
		"gpt-3.5-turbo", "gpt-3.5-turbo-16k",
		"gpt-5", "gpt-5-mini", "o1", "o1-mini",
	}
	
	if c.Model != "" && !contains(validModels, c.Model) {
		return NewValidationError("model", c.Model, "invalid model for OpenAI provider")
	}
	
	return nil
}

// SetDefaults sets default values for the configuration
func (c *AgentConfig) SetDefaults() {
	if c.MaxTokens == 0 {
		c.MaxTokens = 4096
	}
	
	if c.Temperature == 0 {
		c.Temperature = 0.7
	}
	
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	
	if c.RetryAttempts == 0 {
		c.RetryAttempts = 3
	}
	
	// Set provider-specific defaults
	switch strings.ToLower(c.Provider) {
	case "anthropic":
		if c.Model == "" {
			c.Model = "claude-3-5-sonnet-latest"
		}
	case "openai":
		if c.Model == "" {
			c.Model = "gpt-4o"
		}
	}
}

// MultiValidationError represents multiple validation errors
type MultiValidationError struct {
	Errors []error
}

// Error implements the error interface
func (e *MultiValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "no validation errors"
	}
	
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	
	return fmt.Sprintf("multiple validation errors: %s", strings.Join(messages, "; "))
}

// HealthChecker defines interface for health checking
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Status    HealthStatus      `json:"status"`
	Message   string           `json:"message,omitempty"`
	Details   map[string]any   `json:"details,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

// HealthStatus represents the health status
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}