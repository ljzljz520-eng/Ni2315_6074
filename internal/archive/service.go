package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
)

type Service struct {
	repo  *repository.Repository
	clock func() time.Time
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo, clock: func() time.Time { return time.Unix(0, 0).UTC() }}
}

func (s *Service) WithClock(clock func() time.Time) *Service {
	if clock != nil {
		s.clock = clock
	}
	return s
}

func Checksum(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func (s *Service) Archive(entryID, reason string) (domain.ArchiveRecord, error) {
	entry, err := s.repo.Entry(entryID)
	if err != nil {
		return domain.ArchiveRecord{}, err
	}
	if entry.Status != domain.StatusApproved {
		return domain.ArchiveRecord{}, fmt.Errorf("%w: only approved entries can be archived", domain.ErrConflict)
	}
	resources, err := s.repo.ResourcesFor(entryID)
	if err != nil {
		return domain.ArchiveRecord{}, err
	}
	if len(resources) < 2 {
		return domain.ArchiveRecord{}, fmt.Errorf("%w: archive requires resource evidence", domain.ErrInvalid)
	}
	snapshot := struct {
		Entry     domain.JournalEntry `json:"entry"`
		Resources []domain.Resource   `json:"resources"`
	}{Entry: entry, Resources: resources}
	snapshotData, err := json.Marshal(snapshot)
	if err != nil {
		return domain.ArchiveRecord{}, err
	}
	checksum, err := Checksum(snapshot)
	if err != nil {
		return domain.ArchiveRecord{}, err
	}
	record := domain.ArchiveRecord{ID: "archive-" + entryID, EntryID: entryID, Reason: reason, Snapshot: string(snapshotData), ArchivedAt: s.clock(), Checksum: checksum}
	if err := s.repo.SaveArchive(record); err != nil {
		return record, err
	}
	entry.Status = domain.StatusArchived
	entry.UpdatedAt = s.clock()
	entry.Version++
	return record, s.repo.SaveEntry(entry)
}

func (s *Service) Get(entryID string) ([]domain.ArchiveRecord, error) {
	return s.repo.ArchivesFor(entryID)
}
