package command

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/config"

	_ "modernc.org/sqlite"
)

func TestRunInitCreatesInstallOnly(t *testing.T) {
	tmpDir := t.TempDir()

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	cfg, err := config.Load(tmpDir)
	if err != nil {
		t.Fatalf("loading config failed: %v", err)
	}
	if len(cfg.Projects) != 0 {
		t.Fatalf("expected no registered projects, got %d", len(cfg.Projects))
	}

	configPath := filepath.Join(tmpDir, ".flowforge", "config.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("stat config failed: %v", err)
	}

	dbPath := filepath.Join(tmpDir, ".flowforge", "cache", "flowforge.sqlite")
	info, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("stat sqlite db failed: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("expected sqlite db file, got directory: %s", dbPath)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("opening sqlite db failed: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Errorf("closing sqlite db failed: %v", closeErr)
		}
	})

	var name string
	if err := db.QueryRow(
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'runtime_state'",
	).Scan(&name); err != nil {
		t.Fatalf("querying sqlite_master failed: %v", err)
	}
	if name != "runtime_state" {
		t.Fatalf("expected runtime_state table, got %q", name)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "ff-wiki")); !os.IsNotExist(err) {
		t.Fatalf("expected no default wiki root, stat err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "ff-wiki", "00-STR-HOME.md")); !os.IsNotExist(err) {
		t.Fatalf("expected no default home index, stat err=%v", err)
	}

	for _, skill := range []string{"flowforge-design", "flowforge-implement"} {
		skillPath := filepath.Join(tmpDir, ".agents", "skills", skill, "SKILL.md")
		if _, err := os.Stat(skillPath); err != nil {
			t.Fatalf("expected deployed skill %s: %v", skill, err)
		}
	}

	if _, err := os.Stat(filepath.Join(tmpDir, ".flowforge", "templates")); err != nil {
		t.Fatalf("expected deployed templates directory: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "AGENTS.md")); err != nil {
		t.Fatalf("expected deployed AGENTS.md: %v", err)
	}
}

func TestRunInitAppendsFlowForgeBlockToExistingAgentRules(t *testing.T) {
	tmpDir := t.TempDir()
	agentPath := filepath.Join(tmpDir, "AGENTS.md")
	original := []byte("# Existing Rules\n")
	if err := os.WriteFile(agentPath, original, 0644); err != nil {
		t.Fatalf("writing existing AGENTS.md failed: %v", err)
	}

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	content, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("reading AGENTS.md failed: %v", err)
	}

	if !strings.HasPrefix(string(content), string(original)) {
		t.Fatalf("expected existing rules to be preserved at start, got:\n%s", string(content))
	}

	if !strings.Contains(string(content), "<!-- FLOWFORGE:START -->") {
		t.Fatalf("expected FlowForge block to be appended, got:\n%s", string(content))
	}
}
