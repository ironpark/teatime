// Package recipe provides node data types used by UI and filesystem.
// It minimizes node data when saving and expands minimized data when loading.
package recipe

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/ironpark/teatime/internal/node"
)

// NodeData contains the metadata and configuration for a workflow node.
// It includes all the information needed to display and execute the node.
type NodeData struct {
	Ref           string              `json:"ref"`
	Icon          string              `json:"icon"`
	Label         string              `json:"label"`
	Name          string              `json:"name"`
	NodeType      string              `json:"nodeType"`
	Description   string              `json:"description"`
	Properties    []node.NodeProperty `json:"properties"`
	Outputs       []node.NodeProperty `json:"outputs"`
	OutputHandles []node.OutputHandle `json:"outputHandles"`
}

// Node represents a workflow node instance with position and runtime data.
// It combines NodeData with instance-specific information like ID and position.
type Node struct {
	Id       string   `json:"id"`
	Type     string   `json:"type"`
	Position Position `json:"position"`
	NodeData `json:"data"`
	rawNode  node.Node `json:"-"`
}

type minifiedNode struct {
	Ref        string         `json:"ref"`
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Position   Position       `json:"pos"`
	Properties map[string]any `json:"properties"`
}

// MarshalYAML converts the node to a minified YAML representation.
// It only includes essential data to reduce file size.
func (n Node) MarshalYAML() ([]byte, error) {
	properties := make(map[string]any)
	for _, property := range n.Properties {
		properties[property.Key] = property.Value
	}
	return yaml.Marshal(minifiedNode{
		Ref:        n.Ref,
		ID:         n.Id,
		Label:      n.Label,
		Position:   n.Position,
		Properties: properties,
	})
}

func (n *Node) IsTrigger() bool {
	return n.NodeType == string(node.NodeTypeTrigger)
}

// GetRawNode returns the underlying node implementation.
// This provides access to the node's execution logic.
func (n *Node) GetRawNode() node.Node {
	return n.rawNode
}

// UnmarshalJSON deserializes JSON data into a Node.
// It handles both minified and full node representations.
func (n *Node) UnmarshalJSON(data []byte) error {
	// Parse JSON to detect format type
	var temp map[string]any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Check for minified format (has "pos" field)
	if _, hasPos := temp["pos"]; hasPos {
		minified := minifiedNode{}
		if err := json.Unmarshal(data, &minified); err != nil {
			return err
		}
		return n.unmarshalFromMinified(minified)
	}

	// Check for full format (has "position" field)
	if _, hasPosition := temp["position"]; hasPosition {
		type fullNode struct {
			Id       string   `json:"id"`
			Type     string   `json:"type"`
			Position Position `json:"position"`
			NodeData `json:"data"`
		}
		full := fullNode{}
		if err := json.Unmarshal(data, &full); err != nil {
			return err
		}
		n.Id = full.Id
		n.Type = full.Type
		n.Position = full.Position
		n.NodeData = full.NodeData
		rawNode, err := node.GetNodeByRef(full.Ref)
		if err != nil {
			return err
		}
		n.rawNode = rawNode
		return nil
	}

	return fmt.Errorf("unknown node format: missing both 'pos' and 'position' fields")
}

// UnmarshalYAML deserializes YAML data into a Node.
// It expects minified node format from YAML files.
func (n *Node) UnmarshalYAML(data []byte) error {
	minified := minifiedNode{}
	err := yaml.Unmarshal(data, &minified)
	if err != nil {
		return err
	}
	return n.unmarshalFromMinified(minified)
}

// unmarshalFromMinified populates a Node from its minified representation.
// It retrieves the node definition, resolves properties dynamically, and applies saved values.
func (n *Node) unmarshalFromMinified(minified minifiedNode) error {
	n.Id = minified.ID
	n.Ref = minified.Ref
	n.Label = minified.Label
	n.Position = minified.Position
	rawNode, err := node.GetNodeByRef(n.Ref)
	if err != nil {
		return err
	}
	n.Name = rawNode.Name()
	n.Icon = rawNode.Icon()
	n.Type = string(rawNode.Type())
	n.Description = rawNode.Description()
	n.Properties = rawNode.GetProperties(node.PropertyContext(minified.Properties))
	n.Outputs = rawNode.GetOutput(node.PropertyContext(minified.Properties))
	n.OutputHandles = rawNode.GetOutputHandles(node.PropertyContext(minified.Properties))
	n.rawNode = rawNode
	for i, property := range n.Properties {
		if value, ok := minified.Properties[property.Key]; ok {
			n.Properties[i].Value = value
		}
	}
	return nil
}
