package recipe

import (
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

type Position [2]int

func Open(path string) (*Recipe, error) {
	recipe := &Recipe{}
	yamlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()

	err = yaml.NewDecoder(yamlFile).Decode(&recipe)
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

func (r *Recipe) Save() error {
	yamlFile, err := os.Create(r.Path)
	if err != nil {
		return err
	}
	defer yamlFile.Close()
	err = yaml.NewEncoder(yamlFile).Encode(r)
	if err != nil {
		return err
	}
	return nil
}
