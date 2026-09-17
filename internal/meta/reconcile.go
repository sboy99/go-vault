package meta

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/pkg/logger"
)

// StorageLister lists stored backup keys for reconcile.
type StorageLister interface {
	List(ctx context.Context, prefix string) ([]struct {
		Key          string
		Size         int64
		LastModified time.Time
	}, error)
}

// ObjectRef is a minimal storage object for reconcile.
type ObjectRef struct {
	Key          string
	Size         int64
	LastModified time.Time
}

type ListerFunc func(ctx context.Context, prefix string) ([]ObjectRef, error)

// Reconcile ensures every storage artifact has a metadata row.
func Reconcile(ctx context.Context, list ListerFunc, dbType config.DatabaseEnum, storageType config.StorageEnum) error {
	objects, err := list(ctx, "")
	if err != nil {
		return err
	}
	existing, err := ListAllBackupMeta()
	if err != nil {
		return err
	}
	known := map[string]struct{}{}
	for _, b := range existing {
		known[b.StorageKey] = struct{}{}
		if b.Name != "" {
			known[b.Name] = struct{}{}
		}
	}
	for _, obj := range objects {
		key := obj.Key
		base := filepath.Base(key)
		if _, ok := known[key]; ok {
			continue
		}
		if _, ok := known[base]; ok {
			continue
		}
		if !strings.HasSuffix(base, ".dump") && !strings.HasSuffix(base, ".sql") {
			continue
		}
		m := NewBackupMeta(base, dbType, storageType)
		m.StorageKey = key
		m.Status = StatusSuccess
		m.SizeBytes = obj.Size
		m.CreatedAt = obj.LastModified
		m.StartedAt = obj.LastModified
		finished := obj.LastModified
		m.FinishedAt = &finished
		if err := m.Save(); err != nil {
			logger.Warn("reconcile save failed for %s: %v", key, err)
			continue
		}
		logger.Info("reconciled missing metadata for %s", key)
	}
	return nil
}
