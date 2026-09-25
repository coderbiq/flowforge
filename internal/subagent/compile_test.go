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
	if len(definitions) != 9 {
		t.Fatalf("expected 9 definitions, got %d", len(definitions))
	}
	expectedNames := []string{
		"flowforge-analyst",
		"flowforge-architect",
		"flowforge-batch-analyst",
		"flowforge-executor",
		"flowforge-implementer",
		"flowforge-investigator",
		"flowforge-planner",
		"flowforge-reviewer",
		"flowforge-scribe",
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

func TestCompileOpenCodeModelFallbackPriority(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, def := range definitions {
		cases := []struct {
			name      string
			opts      CompileOptions
			wantModel bool
			want      string
		}{
			{"config model beats fallback", CompileOptions{Model: "A", FallbackModel: "B"}, true, "A"},
			{"config model alone", CompileOptions{Model: "A"}, true, "A"},
			{"fallback fills absent model", CompileOptions{FallbackModel: "B"}, true, "B"},
			{"zero options omit model", CompileOptions{}, false, ""},
		}
		for _, tc := range cases {
			compiled, err := CompileOpenCodeWithOptions(def, tc.opts)
			if err != nil {
				t.Errorf("%s/%s: %v", def.Name, tc.name, err)
				continue
			}
			fm, _, ok := splitFrontmatter(compiled)
			if !ok {
				t.Errorf("%s/%s: frontmatter delimiters missing", def.Name, tc.name)
				continue
			}
			var parsed map[string]interface{}
			if err := yaml.Unmarshal(fm, &parsed); err != nil {
				t.Errorf("%s/%s: frontmatter invalid YAML: %v", def.Name, tc.name, err)
				continue
			}
			got, present := parsed["model"]
			if !tc.wantModel {
				if present {
					t.Errorf("%s/%s: unexpected model field %v", def.Name, tc.name, got)
				}
				continue
			}
			if s, ok := got.(string); !ok || s != tc.want {
				t.Errorf("%s/%s: model = %v (want %q)", def.Name, tc.name, got, tc.want)
			}
		}

		// Zero options must stay byte-identical to the zero-arg delegate.
		base, err := CompileOpenCode(def)
		if err != nil {
			t.Fatalf("CompileOpenCode(%s): %v", def.Name, err)
		}
		zero, err := CompileOpenCodeWithOptions(def, CompileOptions{})
		if err != nil {
			t.Fatalf("CompileOpenCodeWithOptions(%s, zero): %v", def.Name, err)
		}
		if !bytes.Equal(zero, base) {
			t.Errorf("CompileOpenCodeWithOptions(%s, zero) must equal CompileOpenCode byte for byte", def.Name)
		}
	}
}

func TestCompileClaudeCodeModelPriority(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, def := range definitions {
		profileDefault := def.ModelProfile.ClaudeModel()
		cases := []struct {
			name string
			opts CompileOptions
			want string
		}{
			{"zero options keep profile default", CompileOptions{}, profileDefault},
			{"config model beats fallback", CompileOptions{Model: "A", FallbackModel: "B"}, "A"},
			{"fallback beats profile default", CompileOptions{FallbackModel: "B"}, "B"},
			{"config model alone", CompileOptions{Model: "A"}, "A"},
		}
		for _, tc := range cases {
			compiled, err := CompileClaudeCodeWithOptions(def, tc.opts)
			if err != nil {
				t.Errorf("%s/%s: %v", def.Name, tc.name, err)
				continue
			}
			fm, _, ok := splitFrontmatter(compiled)
			if !ok {
				t.Errorf("%s/%s: frontmatter delimiters missing", def.Name, tc.name)
				continue
			}
			var parsed map[string]interface{}
			if err := yaml.Unmarshal(fm, &parsed); err != nil {
				t.Errorf("%s/%s: frontmatter invalid YAML: %v", def.Name, tc.name, err)
				continue
			}
			// Claude always serializes the model field (no omitempty).
			if s, ok := parsed["model"].(string); !ok || s != tc.want {
				t.Errorf("%s/%s: model = %v (want %q)", def.Name, tc.name, parsed["model"], tc.want)
			}
		}

		// Zero options must stay byte-identical to the zero-arg delegate.
		base, err := CompileClaudeCode(def)
		if err != nil {
			t.Fatalf("CompileClaudeCode(%s): %v", def.Name, err)
		}
		zero, err := CompileClaudeCodeWithOptions(def, CompileOptions{})
		if err != nil {
			t.Fatalf("CompileClaudeCodeWithOptions(%s, zero): %v", def.Name, err)
		}
		if !bytes.Equal(zero, base) {
			t.Errorf("CompileClaudeCodeWithOptions(%s, zero) must equal CompileClaudeCode byte for byte", def.Name)
		}
	}
}

func TestCompileOpenCodeStepsBudget(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}
	for _, def := range definitions {
		// Nil MaxSteps: no steps field, byte-identical to the zero-options
		// CompileOpenCode delegate (pre-option behavior).
		base, err := CompileOpenCode(def)
		if err != nil {
			t.Fatalf("CompileOpenCode(%s): %v", def.Name, err)
		}
		zero, err := CompileOpenCodeWithOptions(def, CompileOptions{})
		if err != nil {
			t.Fatalf("CompileOpenCodeWithOptions(%s, zero): %v", def.Name, err)
		}
		if !bytes.Equal(zero, base) {
			t.Errorf("CompileOpenCodeWithOptions(%s, zero) must equal CompileOpenCode byte for byte", def.Name)
		}
		fm, _, ok := splitFrontmatter(zero)
		if !ok {
			t.Fatalf("CompileOpenCodeWithOptions(%s, zero): frontmatter delimiters missing", def.Name)
		}
		if strings.Contains(string(fm), "steps:") {
			t.Errorf("CompileOpenCodeWithOptions(%s, zero): unexpected 'steps:' field", def.Name)
		}

		// Set MaxSteps: frontmatter carries the value.
		steps := 321
		withBudget, err := CompileOpenCodeWithOptions(def, CompileOptions{MaxSteps: &steps})
		if err != nil {
			t.Fatalf("CompileOpenCodeWithOptions(%s, steps): %v", def.Name, err)
		}
		fm, _, ok = splitFrontmatter(withBudget)
		if !ok {
			t.Fatalf("CompileOpenCodeWithOptions(%s, steps): frontmatter delimiters missing", def.Name)
		}
		if !strings.Contains(string(fm), "steps: 321") {
			t.Errorf("CompileOpenCodeWithOptions(%s, steps): missing 'steps: 321'", def.Name)
		}
		var parsed map[string]interface{}
		if err := yaml.Unmarshal(fm, &parsed); err != nil {
			t.Errorf("CompileOpenCodeWithOptions(%s, steps): frontmatter invalid YAML: %v", def.Name, err)
		} else if got, ok := parsed["steps"].(int); !ok || got != 321 {
			t.Errorf("CompileOpenCodeWithOptions(%s, steps): parsed steps = %v (want int 321)", def.Name, parsed["steps"])
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

func TestCompilePiFields(t *testing.T) {
	dir := testAssetsDir(t)
	definitions, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, def := range definitions {
		compiled, err := CompilePi(def)
		if err != nil {
			t.Errorf("CompilePi(%s): %v", def.Name, err)
			continue
		}
		fm, body, ok := splitFrontmatter(compiled)
		if !ok {
			t.Errorf("CompilePi(%s): frontmatter delimiters missing", def.Name)
			continue
		}
		var parsed map[string]interface{}
		if err := yaml.Unmarshal(fm, &parsed); err != nil {
			t.Errorf("CompilePi(%s): frontmatter invalid YAML: %v", def.Name, err)
			continue
		}
		if s, ok := parsed["name"].(string); !ok || s != def.Name {
			t.Errorf("CompilePi(%s): name = %v (want %q)", def.Name, parsed["name"], def.Name)
		}
		if s, ok := parsed["description"].(string); !ok || s != def.Description {
			t.Errorf("CompilePi(%s): description mismatch", def.Name)
		}
		// thinking maps the model profile (high-capability → high, rest → medium).
		wantThinking := "medium"
		if def.ModelProfile == ModelProfileHighCapability {
			wantThinking = "high"
		}
		if s, ok := parsed["thinking"].(string); !ok || s != wantThinking {
			t.Errorf("CompilePi(%s): thinking = %v (want %q)", def.Name, parsed["thinking"], wantThinking)
		}
		// PI omits model: the child inherits the parent session model.
		if _, present := parsed["model"]; present {
			t.Errorf("CompilePi(%s): unexpected 'model' field", def.Name)
		}
		// read-only permission compiles to a strict read-tool allowlist.
		tools, hasTools := parsed["tools"].([]interface{})
		if def.Permission == "read-only" {
			if !hasTools || len(tools) != 4 {
				t.Errorf("CompilePi(%s): read-only tools allowlist = %v", def.Name, parsed["tools"])
			} else {
				if s, ok := tools[0].(string); !ok || s != "read" {
					t.Errorf("CompilePi(%s): tools[0] = %v (want read)", def.Name, tools[0])
				}
			}
		} else if hasTools {
			t.Errorf("CompilePi(%s): unexpected tools allowlist %v", def.Name, parsed["tools"])
		}
		// skills: default skill first, then detour skills.
		skills, ok := parsed["skills"].([]interface{})
		if !ok || len(skills) == 0 {
			t.Errorf("CompilePi(%s): missing skills list", def.Name)
			continue
		}
		if s, ok := skills[0].(string); !ok || s != def.DefaultSkill {
			t.Errorf("CompilePi(%s): skills[0] = %v (want default skill %q)", def.Name, skills[0], def.DefaultSkill)
		}
		for _, detour := range def.DetourSkills {
			found := false
			for _, s := range skills {
				if v, _ := s.(string); v == detour {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("CompilePi(%s): detour skill %q missing from skills %v", def.Name, detour, parsed["skills"])
			}
		}
		if v, ok := parsed["inheritSkills"].(bool); !ok || v {
			t.Errorf("CompilePi(%s): inheritSkills = %v (want false)", def.Name, parsed["inheritSkills"])
		}
		// Body is preserved byte for byte.
		if string(body) != def.Body {
			t.Errorf("CompilePi(%s): body not byte-identical to Definition.Body", def.Name)
		}
		// Idempotence and zero-options equivalence.
		second, err := CompilePi(def)
		if err != nil {
			t.Errorf("CompilePi(%s) second call: %v", def.Name, err)
		} else if !bytes.Equal(compiled, second) {
			t.Errorf("CompilePi(%s) not idempotent", def.Name)
		}
		withOpts, err := CompilePiWithOptions(def, CompileOptions{Model: "ignored", FallbackModel: "ignored", EditDeny: []string{"ignored"}, DenyQuestion: true})
		if err != nil {
			t.Errorf("CompilePiWithOptions(%s, inexpressible opts): %v", def.Name, err)
		} else if !bytes.Equal(compiled, withOpts) {
			t.Errorf("CompilePiWithOptions(%s): inexpressible options must not change output", def.Name)
		}
	}
}
