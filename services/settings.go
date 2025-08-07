package services

import (
	"github.com/ironpark/teatime/internal/database"
	"github.com/ironpark/teatime/stores"
)

type SettingsService struct {
	store stores.SettingsStore
}

func NewSettingsService(store stores.SettingsStore) *SettingsService {
	return &SettingsService{store: store}
}

func (s *SettingsService) GetSettings() (database.Setting, error) {
	return s.store.GetSettings()
}

func (s *SettingsService) UpdateSettings(settings database.Setting) error {
	return s.store.UpdateSettings(settings)
}

func (s *SettingsService) UpdateTheme(theme string) error {
	return s.store.UpdateTheme(theme)
}
