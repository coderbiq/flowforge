package subagent

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func testAssetsDir(t *testing.T) string {
	t.Helper()
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	return filepath.Join(repoRoot, "assets", "subagents")
}

func TestParseDirReturnsSixDefinitions(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}
	if len(definitions) != 6 {
		t.Fatalf("expected 6 definitions, got %d", len(definitions))
	}
	expectedNames := []string{
		"flowforge-analyst",
		"flowforge-architect",
		"flowforge-implementer",
		"flowforge-investigator",
		"flowforge-planner",
		"flowforge-reviewer",
	}
	for i, expected := range expectedNames {
		if definitions[i].Name != expected {
			t.Errorf("definition[%d]: expected name %q, got %q", i, expected, definitions[i].Name)
		}
	}
}

func TestParseRejectsNameMismatch(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test-agent.md")
	content := `---
flowforge_agent:
  name: wrong-name
  description: test
  model_profile: tool-capable
  default_skill: flowforge-align
---
body`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for name mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "does not match filename") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCompileClaudeCodeIncludesSkillsField(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, def := range definitions {
		compiled, err := CompileClaudeCode(def)
		if err != nil {
			t.Errorf("CompileClaudeCode(%s): %v", def.Name, err)
			continue
		}
		content := string(compiled)
		if !strings.Contains(content, "skills:") {
			t.Errorf("CompileClaudeCode(%s): missing 'skills:' field", def.Name)
		}
		if !strings.Contains(content, def.DefaultSkill) {
			t.Errorf("CompileClaudeCode(%s): missing default skill %q", def.Name, def.DefaultSkill)
		}
		// Verify frontmatter is valid YAML
		fm, _, ok := splitFrontmatter(compiled)
		if !ok {
			t.Errorf("CompileClaudeCode(%s): frontmatter delimiters missing", def.Name)
			continue
		}
		var parsed map[string]interface{}
		if err := yaml.Unmarshal(fm, &parsed); err != nil {
			t.Errorf("CompileClaudeCode(%s): frontmatter invalid YAML: %v", def.Name, err)
		}
	}
}

func TestCompileOpenCodeOmitsModelField(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, def := range definitions {
		compiled, err := CompileOpenCode(def)
		if err != nil {
			t.Errorf("CompileOpenCode(%s): %v", def.Name, err)
			continue
		}
		content := string(compiled)
		if strings.Contains(content, "model:") {
			t.Errorf("CompileOpenCode(%s): unexpectedly contains 'model:' field", def.Name)
		}
		if !strings.Contains(content, "mode: subagent") {
			t.Errorf("CompileOpenCode(%s): missing 'mode: subagent'", def.Name)
		}
		// Verify frontmatter is valid YAML
		fm, _, ok := splitFrontmatter(compiled)
		if !ok {
			t.Errorf("CompileOpenCode(%s): frontmatter delimiters missing", def.Name)
			continue
		}
		var parsed map[string]interface{}
		if err := yaml.Unmarshal(fm, &parsed); err != nil {
			t.Errorf("CompileOpenCode(%s): frontmatter invalid YAML: %v", def.Name, err)
		}
	}
}

func TestCompileCodexReplacesSkillInvocation(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, def := range definitions {
		compiled, err := CompileCodex(def)
		if err != nil {
			t.Errorf("CompileCodex(%s): %v", def.Name, err)
			continue
		}
		content := string(compiled)
		if strings.Contains(content, "invoke the Skill tool") {
			t.Errorf("CompileCodex(%s): still contains 'invoke the Skill tool' (should be replaced)", def.Name)
		}
		if !strings.Contains(content, "developer_instructions") {
			t.Errorf("CompileCodex(%s): missing 'developer_instructions'", def.Name)
		}
		expectedDirective := fmt.Sprintf("Read and follow `.agents/skills/%s/SKILL.md`", def.DefaultSkill)
		if !strings.Contains(content, expectedDirective) {
			t.Errorf("CompileCodex(%s): missing expected file read directive %q", def.Name, expectedDirective)
		}
	}
}

func TestCompileIsIdempotent(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, def := range definitions {
		// Test Claude Code idempotency
		cc1, err1 := CompileClaudeCode(def)
		cc2, err2 := CompileClaudeCode(def)
		if err1 != nil || err2 != nil {
			t.Errorf("CompileClaudeCode(%s) errors: %v, %v", def.Name, err1, err2)
		} else if !bytes.Equal(cc1, cc2) {
			t.Errorf("CompileClaudeCode(%s) not idempotent", def.Name)
		}

		// Test OpenCode idempotency
		oc1, err1 := CompileOpenCode(def)
		oc2, err2 := CompileOpenCode(def)
		if err1 != nil || err2 != nil {
			t.Errorf("CompileOpenCode(%s) errors: %v, %v", def.Name, err1, err2)
		} else if !bytes.Equal(oc1, oc2) {
			t.Errorf("CompileOpenCode(%s) not idempotent", def.Name)
		}

		// Test Codex idempotency
		cx1, err1 := CompileCodex(def)
		cx2, err2 := CompileCodex(def)
		if err1 != nil || err2 != nil {
			t.Errorf("CompileCodex(%s) errors: %v, %v", def.Name, err1, err2)
		} else if !bytes.Equal(cx1, cx2) {
			t.Errorf("CompileCodex(%s) not idempotent", def.Name)
		}
	}
}
