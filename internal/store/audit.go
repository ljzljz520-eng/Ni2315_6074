package store

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

type AuditItem struct {
	Bucket  string `json:"bucket"`
	Key     string `json:"key"`
	Size    int    `json:"size"`
	Present bool   `json:"present"`
}

type AuditReport struct {
	Path        string      `json:"path"`
	GeneratedAt time.Time   `json:"generated_at"`
	Items       []AuditItem `json:"items"`
	TotalBytes  int         `json:"total_bytes"`
}

func (s *Store) Audit() (AuditReport, error) {
	report := AuditReport{Path: s.path, GeneratedAt: time.Unix(0, 0).UTC(), Items: make([]AuditItem, 0)}
	err := s.View(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			bucket := tx.Bucket(name)
			if bucket == nil {
				return fmt.Errorf("bucket %s is unavailable", name)
			}
			if err := bucket.ForEach(func(key, value []byte) error {
				if value == nil {
					return nil
				}
				item := AuditItem{Bucket: string(name), Key: string(key), Size: len(value), Present: true}
				report.Items = append(report.Items, item)
				report.TotalBytes += len(value)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})
	sort.Slice(report.Items, func(i, j int) bool {
		if report.Items[i].Bucket == report.Items[j].Bucket {
			return report.Items[i].Key < report.Items[j].Key
		}
		return report.Items[i].Bucket < report.Items[j].Bucket
	})
	return report, err
}

func (r AuditReport) Encode() ([]byte, error) { return json.Marshal(r) }

func (r AuditReport) Empty() bool { return len(r.Items) == 0 }

func (r AuditReport) Buckets() []string {
	seen := make(map[string]bool)
	for _, item := range r.Items {
		seen[item.Bucket] = true
	}
	result := make([]string, 0, len(seen))
	for bucket := range seen {
		result = append(result, bucket)
	}
	sort.Strings(result)
	return result
}

func (r AuditReport) HasKey(bucket, key string) bool {
	for _, item := range r.Items {
		if strings.EqualFold(item.Bucket, bucket) && item.Key == key {
			return true
		}
	}
	return false
}

func (s *Store) Export(bucket string) (map[string]json.RawMessage, error) {
	if bucketFor(bucket) == "" {
		return nil, fmt.Errorf("unknown bucket")
	}
	result := make(map[string]json.RawMessage)
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(key, value []byte) error {
			if value != nil {
				result[string(key)] = append(json.RawMessage(nil), value...)
			}
			return nil
		})
	})
	return result, err
}

func (s *Store) Import(bucket string, values map[string]json.RawMessage) error {
	if bucketFor(bucket) == "" {
		return fmt.Errorf("unknown bucket")
	}
	return s.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		for key, value := range values {
			if strings.TrimSpace(key) == "" {
				continue
			}
			if err := b.Put([]byte(key), value); err != nil {
				return err
			}
		}
		return nil
	})
}
