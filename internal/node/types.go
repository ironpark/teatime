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

// BaseNode provides a default implementation of the Node interface.
// It should be embedded in concrete node implementations.
type BaseNode struct {
	nodeInfo NodeInfo
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

// NewBaseNode creates a new BaseNode with the specified properties.
// If icon is empty, a default icon based on nodeType will be used.
func NewBaseNode(ref string, nodeType NodeType, name string, description string, icon string) *BaseNode {
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
	return &BaseNode{
		nodeInfo: nodeInfo,
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

// Run executes the node logic. This method panics by default and must be
// overridden by concrete implementations.
func (r *BaseNode) Run(ctx context.Context, states map[string]any) (result NodeResult) {
	return NodeResult{
		Output:   nil,
		Error:    fmt.Errorf("%s [%s] not implemented please override Run(ctx, params) method", r.nodeInfo.Name, r.nodeInfo.Ref),
		Continue: false,
	}
}

func (r *BaseNode) ResolveInput(states map[string]any) (map[string]any, error) {
	return states, fmt.Errorf("%s [%s] not implemented please override ResolveInput(states) method", r.nodeInfo.Name, r.nodeInfo.Ref)
}

func (r *BaseNode) Validate(input map[string]any) error {
	return fmt.Errorf("%s [%s] not implemented please override Validate(input) method", r.nodeInfo.Name, r.nodeInfo.Ref)
}

type NodeResult struct {
	Output   map[string]any
	Error    error
	Continue bool
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
	// Properties returns the input properties that can be configured for this node.
	Properties() []NodeProperty
	// Output returns the output properties that this node produces.
	Output() []NodeProperty
	// Info returns the complete metadata for the node.
	Info() NodeInfo
	// Run executes the node with the given context and state.
	Run(ctx context.Context, states map[string]any) (result NodeResult)
	// Validate checks if the provided input is valid for this node.
	Validate(input map[string]any) error
	// ResolveInput evaluates expressions and resolves bindings to produce actual input values.
	ResolveInput(states map[string]any) (map[string]any, error)
}
