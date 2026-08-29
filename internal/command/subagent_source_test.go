package command

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

var expectedSubagentNames = []string{
	"flowforge-analyst",
	"flowforge-architect",
	"flowforge-planner",
	"flowforge-implementer",
	"flowforge-reviewer",
	"flowforge-investigator",
}

var validModelProfiles = map[string]bool{
	"high-capability":        true,
	"tool-capable":           true,
	"tool-capable-read-only": true,
}

var requiredSubagentSections = []string{
	"## Identity",
	"## Boundaries",
	"## Workflow Position",
	"## Default Skill",
	"## Result Contract",
}

type subagentFrontmatter struct {
	FlowforgeAgent struct {
		Name         string   `yaml:"name"`
		Description  string   `yaml:"description"`
		ModelProfile string   `yaml:"model_profile"`
		DefaultSkill string   `yaml:"default_skill"`
		DetourSkills []string `yaml:"detour_skills"`
		Permission   string   `yaml:"permission"`
		After        []string `yaml:"after"`
		Before       []string `yaml:"before"`
		ReturnsTo    []string `yaml:"returns_to"`
	} `yaml:"flowforge_agent"`
}

func subagentsSourceDir(t *testing.T) string {
	t.Helper()
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	return filepath.Join(repoRoot, "assets", "subagents")
}

func splitSubagentFrontmatter(t *testing.T, data []byte) (frontmatter, body string) {
	t.Helper()
	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		t.Fatalf("subagent source is missing opening frontmatter delimiter")
	}
	rest := content[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		t.Fatalf("subagent source is missing closing frontmatter delimiter")
	}
	return rest[:end], rest[end+len("\n---\n"):]
}

func TestSubagentSourceFilesExist(t *testing.T) {
	dir := subagentsSourceDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	var found []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		found = append(found, strings.TrimSuffix(entry.Name(), ".md"))
	}
	sort.Strings(found)
	want := append([]string{}, expectedSubagentNames...)
	sort.Strings(want)
	if len(found) != len(want) {
		t.Fatalf("expected %d subagent source files, found %d: %v", len(want), len(found), found)
	}
	for i := range want {
		if found[i] != want[i] {
			t.Fatalf("subagent source file mismatch: found %v, want %v", found, want)
		}
	}
}

func TestSubagentSourceFrontmatterValid(t *testing.T) {
	dir := subagentsSourceDir(t)
	for _, name := range expectedSubagentNames {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name+".md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}
			frontmatter, _ := splitSubagentFrontmatter(t, data)
			var fm subagentFrontmatter
			if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
				t.Fatalf("parsing frontmatter for %s: %v", path, err)
			}
			if fm.FlowforgeAgent.Name != name {
				t.Errorf("%s: frontmatter name %q does not match filename %q", path, fm.FlowforgeAgent.Name, name)
			}
			if fm.FlowforgeAgent.Description == "" {
				t.Errorf("%s: missing description", path)
			}
			if !validModelProfiles[fm.FlowforgeAgent.ModelProfile] {
				t.Errorf("%s: invalid model_profile %q", path, fm.FlowforgeAgent.ModelProfile)
			}
			if fm.FlowforgeAgent.DefaultSkill == "" {
				t.Errorf("%s: missing default_skill", path)
			}
		})
	}
}

func TestSubagentSourceDefaultSkillResolves(t *testing.T) {
	dir := subagentsSourceDir(t)
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	for _, name := range expectedSubagentNames {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name+".md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}
			frontmatter, _ := splitSubagentFrontmatter(t, data)
			var fm subagentFrontmatter
			if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
				t.Fatalf("parsing frontmatter for %s: %v", path, err)
			}
			skillPath := filepath.Join(repoRoot, "assets", "skills", fm.FlowforgeAgent.DefaultSkill, "SKILL.md")
			if _, err := os.Stat(skillPath); err != nil {
				t.Errorf("%s: default_skill %q does not resolve to %s: %v", path, fm.FlowforgeAgent.DefaultSkill, skillPath, err)
			}
			for _, detour := range fm.FlowforgeAgent.DetourSkills {
				detourPath := filepath.Join(repoRoot, "assets", "skills", detour, "SKILL.md")
				if _, err := os.Stat(detourPath); err != nil {
					t.Errorf("%s: detour_skill %q does not resolve to %s: %v", path, detour, detourPath, err)
				}
			}
		})
	}
}

func TestSubagentSourceHasFiveSections(t *testing.T) {
	dir := subagentsSourceDir(t)
	for _, name := range expectedSubagentNames {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name+".md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}
			_, body := splitSubagentFrontmatter(t, data)
			for _, section := range requiredSubagentSections {
				if !strings.Contains(body, section) {
					t.Errorf("%s: missing required section %q", path, section)
				}
			}
		})
	}
}
