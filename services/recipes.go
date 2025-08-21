package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ironpark/teatime/internal/database"
	"github.com/ironpark/teatime/internal/node"
	rc "github.com/ironpark/teatime/internal/recipe"
	"github.com/ironpark/teatime/internal/runner"
	"github.com/ironpark/teatime/stores"
	"github.com/samber/lo"
)

type RecipesService struct {
	store *stores.Store
}

type RecipeInfo struct {
	database.Recipe
	ExecutionCount      int
	LastExecution       time.Time
	LastExecutionStatus string
	Tags                []string
	NodeTypes           []string
}

func NewRecipesService(store *stores.Store) *RecipesService {
	return &RecipesService{store: store}
}

type CreatedRecipe struct {
	Recipe *rc.Recipe
	ID     string
}

func (s *RecipesService) GetAvailableNodes() []node.NodeInfo {
	return lo.Map(node.GetNodes(), func(n node.Node, _ int) node.NodeInfo {
		return n.Info()
	})
}

func (s *RecipesService) GetAvailableNodesByType(nodeType string) []node.NodeInfo {
	return lo.Map(node.GetNodesByType(node.NodeType(nodeType)), func(n node.Node, _ int) node.NodeInfo {
		return n.Info()
	})
}

func (s *RecipesService) CreateRecipe(name, description string) (*CreatedRecipe, error) {
	id, recipe, err := s.store.CreateRecipe(name, description)
	if err != nil {
		return nil, err
	}
	return &CreatedRecipe{
		Recipe: recipe,
		ID:     id,
	}, nil
}

func (s *RecipesService) SaveRecipe(id string, recipe *rc.Recipe) (*rc.Recipe, error) {
	if err := recipe.Save(); err != nil {
		return nil, err
	}
	s.store.UpdateRecipe(id, recipe)
	return recipe, nil
}

func (s *RecipesService) CreateNode(ref string, x, y int) (rc.Node, error) {
	createdNode, err := node.GetNodeByRef(ref)
	if err != nil {
		return rc.Node{}, err
	}
	current := time.Now().UnixNano()
	currentHex := fmt.Sprintf("%x", current)
	return rc.Node{
		Id:       currentHex,
		Position: rc.Position{x, y},
		Type:     string(createdNode.Type()),
		NodeData: rc.NodeData{
			Ref:           ref,
			Icon:          createdNode.Icon(),
			Label:         createdNode.Name(),
			Name:          createdNode.Name(),
			NodeType:      string(createdNode.Type()),
			Description:   createdNode.Description(),
			Properties:    createdNode.GetProperties(node.PropertyContext{}),
			Outputs:       createdNode.GetOutput(node.PropertyContext{}),
			OutputHandles: createdNode.GetOutputHandles(node.PropertyContext{}),
		},
	}, nil
}

func (s *RecipesService) RunRecipe(recipe *rc.Recipe, startNodeId string, properties map[string]any) error {
	return runner.Run(context.Background(), recipe, startNodeId, properties, func(recipe *rc.Recipe, state runner.NodeExecutionStatus, node rc.Node, output map[string]any, err error) {
		fmt.Println("Recipe", recipe.Name, "Node", node.Id, "State", state, "Output", output, "Error", err)
	})
}

func (s *RecipesService) RunRecipeByID(id string, startNodeId string, properties map[string]any) error {
	recipe, err := s.store.GetRecipe(id)
	if err != nil {
		return err
	}
	triggerNode, err := recipe.GetNodeById(startNodeId)
	if err != nil {
		return err
	}
	if triggerNode.Type != string(node.NodeTypeTrigger) {
		return fmt.Errorf("node %s is not a trigger node", startNodeId)
	}
	return runner.Run(context.Background(), recipe, startNodeId, properties, func(recipe *rc.Recipe, state runner.NodeExecutionStatus, node rc.Node, output map[string]any, err error) {
		fmt.Println("Recipe", recipe.Name, "Node", node.Id, "State", state, "Output", output, "Error", err)
	})
}

func (s *RecipesService) GetRecipe(id string) (*rc.Recipe, error) {
	return s.store.GetRecipe(id)
}

func (s *RecipesService) ListRecipes() ([]RecipeInfo, error) {
	recipes, err := s.store.ListRecipes()
	if err != nil {
		return nil, err
	}
	recipeInfos := make([]RecipeInfo, len(recipes))
	for i, recipe := range recipes {
		info := RecipeInfo{
			Recipe: recipe,
		}
		if recipe.NodeTypes != "" {
			json.Unmarshal([]byte(recipe.NodeTypes), &info.NodeTypes)
			fmt.Println("NodeTypes", info.NodeTypes)
		}
		if recipe.Tags != "" {
			json.Unmarshal([]byte(recipe.Tags), &info.Tags)
			fmt.Println("Tags", info.Tags)
		}
		recipeInfos[i] = info
	}
	return recipeInfos, nil
}

func (s *RecipesService) UpdateRecipe(id string, recipe *rc.Recipe) error {
	return s.store.UpdateRecipe(id, recipe)
}

func (s *RecipesService) DeleteRecipe(id string) error {
	return s.store.DeleteRecipe(id)
}

func (s *RecipesService) Sync() error {
	return s.store.Sync()
}

// ExecuteTriggerNode executes a single trigger node with provided arguments
func (s *RecipesService) ExecuteTriggerNode(recipeID, nodeID string, properties map[string]any, args map[string]any) error {
	recipe, err := s.store.GetRecipe(recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	triggerNode, err := recipe.GetNodeById(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}

	if triggerNode.Type != string(node.NodeTypeTrigger) {
		return fmt.Errorf("node %s is not a trigger node", nodeID)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "args", args)
	fmt.Println("args", args, "properties", properties)
	// Execute the trigger node
	return runner.Run(ctx, recipe, nodeID, properties, func(recipe *rc.Recipe, state runner.NodeExecutionStatus, node rc.Node, output map[string]any, err error) {
		fmt.Printf("Recipe: %s, Node: %s, State: %s, Output: %v, Error: %v\n", recipe.Name, node.Id, state, output, err)
	})
}
