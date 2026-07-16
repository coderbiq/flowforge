package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newCardSplitCmd() *cobra.Command {
	var titles string

	cmd := &cobra.Command{
		Use:   "split <card-id>",
		Short: "Split an oversized FEATURE into parent-child structure",
		Long: `Split a FEATURE card that is too large into a parent container with child features.
The parent retains Design/Constraints/Motivation; Implementation Plan moves to children.

Examples:
  flowforge card split FEAT-xxx --titles "Clone API,子对象复制,前端实现"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]
			if titles == "" {
				return fmt.Errorf("--titles is required (comma-separated child feature titles)")
			}

			childTitles := splitTitles(titles)
			if len(childTitles) < 2 {
				return fmt.Errorf("need at least 2 child features, got %d", len(childTitles))
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			parent, err := store.ReadCard(cardID)
			if err != nil {
				return err
			}
			if parent.Type != core.CardTypeFeature {
				return fmt.Errorf("card %s is not a feature card", cardID)
			}
			if parent.Role == "container" {
				return fmt.Errorf("card %s is already a container feature", cardID)
			}
			if parent.Status != core.CardStatusDesigned && parent.Status != core.CardStatusPlanned {
				return fmt.Errorf("card must be in 'designed' or 'planned' stage to split, current: %s", parent.Status)
			}

			proposalID := parent.Source
			proposalTs := proposalTimestamp(proposalID)
			var childIDs []string

			for _, title := range childTitles {
				child := core.NewCard(core.CardTypeFeature, title)
				child.ID = core.GenerateCardID(core.CardTypeFeature, proposalTs)
				child.Source = proposalID
				child.Status = core.CardStatusDraft
				child.Body = featureTemplateBody(title)
				child.AddLink(parent.ID, "part_of")
				addProposalOwnershipLink(child, proposalID)

				if _, err := store.CreateCard(child, proposalID); err != nil {
					return fmt.Errorf("creating child feature %q: %w", title, err)
				}
				childIDs = append(childIDs, child.ID)
				parent.AddLink(child.ID, "decomposes")
			}

			parent.Role = "container"
			parent.Body = replaceImplementationPlan(parent.Body, childIDs, childTitles)

			if err := store.UpdateCard(parent); err != nil {
				return fmt.Errorf("updating parent card: %w", err)
			}

			out := cmd.OutOrStdout()
			result := struct {
				Parent   string   `json:"parent"`
				Children []string `json:"children"`
			}{
				Parent:   parent.ID,
				Children: childIDs,
			}
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			enc.Encode(result)

			fmt.Fprintf(out, "\nParent %s is now a container feature.\n", parent.ID)
			fmt.Fprintf(out, "Edit child features directly to fill Design and Implementation Plan.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&titles, "titles", "", "Comma-separated child feature titles")
	return cmd
}

func splitTitles(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

func replaceImplementationPlan(body string, childIDs, childTitles []string) string {
	ipIdx := strings.Index(body, "\n## Implementation Plan")
	if ipIdx < 0 {
		return body
	}

	nextSection := strings.Index(body[ipIdx+1:], "\n## ")
	var ipEnd int
	if nextSection < 0 {
		ipEnd = len(body)
	} else {
		ipEnd = ipIdx + 1 + nextSection
	}

	var sf strings.Builder
	sf.WriteString("\n## Sub-Features\n\n")
	for i, id := range childIDs {
		sf.WriteString(fmt.Sprintf("- [%s] %s\n", id, childTitles[i]))
	}

	return body[:ipIdx] + sf.String() + body[ipEnd:]
}
