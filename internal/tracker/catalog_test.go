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

func assertDiagnostic(t *testing.T, diagnostics []tracker.Diagnostic, code tracker.DiagnosticCode, path string) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code && diagnostic.Artifact == path {
			return
		}
	}
	t.Fatalf("missing diagnostic %q for %s: %#v", code, path, diagnostics)
}
