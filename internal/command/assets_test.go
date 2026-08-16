package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"flowforge/internal/core"
)

func TestAssetUpgradeManifestReadFailureIsExplicitAndNonMutating(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	before, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, []byte("version: [invalid"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := applyAssetUpdates(root, false); err == nil {
		t.Fatal("expected manifest read failure")
	} else if !bytes.Contains([]byte(err.Error()), []byte("loading existing project manifest")) {
		t.Fatalf("manifest failure was not explicit: %v", err)
	}
	got, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got, before) {
		t.Fatal("fixture corruption unexpectedly disappeared")
	}
}

func TestSyncDryRunUsesHostPlanWithoutWriting(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifestBefore, err := os.ReadFile(filepath.Join(root, ".flowforge", core.ManifestFileName))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd := newSyncCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := syncProject(cmd, root, syncOptions{forced: []string{core.HostOpenCode}, dryRun: true, explicitEnable: true}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out.Bytes(), []byte(".opencode/agents/flowforge-executor.md")) || !bytes.Contains(out.Bytes(), []byte("AGENTS.md")) {
		t.Fatalf("dry-run omitted host/orchestration plan: %s", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "agents")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote host files: %v", err)
	}
	manifestAfter, err := os.ReadFile(filepath.Join(root, ".flowforge", core.ManifestFileName))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(manifestBefore, manifestAfter) {
		t.Fatal("dry-run changed manifest")
	}
}

func TestAssetUpgradeErrorIncludesCause(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, ".flowforge", core.ManifestFileName)); err != nil {
		t.Fatal(err)
	}
	_, err := applyAssetUpdates(root, false)
	if !errContains(err, "loading existing project manifest") {
		t.Fatalf("expected explicit missing-manifest error, got %v", err)
	}
}

func errContains(err error, want string) bool {
	return err != nil && bytes.Contains([]byte(err.Error()), []byte(want))
}

func TestAssetUpgradeV1MigratesDormantWithoutGeneratingHosts(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest.Version = core.ManifestVersionV1
	manifest.Mode = ""
	manifest.HostIntent = core.HostIntent{}
	manifest.Renderer = core.RendererMetadata{}
	manifest.Files = append(manifest.Files, core.FileEntry{
		Source: "generated/opencode/legacy.md",
		Target: ".opencode/agents/legacy.md",
		SHA256: "legacy-baseline",
		Type:   "opencode_agent",
	})
	if err := manifest.Save(root); err == nil {
		t.Fatal("v1 manifest unexpectedly passed v2 save validation")
	}
	// Write the v1 fixture directly because SaveManifestAtomic intentionally
	// accepts only the current schema.
	data := []byte("version: 1\ncli_version: v1\nfiles:\n- source: assets/AGENTS.md\n  target: AGENTS.md\n  sha256: " + manifest.Files[0].SHA256 + "\n  type: agents_block\n- source: generated/opencode/legacy.md\n  target: .opencode/agents/legacy.md\n  sha256: legacy-baseline\n  type: opencode_agent\n")
	if err := os.WriteFile(filepath.Join(root, ".flowforge", core.ManifestFileName), data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	got, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != core.ManifestVersionV2 || got.Mode != core.ManifestModeNonSubagent || got.HostIntent.OpenCode != core.HostDisabled {
		t.Fatalf("v1 upgrade did not remain disabled: %#v", got)
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode/agents/legacy.md")); !os.IsNotExist(err) {
		t.Fatalf("v1 upgrade generated a host file: %v", err)
	}
	for _, entry := range got.Files {
		if entry.Type == "opencode_agent" && (!entry.Dormant || entry.Host != core.HostOpenCode) {
			t.Fatalf("legacy dynamic entry was not dormant: %#v", entry)
		}
	}
}

func TestAssetUpgradePreservesV2IntentAndDynamicEntries(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest.Mode = core.ManifestModeSubagent
	manifest.HostIntent = core.HostIntent{OpenCode: core.HostEnabled, Codex: core.HostDisabled}
	manifest.Renderer = core.RendererMetadata{PolicyDigest: "policy", Hosts: map[string]string{core.HostOpenCode: "renderer:old"}}
	dynamic := core.FileEntry{Source: "generated/opencode/legacy.md", Target: ".opencode/agents/legacy.md", SHA256: "baseline", Type: "opencode_agent", Host: core.HostOpenCode}
	manifest.Files = append(manifest.Files, dynamic)
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if _, err := applyAssetUpdates(root, false); err != nil {
		t.Fatal(err)
	}
	got, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mode != manifest.Mode || got.HostIntent != manifest.HostIntent || got.Renderer.PolicyDigest != "policy" {
		t.Fatalf("v2 metadata changed: %#v", got)
	}
	if len(got.DynamicEntriesForHost(core.HostOpenCode)) != 1 || got.DynamicEntriesForHost(core.HostOpenCode)[0].SHA256 != dynamic.SHA256 {
		t.Fatalf("dynamic baseline changed: %#v", got.DynamicEntriesForHost(core.HostOpenCode))
	}
}

func TestAssetUpgradeConflictKeepsStaticBaseline(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	var target string
	for i := range manifest.Files {
		if manifest.Files[i].Type == "skill" {
			target = manifest.Files[i].Target
			manifest.Files[i].SHA256 = "old-baseline"
			break
		}
	}
	if target == "" {
		t.Fatal("test fixture has no static skill")
	}
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, target), []byte("user modified"), 0644); err != nil {
		t.Fatal(err)
	}
	report, err := applyAssetUpdates(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if report == nil || len(report.Conflict) == 0 {
		t.Fatal("expected static asset conflict")
	}
	got, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range got.Files {
		if entry.Target == target && entry.SHA256 != "old-baseline" {
			t.Fatalf("conflict advanced baseline: %#v", entry)
		}
	}
}

func TestAssetUpgradeRepeatIsIdempotent(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	var target string
	for i := range manifest.Files {
		if manifest.Files[i].Type == "skill" {
			target = manifest.Files[i].Target
			manifest.Files[i].SHA256 = core.SHA256Hex([]byte("old content"))
			break
		}
	}
	if target == "" {
		t.Fatal("test fixture has no static skill")
	}
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, target), []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}
	first, err := applyAssetUpdates(root, false)
	if err != nil || first == nil || len(first.Updated) == 0 {
		t.Fatalf("first upgrade did not update: %#v, %v", first, err)
	}
	updated, err := os.ReadFile(filepath.Join(root, target))
	if err != nil {
		t.Fatal(err)
	}
	second, err := applyAssetUpdates(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if second != nil {
		t.Fatalf("repeat upgrade reported changes: %#v", second)
	}
	repeated, err := os.ReadFile(filepath.Join(root, target))
	if err != nil {
		t.Fatal(err)
	}
	if string(updated) != string(repeated) {
		t.Fatal("repeat upgrade changed asset bytes")
	}
}

func TestApplyUpgradeNeverCreatesDynamicEntries(t *testing.T) {
	root := t.TempDir()
	diff := &core.DiffResult{Added: []core.FileEntry{{Source: "generated/opencode/a.md", Target: ".opencode/agents/a.md", Type: "opencode_agent"}}}
	report := core.ApplyUpgrade(diff, &core.ProjectManifest{}, root, fstest.MapFS{"generated/opencode/a.md": &fstest.MapFile{Data: []byte("agent")}}, filepath.Join(root, ".flowforge", "backup", "v1"))
	if report.Error != nil {
		t.Fatal(report.Error)
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode/agents/a.md")); !os.IsNotExist(err) {
		t.Fatalf("asset upgrade created dynamic entry: %v", err)
	}
}
