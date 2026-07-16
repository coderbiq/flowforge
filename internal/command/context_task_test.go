package command

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestContextTaskIncludesLinkedCardsAndBacklinkEvidence(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Task context proposal")

	store := testCardStore(t, tmpDir)

	req := core.NewCard(core.CardTypeRequirement, "Linked requirement")
	req.ID = "REQ-CTX-1"
	req.Body = "## Summary\n\nRequirement summary."
	if _, err := store.CreateCard(req, proposalID); err != nil {
		t.Fatalf("creating requirement failed: %v", err)
	}

	conv := core.NewCard(core.CardTypeConvention, "Runtime convention")
	conv.ID = "CONV-CTX-1"
	conv.Body = "## Rules\n\nUse explicit errors."
	if _, err := store.CreateCard(conv, ""); err != nil {
		t.Fatalf("creating convention failed: %v", err)
	}

	task := core.NewCard(core.CardTypeTask, "Implement linked behavior")
	task.ID = "TASK-CTX-i-1"
	task.Status = core.CardStatusReady
	task.Body = "## Goal\n\nImplement the behavior."
	task.AddLink(req.ID, "satisfies")
	task.AddLink(conv.ID, "constrains")
	task.AddLink("REQ-MISSING", "references")
	if _, err := store.CreateCard(task, proposalID); err != nil {
		t.Fatalf("creating task failed: %v", err)
	}

	log := core.NewCard(core.CardTypeLog, "Implementation note")
	log.ID = "LOG-CTX-1"
	log.Body = "## Summary\n\nStarted implementation."
	log.AddLink(task.ID, "records")
	if _, err := store.CreateCard(log, proposalID); err != nil {
		t.Fatalf("creating log failed: %v", err)
	}

	cmd := newContextTaskCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--task", task.ID})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("context task failed: %v", err)
	}

	text := out.String()
	for _, want := range []string{
		"## Task Context: TASK-CTX-i-1",
		"Linked requirement",
		"Runtime convention",
		"Implementation note",
		"## Warnings",
		"linked card REQ-MISSING could not be read",
		"flowforge card read REQ-CTX-1 --summary",
		"flowforge card read LOG-CTX-1 --summary",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("context task output missing %q:\n%s", want, text)
		}
	}
}
