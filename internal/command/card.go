package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/core"
)

func newCardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Manage knowledge cards",
		Long:  "Create, read, update, delete, and list knowledge cards.",
	}

	cmd.AddCommand(newCardCreateCmd())
	cmd.AddCommand(newCardReadCmd())
	cmd.AddCommand(newCardListCmd())
	cmd.AddCommand(newCardUpdateCmd())
	cmd.AddCommand(newCardDeleteCmd())
	cmd.AddCommand(newCardRelatedCmd())
	cmd.AddCommand(newCardDependentsCmd())

	return cmd
}

func newCardCreateCmd() *cobra.Command {
	var (
		cardType   string
		title      string
		body       string
		proposalID string
		links      []string
		tags       []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new card",
		Long: `Create a new knowledge card.

Card types: requirement, decision, design, task, log, convention, finding, module, structure

Examples:
  flowforge card create --type requirement --title "User login feature"
  flowforge card create --type decision --title "Use PostgreSQL" --proposal CR24010101
  flowforge card create --type task --title "Implement API" --links "DEC-abc123"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cardType == "" {
				return fmt.Errorf("--type is required")
			}
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			ct := core.CardType(cardType)
			if !ct.Valid() {
				return fmt.Errorf("invalid card type: %s", cardType)
			}

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))

			card := core.NewCard(ct, title)
			card.Body = body
			card.Tags = tags

			var proposalTs string
			if proposalID != "" {
				parts := strings.Split(proposalID, "-")
				if len(parts) >= 2 {
					proposalTs = parts[1]
				}
			}

			if ct == core.CardTypeTask {
				taskType := "i"
				if len(proposalID) > 0 {
					card.ID = core.GenerateTaskID(proposalTs, taskType)
				} else {
					card.ID = core.GenerateTaskID("", taskType)
				}
			} else {
				card.ID = core.GenerateCardID(ct, proposalTs)
			}

			for _, linkStr := range links {
				parts := strings.Split(linkStr, ":")
				target := parts[0]
				relation := "related"
				if len(parts) > 1 {
					relation = parts[1]
				}
				card.AddLink(target, relation)
			}

			filePath, err := store.CreateCard(card, proposalID)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Created card %s\n", card.ID)
			fmt.Printf("  Type: %s\n", card.Type)
			fmt.Printf("  Title: %s\n", card.Title)
			fmt.Printf("  File: %s\n", filePath)

			return nil
		},
	}

	cmd.Flags().StringVar(&cardType, "type", "", "Card type (requirement/decision/design/task/log/convention/finding/module/structure)")
	cmd.Flags().StringVar(&title, "title", "", "Card title")
	cmd.Flags().StringVar(&body, "body", "", "Card body content")
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID to associate with")
	cmd.Flags().StringSliceVar(&links, "links", nil, "Links to other cards (format: CARD_ID or CARD_ID:relation)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags for the card")

	return cmd
}

func newCardReadCmd() *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "read <card-id>",
		Short: "Read a card's content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))
			card, err := store.ReadCard(cardID)
			if err != nil {
				return err
			}

			if outputJSON {
				data, _ := json.MarshalIndent(card, "", "  ")
				fmt.Println(string(data))
			} else {
				fmt.Printf("ID: %s\n", card.ID)
				fmt.Printf("Type: %s\n", card.Type)
				fmt.Printf("Title: %s\n", card.Title)
				fmt.Printf("Status: %s\n", card.Status)
				fmt.Printf("Importance: %s\n", card.Importance)
				if len(card.Tags) > 0 {
					fmt.Printf("Tags: %s\n", strings.Join(card.Tags, ", "))
				}
				if len(card.Links) > 0 {
					fmt.Println("Links:")
					for _, link := range card.Links {
						fmt.Printf("  - %s (%s)\n", link.Target, link.Relation)
					}
				}
				fmt.Printf("Created: %s\n", card.Created.Format("2006-01-02 15:04:05"))
				fmt.Printf("Updated: %s\n", card.Updated.Format("2006-01-02 15:04:05"))
				fmt.Printf("File: %s\n", card.FilePath)
				if card.Body != "" {
					fmt.Println("\n--- Body ---")
					fmt.Println(card.Body)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	return cmd
}

func newCardListCmd() *cobra.Command {
	var (
		cardType string
		status   string
		proposal string
		jsonOut  bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cards",
		Long: `List cards with optional filters.

Examples:
  flowforge card list
  flowforge card list --type task
  flowforge card list --type requirement --status draft
  flowforge card list --proposal CR24010101
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))

			var cards []*core.Card

			if proposal != "" {
				cardsDir := store.ProposalCardsDir(proposal)
				cards, err = store.ListCards(cardsDir)
				if err != nil {
					return err
				}
			} else if cardType != "" {
				ct := core.CardType(cardType)
				if !ct.Valid() {
					return fmt.Errorf("invalid card type: %s", cardType)
				}
				cards, err = store.ListCardsByType(ct)
				if err != nil {
					return err
				}
			} else {
				allDirs := []string{
					store.ActiveDir(),
					store.IntakeDir(),
					store.LibraryDir(),
				}
				for _, dir := range allDirs {
					dirCards, _ := store.ListCards(dir)
					cards = append(cards, dirCards...)
				}
			}

			if status != "" {
				var filtered []*core.Card
				for _, card := range cards {
					if string(card.Status) == status {
						filtered = append(filtered, card)
					}
				}
				cards = filtered
			}

			if jsonOut {
				data, _ := json.MarshalIndent(cards, "", "  ")
				fmt.Println(string(data))
			} else {
				if len(cards) == 0 {
					fmt.Println("No cards found.")
					return nil
				}

				fmt.Printf("Found %d card(s):\n\n", len(cards))
				for _, card := range cards {
					fmt.Printf("  %s [%s] %s\n", card.ID, card.Type, card.Title)
					fmt.Printf("    Status: %s | Importance: %s\n", card.Status, card.Importance)
					if card.Source != "" {
						fmt.Printf("    Proposal: %s\n", card.Source)
					}
					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&cardType, "type", "", "Filter by card type")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&proposal, "proposal", "", "Filter by proposal ID")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}

func newCardUpdateCmd() *cobra.Command {
	var (
		title       string
		status      string
		importance  string
		body        string
		addLinks    []string
		removeLinks []string
	)

	cmd := &cobra.Command{
		Use:   "update <card-id>",
		Short: "Update a card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))
			card, err := store.ReadCard(cardID)
			if err != nil {
				return err
			}

			if title != "" {
				card.Title = title
			}
			if status != "" {
				card.Status = core.CardStatus(status)
			}
			if importance != "" {
				card.Importance = core.Importance(importance)
			}
			if body != "" {
				card.Body = body
			}

			for _, linkStr := range addLinks {
				parts := strings.Split(linkStr, ":")
				target := parts[0]
				relation := "related"
				if len(parts) > 1 {
					relation = parts[1]
				}
				card.AddLink(target, relation)
			}

			for _, linkStr := range removeLinks {
				parts := strings.Split(linkStr, ":")
				target := parts[0]
				relation := "related"
				if len(parts) > 1 {
					relation = parts[1]
				}
				card.RemoveLink(target, relation)
			}

			if err := store.UpdateCard(card); err != nil {
				return err
			}

			fmt.Printf("✓ Updated card %s\n", card.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().StringVar(&status, "status", "", "New status (draft/active/accepted/deprecated/superseded)")
	cmd.Flags().StringVar(&importance, "importance", "", "New importance (must/should/may)")
	cmd.Flags().StringVar(&body, "body", "", "New body content")
	cmd.Flags().StringSliceVar(&addLinks, "add-link", nil, "Add link (format: CARD_ID or CARD_ID:relation)")
	cmd.Flags().StringSliceVar(&removeLinks, "remove-link", nil, "Remove link (format: CARD_ID or CARD_ID:relation)")

	return cmd
}

func newCardDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <card-id>",
		Short: "Delete a card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))

			if !force {
				card, err := store.ReadCard(cardID)
				if err != nil {
					return err
				}
				if card.Status != core.CardStatusDraft {
					return fmt.Errorf("only draft cards can be deleted (current status: %s). Use --force to override", card.Status)
				}

				fmt.Printf("Delete card %s [%s] %s? [y/N] ", card.ID, card.Type, card.Title)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Aborted.")
					return nil
				}
			}

			if err := store.DeleteCard(cardID); err != nil {
				return err
			}

			fmt.Printf("✓ Deleted card %s\n", cardID)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force delete even if not draft")

	return cmd
}

func newCardRelatedCmd() *cobra.Command {
	var (
		relation string
		depth    int
	)

	cmd := &cobra.Command{
		Use:   "related <card-id>",
		Short: "Show related cards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))
			related, err := store.GetRelated(cardID, relation, depth)
			if err != nil {
				return err
			}

			if len(related) == 0 {
				fmt.Println("No related cards found.")
				return nil
			}

			fmt.Printf("Related cards for %s:\n\n", cardID)
			for _, card := range related {
				fmt.Printf("  %s [%s] %s\n", card.ID, card.Type, card.Title)
				fmt.Printf("    Status: %s\n", card.Status)
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&relation, "relation", "", "Filter by relation type")
	cmd.Flags().IntVar(&depth, "depth", 1, "Traversal depth")

	return cmd
}

func newCardDependentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dependents <card-id>",
		Short: "Show cards that depend on this card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			store := core.NewCardStore(cfg.WikiRoot(projectRoot))
			dependents, err := store.GetDependents(cardID)
			if err != nil {
				return err
			}

			if len(dependents) == 0 {
				fmt.Println("No cards depend on this card.")
				return nil
			}

			fmt.Printf("Cards depending on %s:\n\n", cardID)
			for _, card := range dependents {
				fmt.Printf("  %s [%s] %s\n", card.ID, card.Type, card.Title)
				fmt.Printf("    Status: %s\n", card.Status)
				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}

func getOutputFormat() string {
	format := os.Getenv("FLOWFORGE_OUTPUT")
	if format == "" {
		format = "text"
	}
	return format
}
