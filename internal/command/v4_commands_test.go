package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func setupV4TestEnv(t *testing.T) (string, *core.CardStore) {
	t.Helper()
	tmpDir := t.TempDir()
	wikiRoot := filepath.Join(tmpDir, "ff-wiki-test")
	_ = os.MkdirAll(filepath.Join(wikiRoot, "01-workspace"), 0755)
	store := core.NewCardStore(wikiRoot)
	return tmpDir, store
}

func TestV4MemoryInitAndShow(t *testing.T) {
	tmpDir, store := setupV4TestEnv(t)
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	// Write config.yaml
	cfgDir := filepath.Join(tmpDir, ".flowforge")
	_ = os.MkdirAll(cfgDir, 0755)
	_ = os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("version: 2.0.0\nprojects:\n  - id: test\n    wikiRoot: ff-wiki-test\n"), 0644)

	// Test Memory Init Proposal
	rootCmd := NewRootCmd()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"memory", "init", "--proposal", "CR26082301_test", "--hierarchical"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("memory init failed: %v, output: %s", err, buf.String())
	}

	readmePath := filepath.Join(store.WorkspaceDir(), "CR26082301_test", "README.md")
	if _, err := os.Stat(readmePath); err != nil {
		t.Fatalf("expected proposal README at %s", readmePath)
	}

	// Verify modules directory created due to --hierarchical
	modPath := filepath.Join(store.WorkspaceDir(), "CR26082301_test", "modules")
	if _, err := os.Stat(modPath); err != nil {
		t.Fatalf("expected modules dir at %s", modPath)
	}

	// Test Memory Show
	buf.Reset()
	rootCmd.SetArgs([]string{"memory", "show", "--proposal", "CR26082301_test"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("memory show failed: %v", err)
	}
	if !strings.Contains(buf.String(), "Tier 2: Proposal Scratchpad") {
		t.Fatalf("unexpected memory show output:\n%s", buf.String())
	}
}

func TestV4ContextSlice(t *testing.T) {
	tmpDir, store := setupV4TestEnv(t)
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	cfgDir := filepath.Join(tmpDir, ".flowforge")
	_ = os.MkdirAll(cfgDir, 0755)
	_ = os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("version: 2.0.0\nprojects:\n  - id: test\n    wikiRoot: ff-wiki-test\n"), 0644)

	propDir := filepath.Join(store.WorkspaceDir(), "CR26082301_feat")
	_ = os.MkdirAll(propDir, 0755)
	readmeContent := `# CR26082301: Feature X

## Slices

### Slice 1: Database Migration
- **Seams**: internal/db/migrate.go
- **Test**: go test ./internal/db/...
Implementation details...
`
	_ = os.WriteFile(filepath.Join(propDir, "README.md"), []byte(readmeContent), 0644)

	rootCmd := NewRootCmd()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"context", "slice", "--proposal", "CR26082301_feat", "--slice", "1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("context slice failed: %v, out: %s", err, buf.String())
	}

	out := buf.String()
	if !strings.Contains(out, "Slice 1: Database Migration") || !strings.Contains(out, "go test ./internal/db/...") {
		t.Fatalf("context slice output missing expected content:\n%s", out)
	}
}
