package stores

import (
	"context"

	"github.com/ironpark/teatime/internal/database"
	_ "github.com/ironpark/teatime/nodes"
)

type settingsStore struct {
	db *database.Client
}

func NewSettingsStore(db *database.Client) *settingsStore {
	return &settingsStore{db: db}
}

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

func (s *settingsStore) UpdateSettings(settings database.Setting) error {
	return s.db.UpdateSettings(context.Background(), database.UpdateSettingsParams{
		Theme:     settings.Theme,
		AutoStart: settings.AutoStart,
		Language:  settings.Language,
	})
}

func (s *settingsStore) UpdateTheme(theme string) error {
	return s.db.UpdateTheme(context.Background(), theme)
}

func (s *settingsStore) UpdateAutoStart(autoStart bool) error {
	return s.db.UpdateAutoStart(context.Background(), autoStart)
}

func (s *settingsStore) UpdateLanguage(language string) error {
	return s.db.UpdateLanguage(context.Background(), language)
}
