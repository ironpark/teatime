// Package stores provides data access and management functionality for the Teatime workflow automation system.
// It manages workflow node types, recipe storage, application settings, and execution records.
//
// The package includes three main store types:
//   - NodeStore: manages workflow node definitions and creates node instances
//   - RecipesStore: handles recipe persistence across database and YAML files
//   - SettingsStore: manages user preferences and application configuration
//
// All stores provide consistent interfaces for CRUD operations and maintain data integrity
// between different storage backends (SQLite database and file system).
package stores

import (
	"github.com/ironpark/teatime/internal/database"
	node2 "github.com/ironpark/teatime/internal/node"
	rc "github.com/ironpark/teatime/internal/recipe"
)

// Store provides a unified interface to all application data stores.
// It combines node management, recipe storage, and settings management
// into a single convenient interface for the application services.
//
// The embedded interfaces provide access to workflow nodes, recipe definitions,
// and user preferences through their respective store implementations.
type Store struct {
	NodeStore
	RecipesStore
	SettingsStore
}

// NewStore creates a new Store with all required data stores initialized.
// It combines node, recipe, and settings stores into a single interface.
//
// The db parameter is used for recipe and settings storage, while recipesDir
// specifies where recipe YAML files are stored on the file system.
//
// All embedded stores are initialized and ready for use after this call returns.
func NewStore(db *database.Client, recipesDir string) *Store {
	return &Store{
		NodeStore:     NewNodeStore(),
		RecipesStore:  NewRecipesStore(db, recipesDir),
		SettingsStore: NewSettingsStore(db),
	}
}

// NodeStore defines the interface for managing workflow node types and instances.
// It provides access to registered node definitions and creates node instances
// for use in workflows.
type NodeStore interface {
	// GetNodeInfos returns information about all registered node types.
	GetNodeInfos() []node2.NodeInfo

	// GetNodeInfosByType returns nodes filtered by the specified type category.
	GetNodeInfosByType(nodeType string) []node2.NodeInfo

	// GetNodeInfo retrieves information for a specific node type by ID.
	GetNodeInfo(id string) node2.NodeInfo

	// CreateNode creates a new node instance with a unique identifier.
	CreateNode(nodeId string) Node
}

// RecipesStore defines the interface for managing recipes and their executions.
// It handles both recipe definitions (stored as YAML files) and execution records
// (stored in the database).
type RecipesStore interface {
	// Sync synchronizes recipe metadata between database and file system.
	Sync() error

	// CreateRecipe creates a new recipe with the given name and description.
	CreateRecipe(name, description string) (id string, recipe *rc.Recipe, err error)

	// GetRecipe retrieves a recipe by its ID.
	GetRecipe(id string) (*rc.Recipe, error)

	// ListRecipes returns all recipe metadata from the database.
	ListRecipes() ([]database.Recipe, error)

	// UpdateRecipe updates both recipe file and database metadata.
	UpdateRecipe(id string, recipe *rc.Recipe) error

	// DeleteRecipe removes a recipe from both database and file system.
	DeleteRecipe(id string) error

	// ExecuteRecipe creates a new execution record for a recipe.
	ExecuteRecipe(recipeID string, inputData map[string]interface{}) (*database.Execution, error)

	// UpdateExecutionStatus updates the status and results of an execution.
	UpdateExecutionStatus(executionID, status string, outputData map[string]interface{}, errorMsg string) error

	// GetExecutionsByRecipe retrieves all execution records for a recipe.
	GetExecutionsByRecipe(recipeID string) ([]database.Execution, error)
}

// SettingsStore defines the interface for managing application settings.
// It provides access to user preferences including theme, language, and
// startup behavior settings.
type SettingsStore interface {
	// GetSettings retrieves current application settings, creating defaults if needed.
	GetSettings() (database.Setting, error)

	// UpdateSettings updates all settings with the provided values.
	UpdateSettings(settings database.Setting) error

	// UpdateTheme updates only the theme setting.
	UpdateTheme(theme string) error

	// UpdateAutoStart updates only the auto-start setting.
	UpdateAutoStart(autoStart bool) error

	// UpdateLanguage updates only the language setting.
	UpdateLanguage(language string) error
}
