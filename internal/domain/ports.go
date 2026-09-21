package domain

import (
	"context"
	"io"
	"time"
)

// BackupRepository persists backup metadata.
type BackupRepository interface {
	Save(ctx context.Context, b *Backup) error
	FindByID(ctx context.Context, id string) (*Backup, error)
	FindAll(ctx context.Context) ([]*Backup, error)
	List(ctx context.Context, limit, offset int) ([]*Backup, error)
	Delete(ctx context.Context, id string) error
}

// JobRepository persists job metadata.
type JobRepository interface {
	Save(ctx context.Context, j *Job) error
	FindByID(ctx context.Context, id string) (*Job, error)
	List(ctx context.Context, limit, offset int) ([]*Job, error)
}

// ArtifactStore is a streaming storage backend for backup artifacts.
type ArtifactStore interface {
	Save(ctx context.Context, key string, r io.Reader) (ObjectInfo, error)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]ObjectInfo, error)
}

// DumpEngine runs database dump/restore/verify operations.
type DumpEngine interface {
	DumpTo(ctx context.Context, p ConnParams, w io.Writer) (*DumpResult, error)
	Verify(ctx context.Context, dumpPath string, serverMajor int) error
	Restore(ctx context.Context, p ConnParams, dumpPath string, jobs int) error
	ServerVersion(ctx context.Context, p ConnParams) (string, int, error)
}

// Notifier sends failure notifications.
type Notifier interface {
	NotifyFailure(ctx context.Context, event, message string)
}

// MetricsRecorder records backup/restore operational metrics.
type MetricsRecorder interface {
	BackupSucceeded(started time.Time, size int64)
	BackupFailed()
	RestoreFinished(status string)
	BackupPruned()
}

// SettingsStore persists setup wizard answers.
type SettingsStore interface {
	Save(s Settings) error
}
