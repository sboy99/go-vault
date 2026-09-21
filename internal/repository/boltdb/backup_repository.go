package boltdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sboy99/go-vault/internal/domain"
	boltclient "github.com/sboy99/go-vault/pkg/boltdb"
)

const bucketBackup = "backup_meta"

// BackupRepository persists Backup entities in BoltDB.
type BackupRepository struct {
	client *boltclient.Client
}

var _ domain.BackupRepository = (*BackupRepository)(nil)

func NewBackupRepository(client *boltclient.Client) (*BackupRepository, error) {
	if client == nil {
		return nil, errors.New("boltdb client is required")
	}
	if err := client.CreateBucket(bucketBackup); err != nil {
		return nil, fmt.Errorf("create backup bucket: %w", err)
	}
	return &BackupRepository{client: client}, nil
}

func (r *BackupRepository) Save(ctx context.Context, b *domain.Backup) error {
	_ = ctx
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return r.client.Save(bucketBackup, b.BackupId, data)
}

func (r *BackupRepository) FindByID(ctx context.Context, id string) (*domain.Backup, error) {
	_ = ctx
	data, err := r.client.Get(bucketBackup, id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, domain.ErrBackupNotFound
	}
	var b domain.Backup
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BackupRepository) FindAll(ctx context.Context) ([]*domain.Backup, error) {
	_ = ctx
	rows, err := r.client.ListAll(bucketBackup)
	if err != nil {
		return nil, err
	}
	return decodeBackupList(rows)
}

func (r *BackupRepository) List(ctx context.Context, limit, offset int) ([]*domain.Backup, error) {
	_ = ctx
	rows, err := r.client.ListReverse(bucketBackup, limit, offset)
	if err != nil {
		return nil, err
	}
	return decodeBackupList(rows)
}

func (r *BackupRepository) Delete(ctx context.Context, id string) error {
	_ = ctx
	return r.client.Delete(bucketBackup, id)
}

func decodeBackupList(rows [][]byte) ([]*domain.Backup, error) {
	out := make([]*domain.Backup, 0, len(rows))
	for _, row := range rows {
		var b domain.Backup
		if err := json.Unmarshal(row, &b); err != nil {
			return nil, err
		}
		out = append(out, &b)
	}
	return out, nil
}
