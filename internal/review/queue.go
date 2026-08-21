package review

import (
	"sort"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/repository"
)

type QueueItem struct {
	Entry   domain.JournalEntry
	Reviews int
	Age     int
}

func (s *Service) Queue() ([]QueueItem, error) {
	entries, err := s.repo.FindEntries(repository.EntryFilter{Status: domain.StatusSubmitted})
	if err != nil {
		return nil, err
	}
	queue := make([]QueueItem, 0, len(entries))
	for _, entry := range entries {
		reviews, reviewErr := s.repo.ReviewsFor(entry.ID)
		if reviewErr != nil {
			return nil, reviewErr
		}
		queue = append(queue, QueueItem{Entry: entry, Reviews: len(reviews), Age: entry.Version})
	}
	sort.Slice(queue, func(i, j int) bool {
		if queue[i].Age == queue[j].Age {
			return queue[i].Entry.ID < queue[j].Entry.ID
		}
		return queue[i].Age < queue[j].Age
	})
	return queue, nil
}

func (s *Service) ReviewCount(entryID string) (int, error) {
	records, err := s.History(entryID)
	return len(records), err
}

func (s *Service) HasDecision(entryID string, decision domain.EntryStatus) (bool, error) {
	records, err := s.History(entryID)
	if err != nil {
		return false, err
	}
	for _, record := range records {
		if record.Decision == decision {
			return true, nil
		}
	}
	return false, nil
}
