package repository

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"independentweeklylog/internal/domain"
	"independentweeklylog/internal/store"
)

type Repository struct{ db *store.Store }

func New(db *store.Store) *Repository { return &Repository{db: db} }

func (r *Repository) SaveEntry(entry domain.JournalEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	return r.db.Put(store.EntriesBucket, entry.ID, entry)
}

func (r *Repository) Entry(id string) (domain.JournalEntry, error) {
	var entry domain.JournalEntry
	err := r.db.Get(store.EntriesBucket, id, &entry)
	if errors.Is(err, os.ErrNotExist) {
		return entry, domain.ErrNotFound
	}
	return entry, err
}

func (r *Repository) DeleteEntry(id string) error { return r.db.Delete(store.EntriesBucket, id) }

func (r *Repository) Entries() ([]domain.JournalEntry, error) {
	keys, err := r.db.Keys(store.EntriesBucket)
	if err != nil {
		return nil, err
	}
	entries := make([]domain.JournalEntry, 0, len(keys))
	for _, key := range keys {
		entry, getErr := r.Entry(key)
		if getErr != nil {
			return nil, getErr
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Week == entries[j].Week {
			return entries[i].ID < entries[j].ID
		}
		return entries[i].Week < entries[j].Week
	})
	return entries, nil
}

func (r *Repository) SaveReview(review domain.ReviewRecord) error {
	if err := review.Validate(); err != nil {
		return err
	}
	return r.db.Put(store.ReviewsBucket, review.ID, review)
}

func (r *Repository) ReviewsFor(entryID string) ([]domain.ReviewRecord, error) {
	keys, err := r.db.Keys(store.ReviewsBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ReviewRecord, 0)
	for _, key := range keys {
		var review domain.ReviewRecord
		if err := r.db.Get(store.ReviewsBucket, key, &review); err != nil {
			return nil, err
		}
		if review.EntryID == entryID {
			result = append(result, review)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
	return result, nil
}

func (r *Repository) SaveResource(resource domain.Resource) error {
	if err := resource.Validate(); err != nil {
		return err
	}
	return r.db.Put(store.ResourcesBucket, resource.ID, resource)
}

func (r *Repository) Resource(id string) (domain.Resource, error) {
	var resource domain.Resource
	err := r.db.Get(store.ResourcesBucket, id, &resource)
	if errors.Is(err, os.ErrNotExist) {
		return resource, domain.ErrNotFound
	}
	return resource, err
}

func (r *Repository) ResourcesFor(entryID string) ([]domain.Resource, error) {
	keys, err := r.db.Keys(store.ResourcesBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Resource, 0)
	for _, key := range keys {
		resource, getErr := r.Resource(key)
		if getErr != nil {
			return nil, getErr
		}
		if resource.EntryID == entryID {
			result = append(result, resource)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (r *Repository) SaveArchive(record domain.ArchiveRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return r.db.Put(store.ArchivesBucket, record.ID, record)
}

func (r *Repository) ArchivesFor(entryID string) ([]domain.ArchiveRecord, error) {
	keys, err := r.db.Keys(store.ArchivesBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ArchiveRecord, 0)
	for _, key := range keys {
		var record domain.ArchiveRecord
		if err := r.db.Get(store.ArchivesBucket, key, &record); err != nil {
			return nil, err
		}
		if record.EntryID == entryID {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ArchivedAt.Before(result[j].ArchivedAt) })
	return result, nil
}

func (r *Repository) SaveProfile(profile domain.DeveloperProfile) error {
	if profile.ID == "" || profile.DisplayName == "" {
		return fmt.Errorf("%w: profile is incomplete", domain.ErrInvalid)
	}
	return r.db.Put(store.ProfilesBucket, profile.ID, profile)
}

func (r *Repository) Profile(id string) (domain.DeveloperProfile, error) {
	var profile domain.DeveloperProfile
	err := r.db.Get(store.ProfilesBucket, id, &profile)
	if errors.Is(err, os.ErrNotExist) {
		return profile, domain.ErrNotFound
	}
	return profile, err
}
