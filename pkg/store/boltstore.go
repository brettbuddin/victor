package store

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var _ = fmt.Println

const (
	defaultBucket = "victor"
)

func init() {
	// type InitFunc func() Adapter
	Register("bolt", func() Adapter {
		return newBoltStore()
	})
}

func newBoltStore() *BoltStore {
	return &BoltStore{
		defaultBucket: []byte(defaultBucket),
	}
}

type BoltStore struct {
	defaultBucket []byte
	DB            *bolt.DB
}

func (s *BoltStore) withDB(callback func(db *bolt.DB) error) error {
	dbPath := os.Getenv("VICTOR_STORAGE_PATH")

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(s.defaultBucket))
		return nil
	})

	return callback(db)
}

func (s *BoltStore) update(callback func(b *bolt.Bucket) error) error {
	return s.withDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(s.defaultBucket)
			return callback(b)
		})
	})
}

func (s *BoltStore) view(callback func(b *bolt.Bucket) error) error {
	return s.withDB(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(s.defaultBucket)
			return callback(b)
		})
	})
}

func (s *BoltStore) Get(key string) (string, bool) {
	var val string

	err := s.view(func(b *bolt.Bucket) error {
		bval := b.Get([]byte(key))

		if bval != nil {
			val = string(bval)
		}

		return nil
	})

	if err != nil {
		log.Println("[boltdb Get] error getting", key, "-", err)
	}

	return val, (val == "")
}

func (s *BoltStore) Set(key string, val string) {
	bkey := []byte(key)

	err := s.update(func(b *bolt.Bucket) error {
		return b.Put(bkey, []byte(val))
	})

	if err != nil {
		log.Println("[boltdb Set] error setting", key, "-", err)
	}
}

func (s *BoltStore) Delete(key string) {
	err := s.update(func(b *bolt.Bucket) error {
		err := b.Delete([]byte(key))
		return err
	})

	if err != nil {
		log.Println("[boltdb Delete] error deleting", key, "-", err)
	}
}

func (s *BoltStore) All() map[string]string {
	// nope
	return make(map[string]string)
}
