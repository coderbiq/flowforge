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
