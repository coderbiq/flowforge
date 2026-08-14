package state

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/core"
)

func TestRebuildAllAndDerivedIndexIgnoreLegacyFilesWithoutRewritingThem(t *testing.T) {
	root := t.TempDir()
	legacyPath := filepath.Join(root, "TASK-legacy.md")
	strPath := filepath.Join(root, "STR-proposal-REQ.md")
	legacy := []byte("# legacy task\n\nlegacy links: [[FEAT-current]]\n")
	str := []byte("# STR metadata\n\nlegacy proposal metadata\n")
	for path, data := range map[string][]byte{legacyPath: legacy, strPath: str} {
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	dbPath := filepath.Join(root, "state.sqlite")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.EnsureSchema(); err != nil {
		t.Fatal(err)
	}
	service := NewCardSyncService(store.DB())
	current := core.NewCard(core.CardTypeFeature, "Current")
	current.ID = "FEAT-current"
	current.FilePath = filepath.Join(root, "FEAT-current.md")
	current.AddLink("TASK-legacy", "references")
	list := func(string) ([]*core.Card, error) { return []*core.Card{current}, nil }
	if cards, links, rebuildErr := service.RebuildAll(list, []string{"ignored"}); rebuildErr != nil || cards != 1 || links != 1 {
		t.Fatalf("RebuildAll = cards %d links %d err %v", cards, links, rebuildErr)
	}
	legacyCard := &core.Card{ID: "TASK-legacy", Type: core.CardType("task")}
	strCard := &core.Card{ID: "STR-proposal-REQ", Type: core.CardType("structure")}
	if cards, links, rebuildErr := store.RebuildDerivedIndex([]*core.Card{current, legacyCard, strCard}); rebuildErr != nil || cards != 1 || links != 1 {
		t.Fatalf("RebuildDerivedIndex = cards %d links %d err %v", cards, links, rebuildErr)
	}
	for _, id := range []string{"TASK-legacy", "STR-proposal-REQ"} {
		var count int
		if err := store.DB().QueryRow("SELECT COUNT(*) FROM card_index WHERE id = ?", id).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatalf("legacy/control-plane card %s was written to derived card index", id)
		}
	}
	for path, want := range map[string][]byte{legacyPath: legacy, strPath: str} {
		got, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatal(readErr)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("legacy file %s was rewritten", path)
		}
	}
}

func TestRebuildAllScansEachPhysicalDirectoryOnce(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "state.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.EnsureSchema(); err != nil {
		t.Fatal(err)
	}

	calls := 0
	list := func(string) ([]*core.Card, error) {
		calls++
		return nil, nil
	}
	service := NewCardSyncService(store.DB())
	if _, _, err := service.RebuildAll(list, []string{"workspace", "workspace", "./workspace"}); err != nil {
		t.Fatalf("RebuildAll failed: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected one physical scan, got %d", calls)
	}
}
