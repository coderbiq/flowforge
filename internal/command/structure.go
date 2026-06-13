package command

import (
	"fmt"
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

			structureCard, err := store.ReadCard(structureID)
			if err != nil {
				return err
			}
			if structureCard.Type != core.CardTypeStructure {
				return fmt.Errorf("card %s is not a structure card (type: %s)", structureID, structureCard.Type)
			}

			if _, err := store.ReadCard(cardID); err != nil {
				return err
			}

			before := len(structureIndexedCardIDs(structureCard))
			structureCard.AddLink(cardID, "indexes")
			if err := store.UpdateCard(structureCard); err != nil {
				return err
			}
			indexedCount := len(structureIndexedCardIDs(structureCard))

			out := cmd.OutOrStdout()
			if indexedCount == before {
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

			structureCard, err := store.ReadCard(structureID)
			if err != nil {
				return err
			}
			if structureCard.Type != core.CardTypeStructure {
				return fmt.Errorf("card %s is not a structure card (type: %s)", structureID, structureCard.Type)
			}

			removed := structureCard.RemoveLink(cardID, "indexes")
			if !removed {
				fmt.Fprintf(cmd.OutOrStdout(), "No change: %s does not index %s\n", structureID, cardID)
				return nil
			}

			if err := store.UpdateCard(structureCard); err != nil {
				return err
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
