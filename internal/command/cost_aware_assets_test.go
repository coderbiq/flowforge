package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// Guards delivery integrity; model behavior needs actual controlled execution.
func TestCostAwareRulePackageSourceEmbeddedDeployed(t *testing.T) {
	target := t.TempDir()
	docsRoot := filepath.Join(target, "wiki")
	if err := deployManagedAssets(target, docsRoot); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"skills/_shared/ARTIFACT-CONTRACT.md",
		"skills/flowforge-solution-design/SKILL.md",
		"skills/flowforge-plan/SKILL.md",
		"skills/flowforge-refine-ticket/SKILL.md",
		"skills/flowforge-implement/SKILL.md",
		"skills/flowforge-review/SKILL.md",
		"agents/issue-tracker.md",
	} {
		t.Run(rel, func(t *testing.T) {
			source, err := os.ReadFile(filepath.Join("..", "..", "assets", filepath.FromSlash(rel)))
			if err != nil {
				t.Fatal(err)
			}
			embedded, err := embeddedAssets.ReadFile("assets/" + rel)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(source, embedded) {
				t.Error("source and embedded rule differ; regenerate embedded assets from source")
			}
			deployedPath := filepath.Join(target, ".agents", filepath.FromSlash(rel))
			if rel == "agents/issue-tracker.md" {
				deployedPath = filepath.Join(docsRoot, filepath.FromSlash(rel))
			}
			deployed, err := os.ReadFile(deployedPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(embedded, deployed) {
				t.Error("embedded and freshly deployed rule differ")
			}
		})
	}
}
