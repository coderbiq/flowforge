package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func TestIndexRebuildStatusAndBacklinks(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")

	store, err := state.Open(filepath.Join(tmpDir, ".flowforge", "cache", "flowforge.sqlite"))
	if err != nil {
		t.Fatalf("opening runtime store failed: %v", err)
	}
	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("ensuring runtime schema failed: %v", err)
	}
	if err := store.SetCurrentProjectID("default"); err != nil {
		t.Fatalf("setting current project failed: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}

	cardStore := testCardStore(t, tmpDir)
	legacyPath := filepath.Join(cardStore.LibraryDir(), "40-tasks", "TASK-legacy.md")
	historyPath := filepath.Join(tmpDir, "ff-wiki", "historical-wiki.md")
	legacy := []byte("# legacy task\n\nold link [[OLD-123]]\n")
	history := []byte("# historical wiki\n\nold content\n")
	for path, data := range map[string][]byte{legacyPath: legacy, historyPath: history} {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	alpha := core.NewCard(core.CardTypeFeature, "Alpha requirement")
	alpha.ID = "FEAT-ALPHA"
	alpha.Status = core.CardStatusActive
	alpha.Source = "default"
	alpha.Domain = "platform"
	alpha.AddLink("FEAT-BETA", "references")
	if err := alpha.Save(filepath.Join(cardStore.ActiveDir(), "FEAT-ALPHA-alpha-requirement.md")); err != nil {
		t.Fatalf("saving alpha card failed: %v", err)
	}

	beta := core.NewCard(core.CardTypeFeature, "Beta design")
	beta.ID = "FEAT-BETA"
	beta.Status = core.CardStatusReady
	beta.Source = "default"
	if err := beta.Save(filepath.Join(cardStore.IntakeDir(), "FEAT-BETA-beta-design.md")); err != nil {
		t.Fatalf("saving beta card failed: %v", err)
	}

	completedLogDir := filepath.Join(cardStore.CompletedDir(), "CR260612", "90-cards")
	if err := os.MkdirAll(completedLogDir, 0755); err != nil {
		t.Fatalf("creating completed log dir failed: %v", err)
	}
	completedLog := core.NewCard(core.CardTypeFeature, "Completed log")
	completedLog.ID = "FEAT-COMPLETED"
	completedLog.Status = core.CardStatusDone
	completedLog.Source = "CR260612"
	completedLog.AddLink("FEAT-BETA", "records")
	if err := completedLog.Save(filepath.Join(completedLogDir, "FEAT-COMPLETED-completed-log.md")); err != nil {
		t.Fatalf("saving completed log failed: %v", err)
	}

	rebuildCmd := newIndexRebuildCmd()
	var rebuildOut bytes.Buffer
	rebuildCmd.SetOut(&rebuildOut)
	if err := rebuildCmd.Execute(); err != nil {
		t.Fatalf("index rebuild failed: %v", err)
	}
	rebuildText := rebuildOut.String()
	for _, want := range []string{"✓ Rebuilt index for project default", "cards: 3", "links: 2"} {
		if !strings.Contains(rebuildText, want) {
			t.Fatalf("rebuild output missing %q:\n%s", want, rebuildText)
		}
	}

	statusCmd := newIndexStatusCmd()
	var statusOut bytes.Buffer
	statusCmd.SetOut(&statusOut)
	if err := statusCmd.Execute(); err != nil {
		t.Fatalf("index status failed: %v", err)
	}
	statusText := statusOut.String()
	for _, want := range []string{"Project: default", "card_index: 3", "card_link: 2"} {
		if !strings.Contains(statusText, want) {
			t.Fatalf("status output missing %q:\n%s", want, statusText)
		}
	}

	backlinksCmd := newIndexBacklinksCmd()
	var backlinksOut bytes.Buffer
	backlinksCmd.SetOut(&backlinksOut)
	backlinksCmd.SetArgs([]string{"FEAT-BETA"})
	if err := backlinksCmd.Execute(); err != nil {
		t.Fatalf("index backlinks failed: %v", err)
	}
	backlinksText := backlinksOut.String()
	for _, want := range []string{"Backlinks for FEAT-BETA:", "FEAT-COMPLETED records", "FEAT-ALPHA references"} {
		if !strings.Contains(backlinksText, want) {
			t.Fatalf("backlinks output missing %q:\n%s", want, backlinksText)
		}
	}
	for path, want := range map[string][]byte{legacyPath: legacy, historyPath: history} {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("index rebuild rewrote historical file %s", path)
		}
	}
}

func TestIndexBacklinksSuggestsRebuildWhenEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")

	cmd := newIndexBacklinksCmd()
	cmd.SetArgs([]string{"FEAT-MISSING"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "run `flowforge index rebuild` first") {
		t.Fatalf("expected helpful rebuild message, got %v", err)
	}
}
