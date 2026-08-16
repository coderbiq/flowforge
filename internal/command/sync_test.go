package command

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"flowforge/internal/core"
	"flowforge/internal/orchestration"
)

func TestSyncDetectsOpenCodeAndUpdatesRoutingBlock(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"flowforge-coordinator.md", "flowforge-design-analyst.md", "flowforge-investigator.md", "flowforge-executor.md"} {
		if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"### Subagent Orchestration", orchestrationBlockStart, "execution scheduler", "flowforge-design-analyst", "flowforge-investigator", "Delegation depth is one", "structured analysis revision/readiness/re-entry state", "context preflight", "context risk-review"} {
		if !strings.Contains(string(agents), want) {
			t.Fatalf("AGENTS missing %q", want)
		}
	}
}

func TestSyncRendererFailureLeavesProjectBytesUntouched(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	manifestBefore, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	agentsBefore, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	previous := renderOpenCode
	renderOpenCode = func(orchestration.Policy) (orchestration.RenderOutput, error) {
		return orchestration.RenderOutput{}, fmt.Errorf("injected renderer failure")
	}
	t.Cleanup(func() { renderOpenCode = previous })
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{core.HostOpenCode}}); err == nil {
		t.Fatal("expected renderer failure")
	}
	manifestAfter, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	agentsAfter, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(manifestBefore, manifestAfter) || !bytes.Equal(agentsBefore, agentsAfter) {
		t.Fatal("renderer failure changed project bytes")
	}
}

func TestSyncRejectsDynamicTargetThroughOutsideSymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, ".opencode")); err != nil {
		if os.IsNotExist(err) || strings.Contains(strings.ToLower(err.Error()), "not supported") {
			t.Skipf("symlinks are not supported: %v", err)
		}
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest.Mode = core.ManifestModeSubagent
	manifest.HostIntent.OpenCode = core.HostEnabled
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err == nil {
		t.Fatal("sync accepted a dynamic target escaping through a symlink")
	}
	if _, err := os.Stat(filepath.Join(outside, "agents")); !os.IsNotExist(err) {
		t.Fatalf("sync wrote outside the project: %v", err)
	}
}

func TestSyncManifestSaveFailureRollsBackHostAndAgents(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	manifestBefore, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	agentsBefore, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	previous := saveManifest
	saveManifest = func(string, *core.ProjectManifest) error { return fmt.Errorf("injected manifest save failure") }
	t.Cleanup(func() { saveManifest = previous })
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{core.HostOpenCode}}); err == nil {
		t.Fatal("expected manifest save failure")
	}
	manifestAfter, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	agentsAfter, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(manifestBefore, manifestAfter) || !bytes.Equal(agentsBefore, agentsAfter) {
		t.Fatal("manifest save failure changed manifest or AGENTS")
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")); !os.IsNotExist(err) {
		t.Fatalf("manifest save failure left host file: %v", err)
	}
}

func TestSyncLeavesLegacyCardAndHistoryWikiBytesUntouched(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	legacyPath := filepath.Join(root, "ff-wiki", "02-library", "40-tasks", "TASK-legacy.md")
	historyPath := filepath.Join(root, "ff-wiki", "03-proposal", "history-wiki.md")
	legacy := []byte("# legacy task\n\nlegacy body and [[OLD-123]]\n")
	history := []byte("# historical wiki\n\nold links remain byte-for-byte\n")
	for path, data := range map[string][]byte{legacyPath: legacy, historyPath: history} {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string][]byte{legacyPath: legacy, historyPath: history} {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("sync rewrote historical file %s", path)
		}
	}
}

func TestSyncDetectsCodexAndIsIdempotent(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".codex"), 0755); err != nil {
		t.Fatal(err)
	}
	cmd := newSyncCmd()
	var first bytes.Buffer
	cmd.SetOut(&first)
	if err := syncProject(cmd, root, syncOptions{forced: []string{"codex"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".codex", "agents", "flowforge-executor.toml")); err != nil {
		t.Fatal(err)
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"current main thread IS the FlowForge Coordinator", "manual fallback", "user explicitly requests `flowforge-coordinator`", "delegation is mandatory", "MUST spawn `flowforge-design-analyst`", "MUST spawn `flowforge-executor`", "return `BLOCKED`"} {
		if !strings.Contains(string(agents), want) {
			t.Fatalf("Codex AGENTS missing %q", want)
		}
	}
	var second bytes.Buffer
	cmd.SetOut(&second)
	if err := syncProject(cmd, root, syncOptions{forced: []string{"codex"}}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(second.String(), "~ .codex/agents") {
		t.Fatalf("sync rewrote unchanged agents: %s", second.String())
	}
}

func TestSyncPreservesManagedOpenCodeModelDuringUpdate(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("mode: subagent\n"), []byte("mode: subagent\nmodel: local/executor\n"), 1)
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("model: \"local/executor\"")) {
		t.Fatalf("OpenCode model was not preserved:\n%s", data)
	}
	if !bytes.Contains(data, []byte("## Result Contract")) {
		t.Fatal("OpenCode prompt was not regenerated")
	}
}

func TestSyncDoesNotFailOnMalformedOpenCodeFrontmatter(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".opencode", "agents", "flowforge-design-analyst.md")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("mode: subagent\n"), []byte("mode: subagent\nmodel: local/analyst\ninvalid: [legacy: value\n"), 1)
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatalf("malformed legacy frontmatter blocked sync: %v", err)
	}
	updated, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(updated, []byte("model: \"local/analyst\"")) {
		t.Fatal("model was not preserved from malformed frontmatter")
	}
}

func TestSyncAdoptsGeneratedAgentMissingFromManifest(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	filtered := manifest.Files[:0]
	for _, entry := range manifest.Files {
		if entry.Type != "opencode_agent" {
			filtered = append(filtered, entry)
		}
	}
	manifest.Files = filtered
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	manifest, err = core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, entry := range manifest.Files {
		if entry.Type == "opencode_agent" {
			count++
		}
	}
	if count != 4 {
		t.Fatalf("expected all generated agents to be adopted, got %d", count)
	}
}

func TestSyncReportsHostEnforcement(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	cmd := newSyncCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := syncProject(cmd, root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Enforcement opencode:", "write_scope=soft", "external_sources=unsupported"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("sync output missing %q:\n%s", want, out.String())
		}
	}
}

func TestSyncPreservesManagedCodexModelDuringUpdate(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".codex"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"codex"}}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".codex", "agents", "flowforge-executor.toml")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("name = \"flowforge-executor\"\n"), []byte("name = \"flowforge-executor\"\nmodel = \"local/executor\"\n"), 1)
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"codex"}}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("model = 'local/executor'")) && !bytes.Contains(data, []byte("model = \"local/executor\"")) {
		t.Fatalf("Codex model was not preserved:\n%s", data)
	}
	if !bytes.Contains(data, []byte("## Result Contract")) {
		t.Fatal("Codex prompt was not regenerated")
	}
}

func TestSyncPreservesExplicitOpenCodePermissions(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".opencode", "agents", "flowforge-investigator.md")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("question: deny"), []byte("question: ask"), 1)
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	updated, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(updated, []byte("question: ask")) {
		t.Fatalf("OpenCode permission was not preserved:\n%s", updated)
	}
	if !bytes.Contains(updated, []byte("## Investigation Contract")) {
		t.Fatal("OpenCode investigator prompt was not regenerated")
	}
}

func TestSyncPreservesExplicitCodexPermissions(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".codex"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"codex"}}); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".codex", "agents", "flowforge-executor.toml")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("sandbox_mode = \"workspace-write\""), []byte("sandbox_mode = \"danger-full-access\"\napproval_policy = \"on-request\""), 1)
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"codex"}}); err != nil {
		t.Fatal(err)
	}
	updated, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"danger-full-access", "approval_policy", "on-request", "## Result Contract"} {
		if !bytes.Contains(updated, []byte(want)) {
			t.Fatalf("Codex config missing %q:\n%s", want, updated)
		}
	}
}

func TestSyncWithoutHostRemainsDisabled(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	cmd := newSyncCmd()
	if err := syncProject(cmd, root, syncOptions{removed: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")); !os.IsNotExist(err) {
		t.Fatalf("managed agent not removed: %v", err)
	}
	if err := syncProject(cmd, root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")); !os.IsNotExist(err) {
		t.Fatalf("disabled host was reinstalled: %v", err)
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(agents), orchestrationBlockStart) {
		t.Fatal("routing block not removed")
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Mode != core.ManifestModeNonSubagent || manifest.HostIntent.OpenCode != core.HostDisabled || manifest.HostIntent.Codex != core.HostDisabled {
		t.Fatalf("disabled sync changed v2 intent: mode=%s intent=%#v", manifest.Mode, manifest.HostIntent)
	}
}

func TestSyncAdoptsKnownV310Skeleton(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, ".opencode", "agents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	legacy := legacyOpenCodeAgent("executor", "Executor", "tool-capable", "flowforge-implement", "allow")
	if err := os.WriteFile(filepath.Join(dir, "flowforge-executor.md"), legacy, 0644); err != nil {
		t.Fatal(err)
	}
	cmd := newSyncCmd()
	if err := syncProject(cmd, root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "flowforge-executor.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == string(legacy) {
		t.Fatal("legacy skeleton not upgraded")
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range manifest.Files {
		if entry.Target == filepath.Join(".opencode", "agents", "flowforge-executor.md") {
			return
		}
	}
	t.Fatal("adopted agent missing from manifest")
}

func TestSyncPreservesUnknownAgentUnlessAdopted(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, ".opencode", "agents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "flowforge-executor.md")
	if err := os.WriteFile(target, []byte("user agent"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := newSyncCmd()
	if err := syncProject(cmd, root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "user agent" {
		t.Fatal("unknown agent was overwritten")
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(agents), orchestrationBlockStart) {
		t.Fatal("routing advertised an incomplete host installation")
	}
	if err := syncProject(cmd, root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "user agent" {
		t.Fatal("pending conflict was overwritten on the next sync")
	}
	if err := syncProject(cmd, root, syncOptions{forced: []string{"opencode"}, adopt: true}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "user agent" {
		t.Fatal("--adopt did not replace existing agent")
	}
}

func TestSyncDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}, dryRun: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote agent: %v", err)
	}
}

func TestSyncFreshProjectDoesNotInferHostFromEvidence(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	before, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("host evidence changed the manifest")
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "agents")); !os.IsNotExist(err) {
		t.Fatalf("sync inferred opencode from evidence: %v", err)
	}
}

func TestSyncLegacyHostFlagsDoNotMutateIntent(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	before, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	cmd := newSyncCmd()
	cmd.SetArgs([]string{"--host", "opencode"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected legacy flag migration error")
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("legacy sync flag changed intent")
	}
}

func TestSyncEnabledIntentIsIdempotentWithoutHostEvidence(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest.Mode = core.ManifestModeSubagent
	manifest.HostIntent.OpenCode = core.HostEnabled
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	first, err := filesystemSnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	second, err := filesystemSnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("enabled sync was not idempotent")
	}
}

func filesystemSnapshot(root string) (map[string][]byte, error) {
	snapshot := make(map[string][]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		snapshot[rel] = data
		return nil
	})
	return snapshot, err
}

func TestExplicitEnableSupportsSingleAndDualHostsWithoutEvidence(t *testing.T) {
	for _, hosts := range [][]string{{core.HostOpenCode}, {core.HostCodex}, {core.HostOpenCode, core.HostCodex}} {
		t.Run(strings.Join(hosts, "-"), func(t *testing.T) {
			root := t.TempDir()
			if err := runInit(root, true, "default"); err != nil {
				t.Fatal(err)
			}
			if err := syncProject(newSyncCmd(), root, syncOptions{forced: hosts, explicitEnable: true}); err != nil {
				t.Fatal(err)
			}
			manifest, err := core.LoadProjectManifest(root)
			if err != nil {
				t.Fatal(err)
			}
			for _, host := range hosts {
				if hostIntent(manifest, host) != core.HostEnabled {
					t.Fatalf("%s was not enabled", host)
				}
			}
			if len(manifest.DynamicEntriesForHost(core.HostOpenCode)) == 0 && containsHost(hosts, core.HostOpenCode) {
				t.Fatal("opencode dynamic entries were not registered")
			}
			if len(manifest.DynamicEntriesForHost(core.HostCodex)) == 0 && containsHost(hosts, core.HostCodex) {
				t.Fatal("codex dynamic entries were not registered")
			}
		})
	}
}

func TestExplicitEnableConflictKeepsIntentAndEntriesUnchanged(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("unmanaged"), 0644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	before, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{core.HostOpenCode}, explicitEnable: true}); err == nil {
		t.Fatal("expected unmanaged conflict")
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("explicit conflict committed manifest intent or entries")
	}
	if data, err := os.ReadFile(target); err != nil || string(data) != "unmanaged" {
		t.Fatalf("unmanaged file changed: %q, %v", data, err)
	}
}

func TestExplicitEnableDoesNotReconcileUnlistedEnabledHost(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest.Mode = core.ManifestModeSubagent
	manifest.HostIntent.Codex = core.HostEnabled
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{core.HostOpenCode}, explicitEnable: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".codex", "agents")); !os.IsNotExist(err) {
		t.Fatalf("unlisted codex host was reconciled: %v", err)
	}
	manifest, err = core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.HostIntent.Codex != core.HostEnabled {
		t.Fatal("explicit opencode enable changed codex intent")
	}
}

func containsHost(hosts []string, wanted string) bool {
	for _, host := range hosts {
		if host == wanted {
			return true
		}
	}
	return false
}

func TestSyncPreservesModifiedOrchestrationBlock(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{"opencode"}}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "AGENTS.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	modified := strings.Replace(string(data), "Installed hosts: opencode", "Installed hosts: user-edited", 1)
	if err := os.WriteFile(path, []byte(modified), 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Installed hosts: user-edited") {
		t.Fatal("modified orchestration block was overwritten")
	}
}

func TestSyncPreservesModifiedStaticAssetManifestBaseline(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".agents", "skills", "flowforge-design", "SKILL.md")
	if err := os.WriteFile(target, []byte("user modification\n"), 0644); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	for i := range manifest.Files {
		if manifest.Files[i].Target == filepath.Join(".agents", "skills", "flowforge-design", "SKILL.md") {
			manifest.Files[i].SHA256 = strings.Repeat("0", 64)
		}
	}
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "user modification\n" {
		t.Fatal("modified static asset was overwritten")
	}
	manifest, err = core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range manifest.Files {
		if entry.Target == filepath.Join(".agents", "skills", "flowforge-design", "SKILL.md") && entry.SHA256 != strings.Repeat("0", 64) {
			t.Fatalf("conflict advanced the manifest baseline to %s", entry.SHA256)
		}
	}
}

func TestSyncDoesNotAdoptUntrustedStaticAsset(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".agents", "skills", "flowforge-design", "SKILL.md")
	userData := []byte("user-owned skill\n")
	if err := os.WriteFile(target, userData, 0644); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	filtered := manifest.Files[:0]
	for _, entry := range manifest.Files {
		if entry.Target != filepath.Join(".agents", "skills", "flowforge-design", "SKILL.md") {
			filtered = append(filtered, entry)
		}
	}
	manifest.Files = filtered
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{adopt: true}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, userData) {
		t.Fatal("--adopt replaced an untrusted static asset")
	}
	updated, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range updated.Files {
		if entry.Target == filepath.Join(".agents", "skills", "flowforge-design", "SKILL.md") {
			t.Fatal("untrusted conflict advanced the manifest baseline")
		}
	}
}

func TestSyncFromSubdirectoryUsesProjectRootForConflicts(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".agents", "skills", "flowforge-design", "SKILL.md")
	if err := os.WriteFile(target, []byte("user modification\n"), 0644); err != nil {
		t.Fatal(err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	for i := range manifest.Files {
		if manifest.Files[i].Target == filepath.Join(".agents", "skills", "flowforge-design", "SKILL.md") {
			manifest.Files[i].SHA256 = strings.Repeat("0", 64)
		}
	}
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	subdir := filepath.Join(root, "nested")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(subdir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Errorf("restoring working directory: %v", err)
		}
	})
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "user modification\n" {
		t.Fatal("subdirectory sync overwrote project-root conflict")
	}
}

func legacyOpenCodeAgent(id, display, profile, skill, edit string) []byte {
	permissions := "read: allow\n  edit: " + edit + "\n  task: deny\n  question: deny\n  skill: allow"
	return []byte(fmt.Sprintf("---\nname: flowforge-%s\ndescription: FlowForge %s role.\nmode: subagent\npermission:\n  %s\n---\n\n# FlowForge %s\n\nActive Role: %s\nModel Profile: %s\nDefault Skill: %s\n\nRead the Proposal Journal and referenced artifacts. Follow the installed FlowForge Skill and return control to the Coordinator. Do not delegate or ask the user directly.\n", id, display, strings.ReplaceAll(permissions, "\n", "\n  "), display, display, profile, skill))
}
