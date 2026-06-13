package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Version != "2.0.0" {
		t.Errorf("expected version 2.0.0, got %s", cfg.Version)
	}

	if len(cfg.Projects) != 1 {
		t.Fatalf("expected 1 default project, got %d", len(cfg.Projects))
	}

	if cfg.Projects[0].ID != "default" {
		t.Errorf("expected project id default, got %s", cfg.Projects[0].ID)
	}

	if cfg.Projects[0].WikiRoot != "ff-wiki" {
		t.Errorf("expected wiki root ff-wiki, got %s", cfg.Projects[0].WikiRoot)
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	configContent := `version: "2.0.0"
projects:
  - id: "default"
    wikiRoot: "docs"
    srcDirs:
      - "src"
      - "app"
`
	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if len(cfg.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(cfg.Projects))
	}

	if cfg.Projects[0].ID != "default" {
		t.Errorf("expected project id default, got %s", cfg.Projects[0].ID)
	}

	if cfg.Projects[0].WikiRoot != "docs" {
		t.Errorf("expected wiki root docs, got %s", cfg.Projects[0].WikiRoot)
	}

	if len(cfg.Projects[0].SrcDirs) != 2 {
		t.Fatalf("expected 2 source dirs, got %d", len(cfg.Projects[0].SrcDirs))
	}
}

func TestLoadConfigMissing(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("expected no error for missing config, got: %v", err)
	}

	if cfg.Version != "2.0.0" {
		t.Errorf("expected default version, got %s", cfg.Version)
	}

	if len(cfg.Projects) != 1 {
		t.Fatalf("expected 1 default project, got %d", len(cfg.Projects))
	}

	if cfg.Projects[0].ID != "default" {
		t.Errorf("expected default project id, got %s", cfg.Projects[0].ID)
	}
}

func TestFindProjectRoot(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte("version: 2.0.0"), 0644); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(tmpDir, "sub", "dir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	root, err := FindProjectRoot(subDir)
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	if root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, root)
	}
}

func TestFindProjectRootNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := FindProjectRoot(tmpDir)
	if err == nil {
		t.Error("expected error for missing project root")
	}
}

func TestWikiRoot(t *testing.T) {
	cfg := &Config{
		Projects: []ProjectConfig{{ID: "default", WikiRoot: "ff-wiki"}},
	}

	wikiRoot := cfg.WikiRoot("/project")
	expected := "/project/ff-wiki"
	if wikiRoot != expected {
		t.Errorf("expected %s, got %s", expected, wikiRoot)
	}
}

func TestWikiRootAbsolute(t *testing.T) {
	cfg := &Config{
		Projects: []ProjectConfig{{ID: "default", WikiRoot: "/absolute/path"}},
	}

	wikiRoot := cfg.WikiRoot("/project")
	expected := "/absolute/path"
	if wikiRoot != expected {
		t.Errorf("expected %s, got %s", expected, wikiRoot)
	}
}
