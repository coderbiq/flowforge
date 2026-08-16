package uninstall

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestCleanProjectDeleteFailureRollsBackAndRetryRecovers(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, ".opencode", "agents", "first.md")
	second := filepath.Join(root, ".opencode", "agents", "second.md")
	writeFile(t, first, "first")
	writeFile(t, second, "second")
	manifest := &core.ProjectManifest{Version: core.ManifestVersionV2, Mode: core.ManifestModeSubagent, HostIntent: core.HostIntent{OpenCode: core.HostEnabled, Codex: core.HostDisabled}, Files: []core.FileEntry{
		{Source: "generated/opencode/first.md", Target: ".opencode/agents/first.md", SHA256: core.SHA256Hex([]byte("first")), Type: "opencode_agent", Host: core.HostOpenCode},
		{Source: "generated/opencode/second.md", Target: ".opencode/agents/second.md", SHA256: core.SHA256Hex([]byte("second")), Type: "opencode_agent", Host: core.HostOpenCode},
	}}
	if err := core.SaveManifestAtomic(root, manifest); err != nil {
		t.Fatal(err)
	}
	manifestBefore := readFile(t, filepath.Join(root, ".flowforge", core.ManifestFileName))
	oldRemove := projectRemove
	count := 0
	projectRemove = func(path string) error {
		count++
		if count == 2 {
			return fmt.Errorf("injected delete failure")
		}
		return oldRemove(path)
	}
	t.Cleanup(func() { projectRemove = oldRemove })
	if _, err := CleanProject(root); err == nil {
		t.Fatal("expected injected delete failure")
	}
	if string(readFile(t, first)) != "first" || string(readFile(t, second)) != "second" {
		t.Fatal("delete failure did not restore all files")
	}
	if got := readFile(t, filepath.Join(root, ".flowforge", core.ManifestFileName)); string(got) != string(manifestBefore) {
		t.Fatal("delete failure changed manifest")
	}
	projectRemove = oldRemove
	if _, err := CleanProject(root); err != nil {
		t.Fatalf("retry did not recover: %v", err)
	}
}

func TestCleanProjectIsManifestScopedAndIdempotent(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".agents", "skills", "managed.md"), "managed skill")
	writeFile(t, filepath.Join(root, ".agents", "skills", "unknown.md"), "keep unknown")
	writeFile(t, filepath.Join(root, ".opencode", "agents", "managed.md"), "managed agent")
	agents := "user heading\n\n<!-- FLOWFORGE:START -->\nbase rules\n<!-- FLOWFORGE:END -->\n\n<!-- FLOWFORGE:ORCHESTRATION:START -->\norchestration\n<!-- FLOWFORGE:ORCHESTRATION:END -->\n\nother tool block\n"
	writeFile(t, filepath.Join(root, "AGENTS.md"), agents)

	manifest := &core.ProjectManifest{
		Version:    core.ManifestVersionV2,
		CLIVersion: "v3.0.0",
		Mode:       core.ManifestModeSubagent,
		HostIntent: core.HostIntent{OpenCode: core.HostEnabled, Codex: core.HostDisabled},
		Files: []core.FileEntry{
			{Source: "assets/skills/managed.md", Target: ".agents/skills/managed.md", SHA256: core.SHA256Hex([]byte("managed skill")), Type: "skill"},
			{Source: "assets/AGENTS.md", Target: "AGENTS.md", SHA256: "base", Type: "agents_block", Markers: &core.BlockMarkers{Start: "<!-- FLOWFORGE:START -->", End: "<!-- FLOWFORGE:END -->"}},
			{Source: "generated/opencode/managed.md", Target: ".opencode/agents/managed.md", SHA256: core.SHA256Hex([]byte("managed agent")), Type: "opencode_agent", Host: core.HostOpenCode},
			{Source: "generated/AGENTS.orchestration.md", Target: "AGENTS.md", SHA256: core.SHA256Hex([]byte("orchestration\n")), Type: "orchestration_block", Markers: &core.BlockMarkers{Start: "<!-- FLOWFORGE:ORCHESTRATION:START -->", End: "<!-- FLOWFORGE:ORCHESTRATION:END -->"}},
		},
	}
	if err := core.SaveManifestAtomic(root, manifest); err != nil {
		t.Fatal(err)
	}

	result, err := CleanProject(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("unexpected cleanup warnings: %v", result.Errors)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "skills", "managed.md")); !os.IsNotExist(err) {
		t.Fatalf("managed static file was not removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", "managed.md")); !os.IsNotExist(err) {
		t.Fatalf("managed dynamic file was not removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "skills", "unknown.md")); err != nil {
		t.Fatalf("unknown file was removed: %v", err)
	}
	agentsAfter := string(readFile(t, filepath.Join(root, "AGENTS.md")))
	if !strings.Contains(agentsAfter, "base rules") || !strings.Contains(agentsAfter, "user heading") || !strings.Contains(agentsAfter, "other tool block") || strings.Contains(agentsAfter, "orchestration") {
		t.Fatalf("AGENTS content was not scoped correctly:\n%s", agentsAfter)
	}
	backupRoot := filepath.Join(root, ".flowforge", "backups", "uninstall")
	if _, err := os.Stat(backupRoot); err != nil {
		t.Fatalf("dynamic backup was not created: %v", err)
	}

	if _, err := CleanProject(root); err != nil {
		t.Fatal(err)
	}
}

func TestCleanProjectPreservesModifiedStaticAndBacksUpModifiedDynamic(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".agents", "skills", "modified.md"), "user version")
	writeFile(t, filepath.Join(root, ".codex", "agents", "modified.toml"), "user agent")
	manifest := &core.ProjectManifest{
		Version:    core.ManifestVersionV2,
		Mode:       core.ManifestModeSubagent,
		HostIntent: core.HostIntent{OpenCode: core.HostDisabled, Codex: core.HostEnabled},
		Files: []core.FileEntry{
			{Source: "assets/skills/modified.md", Target: ".agents/skills/modified.md", SHA256: "old", Type: "skill"},
			{Source: "generated/codex/modified.toml", Target: ".codex/agents/modified.toml", SHA256: "old", Type: "codex_agent", Host: core.HostCodex},
		},
	}
	if err := core.SaveManifestAtomic(root, manifest); err != nil {
		t.Fatal(err)
	}
	result, err := CleanProject(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Errors) != 1 || !strings.Contains(result.Errors[0].Error(), "modified static") {
		t.Fatalf("expected modified static warning, got %v", result.Errors)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "skills", "modified.md")); err != nil {
		t.Fatalf("modified static file was removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".codex", "agents", "modified.toml")); !os.IsNotExist(err) {
		t.Fatalf("modified dynamic file was not removed: %v", err)
	}
}

func TestCleanConfigDoesNotTouchProject(t *testing.T) {
	home := t.TempDir()
	project := filepath.Join(home, "project")
	writeFile(t, filepath.Join(home, ".flowforge", "config.yaml"), "config")
	writeFile(t, filepath.Join(project, "AGENTS.md"), "project")
	if _, err := CleanConfig(home); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(project, "AGENTS.md")); err != nil {
		t.Fatalf("global cleanup touched project: %v", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
