package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/sboy99/go-vault/internal/domain"
)

// BackupRepository is an in-memory BackupRepository for tests.
type BackupRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Backup
}

var _ domain.BackupRepository = (*BackupRepository)(nil)

func NewBackupRepository() *BackupRepository {
	return &BackupRepository{data: map[string]*domain.Backup{}}
}

func (r *BackupRepository) Save(ctx context.Context, b *domain.Backup) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *b
	r.data[b.BackupId] = &cp
	return nil
}

func (r *BackupRepository) FindByID(ctx context.Context, id string) (*domain.Backup, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.data[id]
	if !ok {
		return nil, domain.ErrBackupNotFound
	}
	cp := *b
	return &cp, nil
}

func (r *BackupRepository) FindAll(ctx context.Context) ([]*domain.Backup, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Backup, 0, len(r.data))
	for _, b := range r.data {
		cp := *b
		out = append(out, &cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *BackupRepository) List(ctx context.Context, limit, offset int) ([]*domain.Backup, error) {
	all, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	if offset > len(all) {
		return []*domain.Backup{}, nil
	}
	all = all[offset:]
	if limit > 0 && limit < len(all) {
		all = all[:limit]
	}
	return all, nil
}

func (r *BackupRepository) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}
