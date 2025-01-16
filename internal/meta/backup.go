package meta

import (
	"time"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/pkg/utils"
)

type BackupMeta struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	DatabaseType config.DatabaseEnum `json:"database_type"`
	StorageType  config.StorageEnum  `json:"storage_type"`
	CreatedAt    time.Time           `json:"created_at"`
}

func NewBackupMeta(name string, dbType config.DatabaseEnum, storageType config.StorageEnum) *BackupMeta {
	return &BackupMeta{
		ID:           utils.GenerateUUID(),
		CreatedAt:    utils.GetNow(),
		Name:         name,
		DatabaseType: dbType,
		StorageType:  storageType,
	}
}

func (b *BackupMeta) Save() error {
	backupMetaJSON, err := utils.MarshalJSON(b)
	if err != nil {
		return err
	}
	return saveMetaData(_BOLT_BACKUP_BUCKET, b.ID, backupMetaJSON)
}

func GetBackupMeta(id string) (*BackupMeta, error) {
	backupMetaJSON, err := getMetaData(_BOLT_BACKUP_BUCKET, id)
	if err != nil {
		return nil, err
	}
	var backupMeta BackupMeta
	if err := utils.UnmarshalJSON(backupMetaJSON, &backupMeta); err != nil {
		return nil, err
	}
	return &backupMeta, nil
}

func ListBackupMeta(size int) ([]*BackupMeta, error) {
	backupMetaJSONs, err := listMetaData(_BOLT_BACKUP_BUCKET, size)
	if err != nil {
		return nil, err
	}
	var backupMetas []*BackupMeta
	for _, backupMetaJSON := range backupMetaJSONs {
		var backupMeta BackupMeta
		if err := utils.UnmarshalJSON(backupMetaJSON, &backupMeta); err != nil {
			return nil, err
		}
		backupMetas = append(backupMetas, &backupMeta)
	}
	return backupMetas, nil
}
