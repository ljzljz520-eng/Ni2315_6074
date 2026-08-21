package review

import (
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

func (s *Service) Decide(entryID, reviewer string, decision domain.EntryStatus, note string) (domain.JournalEntry, domain.ReviewRecord, error) {
	entry, err := s.repo.Entry(entryID)
	if err != nil {
		return entry, domain.ReviewRecord{}, err
	}
	if entry.Status != domain.StatusSubmitted {
		return entry, domain.ReviewRecord{}, fmt.Errorf("%w: entry is not submitted", domain.ErrConflict)
	}
	review := domain.ReviewRecord{ID: fmt.Sprintf("review-%s-%s", entryID, reviewer), EntryID: entryID, Reviewer: reviewer, Decision: decision, Note: note, CreatedAt: s.clock()}
	if err := review.Validate(); err != nil {
		return entry, review, err
	}
	entry.Status = decision
	entry.Version++
	entry.UpdatedAt = s.clock()
	if err := s.repo.SaveReview(review); err != nil {
		return entry, review, err
	}
	if err := s.repo.SaveEntry(entry); err != nil {
		return entry, review, err
	}
	return entry, review, nil
}

func (s *Service) History(entryID string) ([]domain.ReviewRecord, error) {
	return s.repo.ReviewsFor(entryID)
}

func CanReview(entry domain.JournalEntry, reviewer string) bool {
	return entry.Status == domain.StatusSubmitted && reviewer != "" && reviewer != entry.Author
}
