// Package node provides the core node system for workflow automation.
// It defines the base interfaces and types that all workflow nodes must implement.
package node

import (
	"context"
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
// TODO: Add a way to specify the next node to execute (edges)
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
}
