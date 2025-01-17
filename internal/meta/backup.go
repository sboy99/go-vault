package meta

import (
	"time"

	"github.com/sboy99/go-vault/config"
	"github.com/sboy99/go-vault/pkg/boltdb"
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
	return boltdb.Save(_BACKUP_META, b.ID, backupMetaJSON)
}

func GetBackupMeta(id string) (*BackupMeta, error) {
	backupMetaJSON, err := boltdb.Get(_BACKUP_META, id)
	if err != nil {
		return nil, err
	}
	var backupMeta BackupMeta
	if err := utils.UnmarshalJSON(backupMetaJSON, &backupMeta); err != nil {
		return nil, err
	}
	return &backupMeta, nil
}

func ListBackupMeta(size, offset int) ([]*BackupMeta, error) {
	backupMetaJSONs, err := boltdb.List(_BACKUP_META, size, offset)
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

func DeleteBackupMeta(id string) error {
	return boltdb.Delete(_BACKUP_META, id)
}
