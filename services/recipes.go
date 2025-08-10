package services

import (
	"fmt"
	"time"

	"github.com/ironpark/teatime/internal/database"
	"github.com/ironpark/teatime/internal/node"
	rc "github.com/ironpark/teatime/internal/recipe"
	"github.com/ironpark/teatime/stores"
	"github.com/samber/lo"
)

type RecipesService struct {
	store stores.RecipesStore
}

func NewRecipesService(store stores.RecipesStore) *RecipesService {
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

func (s *RecipesService) SaveRecipe(recipe *rc.Recipe) (*rc.Recipe, error) {
	if err := recipe.Save(); err != nil {
		return nil, err
	}
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
		Id:          currentHex,
		Ref:         ref,
		Position:    rc.Position{x, y},
		Icon:        createdNode.Icon(),
		Properties:  createdNode.Properties(),
		Output:      createdNode.Output(),
		Name:        createdNode.Name(),
		Description: createdNode.Description(),
		Type:        string(createdNode.Type()),
	}, nil
}

func (s *RecipesService) GetRecipe(id string) (*rc.Recipe, error) {
	return s.store.GetRecipe(id)
}

func (s *RecipesService) ListRecipes() ([]database.Recipe, error) {
	return s.store.ListRecipes()
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
