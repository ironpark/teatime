package recipe

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/ironpark/teatime/internal/node"
)

type NodeData struct {
	Ref         string              `json:"ref"`
	Icon        string              `json:"icon"`
	Label       string              `json:"label"`
	Name        string              `json:"name"`
	NodeType    string              `json:"nodeType"`
	Description string              `json:"description"`
	Properties  []node.NodeProperty `json:"properties"`
	Output      []node.NodeProperty `json:"output"`
}

type Node struct {
	Id       string   `json:"id"`
	Type     string   `json:"type"`
	Position Position `json:"position"`
	NodeData `json:"data"`
}

type minifiedNode struct {
	Ref        string         `json:"ref"`
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Position   Position       `json:"pos"`
	Properties map[string]any `json:"properties"`
}

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

func (n *Node) UnmarshalYAML(data []byte) error {
	fmt.Println(string(data))
	minified := minifiedNode{}
	err := yaml.Unmarshal(data, &minified)
	if err != nil {
		return err
	}
	n.Id = minified.ID
	n.Ref = minified.Ref
	n.Label = minified.Label
	n.Position = minified.Position
	node, err := node.GetNodeByRef(n.Ref)
	if err != nil {
		return err
	}
	n.Name = node.Name()
	n.Icon = node.Icon()
	n.Type = string(node.Type())
	n.Description = node.Description()
	n.Properties = node.Properties()
	n.Output = node.Output()

	for i, property := range n.Properties {
		if value, ok := minified.Properties[property.Key]; ok {
			n.Properties[i].Value = value
		}
	}
	return nil
}
