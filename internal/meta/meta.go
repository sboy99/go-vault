package meta

import (
	bolt "go.etcd.io/bbolt"
)

const (
	_BOLT_DB_PATH        = "./go-vault.db"
	_BOLT_DB_MODE        = 0600 // Read and write permission
	_BOLT_BACKUP_BUCKET  = "backups"
	_BOLT_RESTORE_BUCKET = "restores"
)

var (
	db      *bolt.DB
	buckets = []string{_BOLT_BACKUP_BUCKET, _BOLT_RESTORE_BUCKET}
)

func Init() error {
	// Init bolt db
	conn, err := bolt.Open(_BOLT_DB_PATH, _BOLT_DB_MODE, nil)
	if err != nil {
		return err
	}
	db = conn
	// Create buckets
	if err := createBucketsIfNotExists(); err != nil {
		return err
	}
	return nil
}

func Destroy() {
	db.Close()
}

func createBucketsIfNotExists() error {
	return db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func getMetaData(bucket string, key string) ([]byte, error) {
	var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get([]byte(key))
		return nil
	})
	return value, err
}

func listMetaData(bucket string, size int) ([][]byte, error) {
	var values [][]byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			values = append(values, v)
			if len(values) == size {
				break
			}
		}
		return nil
	})
	return values, err
}

func saveMetaData(bucket string, key string, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), value)
	})
}

func deleteMetaData(bucket string, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})
}
