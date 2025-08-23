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
)

// Store provides a unified interface to all application data stores.
// It combines node management, recipe storage, settings management,
// secrets storage, and environment variables into a single convenient interface for the application services.
//
// The embedded concrete stores provide access to workflow nodes, recipe definitions,
// user preferences, secure secrets storage, and environment variables through their respective store implementations.
type Store struct {
	*nodeStore
	*recipesStore
	*settingsStore
	*secretsStore
	*environmentVariablesStore
}

// NewStore creates a new Store with all required data stores initialized.
// It combines node, recipe, settings, secrets, and environment variables stores into a single interface.
//
// The db parameter is used for recipe, settings, secrets, and environment variables storage, while recipesDir
// specifies where recipe YAML files are stored on the file system.
//
// All embedded stores are initialized and ready for use after this call returns.
func NewStore(db *database.Client, recipesDir string) *Store {
	return &Store{
		nodeStore:                 NewNodeStore(),
		recipesStore:              NewRecipesStore(db, recipesDir),
		settingsStore:             NewSettingsStore(db),
		secretsStore:              NewSecretsStore(db),
		environmentVariablesStore: NewEnvironmentVariablesStore(db),
	}
}

