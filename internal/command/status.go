package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

var (
	statusDir string
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display progress overview of all features and tickets in proposals",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := statusDir
			if dir == "" {
				var err error
				dir, err = config.ResolveProposalsDir(".")
				if err != nil {
					return err
				}
			}

			issues, err := tracker.DiscoverIssues(dir)
			if err != nil {
				return fmt.Errorf("discovering issues: %w", err)
			}

			if len(issues) == 0 {
				cmd.Println("No features or issues found in", dir)
				return nil
			}

			g := tracker.BuildGraph(issues)
			frontier := g.ComputeFrontier()

			// Group by feature
			byFeature := make(map[string][]*tracker.Issue)
			for _, issue := range issues {
				f := issue.Feature
				if f == "" {
					f = "(root)"
				}
				byFeature[f] = append(byFeature[f], issue)
			}

			totalResolved := 0
			for _, issue := range issues {
				if issue.Status.IsTerminal() {
					totalResolved++
				}
			}

			pct := 0
			if len(issues) > 0 {
				pct = (totalResolved * 100) / len(issues)
			}

			cmd.Printf("FlowForge Local Tracker Status: %d/%d resolved (%d%%)\n", totalResolved, len(issues), pct)
			cmd.Printf("Frontier Ready: %d | In Progress: %d | Blocked: %d\n\n", len(frontier.Ready), len(frontier.Claimed), len(frontier.Blocked))

			for feat, list := range byFeature {
				resCount := 0
				for _, it := range list {
					if it.Status.IsTerminal() {
						resCount++
					}
				}
				cmd.Printf("• Feature: %s [%d/%d done]\n", feat, resCount, len(list))
				for _, it := range list {
					icon := "⚪"
					if it.Status.IsTerminal() {
						icon = "✅"
					} else if it.Status == tracker.StatusClaimed {
						icon = "⏳"
					} else {
						// check if ready or blocked
						isReady := false
						for _, r := range frontier.Ready {
							if r.ID == it.ID && r.Feature == it.Feature {
								isReady = true
								break
							}
						}
						if isReady {
							icon = "🟢"
						} else {
							icon = "⛔"
						}
					}
					cmd.Printf("    %s #%s %s (%s)\n", icon, it.ID, it.Title, it.Status)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&statusDir, "dir", "d", "", "Directory to scan for issues (default: <docs_dir>/proposals)")
	return cmd
}
