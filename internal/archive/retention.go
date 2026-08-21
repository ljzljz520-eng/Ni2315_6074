package archive

import (
	"sort"
	"strings"
	"time"

	"independentweeklylog/internal/domain"
)

type RetentionPolicy struct {
	KeepLatest   int
	KeepApproved bool
	Before       time.Time
}

func DefaultRetention() RetentionPolicy { return RetentionPolicy{KeepLatest: 3, KeepApproved: true} }

func (p RetentionPolicy) Select(records []domain.ArchiveRecord) []domain.ArchiveRecord {
	items := append([]domain.ArchiveRecord(nil), records...)
	sort.Slice(items, func(i, j int) bool { return items[i].ArchivedAt.After(items[j].ArchivedAt) })
	selected := make([]domain.ArchiveRecord, 0, len(items))
	for _, record := range items {
		if !p.Before.IsZero() && record.ArchivedAt.Before(p.Before) {
			continue
		}
		if p.KeepLatest > 0 && len(selected) >= p.KeepLatest {
			break
		}
		selected = append(selected, record)
	}
	return selected
}

func (p RetentionPolicy) Validate() bool { return p.KeepLatest != 0 || p.KeepApproved }

func Summarize(records []domain.ArchiveRecord) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		reason := strings.TrimSpace(record.Reason)
		if reason == "" {
			reason = "unspecified"
		}
		counts[reason]++
	}
	return counts
}

func NewArchiveID(entryID string, at time.Time) string {
	return entryID + "-" + at.UTC().Format("20060102T150405")
}

func IsFresh(record domain.ArchiveRecord, now time.Time, window time.Duration) bool {
	return !record.ArchivedAt.IsZero() && now.Sub(record.ArchivedAt) <= window
}
