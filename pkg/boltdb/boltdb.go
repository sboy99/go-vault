package boltdb

import (
	"fmt"
	"sync"

	bolt "go.etcd.io/bbolt"
)

const _BOLT_DB_MODE = 0600

var (
	db     *bolt.DB
	dbPath = "./go-vault.db"
	mu     sync.Mutex
)

// SetPath configures the BoltDB file path. Must be called before Connect.
func SetPath(path string) {
	mu.Lock()
	defer mu.Unlock()
	if path != "" {
		dbPath = path
	}
}

func Connect() error {
	mu.Lock()
	defer mu.Unlock()
	conn, err := bolt.Open(dbPath, _BOLT_DB_MODE, nil)
	if err != nil {
		return fmt.Errorf("open boltdb %s: %w", dbPath, err)
	}
	db = conn
	return nil
}

func Disconnect() error {
	mu.Lock()
	defer mu.Unlock()
	if db == nil {
		return nil
	}
	err := db.Close()
	db = nil
	return err
}

func GetDB() *bolt.DB {
	return db
}

func CreateBucket(bucket string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}

func Save(bucket, key string, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		return b.Put([]byte(key), value)
	})
}

func Get(bucket, key string) ([]byte, error) {
	var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		data := b.Get([]byte(key))
		if data != nil {
			value = append([]byte(nil), data...)
		}
		return nil
	})
	return value, err
}

// List returns values in key order (oldest first for time-sortable keys).
func List(bucket string, size, offset int) ([][]byte, error) {
	var values [][]byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		c := b.Cursor()
		i := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if i >= offset {
				values = append(values, append([]byte(nil), v...))
				if size > 0 && len(values) == size {
					break
				}
			}
			i++
		}
		return nil
	})
	return values, err
}

// ListReverse returns values newest-first for time-sortable keys.
func ListReverse(bucket string, size, offset int) ([][]byte, error) {
	var values [][]byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		c := b.Cursor()
		i := 0
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if i >= offset {
				values = append(values, append([]byte(nil), v...))
				if size > 0 && len(values) == size {
					break
				}
			}
			i++
		}
		return nil
	})
	return values, err
}

// ListAll returns every value in key order.
func ListAll(bucket string) ([][]byte, error) {
	return List(bucket, 0, 0)
}

func Delete(bucket, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		return b.Delete([]byte(key))
	})
}
