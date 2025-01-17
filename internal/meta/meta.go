package meta

import (
	"fmt"

	"github.com/sboy99/go-vault/pkg/boltdb"
)

const (
	_BACKUP_META  = "backup_meta"
	_RESTORE_META = "restore_meta"
)

func Init() error {
	if err := boltdb.Connect(); err != nil {
		return err
	}
	buckets := []string{_BACKUP_META, _RESTORE_META}
	return createBucketsIfNotExists(buckets)
}

func Cleanup() error {
	return boltdb.Disconnect()
}

func createBucketsIfNotExists(buckets []string) error {
	goroutineCount := len(buckets)
	errChannel := make(chan error, len(buckets))
	doneChannel := make(chan struct{}, len(buckets))

	for _, bucket := range buckets {
		go func(bucket string) {
			errChannel <- boltdb.CreateBucket(bucket)
			doneChannel <- struct{}{}
		}(bucket)
	}

	// Collect errors from all goroutines
	var combinedErrors error
	for {
		select {
		case err := <-errChannel:
			if err == nil {
				combinedErrors = err
			} else {
				combinedErrors = fmt.Errorf("%v; %v", combinedErrors, err)
			}
		case <-doneChannel:
			goroutineCount--
			if goroutineCount == 0 {
				close(errChannel)
				close(doneChannel)
				return combinedErrors
			}
		}
	}
}
