package command

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"flowforge/internal/core"
)

func TestJournalAppendAndRecentUseCurrentProposal(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Journal proposal")

	appendCmd := newJournalAppendCmd()
	var appendOut bytes.Buffer
	appendCmd.SetOut(&appendOut)
	appendCmd.SetArgs([]string{
		"--actor", "design-analyst",
		"--message", "Completed proposal design",
		"--status", "designed",
		"--next", "Wait for approval",
	})
	if err := appendCmd.Execute(); err != nil {
		t.Fatalf("journal append failed: %v", err)
	}
	if !strings.Contains(appendOut.String(), "Journal entry appended to "+proposalID) {
		t.Fatalf("unexpected append output:\n%s", appendOut.String())
	}

	recentCmd := newJournalRecentCmd()
	var recentOut bytes.Buffer
	recentCmd.SetOut(&recentOut)
	if err := recentCmd.Execute(); err != nil {
		t.Fatalf("journal recent failed: %v", err)
	}
	for _, want := range []string{
		"design-analyst",
		"- Summary: Completed proposal design",
		"- Status: designed",
		"- Next: Wait for approval",
	} {
		if !strings.Contains(recentOut.String(), want) {
			t.Fatalf("recent output missing %q:\n%s", want, recentOut.String())
		}
	}
}

func TestJournalAppendValidatesReferencesAndRequiredArguments(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Journal proposal")

	missingActorCmd := newJournalAppendCmd()
	missingActorCmd.SetArgs([]string{"--message", "Missing actor"})
	if err := missingActorCmd.Execute(); err == nil || !strings.Contains(err.Error(), "--actor is required") {
		t.Fatalf("expected actor validation error, got %v", err)
	}

	invalidReferenceCmd := newJournalAppendCmd()
	invalidReferenceCmd.SetArgs([]string{
		"--proposal", proposalID,
		"--actor", "executor",
		"--message", "Invalid reference",
		"--references", "FEAT-missing",
	})
	if err := invalidReferenceCmd.Execute(); err == nil || !strings.Contains(err.Error(), "reading journal reference FEAT-missing") {
		t.Fatalf("expected reference validation error, got %v", err)
	}
}

func TestJournalRecentJSONOutputAndLimit(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Journal proposal")
	store := testCardStore(t, tmpDir)
	for _, message := range []string{"first", "second"} {
		if err := store.AppendProposalJournal(proposalID, core.JournalEntry{
			Actor:   "executor",
			Message: message,
		}); err != nil {
			t.Fatalf("creating journal entry: %v", err)
		}
	}
	viper.Set("output", "json")
	t.Cleanup(func() {
		viper.Set("output", "text")
	})

	cmd := newJournalRecentCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--proposal", proposalID, "--limit", "1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("journal recent failed: %v", err)
	}
	var entries []journalEntryOutput
	if err := json.Unmarshal(out.Bytes(), &entries); err != nil {
		t.Fatalf("decoding JSON output: %v\n%s", err, out.String())
	}
	if len(entries) != 1 || entries[0].Message != "second" {
		t.Fatalf("unexpected JSON output: %#v", entries)
	}
}

func TestJournalRecentRejectsNegativeLimit(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Journal proposal")

	cmd := newJournalRecentCmd()
	cmd.SetArgs([]string{"--proposal", proposalID, "--limit", "-1"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "journal entry limit cannot be negative") {
		t.Fatalf("expected negative limit error, got %v", err)
	}
}
