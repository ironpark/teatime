package services

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/ironpark/teatime/internal/database"
	"github.com/ironpark/teatime/stores"
)

// SecretsService provides secure secret management functionality to the frontend.
// It handles CRUD operations for encrypted secrets including API keys, tokens,
// and other sensitive data with support for multiple storage backends.
//
// All sensitive data is automatically encrypted before storage and decrypted on retrieval.
// The service is safe for concurrent use and provides validation for secret operations.
type SecretsService struct {
	store *stores.Store
}

// SecretInfo represents basic secret information without sensitive data.
// Used for listing and displaying secrets in the UI safely.
type SecretInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StorageType string `json:"storage_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	LastUsedAt  string `json:"last_used_at"`
}

// SecretCreateRequest represents the data needed to create a new secret.
type SecretCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// SecretUpdateRequest represents the data needed to update an existing secret.
type SecretUpdateRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// SecretValueResponse contains a decrypted secret value.
// Used for retrieving specific secret values for workflow execution.
type SecretValueResponse struct {
	Value string `json:"value"`
	Found bool   `json:"found"`
}

// NewSecretsService creates a new secrets service with the given store.
// The service provides methods for secure secret management with automatic
// encryption and decryption of sensitive data.
func NewSecretsService(store *stores.Store) *SecretsService {
	return &SecretsService{store: store}
}

// CreateSecret creates a new encrypted secret with the provided data.
// The secret data is automatically encrypted using AES-256-GCM with a unique salt.
//
// Returns an error if:
// - A secret with the same name already exists
// - The storage type is invalid
// - Encryption fails
// - Database operation fails
func (s *SecretsService) CreateSecret(req SecretCreateRequest) (*SecretInfo, error) {
	// Validate secret name
	if err := validateSecretName(req.Name); err != nil {
		return nil, err
	}
	if req.Value == "" {
		return nil, fmt.Errorf("secret value cannot be empty")
	}

	// Check if name already exists
	exists, err := s.store.CheckSecretNameExists(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check secret name: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("secret with name '%s' already exists", req.Name)
	}

	// Use keychain storage only
	storageType := stores.StorageKeychain

	// Create secret
	secret, err := s.store.CreateSecret(req.Name, req.Description, req.Value, storageType)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}

	return s.convertToSecretInfo(secret), nil
}

// GetSecret retrieves a secret by ID with basic information only.
// Sensitive data is not included in the response for security.
//
// To retrieve decrypted secret data, use GetSecretValue.
func (s *SecretsService) GetSecret(id string) (*SecretInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("secret ID cannot be empty")
	}

	secret, err := s.store.GetSecret(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return s.convertToSecretInfo(secret), nil
}

// ListSecrets returns all secrets with basic information only.
// Sensitive data is not included in the response for security.
//
// This method is safe to use in UI contexts where secret lists are displayed.
func (s *SecretsService) ListSecrets() ([]SecretInfo, error) {
	secrets, err := s.store.ListSecrets()
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	result := make([]SecretInfo, len(secrets))
	for i, secret := range secrets {
		result[i] = SecretInfo{
			ID:          secret.ID,
			Name:        secret.Name,
			Description: secret.Description.String,
			StorageType: secret.StorageType,
			CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
			LastUsedAt:  formatNullTime(secret.LastUsedAt),
		}
	}

	return result, nil
}

// UpdateSecret updates an existing secret with new data.
// The secret data is re-encrypted with a new salt for enhanced security.
//
// Returns an error if the secret doesn't exist or update fails.
func (s *SecretsService) UpdateSecret(req SecretUpdateRequest) (*SecretInfo, error) {
	// Validate input
	if req.ID == "" {
		return nil, fmt.Errorf("secret ID cannot be empty")
	}
	if err := validateSecretName(req.Name); err != nil {
		return nil, err
	}
	if req.Value == "" {
		return nil, fmt.Errorf("secret value cannot be empty")
	}

	// Use keychain storage only
	storageType := stores.StorageKeychain

	// Update secret
	secret, err := s.store.UpdateSecret(req.ID, req.Name, req.Description, req.Value, storageType)
	if err != nil {
		return nil, fmt.Errorf("failed to update secret: %w", err)
	}

	return s.convertToSecretInfo(secret), nil
}

// DeleteSecret removes a secret from storage permanently.
// This operation cannot be undone.
//
// For keychain storage, this also removes the secret from the system keychain.
func (s *SecretsService) DeleteSecret(id string) error {
	if id == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	err := s.store.DeleteSecret(id)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

// GetSecretValue retrieves and decrypts a specific secret by name.
// This method marks the secret as used and updates the last_used_at timestamp.
//
// Use this method for retrieving secret values during workflow execution.
func (s *SecretsService) GetSecretValue(secretName string) (*SecretValueResponse, error) {
	if secretName == "" {
		return nil, fmt.Errorf("secret name cannot be empty")
	}

	// Get secret by name (this updates last_used_at)
	secret, err := s.store.GetSecretByName(secretName)
	if err != nil {
		return &SecretValueResponse{Found: false}, nil // Return not found instead of error
	}

	// Get the secret value
	value, err := s.store.GetSecretValue(secret.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret value: %w", err)
	}

	return &SecretValueResponse{
		Value: value,
		Found: true,
	}, nil
}

// GetSecretData retrieves and decrypts the value from a secret.
// This method is primarily used for secret editing in the UI.
//
// WARNING: This returns sensitive decrypted data. Use with caution.
func (s *SecretsService) GetSecretData(id string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("secret ID cannot be empty")
	}

	value, err := s.store.GetSecretValue(id)
	if err != nil {
		return "", fmt.Errorf("failed to get secret data: %w", err)
	}

	return value, nil
}

// SearchSecrets searches for secrets by name pattern.
// The pattern supports SQL LIKE syntax (% for wildcards).
func (s *SecretsService) SearchSecrets(pattern string) ([]SecretInfo, error) {
	if pattern == "" {
		return nil, fmt.Errorf("search pattern cannot be empty")
	}

	secrets, err := s.store.SearchSecretsByName(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search secrets: %w", err)
	}

	result := make([]SecretInfo, len(secrets))
	for i, secret := range secrets {
		result[i] = SecretInfo{
			ID:          secret.ID,
			Name:        secret.Name,
			Description: secret.Description.String,
			StorageType: secret.StorageType,
			CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
			LastUsedAt:  formatNullTime(secret.LastUsedAt),
		}
	}

	return result, nil
}

// GetUnusedSecrets returns secrets that haven't been used recently.
// Useful for cleanup and maintenance of stale secrets.
func (s *SecretsService) GetUnusedSecrets() ([]SecretInfo, error) {
	secrets, err := s.store.GetUnusedSecrets()
	if err != nil {
		return nil, fmt.Errorf("failed to get unused secrets: %w", err)
	}

	result := make([]SecretInfo, len(secrets))
	for i, secret := range secrets {
		result[i] = SecretInfo{
			ID:          secret.ID,
			Name:        secret.Name,
			Description: secret.Description.String,
			StorageType: secret.StorageType,
			CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
			LastUsedAt:  formatNullTime(secret.LastUsedAt),
		}
	}

	return result, nil
}

// GetAvailableStorageTypes returns a list of supported storage types.
func (s *SecretsService) GetAvailableStorageTypes() []string {
	return []string{
		string(stores.StorageKeychain),
	}
}

// Helper functions

func (s *SecretsService) convertToSecretInfo(secret database.Credential) *SecretInfo {
	return &SecretInfo{
		ID:          secret.ID,
		Name:        secret.Name,
		Description: secret.Description.String,
		StorageType: secret.StorageType,
		CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
		LastUsedAt:  formatNullTime(secret.LastUsedAt),
	}
}

func formatNullTime(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

// validateSecretName validates that secret name contains only alphanumeric characters and underscores
func validateSecretName(name string) error {
	if name == "" {
		return fmt.Errorf("secret name cannot be empty")
	}
	
	// Only allow alphanumeric characters, underscores, and hyphens
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	if err != nil {
		return fmt.Errorf("failed to validate secret name: %w", err)
	}
	
	if !matched {
		return fmt.Errorf("secret name can only contain letters, numbers, underscores, and hyphens")
	}
	
	return nil
}

