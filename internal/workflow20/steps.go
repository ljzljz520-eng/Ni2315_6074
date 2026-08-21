package workflow20

import (
	"fmt"
	"strings"

	"independentweeklylog/internal/domain"
)

type Step struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

func DefaultSteps() []Step {
	return []Step{
		{Number: 1, Name: "记录本周目标"},
		{Number: 2, Name: "整理实现进展"},
		{Number: 3, Name: "标记阻塞事项"},
		{Number: 4, Name: "准备下周计划"},
		{Number: 5, Name: "绑定可追溯资源"},
		{Number: 6, Name: "提交同行审核"},
		{Number: 7, Name: "关闭父子资源并保存"},
		{Number: 8, Name: "归档可复用结论"},
	}
}

func CompleteStep(steps []Step, number int) ([]Step, error) {
	updated := append([]Step(nil), steps...)
	found := false
	for index := range updated {
		if updated[index].Number == number {
			if index > 0 && !updated[index-1].Complete {
				return steps, fmt.Errorf("step %d must follow step %d", number, updated[index-1].Number)
			}
			updated[index].Complete = true
			found = true
		}
	}
	if !found {
		return steps, fmt.Errorf("unknown workflow step %d", number)
	}
	return updated, nil
}

func ValidateWorkflow(entry domain.JournalEntry, steps []Step) error {
	if entry.Status == domain.StatusDraft && len(steps) == 0 {
		return fmt.Errorf("draft needs workflow steps")
	}
	for i, step := range steps {
		if step.Number != i+1 {
			return fmt.Errorf("workflow steps must be contiguous")
		}
		if strings.TrimSpace(step.Name) == "" {
			return fmt.Errorf("workflow step name is required")
		}
	}
	return nil
}
