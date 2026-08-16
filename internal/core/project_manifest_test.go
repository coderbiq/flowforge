package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateManifestV1(t *testing.T) {
	tests := []struct {
		name string
		v1   *ProjectManifest
		want string
	}{
		{"legacy hosts remain disabled", &ProjectManifest{Version: 1, DisabledHosts: []string{"opencode"}, PendingHosts: []string{"codex"}}, ""},
		{"dynamic entries become dormant", &ProjectManifest{Version: 1, Files: []FileEntry{{Source: "generated/opencode/a.md", Target: ".opencode/agents/a.md", Type: "opencode_agent"}}}, HostOpenCode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MigrateManifestV1(tt.v1)
			if err != nil {
				t.Fatal(err)
			}
			if got.Version != ManifestVersionV2 || got.Mode != ManifestModeNonSubagent || got.HostIntent.OpenCode != HostDisabled || got.HostIntent.Codex != HostDisabled {
				t.Fatalf("unexpected migration: %#v", got)
			}
			if tt.want != "" && (got.Files[0].Host != tt.want || !got.Files[0].Dormant) {
				t.Fatalf("dynamic entry was not dormant: %#v", got.Files[0])
			}
		})
	}
}

func TestManifestValidationRejectsEntriesWithoutSaving(t *testing.T) {
	tests := []struct {
		name  string
		files []FileEntry
	}{
		{"path traversal", []FileEntry{{Source: "assets/a", Target: "../outside", Type: "skill"}}},
		{"unknown host", []FileEntry{{Source: "generated/a", Target: ".agents/a", Host: "other", Type: "skill"}}},
		{"duplicate source", []FileEntry{{Source: "assets/a", Target: ".agents/a", Type: "skill"}, {Source: "assets/a", Target: ".agents/b", Type: "skill"}}},
		{"duplicate target", []FileEntry{{Source: "assets/a", Target: ".agents/a", Type: "skill"}, {Source: "assets/b", Target: ".agents/a", Type: "skill"}}},
		{"duplicate orchestration target", []FileEntry{{Source: "generated/orchestration/a", Target: "AGENTS.md", Type: "orchestration_block"}, {Source: "generated/orchestration/b", Target: "AGENTS.md", Type: "orchestration_block"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, configDirName), 0755); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, configDirName, ManifestFileName)
			old := []byte("sentinel\n")
			if err := os.WriteFile(path, old, 0644); err != nil {
				t.Fatal(err)
			}
			m := &ProjectManifest{Version: ManifestVersionV2, Mode: ManifestModeNonSubagent, HostIntent: HostIntent{OpenCode: HostDisabled, Codex: HostDisabled}, Files: tt.files}
			if err := m.Save(root); err == nil {
				t.Fatal("expected validation error")
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(old) {
				t.Fatalf("manifest changed: %q", got)
			}
		})
	}
}

func TestManifestRoundTripNormalizesFileOrder(t *testing.T) {
	root := t.TempDir()
	m := &ProjectManifest{Version: ManifestVersionV2, Mode: ManifestModeNonSubagent, HostIntent: HostIntent{OpenCode: HostDisabled, Codex: HostDisabled}, Files: []FileEntry{{Source: "assets/z", Target: ".agents/z", Type: "skill"}, {Source: "assets/a", Target: ".agents/a", Type: "skill"}}}
	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	got, err := LoadProjectManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.Files[0].Source != "assets/a" || got.Files[1].Source != "assets/z" {
		t.Fatalf("files were not sorted: %#v", got.Files)
	}
	data, err := os.ReadFile(filepath.Join(root, configDirName, ManifestFileName))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "disabled_hosts") || !strings.Contains(string(data), "version: 2") {
		t.Fatalf("unexpected v2 serialization: %s", data)
	}
}

func TestDesiredHostSetUsesIntentAndStableRegistration(t *testing.T) {
	m := &ProjectManifest{Version: ManifestVersionV2, Mode: ManifestModeSubagent, HostIntent: HostIntent{OpenCode: HostEnabled, Codex: HostDisabled}, Files: []FileEntry{
		{Source: "generated/opencode/z", Target: ".opencode/agents/z", Type: "opencode_agent", Host: HostOpenCode},
		{Source: "generated/opencode/a", Target: ".opencode/agents/a", Type: "opencode_agent", Host: HostOpenCode},
		{Source: "generated/codex/a", Target: ".codex/agents/a", Type: "codex_agent", Host: HostCodex},
	}}
	got := m.DesiredHostSet(HostOpenCode)
	if len(got) != 2 || got[0].Target > got[1].Target {
		t.Fatalf("unexpected stable desired set: %#v", got)
	}
	if len(m.DesiredHostSet(HostCodex)) != 0 {
		t.Fatal("disabled host produced desired entries")
	}
	m.Mode = ManifestModeNonSubagent
	if len(m.DesiredHostSet(HostOpenCode)) != 0 {
		t.Fatal("non_subagent mode produced desired entries")
	}
}

func TestValidateManifestTargetRequiresRegisteredDynamicEntry(t *testing.T) {
	m := &ProjectManifest{Files: []FileEntry{{Source: "generated/opencode/a", Target: ".opencode/agents/a", Type: "opencode_agent", Host: HostOpenCode}}}
	if _, err := ValidateManifestTarget(m, ".opencode/agents/missing"); err == nil {
		t.Fatal("unregistered target was accepted")
	}
	entry, err := ValidateManifestTarget(m, ".opencode/agents/a")
	if err != nil || entry.Source != "generated/opencode/a" {
		t.Fatalf("registered target rejected: %#v, %v", entry, err)
	}
}
