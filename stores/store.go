package stores

import (
	"github.com/ironpark/teatime/internal/database"
	node2 "github.com/ironpark/teatime/internal/node"
	rc "github.com/ironpark/teatime/internal/recipe"
)

type Store struct {
	NodeStore
	RecipesStore
	SettingsStore
}

func NewStore(db *database.Client, recipesDir string) *Store {
	return &Store{
		NodeStore:     NewNodeStore(),
		RecipesStore:  NewRecipesStore(db, recipesDir),
		SettingsStore: NewSettingsStore(db),
	}
}

type NodeStore interface {
	GetNodeInfos() []node2.NodeInfo
	GetNodeInfosByType(nodeType string) []node2.NodeInfo
	GetNodeInfo(id string) node2.NodeInfo
	CreateNode(nodeId string) Node
}

type RecipesStore interface {
	Sync() error
	CreateRecipe(name, description string) (id string, recipe *rc.Recipe, err error)
	GetRecipe(id string) (*rc.Recipe, error)
	ListRecipes() ([]database.Recipe, error)
	UpdateRecipe(id string, recipe *rc.Recipe) error
	DeleteRecipe(id string) error
	ExecuteRecipe(recipeID string, inputData map[string]interface{}) (*database.Execution, error)
	UpdateExecutionStatus(executionID, status string, outputData map[string]interface{}, errorMsg string) error
	GetExecutionsByRecipe(recipeID string) ([]database.Execution, error)
}

type SettingsStore interface {
	GetSettings() (database.Setting, error)
	UpdateSettings(settings database.Setting) error
	UpdateTheme(theme string) error
	UpdateAutoStart(autoStart bool) error
	UpdateLanguage(language string) error
}
