package command

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"flowforge/internal/core"
	"flowforge/internal/state"
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

func TestJournalEventSealAndAnalysisStatusJSON(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Analysis journal")

	initCmd := newJournalEventInitCmd()
	initCmd.SetArgs([]string{"--proposal", proposalID, "--kind", "analysis.plan_published", "--id", "JEV-plan"})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("initializing event: %v", err)
	}
	store := testCardStore(t, tmpDir)
	path := store.ProposalJournalPath(proposalID)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading journal: %v", err)
	}
	updated := strings.Replace(string(content), `"cycleId": ""`, `"cycleId": "cycle-1"`, 1)
	updated = strings.Replace(updated, `"work": []`, `"work": [{"workId":"w1","question":"Inspect code","role":"flowforge-investigator","required":true}]`, 1)
	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		t.Fatalf("editing event: %v", err)
	}

	sealCmd := newJournalEventSealCmd()
	sealCmd.SetArgs([]string{"JEV-plan", "--proposal", proposalID})
	if err := sealCmd.Execute(); err != nil {
		t.Fatalf("sealing event: %v", err)
	}

	viper.Set("output", "json")
	t.Cleanup(func() { viper.Set("output", "text") })
	statusCmd := newAnalysisStatusCmd()
	var out bytes.Buffer
	statusCmd.SetOut(&out)
	statusCmd.SetArgs([]string{"--proposal", proposalID})
	if err := statusCmd.Execute(); err != nil {
		t.Fatalf("reading analysis status: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("decoding status JSON: %v\n%s", err, out.String())
	}
	for _, key := range []string{"schemaVersion", "proposalId", "sourceRevision", "state", "activePlan", "readyWork", "runningWork", "reentry", "issues", "nextAction"} {
		if _, ok := result[key]; !ok {
			t.Fatalf("status JSON missing %s: %s", key, out.String())
		}
	}
	if result["nextAction"] != "dispatch_ready_work" {
		t.Fatalf("unexpected next action: %v", result["nextAction"])
	}

	inspectCmd := newProposalInspectCmd()
	var inspectOut bytes.Buffer
	inspectCmd.SetOut(&inspectOut)
	inspectCmd.SetArgs([]string{proposalID})
	if err := inspectCmd.Execute(); err != nil {
		t.Fatalf("inspecting proposal: %v", err)
	}
	for _, want := range []string{"## Active Analysis", "- Revision: 1", "- Ready Work: w1", "## Next Action", "dispatch_ready_work"} {
		if !strings.Contains(inspectOut.String(), want) {
			t.Fatalf("proposal inspect missing %q:\n%s", want, inspectOut.String())
		}
	}

	viper.Set("output", "text")
	contextCmd := newContextProposalCmd()
	var contextOut bytes.Buffer
	contextCmd.SetOut(&contextOut)
	contextCmd.SetArgs([]string{"--proposal", proposalID})
	if err := contextCmd.Execute(); err != nil {
		t.Fatalf("reading proposal context: %v", err)
	}
	for _, want := range []string{"## Active Analysis", "- Ready Work: w1", "- Next Action: dispatch_ready_work"} {
		if !strings.Contains(contextOut.String(), want) {
			t.Fatalf("proposal context missing %q:\n%s", want, contextOut.String())
		}
	}
}

func TestStructuredCardMutationPreservesDirectBodyEditWhenIndexIsStale(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Stale index")
	cardStore := testCardStore(t, tmpDir)

	from := core.NewCard(core.CardTypeFeature, "From feature")
	from.ID = "FEAT-" + proposalID + "-from"
	from.Source = proposalID
	from.Body = "# From feature\n\n## Notes\n\nOriginal body.\n"
	fromPath, err := cardStore.CreateCard(from, proposalID)
	if err != nil {
		t.Fatalf("creating from card: %v", err)
	}
	to := core.NewCard(core.CardTypeFeature, "To feature")
	to.ID = "FEAT-" + proposalID + "-to"
	to.Source = proposalID
	to.Body = "# To feature\n"
	if _, err := cardStore.CreateCard(to, proposalID); err != nil {
		t.Fatalf("creating to card: %v", err)
	}

	_, _, runtimeStore, err := currentProposalStoreWithState()
	if err != nil {
		t.Fatalf("opening runtime store: %v", err)
	}
	defer closeStateStore(runtimeStore)
	syncService := state.NewCardSyncService(runtimeStore.DB())
	for _, id := range []string{from.ID, to.ID} {
		card, readErr := cardStore.ReadCard(id)
		if readErr != nil {
			t.Fatalf("reading %s: %v", id, readErr)
		}
		if syncErr := syncService.SyncCard(card); syncErr != nil {
			t.Fatalf("syncing stale snapshot: %v", syncErr)
		}
	}
	content, err := os.ReadFile(fromPath)
	if err != nil {
		t.Fatalf("reading from card: %v", err)
	}
	latest := strings.Replace(string(content), "Original body.", "Latest direct body edit.", 1)
	if err := os.WriteFile(fromPath, []byte(latest), 0644); err != nil {
		t.Fatalf("direct editing card: %v", err)
	}

	linkCmd := newCardLinkCmd()
	linkCmd.SetArgs([]string{from.ID, to.ID, "--relation", "references"})
	if err := linkCmd.Execute(); err != nil {
		t.Fatalf("linking cards: %v", err)
	}
	updated, err := os.ReadFile(fromPath)
	if err != nil {
		t.Fatalf("reading updated card: %v", err)
	}
	if !strings.Contains(string(updated), "Latest direct body edit.") {
		t.Fatalf("structured mutation overwrote latest body:\n%s", updated)
	}
}
