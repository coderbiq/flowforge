package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Create proposal lifecycle log cards",
	}

	cmd.AddCommand(newLogCreateCmd())

	return cmd
}

func newLogCreateCmd() *cobra.Command {
	var (
		kind       string
		title      string
		summary    string
		proposalID string
		forCards   []string
		tags       []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a log card",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("--title is required")
			}
			if err := validateLogKind(kind); err != nil {
				return err
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			resolvedProposalID, err := resolveLogProposalID(store, proposalID, forCards)
			if err != nil {
				return err
			}
			if resolvedProposalID == "" {
				return fmt.Errorf("log requires proposal context; run flowforge proposal use <id> or pass --proposal")
			}
			for _, cardID := range forCards {
				if _, err := store.ReadCard(cardID); err != nil {
					return fmt.Errorf("reading context card %s: %w", cardID, err)
				}
			}

			logCard := core.NewCard(core.CardTypeLog, title)
			logCard.ID = core.GenerateCardID(core.CardTypeLog, proposalTimestamp(resolvedProposalID))
			logCard.Body = renderLogBody(kind, summary)
			logCard.Tags = append([]string{kind}, tags...)
			if len(forCards) == 0 {
				logCard.AddLink("PROP-"+resolvedProposalID, "records")
			}
			for _, cardID := range forCards {
				logCard.AddLink(cardID, "records")
			}

			upsertLinksSection(store, logCard)

			_, err = store.CreateCard(logCard, resolvedProposalID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Created log %s\n", logCard.ID)
			fmt.Fprintf(out, "  Kind: %s\n", kind)
			fmt.Fprintf(out, "  Proposal: %s\n", resolvedProposalID)
			if len(forCards) > 0 {
				fmt.Fprintf(out, "  Records: %s\n", strings.Join(forCards, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&kind, "kind", "progress", "Log kind: progress/bug/finding/knowledge/blocked")
	cmd.Flags().StringVar(&title, "title", "", "Log title")
	cmd.Flags().StringVar(&summary, "summary", "", "Log summary")
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID to associate with")
	cmd.Flags().StringSliceVar(&forCards, "for", nil, "Context card IDs this log records")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Additional tags for the log")

	return cmd
}

func validateLogKind(kind string) error {
	switch kind {
	case "progress", "bug", "finding", "knowledge", "blocked":
		return nil
	default:
		return fmt.Errorf("invalid log kind %q (expected progress/bug/finding/knowledge/blocked)", kind)
	}
}

func resolveLogProposalID(store *core.CardStore, explicitProposalID string, forCards []string) (string, error) {
	if explicitProposalID != "" {
		return explicitProposalID, nil
	}

	for _, cardID := range forCards {
		card, err := store.ReadCard(cardID)
		if err != nil {
			return "", fmt.Errorf("reading context card %s: %w", cardID, err)
		}
		if card.Source != "" {
			return card.Source, nil
		}
	}

	return resolveDefaultProposalID("", core.CardTypeLog)
}

func renderLogBody(kind string, summary string) string {
	var body strings.Builder
	body.WriteString("## Kind\n\n")
	body.WriteString(kind)
	body.WriteString("\n")
	if summary != "" {
		body.WriteString("\n## Summary\n\n")
		body.WriteString(summary)
		body.WriteString("\n")
	}

	return body.String()
}
