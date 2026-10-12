package tracker_test

import (
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/tracker"
)

const standardsTicketFrontmatter = "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"

func writeStandardsTicket(t *testing.T, root, feature, name, body string) string {
	t.Helper()
	dir := filepath.Join(root, feature, "issues")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(standardsTicketFrontmatter+body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func standardsTicketBody(status, constraints string) string {
	return "# 01: sample\n\n**Blocked by:** None\n\n**Status:** " + status +
		"\n\n## Changes\n\n- [ ] 1. Do the thing\n\n## Constraints\n\n" + constraints + "\n"
}

// standardsDiagnostics returns the standards-gate diagnostics recorded for one
// artifact path, ignoring unrelated diagnostics such as the execution-contract
// gap or completion-evidence warnings.
func standardsDiagnostics(catalog *tracker.Catalog, path string) []tracker.Diagnostic {
	var found []tracker.Diagnostic
	for _, diagnostic := range catalog.Diagnostics {
		if diagnostic.Artifact == path && isStandardsGateCode(diagnostic.Code) {
			found = append(found, diagnostic)
		}
	}
	return found
}

func isStandardsGateCode(code tracker.DiagnosticCode) bool {
	switch code {
	case tracker.DiagnosticStandardsMarkerMissing, tracker.DiagnosticStandardsPending, tracker.DiagnosticStandardsClausesMissing:
		return true
	}
	return false
}

// TestStandardsMarkerGateMatrix covers the three-state marker vocabulary in
// Constraints (must clause / none found / pending / all missing), the two skip
// paths (pure documentation ticket, terminal ticket), the legacy exemption
// channel, and near-miss spellings that must not satisfy the gate.
func TestStandardsMarkerGateMatrix(t *testing.T) {
	const writeSet = "- Write set: internal/tracker/"
	cases := []struct {
		name         string
		status       string
		constraints  string
		exempt       bool
		want         []tracker.DiagnosticCode
		wantSeverity tracker.DiagnosticSeverity
	}{
		{"must-clause", "open", "- must keep behavior\n" + writeSet, false, nil, ""},
		{"must-not-clause", "open", "- must not shell out\n" + writeSet, false, nil, ""},
		{"none-found", "open", "- standards: none found per guide\n" + writeSet, false, nil, ""},
		{"pending", "open", "- standards: pending\n" + writeSet, false, []tracker.DiagnosticCode{tracker.DiagnosticStandardsPending}, tracker.SeverityWarning},
		{"all-missing", "open", "- keep the change small\n" + writeSet, false, []tracker.DiagnosticCode{tracker.DiagnosticStandardsMarkerMissing}, tracker.SeverityBlocker},
		{"claimed-missing", "claimed", "- keep the change small\n" + writeSet, false, []tracker.DiagnosticCode{tracker.DiagnosticStandardsMarkerMissing}, tracker.SeverityBlocker},
		{"near-miss-none", "open", "- standards: none\n" + writeSet, false, []tracker.DiagnosticCode{tracker.DiagnosticStandardsMarkerMissing}, tracker.SeverityBlocker},
		{"near-miss-must-prefix", "open", "- mustang oil\n" + writeSet, false, []tracker.DiagnosticCode{tracker.DiagnosticStandardsMarkerMissing}, tracker.SeverityBlocker},
		{"pure-doc", "open", "- standards rules live in the skill contract", false, nil, ""},
		{"terminal", "closed", "- keep the change small\n" + writeSet, false, nil, ""},
		{"exempt", "open", "- keep the change small\n" + writeSet, true, nil, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			path := writeStandardsTicket(t, root, "feature", "01-sample.md", standardsTicketBody(tc.status, tc.constraints))
			options := tracker.Options{}
			if tc.exempt {
				options.StandardsExemptProposals = []string{"feature"}
			}
			catalog, err := tracker.DiscoverArtifactsWithConfig(root, options)
			if err != nil {
				t.Fatal(err)
			}
			found := standardsDiagnostics(catalog, path)
			if len(found) != len(tc.want) {
				t.Fatalf("standards diagnostics = %#v, want %v", found, tc.want)
			}
			for i, diagnostic := range found {
				if diagnostic.Code != tc.want[i] {
					t.Fatalf("standards diagnostic[%d] = %q, want %q", i, diagnostic.Code, tc.want[i])
				}
				if tc.wantSeverity != "" && diagnostic.Severity != tc.wantSeverity {
					t.Fatalf("standards diagnostic %q severity = %q, want %q", diagnostic.Code, diagnostic.Severity, tc.wantSeverity)
				}
				if diagnostic.Message == "" || diagnostic.Source.Path != path {
					t.Fatalf("standards diagnostic %q must carry a message and source path: %#v", diagnostic.Code, diagnostic)
				}
			}
			if len(tc.want) == 1 {
				assertDiagnostic(t, catalog.Diagnostics, tc.want[0], path)
			}
		})
	}
}

// TestStandardsDesignClausesDiagnostic covers the design-side check: a design
// authority consumed by a non-terminal Write-set ticket must carry either a
// `## Standards clauses` section or an explicit `standards: none found per
// guide` declaration.
func TestStandardsDesignClausesDiagnostic(t *testing.T) {
	designBody := func(section string) string {
		return "---\nflowforge:\n  schema: 1\n  role: design\n  id: feature-design\n  revision: 1\n---\n<a id=\"feature-design\"></a>\n# Design\n\n" + section
	}
	ticketBody := func(status, constraints string, consumesDesign bool) string {
		metadata := "---\nflowforge:\n  schema: 1\n  role: ticket\n"
		if consumesDesign {
			metadata += "  consumes:\n    design:\n      feature-design: 1\n"
		}
		metadata += "---\n"
		return metadata + standardsTicketBody(status, constraints)
	}
	const writeSet = "- must keep behavior\n- Write set: internal/tracker/"
	cases := []struct {
		name         string
		design       string
		ticketStatus string
		constraints  string
		consumes     bool
		want         bool
	}{
		{"missing-clauses", designBody("# Design body\n"), "open", writeSet, true, true},
		{"clauses-present", designBody("## Standards clauses\n\n- must keep behavior\n"), "open", writeSet, true, false},
		{"none-found-declared", designBody("standards: none found per guide\n"), "open", writeSet, true, false},
		{"terminal-consumer", designBody("# Design body\n"), "closed", writeSet, true, false},
		{"pure-doc-consumer", designBody("# Design body\n"), "open", "- standards rules live in the skill contract", true, false},
		{"not-consumed", designBody("# Design body\n"), "open", writeSet, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			featureDir := filepath.Join(root, "feature")
			if err := os.MkdirAll(filepath.Join(featureDir, "issues"), 0o755); err != nil {
				t.Fatal(err)
			}
			designPath := filepath.Join(featureDir, "design.md")
			if err := os.WriteFile(designPath, []byte(tc.design), 0o644); err != nil {
				t.Fatal(err)
			}
			ticketPath := filepath.Join(featureDir, "issues", "01-sample.md")
			if err := os.WriteFile(ticketPath, []byte(ticketBody(tc.ticketStatus, tc.constraints, tc.consumes)), 0o644); err != nil {
				t.Fatal(err)
			}
			catalog, err := tracker.DiscoverArtifacts(root)
			if err != nil {
				t.Fatal(err)
			}
			found := standardsDiagnostics(catalog, designPath)
			if !tc.want {
				if len(found) != 0 {
					t.Fatalf("unexpected design clauses diagnostic: %#v", found)
				}
				return
			}
			if len(found) != 1 || found[0].Code != tracker.DiagnosticStandardsClausesMissing || found[0].Severity != tracker.SeverityWarning {
				t.Fatalf("design clauses diagnostic = %#v, want one warning %q", found, tracker.DiagnosticStandardsClausesMissing)
			}
			assertDiagnostic(t, catalog.Diagnostics, tracker.DiagnosticStandardsClausesMissing, designPath)
		})
	}
}

// TestStandardsDesignClausesDiagnosticExempt covers the legacy exemption
// channel on the design-side check (Fix Change 8): a design sitting inside an
// exempt proposal, or a design whose Write-set consuming tickets all sit
// inside exempt proposals, produces no standards-clauses-missing diagnostic.
func TestStandardsDesignClausesDiagnosticExempt(t *testing.T) {
	const writeSet = "- must keep behavior\n- Write set: internal/tracker/"
	consumer := func(consumedKey string) string {
		return "---\nflowforge:\n  schema: 1\n  role: ticket\n  consumes:\n    design:\n      " + consumedKey + ": 1\n---\n" +
			standardsTicketBody("open", writeSet)
	}
	// Fixture: the design lives in feature-a; feature-b holds a remote consumer,
	// and feature-a optionally holds a local one.
	fixture := func(t *testing.T, withLocalConsumer bool) (string, string) {
		t.Helper()
		root := t.TempDir()
		for _, dir := range []string{filepath.Join(root, "feature-a", "issues"), filepath.Join(root, "feature-b", "issues")} {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
		}
		designPath := filepath.Join(root, "feature-a", "design.md")
		design := "---\nflowforge:\n  schema: 1\n  role: design\n  id: feature-a-design\n  revision: 1\n---\n<a id=\"feature-a-design\"></a>\n# Design\n\n## Scope\n\nNo standards clauses section here.\n"
		if err := os.WriteFile(designPath, []byte(design), 0o644); err != nil {
			t.Fatal(err)
		}
		if withLocalConsumer {
			local := filepath.Join(root, "feature-a", "issues", "01-local.md")
			if err := os.WriteFile(local, []byte(consumer("feature-a-design")), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		remote := filepath.Join(root, "feature-b", "issues", "01-remote.md")
		if err := os.WriteFile(remote, []byte(consumer("feature-a/feature-a-design")), 0o644); err != nil {
			t.Fatal(err)
		}
		return root, designPath
	}
	cases := []struct {
		name      string
		exempt    []string
		withLocal bool
		want      bool
	}{
		{"no-exemption-control", nil, true, true},
		{"remote-only-control", nil, false, true},
		{"exempt-consumer", []string{"feature-b"}, false, false},
		{"mixed-consumers", []string{"feature-b"}, true, true},
		{"exempt-design", []string{"feature-a"}, true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, designPath := fixture(t, tc.withLocal)
			catalog, err := tracker.DiscoverArtifactsWithConfig(root, tracker.Options{StandardsExemptProposals: tc.exempt})
			if err != nil {
				t.Fatal(err)
			}
			found := standardsDiagnostics(catalog, designPath)
			if !tc.want {
				if len(found) != 0 {
					t.Fatalf("exempt proposal must skip the design clauses diagnostic: %#v", found)
				}
				return
			}
			if len(found) != 1 || found[0].Code != tracker.DiagnosticStandardsClausesMissing {
				t.Fatalf("design clauses diagnostic = %#v, want one %q", found, tracker.DiagnosticStandardsClausesMissing)
			}
		})
	}
}
