package domain

import (
	"fmt"
	"strings"
)

type AccessPolicy struct {
	Reviewers  []string
	Editors    []string
	Archivists []string
}

func (p AccessPolicy) CanEdit(actor string, entry JournalEntry) bool {
	if actor == "" || !entry.IsEditable() {
		return false
	}
	if actor == entry.Author {
		return true
	}
	for _, editor := range p.Editors {
		if strings.EqualFold(editor, actor) {
			return true
		}
	}
	return false
}

func (p AccessPolicy) CanApprove(actor string) bool {
	for _, reviewer := range p.Reviewers {
		if strings.EqualFold(reviewer, actor) {
			return true
		}
	}
	return false
}

func (p AccessPolicy) CanArchive(actor string) bool {
	for _, archivist := range p.Archivists {
		if strings.EqualFold(archivist, actor) {
			return true
		}
	}
	return false
}

func (p AccessPolicy) Validate() error {
	if len(p.Reviewers) == 0 {
		return fmt.Errorf("at least one reviewer is required")
	}
	if len(p.Editors) == 0 {
		return fmt.Errorf("at least one editor is required")
	}
	if len(p.Archivists) == 0 {
		return fmt.Errorf("at least one archivist is required")
	}
	return nil
}

func (p AccessPolicy) Roles(actor string) []string {
	roles := make([]string, 0, 3)
	if contains(p.Reviewers, actor) {
		roles = append(roles, "reviewer")
	}
	if contains(p.Editors, actor) {
		roles = append(roles, "editor")
	}
	if contains(p.Archivists, actor) {
		roles = append(roles, "archivist")
	}
	return roles
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}
	return false
}
