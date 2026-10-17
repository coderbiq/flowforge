package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// checkOutputFixture builds a notebook whose diagnostics exercise every
// convergence path: a standards blocker and a clean ticket in featureA, an
// execution-contract gap in featureB.
func checkOutputFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	fm := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	execContract := "\n\n## Execution detail\n\n### Verified contracts\n\n- check.go owns validity.\n\n### Execution scenarios\n\n- Success: summary carries counts.\n- Failure: exit code unchanged.\n\n### Expected tests\n\n- go test ./internal/command/...\n\n### Generated artifacts\n\n- Not applicable.\n\n### Conventions\n\n- Reuse output_budget helpers.\n"
	write := func(feature, name, constraints string) {
		issues := filepath.Join(root, feature, "issues")
		if err := os.MkdirAll(issues, 0o755); err != nil {
			t.Fatal(err)
		}
		body := fm + "# " + strings.TrimSuffix(name, ".md") + "\n**Status:** open\n**Blocked by:** None\n" + execContract
		if constraints != "" {
			body += "\n## Constraints\n\n" + constraints
		}
		if err := os.WriteFile(filepath.Join(issues, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("featureA", "01-clean.md", "- must keep behavior\n- Write set: app/\n")
	write("featureA", "02-gated.md", "- Write set: app/\n") // Write set without a standards marker -> blocker
	noContract := filepath.Join(root, "featureB", "issues")
	if err := os.MkdirAll(noContract, 0o755); err != nil {
		t.Fatal(err)
	}
	broken := fm + "# 03-noisy.md\n**Status:** open\n**Blocked by:** None\n"
	if err := os.WriteFile(filepath.Join(noContract, "03-noisy.md"), []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func runCheckCaptured(t *testing.T, root string, jsonMode, full bool, feature string) (string, string, error) {
	t.Helper()
	cmd := newCheckCmd()
	checkDir, checkJSON, checkStrict, checkFeature, checkFull = root, jsonMode, false, feature, full
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	err := cmd.RunE(cmd, nil)
	return stdout.String(), stderr.String(), err
}

func TestCheckTextModeSummary(t *testing.T) {
	stdout, stderr, err := runCheckCaptured(t, checkOutputFixture(t), false, false, "")
	if err == nil {
		t.Fatal("blocker fixture must exit non-zero after convergence")
	}
	for _, want := range []string{
		"diagnostics blocker:standards-marker-missing: 1",
		"diagnostics gap:execution-contract-incomplete: 1",
		"--full for all diagnostics",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("summary missing %q, stderr=%q stdout=%q", want, stderr, stdout)
		}
	}

	_, stderrFull, errFull := runCheckCaptured(t, checkOutputFixture(t), false, true, "")
	if errFull == nil {
		t.Fatal("--full must keep the blocker exit code")
	}
	if !strings.Contains(stderrFull, "standards-marker-missing") || !strings.Contains(stderrFull, " — ") {
		t.Fatalf("--full must print verbatim diagnostics with messages, stderr=%q", stderrFull)
	}
}

func TestCheckJSONSlim(t *testing.T) {
	stdout, _, err := runCheckCaptured(t, checkOutputFixture(t), true, false, "")
	if err == nil {
		t.Fatal("blocker fixture must exit non-zero in JSON mode")
	}
	for _, want := range []string{`"valid":false`, `"counts"`, `"blockers"`, `"detail"`} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("slim JSON missing %s: %s", want, stdout)
		}
	}
	if strings.Contains(stdout, "\n  \"") {
		t.Fatalf("default JSON must marshal compact: %s", stdout[:120])
	}

	stdoutFull, _, errFull := runCheckCaptured(t, checkOutputFixture(t), true, true, "")
	if errFull == nil {
		t.Fatal("--full must keep the blocker exit code")
	}
	if !strings.Contains(stdoutFull, "\n  \"") || !strings.Contains(stdoutFull, `"message"`) {
		t.Fatalf("--full JSON must stay indented with per-diagnostic messages: %s", stdoutFull[:120])
	}
}

func TestCheckFeatureScoping(t *testing.T) {
	stdout, _, err := runCheckCaptured(t, checkOutputFixture(t), true, false, "featureA")
	if err == nil {
		t.Fatal("scoped view must not relax validity: the global blocker keeps the exit non-zero")
	}
	if !strings.Contains(stdout, `"valid":false`) {
		t.Fatalf("scoped JSON must report validity, stdout=%q", stdout)
	}
	if strings.Contains(stdout, "execution-contract-incomplete") {
		t.Fatalf("scoped diagnostics must exclude featureB gaps: %s", stdout)
	}
	if !strings.Contains(stdout, "standards-marker-missing") {
		t.Fatalf("scoped diagnostics must keep featureA blockers: %s", stdout)
	}
	if !strings.Contains(stdout, `"issues_count":3`) {
		t.Fatalf("issues_count stays global, stdout=%q", stdout)
	}

	textStdout, _, _ := runCheckCaptured(t, checkOutputFixture(t), false, false, "featureA")
	if !strings.Contains(textStdout, "(diagnostics scoped to featureA; validity computed globally)") {
		t.Fatalf("text mode must announce the scope, stdout=%q", textStdout)
	}
}
