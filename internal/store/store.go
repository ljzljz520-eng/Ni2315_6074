package store

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"

	"go.etcd.io/bbolt"
)

type Store struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	if err := os.MkdirAll(filepathDir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open bbolt: %w", err)
	}
	store := &Store{db: db, path: path}
	if err := store.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func filepathDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			if i == 0 {
				return "/"
			}
			return path[:i]
		}
	}
	return "."
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", string(name), err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func (s *Store) Put(bucket string, key string, value any) error {
	if bucketFor(bucket) == "" || key == "" {
		return errors.New("bucket and key are required")
	}
	data, err := encode(value)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Put([]byte(key), data)
	})
}

func (s *Store) Get(bucket string, key string, target any) error {
	if bucketFor(bucket) == "" || key == "" {
		return errors.New("bucket and key are required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if value == nil {
			return os.ErrNotExist
		}
		return decode(value, target)
	})
}

func (s *Store) Delete(bucket string, key string) error {
	if bucketFor(bucket) == "" || key == "" {
		return errors.New("bucket and key are required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete([]byte(key))
	})
}

func (s *Store) Keys(bucket string) ([]string, error) {
	if bucketFor(bucket) == "" {
		return nil, errors.New("unknown bucket")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	keys := make([]string, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(key, _ []byte) error {
			keys = append(keys, string(key))
			return nil
		})
	})
	sort.Strings(keys)
	return keys, err
}

func (s *Store) Update(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(fn)
}

func (s *Store) View(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.View(fn)
}
