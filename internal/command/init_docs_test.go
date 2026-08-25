package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()
	oldWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWorkingDir) })
}

func TestInitPreservesConfiguredDocsDirUnderForce(t *testing.T) {
	projectRoot := t.TempDir()
	configDir := filepath.Join(projectRoot, ".flowforge")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	configured := "version: 5.0.0\nversion_check: true\ndocs_dir: wiki\n"
	if err := os.WriteFile(configPath, []byte(configured), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := newInitCmd()
	t.Cleanup(func() { initForce = false })
	if err := cmd.Flags().Set("force", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.RunE(cmd, []string{projectRoot}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != configured {
		t.Fatalf("init --force replaced project config:\n%s", data)
	}
	for _, path := range []string{
		filepath.Join(projectRoot, "wiki", "proposals"),
		filepath.Join(projectRoot, "wiki", "adr"),
		filepath.Join(projectRoot, "wiki", "CONTEXT.md"),
		filepath.Join(projectRoot, "wiki", "agents", "issue-tracker.md"),
		filepath.Join(projectRoot, ".agents", "skills"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected initialized path %s: %v", path, err)
		}
	}
	for _, skill := range []string{"flowforge-review", "flowforge-setup"} {
		data, err := os.ReadFile(filepath.Join(projectRoot, ".agents", "skills", skill, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "<docs_dir>/agents/") || strings.Contains(string(data), "`docs/agents/") {
			t.Fatalf("%s does not resolve configured agent rules:\n%s", skill, data)
		}
	}
}

func TestConfigListIsStableAndIncludesDocsDir(t *testing.T) {
	projectRoot := t.TempDir()
	configDir := filepath.Join(projectRoot, ".flowforge")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("version: 5.0.0\ndocs_dir: wiki\nversion_check: true\nwiki:\n  root: legacy-wiki\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	withWorkingDir(t, projectRoot)

	cmd := newConfigListCmd()
	var output strings.Builder
	cmd.SetOut(&output)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "Project Config:\n  docs_dir = wiki\n  version_check = true\n"; got != want {
		t.Fatalf("config list output:\n%s\nwant:\n%s", got, want)
	}

	setCmd := newConfigSetCmd()
	if err := setCmd.RunE(setCmd, []string{"docs_dir", ""}); err == nil {
		t.Fatal("empty docs_dir must be rejected")
	}
	if err := setCmd.RunE(setCmd, []string{"docs_dir", "custom-docs"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, retained := range []string{"docs_dir: custom-docs", "root: legacy-wiki"} {
		if !strings.Contains(string(data), retained) {
			t.Fatalf("config set failed to preserve %q:\n%s", retained, data)
		}
	}
}

func TestGraphCommandsResolveConfiguredDocsDirFromNestedPath(t *testing.T) {
	projectRoot := t.TempDir()
	configDir := filepath.Join(projectRoot, ".flowforge")
	issuesDir := filepath.Join(projectRoot, "wiki", "proposals", "sample", "issues")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(issuesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("version: 5.0.0\ndocs_dir: wiki\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ticket := "# 01: Configured ticket\n**Blocked by:** None\n**Status:** open\n"
	if err := os.WriteFile(filepath.Join(issuesDir, "01-configured.md"), []byte(ticket), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(projectRoot, "src", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, nested)

	checkCmd := newCheckCmd()
	checkDir, checkJSON, checkStrict = "", false, false
	var checkOutput bytes.Buffer
	checkCmd.SetOut(&checkOutput)
	if err := checkCmd.RunE(checkCmd, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(checkOutput.String(), "Checked 1 issues") {
		t.Fatalf("check did not use configured docs_dir: %s", checkOutput.String())
	}

	frontierCmd := newFrontierCmd()
	frontierDir, frontierJSON, frontierQuiet, frontierStrict, frontierIncludeGaps = "", false, true, false, false
	var frontierOutput bytes.Buffer
	frontierCmd.SetOut(&frontierOutput)
	if err := frontierCmd.RunE(frontierCmd, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(frontierOutput.String(), "01-configured.md") {
		t.Fatalf("frontier did not use configured docs_dir: %s", frontierOutput.String())
	}

	statusCmd := newStatusCmd()
	statusDir = ""
	var statusOutput bytes.Buffer
	statusCmd.SetOut(&statusOutput)
	if err := statusCmd.RunE(statusCmd, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(statusOutput.String(), "sample") {
		t.Fatalf("status did not use configured docs_dir: %s", statusOutput.String())
	}
}

func TestGraphCommandRejectsMalformedProjectConfig(t *testing.T) {
	projectRoot := t.TempDir()
	configDir := filepath.Join(projectRoot, ".flowforge")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("docs_dir: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, projectRoot)
	cmd := newCheckCmd()
	checkDir, checkJSON, checkStrict = "", false, false
	if err := cmd.RunE(cmd, nil); err == nil || !strings.Contains(err.Error(), "loading project configuration") {
		t.Fatalf("malformed config should fail graph resolution, got %v", err)
	}
}

func TestUpgradeAssetSyncUsesConfiguredDocsRoot(t *testing.T) {
	projectRoot := t.TempDir()
	configDir := filepath.Join(projectRoot, ".flowforge")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("version: 5.0.0\ndocs_dir: wiki\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(projectRoot, "src")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, nested)
	cmd := newUpgradeCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	syncProjectAssets(cmd, "synced")
	if _, err := os.Stat(filepath.Join(projectRoot, "wiki", "agents", "issue-tracker.md")); err != nil {
		t.Fatalf("upgrade sync ignored configured docs root: %v", err)
	}
	if !strings.Contains(output.String(), "synced") {
		t.Fatalf("upgrade sync did not report success: %s", output.String())
	}
}
