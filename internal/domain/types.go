package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNotFound  = errors.New("weekly log not found")
	ErrInvalid   = errors.New("invalid weekly log")
	ErrConflict  = errors.New("weekly log conflict")
	ErrForbidden = errors.New("operation not permitted")
)

type EntryStatus string

const (
	StatusDraft     EntryStatus = "draft"
	StatusSubmitted EntryStatus = "submitted"
	StatusApproved  EntryStatus = "approved"
	StatusRejected  EntryStatus = "rejected"
	StatusArchived  EntryStatus = "archived"
)

type JournalEntry struct {
	ID           string      `json:"id"`
	Week         int         `json:"week"`
	Title        string      `json:"title"`
	Summary      string      `json:"summary"`
	Achievements []string    `json:"achievements"`
	Blockers     []string    `json:"blockers"`
	NextSteps    []string    `json:"next_steps"`
	Status       EntryStatus `json:"status"`
	Author       string      `json:"author"`
	Tags         []string    `json:"tags"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Version      int         `json:"version"`
}

type ReviewRecord struct {
	ID        string      `json:"id"`
	EntryID   string      `json:"entry_id"`
	Reviewer  string      `json:"reviewer"`
	Decision  EntryStatus `json:"decision"`
	Note      string      `json:"note"`
	CreatedAt time.Time   `json:"created_at"`
}

type Resource struct {
	ID       string `json:"id"`
	EntryID  string `json:"entry_id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	ParentID string `json:"parent_id"`
	Payload  string `json:"payload"`
	Required bool   `json:"required"`
	Closed   bool   `json:"closed"`
}

type ArchiveRecord struct {
	ID         string    `json:"id"`
	EntryID    string    `json:"entry_id"`
	Reason     string    `json:"reason"`
	Snapshot   string    `json:"snapshot"`
	ArchivedAt time.Time `json:"archived_at"`
	Checksum   string    `json:"checksum"`
}

type DeveloperProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Timezone    string `json:"timezone"`
	Active      bool   `json:"active"`
}

type WeeklySummary struct {
	Week     int      `json:"week"`
	Total    int      `json:"total"`
	Approved int      `json:"approved"`
	Pending  int      `json:"pending"`
	Archived int      `json:"archived"`
	Authors  []string `json:"authors"`
}

func (e JournalEntry) Validate() error {
	if strings.TrimSpace(e.ID) == "" || strings.TrimSpace(e.Title) == "" {
		return fmt.Errorf("%w: id and title are required", ErrInvalid)
	}
	if e.Week < 1 || e.Week > 53 {
		return fmt.Errorf("%w: week must be between 1 and 53", ErrInvalid)
	}
	if strings.TrimSpace(e.Author) == "" {
		return fmt.Errorf("%w: author is required", ErrInvalid)
	}
	if len(strings.TrimSpace(e.Summary)) < 12 {
		return fmt.Errorf("%w: summary is too short", ErrInvalid)
	}
	if len(e.Achievements) == 0 {
		return fmt.Errorf("%w: add at least one achievement", ErrInvalid)
	}
	if len(e.NextSteps) == 0 {
		return fmt.Errorf("%w: add at least one next step", ErrInvalid)
	}
	return nil
}

func (e JournalEntry) IsEditable() bool {
	return e.Status == StatusDraft || e.Status == StatusRejected
}

func (e JournalEntry) IsVisible() bool {
	return e.Status == StatusApproved || e.Status == StatusArchived
}

func (e JournalEntry) Clone() JournalEntry {
	copyEntry := e
	copyEntry.Achievements = append([]string(nil), e.Achievements...)
	copyEntry.Blockers = append([]string(nil), e.Blockers...)
	copyEntry.NextSteps = append([]string(nil), e.NextSteps...)
	copyEntry.Tags = append([]string(nil), e.Tags...)
	return copyEntry
}

func (r ReviewRecord) Validate() error {
	if r.EntryID == "" || r.Reviewer == "" {
		return fmt.Errorf("%w: review identity is required", ErrInvalid)
	}
	if r.Decision != StatusApproved && r.Decision != StatusRejected {
		return fmt.Errorf("%w: decision must approve or reject", ErrInvalid)
	}
	if strings.TrimSpace(r.Note) == "" {
		return fmt.Errorf("%w: review note is required", ErrInvalid)
	}
	return nil
}

func (r Resource) Validate() error {
	if r.ID == "" || r.EntryID == "" || r.Name == "" {
		return fmt.Errorf("%w: resource identity is required", ErrInvalid)
	}
	if r.ParentID == r.ID {
		return fmt.Errorf("%w: resource cannot parent itself", ErrInvalid)
	}
	return nil
}

func (a ArchiveRecord) Validate() error {
	if a.ID == "" || a.EntryID == "" || a.Snapshot == "" || a.Checksum == "" {
		return fmt.Errorf("%w: archive record is incomplete", ErrInvalid)
	}
	return nil
}
