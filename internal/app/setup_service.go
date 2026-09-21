package app

import (
	"fmt"

	"github.com/sboy99/go-vault/internal/domain"
)

// SetupService validates and persists setup wizard answers.
type SetupService struct {
	store domain.SettingsStore
}

func NewSetupService(store domain.SettingsStore) *SetupService {
	return &SetupService{store: store}
}

func (s *SetupService) SaveSettings(settings domain.Settings) error {
	if settings.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if settings.DBHost == "" {
		return fmt.Errorf("database host is required")
	}
	if settings.DBUsername == "" {
		return fmt.Errorf("database username is required")
	}
	if settings.StorageType == "" {
		return fmt.Errorf("storage type is required")
	}
	if settings.DBSSLMode == "" {
		settings.DBSSLMode = "require"
	}
	return s.store.Save(settings)
}
