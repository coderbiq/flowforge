package command

import (
	"bytes"
	"encoding/json"
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
	checkDir, checkJSON, checkStrict, checkFull = root, true, false, true
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	checkFull = false
	if output := stdout.String(); !strings.Contains(output, `"issues_count": 0`) || !strings.Contains(output, `"diagnostics"`) || !strings.Contains(output, `"artifacts_count": 1`) {
		t.Fatalf("incomplete JSON: %s", output)
	}
}

func TestStrictCheckValidatesRequirementAndDesignWithoutTickets(t *testing.T) {
	root := t.TempDir()
	previousDir, previousJSON, previousStrict, previousFull := checkDir, checkJSON, checkStrict, checkFull
	defer func() { checkDir, checkJSON, checkStrict, checkFull = previousDir, previousJSON, previousStrict, previousFull }()
	requirements := `---
flowforge:
  schema: 1
  role: requirement
  id: intake-requirements
  revision: 1
---
<a id="intake-requirements"></a>
# Requirements
`
	design := `---
flowforge:
  schema: 1
  role: design
  id: intake-design
  revision: 1
  areas:
    publication: {revision: 1, anchor: publication}
  consumes:
    requirements:
      intake-requirements: 1
---
<a id="intake-design"></a>
<a id="publication"></a>
# Design

See [requirements](requirements.md#intake-requirements).
`
	if err := os.WriteFile(filepath.Join(root, "requirements.md"), []byte(requirements), 0o644); err != nil {
		t.Fatal(err)
	}
	run := func(body string, strict, json bool) (string, error) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, "design.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		cmd := newCheckCmd()
		checkDir, checkJSON, checkStrict, checkFull = root, json, strict, true
		var stdout bytes.Buffer
		cmd.SetOut(&stdout)
		err := cmd.RunE(cmd, nil)
		return stdout.String(), err
	}
	if _, err := run(design, true, false); err != nil {
		t.Fatalf("valid authority-only feature rejected: %v", err)
	}

	invalid := map[tracker.DiagnosticCode]string{
		tracker.DiagnosticMissingAnchor: strings.Replace(design, `anchor: publication`, `anchor: missing-publication`, 1),
		tracker.DiagnosticInvalidOpenItem: strings.Replace(design, `  consumes:`, `  open_items:
    - {id: unresolved-publication, diagnostic: publication-gap, severity: gap, affects: [publication], anchor: missing-open-item}
  consumes:`, 1),
		tracker.DiagnosticFutureRevision:   strings.Replace(design, `intake-requirements: 1`, `intake-requirements: 2`, 1),
		tracker.DiagnosticMissingHumanLink: strings.Replace(design, `requirements.md#intake-requirements`, `unlinked.md#intake-requirements`, 1),
	}
	for diagnostic, body := range invalid {
		output, err := run(body, false, true)
		if err != nil {
			t.Fatalf("non-strict check rejected %s: %v", diagnostic, err)
		}
		if !strings.Contains(output, `"issues_count": 0`) || !strings.Contains(output, `"code": "`+string(diagnostic)+`"`) {
			t.Fatalf("non-strict JSON omitted authority diagnostic %s: %s", diagnostic, output)
		}
		output, err = run(body, true, true)
		if err == nil || !strings.Contains(output, `"code": "`+string(diagnostic)+`"`) {
			t.Fatalf("strict check accepted or hid authority diagnostic %s: output=%s err=%v", diagnostic, output, err)
		}
	}
}

func TestCheckReportsClosedTicketWithoutCompletionEvidence(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	ticket := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 01: Closed too early\n\n**Blocked by:** None\n**Status:** closed\n"
	if err := os.WriteFile(filepath.Join(issues, "01-closed.md"), []byte(ticket), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := newCheckCmd()
	checkDir, checkJSON, checkStrict, checkFull = root, true, false, true
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	checkFull = false
	if output := stdout.String(); !strings.Contains(output, `"code": "missing-completion-evidence"`) {
		t.Fatalf("check omitted completion evidence diagnostic: %s", output)
	}

	strict := newCheckCmd()
	checkDir, checkJSON, checkStrict = root, false, true
	strict.SetOut(&bytes.Buffer{})
	strict.SetErr(&bytes.Buffer{})
	if err := strict.RunE(strict, nil); err == nil {
		t.Fatal("strict check must reject a closed ticket without completion evidence")
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

func TestFrontierExcludesIncompleteExecutionContractsUnlessGapsAreIncluded(t *testing.T) {
	root := executionContractFrontierFixture(t)
	run := func(include bool) (string, string, error) {
		cmd := newFrontierCmd()
		frontierDir, frontierQuiet, frontierJSON, frontierStrict, frontierIncludeGaps = root, true, false, false, include
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		err := cmd.RunE(cmd, nil)
		return stdout.String(), stderr.String(), err
	}

	stdout, stderr, err := run(false)
	if err != nil || !strings.Contains(stdout, "01-complete.md") || strings.Contains(stdout, "02-incomplete.md") || !strings.Contains(stderr, "execution-contract-incomplete") {
		t.Fatalf("default execution-contract projection failed stdout=%q stderr=%q err=%v", stdout, stderr, err)
	}
	stdout, stderr, err = run(true)
	if err != nil || !strings.Contains(stdout, "01-complete.md") || !strings.Contains(stdout, "02-incomplete.md") || !strings.Contains(stderr, "execution-contract-incomplete") {
		t.Fatalf("include-gaps execution-contract projection failed stdout=%q stderr=%q err=%v", stdout, stderr, err)
	}
}

func TestFrontierJSONCarriesAllGroupsAndDiagnostics(t *testing.T) {
	cmd := newFrontierCmd()
	frontierDir, frontierJSON, frontierQuiet, frontierStrict, frontierIncludeGaps, frontierFeature, frontierFull = frontierGapFixture(t), true, false, false, false, "", false
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	out := stdout.String()
	for _, key := range []string{`"ready"`, `"ready_with_warnings"`, `"gaps"`, `"claimed"`, `"blocked"`, `"diagnostics"`, `"counts"`, `"detail"`} {
		if !strings.Contains(out, key) {
			t.Fatalf("JSON missing %s: %s", key, out)
		}
	}
	if strings.Contains(out, `"body"`) {
		t.Fatalf("default JSON must not carry ticket bodies: %s", out)
	}

	cmdFull := newFrontierCmd()
	// Flag registration resets package vars to defaults; re-assert both modes.
	frontierJSON, frontierFull = true, true
	var stdoutFull bytes.Buffer
	cmdFull.SetOut(&stdoutFull)
	if err := cmdFull.RunE(cmdFull, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdoutFull.String(), `"body"`) {
		t.Fatalf("--full JSON must carry ticket bodies: %s", stdoutFull.String())
	}
	frontierFull = false
}

func featureProjectionFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	fm := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	execContract := "\n\n## Execution detail\n\n### Verified contracts\n\n- frontier.go owns the frontier projection.\n\n### Execution scenarios\n\n- Success: scoped output lists only the feature's tickets.\n- Failure: cross-feature references stay visible.\n\n### Expected tests\n\n- go test ./internal/command/...\n\n### Generated artifacts\n\n- Not applicable.\n\n### Conventions\n\n- Reuse existing frontier logic.\n"

	write := func(feature, name, blockedBy string) {
		issues := filepath.Join(root, feature, "issues")
		if err := os.MkdirAll(issues, 0o755); err != nil {
			t.Fatal(err)
		}
		body := fm + "# " + strings.TrimSuffix(name, ".md") + "\n**Status:** open\n**Blocked by:** " + blockedBy + "\n" + execContract
		if err := os.WriteFile(filepath.Join(issues, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("featureA", "01-ready.md", "None")
	write("featureA", "02-waiting.md", "01")
	write("featureB", "03-other.md", "None")
	return root
}

func TestFrontierFeatureProjection(t *testing.T) {
	cmd := newFrontierCmd()
	frontierDir, frontierJSON, frontierQuiet, frontierStrict, frontierIncludeGaps, frontierFeature, frontierFull = featureProjectionFixture(t), true, false, false, false, "featureA", false
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("frontier failed: %v", err)
	}
	var result struct {
		Ready []struct {
			ID      string `json:"id"`
			Feature string `json:"feature"`
		} `json:"ready"`
		Blocked []struct {
			ID        string   `json:"id"`
			WaitingOn []string `json:"waiting_on"`
		} `json:"blocked"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v stdout=%q", err, stdout.String())
	}
	for _, entry := range result.Ready {
		if entry.Feature != "featureA" {
			t.Fatalf("scoped ready must contain only featureA tickets, got feature %q", entry.Feature)
		}
		if entry.ID == "03" {
			t.Fatal("scoped ready must not contain featureB tickets")
		}
	}
	foundWaiting := false
	for _, b := range result.Blocked {
		if b.ID == "02" {
			foundWaiting = true
			if len(b.WaitingOn) == 0 || !strings.Contains(b.WaitingOn[0], "01") {
				t.Fatalf("blocked 02 must keep its waiting_on reference to 01, got %v", b.WaitingOn)
			}
		}
		if b.ID == "03" {
			t.Fatal("scoped blocked must not contain featureB tickets")
		}
	}
	if !foundWaiting {
		t.Fatal("in-feature blocked ticket 02 must appear in scoped output")
	}
	if len(stdout.String()) > MaxAgentOutputBytes {
		t.Fatalf("scoped JSON exceeds MaxAgentOutputBytes: %d bytes", len(stdout.String()))
	}
}

func TestFrontierTextModeSummary(t *testing.T) {
	// A ready ticket keeps the text path past its empty-queue early exit; a
	// contract-less ticket in another feature supplies gap diagnostics.
	root := featureProjectionFixture(t)
	brokenDir := filepath.Join(root, "featureB", "issues")
	if err := os.MkdirAll(brokenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	broken := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 04-broken\n**Status:** open\n**Blocked by:** None\n"
	if err := os.WriteFile(filepath.Join(brokenDir, "04-broken.md"), []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := newFrontierCmd()
	frontierDir, frontierJSON, frontierQuiet, frontierStrict, frontierIncludeGaps, frontierFeature, frontierFull = root, false, false, false, false, "", false
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("frontier failed: %v", err)
	}
	errText := stderr.String()
	if !strings.Contains(errText, "diagnostics gap:execution-contract-incomplete:") {
		t.Fatalf("summary must carry a count line, stderr=%q", errText)
	}
	if !strings.Contains(errText, "--full for all diagnostics") {
		t.Fatalf("summary must carry the --full hint, stderr=%q", errText)
	}
	if strings.Contains(errText, " — ") {
		t.Fatalf("default text mode must not print verbatim diagnostic messages, stderr=%q", errText)
	}

	cmdFull := newFrontierCmd()
	frontierFull = true
	var stdoutFull, stderrFull bytes.Buffer
	cmdFull.SetOut(&stdoutFull)
	cmdFull.SetErr(&stderrFull)
	if err := cmdFull.RunE(cmdFull, nil); err != nil {
		t.Fatalf("frontier --full failed: %v", err)
	}
	if !strings.Contains(stderrFull.String(), " — ") {
		t.Fatalf("--full must print verbatim diagnostics with messages, stderr=%q", stderrFull.String())
	}
	frontierFull = false
}

func TestBlockedEvidenceFrontierClassification(t *testing.T) {
	clean := &tracker.Issue{ID: "01", FilePath: "01-clean.md"}
	blockedEvidence := &tracker.Issue{ID: "02", FilePath: "02-blocked.md"}
	diagnostics := []tracker.Diagnostic{
		{Code: tracker.DiagnosticBlockedEvidencePresent, Severity: tracker.SeverityWarning, Artifact: blockedEvidence.FilePath},
	}
	gotClean, gotWarnings, gotGaps, gotBlocked := classifyReady([]*tracker.Issue{clean, blockedEvidence}, diagnostics)
	if len(gotClean) != 1 || gotClean[0] != clean || len(gotWarnings) != 1 || gotWarnings[0] != blockedEvidence || len(gotGaps) != 0 || len(gotBlocked) != 0 {
		t.Fatalf("blocked-evidence warning must land in warnings bucket, not clean ready: clean=%v warnings=%v gaps=%v blocked=%v", gotClean, gotWarnings, gotGaps, gotBlocked)
	}

	root := blockedEvidenceFrontierFixture(t)
	ticketPath := filepath.Join(root, "feature", "issues", "01-blocked.md")
	original, err := os.ReadFile(ticketPath)
	if err != nil {
		t.Fatal(err)
	}
	restore := func() {
		if err := os.WriteFile(ticketPath, original, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	cmd := newFrontierCmd()
	frontierDir, frontierJSON, frontierQuiet, frontierStrict, frontierIncludeGaps = root, true, false, false, false
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("frontier failed: %v stdout=%q", err, stdout.String())
	}
	var result struct {
		Ready            []*tracker.Issue `json:"ready"`
		ReadyWithWarning []*tracker.Issue `json:"ready_with_warnings"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v stdout=%q", err, stdout.String())
	}
	for _, issue := range result.Ready {
		if issue.ID == "01" {
			t.Fatal("blocked-evidence ticket must not enter clean ready")
		}
	}
	var inWarnings bool
	for _, issue := range result.ReadyWithWarning {
		if issue.ID == "01" {
			inWarnings = true
		}
	}
	if !inWarnings {
		t.Fatalf("blocked-evidence ticket missing from ready_with_warnings: %s", stdout.String())
	}

	runFrontier := func(strict bool) (string, string, error) {
		cmd := newFrontierCmd()
		frontierDir, frontierQuiet, frontierJSON, frontierStrict, frontierIncludeGaps = root, true, false, strict, false
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		err := cmd.RunE(cmd, nil)
		return stdout.String(), stderr.String(), err
	}
	stdoutText, stderrText, err := runFrontier(false)
	if err != nil || !strings.Contains(stdoutText, "01-blocked.md") || !strings.Contains(stderrText, "blocked-evidence-present") {
		t.Fatalf("default frontier projection failed stdout=%q stderr=%q err=%v", stdoutText, stderrText, err)
	}
	stdoutText, stderrText, err = runFrontier(true)
	if err == nil || strings.Contains(stdoutText, "01-blocked.md") || !strings.Contains(stderrText, "blocked-evidence-present") {
		t.Fatalf("strict frontier must exclude the blocked-evidence ticket stdout=%q stderr=%q err=%v", stdoutText, stderrText, err)
	}

	runCheck := func(strict bool) (string, error) {
		cmd := newCheckCmd()
		checkDir, checkJSON, checkStrict = root, false, strict
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		err := cmd.RunE(cmd, nil)
		return stderr.String(), err
	}
	stderrText, err = runCheck(false)
	if err != nil || !strings.Contains(stderrText, "blocked-evidence-present") {
		t.Fatalf("default check must print the diagnostic and stay valid stderr=%q err=%v", stderrText, err)
	}
	stderrText, err = runCheck(true)
	if err == nil || !strings.Contains(stderrText, "blocked-evidence-present") {
		t.Fatalf("strict check must fail on blocked-evidence warning stderr=%q err=%v", stderrText, err)
	}

	closed := strings.Replace(string(original), "**Status:** open", "**Status:** closed\n\n## Completion evidence\n\n- go test ./internal/... passed.", 1)
	if err := os.WriteFile(ticketPath, []byte(closed), 0o644); err != nil {
		t.Fatal(err)
	}
	if stderrText, err = runCheck(true); err != nil || strings.Contains(stderrText, "blocked-evidence-present") {
		t.Fatalf("closed ticket must not report blocked evidence stderr=%q err=%v", stderrText, err)
	}

	consumed := strings.Replace(string(original), "\n## Blocked evidence\n\n- Verbatim: exit status 1\n- cmd: go test ./... (exit 1)\n- Next: inspect the failing package\n", "", 1)
	if err := os.WriteFile(ticketPath, []byte(consumed), 0o644); err != nil {
		t.Fatal(err)
	}
	if stderrText, err = runCheck(true); err != nil || strings.Contains(stderrText, "blocked-evidence-present") {
		t.Fatalf("consumed section must clear the diagnostic stderr=%q err=%v", stderrText, err)
	}
	restore()
}

func blockedEvidenceFrontierFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	execContract := "## Execution detail\n\n### Verified contracts\n\n- catalog.go owns ticket diagnostics.\n\n### Execution scenarios\n\n- Success: open ticket reports blocked evidence.\n- Failure: closed or consumed tickets stay silent.\n\n### Expected tests\n\n- go test ./internal/command/...\n\n### Generated artifacts\n\n- Not applicable — diagnostics are computed at runtime.\n\n### Conventions\n\n- Reuse existing warning projection.\n"
	blockedSection := "\n## Blocked evidence\n\n- Verbatim: exit status 1\n- cmd: go test ./... (exit 1)\n- Next: inspect the failing package\n"
	ticket := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 01: Blocked\n**Status:** open\n**Blocked by:** None\n\n" + execContract + blockedSection
	if err := os.WriteFile(filepath.Join(issues, "01-blocked.md"), []byte(ticket), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
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

func executionContractFrontierFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	complete := "## Execution detail\n\n### Verified contracts\n\n- catalog.go supplies diagnostics.\n\n### Execution scenarios\n\n- Success: ticket is ready.\n- Failure: missing sections are gaps.\n\n### Expected tests\n\n- go test ./internal/tracker/...\n\n### Generated artifacts\n\n- Not applicable — no generated output.\n\n### Conventions\n\n- Reuse existing gap projection.\n"
	write := func(name, detail string) {
		t.Helper()
		body := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# " + name + "\n**Status:** open\n**Blocked by:** None\n\n" + detail
		if err := os.WriteFile(filepath.Join(issues, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("01-complete.md", complete)
	write("02-incomplete.md", "## Execution detail\n\n### Verified contracts\n\n- TODO\n")
	return root
}

func TestFrontierExcludesNeedsRepairAndDispatchesRepair(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	execContract := "\n\n## Execution detail\n\n### Verified contracts\n\n- parser.go owns header parsing.\n\n### Execution scenarios\n\n- Success: repair is dispatched.\n- Failure: needs-repair is excluded.\n\n### Expected tests\n\n- go test ./internal/command/...\n\n### Generated artifacts\n\n- Not applicable.\n\n### Conventions\n\n- Reuse existing frontier logic.\n"
	writeTicket := func(name, status, blockedBy, repairOf string) {
		header := "# " + strings.TrimSuffix(name, ".md") + "\n**Status:** " + status + "\n"
		if repairOf != "" {
			header += "**Repair of:** " + repairOf + "\n"
		}
		header += "**Blocked by:** " + blockedBy + "\n"
		body := fm + header + execContract
		if err := os.WriteFile(filepath.Join(issues, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeTicket("01-original.md", "needs-repair", "None", "")
	writeTicket("02-repair.md", "open", "None", "01")
	writeTicket("03-downstream.md", "open", "02", "")

	cmd := newFrontierCmd()
	frontierDir, frontierJSON, frontierQuiet, frontierStrict, frontierIncludeGaps = root, true, false, false, false
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("frontier failed: %v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
	var result struct {
		Ready   []*tracker.Issue `json:"ready"`
		Blocked []struct {
			ID        string   `json:"id"`
			WaitingOn []string `json:"waiting_on"`
		} `json:"blocked"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v stdout=%q", err, stdout.String())
	}
	for _, issue := range result.Ready {
		if issue.ID == "01" {
			t.Fatal("needs-repair ticket must not be in ready")
		}
	}
	var repairReady bool
	for _, issue := range result.Ready {
		if issue.ID == "02" {
			repairReady = true
		}
	}
	if !repairReady {
		t.Fatal("repair ticket must be in ready")
	}
	var downstreamBlocked bool
	for _, b := range result.Blocked {
		if b.ID == "03" {
			downstreamBlocked = true
			found := false
			for _, w := range b.WaitingOn {
				if w == "02" {
					found = true
				}
			}
			if !found {
				t.Fatalf("downstream 03 should wait on 02, got %v", b.WaitingOn)
			}
		}
	}
	if !downstreamBlocked {
		t.Fatal("downstream ticket 03 should be blocked waiting on repair 02")
	}
}

func TestRenderPiWorkflow(t *testing.T) {
	// Titles/paths carry quotes, backticks, and newlines to prove every dynamic
	// string lands in JSON-encoded form and cannot break the script syntax.
	nastyTitle := "quote \" backtick ` newline\nand 'single'"
	nastyPath := "docs/proposals/x\"y/issues/01-a`b.md"
	ready := []*tracker.Issue{
		{ID: "01", FilePath: nastyPath, Title: nastyTitle},
		{ID: "02", FilePath: "docs/proposals/demo/issues/02-clean.md", Title: "Clean title"},
	}
	out := renderPiWorkflow(ready, "docs/proposals")

	if !strings.HasPrefix(out, "const tickets = [") {
		t.Fatalf("output must start with the tickets array, got: %q", out)
	}
	if strings.Contains(out, "function") || strings.Contains(out, "=>") {
		t.Fatalf("workflowScript must not contain function or arrow tokens: %q", out)
	}
	if got := strings.Count(out, "runs.run("); got != len(ready) {
		t.Fatalf("expected exactly %d runs.run calls, got %d", len(ready), got)
	}
	if got := strings.Count(out, `agent: "flowforge-implementer"`); got != len(ready) {
		t.Fatalf("expected each runs.run to pin the flowforge-implementer agent, got %d occurrences", got)
	}
	// context: "fresh" must be pinned on every launch: omitting it lets
	// pi-subagents global defaultSubagentContext/defaultContext settings fork
	// the parent session instead of guaranteeing a fresh execution context.
	if got := strings.Count(out, `context: "fresh"`); got != len(ready) {
		t.Fatalf("expected each runs.run to pin context: \"fresh\", got %d occurrences", got)
	}
	if got := strings.Count(out, `gate: { command: "flowforge check --dir docs/proposals" }`); got != len(ready) {
		t.Fatalf("expected each runs.run to carry the structural gate, got %d occurrences", got)
	}
	// Nasty strings must appear exactly in json.Marshal encoded form.
	for _, s := range []string{nastyTitle, nastyPath, "01", "02"} {
		encoded, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, string(encoded)) {
			t.Fatalf("expected encoded form %s of %q to appear verbatim in output:\n%s", encoded, s, out)
		}
	}
	// A raw newline inside a string literal would break JS syntax; the encoded
	// title must sit on one line, so line count matches structural newlines only.
	if !strings.Contains(out, "newline\\nand") {
		t.Fatalf("newline in title must be encoded as \\n, got:\n%s", out)
	}
	// Guidance bracketing: static prefix/suffix concatenated with the indexed path.
	if !strings.Contains(out, jsString(piWorkflowTaskPrefix)+" + tickets[0].path + "+jsString(piWorkflowTaskSuffix)) {
		t.Fatalf("task must compose prefix + ticket path + suffix, got:\n%s", out)
	}
	// Fail-fast return and final success return are explicit.
	if !strings.Contains(out, "return \"STOPPED at ticket \" + tickets[0].key") {
		t.Fatalf("missing fail-fast STOPPED return:\n%s", out)
	}
	if !strings.HasSuffix(strings.TrimSpace(out), "ticket(s) dispatched and verified.\";") {
		t.Fatalf("missing explicit success return:\n%s", out)
	}
	// Idempotent rendering.
	if again := renderPiWorkflow(ready, "docs/proposals"); again != out {
		t.Fatal("renderPiWorkflow must be deterministic for identical input")
	}
}

func TestRenderPiWorkflowEmptyBatch(t *testing.T) {
	out := renderPiWorkflow(nil, "docs/proposals")
	if !strings.HasPrefix(out, "return 'No ready tickets") {
		t.Fatalf("empty batch must start with the no-ready return, got: %q", out)
	}
	if strings.Contains(out, "for") || strings.Contains(out, "runs.run") {
		t.Fatalf("empty batch must not emit the loop skeleton, got: %q", out)
	}
}

func standardsGateFrontierFixture(t *testing.T, gatedConstraints string) string {
	t.Helper()
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	execContract := "## Execution detail\n\n### Verified contracts\n\n- catalog.go owns ticket diagnostics.\n\n### Execution scenarios\n\n- Success: marked ticket stays ready.\n- Failure: unmarked ticket leaves the queue.\n\n### Expected tests\n\n- go test ./internal/command/...\n\n### Generated artifacts\n\n- Not applicable — diagnostics are computed at runtime.\n\n### Conventions\n\n- Reuse existing blocker projection.\n"
	write := func(name, constraints string) {
		t.Helper()
		body := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# " + name + "\n**Status:** open\n**Blocked by:** None\n\n" + execContract + "\n## Constraints\n\n" + constraints + "\n"
		if err := os.WriteFile(filepath.Join(issues, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("01-marked.md", "- must keep behavior\n- Write set: internal/tracker/\n")
	write("02-gated.md", gatedConstraints)
	return root
}

func runStandardsGateFrontier(t *testing.T, root string) (string, string, error) {
	t.Helper()
	cmd := newFrontierCmd()
	frontierDir, frontierQuiet, frontierJSON, frontierStrict, frontierIncludeGaps = root, true, false, false, false
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	err := cmd.RunE(cmd, nil)
	return stdout.String(), stderr.String(), err
}

// TestFrontierExcludesStandardsMarkerMissingTickets pins the blocker arm: a
// Write-set ticket without any standards marker leaves the dispatch queue and
// frontier prints the repair guidance on stderr.
func TestFrontierExcludesStandardsMarkerMissingTickets(t *testing.T) {
	root := standardsGateFrontierFixture(t, "- keep the change small\n- Write set: internal/tracker/\n")
	stdoutText, stderrText, err := runStandardsGateFrontier(t, root)
	if err != errPolicyViolation {
		t.Fatalf("blocker must fail the frontier policy: err=%v stdout=%q stderr=%q", err, stdoutText, stderrText)
	}
	if !strings.Contains(stdoutText, "01-marked.md") || strings.Contains(stdoutText, "02-gated.md") {
		t.Fatalf("missing-marker ticket must leave the ready queue: stdout=%q", stdoutText)
	}
	if !strings.Contains(stderrText, string(tracker.DiagnosticStandardsMarkerMissing)) || !strings.Contains(stderrText, "return to Plan and rerun standards extraction") {
		t.Fatalf("stderr must carry the standards gate guidance: stderr=%q", stderrText)
	}
}

// TestFrontierExcludesStandardsPendingTickets pins the warning arm: a pending
// ticket is escalated out of ready (the implementer preflight would only
// bounce it) while the diagnostic itself stays a warning, not a blocker.
func TestFrontierExcludesStandardsPendingTickets(t *testing.T) {
	root := standardsGateFrontierFixture(t, "- standards: pending\n- Write set: internal/tracker/\n")
	stdoutText, stderrText, err := runStandardsGateFrontier(t, root)
	if err != nil {
		t.Fatalf("pending ticket must not raise a blocker policy error: err=%v stdout=%q stderr=%q", err, stdoutText, stderrText)
	}
	if !strings.Contains(stdoutText, "01-marked.md") || strings.Contains(stdoutText, "02-gated.md") {
		t.Fatalf("pending ticket must leave the ready queue: stdout=%q", stdoutText)
	}
	if !strings.Contains(stderrText, string(tracker.DiagnosticStandardsPending)) || !strings.Contains(stderrText, "return to Plan and rerun standards extraction") {
		t.Fatalf("stderr must carry the standards gate guidance: stderr=%q", stderrText)
	}
}

// --- 03: journal stage gate exclusion (arming) ---

// journalGateFrontierFixture writes a project config plus one marker-complete
// Write-set ticket and deliberately no journal.md, so the armed stage gate is
// the only gate in play.
func journalGateFrontierFixture(t *testing.T, configBody string) string {
	t.Helper()
	root := t.TempDir()
	configDir := filepath.Join(root, ".flowforge")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("version: 5.0.0\n"+configBody), 0o644); err != nil {
		t.Fatal(err)
	}
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	execContract := "## Execution detail\n\n### Verified contracts\n\n- catalog.go owns ticket diagnostics.\n\n### Execution scenarios\n\n- Success: gated ticket stays out.\n- Failure: missing journal blocks dispatch.\n\n### Expected tests\n\n- go test ./internal/command/...\n\n### Generated artifacts\n\n- Not applicable — diagnostics are computed at runtime.\n\n### Conventions\n\n- Reuse existing blocker projection.\n"
	body := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 01-gated.md\n**Status:** open\n**Blocked by:** None\n\n" + execContract + "\n## Constraints\n\n- must keep behavior\n- Write set: internal/tracker/\n"
	if err := os.WriteFile(filepath.Join(issues, "01-gated.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// writeJournalGateJournal gives the fixture's proposal a complete active-round
// journal (align skipped with reason, design done).
func writeJournalGateJournal(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, "feature")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nflowforge:\n  schema: 1\n  role: journal\n  id: feature-journal\n  events:\n" +
		"  - {round: 1, stage: align, action: skipped, reason: \"fast-track\"}\n" +
		"  - {round: 1, stage: design, action: done}\n" +
		"---\n# Journal\n"
	if err := os.WriteFile(filepath.Join(dir, "journal.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestFrontierExcludesJournalGateBlockers pins the generic blocker channel both
// journal gate codes ride: SeverityBlocker alone moves them out of ready, so
// classifyReady needs no journal-specific branch.
func TestFrontierExcludesJournalGateBlockers(t *testing.T) {
	missing := &tracker.Issue{ID: "01", FilePath: "01.md"}
	incomplete := &tracker.Issue{ID: "02", FilePath: "02.md"}
	clean := &tracker.Issue{ID: "03", FilePath: "03.md"}
	diagnostics := []tracker.Diagnostic{
		{Code: tracker.DiagnosticJournalMissing, Severity: tracker.SeverityBlocker, Artifact: missing.FilePath},
		{Code: tracker.DiagnosticJournalStageIncomplete, Severity: tracker.SeverityBlocker, Artifact: incomplete.FilePath},
	}

	gotClean, gotWarnings, gotGaps, gotBlocked := classifyReady([]*tracker.Issue{missing, incomplete, clean}, diagnostics)
	if len(gotClean) != 1 || gotClean[0] != clean {
		t.Fatalf("clean = %#v, want only %s", gotClean, clean.FilePath)
	}
	if len(gotBlocked) != 2 || gotBlocked[0] != missing || gotBlocked[1] != incomplete {
		t.Fatalf("journal gate blockers must leave ready: blocked=%#v", gotBlocked)
	}
	if len(gotWarnings) != 0 || len(gotGaps) != 0 {
		t.Fatalf("unexpected buckets: warnings=%#v gaps=%#v", gotWarnings, gotGaps)
	}
}

// TestFrontierJournalGateUsesConfigMarker is the end-to-end arming check:
// catalog_discovery passes config journal.migrated into tracker.Options, so a
// migrated project with no journal.md drops the ticket and names the code.
func TestFrontierJournalGateUsesConfigMarker(t *testing.T) {
	root := journalGateFrontierFixture(t, "journal:\n  migrated: true\n  migrated_at: \"2026-10-12T09:00:00Z\"\n")
	stdoutText, stderrText, err := runStandardsGateFrontier(t, root)
	if err != errPolicyViolation {
		t.Fatalf("journal-missing blocker must fail the frontier policy: err=%v stdout=%q stderr=%q", err, stdoutText, stderrText)
	}
	if strings.Contains(stdoutText, "01-gated.md") {
		t.Fatalf("gated ticket must leave the ready queue: stdout=%q", stdoutText)
	}
	if !strings.Contains(stderrText, string(tracker.DiagnosticJournalMissing)) {
		t.Fatalf("stderr must name the journal gate code: stderr=%q", stderrText)
	}
}

// TestFrontierJournalGateSilentWhenMarkerUnset pins the arming lever end to
// end: the same proposal without the marker keeps its ticket ready.
func TestFrontierJournalGateSilentWhenMarkerUnset(t *testing.T) {
	root := journalGateFrontierFixture(t, "journal:\n  migrated: false\n")
	stdoutText, stderrText, err := runStandardsGateFrontier(t, root)
	if err != nil {
		t.Fatalf("unarmed gate must not block: err=%v stdout=%q stderr=%q", err, stdoutText, stderrText)
	}
	if !strings.Contains(stdoutText, "01-gated.md") {
		t.Fatalf("unarmed gate must keep the ticket ready: stdout=%q stderr=%q", stdoutText, stderrText)
	}
	if strings.Contains(stderrText, string(tracker.DiagnosticJournalMissing)) {
		t.Fatalf("unarmed gate must stay silent: stderr=%q", stderrText)
	}
}

// TestFrontierJournalGatePassesCompleteActiveRound pins the success scenario:
// an armed project whose journal finished align and design dispatches normally.
func TestFrontierJournalGatePassesCompleteActiveRound(t *testing.T) {
	root := journalGateFrontierFixture(t, "journal:\n  migrated: true\n")
	writeJournalGateJournal(t, root)
	stdoutText, stderrText, err := runStandardsGateFrontier(t, root)
	if err != nil {
		t.Fatalf("complete active round must stay ready: err=%v stdout=%q stderr=%q", err, stdoutText, stderrText)
	}
	if !strings.Contains(stdoutText, "01-gated.md") {
		t.Fatalf("complete active round must keep the ticket ready: stdout=%q stderr=%q", stdoutText, stderrText)
	}
}
