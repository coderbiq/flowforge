package command

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"flowforge/internal/core"
)

func disableFixture(t *testing.T, modifiedFile, modifiedAgents bool) (string, []byte, []byte) {
	t.Helper()
	root := t.TempDir()
	fileTarget := ".opencode/agents/flowforge-coordinator.md"
	fileData := []byte("managed coordinator\n")
	if modifiedFile {
		fileData = []byte("user modified coordinator\n")
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, fileTarget)), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, fileTarget), fileData, 0644); err != nil {
		t.Fatal(err)
	}
	agents := []byte("user before\n<!-- FLOWFORGE:START -->\nbase\n\n<!-- FLOWFORGE:ORCHESTRATION:START -->\nmanaged rules\n<!-- FLOWFORGE:ORCHESTRATION:END -->\n<!-- tool:START -->\ntool\n<!-- tool:END -->\n<!-- FLOWFORGE:END -->\nuser after\n")
	if modifiedAgents {
		agents = bytes.Replace(agents, []byte("managed rules"), []byte("user changed rules"), 1)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), agents, 0644); err != nil {
		t.Fatal(err)
	}
	m := &core.ProjectManifest{Version: core.ManifestVersionV2, Mode: core.ManifestModeSubagent, HostIntent: core.HostIntent{OpenCode: core.HostEnabled, Codex: core.HostDisabled}, Files: []core.FileEntry{
		{Source: "generated/opencode/flowforge-coordinator.md", Target: fileTarget, SHA256: core.SHA256Hex([]byte("managed coordinator\n")), Type: "opencode_agent", Host: core.HostOpenCode},
		{Source: "generated/AGENTS.orchestration.md", Target: "AGENTS.md", SHA256: core.SHA256Hex([]byte("managed rules\n")), Type: "orchestration_block", Markers: &core.BlockMarkers{Start: "<!-- FLOWFORGE:ORCHESTRATION:START -->", End: "<!-- FLOWFORGE:ORCHESTRATION:END -->"}},
	}}
	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	manifestBytes, err := os.ReadFile(filepath.Join(root, ".flowforge", core.ManifestFileName))
	if err != nil {
		t.Fatal(err)
	}
	return root, manifestBytes, agents
}

func TestDisableCleanModifiedAgentsAndUnknownSentinel(t *testing.T) {
	root, _, originalAgents := disableFixture(t, true, true)
	sentinel := filepath.Join(root, ".opencode", "agents", "sentinel.md")
	if err := os.WriteFile(sentinel, []byte("keep\n"), 0644); err != nil {
		t.Fatal(err)
	}
	plan, err := disableProject(root, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Backups) != 2 || !strings.Contains(strings.Join(plan.Statuses, "\n"), "modified-but-authorized") {
		t.Fatalf("unexpected plan: %#v", plan)
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode/agents/flowforge-coordinator.md")); !os.IsNotExist(err) {
		t.Fatal("registered file was not deleted")
	}
	gotAgents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(gotAgents, originalAgents) || !bytes.Contains(gotAgents, []byte("base")) || !bytes.Contains(gotAgents, []byte("tool")) {
		t.Fatalf("AGENTS preservation failed: %q", gotAgents)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("unregistered sentinel removed: %v", err)
	}
	got, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mode != core.ManifestModeNonSubagent || got.HostIntent.OpenCode != core.HostDisabled || len(got.Files) != 0 {
		t.Fatalf("manifest not disabled: %#v", got)
	}
}

func TestDisableSingleHostPreservesSharedOrchestrationForRemainingEnabledHost(t *testing.T) {
	root, _, _ := disableFixture(t, false, false)
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest.HostIntent.Codex = core.HostEnabled
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}

	if _, err := disableProject(root, []string{core.HostOpenCode}, false); err != nil {
		t.Fatal(err)
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(agents, []byte(orchestrationBlockStart)) || !bytes.Contains(agents, []byte("managed rules")) {
		t.Fatalf("shared orchestration block was removed: %q", agents)
	}
	got, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.DynamicEntriesForHost(core.HostOpenCode)) != 0 || len(got.Files) == 0 {
		t.Fatalf("disabled host entries were not removed: %#v", got.Files)
	}
	if _, err := core.ValidateManifestTarget(got, "AGENTS.md"); err != nil {
		t.Fatalf("shared orchestration entry was not preserved: %v", err)
	}
}

func TestDisableMissingIsIdempotentAndDryRunWritesNothing(t *testing.T) {
	root, before, _ := disableFixture(t, false, false)
	if err := os.Remove(filepath.Join(root, ".opencode/agents/flowforge-coordinator.md")); err != nil {
		t.Fatal(err)
	}
	beforeAgents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := disableProject(root, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(plan.Statuses, "\n"), "missing/already absent") {
		t.Fatalf("missing status absent: %#v", plan.Statuses)
	}
	after, err := os.ReadFile(filepath.Join(root, ".flowforge/manifest.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("dry-run changed manifest")
	}
	afterAgents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(beforeAgents, afterAgents) {
		t.Fatal("dry-run changed AGENTS")
	}
	if _, err := os.Stat(filepath.Join(root, ".flowforge/backups")); !os.IsNotExist(err) {
		t.Fatal("dry-run created backup directory")
	}
	if _, err := disableProject(root, nil, false); err != nil {
		t.Fatal(err)
	}
	if _, err := disableProject(root, nil, false); err != nil {
		t.Fatal(err)
	}
}

func TestDisableBackupFailureDeletesNothing(t *testing.T) {
	root, before, _ := disableFixture(t, false, false)
	oldWrite := disableWriteFile
	disableWriteFile = func(name string, data []byte, perm os.FileMode) error { return fmt.Errorf("injected backup failure") }
	t.Cleanup(func() { disableWriteFile = oldWrite })
	if _, err := disableProject(root, nil, false); err == nil {
		t.Fatal("expected backup failure")
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode/agents/flowforge-coordinator.md")); err != nil {
		t.Fatalf("file deleted after backup failure: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(root, ".flowforge/manifest.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("manifest changed after backup failure")
	}
	backupRoot := filepath.Join(root, ".flowforge", "backups")
	if entries, err := os.ReadDir(backupRoot); err == nil && len(entries) > 0 {
		t.Fatalf("backup failure left backup artifacts: %v", entries)
	} else if err != nil && !os.IsNotExist(err) {
		t.Fatalf("checking backup artifacts: %v", err)
	}
}

func TestDisableBackupsExistBeforeEveryDestructiveMutation(t *testing.T) {
	root, _, originalAgents := disableFixture(t, true, false)
	oldRemove := disableRemove
	var removed []string
	disableRemove = func(name string) error {
		removed = append(removed, name)
		planDir := filepath.Join(root, ".flowforge", "backups", "subagent-disable")
		entries, err := os.ReadDir(planDir)
		if err != nil {
			return fmt.Errorf("backup directory unavailable before deleting %s: %w", name, err)
		}
		if len(entries) != 1 {
			return fmt.Errorf("expected one allocated backup before deleting %s, got %d", name, len(entries))
		}
		rel, err := filepath.Rel(root, name)
		if err != nil {
			return fmt.Errorf("finding relative target %s: %w", name, err)
		}
		backup := filepath.Join(planDir, entries[0].Name(), rel)
		if _, err := os.Stat(backup); err != nil {
			return fmt.Errorf("backup for %s missing before delete: %w", name, err)
		}
		return oldRemove(name)
	}
	t.Cleanup(func() { disableRemove = oldRemove })
	if _, err := disableProject(root, nil, false); err != nil {
		t.Fatal(err)
	}
	if len(removed) != 1 || removed[0] != filepath.Join(root, ".opencode", "agents", "flowforge-coordinator.md") {
		t.Fatalf("unexpected destructive sequence: %v", removed)
	}
	gotAgents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(gotAgents, originalAgents) || !bytes.Contains(gotAgents, []byte("base")) || !bytes.Contains(gotAgents, []byte("tool")) {
		t.Fatalf("AGENTS blocks were not preserved: %q", gotAgents)
	}
}

func TestDisableTimestampCollisionUsesUniqueBackupDirectory(t *testing.T) {
	root, _, _ := disableFixture(t, false, false)
	fixed := time.Date(2026, time.August, 15, 8, 0, 0, 0, time.UTC)
	oldNow := disableNow
	disableNow = func() time.Time { return fixed }
	t.Cleanup(func() { disableNow = oldNow })
	base := filepath.Join(root, ".flowforge", "backups", "subagent-disable")
	first := filepath.Join(base, fixed.Format("20060102T150405Z"))
	if err := os.MkdirAll(first, 0755); err != nil {
		t.Fatal(err)
	}
	plan, err := disableProject(root, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if plan.BackupDir != first+"-1" {
		t.Fatalf("timestamp collision overwrote or skipped expected suffix: %s", plan.BackupDir)
	}
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("pre-existing backup directory was lost: %v", err)
	}
}

func TestSubagentEnableSingleAndDualHostCLI(t *testing.T) {
	for _, hosts := range [][]string{{core.HostOpenCode}, {core.HostCodex}, {core.HostOpenCode, core.HostCodex}} {
		t.Run(strings.Join(hosts, "-"), func(t *testing.T) {
			root := t.TempDir()
			if err := runInit(root, true, "default"); err != nil {
				t.Fatal(err)
			}
			oldWD, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(root); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.Chdir(oldWD); err != nil {
					t.Errorf("restore cwd: %v", err)
				}
			})
			cmd := newSubagentEnableCmd()
			args := make([]string, 0, len(hosts)*2)
			for _, host := range hosts {
				args = append(args, "--host", host)
			}
			cmd.SetArgs(args)
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(stdout.String(), "Synchronized hosts:") || stderr.Len() != 0 {
				t.Fatalf("enable output contract changed: stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			manifest, err := core.LoadProjectManifest(root)
			if err != nil {
				t.Fatal(err)
			}
			for _, host := range hosts {
				if hostIntent(manifest, host) != core.HostEnabled || len(manifest.DynamicEntriesForHost(host)) == 0 {
					t.Fatalf("host was not enabled and registered: %s %#v", host, manifest)
				}
			}
		})
	}
}

func TestDisableDualHostBacksUpAndRemovesOnlyRegisteredFacilities(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{core.HostOpenCode, core.HostCodex}, explicitEnable: true}); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(root, ".codex", "agents", "unmanaged.md")
	if err := os.MkdirAll(filepath.Dir(sentinel), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sentinel, []byte("sentinel\n"), 0644); err != nil {
		t.Fatal(err)
	}
	plan, err := disableProject(root, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Backups) < 3 {
		t.Fatalf("dual-host disable did not plan all registered facilities: %#v", plan)
	}
	for _, backup := range plan.Backups {
		data, readErr := os.ReadFile(backup.Destination)
		if readErr != nil {
			t.Fatalf("missing backup for %s: %v", backup.Target, readErr)
		}
		if len(data) == 0 {
			t.Fatalf("empty backup for %s", backup.Target)
		}
		if _, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(backup.Target))); !os.IsNotExist(statErr) && backup.Target != "AGENTS.md" {
			t.Fatalf("registered host facility survived disable: %s", backup.Target)
		}
	}
	if data, err := os.ReadFile(sentinel); err != nil || string(data) != "sentinel\n" {
		t.Fatalf("unmanaged sentinel changed: %q, %v", data, err)
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Mode != core.ManifestModeNonSubagent || len(manifest.DynamicEntriesForHost(core.HostOpenCode)) != 0 || len(manifest.DynamicEntriesForHost(core.HostCodex)) != 0 {
		t.Fatalf("dual-host disable did not clear dynamic manifest state: %#v", manifest)
	}
}

func TestSubagentStatusCLIIsReadOnlyAndSeparatesStreams(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, true, "default"); err != nil {
		t.Fatal(err)
	}
	before, err := filesystemSnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	})
	cmd := newSubagentStatusCmd()
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "read-only") || !strings.Contains(stdout.String(), "opencode") || !strings.Contains(stdout.String(), "codex") {
		t.Fatalf("status stdout missing exact read-only report: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("status wrote warnings to stderr: %q", stderr.String())
	}
	after, err := filesystemSnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatal("status changed project snapshot")
	}
}

func TestSubagentHostValidationSeparatesErrorFromOutput(t *testing.T) {
	for _, args := range [][]string{{}, {"--host", ""}, {"--host", "invalid"}, {"--host", "opencode", "--host", "opencode"}} {
		cmd := newSubagentEnableCmd()
		cmd.SetArgs(args)
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		err := cmd.Execute()
		if err == nil {
			t.Fatalf("expected non-zero command result for args %#v", args)
		}
		if !strings.Contains(stdout.String(), "Usage:") || !strings.Contains(stderr.String(), "Error:") {
			t.Fatalf("validation stream contract changed for args %#v: stdout=%q stderr=%q", args, stdout.String(), stderr.String())
		}
	}
}

func TestSubagentHelpDocumentsExplicitLifecycle(t *testing.T) {
	root := newSubagentCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"enable", "disable", "status", "explicit"} {
		if !strings.Contains(strings.ToLower(text), want) {
			t.Fatalf("help missing %q:\n%s", want, text)
		}
	}
	disable := newSubagentDisableCmd()
	out.Reset()
	disable.SetOut(&out)
	disable.SetErr(&out)
	disable.SetArgs([]string{"--help"})
	if err := disable.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"delete", "dry-run"} {
		if !strings.Contains(strings.ToLower(out.String()), want) {
			t.Fatalf("disable help missing %q:\n%s", want, out.String())
		}
	}
}

func TestSubagentHostValidationRejectsEmptyInvalidAndDuplicate(t *testing.T) {
	for _, hosts := range [][]string{nil, {}, {"other"}, {"opencode", "opencode"}} {
		if err := validateExplicitHosts(hosts); err == nil {
			t.Fatalf("expected host validation error for %#v", hosts)
		}
	}
	if err := validateHosts([]string{"codex", "bad"}); err == nil {
		t.Fatal("expected invalid host error")
	}
}

func TestDetectHostEvidenceIsReadOnlyForNonSubagent(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".opencode"), 0755); err != nil {
		t.Fatal(err)
	}
	manifest := &core.ProjectManifest{
		Version:    core.ManifestVersionV2,
		Mode:       core.ManifestModeNonSubagent,
		HostIntent: core.HostIntent{OpenCode: core.HostDisabled, Codex: core.HostDisabled},
		Files:      []core.FileEntry{{Source: "assets/skill.md", Target: ".agents/skills/skill.md", Type: "skill"}},
	}
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, ".flowforge", core.ManifestFileName)
	beforeManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	evidence, err := DetectHostEvidence(root, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !evidence[0].Detected || evidence[0].Intent != core.HostDisabled {
		t.Fatalf("unexpected opencode evidence: %#v", evidence[0])
	}
	afterManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(beforeManifest, afterManifest) {
		t.Fatal("read-only host detection changed manifest")
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 { // .flowforge and .opencode only
		t.Fatalf("read-only detection changed filesystem: %#v", entries)
	}
}
