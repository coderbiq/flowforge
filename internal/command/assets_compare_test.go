package command

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareManagedAssetsClassifiesWithoutMutatingProject(t *testing.T) {
	assets := t.TempDir()
	project := t.TempDir()
	writeAsset := func(root, rel, content string) string {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}
	writeAsset(assets, "skills/current/SKILL.md", "current")
	writeAsset(assets, "skills/drifted/SKILL.md", "embedded")
	writeAsset(assets, "skills/missing/SKILL.md", "missing")
	writeAsset(assets, "agents/issue-tracker.md", "agent-rule")
	writeAsset(project, ".agents/skills/current/SKILL.md", "current")
	writeAsset(project, ".agents/skills/drifted/SKILL.md", "project")
	owned := writeAsset(project, ".agents/skills/custom/SKILL.md", "custom")
	ownedMarker := writeAsset(project, ".agents/skills/custom/.gitkeep", "")
	writeAsset(project, "wiki/agents/issue-tracker.md", "agent-rule")

	comparison, err := compareManagedAssets(assets, project, "wiki")
	if err != nil {
		t.Fatal(err)
	}
	states := map[string]managedAssetState{}
	for _, entry := range comparison.Entries {
		states[entry.TargetPath] = entry.State
	}
	for path, want := range map[string]managedAssetState{
		filepath.Join(project, ".agents/skills/current/SKILL.md"): managedAssetCurrent,
		filepath.Join(project, ".agents/skills/drifted/SKILL.md"): managedAssetDrifted,
		filepath.Join(project, ".agents/skills/missing/SKILL.md"): managedAssetMissing,
		filepath.Join(project, "wiki/agents/issue-tracker.md"):    managedAssetCurrent,
		owned:       managedAssetProjectOwned,
		ownedMarker: managedAssetProjectOwned,
	} {
		if got := states[path]; got != want {
			t.Errorf("%s state = %q, want %q", path, got, want)
		}
	}
	if comparison.IsCurrent() {
		t.Fatal("missing and drifted managed assets reported current")
	}
	data, err := os.ReadFile(owned)
	if err != nil || string(data) != "custom" {
		t.Fatalf("comparison mutated project-owned file: %q, %v", data, err)
	}
}

func TestCompareManagedAssetsUsesAbsoluteDocsRoot(t *testing.T) {
	assets := t.TempDir()
	project := t.TempDir()
	docsRoot := filepath.Join(t.TempDir(), "wiki")
	if err := os.MkdirAll(filepath.Join(assets, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(assets, "agents", "issue-tracker.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("rule"), 0o644); err != nil {
		t.Fatal(err)
	}
	comparison, err := compareManagedAssets(assets, project, docsRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(comparison.Entries) != 1 || comparison.Entries[0].TargetPath != filepath.Join(docsRoot, "agents", "issue-tracker.md") || comparison.Entries[0].State != managedAssetMissing {
		t.Fatalf("absolute docs root was not used: %#v", comparison.Entries)
	}
}
