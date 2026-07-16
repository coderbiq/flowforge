package command

import (
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/core"
)

func TestMigrateV3WikiFlattenExecutesAndRemovesOldStructure(t *testing.T) {
	dir := t.TempDir()
	wikiRoot := filepath.Join(dir, "wiki")
	setupOldWikiStructure(t, wikiRoot)

	store := core.NewCardStore(wikiRoot)
	m := allMigrations[0]

	if err := m.run(store, wikiRoot); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	entries, _ := os.ReadDir(filepath.Join(wikiRoot, "01-workspace"))
	propDirs := 0
	for _, e := range entries {
		if e.IsDir() {
			propDirs++
		}
	}
	if propDirs != 2 {
		t.Fatalf("expected 2 proposal dirs after migration, got %d", propDirs)
	}

	if _, err := os.Stat(filepath.Join(wikiRoot, "01-workspace", "01-active")); !os.IsNotExist(err) {
		t.Error("old 01-active directory should be removed")
	}
	if _, err := os.Stat(filepath.Join(wikiRoot, "01-workspace", "03-completed")); !os.IsNotExist(err) {
		t.Error("old 03-completed directory should be removed")
	}

	propCard, err := store.ReadCard("PROP-CR26060101")
	if err != nil {
		t.Fatalf("reading completed proposal card: %v", err)
	}
	if propCard.Status != core.CardStatusCompleted {
		t.Errorf("expected completed proposal status 'completed', got '%s'", propCard.Status)
	}
}

func TestRunPendingMigrationsIdempotent(t *testing.T) {
	dir := t.TempDir()
	wikiRoot := filepath.Join(dir, "wiki")
	setupOldWikiStructure(t, wikiRoot)

	store := core.NewCardStore(wikiRoot)
	m := allMigrations[0]

	if err := m.run(store, wikiRoot); err != nil {
		t.Fatalf("first migration failed: %v", err)
	}
	if err := m.run(store, wikiRoot); err != nil {
		t.Fatalf("second migration failed: %v", err)
	}

	entries, _ := os.ReadDir(filepath.Join(wikiRoot, "01-workspace"))
	propDirs := 0
	for _, e := range entries {
		if e.IsDir() {
			propDirs++
		}
	}
	if propDirs != 2 {
		t.Fatalf("expected 2 proposal dirs after idempotent run, got %d", propDirs)
	}
}

func TestRunPendingMigrationsSkipsWhenAlreadyFlat(t *testing.T) {
	dir := t.TempDir()
	wikiRoot := filepath.Join(dir, "wiki")
	workspaceDir := filepath.Join(wikiRoot, "01-workspace")
	os.MkdirAll(filepath.Join(workspaceDir, "CR-test"), 0755)
	os.MkdirAll(filepath.Join(wikiRoot, "02-library"), 0755)

	store := core.NewCardStore(wikiRoot)
	m := allMigrations[0]

	if err := m.run(store, wikiRoot); err != nil {
		t.Fatalf("migration on flat structure failed: %v", err)
	}
}

func TestCompareVersion(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"3.0.1", "3.0.2", -1},
		{"3.0.2", "3.0.2", 0},
		{"3.0.3", "3.0.2", 1},
		{"3.0.0-alpha", "3.0.2", -1},
		{"3.0.5", "3.0.6", -1},
		{"3.0.6", "3.0.6", 0},
		{"3.0.7", "3.0.6", 1},
		{"2.0.0", "3.0.0", -1},
		{"3.1.0", "3.0.6", 1},
	}
	for _, tt := range tests {
		got := compareVersion(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("compareVersion(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestAllMigrationsCoveredByVersionCheck(t *testing.T) {
	if len(allMigrations) == 0 {
		t.Skip("no migrations registered")
	}
	for _, m := range allMigrations {
		if m.name == "" {
			t.Error("migration has empty name")
		}
		if m.minVersion == "" {
			t.Errorf("migration %q has empty minVersion", m.name)
		}
		if m.run == nil {
			t.Errorf("migration %q has nil run function", m.name)
		}
	}
}

func setupOldWikiStructure(t *testing.T, wikiRoot string) {
	t.Helper()

	workspaceDir := filepath.Join(wikiRoot, "01-workspace")
	os.MkdirAll(filepath.Join(workspaceDir, "01-active", "CR26061201", "90-cards"), 0755)
	os.MkdirAll(filepath.Join(workspaceDir, "02-intake"), 0755)
	os.MkdirAll(filepath.Join(workspaceDir, "03-completed", "CR26060101", "90-cards"), 0755)
	os.MkdirAll(filepath.Join(wikiRoot, "02-library"), 0755)

	propActive := core.NewCard(core.CardTypeProposal, "Active Proposal")
	propActive.ID = "PROP-CR26061201"
	propActive.Status = core.CardStatusActive
	propActive.Save(filepath.Join(workspaceDir, "01-active", "CR26061201", "PROP-CR26061201.md"))

	propCompleted := core.NewCard(core.CardTypeProposal, "Completed Proposal")
	propCompleted.ID = "PROP-CR26060101"
	propCompleted.Status = core.CardStatusActive
	propCompleted.Save(filepath.Join(workspaceDir, "03-completed", "CR26060101", "PROP-CR26060101.md"))
}
