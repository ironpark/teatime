package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ironpark/teatime/internal/database"
	"github.com/ironpark/teatime/internal/node"
	rc "github.com/ironpark/teatime/internal/recipe"
	"github.com/ironpark/teatime/internal/runner"
	"github.com/ironpark/teatime/stores"
	"github.com/samber/lo"
)

type RecipesService struct {
	store           *stores.Store
	editSessions    map[string]*EditSession
	editSessionLock sync.RWMutex
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
	return &RecipesService{
		store:        store,
		editSessions: make(map[string]*EditSession),
	}
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

func (s *RecipesService) getSession(id string) *EditSession {
	s.editSessionLock.Lock()
	defer s.editSessionLock.Unlock()
	session, ok := s.editSessions[id]
	if !ok {
		recipe, err := s.store.GetRecipe(id)
		if err != nil {
			return nil
		}
		session = newSession(recipe)
		s.editSessions[id] = session
	}
	return session
}

func (s *RecipesService) GetRecipe(id string) *EditSession {
	return s.getSession(id)
}

func (s *RecipesService) UpdateRecipe(id string, recipe *rc.Recipe) error {
	return s.getSession(id).UpdateRecipe(recipe)
}

func (s *RecipesService) DeleteRecipe(id string) error {
	return s.store.DeleteRecipe(id)
}

func (s *RecipesService) SaveRecipe(id string) (*EditSession, error) {
	session := s.getSession(id)
	if err := session.Save(); err != nil {
		return nil, err
	}
	s.store.UpdateRecipe(id, session.Recipe)
	return session, nil
}

func (s *RecipesService) CreateNode(id string, ref string, x, y int) (rc.Node, error) {
	session := s.getSession(id)
	return session.CreateNode(ref, x, y)
}

func (s *RecipesService) UpdateNode(id string, nodeId string, x, y int, label string, properties map[string]any) (rc.Node, error) {
	session := s.getSession(id)
	return session.UpdateNode(nodeId, x, y, label, properties)
}

func (s *RecipesService) DeleteNode(id string, nodeId string) error {
	session := s.getSession(id)
	return session.DeleteNode(nodeId)
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

func (s *RecipesService) Sync() error {
	return s.store.Sync()
}

func (s *RecipesService) RunRecipe(recipe *rc.Recipe, startNodeId string, properties map[string]any) error {
	// TODO: implement recipe runner
	return nil
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
	// TODO: implement recipe runner
	return nil
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
	states := runner.NewWorkflowState()
	states.SetExecContext(args)
	err = runner.Run(ctx, recipe, nodeID, states, properties, func(recipe *rc.Recipe, state runner.NodeExecutionStatus, node rc.Node, output map[string]any, err error) {
		fmt.Println("node", node.Id, "output", output, "error", err)
	})
	if err != nil {
		return fmt.Errorf("failed to run recipe: %w", err)
	}
	return nil
}
