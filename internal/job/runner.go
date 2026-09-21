package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sboy99/go-vault/internal/meta"
	"github.com/sboy99/go-vault/pkg/boltdb"
	"github.com/sboy99/go-vault/pkg/utils"
)

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
)

type Type string

const (
	TypeBackup  Type = "backup"
	TypeRestore Type = "restore"
)

type Job struct {
	ID         string     `json:"id"`
	Type       Type       `json:"type"`
	Status     Status     `json:"status"`
	Error      string     `json:"error,omitempty"`
	BackupID   string     `json:"backup_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

// Runner serializes dump/restore/prune so they never overlap.
type Runner struct {
	mu     sync.Mutex
	busy   bool
	cancel context.CancelFunc
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) TryStart(jobType Type, backupID string, fn func(ctx context.Context) error) (*Job, error) {
	r.mu.Lock()
	if r.busy {
		r.mu.Unlock()
		return nil, fmt.Errorf("another job is already running")
	}
	r.busy = true
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.mu.Unlock()

	j := &Job{
		ID:        utils.GenerateUUID(),
		Type:      jobType,
		Status:    StatusRunning,
		BackupID:  backupID,
		CreatedAt: time.Now().UTC(),
	}
	now := time.Now().UTC()
	j.StartedAt = &now
	_ = saveJob(j)

	go func() {
		defer func() {
			r.mu.Lock()
			r.busy = false
			r.cancel = nil
			r.mu.Unlock()
		}()
		err := fn(ctx)
		finished := time.Now().UTC()
		j.FinishedAt = &finished
		if err != nil {
			j.Status = StatusFailed
			j.Error = err.Error()
		} else {
			j.Status = StatusSucceeded
		}
		_ = saveJob(j)
	}()

	return j, nil
}

func (r *Runner) Busy() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.busy
}

func (r *Runner) Cancel() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
}

func Get(id string) (*Job, error) {
	data, err := boltdb.Get(meta.BucketJob(), id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("job %s not found", id)
	}
	var j Job
	if err := utils.UnmarshalJSON(data, &j); err != nil {
		return nil, err
	}
	return &j, nil
}

func List(limit, offset int) ([]*Job, error) {
	rows, err := boltdb.List(meta.BucketJob(), limit, offset)
	if err != nil {
		return nil, err
	}
	return decodeJobList(rows)
}

func decodeJobList(rows [][]byte) ([]*Job, error) {
	var jobs []*Job
	for _, row := range rows {
		var j Job
		if err := utils.UnmarshalJSON(row, &j); err != nil {
			return nil, err
		}
	}
	return jobs, nil
}

func saveJob(j *Job) error {
	data, err := utils.MarshalJSON(j)
	if err != nil {
		return err
	}
	return boltdb.Save(meta.BucketJob(), j.ID, data)
}
