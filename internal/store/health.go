package store

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type Health struct {
	Open    bool   `json:"open"`
	Path    string `json:"path"`
	Buckets int    `json:"buckets"`
}

func (s *Store) Health() (Health, error) {
	if s == nil {
		return Health{}, fmt.Errorf("store is nil")
	}
	if s.db == nil {
		return Health{Open: false, Path: s.path}, nil
	}
	count := 0
	err := s.View(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if tx.Bucket(name) == nil {
				return fmt.Errorf("bucket %s missing", name)
			}
		}
		return nil
	})
	if err != nil {
		return Health{}, err
	}
	for range bucketNames {
		count++
	}
	return Health{Open: true, Path: s.path, Buckets: count}, nil
}
