package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newStructureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "structure",
		Short: "Manage STR index links",
		Long:  "Maintain indexes relations on structure cards.",
	}

	cmd.AddCommand(newStructureAddCmd())
	cmd.AddCommand(newStructureRemoveCmd())
	cmd.AddCommand(newStructureListCmd())
	cmd.AddCommand(newStructureRefreshCmd())

	return cmd
}

func newStructureAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <structure-id> <card-id>",
		Short: "Add an indexed card to a structure",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			structureID := args[0]
			cardID := args[1]

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			var indexedCount int
			var changed bool
			if err := store.UpdateCardWithLock(structureID, func(card *core.Card) error {
				if card.Type != core.CardTypeStructure {
					return fmt.Errorf("card %s is not a structure card (type: %s)", structureID, card.Type)
				}
				indexedCard, err := store.ReadCard(cardID)
				if err != nil {
					return err
				}
				if err := validateStructureIndexedCard(card, indexedCard); err != nil {
					return err
				}
				before := len(structureIndexedCardIDs(card))
				card.AddLink(cardID, "indexes")
				refreshedBody, err := refreshStructureEntriesBody(store, card)
				if err != nil {
					return err
				}
				card.Body = refreshedBody
				indexedCount = len(structureIndexedCardIDs(card))
				changed = indexedCount > before
				return nil
			}); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if !changed {
				fmt.Fprintf(out, "No change: %s already indexes %s\n", structureID, cardID)
				return nil
			}

			fmt.Fprintf(out, "✓ Added %s to %s\n", cardID, structureID)
			fmt.Fprintf(out, "  relation: indexes\n")
			if indexedCount > 15 {
				fmt.Fprintf(out, "  warning: %s now has %d direct indexed cards; consider splitting the structure\n", structureID, indexedCount)
			}

			return nil
		},
	}

	return cmd
}

func newStructureRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <structure-id> <card-id>",
		Short: "Remove an indexed card from a structure",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			structureID := args[0]
			cardID := args[1]

			store, err := currentCardStore()
			if err != nil {
				return err
			}

		var removed bool
		if err := store.UpdateCardWithLock(structureID, func(card *core.Card) error {
			if card.Type != core.CardTypeStructure {
				return fmt.Errorf("card %s is not a structure card (type: %s)", structureID, card.Type)
			}

			if err := guardOrphanCardOnStructureRemove(store, cardID, structureID); err != nil {
				return err
			}

			removed = card.RemoveLink(cardID, "indexes")
			refreshedBody, err := refreshStructureEntriesBody(store, card)
			if err != nil {
				return err
			}
			card.Body = refreshedBody
			return nil
		}); err != nil {
			return err
		}

			if !removed {
				fmt.Fprintf(cmd.OutOrStdout(), "No change: %s does not index %s\n", structureID, cardID)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Removed %s from %s\n", cardID, structureID)
			return nil
		},
	}

	return cmd
}

func newStructureListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <structure-id>",
		Short: "List indexed cards for a structure",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			structureID := args[0]

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			structureCard, err := store.ReadCard(structureID)
			if err != nil {
				return err
			}
			if structureCard.Type != core.CardTypeStructure {
				return fmt.Errorf("card %s is not a structure card (type: %s)", structureID, structureCard.Type)
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Structure: %s\n", structureCard.ID)
			fmt.Fprintf(out, "Title: %s\n", structureCard.Title)

			indexedIDs := structureIndexedCardIDs(structureCard)
			if len(indexedIDs) == 0 {
				fmt.Fprintln(out, "No indexed cards.")
				return nil
			}

			fmt.Fprintf(out, "Indexed cards (%d):\n", len(indexedIDs))
			for _, cardID := range indexedIDs {
				card, err := store.ReadCard(cardID)
				if err != nil {
					fmt.Fprintf(out, "  - %s\n", cardID)
					continue
				}
				fmt.Fprintf(out, "  - %s [%s] %s\n", card.ID, card.Type, card.Title)
			}

			return nil
		},
	}

	return cmd
}

func newStructureRefreshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh <structure-id>",
		Short: "Refresh a structure card's readable entries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			structureID := args[0]

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			var indexedCount int
			if err := store.UpdateCardWithLock(structureID, func(card *core.Card) error {
				if card.Type != core.CardTypeStructure {
					return fmt.Errorf("card %s is not a structure card (type: %s)", structureID, card.Type)
				}
				refreshedBody, err := refreshStructureEntriesBody(store, card)
				if err != nil {
					return err
				}
				card.Body = refreshedBody
				indexedCount = len(structureIndexedCardIDs(card))
				return nil
			}); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Refreshed %s\n", structureID)
			fmt.Fprintf(cmd.OutOrStdout(), "  entries: %d\n", indexedCount)
			return nil
		},
	}

	return cmd
}

func structureIndexedCardIDs(card *core.Card) []string {
	ids := make([]string, 0, len(card.Links))
	for _, link := range card.Links {
		if link.Relation != "indexes" {
			continue
		}
		ids = append(ids, strings.TrimSpace(link.Target))
	}
	return ids
}

func refreshStructureEntriesBody(store *core.CardStore, card *core.Card) (string, error) {
	entries, err := renderStructureEntries(store, card)
	if err != nil {
		return "", err
	}
	body := strings.TrimSpace(card.Body)
	if body == "" {
		body = "# " + card.Title + "\n\n## Purpose\n\nStructure index."
	}
	return upsertMarkdownSection(body, "Entries", entries), nil
}

func renderStructureEntries(store *core.CardStore, card *core.Card) (string, error) {
	indexedIDs := structureIndexedCardIDs(card)
	if len(indexedIDs) == 0 {
		return "- None", nil
	}

	var lines []string
	for _, cardID := range indexedIDs {
		linkedCard, err := store.ReadCard(cardID)
		if err != nil {
			return "", fmt.Errorf("indexed card %s could not be read: %w", cardID, err)
		}
		if err := validateStructureIndexedCard(card, linkedCard); err != nil {
			return "", err
		}
		target, err := markdownLinkTarget(card, linkedCard)
		if err != nil {
			return "", err
		}
		lines = append(lines, fmt.Sprintf("- [%s](%s) (%s, %s) - %s", linkedCard.ID, target, linkedCard.Type, linkedCard.Status, linkedCard.Title))
	}
	return strings.Join(lines, "\n"), nil
}

func validateStructureIndexedCard(structureCard, indexedCard *core.Card) error {
	if strings.HasPrefix(structureCard.ID, "STR-") && strings.Contains(structureCard.ID, "-REQ") {
		if indexedCard.Type != core.CardTypeRequirement && indexedCard.Type != core.CardTypeStructure {
			return fmt.Errorf("proposal requirement index %s can only index requirement or structure cards, got %s (%s)", structureCard.ID, indexedCard.Type, indexedCard.ID)
		}
	}
	return nil
}

func markdownLinkTarget(fromCard, toCard *core.Card) (string, error) {
	if fromCard.FilePath == "" || toCard.FilePath == "" {
		return "", fmt.Errorf("cannot render markdown link without file paths")
	}
	rel, err := filepath.Rel(filepath.Dir(fromCard.FilePath), toCard.FilePath)
	if err != nil {
		return "", fmt.Errorf("computing relative link: %w", err)
	}
	return filepath.ToSlash(rel), nil
}

func upsertMarkdownSection(body string, section string, content string) string {
	heading := "## " + section
	replacement := heading + "\n\n" + strings.TrimSpace(content)
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return replacement + "\n"
	}

	idx := strings.Index(trimmed, heading)
	if idx >= 0 {
		before := strings.TrimRight(trimmed[:idx], "\n")
		afterStart := idx + len(heading)
		after := trimmed[afterStart:]
		next := strings.Index(after, "\n## ")
		if next >= 0 {
			after = after[next:]
		} else {
			after = ""
		}
		parts := []string{}
		if before != "" {
			parts = append(parts, before)
		}
		parts = append(parts, replacement)
		if strings.TrimSpace(after) != "" {
			parts = append(parts, strings.TrimLeft(after, "\n"))
		}
		return strings.Join(parts, "\n\n") + "\n"
	}

	openQuestions := "\n## Open Questions"
	if idx := strings.Index(trimmed, openQuestions); idx >= 0 {
		before := strings.TrimRight(trimmed[:idx], "\n")
		after := strings.TrimLeft(trimmed[idx:], "\n")
		return before + "\n\n" + replacement + "\n\n" + after + "\n"
	}

	return trimmed + "\n\n" + replacement + "\n"
}

func guardOrphanCardOnStructureRemove(store *core.CardStore, cardID, structureID string) error {
	dependents, err := store.GetDependents(cardID)
	if err != nil {
		return fmt.Errorf("checking card dependents: %w", err)
	}

	hasOtherStructureIndex := false
	for _, dep := range dependents {
		if dep.ID == structureID {
			continue
		}
		if dep.Type != core.CardTypeStructure {
			continue
		}
		for _, link := range dep.Links {
			if link.Target == cardID && link.Relation == "indexes" {
				hasOtherStructureIndex = true
				break
			}
		}
		if hasOtherStructureIndex {
			break
		}
	}

	if !hasOtherStructureIndex {
		return fmt.Errorf(
			"removing %s from %s would leave the card without any structure index; use `card delete %s --force` to remove it permanently, or add it to another structure first",
			cardID, structureID, cardID,
		)
	}

	return nil
}
