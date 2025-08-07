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

type recipesStore struct {
	db         *database.Client
	recipesDir string // Directory to store recipe JSON files
}

// NewRecipesStore creates a new recipes store
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

// Sync syncs recipes between database and file system
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

func (r *recipesStore) generateRecipePath(name string, number int) (string, error) {
	normalizedName := strings.TrimSpace(strings.ReplaceAll(name, " ", "-"))
	normalizedName = strings.ToLower(normalizedName)
	normalizedName = norm.NFC.String(normalizedName)
	if number == 0 {
		recipeFileName := fmt.Sprintf("%s.yaml", normalizedName)
		recipePath := filepath.Join(r.recipesDir, recipeFileName)
		if _, err := os.Stat(recipePath); err == nil {
			return r.generateRecipePath(name, number+1)
		}
		return recipePath, nil
	}
	recipeFileName := fmt.Sprintf("%s-%d.yaml", normalizedName, number)
	return filepath.Join(r.recipesDir, recipeFileName), nil
}

// CreateRecipe creates a new recipe
func (r *recipesStore) CreateRecipe(name, description string) (id string, recipe *rc.Recipe, err error) {
	// Generate unique ID for recipe
	recipeID := uuid.New().String()
	// Create recipe file path
	recipePath, err := r.generateRecipePath(name, 0)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get recipe path: %w", err)
	}
	// Save flow to YAML file using recipe.Save
	recipe, err = rc.Create(recipePath, name, description)
	if err != nil {
		return "", nil, fmt.Errorf("failed to save recipe file: %w", err)
	}

	// Save recipe metadata to database
	_, err = r.db.CreateRecipe(context.Background(), database.CreateRecipeParams{
		ID:          recipeID,
		Name:        name,
		Description: description,
		RecipePath:  recipePath,
	})
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(recipePath)
		return "", nil, fmt.Errorf("failed to create recipe in database: %w", err)
	}
	return recipeID, recipe, nil
}

// GetRecipe retrieves a recipe by ID
func (r *recipesStore) GetRecipe(id string) (*rc.Recipe, error) {
	dbRecipe, err := r.db.GetRecipe(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("recipe not found: %w", err)
	}
	flow, err := rc.Open(dbRecipe.RecipePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open recipe file: %w", err)
	}
	return flow, nil
}

// ListRecipes returns all recipes
func (r *recipesStore) ListRecipes() ([]database.Recipe, error) {
	recipes, err := r.db.ListRecipes(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list recipes: %w", err)
	}
	return recipes, nil
}

// UpdateRecipe updates recipe metadata and flow
func (r *recipesStore) UpdateRecipe(id string, recipe *rc.Recipe) error {
	// Update flow file if provided
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

// DeleteRecipe deletes a recipe and its flow file
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

// ExecuteRecipe creates a new execution record for a recipe
func (r *recipesStore) ExecuteRecipe(recipeID string, inputData map[string]interface{}) (*database.Execution, error) {
	panic("not implemented")
}

// UpdateExecutionStatus updates the status of an execution
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

// GetExecutionsByRecipe returns all executions for a recipe
func (r *recipesStore) GetExecutionsByRecipe(recipeID string) ([]database.Execution, error) {
	executions, err := r.db.ListExecutionsByRecipe(context.Background(), recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions: %w", err)
	}
	return executions, nil
}
