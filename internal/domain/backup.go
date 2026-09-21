package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DatabaseType string

const (
	DatabasePostgreSQL DatabaseType = "POSTGRESQL"
)

type StorageType string

const (
	StorageLocal StorageType = "LOCAL"
	StorageCloud StorageType = "CLOUD"
)

type BackupStatus string

const (
	StatusRunning BackupStatus = "running"
	StatusSuccess BackupStatus = "success"
	StatusFailed  BackupStatus = "failed"
)

// Backup is the domain entity for a backup artifact and its metadata.
type Backup struct {
	BackupId     string       `json:"id"`
	Name         string       `json:"name"`
	StorageKey   string       `json:"storage_key"`
	DatabaseType DatabaseType `json:"database_type"`
	StorageType  StorageType  `json:"storage_type"`
	Status       BackupStatus `json:"status"`
	Error        string       `json:"error,omitempty"`
	StartedAt    time.Time    `json:"started_at"`
	FinishedAt   *time.Time   `json:"finished_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	SizeBytes    int64        `json:"size_bytes"`
	SHA256       string       `json:"sha256,omitempty"`
	PgVersion    string       `json:"pg_version,omitempty"`
	Format       string       `json:"format"`
	Verified     bool         `json:"verified"`
}

// NewBackup creates a running backup entity with a time-sortable ID.
func NewBackup(name string, dbType DatabaseType, storageType StorageType) *Backup {
	now := time.Now().UTC()
	parts := strings.Split(uuid.New().String(), "-")
	uid := parts[len(parts)-1]
	id := fmt.Sprintf("%s_%s", now.Format(time.RFC3339Nano), uid)
	return &Backup{
		BackupId:     id,
		Name:         name,
		StorageKey:   name,
		DatabaseType: dbType,
		StorageType:  storageType,
		Status:       StatusRunning,
		StartedAt:    now,
		CreatedAt:    now,
		Format:       "custom",
	}
}

// MarkSuccess returns a copy marked successful.
func (b Backup) MarkSuccess(size int64, sha, pgVersion string, verified bool) Backup {
	now := time.Now().UTC()
	b.Status = StatusSuccess
	b.FinishedAt = &now
	b.SizeBytes = size
	b.SHA256 = sha
	b.PgVersion = pgVersion
	b.Verified = verified
	b.Error = ""
	return b
}

// MarkFailed returns a copy marked failed.
func (b Backup) MarkFailed(errMsg string) Backup {
	now := time.Now().UTC()
	b.Status = StatusFailed
	b.FinishedAt = &now
	b.Error = errMsg
	return b
}
