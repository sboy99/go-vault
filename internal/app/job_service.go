package app

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sboy99/go-vault/internal/domain"
)

// JobService serializes dump/restore so they never overlap and tracks job status.
type JobService struct {
	mu     sync.Mutex
	busy   bool
	cancel context.CancelFunc
	jobs   domain.JobRepository
}

func NewJobService(jobs domain.JobRepository) *JobService {
	return &JobService{jobs: jobs}
}

func (s *JobService) TryStart(jobType domain.JobType, backupID string, fn func(ctx context.Context) error) (*domain.Job, error) {
	s.mu.Lock()
	if s.busy {
		s.mu.Unlock()
		return nil, domain.ErrJobAlreadyRunning
	}
	s.busy = true
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.mu.Unlock()

	j := &domain.Job{
		ID:        uuid.New().String(),
		Type:      jobType,
		Status:    domain.JobStatusRunning,
		BackupID:  backupID,
		CreatedAt: time.Now().UTC(),
	}
	now := time.Now().UTC()
	j.StartedAt = &now
	_ = s.jobs.Save(context.Background(), j)

	go func() {
		defer func() {
			s.mu.Lock()
			s.busy = false
			s.cancel = nil
			s.mu.Unlock()
		}()
		err := fn(ctx)
		finished := time.Now().UTC()
		j.FinishedAt = &finished
		if err != nil {
			j.Status = domain.JobStatusFailed
			j.Error = err.Error()
		} else {
			j.Status = domain.JobStatusSucceeded
		}
		_ = s.jobs.Save(context.Background(), j)
	}()

	return j, nil
}

func (s *JobService) Busy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.busy
}

func (s *JobService) Cancel() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *JobService) Get(ctx context.Context, id string) (*domain.Job, error) {
	return s.jobs.FindByID(ctx, id)
}

func (s *JobService) List(ctx context.Context, limit, offset int) ([]*domain.Job, error) {
	return s.jobs.List(ctx, limit, offset)
}
