package tracker_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/tracker"
)

func TestDiscoverArtifactsProjectsOnlyIssueTickets(t *testing.T) {
	root := t.TempDir()
	feature := filepath.Join(root, "catalog-refinement")
	issues := filepath.Join(feature, "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}

	write := func(path, content string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(filepath.Join(feature, "spec.md"), `---
flowforge:
  schema: 1
  role: requirement
  id: catalog-requirements
  revision: 1
---
# Catalog requirements
`)
	write(filepath.Join(issues, "01-catalog.md"), `---
flowforge:
  schema: 1
  role: ticket
---
# 01: Build catalog

**Status:** open
**Blocked by:** None
`)
	write(filepath.Join(issues, "02-consumer.md"), `# 02: Keep legacy tickets

**Status:** open
**Blocked by:** 01
`)
	write(filepath.Join(issues, "03-evidence.md"), `---
flowforge:
  schema: 1
  role: evidence
---
# Evidence in the wrong directory
`)

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatalf("DiscoverArtifacts failed: %v", err)
	}
	if len(catalog.Artifacts) != 4 {
		t.Fatalf("expected 4 artifacts, got %d", len(catalog.Artifacts))
	}

	tickets := catalog.ExecutableTickets()
	if len(tickets) != 2 {
		t.Fatalf("expected 2 executable tickets, got %d", len(tickets))
	}
	if tickets[0].ID != "01" || tickets[0].Title != "Build catalog" {
		t.Fatalf("unexpected schema ticket: %#v", tickets[0])
	}
	if tickets[1].ID != "02" || len(tickets[1].BlockedBy) != 1 || tickets[1].BlockedBy[0] != "01" {
		t.Fatalf("legacy ticket body metadata was not preserved: %#v", tickets[1])
	}

	assertDiagnostic(t, catalog.Diagnostics, "legacy-metadata", filepath.Join(issues, "02-consumer.md"))
	assertDiagnostic(t, catalog.Diagnostics, "role-location-conflict", filepath.Join(issues, "03-evidence.md"))
}

func TestDiscoverArtifactsRejectsMalformedMetadataSafely(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"01-unclosed.md":    "---\nflowforge:\n  schema: 1\n  role: ticket\n# no closing delimiter\n",
		"02-false-close.md": "---\nflowforge:\n  schema: 1\n  role: ticket\n---not-a-close\n# no closing delimiter\n",
		"03-bad-role.md":    "---\nflowforge:\n  schema: 1\n  role: imaginary\n---\n# bad role\n",
		"04-bad-schema.md":  "---\nflowforge:\n  schema: -1\n  role: ticket\n---\n# bad schema\n",
	}
	for name, content := range cases {
		if err := os.WriteFile(filepath.Join(issues, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Tickets) != 2 {
		t.Fatalf("malformed YAML may legacy-fallback, invalid known metadata may not execute: %#v", catalog.Tickets)
	}
	assertDiagnostic(t, catalog.Diagnostics, "invalid-frontmatter", filepath.Join(issues, "01-unclosed.md"))
	assertDiagnostic(t, catalog.Diagnostics, "invalid-frontmatter", filepath.Join(issues, "02-false-close.md"))
	assertDiagnostic(t, catalog.Diagnostics, "invalid-metadata", filepath.Join(issues, "03-bad-role.md"))
	assertDiagnostic(t, catalog.Diagnostics, "invalid-metadata", filepath.Join(issues, "04-bad-schema.md"))
	if catalog.Diagnostics[0].Source.Path == "" {
		t.Fatal("diagnostics must expose a source path")
	}
}

func TestArtifactRoleLocationMatrix(t *testing.T) {
	root := t.TempDir()
	feature := filepath.Join(root, "feature")
	issues := filepath.Join(feature, "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}

	roles := []string{"requirement", "design", "spec", "evidence", "research", "map"}
	for _, role := range roles {
		body := "---\nflowforge:\n  schema: 1\n  role: " + role + "\n---\n# " + role + "\n"
		if err := os.WriteFile(filepath.Join(feature, role+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		misplaced := filepath.Join(issues, role+".md")
		if err := os.WriteFile(misplaced, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	outsideTicket := filepath.Join(feature, "ticket.md")
	if err := os.WriteFile(outsideTicket, []byte("---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# Ticket\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Tickets) != 0 {
		t.Fatalf("role/location matrix projected non-executable artifacts: %#v", catalog.Tickets)
	}
	assertDiagnostic(t, catalog.Diagnostics, "role-location-conflict", outsideTicket)
	for _, role := range roles {
		assertDiagnostic(t, catalog.Diagnostics, "role-location-conflict", filepath.Join(issues, role+".md"))
	}
}

func TestClosedTicketRequiresObservableCompletionEvidence(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}

	missing := filepath.Join(issues, "01-missing.md")
	empty := filepath.Join(issues, "02-empty.md")
	present := filepath.Join(issues, "03-present.md")
	for path, body := range map[string]string{
		missing: "# 01: Missing\n\n**Blocked by:** None\n**Status:** closed\n\n## Changes\n\nDone.\n",
		empty:   "# 02: Empty\n\n**Blocked by:** None\n**Status:** closed\n\n## Completion evidence\n\n## Notes\n\nNone.\n",
		present: "# 03: Present\n\n**Blocked by:** None\n**Status:** closed\n\n## Completion evidence\n\n- Verification: `go test ./...` passed.\n",
	} {
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	assertDiagnostic(t, catalog.Diagnostics, tracker.DiagnosticMissingEvidence, missing)
	assertDiagnostic(t, catalog.Diagnostics, tracker.DiagnosticMissingEvidence, empty)
	for _, diagnostic := range catalog.Diagnostics {
		if diagnostic.Code == tracker.DiagnosticMissingEvidence && diagnostic.Artifact == present {
			t.Fatalf("ticket with observable evidence was rejected: %#v", diagnostic)
		}
	}
}

func TestBlockedEvidenceDiagnostic(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	blockedSection := "\n## Blocked evidence\n\n- Verbatim: exit status 1\n- cmd: go test ./... (exit 1)\n- Next: inspect the failing package\n"

	writeTicket := func(name, status, body string) string {
		t.Helper()
		path := filepath.Join(issues, name)
		content := fm + "# " + strings.TrimSuffix(name, ".md") + "\n**Status:** " + status + "\n**Blocked by:** None\n" + body
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	reported := writeTicket("01-reported.md", "open", blockedSection)
	headingOnly := writeTicket("02-heading-only.md", "open", "\n## Blocked evidence\n")
	caseVariant := writeTicket("03-case-variant.md", "open", "\n##   BLOCKED evidence   \n- fact\n")
	readyAgent := writeTicket("04-ready-for-agent.md", "ready-for-agent", blockedSection)
	closedTicket := writeTicket("05-closed.md", "closed", "\n## Completion evidence\n\n- go test ./... passed.\n"+blockedSection)
	removed := writeTicket("06-removed.md", "open", "\n## Notes\n\n- blocked section was consumed\n")
	nested := writeTicket("07-nested.md", "open", "\n### Blocked evidence\n- not the pinned heading level\n")
	suffixed := writeTicket("08-suffixed.md", "open", "\n## Blocked evidence: extra\n- not the exact pinned heading\n")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{reported, headingOnly, caseVariant, readyAgent} {
		assertDiagnostic(t, catalog.Diagnostics, tracker.DiagnosticBlockedEvidencePresent, path)
	}
	for _, path := range []string{closedTicket, removed, nested, suffixed} {
		for _, diagnostic := range catalog.Diagnostics {
			if diagnostic.Code == tracker.DiagnosticBlockedEvidencePresent && diagnostic.Artifact == path {
				t.Fatalf("unexpected blocked-evidence diagnostic for %s: %#v", path, diagnostic)
			}
		}
	}
	for _, diagnostic := range catalog.Diagnostics {
		if diagnostic.Code != tracker.DiagnosticBlockedEvidencePresent || diagnostic.Artifact != reported {
			continue
		}
		if diagnostic.Severity != tracker.SeverityWarning {
			t.Fatalf("severity = %q, want warning", diagnostic.Severity)
		}
		if diagnostic.Source.Path != reported {
			t.Fatalf("source path = %q, want %s", diagnostic.Source.Path, reported)
		}
		if !strings.Contains(diagnostic.Message, "flowforge-refine-ticket") {
			t.Fatalf("message must point dispatchers at the refine consumption path: %q", diagnostic.Message)
		}
	}
}

func TestExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}

	complete := `## Execution detail

### Verified contracts

- internal/tracker/catalog.go:discoverArtifact owns ticket diagnostics.

### Execution scenarios

- Success: a complete ticket is eligible.
- Failure: a missing section is diagnosed.

### Expected tests

- go test ./internal/tracker/...

### Generated artifacts

- Not applicable — this ticket has no generated output.

### Conventions

- Reuse catalog diagnostics.
`
	sectionContent := map[string]string{
		"Verified contracts":  "- internal/tracker/catalog.go:discoverArtifact owns ticket diagnostics.",
		"Execution scenarios": "- Success: a complete ticket is eligible.\n- Failure: a missing section is diagnosed.",
		"Expected tests":      "- go test ./internal/tracker/...",
		"Generated artifacts": "- Not applicable — this ticket has no generated output.",
		"Conventions":         "- Reuse catalog diagnostics.",
	}
	placeholderOnly := strings.Replace(complete, "- go test ./internal/tracker/...", "- [ ]", 1)
	headingsOutsideDetail := "## Execution detail\n\n## Notes\n\n" + strings.TrimPrefix(complete, "## Execution detail\n\n")
	writeTicket := func(name, status, detail string, managed bool) string {
		t.Helper()
		frontmatter := ""
		if managed {
			frontmatter = "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
		}
		path := filepath.Join(issues, name)
		body := frontmatter + "# " + strings.TrimSuffix(name, ".md") + "\n**Status:** " + status + "\n**Blocked by:** None\n\n" + detail
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	completePath := writeTicket("01-complete.md", "open", complete, true)
	readyPath := writeTicket("02-ready.md", "ready-for-agent", complete, true)
	missingPaths := make([]string, 0, len(sectionContent))
	for index, section := range []string{"Verified contracts", "Execution scenarios", "Expected tests", "Generated artifacts", "Conventions"} {
		missing := strings.Replace(complete, "### "+section+"\n\n"+sectionContent[section], "", 1)
		missingPaths = append(missingPaths, writeTicket("0"+string(rune('3'+index))+"-missing.md", "open", missing, true))
	}
	placeholderPath := writeTicket("08-placeholder.md", "ready-for-agent", placeholderOnly, true)
	outsidePath := writeTicket("09-outside.md", "open", headingsOutsideDetail, true)
	closedPath := writeTicket("10-closed.md", "closed", "", true)
	legacyPath := writeTicket("11-legacy.md", "open", "", false)

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range append(missingPaths, placeholderPath, outsidePath) {
		assertDiagnostic(t, catalog.Diagnostics, tracker.DiagnosticCode("execution-contract-incomplete"), path)
	}
	for _, path := range []string{completePath, readyPath, closedPath, legacyPath} {
		for _, diagnostic := range catalog.Diagnostics {
			if diagnostic.Code == tracker.DiagnosticCode("execution-contract-incomplete") && diagnostic.Artifact == path {
				t.Fatalf("unexpected execution-contract diagnostic for %s: %#v", path, diagnostic)
			}
		}
	}
}

func TestNeedsRepairIsNonExecutableAndRepairProvenanceIsParsed(t *testing.T) {
	if tracker.StatusNeedsRepair.IsExecutable() {
		t.Fatal("needs-repair must not be executable")
	}
	if tracker.StatusNeedsRepair.IsTerminal() {
		t.Fatal("needs-repair must not be terminal")
	}

	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	writeTicket := func(name, header string) string {
		path := filepath.Join(issues, name)
		body := fm + "# " + strings.TrimSuffix(name, ".md") + "\n" + header + "\n**Blocked by:** None\n"
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}
	original := writeTicket("01-original.md", "**Status:** needs-repair\n**Repair of:** None")
	repair := writeTicket("02-repair.md", "**Status:** open\n**Repair of:** 01")
	legacy := writeTicket("03-legacy.md", "**Status:** open")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	var origTicket, repairTicket *tracker.Issue
	for _, ticket := range catalog.Tickets {
		switch ticket.FilePath {
		case original:
			origTicket = ticket
		case repair:
			repairTicket = ticket
		}
	}
	if origTicket == nil || repairTicket == nil {
		t.Fatalf("missing tickets: original=%v repair=%v", origTicket, repairTicket)
	}
	if origTicket.Status != tracker.StatusNeedsRepair {
		t.Fatalf("original status = %q, want needs-repair", origTicket.Status)
	}
	if repairTicket.RepairOf != "01" {
		t.Fatalf("repair provenance = %q, want 01", repairTicket.RepairOf)
	}
	if origTicket.RepairOf != "" {
		t.Fatalf("original with 'Repair of: None' should have empty RepairOf, got %q", origTicket.RepairOf)
	}
	for _, ticket := range catalog.Tickets {
		if ticket.FilePath == legacy && ticket.RepairOf != "" {
			t.Fatalf("legacy ticket should have empty RepairOf, got %q", ticket.RepairOf)
		}
	}
}

func TestRepairReciprocalReferencesAreValidated(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	writeTicket := func(name, status, repairOf string) string {
		path := filepath.Join(issues, name)
		header := "# " + strings.TrimSuffix(name, ".md") + "\n**Status:** " + status + "\n"
		if repairOf != "" {
			header += "**Repair of:** " + repairOf + "\n"
		}
		header += "**Blocked by:** None\n"
		body := fm + header
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}
	orphanRepair := writeTicket("01-orphan-repair.md", "open", "99")
	originalNoRepair := writeTicket("02-no-reciprocal.md", "needs-repair", "")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	assertDiagnostic(t, catalog.Diagnostics, tracker.DiagnosticCode("dangling-repair-reference"), orphanRepair)
	assertDiagnostic(t, catalog.Diagnostics, tracker.DiagnosticCode("missing-reciprocal-repair"), originalNoRepair)
}

func TestEndToEndRepairDAGPath(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"
	execContract := "\n\n## Execution detail\n\n### Verified contracts\n\n- model.go owns status.\n\n### Execution scenarios\n\n- Success: repair completes.\n- Failure: original stays needs-repair.\n\n### Expected tests\n\n- go test ./internal/tracker/...\n\n### Generated artifacts\n\n- Not applicable.\n\n### Conventions\n\n- Reuse existing DAG logic.\n"
	writeTicket := func(name, status, blockedBy, repairOf string) string {
		path := filepath.Join(issues, name)
		header := "# " + strings.TrimSuffix(name, ".md") + "\n**Status:** " + status + "\n"
		if repairOf != "" {
			header += "**Repair of:** " + repairOf + "\n"
		}
		header += "**Blocked by:** " + blockedBy + "\n"
		body := fm + header + execContract
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	original := writeTicket("01-original.md", "needs-repair", "None", "")
	repair := writeTicket("02-repair.md", "open", "None", "01")
	downstream := writeTicket("03-downstream.md", "open", "02", "")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}

	var origIssue, repairIssue, downstreamIssue *tracker.Issue
	for _, ticket := range catalog.Tickets {
		switch ticket.FilePath {
		case original:
			origIssue = ticket
		case repair:
			repairIssue = ticket
		case downstream:
			downstreamIssue = ticket
		}
	}
	if origIssue == nil || repairIssue == nil || downstreamIssue == nil {
		t.Fatalf("missing tickets: orig=%v repair=%v downstream=%v", origIssue, repairIssue, downstreamIssue)
	}

	if origIssue.Status != tracker.StatusNeedsRepair {
		t.Fatalf("original status = %q, want needs-repair", origIssue.Status)
	}
	if origIssue.Status.IsExecutable() {
		t.Fatal("needs-repair must be non-executable")
	}
	if repairIssue.RepairOf != "01" {
		t.Fatalf("repair provenance = %q, want 01", repairIssue.RepairOf)
	}
	if !repairIssue.Status.IsExecutable() {
		t.Fatal("repair ticket must be executable")
	}

	g := tracker.BuildGraph(catalog.Tickets)
	frontier := g.ComputeFrontier()
	for _, ready := range frontier.Ready {
		if ready.ID == "01" {
			t.Fatal("needs-repair original must not be in ready")
		}
		if ready.ID == "03" {
			t.Fatal("downstream must not be in ready while repair is open")
		}
	}
	var repairInReady bool
	for _, ready := range frontier.Ready {
		if ready.ID == "02" {
			repairInReady = true
		}
	}
	if !repairInReady {
		t.Fatal("repair ticket must be in ready")
	}
	var downstreamBlocked bool
	for _, b := range frontier.Blocked {
		if b.Issue != nil && b.Issue.ID == "03" {
			downstreamBlocked = true
		}
	}
	if !downstreamBlocked {
		t.Fatal("downstream must be blocked waiting on repair")
	}

	noRepairDiagnostics := true
	for _, d := range catalog.Diagnostics {
		if d.Code == tracker.DiagnosticDanglingRepairReference || d.Code == tracker.DiagnosticMissingReciprocalRepair {
			if d.Artifact == original || d.Artifact == repair {
				noRepairDiagnostics = false
			}
		}
	}
	if !noRepairDiagnostics {
		t.Fatal("valid repair pair must not produce reciprocal-repair diagnostics")
	}
}

func TestRepositoryAndTangramProposalLayouts(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	flowforge, err := tracker.DiscoverArtifacts(filepath.Join(repoRoot, "docs", "proposals"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ticket := range flowforge.Tickets {
		if filepath.Base(ticket.FilePath) == "spec.md" {
			t.Fatal("FlowForge spec.md entered executable projection")
		}
	}

	tangramRoot := "/vol3/1000/develop/tangram-v2/ff-wiki-v5/proposals"
	if _, err := os.Stat(tangramRoot); err != nil {
		t.Skip("Tangram compatibility checkout is not available")
	}
	tangram, err := tracker.DiscoverArtifacts(tangramRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(tangram.Tickets) == 0 {
		t.Fatal("expected non-empty Tangram tickets")
	}
	for _, ticket := range tangram.Tickets {
		if filepath.Base(ticket.FilePath) == "spec.md" || filepath.Base(filepath.Dir(ticket.FilePath)) != "issues" {
			t.Fatalf("Tangram non-ticket entered executable projection: %s", ticket.FilePath)
		}
	}
}

func TestSemanticDiagnosticsAndExactWaiver(t *testing.T) {
	root := t.TempDir()
	feature := filepath.Join(root, "feature")
	issues := filepath.Join(feature, "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(path, body string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(feature, "design.md"), `---
flowforge:
  schema: 1
  role: design
  id: feature-design
  areas:
    route-seam: {revision: 3, anchor: route-seam}
  open_items:
    - {id: verify-route, diagnostic: verification-unselected, severity: gap, affects: [02], anchor: verify-route}
---
<a id="route-seam"></a>
## Route seam
<a id="verify-route"></a>
## Verification remains open
`)
	write(filepath.Join(issues, "02-consume.md"), `---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design: {route-seam: 2, missing-seam: 1}
  waivers:
    - {diagnostic: upstream-changed, target: route-seam, reason: reviewed unchanged request shape}
---
# 02: Consume
**Status:** open
**Blocked by:** None
See [route seam](../design.md#route-seam).
`)
	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	assertDiagnostic(t, catalog.Diagnostics, "missing-authority", filepath.Join(issues, "02-consume.md"))
	assertDiagnostic(t, catalog.Diagnostics, "verification-unselected", filepath.Join(issues, "02-consume.md"))
	for _, d := range catalog.Diagnostics {
		if d.Code == "upstream-changed" {
			if d.Waiver == nil || d.Waiver.Reason == "" {
				t.Fatal("exact waiver was not attached")
			}
			return
		}
	}
	t.Fatal("missing upstream-changed diagnostic")
}

func TestSemanticDiagnosticFailureMatrix(t *testing.T) {
	root := t.TempDir()
	feature := filepath.Join(root, "feature")
	issues := filepath.Join(feature, "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(path, body string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(feature, "design-a.md"), "---\nflowforge:\n  schema: 1\n  role: design\n  areas:\n    shared: {revision: 2, anchor: shared}\n  open_items:\n    - {id: bad-scope, diagnostic: design-gap, severity: gap, affects: [99], anchor: shared}\n---\n<a id=\"shared\"></a>\n")
	write(filepath.Join(feature, "design-b.md"), "---\nflowforge:\n  schema: 1\n  role: design\n  areas:\n    shared: {revision: 3, anchor: shared}\n    invalid: {revision: 0, anchor: missing}\n---\n<a id=\"shared\"></a>\n")
	write(filepath.Join(feature, "design-c.md"), "---\nflowforge:\n  schema: 1\n  role: design\n  areas:\n    future: {revision: 2, anchor: future}\n    scoped: {revision: 1, anchor: scoped}\n  open_items:\n    - {id: scoped-gap, diagnostic: scoped-design-gap, severity: gap, affects: [scoped], anchor: scoped}\n---\n<a id=\"future\"></a>\n<a id=\"scoped\"></a>\n")
	write(filepath.Join(issues, "01-consumer.md"), "---\nflowforge:\n  schema: 1\n  role: ticket\n  consumes:\n    design: {shared: 4, future: 4, scoped: 1}\n  waivers:\n    - {diagnostic: upstream-changed, target: '*', reason: blanket}\n    - {diagnostic: absent-code, target: shared, reason: old}\n---\n# 01: Consumer\n**Status:** open\n**Blocked by:** None\nSee [scoped](../design-c.md#scoped).\n")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(issues, "01-consumer.md")
	assertDiagnostic(t, catalog.Diagnostics, "duplicate-semantic-id", filepath.Join(feature, "design-b.md"))
	assertDiagnostic(t, catalog.Diagnostics, "duplicate-semantic-id", path)
	assertDiagnostic(t, catalog.Diagnostics, "missing-anchor", filepath.Join(feature, "design-b.md"))
	assertDiagnostic(t, catalog.Diagnostics, "invalid-open-item", filepath.Join(feature, "design-a.md"))
	assertDiagnostic(t, catalog.Diagnostics, "invalid-waiver", path)
	assertDiagnostic(t, catalog.Diagnostics, "stale-waiver", path)
	assertDiagnostic(t, catalog.Diagnostics, "future-consumed-revision", path)
	assertDiagnostic(t, catalog.Diagnostics, "missing-human-link", path)
	assertDiagnostic(t, catalog.Diagnostics, "scoped-design-gap", path)
}

func TestCrossFeatureSemanticLinksAndReverseMismatch(t *testing.T) {
	root := t.TempDir()
	for _, feature := range []string{"producer", "consumer"} {
		if err := os.MkdirAll(filepath.Join(root, feature, "issues"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	producer := filepath.Join(root, "producer", "design.md")
	if err := os.WriteFile(producer, []byte("---\nflowforge:\n  schema: 1\n  role: design\n  areas:\n    route: {revision: 2, anchor: route}\n---\n<a id=\"route\"></a>\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(root, "consumer", "issues", "01-linked.md")
	if err := os.WriteFile(linked, []byte("---\nflowforge:\n  schema: 1\n  role: ticket\n  consumes:\n    design: {producer/route: 2}\n---\n# 01: Linked\n**Status:** open\n**Blocked by:** None\nSee [route](../../producer/design.md#route).\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	untracked := filepath.Join(root, "consumer", "issues", "02-untracked.md")
	if err := os.WriteFile(untracked, []byte("---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 02: Untracked\n**Status:** open\n**Blocked by:** None\nSee [route](../../producer/design.md#route).\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range catalog.Diagnostics {
		if d.Artifact == linked && d.Code == "missing-human-link" {
			t.Fatalf("valid cross-feature link rejected: %#v", d)
		}
	}
	assertDiagnostic(t, catalog.Diagnostics, "untracked-upstream", untracked)
}

func assertDiagnostic(t *testing.T, diagnostics []tracker.Diagnostic, code tracker.DiagnosticCode, path string) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code && diagnostic.Artifact == path {
			return
		}
	}
	t.Fatalf("missing diagnostic %q for %s: %#v", code, path, diagnostics)
}
