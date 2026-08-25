package command

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var localMarkdownLink = regexp.MustCompile(`\[[^]]+\]\(([^)]+)\)`)

func TestPackagedSkillPointersResolve(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	assertSkillPointersResolve(t, filepath.Join(repoRoot, "assets", "skills"))
	assertRequiredArtifactContractPointers(t, filepath.Join(repoRoot, "assets", "skills"))

	target := t.TempDir()
	if err := deployManagedAssets(target, "docs"); err != nil {
		t.Fatal(err)
	}
	assertSkillPointersResolve(t, filepath.Join(target, ".agents", "skills"))
	assertRequiredArtifactContractPointers(t, filepath.Join(target, ".agents", "skills"))
}

func assertRequiredArtifactContractPointers(t *testing.T, root string) {
	t.Helper()
	required := map[string][]string{
		"flowforge-align":     {"roles-and-authority", "hand-offs", "information-value"},
		"flowforge-route":     {"roles-and-authority", "hand-offs"},
		"flowforge-to-spec":   {"roles-and-authority", "hand-offs", "information-value"},
		"flowforge-plan":      {"packaging", "hand-offs", "information-value"},
		"flowforge-wayfinder": {"packaging", "hand-offs"},
		"flowforge-implement": {"hand-offs", "diagnostics"},
		"flowforge-tdd":       {"hand-offs", "diagnostics"},
		"flowforge-review":    {"roles-and-authority", "information-value"},
		"flowforge-handoff":   {"hand-offs"},
		"flowforge-solution-design": {"roles-and-authority", "packaging", "hand-offs", "diagnostics", "information-value"},
	}
	for skill, anchors := range required {
		data, err := os.ReadFile(filepath.Join(root, skill, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		for _, anchor := range anchors {
			needle := "../_shared/ARTIFACT-CONTRACT.md#" + anchor
			if !strings.Contains(string(data), needle) {
				t.Errorf("%s is missing required contract pointer %s", skill, needle)
			}
		}
	}
}

func assertSkillPointersResolve(t *testing.T, root string) {
	t.Helper()
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, match := range localMarkdownLink.FindAllStringSubmatch(string(data), -1) {
			parts := strings.SplitN(match[1], "#", 2)
			target := parts[0]
			if target == "" || target == "link" || strings.Contains(target, "://") {
				continue
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(target)))
			if _, err := os.Stat(resolved); err != nil {
				t.Errorf("broken Skill pointer %s -> %s", path, match[1])
				continue
			}
			if len(parts) == 2 && parts[1] != "" {
				body, err := os.ReadFile(resolved)
				if err != nil {
					return err
				}
				if !markdownHasAnchor(string(body), parts[1]) {
					t.Errorf("broken Skill anchor %s -> %s", path, match[1])
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func markdownHasAnchor(body, wanted string) bool {
	for _, line := range strings.Split(body, "\n") {
		title := strings.TrimSpace(strings.TrimLeft(line, "#"))
		if title == line {
			continue
		}
		slug := strings.ToLower(title)
		slug = strings.ReplaceAll(slug, " ", "-")
		slug = strings.ReplaceAll(slug, "/", "")
		if slug == wanted {
			return true
		}
	}
	return false
}
