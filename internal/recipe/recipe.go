// Package recipe provides workflow definition and node data structures for the Teatime automation platform.
//
// This package handles the complete lifecycle of workflow recipes, including:
//   - Loading and saving recipes from/to YAML files
//   - Managing workflow nodes with their properties and connections
//   - Serializing nodes in both full and minified formats for storage efficiency
//   - Providing graph operations like dependency resolution and connected node traversal
//
// The main types are Recipe (workflow definition), Node (workflow step), and FlowEdge (node connection).
// Nodes support dual serialization formats: full format for runtime and minified format for persistence.
package recipe

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// Recipe represents a complete workflow definition with nodes and connections.
// It contains all the information needed to execute a workflow.
type Recipe struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Description string     `json:"description"`
	Nodes       []Node     `json:"nodes"`
	Edges       []FlowEdge `json:"edges"`
}

// FlowEdge represents a connection between two nodes in the workflow.
// It defines the data flow from source node to target node through specific handles.
type FlowEdge struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	SourceHandle string `json:"sourceHandle"`         // Output handle ID from source node
	TargetHandle string `json:"targetHandle"`         // Input handle ID to target node (optional)
	Type         string `json:"type,omitempty"`
}

// Open loads a recipe from a YAML file at the specified path.
// It returns an error if the file cannot be read or parsed.
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

// Create creates a new recipe with the given path, name, and description.
// It saves the recipe to disk and returns the created recipe instance.
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

// GetNodeById finds and returns a node by its ID.
// It returns an error if no node with the specified ID is found.
func (r *Recipe) GetNodeById(id string) (Node, error) {
	for _, node := range r.Nodes {
		if node.Id == id {
			return node, nil
		}
	}
	return Node{}, fmt.Errorf("node not found: %s", id)
}

// GetConnectedNodes returns all nodes that are connected as targets from the specified source node.
// It follows outgoing edges from the given node ID.
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

// GetConnectedNodesByHandles returns nodes connected through specific output handles.
// It only follows edges that match the specified output handle IDs.
// If outputHandles is empty or contains "default", it returns all connected nodes for backward compatibility.
func (r *Recipe) GetConnectedNodesByHandles(sourceId string, outputHandles []string) ([]Node, error) {
	// Handle empty or default case
	if len(outputHandles) == 0 || (len(outputHandles) == 1 && outputHandles[0] == "default") {
		return r.GetConnectedNodes(sourceId)
	}
	
	// Create a map for quick handle lookup
	handleMap := make(map[string]bool)
	for _, handle := range outputHandles {
		handleMap[handle] = true
	}
	
	nodes := []Node{}
	for _, edge := range r.Edges {
		if edge.Source == sourceId {
			// If edge has no sourceHandle specified, treat as "default"
			edgeHandle := edge.SourceHandle
			if edgeHandle == "" {
				edgeHandle = "default"
			}
			
			// Check if this edge's handle is in the active handles
			if handleMap[edgeHandle] {
				node, err := r.GetNodeById(edge.Target)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, node)
			}
		}
	}
	return nodes, nil
}

// GetNodeDependencies returns the IDs of all nodes that the specified node depends on.
// It follows incoming edges to the given node ID.
func (r *Recipe) GetNodeDependencies(id string) (ids []string, err error) {
	ids = []string{}
	for _, edge := range r.Edges {
		if edge.Target == id {
			ids = append(ids, edge.Source)
		}
	}
	return ids, nil
}

// Save persists the recipe to its associated file path in YAML format.
// It creates or overwrites the file with the current recipe data.
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
