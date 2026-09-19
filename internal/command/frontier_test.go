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

func TestStrictCheckValidatesRequirementAndDesignWithoutTickets(t *testing.T) {
	root := t.TempDir()
	previousDir, previousJSON, previousStrict := checkDir, checkJSON, checkStrict
	defer func() { checkDir, checkJSON, checkStrict = previousDir, previousJSON, previousStrict }()
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
		checkDir, checkJSON, checkStrict = root, json, strict
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
	checkDir, checkJSON, checkStrict = root, true, false
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
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
			Issue     *tracker.Issue `json:"issue"`
			WaitingOn []string       `json:"waiting_on"`
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
		if b.Issue != nil && b.Issue.ID == "03" {
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
