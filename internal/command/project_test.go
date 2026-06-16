package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/config"
	"flowforge/internal/core"
	"flowforge/internal/state"
)

func TestProjectCreateBootstrapsDerivedWikiRoot(t *testing.T) {
	tmpDir := t.TempDir()
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	cmd := newProjectCreateCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"frontend"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("project create failed: %v", err)
	}
	if !strings.Contains(out.String(), "✓ Project created: frontend") {
		t.Fatalf("project create output missing success message:\n%s", out.String())
	}

	cfg, err := config.Load(tmpDir)
	if err != nil {
		t.Fatalf("loading config failed: %v", err)
	}
	if len(cfg.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(cfg.Projects))
	}

	project, ok := cfg.ProjectByID("frontend")
	if !ok {
		t.Fatalf("expected frontend project to be registered")
	}
	if project.WikiRoot != "ff-wiki-frontend" {
		t.Fatalf("expected derived wiki root ff-wiki-frontend, got %s", project.WikiRoot)
	}

	expectedDirs := []string{
		filepath.Join(tmpDir, "ff-wiki-frontend", "01-workspace", "01-active"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "01-workspace", "02-intake"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "01-workspace", "03-completed"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "10-requirements"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "20-decisions"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "30-designs"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "40-tasks"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "50-logs"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "60-conventions"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "70-findings"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "02-library", "80-modules"),
		filepath.Join(tmpDir, "ff-wiki-frontend", "03-proposal"),
	}
	for _, dir := range expectedDirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("expected directory %s: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected directory, got file: %s", dir)
		}
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "ff-wiki-frontend", "00-STR-HOME.md")); err != nil {
		t.Fatalf("expected home index file: %v", err)
	}

	store, err := state.Open(filepath.Join(tmpDir, ".flowforge", "cache", "flowforge.sqlite"))
	if err != nil {
		t.Fatalf("opening runtime store failed: %v", err)
	}
	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("ensuring runtime schema failed: %v", err)
	}
	currentID, ok, err := store.CurrentProjectID()
	if err != nil {
		t.Fatalf("reading current project failed: %v", err)
	}
	if !ok || currentID != "frontend" {
		t.Fatalf("expected current project frontend, got %q, ok=%v", currentID, ok)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}
}

func TestProjectCreateSupportsFlagsAndDuplicateCheck(t *testing.T) {
	tmpDir := t.TempDir()
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	cmd := newProjectCreateCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"backend", "--wiki-root", "custom-wiki", "--src-dir", "api", "--src-dir", "worker", "--default"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("project create with flags failed: %v", err)
	}

	cfg, err := config.Load(tmpDir)
	if err != nil {
		t.Fatalf("loading config failed: %v", err)
	}
	project, ok := cfg.ProjectByID("backend")
	if !ok {
		t.Fatalf("expected backend project to be registered")
	}
	if project.WikiRoot != "custom-wiki" {
		t.Fatalf("expected wiki root custom-wiki, got %s", project.WikiRoot)
	}
	if len(project.SrcDirs) != 2 || project.SrcDirs[0] != "api" || project.SrcDirs[1] != "worker" {
		t.Fatalf("unexpected srcDirs: %#v", project.SrcDirs)
	}

	store, err := state.Open(filepath.Join(tmpDir, ".flowforge", "cache", "flowforge.sqlite"))
	if err != nil {
		t.Fatalf("opening runtime store failed: %v", err)
	}
	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("ensuring runtime schema failed: %v", err)
	}
	currentID, ok, err := store.CurrentProjectID()
	if err != nil {
		t.Fatalf("reading current project failed: %v", err)
	}
	if !ok || currentID != "backend" {
		t.Fatalf("expected current project backend, got %q, ok=%v", currentID, ok)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "custom-wiki", "00-STR-HOME.md")); err != nil {
		t.Fatalf("expected custom home index file: %v", err)
	}

	duplicateCmd := newProjectCreateCmd()
	duplicateCmd.SetArgs([]string{"backend"})
	if err := duplicateCmd.Execute(); err == nil {
		t.Fatalf("expected duplicate project create to fail")
	}
}

func TestProjectCreateUsesDefaultWikiRootForReservedID(t *testing.T) {
	tmpDir := t.TempDir()
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	cmd := newProjectCreateCmd()
	cmd.SetArgs([]string{"default"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("project create failed: %v", err)
	}

	cfg, err := config.Load(tmpDir)
	if err != nil {
		t.Fatalf("loading config failed: %v", err)
	}
	project, ok := cfg.ProjectByID("default")
	if !ok {
		t.Fatalf("expected default project to be registered")
	}
	if project.WikiRoot != "ff-wiki" {
		t.Fatalf("expected wiki root ff-wiki, got %s", project.WikiRoot)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "ff-wiki", "00-STR-HOME.md")); err != nil {
		t.Fatalf("expected default home index file: %v", err)
	}
}

func TestProjectListCurrentAndUseCommands(t *testing.T) {
	tmpDir := createMultiProjectFixture(t)
	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	currentCmd := newProjectCurrentCmd()
	if err := currentCmd.Execute(); err == nil {
		t.Fatalf("expected project current to require project use with multiple projects")
	}

	useCmd := newProjectUseCmd()
	var useOut bytes.Buffer
	useCmd.SetOut(&useOut)
	useCmd.SetArgs([]string{"backend"})
	if err := useCmd.Execute(); err != nil {
		t.Fatalf("project use failed: %v", err)
	}
	if !strings.Contains(useOut.String(), "✓ Current project: backend") {
		t.Fatalf("project use output missing current project:\n%s", useOut.String())
	}

	listCmd := newProjectListCmd()
	var listOut bytes.Buffer
	listCmd.SetOut(&listOut)
	if err := listCmd.Execute(); err != nil {
		t.Fatalf("project list failed: %v", err)
	}
	listText := listOut.String()
	for _, want := range []string{
		"Projects:",
		"  frontend",
		"* backend",
		"wikiRoot: ff-wiki-be",
	} {
		if !strings.Contains(listText, want) {
			t.Fatalf("project list output missing %q:\n%s", want, listText)
		}
	}

	var currentOut bytes.Buffer
	currentOut.Reset()
	currentCmd = newProjectCurrentCmd()
	currentCmd.SetOut(&currentOut)
	if err := currentCmd.Execute(); err != nil {
		t.Fatalf("project current after use failed: %v", err)
	}
	if !strings.Contains(currentOut.String(), "Project: backend") {
		t.Fatalf("project current did not use runtime state:\n%s", currentOut.String())
	}
}

func TestProjectCurrentUsesSingleProjectFallback(t *testing.T) {
	tmpDir := createSingleProjectFixture(t)
	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	currentCmd := newProjectCurrentCmd()
	var currentOut bytes.Buffer
	currentCmd.SetOut(&currentOut)
	if err := currentCmd.Execute(); err != nil {
		t.Fatalf("project current failed: %v", err)
	}
	for _, want := range []string{
		"Project: frontend",
		"Source: single-project",
		"WikiRoot: ff-wiki-fe",
	} {
		if !strings.Contains(currentOut.String(), want) {
			t.Fatalf("project current output missing %q:\n%s", want, currentOut.String())
		}
	}
}

func TestProjectUseRejectsUnknownProject(t *testing.T) {
	tmpDir := createMultiProjectFixture(t)
	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	cmd := newProjectUseCmd()
	cmd.SetArgs([]string{"missing"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected project use to reject unknown project")
	}
}

func TestCurrentCardStoreUsesRuntimeProject(t *testing.T) {
	tmpDir := createMultiProjectFixture(t)
	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	useCmd := newProjectUseCmd()
	useCmd.SetArgs([]string{"backend"})
	if err := useCmd.Execute(); err != nil {
		t.Fatalf("project use failed: %v", err)
	}

	store, err := currentCardStore()
	if err != nil {
		t.Fatalf("currentCardStore failed: %v", err)
	}

	card := core.NewCard(core.CardTypeTask, "Backend task")
	card.ID = core.GenerateTaskID("260613", "i")
	if _, err := store.CreateCard(card, "CR26061301"); err != nil {
		t.Fatalf("creating backend card failed: %v", err)
	}

	backendCardsDir := filepath.Join(tmpDir, "ff-wiki-be", "01-workspace", "01-active", "CR26061301", "90-cards")
	entries, err := os.ReadDir(backendCardsDir)
	if err != nil {
		t.Fatalf("reading backend cards dir failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected backend card dir to contain 1 entry, got %d", len(entries))
	}

	frontendCardsDir := filepath.Join(tmpDir, "ff-wiki-fe", "01-workspace", "01-active", "CR26061301", "90-cards")
	if _, err := os.Stat(frontendCardsDir); !os.IsNotExist(err) {
		t.Fatalf("expected frontend cards dir to not exist, stat err=%v", err)
	}
}

func createMultiProjectFixture(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".flowforge")
	if err := os.MkdirAll(filepath.Join(configDir, "cache"), 0755); err != nil {
		t.Fatalf("creating config dir failed: %v", err)
	}

	configContent := `version: "2.0.0"
projects:
  - id: "frontend"
    wikiRoot: "ff-wiki-fe"
    srcDirs:
      - "web"
  - id: "backend"
    wikiRoot: "ff-wiki-be"
    srcDirs:
      - "api"
`
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("writing config failed: %v", err)
	}

	return tmpDir
}

func createSingleProjectFixture(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".flowforge")
	if err := os.MkdirAll(filepath.Join(configDir, "cache"), 0755); err != nil {
		t.Fatalf("creating config dir failed: %v", err)
	}

	configContent := `version: "2.0.0"
projects:
  - id: "frontend"
    wikiRoot: "ff-wiki-fe"
    srcDirs:
      - "web"
`
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatalf("writing config failed: %v", err)
	}

	return tmpDir
}
