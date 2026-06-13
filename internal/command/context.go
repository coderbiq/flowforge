package command

import (
	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Inspect proposal context",
	}

	cmd.AddCommand(newContextProposalCmd())

	return cmd
}

func newContextProposalCmd() *cobra.Command {
	var (
		proposalID string
		cardID     string
		taskID     string
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

			report, err := buildProposalContextReport(store, resolvedProposalID, cardID, taskID)
			if err != nil {
				return err
			}

			return renderProposalContextReport(cmd.OutOrStdout(), report)
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID")
	cmd.Flags().StringVar(&cardID, "cards", "", "Focus card ID")
	cmd.Flags().StringVar(&taskID, "task", "", "Focus task ID")

	return cmd
}

func focusCardFromFlags(report *proposalSnapshot, cardID string, taskID string) *core.Card {
	if report == nil {
		return nil
	}
	if taskID != "" {
		if card, ok := report.cardByID[taskID]; ok {
			return card
		}
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
	for _, card := range report.cards {
		if isAnalysisTask(card) && isActiveTaskStatus(card.Status) {
			return card
		}
	}
	return nil
}
