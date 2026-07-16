package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

var initAllowedTypes = map[core.CardType]bool{
	core.CardTypeFeature:    true,
	core.CardTypeConvention: true,
	core.CardTypeDecision:   true,
	core.CardTypeModule:     true,
	core.CardTypeFinding:    true,
}

func newCardInitCmd() *cobra.Command {
	var (
		cardType   string
		title      string
		proposalID string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a card skeleton for direct editing",
		Long: `Create a new card with generated ID, frontmatter, and template skeleton.
Returns the card ID and file path for subsequent direct editing.

Supported types: feature, convention, decision, module, finding

Examples:
  flowforge card init --type feature --title "FileProcessor Clone" --proposal CR26070801
  flowforge card init --type convention --title "clone pattern"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cardType == "" {
				return fmt.Errorf("--type is required")
			}
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			ct := core.CardType(cardType)
			if !initAllowedTypes[ct] {
				return fmt.Errorf("card init does not support type %q; use: feature, convention, decision, module, finding", cardType)
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			card := core.NewCard(ct, title)

			resolvedProposalID, err := resolveDefaultProposalID(proposalID, ct)
			if err != nil {
				return err
			}

			proposalTs := proposalTimestamp(resolvedProposalID)
			card.ID = core.GenerateCardID(ct, proposalTs)
			card.Source = resolvedProposalID
			card.Body = templateBody(ct, title)
			addProposalOwnershipLink(card, resolvedProposalID)

			savedPath, err := store.CreateCard(card, resolvedProposalID)
			if err != nil {
				return err
			}
			card.FilePath = savedPath

			out := cmd.OutOrStdout()
			result := struct {
				ID   string `json:"id"`
				Path string `json:"path"`
				Type string `json:"type"`
			}{
				ID:   card.ID,
				Path: savedPath,
				Type: string(card.Type),
			}
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		},
	}

	cmd.Flags().StringVar(&cardType, "type", "", "Card type (feature/convention/decision/module/finding)")
	cmd.Flags().StringVar(&title, "title", "", "Card title")
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID to associate with")

	return cmd
}

func templateBody(ct core.CardType, title string) string {
	switch ct {
	case core.CardTypeFeature:
		return featureTemplateBody(title)
	case core.CardTypeConvention:
		return conventionTemplateBody(title)
	case core.CardTypeDecision:
		return decisionTemplateBody(title)
	case core.CardTypeModule:
		return moduleTemplateBody(title)
	case core.CardTypeFinding:
		return findingTemplateBody(title)
	default:
		return ""
	}
}

func featureTemplateBody(title string) string {
	sections := []struct{ heading, placeholder string }{
		{"Summary", "1-3 sentences: what user problem does this feature solve?"},
		{"Motivation", "Why is this needed? What happens if we don't do it?"},
		{"Design", ""},
		{"Constraints", ""},
		{"Implementation Plan", ""},
		{"Verification", ""},
		{"History", ""},
		{"Open Questions", ""},
		{"Dependencies", ""},
	}
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	for _, s := range sections {
		b.WriteString("## " + s.heading + "\n\n")
		if s.heading == "Design" {
			b.WriteString("### Key Decisions\n\n<!-- TBD -->\n\n")
			b.WriteString("### Architecture\n\n<!-- TBD -->\n\n")
			b.WriteString("### Alternatives Considered\n\n<!-- TBD -->\n\n")
		} else if s.placeholder != "" {
			b.WriteString("<!-- " + s.placeholder + " -->\n\n")
		} else {
			b.WriteString("<!-- TBD -->\n\n")
		}
	}
	return strings.TrimSpace(b.String())
}

func conventionTemplateBody(title string) string {
	sections := []string{"Rule", "Rationale", "Applies To", "Examples"}
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	for _, s := range sections {
		b.WriteString("## " + s + "\n\n<!-- TBD -->\n\n")
	}
	return strings.TrimSpace(b.String())
}

func decisionTemplateBody(title string) string {
	sections := []string{"Context", "Decision", "Rationale", "Consequences", "Alternatives"}
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	for _, s := range sections {
		b.WriteString("## " + s + "\n\n<!-- TBD -->\n\n")
	}
	return strings.TrimSpace(b.String())
}

func moduleTemplateBody(title string) string {
	sections := []string{"Purpose", "Responsibilities", "Dependencies", "Public Interface"}
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	for _, s := range sections {
		b.WriteString("## " + s + "\n\n<!-- TBD -->\n\n")
	}
	return strings.TrimSpace(b.String())
}

func findingTemplateBody(title string) string {
	sections := []string{"Summary", "Source", "Evidence", "Impact", "Open Questions"}
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	for _, s := range sections {
		b.WriteString("## " + s + "\n\n<!-- TBD -->\n\n")
	}
	return strings.TrimSpace(b.String())
}
