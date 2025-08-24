package services

import (
	"fmt"
	"regexp"

	"github.com/ironpark/teatime/internal/database"
	"github.com/ironpark/teatime/stores"
)

// EnvironmentVariablesService provides environment variable management functionality to the frontend.
// It handles CRUD operations for non-sensitive configuration data stored in plaintext.
//
// Environment variables are used for configuration data that doesn't require encryption,
// such as API base URLs, timeouts, and other application settings.
type EnvironmentVariablesService struct {
	store *stores.Store
}

// EnvironmentVariableInfo represents basic environment variable information.
type EnvironmentVariableInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// EnvironmentVariableCreateRequest represents the data needed to create a new environment variable.
type EnvironmentVariableCreateRequest struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// EnvironmentVariableUpdateRequest represents the data needed to update an existing environment variable.
type EnvironmentVariableUpdateRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// NewEnvironmentVariablesService creates a new environment variables service with the given store.
func NewEnvironmentVariablesService(store *stores.Store) *EnvironmentVariablesService {
	return &EnvironmentVariablesService{store: store}
}

// CreateEnvironmentVariable creates a new environment variable with the provided data.
//
// Returns an error if:
// - An environment variable with the same name already exists
// - Database operation fails
func (s *EnvironmentVariablesService) CreateEnvironmentVariable(req EnvironmentVariableCreateRequest) (*EnvironmentVariableInfo, error) {
	// Validate environment variable name
	if err := validateEnvironmentVariableName(req.Name); err != nil {
		return nil, err
	}
	if req.Value == "" {
		return nil, fmt.Errorf("environment variable value cannot be empty")
	}

	// Check if name already exists
	exists, err := s.store.CheckEnvironmentVariableNameExists(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check environment variable name: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("environment variable with name '%s' already exists", req.Name)
	}

	// Create environment variable
	envVar, err := s.store.CreateEnvironmentVariable(req.Name, req.Value, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment variable: %w", err)
	}

	return s.convertToEnvironmentVariableInfo(envVar), nil
}

// GetEnvironmentVariable retrieves an environment variable by ID.
func (s *EnvironmentVariablesService) GetEnvironmentVariable(id string) (*EnvironmentVariableInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("environment variable ID cannot be empty")
	}

	envVar, err := s.store.GetEnvironmentVariable(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variable: %w", err)
	}

	return s.convertToEnvironmentVariableInfo(envVar), nil
}

// ListEnvironmentVariables returns all environment variables.
func (s *EnvironmentVariablesService) ListEnvironmentVariables() ([]EnvironmentVariableInfo, error) {
	envVars, err := s.store.ListEnvironmentVariables()
	if err != nil {
		return nil, fmt.Errorf("failed to list environment variables: %w", err)
	}

	result := make([]EnvironmentVariableInfo, len(envVars))
	for i, envVar := range envVars {
		result[i] = EnvironmentVariableInfo{
			ID:          envVar.ID,
			Name:        envVar.Name,
			Value:       envVar.Value,
			Description: envVar.Description.String,
			CreatedAt:   envVar.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   envVar.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return result, nil
}

// UpdateEnvironmentVariable updates an existing environment variable with new data.
func (s *EnvironmentVariablesService) UpdateEnvironmentVariable(req EnvironmentVariableUpdateRequest) (*EnvironmentVariableInfo, error) {
	// Validate input
	if req.ID == "" {
		return nil, fmt.Errorf("environment variable ID cannot be empty")
	}
	if err := validateEnvironmentVariableName(req.Name); err != nil {
		return nil, err
	}
	if req.Value == "" {
		return nil, fmt.Errorf("environment variable value cannot be empty")
	}

	// Update environment variable
	envVar, err := s.store.UpdateEnvironmentVariable(req.ID, req.Name, req.Value, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update environment variable: %w", err)
	}

	return s.convertToEnvironmentVariableInfo(envVar), nil
}

// DeleteEnvironmentVariable removes an environment variable from storage permanently.
func (s *EnvironmentVariablesService) DeleteEnvironmentVariable(id string) error {
	if id == "" {
		return fmt.Errorf("environment variable ID cannot be empty")
	}

	err := s.store.DeleteEnvironmentVariable(id)
	if err != nil {
		return fmt.Errorf("failed to delete environment variable: %w", err)
	}

	return nil
}

// GetEnvironmentVariableValue retrieves a specific environment variable value by name.
func (s *EnvironmentVariablesService) GetEnvironmentVariableValue(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("environment variable name cannot be empty")
	}

	envVar, err := s.store.GetEnvironmentVariableByName(name)
	if err != nil {
		return "", fmt.Errorf("failed to get environment variable: %w", err)
	}

	return envVar.Value, nil
}

// SearchEnvironmentVariables searches for environment variables by name pattern.
func (s *EnvironmentVariablesService) SearchEnvironmentVariables(pattern string) ([]EnvironmentVariableInfo, error) {
	if pattern == "" {
		return nil, fmt.Errorf("search pattern cannot be empty")
	}

	envVars, err := s.store.SearchEnvironmentVariablesByName(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search environment variables: %w", err)
	}

	result := make([]EnvironmentVariableInfo, len(envVars))
	for i, envVar := range envVars {
		result[i] = EnvironmentVariableInfo{
			ID:          envVar.ID,
			Name:        envVar.Name,
			Value:       envVar.Value,
			Description: envVar.Description.String,
			CreatedAt:   envVar.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   envVar.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return result, nil
}

// Helper functions

func (s *EnvironmentVariablesService) convertToEnvironmentVariableInfo(envVar database.EnvironmentVariable) *EnvironmentVariableInfo {
	return &EnvironmentVariableInfo{
		ID:          envVar.ID,
		Name:        envVar.Name,
		Value:       envVar.Value,
		Description: envVar.Description.String,
		CreatedAt:   envVar.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   envVar.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// validateEnvironmentVariableName validates that environment variable name contains only valid characters
func validateEnvironmentVariableName(name string) error {
	if name == "" {
		return fmt.Errorf("environment variable name cannot be empty")
	}
	
	// Only allow alphanumeric characters, underscores, and hyphens
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	if err != nil {
		return fmt.Errorf("failed to validate environment variable name: %w", err)
	}
	
	if !matched {
		return fmt.Errorf("environment variable name can only contain letters, numbers, underscores, and hyphens")
	}
	
	return nil
}