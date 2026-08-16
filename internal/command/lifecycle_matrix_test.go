package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/core"
)

// TestLifecycleMatrix keeps the lifecycle contract in one table. In
// particular, host directories are evidence only; manifest intent remains the
// sole authorization for sync to manage dynamic files.
func TestLifecycleMatrix(t *testing.T) {
	tests := []struct {
		name     string
		evidence bool
	}{
		{name: "fresh init without host evidence"},
		{name: "fresh init with host evidence", evidence: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if tt.evidence {
				if err := os.MkdirAll(filepath.Join(root, ".opencode", "agents"), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, ".opencode", "agents", "unmanaged.md"), []byte("host evidence\n"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			if err := runInit(root, true, "default"); err != nil {
				t.Fatal(err)
			}
			manifest, err := core.LoadProjectManifest(root)
			if err != nil {
				t.Fatal(err)
			}
			if manifest.Mode != core.ManifestModeNonSubagent || manifest.HostIntent != (core.HostIntent{OpenCode: core.HostDisabled, Codex: core.HostDisabled}) {
				t.Fatalf("host evidence changed initial intent: %#v", manifest)
			}
			if len(manifest.DynamicEntriesForHost(core.HostOpenCode)) != 0 || len(manifest.DynamicEntriesForHost(core.HostCodex)) != 0 {
				t.Fatalf("fresh init registered dynamic entries: %#v", manifest.Files)
			}
			if _, err := os.Stat(filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")); !os.IsNotExist(err) {
				t.Fatalf("fresh init rendered a host file: %v", err)
			}
		})
	}

	t.Run("disabled sync is a no-op despite host evidence", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, ".codex", "agents"), 0755); err != nil {
			t.Fatal(err)
		}
		unmanaged := filepath.Join(root, ".codex", "agents", "unmanaged.toml")
		if err := os.WriteFile(unmanaged, []byte("unmanaged\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := runInit(root, true, "default"); err != nil {
			t.Fatal(err)
		}
		before, err := filesystemSnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if err := syncProject(newSyncCmd(), root, syncOptions{}); err != nil {
			t.Fatal(err)
		}
		after, err := filesystemSnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !sameSnapshot(before, after) {
			t.Fatal("disabled sync changed project files")
		}
		if got, readErr := os.ReadFile(unmanaged); readErr != nil || string(got) != "unmanaged\n" {
			t.Fatalf("unmanaged host evidence changed: %q, %v", got, readErr)
		}
	})

	t.Run("enabled sync repeat is byte stable", func(t *testing.T) {
		root := t.TempDir()
		if err := runInit(root, true, "default"); err != nil {
			t.Fatal(err)
		}
		if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{core.HostOpenCode}, explicitEnable: true}); err != nil {
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
		if !sameSnapshot(first, second) {
			t.Fatal("repeated enabled sync changed project files")
		}
		manifest, err := core.LoadProjectManifest(root)
		if err != nil {
			t.Fatal(err)
		}
		if manifest.Mode != core.ManifestModeSubagent || manifest.HostIntent.OpenCode != core.HostEnabled {
			t.Fatalf("enabled intent was not retained: %#v", manifest)
		}
	})

	t.Run("unmanaged conflict preserves complete snapshot", func(t *testing.T) {
		root := t.TempDir()
		if err := runInit(root, true, "default"); err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte("unmanaged\n"), 0644); err != nil {
			t.Fatal(err)
		}
		before, err := filesystemSnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if err := syncProject(newSyncCmd(), root, syncOptions{forced: []string{core.HostOpenCode}, explicitEnable: true}); err == nil {
			t.Fatal("expected unmanaged conflict")
		}
		after, err := filesystemSnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !sameSnapshot(before, after) {
			t.Fatal("conflict changed files or created backup artifacts")
		}
		manifest, err := core.LoadProjectManifest(root)
		if err != nil {
			t.Fatal(err)
		}
		if manifest.Mode != core.ManifestModeNonSubagent || manifest.HostIntent.OpenCode != core.HostDisabled {
			t.Fatalf("conflict changed manifest intent: %#v", manifest)
		}
	})
}

func sameSnapshot(left, right map[string][]byte) bool {
	if len(left) != len(right) {
		return false
	}
	for path, want := range left {
		got, ok := right[path]
		if !ok || !bytes.Equal(want, got) {
			return false
		}
	}
	return true
}
