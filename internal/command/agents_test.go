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
