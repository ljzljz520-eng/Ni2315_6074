package archive

import (
	"encoding/json"
	"fmt"
	"strings"

	"independentweeklylog/internal/domain"
)

type snapshot struct {
	Entry     domain.JournalEntry `json:"entry"`
	Resources []domain.Resource   `json:"resources"`
}

func (s *Service) Verify(record domain.ArchiveRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	var value snapshot
	if err := json.Unmarshal([]byte(record.Snapshot), &value); err != nil {
		return fmt.Errorf("decode archive: %w", err)
	}
	checksum, err := Checksum(value)
	if err != nil {
		return err
	}
	if checksum != record.Checksum {
		return fmt.Errorf("archive checksum mismatch")
	}
	return nil
}

func (s *Service) Restore(record domain.ArchiveRecord) (domain.JournalEntry, []domain.Resource, error) {
	if err := s.Verify(record); err != nil {
		return domain.JournalEntry{}, nil, err
	}
	var value snapshot
	if err := json.Unmarshal([]byte(record.Snapshot), &value); err != nil {
		return domain.JournalEntry{}, nil, err
	}
	value.Entry.Status = domain.StatusArchived
	if err := s.repo.SaveEntry(value.Entry); err != nil {
		return domain.JournalEntry{}, nil, err
	}
	if err := s.repo.SaveResources(value.Resources); err != nil {
		return domain.JournalEntry{}, nil, err
	}
	return value.Entry, value.Resources, nil
}

func CanonicalReason(reason string) string { return strings.TrimSpace(reason) }

func (s *Service) Latest(entryID string) (domain.ArchiveRecord, error) {
	records, err := s.Get(entryID)
	if err != nil {
		return domain.ArchiveRecord{}, err
	}
	if len(records) == 0 {
		return domain.ArchiveRecord{}, domain.ErrNotFound
	}
	latest := records[0]
	for _, record := range records[1:] {
		if record.ArchivedAt.After(latest.ArchivedAt) {
			latest = record
		}
	}
	return latest, nil
}

func (s *Service) History(entryID string) ([]domain.ArchiveRecord, error) { return s.Get(entryID) }

func (s *Service) HasArchive(entryID string) (bool, error) {
	records, err := s.Get(entryID)
	return len(records) > 0, err
}

func ArchiveCount(records []domain.ArchiveRecord) int { return len(records) }
