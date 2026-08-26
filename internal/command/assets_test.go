package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssetsVerifyReportsCurrentAndDriftWithoutMutation(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".flowforge"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".flowforge", "config.yaml"), []byte("docs_dir: wiki\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := deployManagedAssets(project, filepath.Join(project, "wiki")); err != nil {
		t.Fatal(err)
	}
	run := func(json bool) (string, error) {
		cmd := newAssetsCmd()
		cmd.SetArgs([]string{"verify", "--json", project})
		if !json {
			cmd.SetArgs([]string{"verify", project})
		}
		var stdout bytes.Buffer
		cmd.SetOut(&stdout)
		err := cmd.Execute()
		return stdout.String(), err
	}
	output, err := run(true)
	if err != nil || !strings.Contains(output, `"current": true`) {
		t.Fatalf("current project verification failed: output=%s err=%v", output, err)
	}
	owned := filepath.Join(project, ".agents", "skills", "project-only.md")
	if err := os.WriteFile(owned, []byte("owned"), 0o644); err != nil {
		t.Fatal(err)
	}
	output, err = run(true)
	if err != nil || !strings.Contains(output, `"current": true`) || !strings.Contains(output, `"state": "project-owned"`) {
		t.Fatalf("project-owned-only verification failed: output=%s err=%v", output, err)
	}
	drifted := filepath.Join(project, ".agents", "skills", "flowforge-align", "SKILL.md")
	if err := os.WriteFile(drifted, []byte("custom"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(project, ".agents", "skills", "flowforge-route", "SKILL.md")
	if err := os.Remove(missing); err != nil {
		t.Fatal(err)
	}
	output, err = run(true)
	if err == nil || !strings.Contains(output, `"state": "missing"`) || !strings.Contains(output, `"state": "drifted"`) || !strings.Contains(output, `"state": "project-owned"`) {
		t.Fatalf("drift/project-owned projection failed: output=%s err=%v", output, err)
	}
	output, err = run(false)
	if err == nil || !strings.Contains(output, "missing "+missing) || !strings.Contains(output, "drifted "+drifted) || !strings.Contains(output, "project-owned "+owned) {
		t.Fatalf("human projection failed: output=%s err=%v", output, err)
	}
	data, readErr := os.ReadFile(owned)
	if readErr != nil || string(data) != "owned" {
		t.Fatalf("verify mutated project-owned asset: %q, %v", data, readErr)
	}
}
