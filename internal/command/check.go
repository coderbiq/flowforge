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
)

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Validate issue dependency graph for cycles, deadlocks, and dangling links",
		Long: `Scans proposal issues, constructs the dependency DAG, and detects:
1. Circular dependencies (deadlocks)
2. Dangling references (blocked by non-existent tickets)
3. Self-dependencies`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := checkDir
			if dir == "" {
				dir = config.ResolveProposalsDir(".")
			}

			catalog, err := tracker.DiscoverArtifacts(dir)
			if err != nil {
				return fmt.Errorf("discovering artifacts in %s: %w", dir, err)
			}
			issues := catalog.Tickets

			g := tracker.BuildGraph(issues)
			result := g.CheckDependencies()

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
				out := map[string]interface{}{
					"valid":           isValid,
					"issues_count":    len(issues),
					"cycles":          result.Cycles,
					"dangling":        result.Dangling,
					"self_blocked":    result.SelfBlocked,
					"diagnostics":     catalog.Diagnostics,
					"artifacts_count": len(catalog.Artifacts),
				}
				data, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(data))
				if !isValid {
					return errPolicyViolation
				}
				return nil
			}

			fmt.Printf("Checked %d issues in %s\n", len(issues), dir)
			printDiagnostics(cmd, catalog.Diagnostics)

			if isValid {
				fmt.Println("✓ Dependency graph is healthy. No cycles or dangling references found.")
				return nil
			}

			if result.HasCycles {
				fmt.Println("\n❌ ERROR: Circular dependency detected:")
				for i, cycle := range result.Cycles {
					fmt.Printf("  Cycle #%d: %v\n", i+1, cycle)
				}
			}

			if len(result.Dangling) > 0 {
				fmt.Println("\n⚠️ WARNING: Dangling references (blocked by non-existent tickets):")
				for _, d := range result.Dangling {
					fmt.Printf("  Issue %s is blocked by missing #%s\n", d.IssueID, d.BlockerID)
				}
			}

			if len(result.SelfBlocked) > 0 {
				fmt.Println("\n❌ ERROR: Self-dependency detected:")
				for _, s := range result.SelfBlocked {
					fmt.Printf("  Issue %s is blocked by itself\n", s)
				}
			}

			return errPolicyViolation
		},
	}

	cmd.Flags().StringVarP(&checkDir, "dir", "d", "", "Directory to scan for issues (default: <docs_dir>/proposals)")
	cmd.Flags().BoolVar(&checkJSON, "json", false, "Output in JSON format")
	cmd.Flags().BoolVar(&checkStrict, "strict", false, "Make unwaived warnings and gaps fail validation")

	return cmd
}
