package meta

import (
	"fmt"
	"time"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/pkg/boltdb"
	"github.com/sboy99/go-vault/internal/utils"
)

type BackupStatus string

const (
	StatusRunning BackupStatus = "running"
	StatusSuccess BackupStatus = "success"
	StatusFailed  BackupStatus = "failed"
)

type BackupMeta struct {
	BackupId     string              `json:"id"`
	Name         string              `json:"name"`
	StorageKey   string              `json:"storage_key"`
	DatabaseType config.DatabaseEnum `json:"database_type"`
	StorageType  config.StorageEnum  `json:"storage_type"`
	Status       BackupStatus        `json:"status"`
	Error        string              `json:"error,omitempty"`
	StartedAt    time.Time           `json:"started_at"`
	FinishedAt   *time.Time          `json:"finished_at,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	SizeBytes    int64               `json:"size_bytes"`
	SHA256       string              `json:"sha256,omitempty"`
	PgVersion    string              `json:"pg_version,omitempty"`
	Format       string              `json:"format"`
	Verified     bool                `json:"verified"`
}

func NewBackupMeta(name string, dbType config.DatabaseEnum, storageType config.StorageEnum) *BackupMeta {
	now := utils.GetNow().UTC()
	uid := utils.GenerateUID()
	id := fmt.Sprintf("%s_%s", now.Format(time.RFC3339Nano), uid)
	return &BackupMeta{
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

func (b *BackupMeta) Key() string {
	return b.BackupId
}

func (b *BackupMeta) Save() error {
	data, err := utils.MarshalJSON(b)
	if err != nil {
		return err
	}
	return boltdb.Save(_BACKUP_META, b.BackupId, data)
}

func (b *BackupMeta) MarkSuccess(size int64, sha, pgVersion string, verified bool) error {
	now := utils.GetNow().UTC()
	b.Status = StatusSuccess
	b.FinishedAt = &now
	b.SizeBytes = size
	b.SHA256 = sha
	b.PgVersion = pgVersion
	b.Verified = verified
	b.Error = ""
	return b.Save()
}

func (b *BackupMeta) MarkFailed(errMsg string) error {
	now := utils.GetNow().UTC()
	b.Status = StatusFailed
	b.FinishedAt = &now
	b.Error = errMsg
	return b.Save()
}

func GetBackupMeta(id string) (*BackupMeta, error) {
	data, err := boltdb.Get(_BACKUP_META, id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("backup %s not found", id)
	}
	var backupMeta BackupMeta
	if err := utils.UnmarshalJSON(data, &backupMeta); err != nil {
		return nil, err
	}
	return &backupMeta, nil
}

func ListBackupMeta(size, offset int) ([]*BackupMeta, error) {
	rows, err := boltdb.ListReverse(_BACKUP_META, size, offset)
	if err != nil {
		return nil, err
	}
	return decodeBackupList(rows)
}

func ListAllBackupMeta() ([]*BackupMeta, error) {
	rows, err := boltdb.ListAll(_BACKUP_META)
	if err != nil {
		return nil, err
	}
	return decodeBackupList(rows)
}

func DeleteBackupMeta(id string) error {
	return boltdb.Delete(_BACKUP_META, id)
}

func decodeBackupList(rows [][]byte) ([]*BackupMeta, error) {
	var backupMetas []*BackupMeta
	for _, row := range rows {
		var backupMeta BackupMeta
		if err := utils.UnmarshalJSON(row, &backupMeta); err != nil {
			return nil, err
		}
		backupMetas = append(backupMetas, &backupMeta)
	}
	return backupMetas, nil
}
