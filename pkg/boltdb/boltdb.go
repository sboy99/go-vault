package boltdb

import (
	"github.com/sboy99/go-vault/pkg/logger"

	bolt "go.etcd.io/bbolt"
)

const (
	_BOLT_DB_PATH = "./go-vault.db"
	_BOLT_DB_MODE = 0600 // Read and write permission
)

var db *bolt.DB

func Connect() error {
	conn, err := bolt.Open(_BOLT_DB_PATH, _BOLT_DB_MODE, nil)
	if err != nil {
		return err
	}
	db = conn
	logger.Info("Connected to bolt db")
	return nil
}

func Disconnect() error {
	return db.Close()
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
		return b.Put([]byte(key), value)
	})
}

func Get(bucket, key string) ([]byte, error) {
	var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get([]byte(key))
		return nil
	})
	return value, err
}

func List(bucket string, size, offset int) ([][]byte, error) {
	var values [][]byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		i := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if i >= offset {
				values = append(values, v)
				if len(values) == size {
					break
				}
			}
			i++
		}
		return nil
	})
	return values, err
}

func Delete(bucket, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})
}
