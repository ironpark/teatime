package stores

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/ironpark/teatime/internal/database"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/pbkdf2"
)

// StorageType defines where secrets are stored
type StorageType string

const (
	StorageKeychain StorageType = "keychain"
)

// secretsStore manages secure storage and retrieval of secrets.
// It provides encrypted storage for API keys, tokens, and other sensitive values
// with support for multiple storage backends (database, keychain, environment).
//
// The store handles encryption/decryption automatically and is not safe for
// concurrent use. External synchronization is required when accessing from
// multiple goroutines.
type secretsStore struct {
	db *database.Client
}

// NewSecretsStore creates a new secrets store with the given database client.
// The store provides secure methods to create, read, update and delete secrets
// with automatic encryption for sensitive data.
func NewSecretsStore(db *database.Client) *secretsStore {
	return &secretsStore{db: db}
}

// CreateSecret creates a new secret with encrypted data storage.
// The secret data is automatically encrypted using AES-256-GCM with a random salt.
// 
// For keychain storage type, the encrypted data is also stored in the system keychain
// if available (currently macOS only).
//
// Returns the created secret with encrypted data and timestamps populated.
func (s *secretsStore) CreateSecret(name, description string, value string, storageType StorageType) (database.Credential, error) {
	id := uuid.New().String()

	// Prepare data map for single secret value
	data := map[string]string{"value": value}

	// Generate salt for encryption
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return database.Credential{}, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Encrypt the secret data
	encryptedData, err := s.encryptData(data, salt)
	if err != nil {
		return database.Credential{}, fmt.Errorf("failed to encrypt secret data: %w", err)
	}

	// Prepare description as nullable string
	var desc sql.NullString
	if description != "" {
		desc = sql.NullString{String: description, Valid: true}
	}

	// For keychain storage, generate a keychain key and store in system keychain
	var keychainKey sql.NullString
	if storageType == StorageKeychain {
		keychainKey = sql.NullString{
			String: fmt.Sprintf("teatime.secret.%s", id),
			Valid:  true,
		}
		
		// Store the original data in system keychain (no need to encrypt)
		if err := s.storeInKeychain(keychainKey.String, data); err != nil {
			return database.Credential{}, fmt.Errorf("failed to store in keychain: %w", err)
		}
		
		// Clear encrypted data for keychain storage (it's stored in keychain)
		encryptedData = []byte{}
		salt = []byte{}
	}

	secret, err := s.db.CreateCredential(context.Background(), database.CreateCredentialParams{
		ID:            id,
		Name:          name,
		Type:          "secret",
		Description:   desc,
		KeychainKey:   keychainKey,
		EncryptedData: encryptedData,
		Salt:          salt,
		StorageType:   string(storageType),
	})
	if err != nil {
		return database.Credential{}, fmt.Errorf("failed to create secret: %w", err)
	}

	return secret, nil
}

// GetSecret retrieves a secret by ID and decrypts its data.
// Returns the secret with decrypted data accessible via GetSecretValue method.
func (s *secretsStore) GetSecret(id string) (database.Credential, error) {
	secret, err := s.db.GetCredential(context.Background(), id)
	if err != nil {
		return database.Credential{}, fmt.Errorf("failed to get secret: %w", err)
	}

	return secret, nil
}

// GetSecretByName retrieves a secret by name and marks it as used.
// This method updates the last_used_at timestamp for tracking secret usage.
func (s *secretsStore) GetSecretByName(name string) (database.Credential, error) {
	secret, err := s.db.GetCredentialByName(context.Background(), name)
	if err != nil {
		return database.Credential{}, fmt.Errorf("failed to get secret by name: %w", err)
	}

	// Update last used timestamp
	err = s.db.UpdateCredentialLastUsed(context.Background(), secret.ID)
	if err != nil {
		// Log but don't fail on timestamp update error
		// TODO: Add proper logging
	}

	return secret, nil
}

// ListSecrets returns all secrets with basic information (no decrypted data).
// Sensitive data remains encrypted and is not included in the response.
func (s *secretsStore) ListSecrets() ([]database.ListCredentialsRow, error) {
	secrets, err := s.db.ListCredentials(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	return secrets, nil
}

// UpdateSecret updates an existing secret with new data.
// The secret data is re-encrypted with a new salt for enhanced security.
func (s *secretsStore) UpdateSecret(id, name, description string, value string, storageType StorageType) (database.Credential, error) {
	// Prepare data map for single secret value
	data := map[string]string{"value": value}

	// Generate new salt for re-encryption
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return database.Credential{}, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Encrypt the secret data with new salt
	encryptedData, err := s.encryptData(data, salt)
	if err != nil {
		return database.Credential{}, fmt.Errorf("failed to encrypt secret data: %w", err)
	}

	// Prepare description as nullable string
	var desc sql.NullString
	if description != "" {
		desc = sql.NullString{String: description, Valid: true}
	}

	// Handle keychain key for keychain storage
	var keychainKey sql.NullString
	if storageType == StorageKeychain {
		keychainKey = sql.NullString{
			String: fmt.Sprintf("teatime.secret.%s", id),
			Valid:  true,
		}
		
		// Update the original data in system keychain (no need to encrypt)
		if err := s.storeInKeychain(keychainKey.String, data); err != nil {
			return database.Credential{}, fmt.Errorf("failed to update in keychain: %w", err)
		}
		
		// Clear encrypted data for keychain storage (it's stored in keychain)
		encryptedData = []byte{}
		salt = []byte{}
	}

	secret, err := s.db.UpdateCredential(context.Background(), database.UpdateCredentialParams{
		Name:          name,
		Type:          "secret",
		Description:   desc,
		KeychainKey:   keychainKey,
		EncryptedData: encryptedData,
		Salt:          salt,
		StorageType:   string(storageType),
		ID:            id,
	})
	if err != nil {
		return database.Credential{}, fmt.Errorf("failed to update secret: %w", err)
	}

	return secret, nil
}

// DeleteSecret removes a secret from storage.
// For keychain storage, this also removes the secret from the system keychain.
func (s *secretsStore) DeleteSecret(id string) error {
	// Get secret first to check if it's stored in keychain
	secret, err := s.db.GetCredential(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to get secret for deletion: %w", err)
	}
	
	// Remove from keychain if it's keychain storage
	if secret.StorageType == string(StorageKeychain) && secret.KeychainKey.Valid {
		if err := s.deleteFromKeychain(secret.KeychainKey.String); err != nil {
			// Log error but don't fail the deletion
			// The secret might not exist in keychain or keychain might be unavailable
		}
	}
	
	err = s.db.DeleteCredential(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

// CheckSecretNameExists checks if a secret with the given name already exists.
// Returns true if a secret with the same name exists, false otherwise.
func (s *secretsStore) CheckSecretNameExists(name string) (bool, error) {
	count, err := s.db.CheckCredentialNameExists(context.Background(), name)
	if err != nil {
		return false, fmt.Errorf("failed to check secret name: %w", err)
	}

	return count > 0, nil
}

// SearchSecretsByName searches for secrets matching the given name pattern.
// The pattern supports SQL LIKE syntax (% for wildcards).
func (s *secretsStore) SearchSecretsByName(pattern string) ([]database.SearchCredentialsByNameRow, error) {
	secrets, err := s.db.SearchCredentialsByName(context.Background(), pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search secrets: %w", err)
	}

	return secrets, nil
}

// GetUnusedSecrets returns secrets that haven't been used recently.
// Useful for cleanup and maintenance of stale secrets.
func (s *secretsStore) GetUnusedSecrets() ([]database.GetUnusedCredentialsRow, error) {
	secrets, err := s.db.GetUnusedCredentials(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get unused secrets: %w", err)
	}

	return secrets, nil
}

// DecryptSecretData decrypts the secret data using the stored salt.
// Returns a map of key-value pairs containing the decrypted secret information.
func (s *secretsStore) DecryptSecretData(secret database.Credential) (map[string]string, error) {
	// For keychain storage, retrieve from keychain
	if secret.StorageType == string(StorageKeychain) && secret.KeychainKey.Valid {
		return s.retrieveFromKeychain(secret.KeychainKey.String)
	}
	
	// For database storage, decrypt from stored data
	if len(secret.EncryptedData) == 0 || len(secret.Salt) == 0 {
		return make(map[string]string), nil
	}

	return s.decryptData(secret.EncryptedData, secret.Salt)
}

// GetSecretValue retrieves and decrypts a specific field from a secret.
// This is a convenience method for getting the secret value.
func (s *secretsStore) GetSecretValue(id string) (string, error) {
	secret, err := s.GetSecret(id)
	if err != nil {
		return "", err
	}

	data, err := s.DecryptSecretData(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret data: %w", err)
	}

	value, exists := data["value"]
	if !exists {
		return "", fmt.Errorf("secret value not found")
	}

	return value, nil
}

// encryptData encrypts a map of secret data using AES-256-GCM
func (s *secretsStore) encryptData(data map[string]string, salt []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	// Convert map to JSON-like string for encryption
	dataStr := ""
	for k, v := range data {
		if dataStr != "" {
			dataStr += "&"
		}
		dataStr += fmt.Sprintf("%s=%s", k, base64.StdEncoding.EncodeToString([]byte(v)))
	}

	// Derive key from master password + salt
	masterPassword := s.getMasterPassword()
	key := pbkdf2.Key([]byte(masterPassword), salt, 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, []byte(dataStr), nil)
	return ciphertext, nil
}

// decryptData decrypts secret data using AES-256-GCM
func (s *secretsStore) decryptData(encryptedData, salt []byte) (map[string]string, error) {
	if len(encryptedData) == 0 {
		return make(map[string]string), nil
	}

	// Derive key from master password + salt
	masterPassword := s.getMasterPassword()
	key := pbkdf2.Key([]byte(masterPassword), salt, 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Parse back to map
	result := make(map[string]string)
	dataStr := string(plaintext)
	if dataStr == "" {
		return result, nil
	}

	pairs := []rune(dataStr)
	currentPair := ""
	for _, r := range pairs {
		if r == '&' {
			if err := s.parsePair(currentPair, result); err != nil {
				return nil, err
			}
			currentPair = ""
		} else {
			currentPair += string(r)
		}
	}
	
	// Parse the last pair
	if currentPair != "" {
		if err := s.parsePair(currentPair, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// parsePair parses a key=value pair and adds it to the result map
func (s *secretsStore) parsePair(pair string, result map[string]string) error {
	parts := []rune(pair)
	equalIndex := -1
	for i, r := range parts {
		if r == '=' {
			equalIndex = i
			break
		}
	}
	
	if equalIndex == -1 {
		return fmt.Errorf("invalid key-value pair: %s", pair)
	}

	key := string(parts[:equalIndex])
	encodedValue := string(parts[equalIndex+1:])
	
	value, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		return fmt.Errorf("failed to decode value: %w", err)
	}
	
	result[key] = string(value)
	return nil
}

// getMasterPassword gets the master password for encryption
// In production, this should be derived from user input or system keychain
func (s *secretsStore) getMasterPassword() string {
	// TODO: Implement proper master password management
	// This could come from:
	// 1. User input (with secure prompt)
	// 2. System keychain
	// 3. Hardware security module
	// 4. Derived from machine-specific data
	
	// For now, use environment variable or default
	if password := os.Getenv("TEATIME_MASTER_PASSWORD"); password != "" {
		return password
	}
	
	// Fallback to a default (NOT SECURE - only for development)
	return "teatime-default-master-key-please-change-in-production"
}

// Keychain integration methods

// storeInKeychain stores secret data in the system keychain
func (s *secretsStore) storeInKeychain(keychainKey string, data map[string]string) error {
	// Convert data to JSON for storage
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}
	
	// Store in system keychain
	err = keyring.Set("teatime", keychainKey, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to store in keychain: %w", err)
	}
	
	return nil
}

// retrieveFromKeychain retrieves secret data from the system keychain
func (s *secretsStore) retrieveFromKeychain(keychainKey string) (map[string]string, error) {
	// Get from system keychain
	jsonData, err := keyring.Get("teatime", keychainKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve from keychain: %w", err)
	}
	
	// Parse JSON data
	var data map[string]string
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret data: %w", err)
	}
	
	return data, nil
}

// deleteFromKeychain removes secret data from the system keychain
func (s *secretsStore) deleteFromKeychain(keychainKey string) error {
	err := keyring.Delete("teatime", keychainKey)
	if err != nil {
		return fmt.Errorf("failed to delete from keychain: %w", err)
	}
	
	return nil
}