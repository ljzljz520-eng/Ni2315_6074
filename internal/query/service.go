package query

import (
	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
)

type Service struct{ repo *repository.Repository }

func New(repo *repository.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Search(filter repository.EntryFilter) ([]domain.JournalEntry, error) {
	return s.repo.FindEntries(filter)
}

func (s *Service) Summary(week int) (domain.WeeklySummary, error) {
	entries, err := s.repo.Entries()
	if err != nil {
		return domain.WeeklySummary{}, err
	}
	return domain.BuildSummary(entries, week), nil
}

func (s *Service) Timeline(entryID string) ([]domain.DomainEvent, error) {
	events, err := s.repo.Events()
	if err != nil {
		return nil, err
	}
	result := make([]domain.DomainEvent, 0)
	for _, event := range events {
		if event.EntryID == entryID {
			result = append(result, event)
		}
	}
	return result, nil
}
