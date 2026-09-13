package tracker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const evidenceTicketFrontmatter = "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n"

func writeEvidenceTicket(t *testing.T, root, feature, name, body string) string {
	t.Helper()
	dir := filepath.Join(root, feature, "issues")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(evidenceTicketFrontmatter+body), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func evidenceCodes(t *testing.T, catalog *Catalog, path string) []DiagnosticCode {
	t.Helper()
	var codes []DiagnosticCode
	for _, d := range catalog.Diagnostics {
		if d.Artifact == path && isEvidenceCode(d.Code) {
			codes = append(codes, d.Code)
		}
	}
	return codes
}

func isEvidenceCode(code DiagnosticCode) bool {
	switch code {
	case DiagnosticEvidenceMissing, DiagnosticEvidenceIncomplete, DiagnosticEvidenceExitNonzero, DiagnosticEvidenceArtifactMissing, DiagnosticEvidenceRepeatFailure:
		return true
	}
	return false
}

const validQuadruple = "- [x] 1. Add session filter\n    - cmd: ./gradlew :app:test --tests \"*RepoTest*\"\n    - exit: 0\n    - output: \"5 tests completed, 0 failed\"\n    - artifact: app/src/Repo.kt\n- [ ] 2. Not done yet\n"

func ticketBody(changes string) string {
	return "# 01: sample\n\n**Blocked by:** None\n\n**Status:** closed\n\n## Changes\n\n" + changes + "\n## Constraints\n\n- must keep behavior\n- Write set: app/\n"
}

func TestEvidenceQuadrupleParsing(t *testing.T) {
	root := t.TempDir()

	validArtifact := filepath.Join(root, "app", "src")
	if err := os.MkdirAll(validArtifact, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(validArtifact, "Repo.kt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name    string
		changes string
		want    DiagnosticCode
	}{
		{"missing", "- [x] 1. Did the thing\n", DiagnosticEvidenceMissing},
		{"incomplete", "- [x] 1. Did the thing\n    - cmd: `go test ./...`\n    - exit: 0\n", DiagnosticEvidenceIncomplete},
		{"exit-nonzero", "- [x] 1. Did the thing\n    - cmd: `go test ./...`\n    - exit: 1\n    - output: \"1 failed\"\n    - artifact: app/src/Repo.kt\n", DiagnosticEvidenceExitNonzero},
		{"artifact-missing", "- [x] 1. Did the thing\n    - cmd: `go test ./...`\n    - exit: 0\n    - output: \"ok\"\n    - artifact: app/src/Gone.kt\n", DiagnosticEvidenceArtifactMissing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeEvidenceTicket(t, root, "feature-"+tc.name, "01-"+tc.name+".md", ticketBody(tc.changes))
			catalog, err := DiscoverArtifacts(root)
			if err != nil {
				t.Fatal(err)
			}
			codes := evidenceCodes(t, catalog, path)
			if len(codes) != 1 || codes[0] != tc.want {
				t.Fatalf("expected exactly [%s], got %v", tc.want, codes)
			}
		})
	}

	t.Run("valid-quadruple-clean", func(t *testing.T) {
		path := writeEvidenceTicket(t, root, "feature-valid", "01-valid.md", ticketBody(validQuadruple))
		catalog, err := DiscoverArtifacts(root)
		if err != nil {
			t.Fatal(err)
		}
		if codes := evidenceCodes(t, catalog, path); len(codes) != 0 {
			t.Fatalf("expected no evidence diagnostics, got %v", codes)
		}
	})

	t.Run("unchecked-with-quadruple-clean", func(t *testing.T) {
		path := writeEvidenceTicket(t, root, "feature-unchecked", "01-unchecked.md", ticketBody("- [ ] 1. Later\n    - cmd: `go test ./...`\n    - exit: 0\n    - output: \"ok\"\n    - artifact: app/src/Repo.kt\n"))
		catalog, err := DiscoverArtifacts(root)
		if err != nil {
			t.Fatal(err)
		}
		if codes := evidenceCodes(t, catalog, path); len(codes) != 0 {
			t.Fatalf("expected no evidence diagnostics, got %v", codes)
		}
	})

	t.Run("pure-doc-ticket-skipped", func(t *testing.T) {
		path := writeEvidenceTicket(t, root, "feature-doc", "01-doc.md", "# 01: doc\n\n**Status:** closed\n\n## Changes\n\n- [x] 1. Wrote docs\n\n## Constraints\n\n- docs only\n")
		catalog, err := DiscoverArtifacts(root)
		if err != nil {
			t.Fatal(err)
		}
		if codes := evidenceCodes(t, catalog, path); len(codes) != 0 {
			t.Fatalf("expected no evidence diagnostics for pure doc ticket, got %v", codes)
		}
	})
}

func TestEvidenceExemptProposals(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "app"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "app", "a.kt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	exempted := writeEvidenceTicket(t, root, "legacy-proposal", "01-old.md", ticketBody("- [x] 1. Legacy done\n"))
	kept := writeEvidenceTicket(t, root, "fresh-proposal", "01-new.md", ticketBody("- [x] 1. New done\n"))

	catalog, err := DiscoverArtifactsWithConfig(root, Options{ExemptProposals: []string{"legacy-proposal"}})
	if err != nil {
		t.Fatal(err)
	}
	if codes := evidenceCodes(t, catalog, exempted); len(codes) != 0 {
		t.Fatalf("exempt proposal must produce no evidence diagnostics, got %v", codes)
	}
	if codes := evidenceCodes(t, catalog, kept); len(codes) != 1 || codes[0] != DiagnosticEvidenceMissing {
		t.Fatalf("non-exempt proposal must report evidence-missing, got %v", codes)
	}

	plain, err := DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	if codes := evidenceCodes(t, plain, exempted); len(codes) != 1 || codes[0] != DiagnosticEvidenceMissing {
		t.Fatalf("DiscoverArtifacts without options must still report, got %v", codes)
	}
}

func TestEvidenceArtifactBaseResolvesRepositoryRoot(t *testing.T) {
	repoRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repoRoot, "internal", "app"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "internal", "app", "code.go"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	scanRoot := filepath.Join(repoRoot, "docs", "proposals")
	changes := "- [x] 1. Shipped code\n    - cmd: go test ./...\n    - exit: 0\n    - output: \"ok\"\n    - artifact: internal/app/code.go\n"
	path := writeEvidenceTicket(t, scanRoot, "feature-repo", "01-repo.md", ticketBody(changes))

	catalog, err := DiscoverArtifacts(scanRoot)
	if err != nil {
		t.Fatal(err)
	}
	if codes := evidenceCodes(t, catalog, path); len(codes) != 0 {
		t.Fatalf("repository-relative artifact must resolve from repo root, got %v", codes)
	}
}

func TestEvidenceDontAffectUnmatchedTickets(t *testing.T) {
	root := t.TempDir()

	// A non-ticket role artifact carrying checked boxes must be ignored.
	docsDir := filepath.Join(root, "feature-x", "issues")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}
	nonTicket := "---\nflowforge:\n  schema: 1\n  role: design\n---\n<a id=\"d\"></a>\n# Design\n\n## Changes\n\n- [x] 1. Not a ticket Change\n\n## Constraints\n\n- Write set: app/\n"
	if err := os.WriteFile(filepath.Join(docsDir, "design-notes.md"), []byte(nonTicket), 0644); err != nil {
		t.Fatal(err)
	}

	catalog, err := DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range catalog.Diagnostics {
		if isEvidenceCode(d.Code) {
			t.Fatalf("non-ticket artifact must not produce evidence diagnostics: %s %s", d.Code, d.Artifact)
		}
	}
}

// failedEvidence renders a checked Change whose evidence quadruple records a
// failed run of cmd (exit != 0).
func failedEvidence(n int, cmd string) string {
	return fmt.Sprintf("- [x] %d. Attempt %d\n    - cmd: %s\n    - exit: 1\n    - output: \"boom\"\n    - artifact: app/src/a.kt\n", n, n, cmd)
}

func evidenceCodeCount(codes []DiagnosticCode, code DiagnosticCode) int {
	count := 0
	for _, c := range codes {
		if c == code {
			count++
		}
	}
	return count
}

func TestEvidenceRepeatFailureThreshold(t *testing.T) {
	cases := []struct {
		name string
		// changes builds the ticket's Changes section body; returning ""
		// for an n means no Change is rendered for that slot.
		changes func(n int) string
		// wantRepeatFailures is the expected number of
		// evidence-repeat-failure diagnostics for the ticket.
		wantRepeatFailures int
	}{
		{
			name: "three-same-cmd-reports",
			changes: func(n int) string {
				if n <= 3 {
					return failedEvidence(n, "go test ./...")
				}
				return ""
			},
			wantRepeatFailures: 1,
		},
		{
			name: "two-same-cmd-silent",
			changes: func(n int) string {
				if n <= 2 {
					return failedEvidence(n, "go test ./...")
				}
				return ""
			},
			wantRepeatFailures: 0,
		},
		{
			name: "whitespace-variants-count-as-same-cmd",
			changes: func(n int) string {
				switch n {
				case 1:
					return failedEvidence(n, "go test ./...")
				case 2:
					return failedEvidence(n, "`  go  test   ./...`")
				case 3:
					return failedEvidence(n, "go test ./... ")
				default:
					return ""
				}
			},
			wantRepeatFailures: 1,
		},
		{
			name: "different-cmds-do-not-accumulate",
			changes: func(n int) string {
				if n > 4 {
					return ""
				}
				if n%2 == 1 {
					return failedEvidence(n, "go test ./...")
				}
				return failedEvidence(n, "npm test")
			},
			wantRepeatFailures: 0,
		},
		{
			name: "multiple-offending-cmds-report-one-each",
			changes: func(n int) string {
				if n <= 3 {
					return failedEvidence(n, "go test ./...")
				}
				return failedEvidence(n, "npm test")
			},
			wantRepeatFailures: 2,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			var changes string
			for n := 1; n <= 6; n++ {
				changes += tc.changes(n)
			}
			path := writeEvidenceTicket(t, root, "feature-repeat", "01-repeat.md", ticketBody(changes))
			catalog, err := DiscoverArtifacts(root)
			if err != nil {
				t.Fatal(err)
			}
			codes := evidenceCodes(t, catalog, path)
			if got := evidenceCodeCount(codes, DiagnosticEvidenceRepeatFailure); got != tc.wantRepeatFailures {
				t.Fatalf("expected %d evidence-repeat-failure, got %d (codes: %v)", tc.wantRepeatFailures, got, codes)
			}
		})
	}
}

func TestEvidenceRepeatFailureCoexistsWithExitNonzero(t *testing.T) {
	root := t.TempDir()
	var changes string
	for n := 1; n <= 3; n++ {
		changes += failedEvidence(n, "go test ./...")
	}
	path := writeEvidenceTicket(t, root, "feature-coexist", "01-coexist.md", ticketBody(changes))
	catalog, err := DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	codes := evidenceCodes(t, catalog, path)
	if evidenceCodeCount(codes, DiagnosticEvidenceExitNonzero) != 3 {
		t.Fatalf("expected 3 evidence-exit-nonzero (one per failed Change), got %v", codes)
	}
	if evidenceCodeCount(codes, DiagnosticEvidenceRepeatFailure) != 1 {
		t.Fatalf("expected 1 evidence-repeat-failure alongside exit-nonzero, got %v", codes)
	}
	for _, d := range catalog.Diagnostics {
		if d.Artifact == path && d.Code == DiagnosticEvidenceRepeatFailure {
			if d.Severity != SeverityWarning {
				t.Fatalf("evidence-repeat-failure must be a warning, got %s", d.Severity)
			}
			if !strings.Contains(d.Message, "go test ./...") {
				t.Fatalf("diagnostic message should name the offending command, got: %s", d.Message)
			}
		}
	}
}

func TestEvidenceRepeatFailureExemptProposals(t *testing.T) {
	root := t.TempDir()
	var changes string
	for n := 1; n <= 3; n++ {
		changes += failedEvidence(n, "go test ./...")
	}
	exempted := writeEvidenceTicket(t, root, "legacy-proposal", "01-old.md", ticketBody(changes))
	kept := writeEvidenceTicket(t, root, "fresh-proposal", "01-new.md", ticketBody(changes))

	catalog, err := DiscoverArtifactsWithConfig(root, Options{ExemptProposals: []string{"legacy-proposal"}})
	if err != nil {
		t.Fatal(err)
	}
	if codes := evidenceCodes(t, catalog, exempted); len(codes) != 0 {
		t.Fatalf("exempt proposal must produce no evidence diagnostics, got %v", codes)
	}
	if got := evidenceCodeCount(evidenceCodes(t, catalog, kept), DiagnosticEvidenceRepeatFailure); got != 1 {
		t.Fatalf("non-exempt proposal must report evidence-repeat-failure, got %d", got)
	}
}
