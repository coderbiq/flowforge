package tracker_test

import (
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/tracker"
)

func TestParseIssueFile(t *testing.T) {
	tmpDir := t.TempDir()
	issuePath := filepath.Join(tmpDir, "01-auth-db.md")

	content := `# 01: Setup Auth Database

**Status:** open
**Type:** task
**Blocked by:** 00, 99
**Assignee:** alice
**Labels:** backend, security

## What to build
Setup PostgreSQL tables.

## Acceptance criteria
- [ ] Table exists
`
	if err := os.WriteFile(issuePath, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	issue, err := tracker.ParseIssueFile(issuePath)
	if err != nil {
		t.Fatalf("ParseIssueFile failed: %v", err)
	}

	if issue.ID != "01" {
		t.Errorf("expected ID '01', got %q", issue.ID)
	}
	if issue.Title != "Setup Auth Database" {
		t.Errorf("expected Title 'Setup Auth Database', got %q", issue.Title)
	}
	if issue.Status != tracker.StatusOpen {
		t.Errorf("expected Status 'open', got %q", issue.Status)
	}
	if len(issue.BlockedBy) != 2 || issue.BlockedBy[0] != "00" || issue.BlockedBy[1] != "99" {
		t.Errorf("expected BlockedBy ['00', '99'], got %v", issue.BlockedBy)
	}
	if issue.Assignee != "alice" {
		t.Errorf("expected Assignee 'alice', got %q", issue.Assignee)
	}
	if len(issue.Labels) != 2 || issue.Labels[0] != "backend" {
		t.Errorf("expected Labels ['backend', 'security'], got %v", issue.Labels)
	}
}

func TestDAGCycleAndFrontier(t *testing.T) {
	issues := []*tracker.Issue{
		{ID: "01", Feature: "auth", Status: tracker.StatusResolved},
		{ID: "02", Feature: "auth", Status: tracker.StatusOpen, BlockedBy: []string{"01"}},
		{ID: "03", Feature: "auth", Status: tracker.StatusOpen, BlockedBy: []string{"02"}},
		{ID: "04", Feature: "auth", Status: tracker.StatusOpen, BlockedBy: []string{"05"}},
		{ID: "05", Feature: "auth", Status: tracker.StatusOpen, BlockedBy: []string{"04"}}, // Cycle between 04 and 05
	}

	g := tracker.BuildGraph(issues)

	// 1. Check Cycle
	check := g.CheckDependencies()
	if !check.HasCycles {
		t.Errorf("expected cycles, got none")
	}
	if len(check.Cycles) == 0 {
		t.Errorf("expected cycle list not empty")
	}

	// 2. Check Frontier
	frontier := g.ComputeFrontier()
	if len(frontier.Ready) != 1 || frontier.Ready[0].ID != "02" {
		t.Errorf("expected Ready [02], got %v", frontier.Ready)
	}
	if len(frontier.Blocked) < 1 {
		t.Errorf("expected Blocked issues, got %v", frontier.Blocked)
	}
}
