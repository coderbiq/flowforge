package command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Inspect proposal and slice execution context",
	}

	cmd.AddCommand(newContextSliceCmd())
	cmd.AddCommand(newContextProposalCmd())
	cmd.AddCommand(newContextFeatureCmd())
	cmd.AddCommand(newContextPreflightCmd())
	cmd.AddCommand(newContextRiskReviewCmd())

	return cmd
}

func newContextSliceCmd() *cobra.Command {
	var (
		proposalID string
		sliceIndex string
	)

	cmd := &cobra.Command{
		Use:   "slice",
		Short: "Extract minimal execution context for a single slice",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			if proposalID == "" {
				currentID, ok, _ := runtimeStore.CurrentProposalID(projectID)
				if ok && currentID != "" {
					proposalID = currentID
				}
			}

			if proposalID == "" {
				return fmt.Errorf("no active proposal specified. Use --proposal <id>")
			}

			prop, err := store.FindWorkspaceProposal(proposalID)
			if err != nil {
				return err
			}

			slices, err := core.ParseProposalSlices(prop.Readme)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			var targetSlice *core.WorkspaceSlice
			if sliceIndex == "" {
				// Find first uncompleted slice
				for _, s := range slices {
					if !s.Completed {
						targetSlice = &s
						break
					}
				}
				if targetSlice == nil && len(slices) > 0 {
					targetSlice = &slices[0]
				}
			} else {
				targetIdx := strings.ToLower(strings.TrimSpace(sliceIndex))
				for _, s := range slices {
					sIdx := strings.ToLower(strings.TrimSpace(s.Index))
					if sIdx == targetIdx ||
						strings.EqualFold(s.Index, "slice "+sliceIndex) ||
						strings.EqualFold(s.Index, "batch "+sliceIndex) ||
						strings.HasPrefix(sIdx, "batch "+targetIdx) ||
						strings.HasPrefix(sIdx, "slice "+targetIdx) ||
						strings.HasPrefix(sIdx, targetIdx+":") ||
						strings.HasPrefix(sIdx, targetIdx+".") {
						targetSlice = &s
						break
					}
				}
			}

			if targetSlice == nil {
				return fmt.Errorf("slice %q not found in proposal %s", sliceIndex, proposalID)
			}

			fmt.Fprintf(out, "# Slice Context: %s (Proposal: %s)\n\n", targetSlice.Title, prop.ID)
			fmt.Fprintf(out, "**Status**: %s | **Mode**: %s\n\n", map[bool]string{true: "Completed", false: "Pending"}[targetSlice.Completed], prop.Mode)
			fmt.Fprintln(out, "## Execution Target")
			fmt.Fprintln(out, targetSlice.Raw)
			fmt.Fprintln(out)

			if len(targetSlice.TouchPoints) > 0 {
				fmt.Fprintln(out, "## Touchpoints & Physical Specs")
				for _, tp := range targetSlice.TouchPoints {
					fmt.Fprintf(out, "- `%s`\n", tp)
					// If in hierarchical mode, try to load module doc
					if prop.Mode == "hierarchical" {
						modName := strings.TrimSuffix(filepath.Base(tp), filepath.Ext(tp))
						modPath := filepath.Join(prop.Path, "modules", modName+".md")
						if content, err := os.ReadFile(modPath); err == nil {
							fmt.Fprintf(out, "\n<details>\n<summary>Physical Spec: %s</summary>\n\n%s\n</details>\n", modName, string(content))
						}
					}
				}
				fmt.Fprintln(out)
			}

			if targetSlice.TestCommand != "" {
				fmt.Fprintln(out, "## Automated Test Loop")
				fmt.Fprintf(out, "```bash\n%s\n```\n", targetSlice.TestCommand)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID")
	cmd.Flags().StringVar(&sliceIndex, "slice", "", "Slice index (e.g. 1, 2.1)")

	return cmd
}

func newContextProposalCmd() *cobra.Command {
	var (
		proposalID string
		cardID     string
	)

	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "Show minimal proposal context",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, runtimeStore, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			project, _, err := resolveCurrentProject(cfg, runtimeStore)
			if err != nil {
				return err
			}

			resolvedProposalID := proposalID
			if resolvedProposalID == "" {
				resolvedProposalID, err = currentProposalIDForProject(runtimeStore, project.ID)
				if err != nil {
					return err
				}
			}

			wikiRoot, err := cfg.WikiRootForProject(projectRoot, project.ID)
			if err != nil {
				return err
			}
			store := core.NewCardStore(wikiRoot)

			report, err := buildProposalContextReport(store, resolvedProposalID, cardID)
			if err != nil {
				return err
			}

			return renderProposalContextReport(cmd.OutOrStdout(), report)
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID")
	cmd.Flags().StringVar(&cardID, "cards", "", "Focus card ID")
	return cmd
}
func focusCardFromFlags(report *proposalSnapshot, cardID string) *core.Card {
	if report == nil {
		return nil
	}
	if cardID != "" {
		if card, ok := report.cardByID[cardID]; ok {
			return card
		}
	}
	if report.requirementIndex != nil {
		return report.requirementIndex
	}
	if report.rootCard != nil {
		return report.rootCard
	}
	return nil
}

func renderContextCardTable(out interface{ Write([]byte) (int, error) }, cards []*core.Card) {
	if len(cards) == 0 {
		fmt.Fprintln(out, "- None")
		return
	}
	fmt.Fprintln(out, "| ID | Type | Title | Status | Summary |")
	fmt.Fprintln(out, "|----|------|-------|--------|---------|")
	for _, card := range cards {
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s |\n",
			card.ID,
			card.Type,
			escapeTableCell(card.Title),
			card.Status,
			escapeTableCell(summarizeForContext(card)),
		)
	}
}

func summarizeForContext(card *core.Card) string {
	if card == nil {
		return ""
	}
	body := strings.TrimSpace(card.Body)
	if body == "" {
		return ""
	}
	if summary, _ := firstMeaningfulSectionSummary(body); summary != "" {
		return summary
	}
	return ""
}

func sortCardsForContext(cards []*core.Card) {
	sort.SliceStable(cards, func(i, j int) bool {
		if cards[i].Type != cards[j].Type {
			return cards[i].Type < cards[j].Type
		}
		return cards[i].ID < cards[j].ID
	})
}

func escapeTableCell(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", "\\|")
	return strings.TrimSpace(value)
}
