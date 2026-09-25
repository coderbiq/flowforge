package command

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"flowforge/internal/config"
	"flowforge/internal/subagent"
)

func TestAgentsDeployWritesAllHostsForBuiltinRoles(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}

	if len(deployed) != 10 {
		t.Fatalf("expected 10 deployed subagents, got %d", len(deployed))
	}

	expectedRoles := []string{
		"flowforge-analyst",
		"flowforge-architect",
		"flowforge-batch-analyst",
		"flowforge-executor",
		"flowforge-implementer",
		"flowforge-investigator",
		"flowforge-planner",
		"flowforge-reviewer",
		"flowforge-reviewer-lite",
		"flowforge-scribe",
	}

	// Verify all roles were deployed
	for _, role := range expectedRoles {
		found := false
		for _, name := range deployed {
			if name == role {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected role %q not found in deployed list", role)
		}
	}

	// Verify files exist in all three host directories
	for _, role := range expectedRoles {
		claudePath := filepath.Join(projectRoot, ".claude", "agents", role+".md")
		opencodePath := filepath.Join(projectRoot, ".opencode", "agent", role+".md")
		codexPath := filepath.Join(projectRoot, ".codex", "agents", role+".toml")

		for _, path := range []string{claudePath, opencodePath, codexPath} {
			if _, err := os.Stat(path); err != nil {
				t.Errorf("expected file %s does not exist", path)
			}
		}
	}

	// Verify Claude Code file contains skills: field
	claudeContent, err := os.ReadFile(filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(claudeContent), "skills:") {
		t.Error("Claude Code file missing skills: field")
	}

	// Verify OpenCode file contains mode: subagent
	opencodeContent, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(opencodeContent), "mode: subagent") {
		t.Error("OpenCode file missing mode: subagent")
	}

	// Verify Codex file contains developer_instructions and replaced skill invocation
	codexContent, err := os.ReadFile(filepath.Join(projectRoot, ".codex", "agents", "flowforge-analyst.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(codexContent), "developer_instructions") {
		t.Error("Codex file missing developer_instructions")
	}
	if strings.Contains(string(codexContent), "invoke the Skill tool") {
		t.Error("Codex file still contains 'invoke the Skill tool' (should be replaced)")
	}
}

func TestAgentsDeploySingleName(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	deployed, err := deploySubagents(projectRoot, cfg, "flowforge-planner")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}

	if len(deployed) != 1 {
		t.Fatalf("expected 1 deployed subagent, got %d", len(deployed))
	}

	if deployed[0] != "flowforge-planner" {
		t.Errorf("expected flowforge-planner, got %s", deployed[0])
	}

	// Verify only planner files exist
	claudeDir := filepath.Join(projectRoot, ".claude", "agents")
	entries, err := os.ReadDir(claudeDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 file in .claude/agents, got %d", len(entries))
	}
	if entries[0].Name() != "flowforge-planner.md" {
		t.Errorf("expected flowforge-planner.md, got %s", entries[0].Name())
	}
}

func TestAgentsDeployUnknownNameErrors(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	_, err = deploySubagents(projectRoot, cfg, "nonexistent-agent")
	if err == nil {
		t.Fatal("expected error for unknown agent name, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAgentsDeployIsIdempotent(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// First deployment
	_, err = deploySubagents(projectRoot, cfg, "flowforge-analyst")
	if err != nil {
		t.Fatalf("first deployment: %v", err)
	}

	claudePath := filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md")
	firstContent, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}

	// Second deployment
	_, err = deploySubagents(projectRoot, cfg, "flowforge-analyst")
	if err != nil {
		t.Fatalf("second deployment: %v", err)
	}

	secondContent, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(firstContent) != string(secondContent) {
		t.Error("deploySubagents is not idempotent: content differs between runs")
	}
}

func TestAgentsDeployRespectsDisabledList(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Disable two agents
	cfg.Agents.Disabled = []string{"flowforge-analyst", "flowforge-reviewer"}
	if err := cfg.Save(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Reload config
	cfg, err = config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}

	// Should deploy 8 agents (10 - 2 disabled)
	if len(deployed) != 8 {
		t.Fatalf("expected 8 deployed subagents, got %d", len(deployed))
	}

	// Verify disabled agents are not in the list
	for _, name := range deployed {
		if name == "flowforge-analyst" || name == "flowforge-reviewer" {
			t.Errorf("disabled agent %s was deployed", name)
		}
	}

	// Verify disabled agents' files do not exist
	for _, disabled := range []string{"flowforge-analyst", "flowforge-reviewer"} {
		path := filepath.Join(projectRoot, ".claude", "agents", disabled+".md")
		if _, err := os.Stat(path); err == nil {
			t.Errorf("disabled agent %s file exists at %s", disabled, path)
		}
	}
}

// initializeTestProject creates a minimal FlowForge project structure
func initializeTestProject(projectRoot string) error {
	cfg := &config.Config{
		Version:      "5",
		VersionCheck: true,
		DocsDir:      "docs",
	}
	if err := cfg.Save(projectRoot); err != nil {
		return err
	}

	docsDir := filepath.Join(projectRoot, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	return nil
}

func TestAgentsRemoveBuiltinPersistsDisabled(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Deploy all subagents first
	_, err = deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	// Verify analyst exists
	claudePath := filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md")
	if _, err := os.Stat(claudePath); err != nil {
		t.Fatal("expected flowforge-analyst.md to exist before removal")
	}

	// Remove built-in subagent
	isBuiltin, removedPaths, err := removeSubagent(projectRoot, "flowforge-analyst")
	if err != nil {
		t.Fatalf("removeSubagent: %v", err)
	}

	if !isBuiltin {
		t.Error("expected flowforge-analyst to be identified as built-in")
	}

	if len(removedPaths) != 4 {
		t.Errorf("expected 4 removed paths, got %d", len(removedPaths))
	}

	// Verify files were removed
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("expected flowforge-analyst.md to be removed")
	}

	// Verify disabled list was updated
	cfg, err = config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, disabled := range cfg.Agents.Disabled {
		if disabled == "flowforge-analyst" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected flowforge-analyst in cfg.Agents.Disabled")
	}

	// Verify deploy skips disabled agent
	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range deployed {
		if name == "flowforge-analyst" {
			t.Error("disabled subagent flowforge-analyst was redeployed")
		}
	}

	// Verify analyst file still doesn't exist after deploy
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("disabled subagent file should not be recreated by deploy")
	}
}

func TestAgentsRemoveCustomDeletesSourceFile(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Create custom subagent source
	customDir := filepath.Join(projectRoot, config.ConfigDirName, "subagents")
	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatal(err)
	}

	customSource := filepath.Join(customDir, "my-custom.md")
	customContent := `---
flowforge_agent:
  name: my-custom
  description: Test custom agent
  model_profile: tool-capable
  default_skill: flowforge-align
  detour_skills: []
  permission: workspace-write
  after: []
  before: []
  returns_to: []
---

## Identity
Test.

## Boundaries
Test.

## Workflow Position
Test.

## Default Skill
Test.

## Result Contract
Test.
`
	if err := os.WriteFile(customSource, []byte(customContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Deploy custom subagent
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	_, err = deploySubagents(projectRoot, cfg, "my-custom")
	if err != nil {
		t.Fatal(err)
	}

	// Verify deployed files exist
	claudePath := filepath.Join(projectRoot, ".claude", "agents", "my-custom.md")
	if _, err := os.Stat(claudePath); err != nil {
		t.Fatal("expected my-custom.md to exist after deployment")
	}

	// Remove custom subagent
	isBuiltin, removedPaths, err := removeSubagent(projectRoot, "my-custom")
	if err != nil {
		t.Fatalf("removeSubagent: %v", err)
	}

	if isBuiltin {
		t.Error("expected my-custom to be identified as custom (not built-in)")
	}

	// Should remove 4 host files + 1 source file = 5 total
	if len(removedPaths) != 5 {
		t.Errorf("expected 5 removed paths (4 hosts + source), got %d", len(removedPaths))
	}

	// Verify source file was deleted
	if _, err := os.Stat(customSource); err == nil {
		t.Error("expected custom source file to be deleted")
	}

	// Verify host files were deleted
	if _, err := os.Stat(claudePath); err == nil {
		t.Error("expected my-custom.md to be removed from .claude/agents/")
	}
}

func TestAgentsRemoveUnknownCustomNameErrors(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Create custom subagents directory but no files
	customDir := filepath.Join(projectRoot, config.ConfigDirName, "subagents")
	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Try to remove non-existent custom subagent
	_, _, err := removeSubagent(projectRoot, "nonexistent-custom")
	if err == nil {
		t.Fatal("expected error for non-existent custom subagent, got nil")
	}

	if !strings.Contains(err.Error(), "source file not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAgentsStatusReportsCurrentAfterDeploy(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Deploy all subagents
	_, err = deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	result, err := computeSubagentStatus(projectRoot, cfg)
	if err != nil {
		t.Fatalf("computeSubagentStatus: %v", err)
	}

	if !result.Current {
		t.Errorf("expected all current, but got non-current. Entries: %+v", result.Entries)
	}

	for _, entry := range result.Entries {
		if entry.State != string(managedAssetCurrent) {
			t.Errorf("expected %s to be current, got %s", entry.Target, entry.State)
		}
	}
}

func TestAgentsStatusReportsMissing(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Deploy all subagents
	_, err = deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	// Remove one file
	missingPath := filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md")
	if err := os.Remove(missingPath); err != nil {
		t.Fatal(err)
	}

	result, err := computeSubagentStatus(projectRoot, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Current {
		t.Error("expected non-current result after removing a file")
	}

	foundMissing := false
	for _, entry := range result.Entries {
		if entry.Target == missingPath && entry.State == string(managedAssetMissing) {
			foundMissing = true
		}
	}
	if !foundMissing {
		t.Error("expected to find missing entry for removed file")
	}
}

func TestAgentsStatusReportsDrifted(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Deploy all subagents
	_, err = deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	// Drift one file
	driftedPath := filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md")
	if err := os.WriteFile(driftedPath, []byte("modified content"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := computeSubagentStatus(projectRoot, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Current {
		t.Error("expected non-current result after drifting a file")
	}

	foundDrifted := false
	for _, entry := range result.Entries {
		if entry.Target == driftedPath && entry.State == string(managedAssetDrifted) {
			foundDrifted = true
		}
	}
	if !foundDrifted {
		t.Error("expected to find drifted entry for modified file")
	}
}

func TestAgentsStatusReportsProjectOwned(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Deploy all subagents
	_, err = deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatal(err)
	}

	// Add an extra file
	extraPath := filepath.Join(projectRoot, ".claude", "agents", "extra-project-file.md")
	if err := os.WriteFile(extraPath, []byte("project owned"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := computeSubagentStatus(projectRoot, cfg)
	if err != nil {
		t.Fatal(err)
	}

	foundProjectOwned := false
	for _, entry := range result.Entries {
		if entry.Target == extraPath && entry.State == string(managedAssetProjectOwned) {
			foundProjectOwned = true
		}
	}
	if !foundProjectOwned {
		t.Error("expected to find project-owned entry for extra file")
	}

	// Project-owned should not make overall result non-current
	if !result.Current {
		t.Error("expected current=true when only project-owned extra files exist (no missing/drifted)")
	}
}

func TestInitDeploysSubagentsToAllHosts(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate what init does: deploy managed assets + subagents
	if err := deployManagedAssets(projectRoot, filepath.Join(projectRoot, "docs")); err != nil {
		t.Fatal(err)
	}
	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}

	if len(deployed) != 10 {
		t.Fatalf("expected 10 subagents deployed, got %d", len(deployed))
	}

	// Verify all 3 host directories have 10 files each
	for _, dir := range []string{".claude/agents", ".opencode/agent", ".codex/agents"} {
		entries, err := os.ReadDir(filepath.Join(projectRoot, dir))
		if err != nil {
			t.Fatalf("reading %s: %v", dir, err)
		}
		if len(entries) != 10 {
			t.Errorf("expected 10 files in %s, got %d", dir, len(entries))
		}
	}
}

func TestUpgradeSyncDeploysSubagents(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Simulate what upgrade's syncProjectAssets does
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	if err := deployManagedAssets(projectRoot, cfg.DocsRoot(projectRoot)); err != nil {
		t.Fatal(err)
	}

	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}

	if len(deployed) != 10 {
		t.Fatalf("expected 10 subagents, got %d", len(deployed))
	}

	// Verify files exist for representative roles
	for _, role := range []string{"flowforge-analyst", "flowforge-architect"} {
		paths := []string{
			filepath.Join(projectRoot, ".claude", "agents", role+".md"),
			filepath.Join(projectRoot, ".opencode", "agent", role+".md"),
			filepath.Join(projectRoot, ".codex", "agents", role+".toml"),
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err != nil {
				t.Errorf("expected file %s to exist", p)
			}
		}
	}
}

func TestUpgradeSyncSkipsDisabledSubagents(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	// Disable a subagent
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Disabled = []string{"flowforge-planner"}
	if err := cfg.Save(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err = config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate sync
	if err := deployManagedAssets(projectRoot, cfg.DocsRoot(projectRoot)); err != nil {
		t.Fatal(err)
	}
	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}

	// Should deploy 9 (10 - 1 disabled)
	if len(deployed) != 9 {
		t.Fatalf("expected 9 subagents, got %d", len(deployed))
	}

	// Verify disabled agent not deployed
	for _, name := range deployed {
		if name == "flowforge-planner" {
			t.Error("disabled subagent flowforge-planner was deployed")
		}
	}

	// Verify planner files don't exist
	plannerPath := filepath.Join(projectRoot, ".claude", "agents", "flowforge-planner.md")
	if _, err := os.Stat(plannerPath); err == nil {
		t.Error("disabled subagent file should not exist")
	}
}

// piExtensionSource returns the managed extension source from the located
// assets dir, mirroring what deployPiExtension reads.
func piExtensionSource(t *testing.T) []byte {
	t.Helper()
	assetsDir, cleanup, err := locateAssetsDir()
	if err != nil {
		t.Fatalf("locating assets: %v", err)
	}
	defer cleanup()
	data, err := os.ReadFile(filepath.Join(assetsDir, "pi", "flowforge.ts"))
	if err != nil {
		t.Fatalf("reading pi extension source: %v", err)
	}
	return data
}

func TestAgentsDeployPiExtension(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"pi"}

	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}

	extPath := filepath.Join(projectRoot, ".pi", "extensions", "flowforge.ts")
	deployedContent, err := os.ReadFile(extPath)
	if err != nil {
		t.Fatalf("expected %s after pi deploy: %v", extPath, err)
	}
	if want := piExtensionSource(t); !bytes.Equal(deployedContent, want) {
		t.Error("deployed extension content differs from assets/pi/flowforge.ts source")
	}

	// Idempotent redeploy
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatalf("redeploy: %v", err)
	}
	again, err := os.ReadFile(extPath)
	if err != nil {
		t.Fatalf("re-reading %s: %v", extPath, err)
	}
	if !bytes.Equal(again, deployedContent) {
		t.Error("extension deploy is not idempotent")
	}
}

func TestAgentsDeployCleansDeselectedPiExtension(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	// Project-owned file inside .pi/extensions must survive host cleanup
	ownedPath := filepath.Join(projectRoot, ".pi", "extensions", "my-own-extension.ts")
	if err := os.WriteFile(ownedPath, []byte("// mine\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg.Agents.Hosts = []string{"opencode"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatalf("redeploy without pi host: %v", err)
	}

	managed := filepath.Join(projectRoot, ".pi", "extensions", "flowforge.ts")
	if _, err := os.Stat(managed); err == nil {
		t.Errorf("managed extension %s should have been cleaned from deselected host", managed)
	}
	if _, err := os.Stat(ownedPath); err != nil {
		t.Errorf("project-owned file %s must be preserved: %v", ownedPath, err)
	}
}

func TestAgentsDeployHonorsHostSelection(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}

	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}
	if len(deployed) != 10 {
		t.Fatalf("expected 10 deployed subagents, got %d", len(deployed))
	}

	// Selected host received files
	if _, err := os.Stat(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md")); err != nil {
		t.Errorf("expected .opencode/agent/flowforge-analyst.md: %v", err)
	}

	// Deselected host directories must not be created
	for _, dir := range []string{filepath.Join(projectRoot, ".claude", "agents"), filepath.Join(projectRoot, ".codex", "agents")} {
		if _, err := os.Stat(dir); err == nil {
			t.Errorf("deselected host directory %s was created", dir)
		}
	}
}

func TestAgentsDeployCleansDeselectedHosts(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	// Project-owned file inside a host directory that will be deselected
	ownedPath := filepath.Join(projectRoot, ".claude", "agents", "my-own-agent.md")
	if err := os.WriteFile(ownedPath, []byte("# mine\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg.Agents.Hosts = []string{"opencode"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatalf("redeploy with narrowed hosts: %v", err)
	}

	// Managed files removed from deselected hosts
	for _, path := range []string{
		filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md"),
		filepath.Join(projectRoot, ".codex", "agents", "flowforge-analyst.toml"),
	} {
		if _, err := os.Stat(path); err == nil {
			t.Errorf("managed file %s should have been cleaned from deselected host", path)
		}
	}

	// Project-owned file preserved
	if _, err := os.Stat(ownedPath); err != nil {
		t.Errorf("project-owned file %s must be preserved: %v", ownedPath, err)
	}
}

func TestAgentsDeployPiHost(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"pi"}

	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploySubagents: %v", err)
	}
	if len(deployed) != 10 {
		t.Fatalf("expected 10 deployed subagents, got %d", len(deployed))
	}

	// .pi/agents/ holds one native file per definition
	piDir := filepath.Join(projectRoot, ".pi", "agents")
	entries, err := os.ReadDir(piDir)
	if err != nil {
		t.Fatalf("reading %s: %v", piDir, err)
	}
	if len(entries) != 10 {
		t.Errorf("expected 10 files in .pi/agents, got %d", len(entries))
	}

	// Non-read-only role: thinking + skills binding, no tools allowlist
	analystPath := filepath.Join(piDir, "flowforge-analyst.md")
	analystContent, err := os.ReadFile(analystPath)
	if err != nil {
		t.Fatal(err)
	}
	analyst := string(analystContent)
	if !strings.Contains(analyst, "name: flowforge-analyst") {
		t.Error("PI file missing 'name: flowforge-analyst'")
	}
	if !strings.Contains(analyst, "thinking: high") {
		t.Error("PI analyst file missing 'thinking: high'")
	}
	if !strings.Contains(analyst, "inheritSkills: false") {
		t.Error("PI analyst file missing 'inheritSkills: false'")
	}
	if strings.Contains(analyst, "tools:") {
		t.Error("PI analyst file must not carry a tools allowlist")
	}
	if strings.Contains(analyst, "model:") {
		t.Error("PI file must not carry a model field (inherit parent session model)")
	}

	// Read-only role carries the strict read-tool allowlist
	investigatorContent, err := os.ReadFile(filepath.Join(piDir, "flowforge-investigator.md"))
	if err != nil {
		t.Fatal(err)
	}
	investigator := string(investigatorContent)
	if !strings.Contains(investigator, "tools:") {
		t.Error("PI investigator file missing tools allowlist")
	}
	if !strings.Contains(investigator, "- read") {
		t.Error("PI investigator file tools allowlist missing 'read'")
	}
	if !strings.Contains(investigator, "skills:") {
		t.Error("PI investigator file missing skills binding")
	}

	// Other host directories must not be created
	for _, dir := range []string{
		filepath.Join(projectRoot, ".claude", "agents"),
		filepath.Join(projectRoot, ".opencode", "agent"),
		filepath.Join(projectRoot, ".codex", "agents"),
	} {
		if _, err := os.Stat(dir); err == nil {
			t.Errorf("deselected host directory %s was created", dir)
		}
	}
}

func TestAgentsDeployCleansPiWhenDeselected(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	// Default nil hosts deploys all hosts including pi.
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	// Project-owned file inside .pi/agents that will be deselected
	ownedPath := filepath.Join(projectRoot, ".pi", "agents", "my-own-agent.md")
	if err := os.WriteFile(ownedPath, []byte("# mine\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg.Agents.Hosts = []string{"opencode"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatalf("redeploy with narrowed hosts: %v", err)
	}

	// Managed .pi/agents files removed
	for _, role := range []string{"flowforge-analyst", "flowforge-investigator"} {
		path := filepath.Join(projectRoot, ".pi", "agents", role+".md")
		if _, err := os.Stat(path); err == nil {
			t.Errorf("managed file %s should have been cleaned from deselected pi host", path)
		}
	}

	// Project-owned file preserved
	if _, err := os.Stat(ownedPath); err != nil {
		t.Errorf("project-owned file %s must be preserved: %v", ownedPath, err)
	}
}

func TestAgentsHostsValidation(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"vscode"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err == nil {
		t.Error("unknown host name must fail deployment")
	}

	cfg.Agents.Hosts = []string{}
	if _, err := deploySubagents(projectRoot, cfg, ""); err == nil {
		t.Error("empty host list must fail deployment")
	}
}

func TestAgentsStatusScopedToSelectedHosts(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	// Stale managed file in a deselected host must not affect status
	staleDir := filepath.Join(projectRoot, ".claude", "agents")
	if err := os.MkdirAll(staleDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staleDir, "flowforge-analyst.md"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := computeSubagentStatus(projectRoot, cfg)
	if err != nil {
		t.Fatalf("computeSubagentStatus: %v", err)
	}
	if !result.Current {
		t.Error("status must be current when only the selected host is in scope")
	}
	for _, entry := range result.Entries {
		if filepath.Base(filepath.Dir(filepath.Dir(entry.Target))) == ".claude" {
			t.Errorf("status reported deselected host target %s", entry.Target)
		}
	}
}

func TestAgentsRemoveScopedToSelectedHosts(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	if err := cfg.Save(projectRoot); err != nil {
		t.Fatal(err)
	}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	// File in a deselected host: remove must not touch it
	foreignDir := filepath.Join(projectRoot, ".claude", "agents")
	if err := os.MkdirAll(foreignDir, 0755); err != nil {
		t.Fatal(err)
	}
	foreignPath := filepath.Join(foreignDir, "flowforge-analyst.md")
	if err := os.WriteFile(foreignPath, []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	_, removedPaths, err := removeSubagent(projectRoot, "flowforge-analyst")
	if err != nil {
		t.Fatalf("removeSubagent: %v", err)
	}

	for _, path := range removedPaths {
		if filepath.Base(filepath.Dir(filepath.Dir(path))) == ".claude" {
			t.Errorf("remove touched deselected host file %s", path)
		}
	}
	if _, err := os.Stat(foreignPath); err != nil {
		t.Errorf("deselected host file must be left for deploy-time cleanup, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md")); err == nil {
		t.Error("selected host file was not removed")
	}
}

func TestOpenCodeModelPinning(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	cfg.Agents.Models = map[string]string{"tool-capable": "erasebg-gemini/gemini-3.8-flash-high"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	implementer, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(implementer), "model: erasebg-gemini/gemini-3.8-flash-high") {
		t.Error("tool-capable implementer must carry pinned model when configured")
	}

	analyst, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(analyst), "model:") {
		t.Error("unconfigured profile must keep inherit behavior (no model field)")
	}

	// Without Models config nothing carries a model field
	root2 := t.TempDir()
	if err := initializeTestProject(root2); err != nil {
		t.Fatal(err)
	}
	cfg2, err := config.Load(root2)
	if err != nil {
		t.Fatal(err)
	}
	cfg2.Agents.Hosts = []string{"opencode"}
	if _, err := deploySubagents(root2, cfg2, ""); err != nil {
		t.Fatal(err)
	}
	impl2, err := os.ReadFile(filepath.Join(root2, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(impl2), "model:") {
		t.Error("unconfigured Models must produce byte-compatible output (no model field)")
	}
}

func TestOpenCodeTestGuardDefaults(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	implementer, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{"permission:", "edit:", "'**/*_test.go': deny", "'**/src/test/**': deny", "'**/src/integrationTest/**': deny"} {
		if !strings.Contains(string(implementer), needle) {
			t.Errorf("implementer default test guard missing %q", needle)
		}
	}

	analyst, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(analyst), "permission:") {
		t.Error("guard applies only to flowforge-implementer")
	}
}

func TestOpenCodeTestGuardDisabledAndCustom(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	cfg.Agents.DisableTestGuard = true
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}
	impl, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(impl), "edit:") {
		t.Error("DisableTestGuard must remove the edit permission block")
	}
	if !strings.Contains(string(impl), "question: deny") {
		t.Error("implementer must always carry question: deny (subagents never pause for input)")
	}

	root2 := t.TempDir()
	if err := initializeTestProject(root2); err != nil {
		t.Fatal(err)
	}
	cfg2, err := config.Load(root2)
	if err != nil {
		t.Fatal(err)
	}
	cfg2.Agents.Hosts = []string{"opencode"}
	cfg2.Agents.TestFileGlobs = []string{"**/custom/**"}
	if _, err := deploySubagents(root2, cfg2, ""); err != nil {
		t.Fatal(err)
	}
	impl2, err := os.ReadFile(filepath.Join(root2, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(impl2), "'**/custom/**': deny") || strings.Contains(string(impl2), "_test.go") {
		t.Error("TestFileGlobs must override the default glob set")
	}
}

func TestOpenCodeStepsBudget(t *testing.T) {
	// deployOpenCode deploys all subagents to the opencode host with the given
	// agents.max_steps value and returns the implementer and analyst artifacts.
	deployOpenCode := func(t *testing.T, maxSteps int) (string, string) {
		t.Helper()
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		cfg.Agents.MaxSteps = maxSteps
		if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
			t.Fatalf("deploySubagents(maxSteps=%d): %v", maxSteps, err)
		}
		readArtifact := func(name string) string {
			t.Helper()
			data, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", name+".md"))
			if err != nil {
				t.Fatal(err)
			}
			return string(data)
		}
		return readArtifact("flowforge-implementer"), readArtifact("flowforge-analyst")
	}

	t.Run("unconfigured defaults to 200 on implementer", func(t *testing.T) {
		impl, _ := deployOpenCode(t, 0)
		if !strings.Contains(impl, "steps: 200") {
			t.Error("unconfigured max_steps must compile implementer with default steps: 200")
		}
	})

	t.Run("positive value is passed through", func(t *testing.T) {
		impl, _ := deployOpenCode(t, 500)
		if !strings.Contains(impl, "steps: 500") {
			t.Error("max_steps: 500 must compile implementer with steps: 500")
		}
	})

	t.Run("minus one disables the steps field", func(t *testing.T) {
		impl, _ := deployOpenCode(t, -1)
		if strings.Contains(impl, "steps:") {
			t.Error("max_steps: -1 must not produce a steps field")
		}
	})

	t.Run("other negative values are config errors", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		cfg.Agents.MaxSteps = -5
		if _, err := deploySubagents(projectRoot, cfg, ""); err == nil {
			t.Error("max_steps: -5 must fail deployment with a config error")
		} else if !strings.Contains(err.Error(), "max_steps") {
			t.Errorf("error must name the config key: %v", err)
		}
		if _, err := computeSubagentStatus(projectRoot, cfg); err == nil {
			t.Error("max_steps: -5 must fail status with the same config error")
		}
	})

	t.Run("non-implementer agents never carry steps", func(t *testing.T) {
		_, analyst := deployOpenCode(t, 500)
		if strings.Contains(analyst, "steps:") {
			t.Error("max_steps must not affect non-implementer subagents")
		}
	})

	t.Run("unconfigured implementer differs from legacy by exactly the steps line", func(t *testing.T) {
		implDefault, analystDefault := deployOpenCode(t, 0)
		implOff, analystOff := deployOpenCode(t, -1)
		stripped := strings.Replace(implDefault, "steps: 200\n", "", 1)
		if stripped != implOff {
			t.Error("unconfigured implementer output must differ from max_steps: -1 output by exactly one 'steps: 200' line")
		}
		if analystDefault != analystOff {
			t.Error("non-implementer artifacts must be byte-identical regardless of max_steps")
		}
	})
}

func TestAgentsModelsValidation(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Models = map[string]string{"sonnet": "x"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err == nil {
		t.Error("unknown agents.models key must fail")
	}
}

func TestAgentsModelByNameOverridesProfile(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	cfg.Agents.Models = map[string]string{"tool-capable": "profile-model/x"}
	cfg.Agents.ModelOverrides = map[string]string{"flowforge-planner": "name-model/y"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}

	// The name key wins over the profile key for the pinned agent.
	planner, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-planner.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(planner), "model: name-model/y") {
		t.Error("models_by_name key must override the profile key for that agent")
	}
	if strings.Contains(string(planner), "profile-model/x") {
		t.Error("profile model must not leak into a name-pinned agent")
	}

	// A same-profile neighbor keeps the profile-pinned model.
	impl, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(impl), "model: profile-model/x") {
		t.Error("same-profile neighbor without a name key must keep the profile model")
	}
}

func TestAgentsModelByNameUnknownErrors(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	cfg.Agents.ModelOverrides = map[string]string{"no-such-agent": "m/x"}
	_, err = deploySubagents(projectRoot, cfg, "")
	if err == nil {
		t.Fatal("expected error for unknown models_by_name key, got nil")
	}
	if !strings.Contains(err.Error(), `agents.models_by_name: unknown agent "no-such-agent"`) {
		t.Errorf("unexpected error message: %v", err)
	}

	// The failure must precede any artifact write: not even the host
	// directory may exist.
	if _, err := os.Stat(filepath.Join(projectRoot, ".opencode", "agent")); err == nil {
		t.Error("failed deploy must not create host directories or artifacts")
	}

	// Status reports the same config error as deploy (parity with the
	// max_steps unknown-value convention).
	if _, err := computeSubagentStatus(projectRoot, cfg); err == nil {
		t.Error("unknown models_by_name key must fail status with the same config error")
	} else if !strings.Contains(err.Error(), `agents.models_by_name: unknown agent "no-such-agent"`) {
		t.Errorf("status error must match deploy error, got: %v", err)
	}
}

func TestAgentsProfileKeyOnlyUnchanged(t *testing.T) {
	deployProfileOnly := func(t *testing.T, overrides map[string]string) string {
		t.Helper()
		root := t.TempDir()
		if err := initializeTestProject(root); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(root)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		cfg.Agents.Models = map[string]string{"tool-capable": "profile-model/x"}
		cfg.Agents.ModelOverrides = overrides
		if _, err := deploySubagents(root, cfg, ""); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(root, ".opencode", "agent", "flowforge-implementer.md"))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	nilOverrides := deployProfileOnly(t, nil)
	emptyOverrides := deployProfileOnly(t, map[string]string{})
	if nilOverrides != emptyOverrides {
		t.Error("empty models_by_name must behave like an absent one (no output change)")
	}

	// Profile-only artifacts must equal a direct config-driven compile
	// byte for byte: adding the models_by_name mechanism is regression-free
	// for configurations that do not use it.
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	cfg.Agents.Models = map[string]string{"tool-capable": "profile-model/x"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}
	def := findDiscoveredDefinition(t, projectRoot, "flowforge-implementer")
	opts, err := resolveCompileOptions(cfg, def, "opencode")
	if err != nil {
		t.Fatal(err)
	}
	want, err := subagent.CompileOpenCodeWithOptions(def, opts)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("profile-only deploy must equal the config-driven compile\n got: %q\nwant: %q", got, want)
	}
	if nilOverrides != string(want) {
		t.Error("profile-only output must be identical with and without the new config field")
	}
}

func TestStatusUsesSameCompileOptions(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	cfg.Agents.Models = map[string]string{"tool-capable": "erasebg-gemini/gemini-3.8-flash-high"}
	cfg.Agents.MaxSteps = 200
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}
	result, err := computeSubagentStatus(projectRoot, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Current {
		t.Error("status must use the same compile options as deploy (no false drift)")
	}
}

// captureStderr runs fn with os.Stderr redirected into a pipe and returns
// everything fn wrote there.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stderr = w
	defer func() { os.Stderr = old }()
	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("closing stderr pipe: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading stderr pipe: %v", err)
	}
	return string(data)
}

// findDiscoveredDefinition returns the discovered definition with the given
// name, so tests can byte-compare deployed files against fresh compiles.
func findDiscoveredDefinition(t *testing.T, projectRoot, name string) *subagent.Definition {
	t.Helper()
	defs, err := discoverSubagentSources(projectRoot)
	if err != nil {
		t.Fatalf("discoverSubagentSources: %v", err)
	}
	for _, def := range defs {
		if def.Name == name {
			return def
		}
	}
	t.Fatalf("definition %q not found among discovered sources", name)
	return nil
}

func TestDeployPreservesLocalModel(t *testing.T) {
	t.Run("opencode preserves hand-edited model with stderr hint", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		if _, err := deploySubagents(projectRoot, cfg, "flowforge-analyst"); err != nil {
			t.Fatal(err)
		}

		// Hand-edit a local model into the deployed frontmatter.
		path := filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		handEdited := strings.Replace(string(data), "---\n", "---\nmodel: custom-model/x\n", 1)
		if err := os.WriteFile(path, []byte(handEdited), 0644); err != nil {
			t.Fatal(err)
		}

		var deployErr error
		stderr := captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "flowforge-analyst")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		def := findDiscoveredDefinition(t, projectRoot, "flowforge-analyst")
		want, err := subagent.CompileOpenCodeWithOptions(def, subagent.CompileOptions{FallbackModel: "custom-model/x"})
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Errorf("redeployed file must equal fresh compile with the preserved model as fallback\n got: %q\nwant: %q", got, want)
		}

		wantHint := "  info: preserved local model \"custom-model/x\" for .opencode/agent/flowforge-analyst.md (set agents.models in .flowforge/config.yaml to pin explicitly)\n"
		if !strings.Contains(stderr, wantHint) {
			t.Errorf("stderr must carry the preserved hint\ngot:  %q\nwant: %q", stderr, wantHint)
		}
	})

	t.Run("config model wins without hint", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		cfg.Agents.Models = map[string]string{"tool-capable": "pinned-by-config/y"}
		if _, err := deploySubagents(projectRoot, cfg, "flowforge-implementer"); err != nil {
			t.Fatal(err)
		}

		// Replace the pinned model with a local residue value.
		path := filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		handEdited := strings.Replace(string(data), "model: pinned-by-config/y", "model: stale-local/x", 1)
		if handEdited == string(data) {
			t.Fatal("expected the deployed file to carry the pinned model")
		}
		if err := os.WriteFile(path, []byte(handEdited), 0644); err != nil {
			t.Fatal(err)
		}

		var deployErr error
		stderr := captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "flowforge-implementer")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(got), "stale-local/x") {
			t.Error("config-pinned model must win over the file residue")
		}
		def := findDiscoveredDefinition(t, projectRoot, "flowforge-implementer")
		opts, err := resolveCompileOptions(cfg, def, "opencode")
		if err != nil {
			t.Fatal(err)
		}
		want, err := subagent.CompileOpenCodeWithOptions(def, opts)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Errorf("redeployed file must equal the config-driven compile\n got: %q\nwant: %q", got, want)
		}
		if strings.Contains(stderr, "preserved local model") {
			t.Errorf("config-pinned model must not emit a preserved hint, got: %q", stderr)
		}
	})

	t.Run("claude config model wins without hint", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"claude"}
		cfg.Agents.Models = map[string]string{"tool-capable": "pinned-claude/y"}
		if _, err := deploySubagents(projectRoot, cfg, "flowforge-implementer"); err != nil {
			t.Fatal(err)
		}

		// Replace the pinned model with a local residue value.
		path := filepath.Join(projectRoot, ".claude", "agents", "flowforge-implementer.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		handEdited := strings.Replace(string(data), "model: pinned-claude/y", "model: stale-local/w", 1)
		if handEdited == string(data) {
			t.Fatal("expected the deployed claude file to carry the pinned model")
		}
		if err := os.WriteFile(path, []byte(handEdited), 0644); err != nil {
			t.Fatal(err)
		}

		var deployErr error
		stderr := captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "flowforge-implementer")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(got), "stale-local/w") {
			t.Error("config-pinned model must win over the claude file residue")
		}
		def := findDiscoveredDefinition(t, projectRoot, "flowforge-implementer")
		opts, err := resolveCompileOptions(cfg, def, "opencode")
		if err != nil {
			t.Fatal(err)
		}
		want, err := subagent.CompileClaudeCodeWithOptions(def, opts)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Errorf("redeployed claude file must equal the config-driven compile\n got: %q\nwant: %q", got, want)
		}
		if strings.Contains(stderr, "preserved local model") {
			t.Errorf("config-pinned model must not emit a preserved hint, got: %q", stderr)
		}
	})

	t.Run("no spurious model when existing file has none", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		if _, err := deploySubagents(projectRoot, cfg, "flowforge-analyst"); err != nil {
			t.Fatal(err)
		}

		var deployErr error
		stderr := captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "flowforge-analyst")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}

		got, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md"))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(got), "model:") {
			t.Error("redeploy must not introduce a model field out of thin air")
		}
		if strings.Contains(stderr, "preserved local model") {
			t.Errorf("no local model to preserve, unexpected hint: %q", stderr)
		}
	})

	t.Run("claude default model causes no hint but custom model is preserved", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"claude"}
		if _, err := deploySubagents(projectRoot, cfg, "flowforge-analyst"); err != nil {
			t.Fatal(err)
		}
		def := findDiscoveredDefinition(t, projectRoot, "flowforge-analyst")
		profileDefault := def.ModelProfile.ClaudeModel()

		// Redeploying over our own profile-default model must stay silent:
		// the value came from the previous deploy, not from local custom.
		var deployErr error
		stderr := captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "flowforge-analyst")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}
		if strings.Contains(stderr, "preserved local model") {
			t.Errorf("profile-default model is not a local customization, unexpected hint: %q", stderr)
		}

		// Hand-edit a custom model; it must be preserved with a hint.
		path := filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		handEdited := strings.Replace(string(data), "model: "+profileDefault, "model: claude-custom/z", 1)
		if handEdited == string(data) {
			t.Fatalf("expected deployed claude file to carry model: %s", profileDefault)
		}
		if err := os.WriteFile(path, []byte(handEdited), 0644); err != nil {
			t.Fatal(err)
		}

		stderr = captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "flowforge-analyst")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		want, err := subagent.CompileClaudeCodeWithOptions(def, subagent.CompileOptions{FallbackModel: "claude-custom/z"})
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Errorf("redeployed claude file must equal fresh compile with the preserved model as fallback\n got: %q\nwant: %q", got, want)
		}
		wantHint := "  info: preserved local model \"claude-custom/z\" for .claude/agents/flowforge-analyst.md (set agents.models in .flowforge/config.yaml to pin explicitly)\n"
		if !strings.Contains(stderr, wantHint) {
			t.Errorf("stderr must carry the preserved hint\ngot:  %q\nwant: %q", stderr, wantHint)
		}
	})

	t.Run("first deploy emits no hint and keeps fresh output", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		var deployErr error
		stderr := captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}
		if strings.Contains(stderr, "preserved local model") {
			t.Errorf("first deploy has nothing to preserve, unexpected hint: %q", stderr)
		}

		opencodeAnalyst, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md"))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(opencodeAnalyst), "model:") {
			t.Error("first deploy must keep the opencode inherit behavior (no model field)")
		}
		claudeAnalyst, err := os.ReadFile(filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md"))
		if err != nil {
			t.Fatal(err)
		}
		def := findDiscoveredDefinition(t, projectRoot, "flowforge-analyst")
		if !strings.Contains(string(claudeAnalyst), "model: "+def.ModelProfile.ClaudeModel()) {
			t.Error("first deploy must keep the claude profile default model")
		}
	})

	t.Run("codex redeploy stays byte-identical", func(t *testing.T) {
		root1 := t.TempDir()
		if err := initializeTestProject(root1); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(root1)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := deploySubagents(root1, cfg, ""); err != nil {
			t.Fatal(err)
		}

		// Plant a TOML model-looking key in the existing codex file.
		codexPath := filepath.Join(root1, ".codex", "agents", "flowforge-analyst.toml")
		if err := os.WriteFile(codexPath, []byte("model = \"hacked\"\ndeveloper_instructions = \"stale\"\n"), 0644); err != nil {
			t.Fatal(err)
		}

		var deployErr error
		stderr := captureStderr(t, func() {
			_, deployErr = deploySubagents(root1, cfg, "")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}
		if strings.Contains(stderr, ".codex/") && strings.Contains(stderr, "preserved local model") {
			t.Errorf("codex host has no model concept, unexpected hint: %q", stderr)
		}

		root2 := t.TempDir()
		if err := initializeTestProject(root2); err != nil {
			t.Fatal(err)
		}
		cfg2, err := config.Load(root2)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := deploySubagents(root2, cfg2, ""); err != nil {
			t.Fatal(err)
		}
		redeployed, err := os.ReadFile(codexPath)
		if err != nil {
			t.Fatal(err)
		}
		fresh, err := os.ReadFile(filepath.Join(root2, ".codex", "agents", "flowforge-analyst.toml"))
		if err != nil {
			t.Fatal(err)
		}
		if string(redeployed) != string(fresh) {
			t.Errorf("codex redeploy output must be byte-identical to a fresh deploy\n got: %q\nwant: %q", redeployed, fresh)
		}
	})

	t.Run("hand-edited model survives init --force style redeploy", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		// The init/upgrade pipeline: managed assets + subagents to all hosts.
		if err := deployManagedAssets(projectRoot, filepath.Join(projectRoot, "docs")); err != nil {
			t.Fatal(err)
		}
		if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
			t.Fatal(err)
		}

		// Hand-edit models in both frontmatter hosts.
		opencodePath := filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md")
		opencodeData, err := os.ReadFile(opencodePath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(opencodePath, []byte(strings.Replace(string(opencodeData), "---\n", "---\nmodel: custom-model/x\n", 1)), 0644); err != nil {
			t.Fatal(err)
		}
		def := findDiscoveredDefinition(t, projectRoot, "flowforge-analyst")
		claudePath := filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md")
		claudeData, err := os.ReadFile(claudePath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(claudePath, []byte(strings.Replace(string(claudeData), "model: "+def.ModelProfile.ClaudeModel(), "model: claude-custom/z", 1)), 0644); err != nil {
			t.Fatal(err)
		}

		// Re-run the pipeline as init --force does.
		var deployErr error
		stderr := captureStderr(t, func() {
			deployErr = deployManagedAssets(projectRoot, filepath.Join(projectRoot, "docs"))
			if deployErr != nil {
				return
			}
			_, deployErr = deploySubagents(projectRoot, cfg, "")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}

		opencodeGot, err := os.ReadFile(opencodePath)
		if err != nil {
			t.Fatal(err)
		}
		opencodeWant, err := subagent.CompileOpenCodeWithOptions(def, subagent.CompileOptions{FallbackModel: "custom-model/x"})
		if err != nil {
			t.Fatal(err)
		}
		if string(opencodeGot) != string(opencodeWant) {
			t.Errorf("opencode model must survive the init --force pipeline\n got: %q\nwant: %q", opencodeGot, opencodeWant)
		}

		claudeGot, err := os.ReadFile(claudePath)
		if err != nil {
			t.Fatal(err)
		}
		claudeWant, err := subagent.CompileClaudeCodeWithOptions(def, subagent.CompileOptions{FallbackModel: "claude-custom/z"})
		if err != nil {
			t.Fatal(err)
		}
		if string(claudeGot) != string(claudeWant) {
			t.Errorf("claude model must survive the init --force pipeline\n got: %q\nwant: %q", claudeGot, claudeWant)
		}

		for _, hint := range []string{
			"  info: preserved local model \"custom-model/x\" for .opencode/agent/flowforge-analyst.md (set agents.models in .flowforge/config.yaml to pin explicitly)\n",
			"  info: preserved local model \"claude-custom/z\" for .claude/agents/flowforge-analyst.md (set agents.models in .flowforge/config.yaml to pin explicitly)\n",
		} {
			if !strings.Contains(stderr, hint) {
				t.Errorf("stderr must carry the preserved hint\ngot:  %q\nwant: %q", stderr, hint)
			}
		}
	})
}

func TestAgentStatusTreatsPreservedModelAsCurrent(t *testing.T) {
	t.Run("preserved model reports current", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode", "claude"}
		if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
			t.Fatal(err)
		}

		// Hand-edit local models into both frontmatter hosts, then
		// redeploy so the deployed files carry preserve-merged content.
		def := findDiscoveredDefinition(t, projectRoot, "flowforge-analyst")
		opencodePath := filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md")
		opencodeData, err := os.ReadFile(opencodePath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(opencodePath, []byte(strings.Replace(string(opencodeData), "---\n", "---\nmodel: custom-model/x\n", 1)), 0644); err != nil {
			t.Fatal(err)
		}
		claudePath := filepath.Join(projectRoot, ".claude", "agents", "flowforge-analyst.md")
		claudeData, err := os.ReadFile(claudePath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(claudePath, []byte(strings.Replace(string(claudeData), "model: "+def.ModelProfile.ClaudeModel(), "model: claude-custom/z", 1)), 0644); err != nil {
			t.Fatal(err)
		}

		var deployErr error
		captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}

		result, err := computeSubagentStatus(projectRoot, cfg)
		if err != nil {
			t.Fatalf("computeSubagentStatus: %v", err)
		}
		if !result.Current {
			t.Errorf("preserved-model deploy must report current, got entries: %+v", result.Entries)
		}
		for _, entry := range result.Entries {
			if entry.State != string(managedAssetCurrent) {
				t.Errorf("expected %s to be current, got %s", entry.Target, entry.State)
			}
		}
	})

	t.Run("config pinned model still drifts stale file residue", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		cfg.Agents.Models = map[string]string{"tool-capable": "pinned-by-config/y"}
		if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
			t.Fatal(err)
		}

		// Replace the pinned model with a stale local residue value.
		path := filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		handEdited := strings.Replace(string(data), "model: pinned-by-config/y", "model: stale-local/x", 1)
		if handEdited == string(data) {
			t.Fatal("expected the deployed file to carry the pinned model")
		}
		if err := os.WriteFile(path, []byte(handEdited), 0644); err != nil {
			t.Fatal(err)
		}

		result, err := computeSubagentStatus(projectRoot, cfg)
		if err != nil {
			t.Fatalf("computeSubagentStatus: %v", err)
		}
		if result.Current {
			t.Error("config-pinned expectation must drift a stale file residue")
		}
		foundDrifted := false
		for _, entry := range result.Entries {
			if entry.Target == path && entry.State == string(managedAssetDrifted) {
				foundDrifted = true
			}
		}
		if !foundDrifted {
			t.Errorf("expected drifted entry for %s, got entries: %+v", path, result.Entries)
		}

		// The explicit channel wins: a redeploy writes the config value
		// and status returns to current.
		if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
			t.Fatal(err)
		}
		result, err = computeSubagentStatus(projectRoot, cfg)
		if err != nil {
			t.Fatalf("computeSubagentStatus after redeploy: %v", err)
		}
		if !result.Current {
			t.Errorf("redeploy must restore current, got entries: %+v", result.Entries)
		}
	})

	t.Run("hand-edited body still drifts", func(t *testing.T) {
		projectRoot := t.TempDir()
		if err := initializeTestProject(projectRoot); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.Load(projectRoot)
		if err != nil {
			t.Fatal(err)
		}
		cfg.Agents.Hosts = []string{"opencode"}
		if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
			t.Fatal(err)
		}

		// e2e: hand-edit a model, redeploy (preserve-merge) → current.
		path := filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(strings.Replace(string(data), "---\n", "---\nmodel: custom-model/x\n", 1)), 0644); err != nil {
			t.Fatal(err)
		}
		var deployErr error
		captureStderr(t, func() {
			_, deployErr = deploySubagents(projectRoot, cfg, "")
		})
		if deployErr != nil {
			t.Fatal(deployErr)
		}
		result, err := computeSubagentStatus(projectRoot, cfg)
		if err != nil {
			t.Fatalf("computeSubagentStatus: %v", err)
		}
		if !result.Current {
			t.Fatalf("preserved-model deploy must report current before body drift, got entries: %+v", result.Entries)
		}

		// Then hand-edit the body (outside the model line): still drifted.
		bodyData, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(string(bodyData)+"\nlocal body tweak\n"), 0644); err != nil {
			t.Fatal(err)
		}
		result, err = computeSubagentStatus(projectRoot, cfg)
		if err != nil {
			t.Fatalf("computeSubagentStatus after body drift: %v", err)
		}
		if result.Current {
			t.Error("body hand-edit must still drift")
		}
		foundDrifted := false
		for _, entry := range result.Entries {
			if entry.Target == path && entry.State == string(managedAssetDrifted) {
				foundDrifted = true
			}
		}
		if !foundDrifted {
			t.Errorf("expected drifted entry for %s, got entries: %+v", path, result.Entries)
		}
	})
}

func TestOpenCodeQuestionDeny(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode"}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatal(err)
	}
	impl, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-implementer.md"))
	if err != nil {
		t.Fatal(err)
	}
	analyst, err := os.ReadFile(filepath.Join(projectRoot, ".opencode", "agent", "flowforge-analyst.md"))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("implementer carries question deny", func(t *testing.T) {
		if !strings.Contains(string(impl), "question: deny") {
			t.Error("implementer must carry question: deny in frontmatter")
		}
	})

	t.Run("non-implementer subagents do not", func(t *testing.T) {
		if strings.Contains(string(analyst), "question") {
			t.Error("analyst must not carry question deny (implementer-scoped only)")
		}
	})

	t.Run("edit deny and question deny coexist", func(t *testing.T) {
		// Default test guard is on: implementer must carry both edit globs
		// and question deny under one permission block.
		if !strings.Contains(string(impl), "edit:") || !strings.Contains(string(impl), "question: deny") {
			t.Error("implementer with default test guard must carry both edit deny globs and question: deny")
		}
		// Frontmatter must still be valid YAML.
		content := string(impl)
		start := strings.Index(content, "---\n")
		end := strings.Index(content[start+4:], "---\n")
		if start == -1 || end == -1 {
			t.Fatal("frontmatter delimiters missing")
		}
		var parsed map[string]interface{}
		if err := yaml.Unmarshal([]byte(content[start+4:start+4+end]), &parsed); err != nil {
			t.Errorf("frontmatter invalid YAML after adding question deny: %v", err)
		}
	})
}

// TestAgentRulesDescribePerMachineDeployArtifacts anchors the per-machine
// deploy artifact convention in both the template source (assets/AGENTS.md)
// and the bootstrapped repository root AGENTS.md: deploy artifacts and
// .flowforge/config.yaml are per-machine files that init records in
// .gitignore automatically, and user-authored agents meant for the
// repository opt in with `git add -f`.
func TestAgentRulesDescribePerMachineDeployArtifacts(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	for _, path := range []string{
		filepath.Join(repoRoot, "assets", "AGENTS.md"),
		filepath.Join(repoRoot, "AGENTS.md"),
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s not found: %v", path, err)
		}
		body := string(data)
		for _, needle := range []string{
			"## Per-machine deploy artifacts",
			"per-machine",
			".flowforge/config.yaml",
			".gitignore",
			"git add -f",
		} {
			if !strings.Contains(body, needle) {
				t.Errorf("%s missing per-machine deploy artifact anchor %q", path, needle)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Ticket 01: agents.models_by_host — six-level precedence chain + format-level
// validation (design d-model-channels / d-config-validation).
// ---------------------------------------------------------------------------

// readDeployedModel returns the frontmatter model of a deployed agent file
// ("" when the field is omitted).
func readDeployedModel(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return frontmatterModel(data)
}

// writeLocalModelFile pre-seeds a deployed artifact carrying a locally edited
// frontmatter model — the preserve-merge input.
func writeLocalModelFile(t *testing.T, path, model string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("creating dir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte("---\nmodel: "+model+"\n---\nlocally edited body\n"), 0644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

// TestResolveModelPrecedence drives the six-level model chain: each level
// must beat the level below it (per-host name > per-host profile > global
// name > global profile > preserve-merge backfill > host default).
func TestResolveModelPrecedence(t *testing.T) {
	implPath := func(root string) string {
		return filepath.Join(root, ".opencode", "agent", "flowforge-implementer.md")
	}
	rows := []struct {
		name          string
		hostName      string // 1: models_by_host.<host>[agent-name]
		hostProfile   string // 2: models_by_host.<host>[profile-key]
		globalName    string // 3: models_by_name[agent-name]
		globalProfile string // 4: models[profile-key]
		localModel    string // 5: preserve-merge backfill (pre-deployed file)
		want          string // expected frontmatter model ("" = omitted)
	}{
		{name: "host name key beats host profile key", hostName: "prov/a", hostProfile: "prov/b", globalName: "prov/c", globalProfile: "prov/d", localModel: "prov/e", want: "prov/a"},
		{name: "host profile key beats global name key", hostProfile: "prov/b", globalName: "prov/c", globalProfile: "prov/d", localModel: "prov/e", want: "prov/b"},
		{name: "global name key beats global profile key", globalName: "prov/c", globalProfile: "prov/d", localModel: "prov/e", want: "prov/c"},
		{name: "global profile key beats preserve-merge backfill", globalProfile: "prov/d", localModel: "prov/e", want: "prov/d"},
		{name: "preserve-merge backfill beats host default", localModel: "prov/e", want: "prov/e"},
		{name: "host default: model field omitted", want: ""},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			root := t.TempDir()
			if err := initializeTestProject(root); err != nil {
				t.Fatal(err)
			}
			cfg, err := config.Load(root)
			if err != nil {
				t.Fatal(err)
			}
			cfg.Agents.Hosts = []string{"opencode"}
			if row.hostName != "" || row.hostProfile != "" {
				inner := map[string]string{}
				if row.hostName != "" {
					inner["flowforge-implementer"] = row.hostName
				}
				if row.hostProfile != "" {
					inner["tool-capable"] = row.hostProfile
				}
				cfg.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": inner}
			}
			if row.globalName != "" {
				cfg.Agents.ModelOverrides = map[string]string{"flowforge-implementer": row.globalName}
			}
			if row.globalProfile != "" {
				cfg.Agents.Models = map[string]string{"tool-capable": row.globalProfile}
			}
			if row.localModel != "" {
				writeLocalModelFile(t, implPath(root), row.localModel)
			}
			if _, err := deploySubagents(root, cfg, ""); err != nil {
				t.Fatal(err)
			}
			if got := readDeployedModel(t, implPath(root)); got != row.want {
				t.Errorf("model = %q, want %q", got, row.want)
			}
		})
	}
}

// TestModelsByHostPerHostIsolation checks that each enabled host reads only
// its own models_by_host section (name and profile keys).
func TestModelsByHostPerHostIsolation(t *testing.T) {
	root := t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode", "claude"}
	cfg.Agents.ModelHostOverrides = map[string]map[string]string{
		"opencode": {"flowforge-implementer": "prov/opencode-model"},
		"claude":   {"flowforge-implementer": "claude-pinned"},
	}
	if _, err := deploySubagents(root, cfg, ""); err != nil {
		t.Fatal(err)
	}
	if got := readDeployedModel(t, filepath.Join(root, ".opencode", "agent", "flowforge-implementer.md")); got != "prov/opencode-model" {
		t.Errorf("opencode model = %q, want %q", got, "prov/opencode-model")
	}
	if got := readDeployedModel(t, filepath.Join(root, ".claude", "agents", "flowforge-implementer.md")); got != "claude-pinned" {
		t.Errorf("claude model = %q, want %q", got, "claude-pinned")
	}

	// A per-host profile key pins every agent of that profile on that host
	// (claude default for tool-capable is "sonnet"; the pin must win).
	root2 := t.TempDir()
	if err := initializeTestProject(root2); err != nil {
		t.Fatal(err)
	}
	cfg2, err := config.Load(root2)
	if err != nil {
		t.Fatal(err)
	}
	cfg2.Agents.Hosts = []string{"claude"}
	cfg2.Agents.ModelHostOverrides = map[string]map[string]string{
		"claude": {"tool-capable": "haiku"},
	}
	if _, err := deploySubagents(root2, cfg2, ""); err != nil {
		t.Fatal(err)
	}
	if got := readDeployedModel(t, filepath.Join(root2, ".claude", "agents", "flowforge-implementer.md")); got != "haiku" {
		t.Errorf("claude profile-pinned model = %q, want %q", got, "haiku")
	}
}

// TestValidateModelConfig drives the d-config-validation rules one by one.
func TestValidateModelConfig(t *testing.T) {
	root := t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	defs, err := discoverSubagentSources(root)
	if err != nil {
		t.Fatal(err)
	}
	hosts := func(keys ...string) []hostTarget {
		out := make([]hostTarget, 0, len(keys))
		for _, k := range keys {
			out = append(out, hostTarget{key: k})
		}
		return out
	}
	rows := []struct {
		name    string
		mutate  func(*config.Config)
		hosts   []hostTarget
		wantErr string
	}{
		{
			name: "unknown outer host key",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"bogus": {"flowforge-implementer": "p/m"}}
			},
			hosts:   hosts("opencode"),
			wantErr: `agents.models_by_host: unknown host "bogus"`,
		},
		{
			name: "codex section is a config error",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"codex": {"flowforge-implementer": "p/m"}}
			},
			hosts:   hosts("opencode"),
			wantErr: `agents.models_by_host.codex`,
		},
		{
			name: "unknown inner key on a per-host section",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": {"no-such-agent": "p/m"}}
			},
			hosts:   hosts("opencode"),
			wantErr: `agents.models_by_host.opencode: unknown key "no-such-agent"`,
		},
		{
			name: "unknown models_by_name key",
			mutate: func(c *config.Config) {
				c.Agents.ModelOverrides = map[string]string{"no-such-agent": "p/m"}
			},
			hosts:   hosts("opencode"),
			wantErr: `agents.models_by_name: unknown agent "no-such-agent"`,
		},
		{
			name: "profile key is rejected in models_by_name",
			mutate: func(c *config.Config) {
				c.Agents.ModelOverrides = map[string]string{"tool-capable": "p/m"}
			},
			hosts:   hosts("opencode", "claude"),
			wantErr: `agents.models_by_name: unknown agent "tool-capable"`,
		},
		{
			name: "unknown models profile key",
			mutate: func(c *config.Config) {
				c.Agents.Models = map[string]string{"sonnet": "p/m"}
			},
			hosts:   hosts("opencode"),
			wantErr: `agents.models: unknown profile key "sonnet"`,
		},
		{
			name: "opencode value without provider prefix",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": {"flowforge-implementer": "sonnet"}}
			},
			hosts:   hosts("opencode"),
			wantErr: "provider/model",
		},
		{
			name: "opencode value with empty model side",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": {"flowforge-implementer": "prov/"}}
			},
			hosts:   hosts("opencode"),
			wantErr: "provider/model",
		},
		{
			name: "opencode value with two slashes",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": {"flowforge-implementer": "a/b/c"}}
			},
			hosts:   hosts("opencode"),
			wantErr: "provider/model",
		},
		{
			name: "value with inner whitespace",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": {"flowforge-implementer": "pro vider/model"}}
			},
			hosts:   hosts("opencode"),
			wantErr: "whitespace",
		},
		{
			name: "empty value",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": {"flowforge-implementer": ""}}
			},
			hosts:   hosts("opencode"),
			wantErr: "non-empty",
		},
		{
			name: "claude single token is legal",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"claude": {"flowforge-implementer": "sonnet"}}
			},
			hosts: hosts("claude"),
		},
		{
			name: "claude value with whitespace is illegal",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"claude": {"flowforge-implementer": "my model"}}
			},
			hosts:   hosts("claude"),
			wantErr: "whitespace",
		},
		{
			name: "global value failing one enabled host must suggest models_by_host",
			mutate: func(c *config.Config) {
				c.Agents.Models = map[string]string{"tool-capable": "sonnet"}
			},
			hosts:   hosts("opencode", "claude"),
			wantErr: "agents.models_by_host",
		},
		{
			name: "same global value is fine when only claude is enabled",
			mutate: func(c *config.Config) {
				c.Agents.Models = map[string]string{"tool-capable": "sonnet"}
			},
			hosts: hosts("claude"),
		},
		{
			name: "disabled host values do not block deploy",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"pi": {"flowforge-implementer": "no-slash"}}
			},
			hosts: hosts("claude"),
		},
		{
			name: "enabled pi requires provider/model",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"pi": {"flowforge-implementer": "sonnet"}}
			},
			hosts:   hosts("pi"),
			wantErr: "provider/model",
		},
		{
			name: "profile key is legal in a per-host section",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"opencode": {"tool-capable": "p/m"}}
			},
			hosts: hosts("opencode"),
		},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			cfg := &config.Config{}
			row.mutate(cfg)
			err := validateModelConfig(cfg, defs, row.hosts)
			if row.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", row.wantErr)
			}
			if !strings.Contains(err.Error(), row.wantErr) {
				t.Errorf("error %q does not contain %q", err.Error(), row.wantErr)
			}
		})
	}
}

// TestModelsByHostValidationFailFast: a bad model config fails deploy before
// any directory or artifact write, and status reports the identical error.
func TestModelsByHostValidationFailFast(t *testing.T) {
	rows := []struct {
		name   string
		mutate func(*config.Config)
		wantIn string
	}{
		{
			name: "unknown host",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"bogus": {"flowforge-implementer": "p/m"}}
			},
			wantIn: `agents.models_by_host: unknown host "bogus"`,
		},
		{
			name: "codex section",
			mutate: func(c *config.Config) {
				c.Agents.ModelHostOverrides = map[string]map[string]string{"codex": {"flowforge-implementer": "p/m"}}
			},
			wantIn: "codex",
		},
		{
			name: "global value invalid for an enabled host",
			mutate: func(c *config.Config) {
				c.Agents.Models = map[string]string{"tool-capable": "sonnet"}
				c.Agents.Hosts = []string{"opencode", "claude"}
			},
			wantIn: "agents.models_by_host",
		},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			root := t.TempDir()
			if err := initializeTestProject(root); err != nil {
				t.Fatal(err)
			}
			cfg, err := config.Load(root)
			if err != nil {
				t.Fatal(err)
			}
			row.mutate(cfg)

			_, deployErr := deploySubagents(root, cfg, "")
			if deployErr == nil {
				t.Fatal("expected deploy error, got nil")
			}
			if !strings.Contains(deployErr.Error(), row.wantIn) {
				t.Errorf("deploy error %q does not contain %q", deployErr.Error(), row.wantIn)
			}
			if _, err := os.Stat(filepath.Join(root, ".opencode", "agent")); !os.IsNotExist(err) {
				t.Error("failed deploy must not create host directories or artifacts")
			}

			_, statusErr := computeSubagentStatus(root, cfg)
			if statusErr == nil {
				t.Fatal("expected status error, got nil")
			}
			if statusErr.Error() != deployErr.Error() {
				t.Errorf("status and deploy must report the identical config error\n deploy: %q\n status: %q", deployErr.Error(), statusErr.Error())
			}
		})
	}
}

// TestStatusUsesSameCompileOptionsModelsByHost extends the pinned
// deploy/status compile-parity convention to the per-host model layer.
func TestStatusUsesSameCompileOptionsModelsByHost(t *testing.T) {
	root := t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.Hosts = []string{"opencode", "claude"}
	cfg.Agents.ModelHostOverrides = map[string]map[string]string{
		"opencode": {"flowforge-implementer": "prov/oc"},
		"claude":   {"flowforge-implementer": "opus-pinned"},
	}
	if _, err := deploySubagents(root, cfg, ""); err != nil {
		t.Fatal(err)
	}
	result, err := computeSubagentStatus(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Current {
		t.Error("status must use the same per-host compile options as deploy (no false drift)")
	}
	if got := readDeployedModel(t, filepath.Join(root, ".opencode", "agent", "flowforge-implementer.md")); got != "prov/oc" {
		t.Errorf("opencode model = %q, want %q", got, "prov/oc")
	}
	if got := readDeployedModel(t, filepath.Join(root, ".claude", "agents", "flowforge-implementer.md")); got != "opus-pinned" {
		t.Errorf("claude model = %q, want %q", got, "opus-pinned")
	}
}
