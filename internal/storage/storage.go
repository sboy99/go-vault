package storage

import (
	"fmt"

	"github.com/sboy99/go-vault/internal/config"
	"github.com/sboy99/go-vault/internal/domain"
)

// NewStorage builds an ArtifactStore from config.
func NewStorage(cfg *config.Config) (domain.ArtifactStore, error) {
	switch cfg.Storage.Type {
	case config.LOCAL:
		return NewLocalStorage(cfg.Storage.Dest), nil
	case config.CLOUD:
		switch cfg.Storage.Cloud.Type {
		case config.AWS:
			return NewAWSCloudStorage(cfg)
		default:
			return nil, fmt.Errorf("unsupported cloud type %q", cfg.Storage.Cloud.Type)
		}
	default:
		return nil, fmt.Errorf("unsupported storage type %q", cfg.Storage.Type)
	}
}
