package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspaceProposalsAndSlices(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)

	proposalID := "CR26082301_test-feat"
	propDir := filepath.Join(store.WorkspaceDir(), proposalID)
	if err := os.MkdirAll(filepath.Join(propDir, "modules"), 0755); err != nil {
		t.Fatalf("mkdir modules failed: %v", err)
	}

	readmeContent := `# CR26082301: Test Feature

## Objective & Consensus
Test workspace proposal scanning.

## Actionable Slices
- [ ] **Slice 1: Scaffolding**
  - Seams: internal/core/workspace.go
  - Automated Tests: go test ./internal/core/...
- [x] **Slice 2: Delivery**
  - Seams: internal/command/proposal.go
  - Automated Tests: go test ./internal/command/...
`
	readmePath := filepath.Join(propDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		t.Fatalf("write README failed: %v", err)
	}

	proposals, err := store.ListWorkspaceProposals()
	if err != nil {
		t.Fatalf("list proposals failed: %v", err)
	}
	if len(proposals) != 1 {
		t.Fatalf("expected 1 proposal, got %d", len(proposals))
	}
	if proposals[0].ID != proposalID {
		t.Fatalf("expected proposal ID %q, got %q", proposalID, proposals[0].ID)
	}
	if proposals[0].Title != "CR26082301: Test Feature" {
		t.Fatalf("expected title %q, got %q", "CR26082301: Test Feature", proposals[0].Title)
	}
	if proposals[0].Mode != "hierarchical" {
		t.Fatalf("expected mode hierarchical, got %q", proposals[0].Mode)
	}

	slices, err := ParseProposalSlices(readmePath)
	if err != nil {
		t.Fatalf("parse slices failed: %v", err)
	}
	if len(slices) != 2 {
		t.Fatalf("expected 2 slices, got %d", len(slices))
	}
	if !slices[1].Completed {
		t.Fatalf("expected slice 2 to be completed")
	}

	found, err := store.FindWorkspaceProposal(proposalID)
	if err != nil {
		t.Fatalf("find workspace proposal failed: %v", err)
	}
	if found == nil {
		t.Fatalf("expected found proposal not nil")
	}
	if found.Title != "CR26082301: Test Feature" {
		t.Fatalf("expected title %q, got %q", "CR26082301: Test Feature", found.Title)
	}
}
