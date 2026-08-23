package command

import (
	"fmt"

	"flowforge/internal/core"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display project progress, active proposals, and slice status",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			out := cmd.OutOrStdout()
			currentID, _, _ := runtimeStore.CurrentProposalID(projectID)

			proposals, err := store.ListWorkspaceProposals()
			if err != nil {
				return fmt.Errorf("listing proposals: %w", err)
			}

			fmt.Fprintf(out, "Project: %s\n", projectID)
			if currentID != "" {
				fmt.Fprintf(out, "Active Working Context: %s\n", currentID)
			}
			fmt.Fprintln(out)

			if len(proposals) == 0 {
				fmt.Fprintln(out, "No active proposals in workspace.")
				fmt.Fprintln(out, "Run 'flowforge memory init --proposal <id>' to create a new proposal.")
				return nil
			}

			fmt.Fprintln(out, "| Proposal ID | Title | Status | Mode | Slices Progress |")
			fmt.Fprintln(out, "|---|---|---|---|---|")

			for _, p := range proposals {
				slices, _ := core.ParseProposalSlices(p.Readme)
				total := len(slices)
				completed := 0
				for _, s := range slices {
					if s.Completed {
						completed++
					}
				}

				progressStr := "-"
				if total > 0 {
					progressStr = fmt.Sprintf("%d/%d completed", completed, total)
				}

				activeMarker := ""
				if p.ID == currentID {
					activeMarker = " *"
				}

				fmt.Fprintf(out, "| %s%s | %s | %s | %s | %s |\n",
					p.ID, activeMarker, escapeTableCell(p.Title), p.Status, p.Mode, progressStr)
			}

			// If current proposal has slices, display slice breakdown
			if currentID != "" {
				for _, p := range proposals {
					if p.ID == currentID {
						slices, _ := core.ParseProposalSlices(p.Readme)
						if len(slices) > 0 {
							fmt.Fprintf(out, "\n=== Slices Breakdown for %s ===\n", currentID)
							for _, s := range slices {
								statusIcon := "[ ]"
								if s.Completed {
									statusIcon = "[x]"
								}
								fmt.Fprintf(out, "%s %s: %s\n", statusIcon, s.Index, s.Title)
								if s.TestCommand != "" {
									fmt.Fprintf(out, "    Test: %s\n", s.TestCommand)
								}
							}
						}
						break
					}
				}
			}

			return nil
		},
	}

	return cmd
}
