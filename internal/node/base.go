package node

import (
	"context"
	"fmt"
	"strings"
)

// BaseNode provides a default implementation of the Node interface.
// It should be embedded in concrete node implementations.
type BaseNode struct {
	nodeInfo NodeInfo
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
	resolvedStates := make(map[string]any)
	for key, value := range states {
		switch v := value.(type) {
		// it can be a expression, binding, or a value
		// binding is @[bindingName]
		// expression is {{expression}}
		case string:
			if strings.HasPrefix(v, "@[") && strings.HasSuffix(v, "]") {
				// it is a binding
				binding := strings.TrimSpace(v[2 : len(v)-1])
				resolvedStates[key] = states[binding]
			} else if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
				// it is a expression, evaluate it
				expression := v[2 : len(v)-2]
				result, err := Eval(expression, states, r.nodeInfo.Ref)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate expression for property %s: %w", key, err)
				}
				resolvedStates[key] = result
			} else {
				// it is a value
				resolvedStates[key] = v
			}
		}
	}
	return resolvedStates, fmt.Errorf("%s [%s] not implemented please override ResolveInput(states) method", r.nodeInfo.Name, r.nodeInfo.Ref)
}

func (r *BaseNode) Validate(input map[string]any) error {
	return fmt.Errorf("%s [%s] not implemented please override Validate(input) method", r.nodeInfo.Name, r.nodeInfo.Ref)
}
