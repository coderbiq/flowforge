package command

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

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
	cmd.AddCommand(newCardSearchCmd())

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

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			card := core.NewCard(ct, title)
			card.Body = body
			card.Tags = tags

			resolvedProposalID, err := resolveDefaultProposalID(proposalID, ct)
			if err != nil {
				return err
			}

			proposalTs := proposalTimestamp(resolvedProposalID)

			if ct == core.CardTypeTask {
				taskType := "i"
				if len(resolvedProposalID) > 0 {
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

			filePath, err := store.CreateCard(card, resolvedProposalID)
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

			store, err := currentCardStore()
			if err != nil {
				return err
			}
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
			store, err := currentCardStore()
			if err != nil {
				return err
			}

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

			out := cmd.OutOrStdout()
			if jsonOut {
				data, _ := json.MarshalIndent(cards, "", "  ")
				fmt.Fprintln(out, string(data))
			} else {
				if len(cards) == 0 {
					fmt.Fprintln(out, "No cards found.")
					return nil
				}

				fmt.Fprintf(out, "Found %d card(s):\n\n", len(cards))
				for _, card := range cards {
					fmt.Fprintf(out, "  %s [%s] %s\n", card.ID, card.Type, card.Title)
					fmt.Fprintf(out, "    Status: %s | Importance: %s\n", card.Status, card.Importance)
					if card.Source != "" {
						fmt.Fprintf(out, "    Proposal: %s\n", card.Source)
					}
					fmt.Fprintln(out)
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

			store, err := currentCardStore()
			if err != nil {
				return err
			}
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

			store, err := currentCardStore()
			if err != nil {
				return err
			}

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

			store, err := currentCardStore()
			if err != nil {
				return err
			}
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

			store, err := currentCardStore()
			if err != nil {
				return err
			}
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

func newCardSearchCmd() *cobra.Command {
	var (
		scope string
		types string
		limit int
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search cards by keyword",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.TrimSpace(args[0])
			if query == "" {
				return fmt.Errorf("query is required")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			cards, err := searchCards(store, query, scope, types, limit)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if len(cards) == 0 {
				fmt.Fprintln(out, "No cards found.")
				return nil
			}

			fmt.Fprintf(out, "Found %d card(s):\n\n", len(cards))
			for _, result := range cards {
				fmt.Fprintf(out, "  %s [%s] %s\n", result.Card.ID, result.Card.Type, result.Card.Title)
				fmt.Fprintf(out, "    Status: %s", result.Card.Status)
				if result.Card.Source != "" {
					fmt.Fprintf(out, " | Proposal: %s", result.Card.Source)
				}
				fmt.Fprintln(out)
				fmt.Fprintf(out, "    Match: %s\n", result.MatchReason)
				fmt.Fprintln(out)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "all", "Search scope (library/workspace/all)")
	cmd.Flags().StringVar(&types, "type", "", "Comma-separated card types to include")
	cmd.Flags().StringVar(&types, "types", "", "Comma-separated card types to include")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")

	return cmd
}

type cardSearchResult struct {
	Card        *core.Card
	MatchReason string
}

func searchCards(store *core.CardStore, query, scope, types string, limit int) ([]cardSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	typeFilter := map[core.CardType]bool{}
	if strings.TrimSpace(types) != "" {
		for _, raw := range strings.Split(types, ",") {
			ct := core.CardType(strings.TrimSpace(raw))
			if !ct.Valid() {
				return nil, fmt.Errorf("invalid card type: %s", raw)
			}
			typeFilter[ct] = true
		}
	}

	dirs, err := searchScopes(store, scope)
	if err != nil {
		return nil, err
	}

	var cards []*core.Card
	seen := map[string]bool{}
	for _, dir := range dirs {
		dirCards, err := store.ListCards(dir)
		if err != nil {
			return nil, err
		}
		for _, card := range dirCards {
			if seen[card.ID] {
				continue
			}
			seen[card.ID] = true
			cards = append(cards, card)
		}
	}

	results := make([]cardSearchResult, 0, len(cards))
	for _, card := range cards {
		if len(typeFilter) > 0 && !typeFilter[card.Type] {
			continue
		}
		if match, ok := matchCardQuery(card, query); ok {
			results = append(results, cardSearchResult{Card: card, MatchReason: match})
		}
	}

	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Card.ID == results[j].Card.ID {
			return results[i].Card.Title < results[j].Card.Title
		}
		return results[i].Card.ID < results[j].Card.ID
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func searchScopes(store *core.CardStore, scope string) ([]string, error) {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "", "all":
		return []string{store.ActiveDir(), store.IntakeDir(), store.LibraryDir()}, nil
	case "workspace":
		return []string{store.ActiveDir(), store.IntakeDir()}, nil
	case "library":
		return []string{store.LibraryDir()}, nil
	default:
		return nil, fmt.Errorf("invalid scope: %s", scope)
	}
}

func matchCardQuery(card *core.Card, query string) (string, bool) {
	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return "", false
	}

	if strings.Contains(strings.ToLower(card.ID), needle) {
		return "matched id", true
	}
	if strings.Contains(strings.ToLower(card.Title), needle) {
		return "matched title", true
	}
	if strings.Contains(strings.ToLower(card.Body), needle) {
		return "matched body", true
	}
	if strings.Contains(strings.ToLower(card.Domain), needle) {
		return "matched domain", true
	}
	for _, tag := range card.Tags {
		if strings.Contains(strings.ToLower(tag), needle) {
			return "matched tag", true
		}
	}
	return "", false
}

func getOutputFormat() string {
	format := os.Getenv("FLOWFORGE_OUTPUT")
	if format == "" {
		format = "text"
	}
	return format
}
