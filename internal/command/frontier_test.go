package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/tracker"
)

func TestClassifyReadyByUnwaivedDiagnosticFact(t *testing.T) {
	clean := &tracker.Issue{ID: "01", FilePath: "01.md"}
	warned := &tracker.Issue{ID: "02", FilePath: "02.md"}
	gapped := &tracker.Issue{ID: "03", FilePath: "03.md"}
	waived := &tracker.Issue{ID: "04", FilePath: "04.md"}
	legacy := &tracker.Issue{ID: "05", FilePath: "05.md"}
	blocked := &tracker.Issue{ID: "06", FilePath: "06.md"}
	gapThenBlocked := &tracker.Issue{ID: "07", FilePath: "07.md"}
	gapThenWarning := &tracker.Issue{ID: "08", FilePath: "08.md"}
	diagnostics := []tracker.Diagnostic{
		{Code: "risk", Severity: tracker.SeverityWarning, Artifact: warned.FilePath},
		{Code: "missing", Severity: tracker.SeverityGap, Artifact: gapped.FilePath},
		{Code: "reviewed", Severity: tracker.SeverityWarning, Artifact: waived.FilePath, Waiver: &tracker.AppliedWaiver{Reason: "reviewed"}},
		{Code: tracker.DiagnosticLegacyMetadata, Severity: tracker.SeverityWarning, Artifact: legacy.FilePath},
		{Code: "external", Severity: tracker.SeverityBlocker, Artifact: blocked.FilePath},
		{Code: "gap-first", Severity: tracker.SeverityGap, Artifact: gapThenBlocked.FilePath},
		{Code: "blocker-second", Severity: tracker.SeverityBlocker, Artifact: gapThenBlocked.FilePath},
		{Code: "gap-first", Severity: tracker.SeverityGap, Artifact: gapThenWarning.FilePath},
		{Code: "warning-second", Severity: tracker.SeverityWarning, Artifact: gapThenWarning.FilePath},
	}

	gotClean, gotWarnings, gotGaps, gotBlocked := classifyReady([]*tracker.Issue{clean, warned, gapped, waived, legacy, blocked, gapThenBlocked, gapThenWarning}, diagnostics)
	if len(gotClean) != 3 || len(gotWarnings) != 1 || gotWarnings[0] != warned || len(gotGaps) != 2 || gotGaps[0] != gapped || gotGaps[1] != gapThenWarning || len(gotBlocked) != 2 || gotBlocked[0] != blocked || gotBlocked[1] != gapThenBlocked {
		t.Fatalf("unexpected classification: clean=%v warnings=%v gaps=%v", gotClean, gotWarnings, gotGaps)
	}
}

func TestCheckJSONIncludesCatalogWhenThereAreNoTickets(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "design.md"), []byte("---\nflowforge:\n  schema: 1\n---\n# Design\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := newCheckCmd()
	checkDir, checkJSON, checkStrict = root, true, false
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if output := stdout.String(); !strings.Contains(output, `"issues_count": 0`) || !strings.Contains(output, `"diagnostics"`) || !strings.Contains(output, `"artifacts_count": 1`) {
		t.Fatalf("incomplete JSON: %s", output)
	}
}

func TestFrontierPolicyNeverOverridesBlocker(t *testing.T) {
	frontierStrict, frontierIncludeGaps = false, true
	if err := catalogPolicyError([]tracker.Diagnostic{{Severity: tracker.SeverityBlocker}}, false); err == nil {
		t.Fatal("blocker must remain non-executable under include-gaps")
	}
	frontierStrict = true
	if err := catalogPolicyError([]tracker.Diagnostic{{Severity: tracker.SeverityWarning}}, true); err == nil {
		t.Fatal("strict warning must return non-zero policy result")
	}
}

func TestStrictTakesPrecedenceOverGapOverride(t *testing.T) {
	clean := []*tracker.Issue{{ID: "01"}}
	warnings := []*tracker.Issue{{ID: "02"}}
	gaps := []*tracker.Issue{{ID: "03"}}
	ready := effectiveReady(clean, warnings, gaps, true, true)
	if len(ready) != 1 || ready[0].ID != "01" {
		t.Fatalf("strict leaked non-clean work: %#v", ready)
	}
}

func TestFrontierQuietStrictAndGapOverride(t *testing.T) {
	root := frontierGapFixture(t)
	run := func(strict, include bool) (string, string, error) {
		cmd := newFrontierCmd()
		frontierDir, frontierQuiet, frontierJSON, frontierStrict, frontierIncludeGaps = root, true, false, strict, include
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		err := cmd.RunE(cmd, nil)
		return stdout.String(), stderr.String(), err
	}
	stdout, stderr, err := run(true, true)
	if err == nil || stdout != "" || !strings.Contains(stderr, "verification-unselected") {
		t.Fatalf("strict override contract failed stdout=%q stderr=%q err=%v", stdout, stderr, err)
	}
	stdout, stderr, err = run(false, true)
	if err != nil || !strings.Contains(stdout, "01-gap.md") || !strings.Contains(stderr, "verification-unselected") {
		t.Fatalf("explicit gap override failed stdout=%q stderr=%q err=%v", stdout, stderr, err)
	}
}

func TestFrontierJSONCarriesAllGroupsAndDiagnostics(t *testing.T) {
	cmd := newFrontierCmd()
	frontierDir, frontierJSON, frontierQuiet, frontierStrict, frontierIncludeGaps = frontierGapFixture(t), true, false, false, false
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"ready"`, `"ready_with_warnings"`, `"gaps"`, `"claimed"`, `"blocked"`, `"diagnostics"`} {
		if !strings.Contains(stdout.String(), key) {
			t.Fatalf("JSON missing %s: %s", key, stdout.String())
		}
	}
}

func frontierGapFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	feature := filepath.Join(root, "feature")
	issues := filepath.Join(feature, "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	design := `---
flowforge:
  schema: 1
  role: design
  areas:
    seam: {revision: 1, anchor: seam}
  open_items:
    - {id: verification-unselected, diagnostic: verification-unselected, severity: gap, affects: [01], anchor: seam}
---
<a id="seam"></a>
`
	ticket := `---
flowforge:
  schema: 1
  role: ticket
---
# 01: Gap
**Status:** open
**Blocked by:** None
`
	if err := os.WriteFile(filepath.Join(feature, "design.md"), []byte(design), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(issues, "01-gap.md"), []byte(ticket), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
