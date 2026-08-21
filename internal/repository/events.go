package repository

import (
	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/store"
)

func (r *Repository) SaveEvent(event domain.DomainEvent) error {
	if event.ID == "" || event.EntryID == "" {
		return domain.ErrInvalid
	}
	return r.db.Put(store.EventsBucket, event.ID, event)
}

func (r *Repository) Events() ([]domain.DomainEvent, error) {
	keys, err := r.db.Keys(store.EventsBucket)
	if err != nil {
		return nil, err
	}
	events := make([]domain.DomainEvent, 0, len(keys))
	for _, key := range keys {
		var event domain.DomainEvent
		if err := r.db.Get(store.EventsBucket, key, &event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
