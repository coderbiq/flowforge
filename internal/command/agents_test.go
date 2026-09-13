package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/config"
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

	if len(deployed) != 6 {
		t.Fatalf("expected 6 deployed subagents, got %d", len(deployed))
	}

	expectedRoles := []string{
		"flowforge-analyst",
		"flowforge-architect",
		"flowforge-implementer",
		"flowforge-investigator",
		"flowforge-planner",
		"flowforge-reviewer",
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

	// Should deploy 4 agents (6 - 2 disabled)
	if len(deployed) != 4 {
		t.Fatalf("expected 4 deployed subagents, got %d", len(deployed))
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

	if len(removedPaths) != 3 {
		t.Errorf("expected 3 removed paths, got %d", len(removedPaths))
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

	// Should remove 3 host files + 1 source file = 4 total
	if len(removedPaths) != 4 {
		t.Errorf("expected 4 removed paths (3 hosts + source), got %d", len(removedPaths))
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

	if len(deployed) != 6 {
		t.Fatalf("expected 6 subagents deployed, got %d", len(deployed))
	}

	// Verify all 3 host directories have 6 files each
	for _, dir := range []string{".claude/agents", ".opencode/agent", ".codex/agents"} {
		entries, err := os.ReadDir(filepath.Join(projectRoot, dir))
		if err != nil {
			t.Fatalf("reading %s: %v", dir, err)
		}
		if len(entries) != 6 {
			t.Errorf("expected 6 files in %s, got %d", dir, len(entries))
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

	if len(deployed) != 6 {
		t.Fatalf("expected 6 subagents, got %d", len(deployed))
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

	// Should deploy 5 (6 - 1 disabled)
	if len(deployed) != 5 {
		t.Fatalf("expected 5 subagents, got %d", len(deployed))
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
	if len(deployed) != 6 {
		t.Fatalf("expected 6 deployed subagents, got %d", len(deployed))
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
	if strings.Contains(string(impl), "permission:") {
		t.Error("DisableTestGuard must remove the permission block")
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
