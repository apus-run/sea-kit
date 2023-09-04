// Package bolt providers a cache driver that stores the cache using BoltDB.
package bolt

import (
	"encoding/json"
	"errors"
	"time"

	bt "go.etcd.io/bbolt"

	"github.com/apus-run/sea-kit/cache"
)

var boltBucket = []byte("cache")

type (
	bolt struct {
		db *bt.DB
	}

	boltContent struct {
		Duration int64  `json:"duration"`
		Data     string `json:"data,omitempty"`
	}
)

// New creates an instance of BoltDB cache
func New(db *bt.DB) cache.Cache {
	return &bolt{db}
}

func (b *bolt) read(key string) (*boltContent, error) {
	var value []byte

	err := b.db.View(func(tx *bt.Tx) error {
		if bucket := tx.Bucket(boltBucket); bucket != nil {
			value = bucket.Get([]byte(key))
			return nil
		}

		return errors.New("bucket not found")
	})
	if err != nil {
		return nil, err
	}

	content := &boltContent{}
	if err := json.Unmarshal(value, content); err != nil {
		return nil, err
	}

	if content.Duration == 0 {
		return content, nil
	}

	if content.Duration <= time.Now().Unix() {
		_ = b.Delete(key)
		return nil, cache.ErrCacheExpired
	}

	return content, nil
}

// Contains checks if the cached key exists into the BoltDB storage
func (b *bolt) Contains(key string) bool {
	_, err := b.read(key)
	return err == nil
}

// Delete the cached key from BoltDB storage
func (b *bolt) Delete(key string) error {
	return b.db.Update(func(tx *bt.Tx) error {
		if bucket := tx.Bucket(boltBucket); bucket != nil {
			return bucket.Delete([]byte(key))
		}

		return errors.New("bucket not found")
	})
}

// Fetch retrieves the cached value from key of the BoltDB storage
func (b *bolt) Fetch(key string) (string, error) {
	content, err := b.read(key)
	if err != nil {
		return "", err
	}

	return content.Data, nil
}

// FetchMulti retrieve multiple cached values from keys of the BoltDB storage
func (b *bolt) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := b.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the BoltDB storage
func (b *bolt) Flush() error {
	return b.db.Update(func(tx *bt.Tx) error {
		return tx.DeleteBucket(boltBucket)
	})
}

// Save a value in BoltDB storage by key
func (b *bolt) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &boltContent{duration, value}

	data, err := json.Marshal(content)
	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(boltBucket)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), data)
	})
}
