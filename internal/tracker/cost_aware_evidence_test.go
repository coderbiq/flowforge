package tracker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// A shared real run may cover two Changes. A new checked Fix still needs
// its own quadruple; an optional prose schedule cannot waive the old gate.
func TestCostAwareSharedRunPreservesEvidenceGate(t *testing.T) {
	for _, schedule := range []bool{false, true} {
		for _, fixMissing := range []bool{false, true} {
			t.Run(fmt.Sprintf("schedule=%t/fixMissing=%t", schedule, fixMissing), func(t *testing.T) {
				root := t.TempDir()
				if err := os.WriteFile(filepath.Join(root, "run.log"), []byte("group check passed\n"), 0644); err != nil {
					t.Fatal(err)
				}
				quadruple := "    - cmd: go test ./...\n    - exit: 0\n    - output: \"group check passed; shared current run\"\n    - artifact: run.log\n"
				changes := "- [x] 1. Update interface\n" + quadruple + "- [x] 2. Update caller\n" + quadruple
				if fixMissing {
					changes += "- [x] 3. Fix: preserve order\n"
				}
				body := ticketBody(changes)
				if schedule {
					body += "\n## Done and verify\n\n### Verification schedule\n\nG1: Changes 1–2, after both edits, go test ./..., expect exit 0, run.log.\n"
				}
				path := writeEvidenceTicket(t, root, "group", "01-group.md", body)
				catalog, err := DiscoverArtifacts(root)
				if err != nil {
					t.Fatal(err)
				}
				codes := evidenceCodes(t, catalog, path)
				if fixMissing {
					if len(codes) != 1 || codes[0] != DiagnosticEvidenceMissing {
						t.Fatalf("new checked Fix without evidence: got %v", codes)
					}
				} else if len(codes) != 0 {
					t.Fatalf("shared run must remain format-compatible: got %v", codes)
				}
			})
		}
	}
}
