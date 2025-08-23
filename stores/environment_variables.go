package stores

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ironpark/teatime/internal/database"
)

// environmentVariablesStore manages storage and retrieval of environment variables.
// Environment variables are used for non-sensitive configuration data and are stored
// in plaintext in the database.
type environmentVariablesStore struct {
	db *database.Client
}

// NewEnvironmentVariablesStore creates a new environment variables store with the given database client.
func NewEnvironmentVariablesStore(db *database.Client) *environmentVariablesStore {
	return &environmentVariablesStore{db: db}
}

// CreateEnvironmentVariable creates a new environment variable.
func (s *environmentVariablesStore) CreateEnvironmentVariable(name, value, description string) (database.EnvironmentVariable, error) {
	id := uuid.New().String()

	var desc sql.NullString
	if description != "" {
		desc = sql.NullString{String: description, Valid: true}
	}

	envVar, err := s.db.CreateEnvironmentVariable(context.Background(), database.CreateEnvironmentVariableParams{
		ID:          id,
		Name:        name,
		Value:       value,
		Description: desc,
	})
	if err != nil {
		return database.EnvironmentVariable{}, fmt.Errorf("failed to create environment variable: %w", err)
	}

	return envVar, nil
}

// GetEnvironmentVariable retrieves an environment variable by ID.
func (s *environmentVariablesStore) GetEnvironmentVariable(id string) (database.EnvironmentVariable, error) {
	envVar, err := s.db.GetEnvironmentVariable(context.Background(), id)
	if err != nil {
		return database.EnvironmentVariable{}, fmt.Errorf("failed to get environment variable: %w", err)
	}

	return envVar, nil
}

// GetEnvironmentVariableByName retrieves an environment variable by name.
func (s *environmentVariablesStore) GetEnvironmentVariableByName(name string) (database.EnvironmentVariable, error) {
	envVar, err := s.db.GetEnvironmentVariableByName(context.Background(), name)
	if err != nil {
		return database.EnvironmentVariable{}, fmt.Errorf("failed to get environment variable by name: %w", err)
	}

	return envVar, nil
}

// ListEnvironmentVariables returns all environment variables.
func (s *environmentVariablesStore) ListEnvironmentVariables() ([]database.EnvironmentVariable, error) {
	envVars, err := s.db.ListEnvironmentVariables(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list environment variables: %w", err)
	}

	return envVars, nil
}

// UpdateEnvironmentVariable updates an existing environment variable.
func (s *environmentVariablesStore) UpdateEnvironmentVariable(id, name, value, description string) (database.EnvironmentVariable, error) {
	var desc sql.NullString
	if description != "" {
		desc = sql.NullString{String: description, Valid: true}
	}

	envVar, err := s.db.UpdateEnvironmentVariable(context.Background(), database.UpdateEnvironmentVariableParams{
		Name:        name,
		Value:       value,
		Description: desc,
		ID:          id,
	})
	if err != nil {
		return database.EnvironmentVariable{}, fmt.Errorf("failed to update environment variable: %w", err)
	}

	return envVar, nil
}

// DeleteEnvironmentVariable removes an environment variable from storage.
func (s *environmentVariablesStore) DeleteEnvironmentVariable(id string) error {
	err := s.db.DeleteEnvironmentVariable(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to delete environment variable: %w", err)
	}

	return nil
}

// CheckEnvironmentVariableNameExists checks if an environment variable with the given name already exists.
func (s *environmentVariablesStore) CheckEnvironmentVariableNameExists(name string) (bool, error) {
	count, err := s.db.CheckEnvironmentVariableNameExists(context.Background(), name)
	if err != nil {
		return false, fmt.Errorf("failed to check environment variable name: %w", err)
	}

	return count > 0, nil
}

// SearchEnvironmentVariablesByName searches for environment variables matching the given name pattern.
func (s *environmentVariablesStore) SearchEnvironmentVariablesByName(pattern string) ([]database.EnvironmentVariable, error) {
	envVars, err := s.db.SearchEnvironmentVariablesByName(context.Background(), pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search environment variables: %w", err)
	}

	return envVars, nil
}