package meta

import (
	"fmt"
	"sync"

	"github.com/sboy99/go-vault/pkg/boltdb"
)

const (
	_BACKUP_META = "backup_meta"
	_JOB_META    = "job_meta"
)

var (
	clientMu sync.Mutex
	client   *boltdb.Client
)

func Init(metaDBPath string) error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if client != nil {
		_ = client.Close()
		client = nil
	}

	c, err := boltdb.New(metaDBPath)
	if err != nil {
		return err
	}
	client = c

	buckets := []string{_BACKUP_META, _JOB_META}
	return createBucketsIfNotExists(buckets)
}

func Cleanup() error {
	clientMu.Lock()
	defer clientMu.Unlock()
	if client == nil {
		return nil
	}
	err := client.Close()
	client = nil
	return err
}

func createBucketsIfNotExists(buckets []string) error {
	for _, bucket := range buckets {
		if err := client.CreateBucket(bucket); err != nil {
			return fmt.Errorf("create bucket %s: %w", bucket, err)
		}
	}
	return nil
}

// DB returns the temporary package-level BoltDB client.
// Removed once repositories take an injected client.
func DB() *boltdb.Client {
	clientMu.Lock()
	defer clientMu.Unlock()
	return client
}

// BucketJob exposes the job bucket name for the job package.
func BucketJob() string { return _JOB_META }

// BucketBackup exposes the backup bucket name.
func BucketBackup() string { return _BACKUP_META }
