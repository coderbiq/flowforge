package tracker_test

import (
	"os"
	"path/filepath"
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
	if len(tangram.Tickets) != 11 {
		t.Fatalf("expected 11 Tangram tickets, got %d", len(tangram.Tickets))
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
