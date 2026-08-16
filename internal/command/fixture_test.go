package command

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type AGENTSFixture struct {
	Root     string
	Path     string
	Body     []byte
	SHA256   string
	Sentinel UnmanagedSentinel
}

type UnmanagedSentinel struct {
	Path string
	Body []byte
}

func newAGENTSFixture(t *testing.T) AGENTSFixture {
	t.Helper()
	root := t.TempDir()
	agents := []byte("# base FLOWFORGE\n\n<!-- FLOWFORGE:START -->\nmanaged orchestration\n<!-- FLOWFORGE:END -->\n\n# user instructions\nkeep this\n\n<!-- OTHER-TOOL:START -->\nother tool\n<!-- OTHER-TOOL:END -->\n")
	path := filepath.Join(root, "AGENTS.md")
	if err := os.WriteFile(path, agents, 0644); err != nil {
		t.Fatal(err)
	}
	sentinel := UnmanagedSentinel{Path: filepath.Join(root, ".opencode", "agents", "unmanaged.md"), Body: []byte("do not manage\n")}
	if err := os.MkdirAll(filepath.Dir(sentinel.Path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sentinel.Path, sentinel.Body, 0644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(agents)
	return AGENTSFixture{Root: root, Path: path, Body: agents, SHA256: fmt.Sprintf("%x", sum), Sentinel: sentinel}
}

func TestAGENTSFixtureSelfCheck(t *testing.T) {
	fixture := newAGENTSFixture(t)
	got, err := os.ReadFile(fixture.Path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(got)
	if string(got) != string(fixture.Body) || fmt.Sprintf("%x", sum) != fixture.SHA256 {
		t.Fatal("AGENTS fixture baseline is not reproducible")
	}
	sentinel, err := os.ReadFile(fixture.Sentinel.Path)
	if err != nil {
		t.Fatal(err)
	}
	if string(sentinel) != string(fixture.Sentinel.Body) {
		t.Fatal("unmanaged sentinel drifted")
	}
}
