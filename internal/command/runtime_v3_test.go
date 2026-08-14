package command

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/config"
	"flowforge/internal/core"
	"flowforge/internal/state"
)

func TestLegacyCommandsAreNotRegistered(t *testing.T) {
	root := NewRootCmd()
	for _, name := range []string{"task", "structure", "log"} {
		if _, _, err := root.Find([]string{name}); err == nil {
			t.Fatalf("legacy command %q is still registered", name)
		}
	}
}

func restoreWorkingDir(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore working dir failed: %v", err)
		}
	})
}

func testCardStore(t *testing.T, projectRoot string) *core.CardStore {
	t.Helper()
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatalf("loading config failed: %v", err)
	}
	stateStore, err := state.Open(filepath.Join(projectRoot, ".flowforge", "cache", "flowforge.sqlite"))
	if err != nil {
		t.Fatalf("opening state store failed: %v", err)
	}
	t.Cleanup(func() {
		if err := stateStore.Close(); err != nil {
			t.Fatalf("closing state store failed: %v", err)
		}
	})
	return core.NewCardStoreWithSync(cfg.WikiRoot(projectRoot), state.NewCardSyncService(stateStore.DB()))
}

func createProjectForTest(t *testing.T, projectID string) {
	t.Helper()
	cmd := newProjectCreateCmd()
	cmd.SetArgs([]string{projectID})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("project create failed: %v", err)
	}
}

func TestLegacyAndSTRReadUseTheExistingNotFoundError(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	store := testCardStore(t, tmpDir)

	for _, id := range []string{"TASK-legacy", "STR-CR26081201-REQ"} {
		cmd := newCardReadCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{id})
		err := cmd.Execute()
		if err == nil {
			t.Fatalf("card read unexpectedly succeeded for %s", id)
		}
		if got := err.Error(); got != "card not found: "+id {
			t.Fatalf("card read error for %s changed: %q", id, got)
		}
	}
	if _, err := store.ReadCard("STR-CR26081201-REQ"); err == nil {
		t.Fatal("direct STR read unexpectedly succeeded")
	}
}

func TestCardListJSONContainsOnlyCurrentCardsAndStableFields(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	store := testCardStore(t, tmpDir)
	legacyPath := filepath.Join(store.LibraryDir(), "40-tasks", "TASK-legacy.md")
	strPath := filepath.Join(store.IntakeDir(), "STR-CR26081201-REQ.md")
	for _, path := range []string{legacyPath, strPath} {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(legacyPath, []byte("# legacy\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(strPath, []byte("# metadata\n"), 0644); err != nil {
		t.Fatal(err)
	}
	current := core.NewCard(core.CardTypeFeature, "Current card")
	current.ID = "FEAT-current-list"
	if _, err := store.CreateCard(current, ""); err != nil {
		t.Fatalf("creating current card failed: %v", err)
	}

	cmd := newCardListCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card list failed: %v", err)
	}
	var cards []map[string]any
	if err := json.Unmarshal(out.Bytes(), &cards); err != nil {
		t.Fatalf("invalid card list JSON: %v\n%s", err, out.String())
	}
	if len(cards) != 1 || cards[0]["id"] != current.ID {
		t.Fatalf("card list returned legacy or metadata cards: %#v", cards)
	}
	encoded := out.String()
	for _, forbidden := range []string{"TASK-legacy", "STR-CR26081201-REQ", "ignored", "deprecated", "count"} {
		if strings.Contains(encoded, forbidden) {
			t.Fatalf("card list JSON exposed forbidden %q: %s", forbidden, encoded)
		}
	}
}
