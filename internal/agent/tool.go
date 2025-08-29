package agent

import (
	"context"
	"fmt"
)

// ParameterType defines the type of a tool parameter
type ParameterType string

const (
	ParameterTypeString  ParameterType = "string"
	ParameterTypeNumber  ParameterType = "number"
	ParameterTypeInteger ParameterType = "integer"
	ParameterTypeBoolean ParameterType = "boolean"
	ParameterTypeArray   ParameterType = "array"
	ParameterTypeObject  ParameterType = "object"
)

// Parameter describes a tool parameter
type Parameter struct {
	Type        ParameterType         `json:"type"`
	Description string                `json:"description"`
	Required    bool                  `json:"required,omitempty"`
	Enum        []any                 `json:"enum,omitempty"`
	Items       *Parameter            `json:"items,omitempty"`      // For array type
	Properties  map[string]*Parameter `json:"properties,omitempty"` // For object type
	Default     any                   `json:"default,omitempty"`
}

// Tool represents a function that can be called by an AI agent
type Tool struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Parameters  map[string]Parameter `json:"parameters"`
	Handler     ToolHandler          `json:"-"` // Not serialized
}

// ToolHandler is the function signature for tool implementations
type ToolHandler func(ctx context.Context, arguments map[string]any) (*ToolResult, error)

// ToolRegistry manages available tools for agents
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool Tool) error {
	if tool.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	r.tools[tool.Name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (*Tool, bool) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, false
	}
	return &tool, true
}

// List returns all available tools
func (r *ToolRegistry) List() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// Unregister removes a tool from the registry
func (r *ToolRegistry) Unregister(name string) bool {
	if _, exists := r.tools[name]; exists {
		delete(r.tools, name)
		return true
	}
	return false
}

// Execute runs a tool with the given arguments
func (r *ToolRegistry) Execute(ctx context.Context, name string, arguments map[string]any) (*ToolResult, error) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not found", name)
	}
	
	if tool.Handler == nil {
		return nil, fmt.Errorf("tool '%s' has no handler", name)
	}
	
	return tool.Handler(ctx, arguments)
}


// NewParameter creates a new parameter with the given type and description
func NewParameter(paramType ParameterType, description string) Parameter {
	return Parameter{
		Type:        paramType,
		Description: description,
	}
}

// WithRequired marks the parameter as required
func (p Parameter) WithRequired(required bool) Parameter {
	p.Required = required
	return p
}

// WithEnum sets the allowed values for the parameter
func (p Parameter) WithEnum(values ...any) Parameter {
	p.Enum = values
	return p
}

// WithDefault sets the default value for the parameter
func (p Parameter) WithDefault(value any) Parameter {
	p.Default = value
	return p
}

// WithItems sets the item type for array parameters
func (p Parameter) WithItems(items Parameter) Parameter {
	p.Items = &items
	return p
}

// WithProperties sets the properties for object parameters
func (p Parameter) WithProperties(properties map[string]*Parameter) Parameter {
	p.Properties = properties
	return p
}