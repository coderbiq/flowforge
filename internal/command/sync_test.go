package command

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestSyncDetectsOpenCodeAndUpdatesRoutingBlock(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"flowforge-coordinator.md", "flowforge-design-analyst.md", "flowforge-executor.md"} {
		if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{orchestrationBlockStart, "flowforge-design-analyst", "context preflight", "context risk-review"} {
		if !strings.Contains(string(agents), want) {
			t.Fatalf("AGENTS missing %q", want)
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
	if err := syncProject(cmd, root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".codex", "agents", "flowforge-executor.toml")); err != nil {
		t.Fatal(err)
	}
	var second bytes.Buffer
	cmd.SetOut(&second)
	if err := syncProject(cmd, root, syncOptions{}); err != nil {
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
	target := filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("mode: subagent\n"), []byte("mode: subagent\nmodel: local/executor\n"), 1)
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("model: local/executor")) {
		t.Fatalf("OpenCode model was not preserved:\n%s", data)
	}
	if !bytes.Contains(data, []byte("## Result Contract")) {
		t.Fatal("OpenCode prompt was not regenerated")
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
	target := filepath.Join(root, ".codex", "agents", "flowforge-executor.toml")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("name = \"flowforge-executor\"\n"), []byte("name = \"flowforge-executor\"\nmodel = \"local/executor\"\n"), 1)
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
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

func TestSyncPreservesModifiedOrchestrationBlock(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := runInit(root, true, "default"); err != nil {
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

func TestRootHasNoDeprecatedAssetCommands(t *testing.T) {
	for _, cmd := range NewRootCmd().Commands() {
		if cmd.Name() == "assets" || cmd.Name() == "skill" {
			t.Fatalf("deprecated command still registered: %s", cmd.Name())
		}
	}
}

func legacyOpenCodeAgent(id, display, profile, skill, edit string) []byte {
	permissions := "read: allow\n  edit: " + edit + "\n  task: deny\n  question: deny\n  skill: allow"
	return []byte(fmt.Sprintf("---\nname: flowforge-%s\ndescription: FlowForge %s role.\nmode: subagent\npermission:\n  %s\n---\n\n# FlowForge %s\n\nActive Role: %s\nModel Profile: %s\nDefault Skill: %s\n\nRead the Proposal Journal and referenced artifacts. Follow the installed FlowForge Skill and return control to the Coordinator. Do not delegate or ask the user directly.\n", id, display, strings.ReplaceAll(permissions, "\n", "\n  "), display, display, profile, skill))
}
