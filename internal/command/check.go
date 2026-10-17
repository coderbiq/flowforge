package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

var (
	checkDir    string
	checkJSON   bool
	checkStrict bool
	checkFeature string
	checkFull    bool
)

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Validate issue dependency graph for cycles, deadlocks, and dangling links",
		Long: `Scans proposal artifacts, constructs the ticket dependency DAG, and detects:
1. Circular dependencies (deadlocks)
2. Dangling references (blocked by non-existent tickets)
3. Self-dependencies
4. Artifact metadata, authority, semantic-link, waiver, scoped-gap, and completion-evidence diagnostics
5. Checked-change evidence quadruple diagnostics (missing/incomplete/non-zero exit/absent artifact)
6. Repeated failure diagnostics (same command reporting non-zero exit 3+ times in one ticket)
7. Blocked evidence diagnostics (open ticket still carrying a "## Blocked evidence" section awaiting refine-ticket consumption)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := checkDir
			if dir == "" {
				var err error
				dir, err = config.ResolveProposalsDir(".")
				if err != nil {
					return err
				}
			}

			catalog, err := discoverProposalCatalog(dir)
			if err != nil {
				return fmt.Errorf("discovering artifacts in %s: %w", dir, err)
			}
			issues := catalog.Tickets

			g := tracker.BuildGraph(issues)
			result := g.CheckDependencies()
			// --feature scopes only the reported diagnostics; the validity
			// computation below stays on the full unwaived set so a scoped
			// view never relaxes the graph gate.
			scopedDiagnostics := filterDiagnosticsByFeature(catalog.Diagnostics, checkFeature)

			isValid := !result.HasCycles && len(result.Dangling) == 0 && len(result.SelfBlocked) == 0
			for _, diagnostic := range catalog.Diagnostics {
				if diagnostic.Waiver != nil {
					continue
				}
				if diagnostic.Severity == tracker.SeverityBlocker || (checkStrict && (diagnostic.Severity == tracker.SeverityGap || diagnostic.Severity == tracker.SeverityWarning)) {
					isValid = false
				}
			}

			if checkJSON {
				var out map[string]interface{}
				if checkFull {
					out = map[string]interface{}{
						"valid":           isValid,
						"issues_count":    len(issues),
						"cycles":          result.Cycles,
						"dangling":        result.Dangling,
						"self_blocked":    result.SelfBlocked,
						"diagnostics":     catalog.Diagnostics,
						"artifacts_count": len(catalog.Artifacts),
					}
				} else {
					out = map[string]interface{}{
						"valid":           isValid,
						"issues_count":    len(issues),
						"cycles":          result.Cycles,
						"dangling":        result.Dangling,
						"self_blocked":    result.SelfBlocked,
						"diagnostics":     summarizeDiagnostics(scopedDiagnostics),
						"artifacts_count": len(catalog.Artifacts),
						"detail":          "pass --full for complete diagnostics",
					}
				}
				// Compact default keeps the agent-facing channel inside
				// MaxAgentOutputBytes; --full keeps the indented legacy form.
				var data []byte
				if checkFull {
					data, err = json.MarshalIndent(out, "", "  ")
				} else {
					data, err = json.Marshal(out)
				}
				if err != nil {
					return err
				}
				cmd.Println(string(data))
				if !isValid {
					return errPolicyViolation
				}
				return nil
			}

			cmd.Printf("Checked %d issues in %s\n", len(issues), dir)
			if checkFeature != "" {
				cmd.Printf("(diagnostics scoped to %s; validity computed globally)\n", checkFeature)
			}
			printDiagnosticsForMode(cmd, scopedDiagnostics, checkFull)

			if isValid {
				cmd.Println("✓ Dependency graph is healthy. No cycles or dangling references found.")
				return nil
			}

			if result.HasCycles {
				cmd.Println("\n❌ ERROR: Circular dependency detected:")
				for i, cycle := range result.Cycles {
					cmd.Printf("  Cycle #%d: %v\n", i+1, cycle)
				}
			}

			if len(result.Dangling) > 0 {
				cmd.Println("\n⚠️ WARNING: Dangling references (blocked by non-existent tickets):")
				for _, d := range result.Dangling {
					cmd.Printf("  Issue %s is blocked by missing #%s\n", d.IssueID, d.BlockerID)
				}
			}

			if len(result.SelfBlocked) > 0 {
				cmd.Println("\n❌ ERROR: Self-dependency detected:")
				for _, s := range result.SelfBlocked {
					cmd.Printf("  Issue %s is blocked by itself\n", s)
				}
			}

			return errPolicyViolation
		},
	}

	cmd.Flags().StringVarP(&checkDir, "dir", "d", "", "Directory to scan for issues (default: <docs_dir>/proposals)")
	cmd.Flags().BoolVar(&checkJSON, "json", false, "Output in JSON format")
	cmd.Flags().BoolVar(&checkStrict, "strict", false, "Make unwaived warnings and gaps fail validation")
	cmd.Flags().StringVar(&checkFeature, "feature", "", "Proposal slug to scope reported diagnostics to (validity stays global)")
	cmd.Flags().BoolVar(&checkFull, "full", false, "Emit the complete legacy output (every diagnostic verbatim)")

	return cmd
}
