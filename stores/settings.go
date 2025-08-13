package stores

import (
	"context"

	"github.com/ironpark/teatime/internal/database"
	_ "github.com/ironpark/teatime/nodes"
)

// settingsStore manages application settings storage and retrieval.
// It provides access to user preferences like theme, auto-start behavior,
// and language settings stored in the database.
//
// The store is not safe for concurrent use. External synchronization
// is required when accessing from multiple goroutines.
type settingsStore struct {
	db *database.Client
}

// NewSettingsStore creates a new settings store with the given database client.
// The store provides methods to read and update application settings
// including theme, auto-start, and language preferences.
func NewSettingsStore(db *database.Client) *settingsStore {
	return &settingsStore{db: db}
}

// GetSettings retrieves the current application settings from the database.
// If no settings exist, it creates and returns default settings automatically.
//
// The returned Setting contains theme, auto-start, and language preferences.
// Returns an error if database operations fail during settings retrieval or creation.
func (s *settingsStore) GetSettings() (database.Setting, error) {
	// Try to get existing settings first
	settings, err := s.db.GetSettings(context.Background())
	if err != nil {
		// If no settings exist, create default settings
		settings, err = s.db.GetOrCreateSettings(context.Background())
		if err != nil {
			return database.Setting{}, err
		}
	}
	return settings, nil
}

// UpdateSettings updates all application settings with the provided values.
// This method updates theme, auto-start, and language settings in a single operation.
//
// Returns an error if the database update operation fails.
func (s *settingsStore) UpdateSettings(settings database.Setting) error {
	return s.db.UpdateSettings(context.Background(), database.UpdateSettingsParams{
		Theme:     settings.Theme,
		AutoStart: settings.AutoStart,
		Language:  settings.Language,
	})
}

// UpdateTheme updates only the theme setting to the specified value.
// Common theme values include "light", "dark", and "system".
//
// Returns an error if the database update fails.
func (s *settingsStore) UpdateTheme(theme string) error {
	return s.db.UpdateTheme(context.Background(), theme)
}

// UpdateAutoStart updates only the auto-start setting to the specified value.
// When true, the application will start automatically when the system boots.
//
// Returns an error if the database update fails.
func (s *settingsStore) UpdateAutoStart(autoStart bool) error {
	return s.db.UpdateAutoStart(context.Background(), autoStart)
}

// UpdateLanguage updates only the language setting to the specified value.
// The language value should be a valid locale code (e.g., "en", "ko", "ja").
//
// Returns an error if the database update fails.
func (s *settingsStore) UpdateLanguage(language string) error {
	return s.db.UpdateLanguage(context.Background(), language)
}
