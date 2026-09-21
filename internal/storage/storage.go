package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sboy99/go-vault/internal/config"
)

// ObjectInfo describes a stored backup artifact.
type ObjectInfo struct {
	Key          string
	Size         int64
	SHA256       string
	LastModified time.Time
}

// Backend is a streaming storage backend.
type Backend interface {
	Save(ctx context.Context, key string, r io.Reader) (ObjectInfo, error)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]ObjectInfo, error)
}

func NewStorage(cfg *config.Config) (Backend, error) {
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
