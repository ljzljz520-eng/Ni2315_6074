package review

import (
	"fmt"
	"strings"

	"independentweeklylog/internal/domain"
)

type Gate struct {
	MinAchievements        int
	RequireBlockerDecision bool
	RequireNextStep        int
}

func DefaultGate() Gate {
	return Gate{MinAchievements: 1, RequireBlockerDecision: true, RequireNextStep: 1}
}

func (g Gate) Check(entry domain.JournalEntry) []string {
	issues := make([]string, 0)
	if len(entry.Achievements) < g.MinAchievements {
		issues = append(issues, "record a meaningful achievement")
	}
	if len(entry.NextSteps) < g.RequireNextStep {
		issues = append(issues, "record a next step")
	}
	if g.RequireBlockerDecision && len(entry.Blockers) > 0 && strings.TrimSpace(entry.Summary) == "" {
		issues = append(issues, "explain blocker impact")
	}
	return issues
}

func (g Gate) Ready(entry domain.JournalEntry) bool { return len(g.Check(entry)) == 0 }

func RequireReviewer(entry domain.JournalEntry, reviewer string) error {
	if reviewer == "" {
		return fmt.Errorf("reviewer is required")
	}
	if strings.EqualFold(entry.Author, reviewer) {
		return fmt.Errorf("author cannot review own entry")
	}
	if entry.Status != domain.StatusSubmitted {
		return fmt.Errorf("entry must be submitted")
	}
	return nil
}

func DecisionLabel(decision domain.EntryStatus) string {
	switch decision {
	case domain.StatusApproved:
		return "Approved"
	case domain.StatusRejected:
		return "Changes requested"
	default:
		return "Pending"
	}
}

func SortReviews(records []domain.ReviewRecord) []domain.ReviewRecord {
	result := append([]domain.ReviewRecord(nil), records...)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].CreatedAt.Before(result[i].CreatedAt) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

func LatestDecision(records []domain.ReviewRecord) domain.EntryStatus {
	if len(records) == 0 {
		return ""
	}
	sorted := SortReviews(records)
	return sorted[len(sorted)-1].Decision
}

func Notes(records []domain.ReviewRecord) []string {
	notes := make([]string, 0, len(records))
	for _, record := range SortReviews(records) {
		note := strings.TrimSpace(record.Note)
		if note != "" {
			notes = append(notes, note)
		}
	}
	return notes
}
