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

// IStorage is a streaming storage backend.
type IStorage interface {
	Save(ctx context.Context, key string, r io.Reader) (ObjectInfo, error)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]ObjectInfo, error)
}

type Storage struct {
	storageMap map[config.StorageEnum]IStorage
	typ        config.StorageEnum
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	s := &Storage{
		typ:        cfg.Storage.Type,
		storageMap: map[config.StorageEnum]IStorage{},
	}
	switch cfg.Storage.Type {
	case config.LOCAL:
		s.storageMap[config.LOCAL] = NewLocalStorage(cfg.Storage.Dest)
	case config.CLOUD:
		cloud, err := NewCloudStorage(cfg)
		if err != nil {
			return nil, err
		}
		s.storageMap[config.CLOUD] = cloud
	default:
		return nil, fmt.Errorf("unsupported storage type %q", cfg.Storage.Type)
	}
	return s, nil
}

func (s *Storage) Save(ctx context.Context, key string, r io.Reader) (ObjectInfo, error) {
	return s.backend().Save(ctx, key, r)
}

func (s *Storage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.backend().Open(ctx, key)
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.backend().Delete(ctx, key)
}

func (s *Storage) List(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	return s.backend().List(ctx, prefix)
}

func (s *Storage) backend() IStorage {
	return s.storageMap[s.typ]
}
