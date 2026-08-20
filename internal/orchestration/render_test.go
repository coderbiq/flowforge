package orchestration

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/*.golden
var renderGoldens embed.FS

func TestRenderGoldensAndStableInventory(t *testing.T) {
	tests := []struct {
		name   string
		render func(Policy) (RenderOutput, error)
	}{
		{"opencode", RenderOpenCodeOutput},
		{"codex", RenderCodexOutput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			policy := DefaultPolicy()
			first, err := test.render(policy)
			if err != nil {
				t.Fatal(err)
			}
			second, err := test.render(policy)
			if err != nil {
				t.Fatal(err)
			}
			if first.Host != test.name || first.RendererVersion == "" || first.PolicyDigest == "" {
				t.Fatalf("incomplete renderer metadata: %+v", first)
			}
			if fmt.Sprintf("%v", first) != fmt.Sprintf("%v", second) {
				t.Fatal("repeated render inventory differs")
			}
			for i := range first.Files {
				if string(first.Files[i].Content) != string(second.Files[i].Content) {
					t.Fatalf("repeated render differs for %s", first.Files[i].Source)
				}
				if first.Files[i].Host != first.Host || first.Files[i].Type != test.name+"_agent" {
					t.Fatalf("unstable metadata: %+v", first.Files[i])
				}
			}
			golden, err := renderGoldens.ReadFile(filepath.Join("testdata", test.name+".golden"))
			if err != nil {
				t.Fatal(err)
			}
			want := parseGolden(t, string(golden))
			if first.PolicyDigest != want["digest"] || first.RendererVersion != want["version"] {
				t.Fatalf("metadata mismatch: got %s/%s want %s/%s", first.RendererVersion, first.PolicyDigest, want["version"], want["digest"])
			}
			if got := strings.Join(renderedSources(first.Files), ","); got != want["files"] {
				t.Fatalf("files mismatch: got %s want %s", got, want["files"])
			}
		})
	}
}

func parseGolden(t *testing.T, content string) map[string]string {
	t.Helper()
	result := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			t.Fatalf("invalid golden line %q", line)
		}
		result[parts[0]] = parts[1]
	}
	return result
}

func renderedSources(files []RenderedFile) []string {
	result := make([]string, len(files))
	for i, file := range files {
		result[i] = file.Source
	}
	return result
}

func TestRenderCrossHostFormatsAndDigestsDiffer(t *testing.T) {
	opencode, err := RenderOpenCodeOutput(DefaultPolicy())
	if err != nil {
		t.Fatal(err)
	}
	codex, err := RenderCodexOutput(DefaultPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if opencode.PolicyDigest == codex.PolicyDigest {
		t.Fatal("cross-host renderer digests must differ")
	}
	if strings.Contains(string(opencode.Files[0].Content), "sandbox_mode") || strings.Contains(string(codex.Files[0].Content), "permission:") {
		t.Fatal("host-specific fields crossed renderer boundary")
	}
}

func TestRenderRejectsInvalidPolicyWithoutOutput(t *testing.T) {
	policy := DefaultPolicy()
	policy.Roles[1].DefaultSkill = "unknown-skill"
	if output, err := RenderOpenCodeOutput(policy); err == nil || output.Files != nil {
		t.Fatalf("invalid OpenCode policy produced output: %+v, %v", output, err)
	}
	if output, err := RenderCodexOutput(policy); err == nil || output.Files != nil {
		t.Fatalf("invalid Codex policy produced output: %+v, %v", output, err)
	}
}

func TestRenderContentDigestsAreStable(t *testing.T) {
	output, err := RenderOpenCodeOutput(DefaultPolicy())
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range output.Files {
		if file.Source == "" || file.Host == "" || file.Type == "" || len(file.Content) == 0 {
			t.Fatalf("incomplete rendered file: %+v", file)
		}
		hash := sha256.Sum256(file.Content)
		if fmt.Sprintf("%x", hash) == "" {
			t.Fatal("empty content digest")
		}
	}
	if !sort.SliceIsSorted(output.Files, func(i, j int) bool { return output.Files[i].Source < output.Files[j].Source }) {
		t.Fatal("rendered files are not sorted")
	}
}

type RendererGolden struct {
	Host            string
	RendererVersion string
	PolicyDigest    string
	Files           []fixtureBaseline
}

type fixtureBaseline struct {
	Path   string
	Body   []byte
	SHA256 string
}

func TestRendererGoldenFixtureSelfCheck(t *testing.T) {
	for _, test := range []struct {
		name   string
		render func(Policy) (RenderOutput, error)
	}{
		{"opencode", RenderOpenCodeOutput},
		{"codex", RenderCodexOutput},
	} {
		t.Run(test.name, func(t *testing.T) {
			first, err := test.render(DefaultPolicy())
			if err != nil {
				t.Fatal(err)
			}
			second, err := test.render(DefaultPolicy())
			if err != nil {
				t.Fatal(err)
			}
			if first.Host != test.name || first.RendererVersion == "" || first.PolicyDigest == "" || len(first.Files) == 0 {
				t.Fatalf("invalid renderer golden: %+v", first)
			}
			if first.PolicyDigest != second.PolicyDigest || first.RendererVersion != second.RendererVersion || len(first.Files) != len(second.Files) {
				t.Fatal("renderer golden metadata is not stable")
			}
			for i, file := range first.Files {
				if file.Source != second.Files[i].Source || string(file.Content) != string(second.Files[i].Content) {
					t.Fatalf("renderer golden changed for %s", file.Source)
				}
				hash := sha256.Sum256(file.Content)
				baseline := fixtureBaseline{Path: file.Source, Body: append([]byte(nil), file.Content...), SHA256: fmt.Sprintf("%x", hash)}
				if baseline.SHA256 == "" || len(baseline.Body) == 0 {
					t.Fatalf("empty renderer baseline: %+v", baseline)
				}
			}
		})
	}
}

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
		"flowforge-design-analyst.md": {"edit: allow", "flowforge-align", "never modify product code", "evidence acceptance or rejection", "Re-enter"},
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
		t.Fatalf("expected manual coordinator and 3 workers, got %d", len(files))
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
	checks := map[string]string{
		"flowforge-coordinator.toml":    "explicitly requests the FlowForge Coordinator as a manual fallback",
		"flowforge-design-analyst.toml": "Proposal design, FEATURE decomposition",
		"flowforge-investigator.toml":   "one ready registered FlowForge investigation brief",
		"flowforge-executor.toml":       "context preflight returns allow and requires handoff",
	}
	for name, want := range checks {
		if !strings.Contains(string(files[name]), want) {
			t.Fatalf("%s missing activation description %q", name, want)
		}
	}
}
