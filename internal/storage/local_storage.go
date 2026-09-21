package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sboy99/go-vault/internal/domain"
)

var _ domain.ArtifactStore = (*LocalStorage)(nil)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

func (ls *LocalStorage) Save(ctx context.Context, key string, r io.Reader) (domain.ObjectInfo, error) {
	if err := os.MkdirAll(ls.BasePath, 0o755); err != nil {
		return domain.ObjectInfo{}, fmt.Errorf("mkdir backups: %w", err)
	}
	finalPath := ls.filePath(key)
	tmpPath := finalPath + ".tmp"

	f, err := os.Create(tmpPath)
	if err != nil {
		return domain.ObjectInfo{}, fmt.Errorf("create temp file: %w", err)
	}

	hasher := sha256.New()
	writer := io.MultiWriter(f, hasher)

	n, copyErr := io.Copy(writer, readerWithContext(ctx, r))
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmpPath)
		return domain.ObjectInfo{}, fmt.Errorf("write backup: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(tmpPath)
		return domain.ObjectInfo{}, closeErr
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		_ = os.Remove(tmpPath)
		return domain.ObjectInfo{}, fmt.Errorf("rename backup: %w", err)
	}

	return domain.ObjectInfo{
		Key:          key,
		Size:         n,
		SHA256:       hex.EncodeToString(hasher.Sum(nil)),
		LastModified: time.Now().UTC(),
	}, nil
}

func (ls *LocalStorage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	f, err := os.Open(ls.filePath(key))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (ls *LocalStorage) Delete(ctx context.Context, key string) error {
	err := os.Remove(ls.filePath(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (ls *LocalStorage) List(ctx context.Context, prefix string) ([]domain.ObjectInfo, error) {
	if err := os.MkdirAll(ls.BasePath, 0o755); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(ls.BasePath)
	if err != nil {
		return nil, err
	}
	var out []domain.ObjectInfo
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		name := ent.Name()
		if strings.HasSuffix(name, ".tmp") {
			continue
		}
		if prefix != "" && !strings.HasPrefix(name, prefix) {
			continue
		}
		info, err := ent.Info()
		if err != nil {
			continue
		}
		out = append(out, domain.ObjectInfo{
			Key:          name,
			Size:         info.Size(),
			LastModified: info.ModTime().UTC(),
		})
	}
	return out, nil
}

func (ls *LocalStorage) filePath(key string) string {
	return filepath.Join(ls.BasePath, filepath.Base(key))
}

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func readerWithContext(ctx context.Context, r io.Reader) io.Reader {
	return &ctxReader{ctx: ctx, r: r}
}

func (c *ctxReader) Read(p []byte) (int, error) {
	select {
	case <-c.ctx.Done():
		return 0, c.ctx.Err()
	default:
		return c.r.Read(p)
	}
}
