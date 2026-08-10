package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenCodeAdapterStaticSmoke(t *testing.T) {
	root := t.TempDir()
	cmd := newAssetsAdapterInstallCmd()
	if err := applyOpenCodeAdapter(cmd, root, false); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"flowforge-design-analyst.md", "flowforge-executor.md"} {
		data, err := os.ReadFile(filepath.Join(root, ".opencode", "agents", name))
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"mode: subagent", "task: deny", "question: deny", "Default Skill:"} {
			if !strings.Contains(string(data), want) {
				t.Fatalf("%s missing %q", name, want)
			}
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".opencode", "opencode.json")); !os.IsNotExist(err) {
		t.Fatalf("adapter must not create OpenCode config: %v", err)
	}
}
