package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"flowforge/internal/tracker"
)

var (
	checkDir  string
	checkJSON bool
)

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Validate issue dependency graph for cycles, deadlocks, and dangling links",
		Long: `Scans .scratch/ issues, constructs the dependency DAG, and detects:
1. Circular dependencies (deadlocks)
2. Dangling references (blocked by non-existent tickets)
3. Self-dependencies`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := checkDir
			if dir == "" {
				dir = ".scratch"
			}

			issues, err := tracker.DiscoverIssues(dir)
			if err != nil {
				return fmt.Errorf("discovering issues in %s: %w", dir, err)
			}

			if len(issues) == 0 {
				if checkJSON {
					fmt.Println("{\"valid\":true,\"issues_count\":0}")
					return nil
				}
				fmt.Println("No issues found in", dir)
				return nil
			}

			g := tracker.BuildGraph(issues)
			result := g.CheckDependencies()

			isValid := !result.HasCycles && len(result.Dangling) == 0 && len(result.SelfBlocked) == 0

			if checkJSON {
				out := map[string]interface{}{
					"valid":        isValid,
					"issues_count": len(issues),
					"cycles":       result.Cycles,
					"dangling":     result.Dangling,
					"self_blocked": result.SelfBlocked,
				}
				data, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				if !isValid {
					os.Exit(1)
				}
				return nil
			}

			fmt.Printf("Checked %d issues in %s\n", len(issues), dir)

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

			os.Exit(1)
			return nil
		},
	}

	cmd.Flags().StringVarP(&checkDir, "dir", "d", ".scratch", "Directory to scan for issues")
	cmd.Flags().BoolVar(&checkJSON, "json", false, "Output in JSON format")

	return cmd
}
