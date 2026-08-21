package store

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Snapshot struct {
	TakenAt time.Time      `json:"taken_at"`
	Buckets map[string]int `json:"buckets"`
}

func (s *Store) Snapshot() (Snapshot, error) {
	snapshot := Snapshot{TakenAt: time.Unix(0, 0).UTC(), Buckets: make(map[string]int)}
	err := s.View(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			bucket := tx.Bucket(name)
			if bucket == nil {
				return fmt.Errorf("missing bucket %s", name)
			}
			count := 0
			err := bucket.ForEach(func(_, _ []byte) error { count++; return nil })
			if err != nil {
				return err
			}
			snapshot.Buckets[string(name)] = count
		}
		return nil
	})
	return snapshot, err
}

func (s Snapshot) Encode() ([]byte, error) { return json.Marshal(s) }
