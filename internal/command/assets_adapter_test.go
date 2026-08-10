package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCodeAdapterInstallAndPreserveConflict(t *testing.T) {
	root := t.TempDir()
	cmd := newAssetsAdapterInstallCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := applyOpenCodeAdapter(cmd, root, false); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".opencode", "agents", "flowforge-executor.md")
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("adapter target missing: %v", err)
	}
	if err := os.WriteFile(target, []byte("user change"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := applyOpenCodeAdapter(cmd, root, false); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "user change" {
		t.Fatalf("user change overwritten: %s", data)
	}
}

func TestOpenCodeAdapterUpdateIsIdempotent(t *testing.T) {
	root := t.TempDir()
	cmd := newAssetsAdapterInstallCmd()
	var first bytes.Buffer
	cmd.SetOut(&first)
	if err := applyOpenCodeAdapter(cmd, root, false); err != nil {
		t.Fatal(err)
	}
	var second bytes.Buffer
	cmd.SetOut(&second)
	if err := applyOpenCodeAdapter(cmd, root, true); err != nil {
		t.Fatal(err)
	}
	if second.Len() != 0 {
		t.Fatalf("idempotent update should not rewrite files:\n%s", second.String())
	}
}
