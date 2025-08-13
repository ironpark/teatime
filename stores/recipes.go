package stores

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironpark/teatime/internal/database"
	rc "github.com/ironpark/teatime/internal/recipe"
	"golang.org/x/text/unicode/norm"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// recipesStore manages recipe data storage and synchronization between database and file system.
// It maintains recipe metadata in a SQLite database while storing recipe definitions as YAML files.
//
// The store is not safe for concurrent use. Callers must provide external synchronization
// when accessing the store from multiple goroutines.
type recipesStore struct {
	db         *database.Client
	recipesDir string // Directory to store recipe YAML files
}

// NewRecipesStore creates a new recipes store with the given database client and recipes directory.
// It ensures the recipes directory exists and performs an initial synchronization between
// the database and file system. If directory creation or synchronization fails, warnings
// are logged but the store is still returned.
//
// The recipesDir parameter specifies where recipe YAML files will be stored.
func NewRecipesStore(db *database.Client, recipesDir string) *recipesStore {
	// Ensure recipes directory exists
	if err := os.MkdirAll(recipesDir, 0755); err != nil {
		// Log error but continue - will fail on actual operations
		fmt.Printf("Warning: Failed to create recipes directory: %v\n", err)
	}
	store := &recipesStore{
		db:         db,
		recipesDir: recipesDir,
	}
	if err := store.Sync(); err != nil {
		fmt.Printf("Warning: Failed to sync recipes: %v\n", err)
	}
	return store
}

// Sync synchronizes recipe metadata between the database and file system.
// It scans the recipes directory for YAML files, creates database entries for new recipes,
// updates existing entries, and removes database entries for deleted files.
//
// Returns an error if the recipes directory cannot be read or if database operations fail.
// File parsing errors for individual recipes are silently ignored.
func (r *recipesStore) Sync() error {
	// Get all recipes from database
	dbRecipes, err := r.db.ListRecipes(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list recipes: %w", err)
	}

	existDatabaseRecipes := lo.KeyBy(dbRecipes, func(dbRecipe database.Recipe) string {
		return dbRecipe.RecipePath
	})

	needUpdate := make(map[string]*rc.Recipe)
	needCreate := make(map[string]*rc.Recipe)

	// Get all recipes from file system
	files, err := os.ReadDir(r.recipesDir)
	if err != nil {
		return fmt.Errorf("failed to read recipes directory: %w", err)
	}
	recipes := make(map[string]*rc.Recipe)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}
		path := filepath.Join(r.recipesDir, file.Name())
		recipe, err := rc.Open(path)
		if err != nil {
			continue
		}
		recipes[path] = recipe
		if _, ok := existDatabaseRecipes[path]; !ok {
			needCreate[path] = recipe
		} else {
			needUpdate[path] = recipe
		}
	}

	// Sync database and file system
	r.db.WithTx(context.Background(), func(ctx context.Context, queries *database.Queries) error {
		for path, recipe := range needCreate {
			_, err := queries.CreateRecipe(ctx, database.CreateRecipeParams{
				ID:          uuid.New().String(),
				Name:        recipe.Name,
				Description: recipe.Description,
				RecipePath:  path,
			})
			if err != nil {
				return fmt.Errorf("failed to create recipe: %w", err)
			}
		}
		for path, recipe := range needUpdate {
			_, err := queries.UpdateRecipeByPath(ctx, database.UpdateRecipeByPathParams{
				Name:        recipe.Name,
				Description: recipe.Description,
				RecipePath:  path,
			})
			if err != nil {
				return fmt.Errorf("failed to update recipe: %w", err)
			}
		}
		// Delete recipes that are not in the file system
		for _, dbRecipe := range dbRecipes {
			if _, ok := recipes[dbRecipe.RecipePath]; ok {
				continue
			}
			if err := queries.DeleteRecipe(ctx, dbRecipe.ID); err != nil {
				return fmt.Errorf("failed to delete recipe: %w", err)
			}
		}
		return nil
	})
	return nil
}

// generateRecipePath creates a unique file path for a recipe YAML file within the specified recipes directory.
// It normalizes the recipe name by trimming whitespace, replacing spaces with hyphens, converting to lowercase,
// and applying Unicode NFC normalization for consistent character representation.
//
// The function implements collision avoidance by recursively checking for existing files:
// - If number is 0, it attempts to create a path without a numeric suffix (e.g., "my-recipe.yaml")
// - If that file exists, it recursively calls itself with number+1
// - For number > 0, it appends a numeric suffix (e.g., "my-recipe-1.yaml", "my-recipe-2.yaml")
func generateRecipePath(recipesDir, name string, number int) (string, error) {
	normalizedName := strings.TrimSpace(strings.ReplaceAll(name, " ", "-"))
	normalizedName = strings.ToLower(normalizedName)
	normalizedName = norm.NFC.String(normalizedName)
	if number == 0 {
		recipeFileName := fmt.Sprintf("%s.yaml", normalizedName)
		recipePath := filepath.Join(recipesDir, recipeFileName)
		if _, err := os.Stat(recipePath); err == nil {
			return generateRecipePath(recipesDir, name, number+1)
		}
		return recipePath, nil
	}
	recipeFileName := fmt.Sprintf("%s-%d.yaml", normalizedName, number)
	return filepath.Join(recipesDir, recipeFileName), nil
}

// CreateRecipe creates a new recipe with the given name and description.
// It generates a unique file path, saves the recipe to both database and file system,
// and returns the generated recipe ID and the created recipe object.
//
// The name parameter is used to generate the YAML filename after normalization.
// If a file with the same name exists, a numeric suffix is automatically added.
//
// Returns an error if path generation fails, database insertion fails, or file creation fails.
// On database failure, any created files are automatically cleaned up.
func (r *recipesStore) CreateRecipe(name, description string) (id string, recipe *rc.Recipe, err error) {
	// Generate unique ID for recipe
	recipeID := uuid.New().String()
	// Create recipe file path
	recipePath, err := generateRecipePath(r.recipesDir, name, 0)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate recipe path: %w", err)
	}

	// Save recipe metadata to database and file system
	r.db.WithTx(context.Background(), func(ctx context.Context, queries *database.Queries) error {
		_, err = queries.CreateRecipe(ctx, database.CreateRecipeParams{
			ID:          recipeID,
			Name:        name,
			Description: description,
			RecipePath:  recipePath,
		})
		if err != nil {
			// Clean up file if database insert fails
			os.Remove(recipePath)
			return fmt.Errorf("failed to create recipe in database: %w", err)
		}
		// Save recipe to YAML file
		recipe, err = rc.Create(recipePath, name, description)
		if err != nil {
			return fmt.Errorf("failed to save recipe file: %w", err)
		}
		return nil
	})
	return recipeID, recipe, nil
}

// GetRecipe retrieves a recipe by its ID and returns the recipe object.
// It first looks up the recipe metadata in the database, then loads the recipe
// definition from the corresponding YAML file.
//
// Returns an error if the recipe ID is not found in the database or if the
// recipe file cannot be opened or parsed.
func (r *recipesStore) GetRecipe(id string) (*rc.Recipe, error) {
	dbRecipe, err := r.db.GetRecipe(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}
	flow, err := rc.Open(dbRecipe.RecipePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open recipe: %w", err)
	}
	return flow, nil
}

// ListRecipes returns all recipe metadata from the database.
// The returned slice contains database.Recipe objects with ID, name, description,
// and file path information, but not the full recipe definitions.
//
// Returns an error if the database query fails.
func (r *recipesStore) ListRecipes() ([]database.Recipe, error) {
	recipes, err := r.db.ListRecipes(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list recipes: %w", err)
	}
	return recipes, nil
}

// UpdateRecipe updates both the recipe file and database metadata.
// If the recipe parameter is not nil, it saves the recipe to its YAML file
// before updating the database with the recipe's name and description.
//
// Returns an error if file saving fails or if the database update fails.
// The recipe parameter must not be nil.
func (r *recipesStore) UpdateRecipe(id string, recipe *rc.Recipe) error {
	// Update recipe file if provided
	if recipe != nil {
		if err := recipe.Save(); err != nil {
			return fmt.Errorf("failed to save recipe file: %w", err)
		}
	}

	if _, err := r.db.UpdateRecipe(context.Background(), database.UpdateRecipeParams{
		ID:          id,
		Name:        recipe.Name,
		Description: recipe.Description,
	}); err != nil {
		return fmt.Errorf("failed to update recipe in database: %w", err)
	}

	return nil
}

// DeleteRecipe removes a recipe from both the database and file system.
// It first retrieves the recipe metadata to get the file path, removes the
// database entry, then deletes the corresponding YAML file.
//
// Returns an error if the recipe ID is not found, if database deletion fails,
// or if file deletion fails (except when the file doesn't exist).
// The operation is not atomic - database deletion may succeed while file deletion fails.
func (r *recipesStore) DeleteRecipe(id string) error {
	// Get recipe to get file path
	dbRecipe, err := r.db.GetRecipe(context.Background(), id)
	if err != nil {
		return err
	}

	// Delete from database first
	if err := r.db.DeleteRecipe(context.Background(), id); err != nil {
		return fmt.Errorf("failed to delete recipe from database: %w", err)
	}

	// Delete recipe file if exists
	if err := os.Remove(dbRecipe.RecipePath); err != nil {
		if !os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete recipe file: %w", err)
	}

	return nil
}

// ExecuteRecipe creates a new execution record for a recipe with the given input data.
// This method is currently not implemented and will panic when called.
//
// TODO: Implement recipe execution functionality.
func (r *recipesStore) ExecuteRecipe(recipeID string, inputData map[string]interface{}) (*database.Execution, error) {
	panic("not implemented")
}

// UpdateExecutionStatus updates the status, output data, and error message of an execution.
// The outputData is marshaled to JSON before being stored in the database.
// If outputData is nil, an empty string is stored.
//
// Returns an error if JSON marshaling fails or if the database update fails.
func (r *recipesStore) UpdateExecutionStatus(executionID, status string, outputData map[string]interface{}, errorMsg string) error {
	var outputJSON string
	if outputData != nil {
		data, err := json.Marshal(outputData)
		if err != nil {
			return fmt.Errorf("failed to marshal output data: %w", err)
		}
		outputJSON = string(data)
	}

	params := database.UpdateExecutionStatusParams{
		Status:       status,
		ErrorMessage: errorMsg,
		OutputData:   outputJSON,
		ID:           executionID,
	}

	if _, err := r.db.UpdateExecutionStatus(context.Background(), params); err != nil {
		return fmt.Errorf("failed to update execution status: %w", err)
	}

	return nil
}

// GetExecutionsByRecipe retrieves all execution records for the specified recipe ID.
// Returns a slice of database.Execution objects containing execution history,
// status, input/output data, and error messages.
//
// Returns an error if the database query fails.
func (r *recipesStore) GetExecutionsByRecipe(recipeID string) ([]database.Execution, error) {
	executions, err := r.db.ListExecutionsByRecipe(context.Background(), recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions: %w", err)
	}
	return executions, nil
}
