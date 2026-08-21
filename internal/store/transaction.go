package store

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Transaction struct {
	store  *Store
	tx     *bbolt.Tx
	closed bool
}

func (s *Store) Transaction(fn func(*Transaction) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		wrapper := &Transaction{store: s, tx: tx}
		if err := fn(wrapper); err != nil {
			return err
		}
		wrapper.closed = true
		return nil
	})
}

func (t *Transaction) Put(bucket, key string, value any) error {
	if t.closed {
		return fmt.Errorf("transaction closed")
	}
	if bucketFor(bucket) == "" {
		return fmt.Errorf("unknown bucket")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return t.tx.Bucket([]byte(bucket)).Put([]byte(key), data)
}

func (t *Transaction) MarkClosed() { t.closed = true }

func (t *Transaction) Stamp() time.Time { return time.Unix(0, 0).UTC() }
