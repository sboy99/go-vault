package boltdb

import (
	"fmt"
	"sync"

	bolt "go.etcd.io/bbolt"
)

const _BOLT_DB_MODE = 0600

// Client is an instance-scoped BoltDB wrapper.
type Client struct {
	db   *bolt.DB
	path string
	mu   sync.Mutex
}

// New opens a BoltDB file at path and returns a Client.
func New(path string) (*Client, error) {
	if path == "" {
		path = "./go-vault.db"
	}
	conn, err := bolt.Open(path, _BOLT_DB_MODE, nil)
	if err != nil {
		return nil, fmt.Errorf("open boltdb %s: %w", path, err)
	}
	return &Client{db: conn, path: path}, nil
}

// Close closes the underlying BoltDB connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	c.db = nil
	return err
}

// CreateBucket creates a bucket if it does not exist.
func (c *Client) CreateBucket(bucket string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}

func (c *Client) Save(bucket, key string, value []byte) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		return b.Put([]byte(key), value)
	})
}

func (c *Client) Get(bucket, key string) ([]byte, error) {
	var value []byte
	err := c.db.View(func(tx *bolt.Tx) error {
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
func (c *Client) List(bucket string, size, offset int) ([][]byte, error) {
	var values [][]byte
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		cur := b.Cursor()
		i := 0
		for k, v := cur.First(); k != nil; k, v = cur.Next() {
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
func (c *Client) ListReverse(bucket string, size, offset int) ([][]byte, error) {
	var values [][]byte
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		cur := b.Cursor()
		i := 0
		for k, v := cur.Last(); k != nil; k, v = cur.Prev() {
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
func (c *Client) ListAll(bucket string) ([][]byte, error) {
	return c.List(bucket, 0, 0)
}

func (c *Client) Delete(bucket, key string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		return b.Delete([]byte(key))
	})
}
