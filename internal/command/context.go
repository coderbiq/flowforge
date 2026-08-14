package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Inspect proposal context",
	}

	cmd.AddCommand(newContextProposalCmd())
	cmd.AddCommand(newContextFeatureCmd())
	cmd.AddCommand(newContextPreflightCmd())
	cmd.AddCommand(newContextRiskReviewCmd())

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
