package recipe

import "github.com/ironpark/teatime/internal/node"

type Node struct {
	Id          string              `json:"id"`
	Ref         string              `json:"ref"`
	Name        string              `json:"name"`
	Icon        string              `json:"icon"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
	Position    Position            `json:"pos"`
	Properties  []node.NodeProperty `json:"properties"`
	Output      []node.NodeProperty `json:"output"`
}

type minifiedNode struct {
	Ref        string         `json:"ref"`
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Position   Position       `json:"pos"`
	Properties map[string]any `json:"properties"`
}

func (n *Node) MarshalYAML() (interface{}, error) {
	properties := make(map[string]any)
	for _, property := range n.Properties {
		properties[property.Key] = property.Value
	}
	return minifiedNode{
		Ref:        n.Ref,
		ID:         n.Id,
		Label:      n.Name,
		Position:   n.Position,
		Properties: properties,
	}, nil
}

func (n *Node) UnmarshalYAML(unmarshal func(interface{}) error) error {
	minified := minifiedNode{}
	err := unmarshal(&minified)
	if err != nil {
		return err
	}
	n.Id = minified.ID
	n.Ref = minified.Ref
	n.Name = minified.Label

	n.Position = minified.Position
	node, err := node.GetNodeByRef(n.Ref)
	if err != nil {
		return err
	}
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
