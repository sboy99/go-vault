package meta

import (
	"fmt"

	"github.com/sboy99/go-vault/pkg/boltdb"
)

const (
	_BACKUP_META  = "backup_meta"
	_RESTORE_META = "restore_meta"
	_JOB_META     = "job_meta"
)

func Init(metaDBPath string) error {
	if metaDBPath != "" {
		boltdb.SetPath(metaDBPath)
	}
	if err := boltdb.Connect(); err != nil {
		return err
	}
	buckets := []string{_BACKUP_META, _RESTORE_META, _JOB_META}
	return createBucketsIfNotExists(buckets)
}

func Cleanup() error {
	return boltdb.Disconnect()
}

func createBucketsIfNotExists(buckets []string) error {
	for _, bucket := range buckets {
		if err := boltdb.CreateBucket(bucket); err != nil {
			return fmt.Errorf("create bucket %s: %w", bucket, err)
		}
	}
	return nil
}

// BucketJob exposes the job bucket name for the job package.
func BucketJob() string { return _JOB_META }
