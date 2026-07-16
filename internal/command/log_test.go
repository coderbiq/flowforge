package command

import (
	"os"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestLogCreateDefaultsToCurrentProposalAndRecordsContext(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Log proposal")

	taskCmd := newTaskCreateCmd()
	taskCmd.SetArgs([]string{"--title", "Tracked task", "--type", "a"})
	if err := taskCmd.Execute(); err != nil {
		t.Fatalf("task create failed: %v", err)
	}

	store := testCardStore(t, tmpDir)
	taskID := findCardIDByTitle(t, store, proposalID, core.CardTypeTask, "Tracked task")

	logCmd := newLogCreateCmd()
	logCmd.SetArgs([]string{"--kind", "progress", "--title", "Captured progress", "--summary", "Task analysis started", "--for", taskID})
	if err := logCmd.Execute(); err != nil {
		t.Fatalf("log create failed: %v", err)
	}

	logCardID := findCardIDByTitle(t, store, proposalID, core.CardTypeLog, "Captured progress")
	logCard, err := store.ReadCard(logCardID)
	if err != nil {
		t.Fatalf("reading log card failed: %v", err)
	}
	if logCard.Source != proposalID {
		t.Fatalf("expected log source %s, got %q", proposalID, logCard.Source)
	}
	if !hasLink(logCard, taskID, "records") {
		t.Fatalf("expected log to record task %s", taskID)
	}
	if !strings.Contains(logCard.Body, "Task analysis started") {
		t.Fatalf("expected log body to contain summary, got:\n%s", logCard.Body)
	}
}

func findCardIDByTitle(t *testing.T, store *core.CardStore, proposalID string, cardType core.CardType, titlePart string) string {
	t.Helper()

	cards, err := store.ListCards(store.ProposalCardsDir(proposalID))
	if err != nil {
		t.Fatalf("listing proposal cards failed: %v", err)
	}
	for _, card := range cards {
		if card.Type == cardType && strings.Contains(card.Title, titlePart) {
			return card.ID
		}
	}

	t.Fatalf("expected %s card containing title %q", cardType, titlePart)
	return ""
}

func hasLink(card *core.Card, target string, relation string) bool {
	for _, link := range card.Links {
		if link.Target == target && link.Relation == relation {
			return true
		}
	}

	return false
}
