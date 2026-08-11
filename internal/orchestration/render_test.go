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
	for _, name := range []string{"flowforge-coordinator.md", "flowforge-design-analyst.md", "flowforge-investigator.md", "flowforge-executor.md"} {
		if _, ok := files[name]; !ok {
			t.Fatalf("missing %s", name)
		}
	}
	checks := map[string][]string{
		"flowforge-coordinator.md":    {"mode: primary", "\"*\": deny", "flowforge-investigator: allow", "flowforge-executor: allow", "execution-only scheduler", "re-entry condition", "context preflight", "context risk-review"},
		"flowforge-design-analyst.md": {"edit: allow", "flowforge-design", "never modify product code", "evidence acceptance or rejection", "Re-enter"},
		"flowforge-investigator.md":   {"edit: allow", "Investigation Contract", "one writable FIND", "INCONCLUSIVE", "EVIDENCE_CONFLICT", "External research is denied"},
		"flowforge-executor.md":       {"edit: allow", "flowforge-implement", "design_gap", "verification_failed"},
	}
	for name, wants := range map[string][]string{
		"flowforge-design-analyst.md": {"\"flowforge *\": allow", "\"./bin/flowforge *\": allow"},
		"flowforge-investigator.md":   {"\"*\": deny", "\"rg *\": allow", "\"sed *\": allow"},
	} {
		for _, want := range wants {
			if !strings.Contains(string(files[name]), want) {
				t.Fatalf("%s missing permission %q", name, want)
			}
		}
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

func TestEnforcementSummaryReportsSoftAndUnsupportedBoundaries(t *testing.T) {
	for _, host := range []string{"opencode", "codex"} {
		summary := EnforcementSummary(host)
		for _, want := range []string{"soft", "external_sources=unsupported"} {
			if !strings.Contains(summary, want) {
				t.Fatalf("%s enforcement summary missing %q: %s", host, want, summary)
			}
		}
	}
}

func TestRenderOpenCodeOmitsDisabledReviewerFromCoordinator(t *testing.T) {
	policy := DefaultPolicy()
	policy.Roles[4].Enabled = false
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
	if len(files) != 4 {
		t.Fatalf("expected coordinator and 3 workers, got %d", len(files))
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
	for _, name := range []string{"flowforge-design-analyst.toml", "flowforge-investigator.toml", "flowforge-executor.toml"} {
		if !strings.Contains(string(files[name]), `sandbox_mode = "workspace-write"`) {
			t.Fatalf("%s must be workspace-write", name)
		}
	}
}
