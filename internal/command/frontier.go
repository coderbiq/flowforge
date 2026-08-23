package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/tracker"
)

var (
	frontierDir   string
	frontierJSON  bool
	frontierQuiet bool
)

func newFrontierCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontier",
		Short: "Compute unblocked, ready-to-execute tickets from .scratch/",
		Long: `Scans .scratch/ for issues and tickets, evaluates DAG dependencies,
and outputs the exact list of unblocked, executable tasks.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := frontierDir
			if dir == "" {
				dir = ".scratch"
			}

			issues, err := tracker.DiscoverIssues(dir)
			if err != nil {
				return fmt.Errorf("discovering issues in %s: %w", dir, err)
			}

			if len(issues) == 0 {
				if frontierJSON {
					fmt.Println("{\"ready\":[],\"claimed\":[],\"blocked\":[]}")
					return nil
				}
				fmt.Println("No issues found in", dir)
				return nil
			}

			g := tracker.BuildGraph(issues)
			frontier := g.ComputeFrontier()

			if frontierJSON {
				data, err := json.MarshalIndent(frontier, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if frontierQuiet {
				for _, issue := range frontier.Ready {
					fmt.Println(issue.FilePath)
				}
				return nil
			}

			if len(frontier.Ready) == 0 && len(frontier.Claimed) == 0 && len(frontier.Blocked) == 0 {
				fmt.Println("All tasks in", dir, "are resolved/closed.")
				return nil
			}

			if len(frontier.Ready) > 0 {
				fmt.Println("=== READY (Unblocked & Executable) ===")
				for _, issue := range frontier.Ready {
					featurePrefix := ""
					if issue.Feature != "" {
						featurePrefix = "[" + issue.Feature + "] "
					}
					fmt.Printf("✓ %s#%s %s (%s)\n", featurePrefix, issue.ID, issue.Title, issue.FilePath)
				}
			}

			if len(frontier.Claimed) > 0 {
				fmt.Println("\n=== CLAIMED (In Progress) ===")
				for _, issue := range frontier.Claimed {
					assignee := ""
					if issue.Assignee != "" {
						assignee = " by " + issue.Assignee
					}
					fmt.Printf("⏳ #%s %s%s (%s)\n", issue.ID, issue.Title, assignee, issue.FilePath)
				}
			}

			if len(frontier.Blocked) > 0 {
				fmt.Println("\n=== BLOCKED ===")
				for _, b := range frontier.Blocked {
					fmt.Printf("⛔ #%s %s (Waiting on: %v)\n", b.Issue.ID, b.Issue.Title, b.WaitingOn)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&frontierDir, "dir", "d", ".scratch", "Directory to scan for issues")
	cmd.Flags().BoolVar(&frontierJSON, "json", false, "Output in JSON format")
	cmd.Flags().BoolVarP(&frontierQuiet, "quiet", "q", false, "Output only ready file paths")

	return cmd
}
