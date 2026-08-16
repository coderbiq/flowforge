package core

import (
	"os"
	"path/filepath"
	"testing"
)

type fixtureBaseline struct {
	Path   string
	Body   []byte
	SHA256 string
}

type manifestFixture struct {
	Root      string
	Manifest  *ProjectManifest
	Baselines []fixtureBaseline
}

func v1ManifestFixture(t *testing.T) manifestFixture {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, configDirName), 0755); err != nil {
		t.Fatal(err)
	}
	files := []struct {
		source string
		target string
		type_  string
		body   string
	}{
		{"assets/AGENTS.md", "AGENTS.md", "agents_block", "legacy base\n"},
		{"generated/opencode/legacy.md", ".opencode/agents/legacy.md", "opencode_agent", "legacy opencode\n"},
		{"generated/codex/legacy.toml", ".codex/agents/legacy.toml", "codex_agent", "legacy codex\n"},
	}
	manifest := &ProjectManifest{Version: ManifestVersionV1, CLIVersion: "fixture-v1"}
	var baselines []fixtureBaseline
	for _, file := range files {
		path := filepath.Join(root, filepath.FromSlash(file.target))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		body := []byte(file.body)
		if err := os.WriteFile(path, body, 0644); err != nil {
			t.Fatal(err)
		}
		manifest.Files = append(manifest.Files, FileEntry{Source: file.source, Target: file.target, SHA256: SHA256Hex(body), Type: file.type_})
		baselines = append(baselines, fixtureBaseline{Path: file.target, Body: append([]byte(nil), body...), SHA256: SHA256Hex(body)})
	}
	// This is intentionally a literal v1 document: migration tests must not
	// accidentally exercise the v2 serializer while constructing their input.
	const yamlV1 = "version: 1\ncli_version: fixture-v1\nfiles:\n- source: assets/AGENTS.md\n  target: AGENTS.md\n  sha256: 9b4e0a4b8c4f4d11c8b49f0e0be3e4b3c2c1f4f8f3b7a7d8f1b0c5f0c8c5f9d0\n  type: agents_block\n- source: generated/opencode/legacy.md\n  target: .opencode/agents/legacy.md\n  sha256: 2a4e0b5f2cfb5d1f9d0e53b2b6f7c8e4f1e8e2d1f2b3c4d5e6f708192a3b4c5d\n  type: opencode_agent\n- source: generated/codex/legacy.toml\n  target: .codex/agents/legacy.toml\n  sha256: 8d3c1c7f5d8f4a1f6c2b3e4d5a6f708192a3b4c5d6e7f8091a2b3c4d5e6f7081\n  type: codex_agent\n"
	// Replace the illustrative checksums with the actual fixture baselines,
	// while retaining a deterministic v1 serialization shape.
	manifestData := []byte("version: 1\ncli_version: fixture-v1\nfiles:\n")
	for i, file := range files {
		manifestData = append(manifestData, []byte("- source: "+file.source+"\n  target: "+file.target+"\n  sha256: "+baselines[i].SHA256+"\n  type: "+file.type_+"\n")...)
	}
	if string(manifestData) == yamlV1 {
		t.Fatalf("fixture unexpectedly retained placeholder manifest")
	}
	manifestPath := filepath.Join(root, configDirName, ManifestFileName)
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		t.Fatal(err)
	}
	baselines = append(baselines, fixtureBaseline{Path: filepath.Join(configDirName, ManifestFileName), Body: append([]byte(nil), manifestData...), SHA256: SHA256Hex(manifestData)})
	return manifestFixture{Root: root, Manifest: manifest, Baselines: baselines}
}

func v2ManifestFixture(t *testing.T) manifestFixture {
	t.Helper()
	root := t.TempDir()
	manifest := &ProjectManifest{
		Version: ManifestVersionV2, CLIVersion: "fixture-v2", Mode: ManifestModeSubagent,
		HostIntent: HostIntent{OpenCode: HostEnabled, Codex: HostDisabled},
		Files: []FileEntry{
			{Source: "generated/codex/agent.toml", Target: ".codex/agents/agent.toml", SHA256: SHA256Hex([]byte("codex\n")), Type: "codex_agent", Host: HostCodex, Dormant: true},
			{Source: "generated/opencode/agent.md", Target: ".opencode/agents/agent.md", SHA256: SHA256Hex([]byte("opencode\n")), Type: "opencode_agent", Host: HostOpenCode},
		},
	}
	for _, file := range []struct{ target, body string }{{".codex/agents/agent.toml", "codex\n"}, {".opencode/agents/agent.md", "opencode\n"}} {
		path := filepath.Join(root, filepath.FromSlash(file.target))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(file.body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, configDirName, ManifestFileName))
	if err != nil {
		t.Fatal(err)
	}
	return manifestFixture{Root: root, Manifest: manifest, Baselines: []fixtureBaseline{{Path: filepath.Join(configDirName, ManifestFileName), Body: data, SHA256: SHA256Hex(data)}}}
}

func TestManifestFixturesSelfCheckAndRoundTrip(t *testing.T) {
	for _, fixture := range []manifestFixture{v1ManifestFixture(t), v2ManifestFixture(t)} {
		for _, baseline := range fixture.Baselines {
			data, err := os.ReadFile(filepath.Join(fixture.Root, baseline.Path))
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != string(baseline.Body) || SHA256Hex(data) != baseline.SHA256 {
				t.Fatalf("fixture baseline drifted: %s", baseline.Path)
			}
		}
		if _, err := LoadProjectManifest(fixture.Root); err != nil {
			t.Fatalf("fixture manifest does not round-trip: %v", err)
		}
	}
}
