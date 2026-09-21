package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/sboy99/go-vault/internal/domain"
)

// JobRepository is an in-memory JobRepository for tests.
type JobRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Job
}

var _ domain.JobRepository = (*JobRepository)(nil)

func NewJobRepository() *JobRepository {
	return &JobRepository{data: map[string]*domain.Job{}}
}

func (r *JobRepository) Save(ctx context.Context, j *domain.Job) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *j
	r.data[j.ID] = &cp
	return nil
}

func (r *JobRepository) FindByID(ctx context.Context, id string) (*domain.Job, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.data[id]
	if !ok {
		return nil, domain.ErrJobNotFound
	}
	cp := *j
	return &cp, nil
}

func (r *JobRepository) List(ctx context.Context, limit, offset int) ([]*domain.Job, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Job, 0, len(r.data))
	for _, j := range r.data {
		cp := *j
		out = append(out, &cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	if offset > len(out) {
		return []*domain.Job{}, nil
	}
	out = out[offset:]
	if limit > 0 && limit < len(out) {
		out = out[:limit]
	}
	return out, nil
}
