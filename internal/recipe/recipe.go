package recipe

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Recipe struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Description string     `json:"description"`
	Nodes       []Node     `json:"nodes"`
	Edges       []FlowEdge `json:"edges"`
}

// FlowEdge represents a connection between nodes
type FlowEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type,omitempty"`
}

func Open(path string) (*Recipe, error) {
	recipe := &Recipe{}
	yamlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()

	err = yaml.NewDecoder(yamlFile).Decode(recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func Create(path, name, description string) (*Recipe, error) {
	recipe := &Recipe{
		Path:        path,
		Name:        name,
		Description: description,
		Nodes:       []Node{},
	}
	if err := recipe.Save(); err != nil {
		return nil, err
	}
	return recipe, nil
}

func (r *Recipe) GetNodeById(id string) (Node, error) {
	for _, node := range r.Nodes {
		if node.Id == id {
			return node, nil
		}
	}
	return Node{}, fmt.Errorf("node not found: %s", id)
}

func (r *Recipe) GetConnectedNodes(id string) ([]Node, error) {
	nodes := []Node{}
	for _, edge := range r.Edges {
		if edge.Source == id {
			node, err := r.GetNodeById(edge.Target)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func (r *Recipe) GetNodeDependencies(id string) (ids []string, err error) {
	ids = []string{}
	for _, edge := range r.Edges {
		if edge.Target == id {
			ids = append(ids, edge.Source)
		}
	}
	return ids, nil
}

func (r *Recipe) Save() error {
	yamlFile, err := os.Create(r.Path)
	if err != nil {
		return err
	}
	defer yamlFile.Close()
	err = yaml.NewEncoder(yamlFile, yaml.IndentSequence(true), yaml.CustomMarshaler(func(v any) ([]byte, error) {
		if nodes, ok := v.([]Node); ok {
			return yaml.Marshal(nodes)
		}
		return nil, nil
	})).Encode(r)
	if err != nil {
		return err
	}
	return nil
}
