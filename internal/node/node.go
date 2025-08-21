// Package node provides the core node system for workflow automation.
// It defines the base interfaces and types that all workflow nodes must implement.
package node

import (
	"context"
	"fmt"
)

// NodeType represents the category of a workflow node.
type NodeType string

const (
	NodeTypeTrigger NodeType = "trigger" // NodeTypeTrigger initiates workflow execution
	NodeTypeBranch  NodeType = "branch"  // NodeTypeBranch controls workflow flow logic
	NodeTypeAction  NodeType = "action"  // NodeTypeAction performs operations or side effects
	NodeTypeUtil    NodeType = "util"    // NodeTypeUtil provides utility functions
)

// NodeInfo contains metadata about a workflow node.
type NodeInfo struct {
	Ref         string   `json:"ref"`         // Unique identifier for the node type
	Name        string   `json:"name"`        // Human-readable name
	Description string   `json:"description"` // Brief description of node functionality
	Type        NodeType `json:"type"`        // Category of the node
	Icon        string   `json:"icon"`        // Icon name (Lucide icon set)
}

// getDefaultIcon returns the default icon for a given node type.
func getDefaultIcon(nodeType NodeType) string {
	switch nodeType {
	case NodeTypeTrigger:
		return "Zap"
	case NodeTypeBranch:
		return "GitBranch"
	case NodeTypeAction:
		return "Play"
	case NodeTypeUtil:
		return "Settings"
	default:
		return "Activity"
	}
}

// NodeResult is the result of a node execution.
// It includes output data, error status, and the output handles to activate for routing.
type NodeResult struct {
	Output        map[string]any
	Error         error
	Continue      bool
	OutputHandles []string // IDs of output handles to activate (e.g., ["true"], ["success", "log"])
}

// PropertyContext provides context for dynamic property resolution.
// It contains the current property values that nodes can use to
// dynamically adjust their properties and outputs.
type PropertyContext map[string]any

// WorkflowState represents the global state of a workflow execution.
// It stores all node inputs and outputs with structured keys for easy access.
type WorkflowState map[string]any

// GetOutput retrieves an output value from a specific node.
func (ws WorkflowState) GetOutput(nodeId string, key string) any {
	return ws[fmt.Sprintf("%s.output.%s", nodeId, key)]
}

// GetInput retrieves an input value from a specific node.
func (ws WorkflowState) GetInput(nodeId string, key string) any {
	return ws[fmt.Sprintf("%s.input.%s", nodeId, key)]
}

// SetOutput sets an output value for a specific node.
func (ws WorkflowState) SetOutput(nodeId string, key string, value any) {
	ws[fmt.Sprintf("%s.output.%s", nodeId, key)] = value
}

// SetOutputs sets multiple output values for a specific node.
func (ws WorkflowState) SetOutputs(nodeId string, outputs map[string]any) {
	for key, value := range outputs {
		ws.SetOutput(nodeId, key, value)
	}
}

// SetInput sets an input value for a specific node.
func (ws WorkflowState) SetInput(nodeId string, key string, value any) {
	ws[fmt.Sprintf("%s.input.%s", nodeId, key)] = value
}

// SetInputs sets multiple input values for a specific node.
func (ws WorkflowState) SetInputs(nodeId string, inputs map[string]any) {
	for key, value := range inputs {
		ws.SetInput(nodeId, key, value)
	}
}

// OutputHandle represents a connection point from a node.
// Nodes can have multiple output handles for different execution paths.
type OutputHandle struct {
	ID          string `json:"id"`                    // Unique identifier for this handle
	Label       string `json:"label,omitempty"`       // Optional human-readable label
	Description string `json:"description,omitempty"` // Optional description of when this handle is used
}

// Node defines the interface that all workflow nodes must implement.
type Node interface {
	// Ref returns the unique identifier for the node type.
	Ref() string
	// Name returns the human-readable name of the node.
	Name() string
	// Type returns the category of the node.
	Type() NodeType
	// Icon returns the icon name (from Lucide icon set).
	Icon() string
	// Description returns a brief description of the node's functionality.
	Description() string
	// GetProperties returns dynamic properties based on context.
	// Nodes can override this to provide context-aware properties.
	GetProperties(ctx PropertyContext) []NodeProperty
	// GetOutput returns dynamic output properties based on context.
	// Nodes can override this to provide context-aware outputs.
	GetOutput(ctx PropertyContext) []NodeProperty
	// GetOutputHandles returns dynamic output handles based on context.
	// Nodes can override this to provide context-aware handles.
	GetOutputHandles(ctx PropertyContext) []OutputHandle
	// Info returns the complete metadata for the node.
	Info() NodeInfo
	// Run executes the node with the given context and state.
	Run(ctx context.Context, resolvedProperties PropertyContext, states WorkflowState) (result NodeResult)
	// ValidateProperties validates the properties of the node.
	ValidateProperties(input PropertyContext) error
}
