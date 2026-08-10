package orchestration

import (
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func TestRenderOpenCodeIncludesCompleteRoleContracts(t *testing.T) {
	files, err := RenderOpenCode(DefaultPolicy())
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"flowforge-coordinator.md", "flowforge-design-analyst.md", "flowforge-executor.md"} {
		if _, ok := files[name]; !ok {
			t.Fatalf("missing %s", name)
		}
	}
	checks := map[string][]string{
		"flowforge-coordinator.md":    {"mode: primary", "\"*\": deny", "flowforge-executor: allow", "Routing Contract", "context preflight", "context risk-review"},
		"flowforge-design-analyst.md": {"edit: allow", "flowforge-design", "never modify product code", "Journal"},
		"flowforge-executor.md":       {"edit: allow", "flowforge-implement", "design_gap", "verification_failed"},
	}
	for name, wants := range checks {
		for _, want := range wants {
			if !strings.Contains(string(files[name]), want) {
				t.Fatalf("%s missing %q", name, want)
			}
		}
	}
	for name, content := range files {
		parts := strings.SplitN(string(content), "---", 3)
		if len(parts) != 3 {
			t.Fatalf("%s has invalid frontmatter delimiters", name)
		}
		var frontmatter map[string]any
		if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
			t.Fatalf("%s has invalid YAML: %v", name, err)
		}
		if strings.Contains(string(content), "provider/") || strings.Contains(string(content), "model:") {
			t.Fatalf("%s contains provider-specific model selection", name)
		}
	}
}

func TestRenderOpenCodeOmitsDisabledReviewerFromCoordinator(t *testing.T) {
	policy := DefaultPolicy()
	policy.Roles[3].Enabled = false
	files, err := RenderOpenCode(policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := files["flowforge-reviewer.md"]; ok {
		t.Fatal("disabled reviewer rendered")
	}
	if strings.Contains(string(files["flowforge-coordinator.md"]), "flowforge-reviewer: allow") {
		t.Fatal("coordinator allows disabled reviewer")
	}
}

func TestRenderCodexIncludesCompleteWorkerContracts(t *testing.T) {
	files, err := RenderCodex(DefaultPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Fatalf("expected coordinator and 2 workers, got %d", len(files))
	}
	for name, content := range files {
		for _, want := range []string{"developer_instructions = \"\"\"", "Shared Workflow", "Default Skill:"} {
			if !strings.Contains(string(content), want) {
				t.Fatalf("%s missing %q", name, want)
			}
		}
		var document map[string]any
		if err := toml.Unmarshal(content, &document); err != nil {
			t.Fatalf("%s has invalid TOML: %v", name, err)
		}
		if _, ok := document["model"]; ok {
			t.Fatalf("%s contains provider-specific model selection", name)
		}
	}
}
