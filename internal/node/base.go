package node

import (
	"context"
	"fmt"
)

// BaseNode provides a default implementation of the Node interface.
// It should be embedded in concrete node implementations.
type BaseNode struct {
	nodeInfo      NodeInfo
	properties    []NodeProperty
	output        []NodeProperty
	outputHandles []OutputHandle
}

// NewBaseNode creates a new BaseNode with properties, outputs, and handles.
func NewBaseNode(ref string, nodeType NodeType, name string, description string, icon string, properties []NodeProperty, output []NodeProperty, outputHandles []OutputHandle) BaseNode {
	nodeInfo := NodeInfo{
		Ref:         ref,
		Type:        nodeType,
		Name:        name,
		Description: description,
		Icon:        getDefaultIcon(nodeType),
	}
	if icon != "" {
		nodeInfo.Icon = icon
	}
	
	// If no output handles specified, create a default one
	if outputHandles == nil {
		outputHandles = []OutputHandle{
			{ID: "default"},
		}
	}
	
	return BaseNode{
		nodeInfo:      nodeInfo,
		properties:    properties,
		output:        output,
		outputHandles: outputHandles,
	}
}

// Ref returns the unique identifier of the node.
func (r *BaseNode) Ref() string {
	return r.nodeInfo.Ref
}

// Name returns the human-readable name of the node.
func (r *BaseNode) Name() string {
	return r.nodeInfo.Name
}

// Icon returns the icon name for the node.
// If no icon is set, returns a default based on node type.
func (r *BaseNode) Icon() string {
	if r.nodeInfo.Icon == "" {
		switch r.nodeInfo.Type {
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
	return r.nodeInfo.Icon
}

// Type returns the category of the node.
func (r *BaseNode) Type() NodeType {
	return r.nodeInfo.Type
}

// Description returns the description of the node's functionality.
func (r *BaseNode) Description() string {
	return r.nodeInfo.Description
}

// Info returns the complete node metadata.
func (r *BaseNode) Info() NodeInfo {
	return r.nodeInfo
}

// Properties returns the static input properties.
// This provides access to the base properties defined in NewBaseNode.
func (r *BaseNode) Properties() []NodeProperty {
	return r.properties
}

// Output returns the static output properties.
// This provides access to the base output defined in NewBaseNode.
func (r *BaseNode) Output() []NodeProperty {
	return r.output
}

// OutputHandles returns the static output handles.
// This provides access to the base output handles defined in NewBaseNode.
func (r *BaseNode) OutputHandles() []OutputHandle {
	return r.outputHandles
}

// GetProperties returns dynamic properties based on context.
// Default implementation returns the static properties.
// Concrete nodes can override this for dynamic behavior.
func (r *BaseNode) GetProperties(ctx PropertyContext) []NodeProperty {
	return r.properties
}

// GetOutput returns dynamic output properties based on context.
// Default implementation returns the static output.
// Concrete nodes can override this for dynamic behavior.
func (r *BaseNode) GetOutput(ctx PropertyContext) []NodeProperty {
	return r.output
}

// GetOutputHandles returns dynamic output handles based on context.
// Default implementation returns the static output handles.
// Concrete nodes can override this for dynamic behavior.
func (r *BaseNode) GetOutputHandles(ctx PropertyContext) []OutputHandle {
	return r.outputHandles
}

// ValidateProperties validates the properties of the node.
// 기본적으로는 아무것도 하지 않지만 각 노드 구현에서 재정의 가능
func (r *BaseNode) ValidateProperties(input map[string]any) error {
	return nil
}

// Run executes the node logic. This method panics by default and must be
// overridden by concrete implementations.
func (r *BaseNode) Run(ctx context.Context, states map[string]any) (result NodeResult) {
	return NodeResult{
		Output:   nil,
		Error:    fmt.Errorf("%s [%s] not implemented please override Run(ctx, params) method", r.nodeInfo.Name, r.nodeInfo.Ref),
		Continue: false,
	}
}
