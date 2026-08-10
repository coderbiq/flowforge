package orchestration

import (
	"strings"
	"testing"
)

func TestRenderOpenCodeCreatesOnlyEnabledWorkers(t *testing.T) {
	files, err := RenderOpenCode(DefaultPolicy())
	if err != nil {
		t.Fatalf("rendering adapter: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected analyst and executor, got %d", len(files))
	}
	if _, ok := files["flowforge-reviewer.md"]; ok {
		t.Fatal("disabled reviewer must not render")
	}
	for name, want := range map[string]string{
		"flowforge-design-analyst.md": "edit: deny",
		"flowforge-executor.md":       "edit: allow",
	} {
		if !strings.Contains(string(files[name]), want) {
			t.Fatalf("%s missing %q:\n%s", name, want, files[name])
		}
	}
}

func TestRenderOpenCodeIncludesEnabledReviewer(t *testing.T) {
	policy := DefaultPolicy()
	policy.Roles[3].Enabled = true
	files, err := RenderOpenCode(policy)
	if err != nil {
		t.Fatal(err)
	}
	content, ok := files["flowforge-reviewer.md"]
	if !ok {
		t.Fatal("enabled reviewer missing")
	}
	if !strings.Contains(string(content), "Default Skill: flowforge-review") {
		t.Fatalf("reviewer missing skill:\n%s", content)
	}
}
