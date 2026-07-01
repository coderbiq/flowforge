package command

import (
	"encoding/json"
	"fmt"
	"io"
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
	cmd.AddCommand(newCardRefreshCmd())
	cmd.AddCommand(newCardDeleteCmd())
	cmd.AddCommand(newCardRelatedCmd())
	cmd.AddCommand(newCardDependentsCmd())
	cmd.AddCommand(newCardLinkCmd())
	cmd.AddCommand(newCardUnlinkCmd())
	cmd.AddCommand(newCardSearchCmd())
	cmd.AddCommand(newCardCreateBatchCmd())

	return cmd
}
func newCardRefreshCmd() *cobra.Command {
	var proposalID string

	cmd := &cobra.Command{
		Use:   "refresh [card-id]",
		Short: "Rebuild auto-generated card body content",
		Long: `Rebuild CLI-managed sections in the card body from frontmatter links.

  Structure cards: ## Entries - indexed card list
  All other cards:  ## Links - outgoing and incoming, grouped by relation

Use this to repair historical cards after upgrading flowforge or when
auto-generated sections are out of sync with frontmatter links.

Examples:
  flowforge card refresh DES-xxx
  flowforge card refresh --proposal CR26061601`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if proposalID == "" && len(args) == 0 {
				return fmt.Errorf("specify a card-id or --proposal")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			var cardIDs []string
			if proposalID != "" {
				cards, listErr := store.ListCards(store.ProposalCardsDir(proposalID))
				if listErr != nil {
					return fmt.Errorf("listing proposal cards: %w", listErr)
				}
				proposalCards, _ := store.ListCards(store.ProposalCardDir())
				for _, card := range proposalCards {
					if card.Source == proposalID || strings.HasPrefix(card.ID, "PROP-") {
						cards = append(cards, card)
					}
				}
				cards = append(cards,
					mustReadCard(store, "STR-"+proposalID+"-REQ"),
				)
				for _, card := range cards {
					if card != nil {
						cardIDs = append(cardIDs, card.ID)
					}
				}
			} else {
				cardIDs = []string{args[0]}
			}

			out := cmd.OutOrStdout()
			var refreshed, skipped int
			for _, cardID := range cardIDs {
				var changed bool
				if err := store.UpdateCardWithLock(cardID, func(card *core.Card) error {
					body, c, err := refreshCardGeneratedNavigation(store, card)
					if err != nil {
						return err
					}
					if c {
						card.Body = body
						changed = true
					}
					return nil
				}); err != nil {
					fmt.Fprintf(out, "✗ %s: %v\n", cardID, err)
					skipped++
					continue
				}
				if changed {
					fmt.Fprintf(out, "✓ %s\n", cardID)
					refreshed++
				} else {
					skipped++
				}
			}

			fmt.Fprintf(out, "\nRefreshed: %d, Skipped: %d\n", refreshed, skipped)
			return nil
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Refresh all cards in a proposal")
	return cmd
}

func mustReadCard(store *core.CardStore, cardID string) *core.Card {
	card, err := store.ReadCard(cardID)
	if err != nil {
		return nil
	}
	return card
}

func newCardCreateCmd() *cobra.Command {
	var (
		cardType   string
		title      string
		body       string
		status     string
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
  flowforge card create --type structure --title "CLI Architecture" --status active
  flowforge card create --type design --title "Init command" --status draft --body "## Goal\n\nDesign the init command."
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

			body, err := readBody(body)
			if err != nil {
				return err
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			card := core.NewCard(ct, title)
			card.Body = body
			card.Tags = tags
			if status != "" {
				card.Status = core.CardStatus(status)
				if !card.Status.Valid() {
					return fmt.Errorf("invalid status: %s", status)
				}
			}

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

			parsedLinks, err := parseLinkArgs(links)
			if err != nil {
				return err
			}
			addProposalOwnershipLink(card, resolvedProposalID)
			if len(card.Links) == 0 && len(parsedLinks) == 0 {
				return fmt.Errorf("card requires at least one outbound link; pass --links or create it under a proposal")
			}
			if err := ensureLinkTargetsExist(store, parsedLinks); err != nil {
				return err
			}
			for _, link := range parsedLinks {
				card.AddLink(link.target, link.relation)
			}

			upsertLinksSection(store, card)

			_, err = store.CreateCard(card, resolvedProposalID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			printResult(cmd, out, CommandResult{
				ID:    card.ID,
				Type:  string(card.Type),
				Title: card.Title,
			})

			return nil
		},
	}

	cmd.Flags().StringVar(&cardType, "type", "", "Card type (requirement/decision/design/task/log/convention/finding/module/structure)")
	cmd.Flags().StringVar(&title, "title", "", "Card title")
	cmd.Flags().StringVar(&body, "body", "", "Card body content; use '-' to read from stdin")
	cmd.Flags().StringVar(&status, "status", string(core.CardStatusDraft), "Card status (draft/active/accepted/deprecated/superseded)")
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID to associate with")
	cmd.Flags().StringSliceVar(&links, "links", nil, "Links to other cards (format: CARD_ID or CARD_ID:relation)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags for the card")

	return cmd
}

func newCardReadCmd() *cobra.Command {
	var (
		outputJSON bool
		summary    bool
		section    string
	)

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

			if summary && section != "" {
				return fmt.Errorf("--summary and --section cannot be used together")
			}

			out := cmd.OutOrStdout()
			if outputJSON {
				data, _ := json.MarshalIndent(card, "", "  ")
				fmt.Fprintln(out, string(data))
			} else if summary {
				printCardSummary(out, card)
			} else if section != "" {
				body, ok := extractCardSection(card.Body, section)
				if !ok {
					return fmt.Errorf("section %q not found in card %s", section, card.ID)
				}
				printCardFrontmatterSummary(out, card)
				if body != "" {
					fmt.Fprintln(out)
					fmt.Fprintln(out, body)
				}
			} else {
				fmt.Fprintf(out, "ID: %s\n", card.ID)
				fmt.Fprintf(out, "Type: %s\n", card.Type)
				fmt.Fprintf(out, "Title: %s\n", card.Title)
				fmt.Fprintf(out, "Status: %s\n", card.Status)
				fmt.Fprintf(out, "Importance: %s\n", card.Importance)
				if len(card.Tags) > 0 {
					fmt.Fprintf(out, "Tags: %s\n", strings.Join(card.Tags, ", "))
				}
				if len(card.Links) > 0 {
					fmt.Fprintln(out, "Links:")
					for _, link := range card.Links {
						fmt.Fprintf(out, "  - %s (%s)\n", link.Target, link.Relation)
					}
				}
				fmt.Fprintf(out, "Created: %s\n", card.Created.Format("2006-01-02 15:04:05"))
				fmt.Fprintf(out, "Updated: %s\n", card.Updated.Format("2006-01-02 15:04:05"))
				if card.Body != "" {
					fmt.Fprintln(out, "\n--- Body ---")
					fmt.Fprintln(out, card.Body)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&summary, "summary", false, "Output a card summary")
	cmd.Flags().StringVar(&section, "section", "", "Output a single markdown section by heading")

	return cmd
}

func printCardFrontmatterSummary(out io.Writer, card *core.Card) {
	fmt.Fprintf(out, "ID: %s\n", card.ID)
	fmt.Fprintf(out, "Type: %s\n", card.Type)
	fmt.Fprintf(out, "Title: %s\n", card.Title)
	fmt.Fprintf(out, "Status: %s\n", card.Status)
	fmt.Fprintf(out, "Importance: %s\n", card.Importance)
	if len(card.Tags) > 0 {
		fmt.Fprintf(out, "Tags: %s\n", strings.Join(card.Tags, ", "))
	}
	if card.Source != "" {
		fmt.Fprintf(out, "Source: %s\n", card.Source)
	}
	if card.Domain != "" {
		fmt.Fprintf(out, "Domain: %s\n", card.Domain)
	}
}

func printCardSummary(out io.Writer, card *core.Card) {
	printCardFrontmatterSummary(out, card)
	if summary, heading := firstMeaningfulSectionSummary(card.Body); summary != "" {
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Summary: %s\n", summary)
		if heading != "" {
			fmt.Fprintf(out, "First section: %s\n", heading)
		}
	}
	if len(card.Links) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Links:")
		for _, link := range card.Links {
			fmt.Fprintf(out, "  - %s (%s)\n", link.Target, link.Relation)
		}
	}
	headings := collectMarkdownHeadings(card.Body)
	if len(headings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Sections:")
		for _, heading := range headings {
			fmt.Fprintf(out, "  - %s\n", heading)
		}
	}
}

func firstMeaningfulSectionSummary(body string) (string, string) {
	preamble := summarizeSectionBody(bodyBeforeFirstHeading(body))
	if preamble != "" {
		return preamble, ""
	}
	sections := parseMarkdownSections(body)
	for _, section := range sections {
		summary := summarizeSectionBody(section.Body)
		if summary != "" {
			return summary, section.Heading
		}
	}
	return "", ""
}

func bodyBeforeFirstHeading(body string) string {
	lines := strings.Split(body, "\n")
	var collected []string
	for _, raw := range lines {
		if _, _, ok := parseMarkdownHeading(strings.TrimSpace(raw)); ok {
			break
		}
		collected = append(collected, raw)
	}
	return strings.TrimSpace(strings.Join(collected, "\n"))
}

type markdownSection struct {
	Level   int
	Heading string
	Body    string
}

func parseMarkdownSections(body string) []markdownSection {
	var sections []markdownSection
	lines := strings.Split(body, "\n")
	var current *markdownSection
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if level, heading, ok := parseMarkdownHeading(line); ok {
			sections = append(sections, markdownSection{Level: level, Heading: heading})
			current = &sections[len(sections)-1]
			continue
		}
		if current != nil {
			if current.Body == "" {
				current.Body = line
			} else {
				current.Body += "\n" + raw
			}
		}
	}
	return sections
}

func extractCardSection(body, section string) (string, bool) {
	target := strings.ToLower(strings.TrimSpace(section))
	if target == "" {
		return "", false
	}

	lines := strings.Split(body, "\n")
	var capture []string
	var active bool
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if _, heading, ok := parseMarkdownHeading(trimmed); ok {
			heading = strings.ToLower(heading)
			if active {
				break
			}
			if heading == target {
				active = true
				continue
			}
		}
		if active {
			capture = append(capture, line)
		}
	}
	if !active {
		return "", false
	}
	return strings.TrimSpace(strings.Join(capture, "\n")), true
}

func collectMarkdownHeadings(body string) []string {
	var headings []string
	for _, section := range parseMarkdownSections(body) {
		headings = append(headings, section.Heading)
	}
	return headings
}

func summarizeSectionBody(body string) string {
	for _, line := range strings.Split(strings.TrimSpace(body), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			continue
		case strings.HasPrefix(trimmed, "#"):
			continue
		case strings.HasPrefix(trimmed, "- "):
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
		default:
			return trimmed
		}
	}
	return ""
}

func parseMarkdownHeading(line string) (int, string, bool) {
	if !strings.HasPrefix(line, "#") {
		return 0, "", false
	}

	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level == 0 || level > len(line) || (len(line) > level && line[level] != ' ') {
		return 0, "", false
	}
	heading := strings.TrimSpace(line[level:])
	if heading == "" {
		return 0, "", false
	}
	return level, heading, true
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
		sectionName string
		addLinks    []string
		removeLinks []string
	)

	cmd := &cobra.Command{
		Use:   "update <card-id>",
		Short: "Update a card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]

			body, err := readBody(body)
			if err != nil {
				return err
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			if err := store.UpdateCardWithLock(cardID, func(card *core.Card) error {
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
					if sectionName != "" {
						card.Body = upsertMarkdownSection(card.Body, sectionName, body)
					} else {
						card.Body = body
					}
				}

				parsedAddLinks, err := parseLinkArgs(addLinks)
				if err != nil {
					return err
				}
				if err := ensureLinkTargetsExist(store, parsedAddLinks); err != nil {
					return err
				}
				for _, link := range parsedAddLinks {
					card.AddLink(link.target, link.relation)
				}

				parsedRemoveLinks, err := parseLinkArgs(removeLinks)
				if err != nil {
					return err
				}
				for _, link := range parsedRemoveLinks {
					card.RemoveLink(link.target, link.relation)
				}

				upsertLinksSection(store, card)

				return nil
			}); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			printResult(cmd, out, CommandResult{
				ID:      cardID,
				Updated: true,
			})
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().StringVar(&status, "status", "", "New status (draft/active/accepted/deprecated/superseded)")
	cmd.Flags().StringVar(&importance, "importance", "", "New importance (must/should/may)")
	cmd.Flags().StringVar(&body, "body", "", "New body content; use '-' to read from stdin")
	cmd.Flags().StringVar(&sectionName, "section", "", "Replace only this markdown section (requires --body)")
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

				if err := store.DeleteCard(cardID); err != nil {
					return err
				}
			
			} else {
				if err := store.ForceDeleteCard(cardID); err != nil {
					return err
				}
			}

			out := cmd.OutOrStdout()
			printResult(cmd, out, CommandResult{
				ID:      cardID,
				Updated: true,
			})
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force delete even if not draft")

	return cmd
}

func newCardRelatedCmd() *cobra.Command {
	var (
		relation  string
		depth     int
		direction string
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
			if direction != "forward" && direction != "backlinks" {
				return fmt.Errorf("invalid direction %q (expected forward/backlinks)", direction)
			}
			var related []*core.Card
			if direction == "backlinks" {
				dependents, err := store.GetDependents(cardID)
				if err != nil {
					return err
				}
				for _, dependent := range dependents {
					if relation == "" || hasLinkRelation(dependent, cardID, relation) {
						related = append(related, dependent)
					}
				}
			} else {
				related, err = store.GetRelated(cardID, relation, depth)
				if err != nil {
					return err
				}
			}

			out := cmd.OutOrStdout()
			if len(related) == 0 {
				fmt.Fprintln(out, "No related cards found.")
				return nil
			}

			fmt.Fprintf(out, "Related cards for %s:\n\n", cardID)
			for _, card := range related {
				fmt.Fprintf(out, "  %s [%s] %s\n", card.ID, card.Type, card.Title)
				fmt.Fprintf(out, "    Status: %s\n", card.Status)
				fmt.Fprintln(out)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&relation, "relation", "", "Filter by relation type")
	cmd.Flags().IntVar(&depth, "depth", 1, "Traversal depth")
	cmd.Flags().StringVar(&direction, "direction", "forward", "Traversal direction (forward/backlinks)")

	return cmd
}

func hasLinkRelation(card *core.Card, target, relation string) bool {
	for _, link := range card.Links {
		if link.Target == target && link.Relation == relation {
			return true
		}
	}
	return false
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
				fmt.Fprintln(cmd.OutOrStdout(), "No cards depend on this card.")
				return nil
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Cards depending on %s:\n\n", cardID)
			for _, card := range dependents {
				fmt.Fprintf(out, "  %s [%s] %s\n", card.ID, card.Type, card.Title)
				fmt.Fprintf(out, "    Status: %s\n", card.Status)
				fmt.Fprintln(out)
			}

			return nil
		},
	}

	return cmd
}

func newCardLinkCmd() *cobra.Command {
	var relation string

	cmd := &cobra.Command{
		Use:   "link <from-id> <to-id>",
		Short: "Add a link between cards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromID := args[0]
			toID := args[1]

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			if err := store.UpdateCardWithLock(fromID, func(card *core.Card) error {
				if !core.IsValidRelation(relation) {
					return fmt.Errorf("invalid relation: %s", relation)
				}
				if _, err := store.ReadCard(toID); err != nil {
					return fmt.Errorf("reading target card %s: %w", toID, err)
				}
				card.AddLink(toID, relation)
				upsertLinksSection(store, card)
				return nil
			}); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Linked %s -> %s (%s)\n", fromID, toID, relation)
			return nil
		},
	}

	cmd.Flags().StringVar(&relation, "relation", "related", "Link relation type")
	return cmd
}

func ensureLinkTargetsExist(store *core.CardStore, links []parsedLinkArg) error {
	for _, link := range links {
		if _, err := store.ReadCard(link.target); err != nil {
			return fmt.Errorf("reading target card %s: %w", link.target, err)
		}
	}
	return nil
}

func refreshCardGeneratedNavigation(store *core.CardStore, card *core.Card) (string, bool, error) {
	body := card.Body

	if card.Type == core.CardTypeStructure {
		entries, err := renderStructureEntries(store, card)
		if err != nil {
			return "", false, err
		}
		body = upsertMarkdownSection(card.Body, "Entries", entries)
		return body, body != card.Body, nil
	}

	linksContent := renderUnifiedLinksSection(store, card)
	body = upsertMarkdownSection(card.Body, "Links", linksContent)

	return body, body != card.Body, nil
}

func renderUnifiedLinksSection(store *core.CardStore, card *core.Card) string {
	var parts []string

	outgoing := renderOutgoingLinks(store, card)
	if outgoing != "" {
		parts = append(parts, "### Outgoing\n\n"+outgoing)
	}

	incoming := renderIncomingLinks(store, card)
	if incoming != "" {
		parts = append(parts, "### Incoming\n\n"+incoming)
	}

	if len(parts) == 0 {
		return "- None"
	}
	return strings.Join(parts, "\n\n")
}

func renderOutgoingLinks(store *core.CardStore, card *core.Card) string {
	if len(card.Links) == 0 {
		return ""
	}

	grouped := map[string][]string{}
	for _, link := range card.Links {
		targetCard, err := store.ReadCard(link.Target)
		line := ""
		if err != nil {
			line = fmt.Sprintf("- `%s` [%s]", link.Target, link.Relation)
		} else {
			targetPath, _ := markdownLinkTarget(card, targetCard)
			line = fmt.Sprintf("- [%s](%s) [%s] - %s", targetCard.ID, targetPath, targetCard.Type, targetCard.Title)
		}
		grouped[link.Relation] = append(grouped[link.Relation], line)
	}

	return renderGroupedLinks(grouped)
}

func renderIncomingLinks(store *core.CardStore, card *core.Card) string {
	dependents, err := store.GetDependents(card.ID)
	if err != nil || len(dependents) == 0 {
		return ""
	}

	grouped := map[string][]string{}
	for _, dep := range dependents {
		if dep.Type == core.CardTypeProposal || dep.Type == core.CardTypeStructure {
			continue
		}
		relation := ""
		for _, link := range dep.Links {
			if link.Target == card.ID {
				relation = link.Relation
				break
			}
		}
		targetPath, err := markdownLinkTarget(card, dep)
		line := ""
		if err != nil {
			line = fmt.Sprintf("- `%s` [%s]", dep.ID, dep.Type)
		} else {
			line = fmt.Sprintf("- [%s](%s) [%s] - %s", dep.ID, targetPath, dep.Type, dep.Title)
		}
		grouped[relation] = append(grouped[relation], line)
	}

	return renderGroupedLinks(grouped)
}

func renderGroupedLinks(grouped map[string][]string) string {
	if len(grouped) == 0 {
		return ""
	}

	keys := make([]string, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, relation := range keys {
		lines := grouped[relation]
		if len(lines) == 1 {
			parts = append(parts, lines[0])
		} else {
			parts = append(parts, "#### "+relation)
			parts = append(parts, lines...)
		}
	}
	return strings.Join(parts, "\n")
}

func linksToAny(card *core.Card, targets []string, relations []string) bool {
	targetSet := map[string]bool{}
	for _, target := range targets {
		targetSet[target] = true
	}
	return linksToAnyMap(card, targetSet, relations)
}

func linksToAnyMap(card *core.Card, targets map[string]bool, relations []string) bool {
	relationSet := map[string]bool{}
	for _, relation := range relations {
		relationSet[relation] = true
	}
	for _, link := range card.Links {
		if targets[link.Target] && relationSet[link.Relation] {
			return true
		}
	}
	return false
}

func cardInList(cards []*core.Card, cardID string) bool {
	for _, card := range cards {
		if card.ID == cardID {
			return true
		}
	}
	return false
}

func addProposalOwnershipLink(card *core.Card, proposalID string) {
	if card == nil || proposalID == "" || card.Type == core.CardTypeProposal {
		return
	}
	rootID := "PROP-" + proposalID
	for _, link := range card.Links {
		if link.Target == rootID && link.Relation == "belongs_to" {
			return
		}
	}
	card.AddLink(rootID, "belongs_to")
}

func newCardUnlinkCmd() *cobra.Command {
	var relation string

	cmd := &cobra.Command{
		Use:   "unlink <from-id> <to-id>",
		Short: "Remove a link between cards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromID := args[0]
			toID := args[1]

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			if err := store.UpdateCardWithLock(fromID, func(card *core.Card) error {
				if !card.RemoveLink(toID, relation) {
					return fmt.Errorf("link not found: %s -> %s (%s)", fromID, toID, relation)
				}
				upsertLinksSection(store, card)
				return nil
			}); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Unlinked %s -> %s (%s)\n", fromID, toID, relation)
			return nil
		},
	}

	cmd.Flags().StringVar(&relation, "relation", "related", "Link relation type")
	return cmd
}

func newCardSearchCmd() *cobra.Command {
	var (
		scope      string
		proposalID string
		types      string
		status     string
		domain     string
		tag        string
		limit      int
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

			cards, err := searchCards(store, query, scope, proposalID, types, status, domain, tag, limit)
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

	cmd.Flags().StringVar(&scope, "scope", "all", "Search scope (library/workspace/proposal/all)")
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (required when scope=proposal)")
	cmd.Flags().StringVar(&types, "type", "", "Comma-separated card types to include")
	cmd.Flags().StringVar(&types, "types", "", "Comma-separated card types to include")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&domain, "domain", "", "Filter by domain")
	cmd.Flags().StringVar(&tag, "tag", "", "Comma-separated tags; match any tag")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")

	return cmd
}

type cardSearchResult struct {
	Card        *core.Card
	MatchReason string
}

func searchCards(store *core.CardStore, query, scope, proposalID, types, status, domain, tag string, limit int) ([]cardSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	if scope == "proposal" && proposalID == "" {
		return nil, fmt.Errorf("--proposal is required when scope=proposal")
	}

	typeFilter := map[core.CardType]bool{}
	if strings.TrimSpace(types) != "" {
		for _, raw := range strings.Split(types, ",") {
			trimmed := strings.TrimSpace(raw)
			ct := core.CardType(trimmed)
			if !ct.Valid() {
				return nil, fmt.Errorf("invalid card type: %s", trimmed)
			}
			typeFilter[ct] = true
		}
	}

	tagFilter := splitNonEmptyCSV(tag)

	dirs, err := searchScopes(store, scope, proposalID)
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
		if match, ok := matchCardQuery(card, query); ok {
			filterReason, ok := matchCardSearchFilters(card, typeFilter, status, domain, tagFilter)
			if !ok {
				continue
			}
			reasonParts := []string{match}
			if filterReason != "" {
				reasonParts = append(reasonParts, filterReason)
			}
			results = append(results, cardSearchResult{Card: card, MatchReason: strings.Join(reasonParts, " | ")})
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

func searchScopes(store *core.CardStore, scope, proposalID string) ([]string, error) {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "", "all":
		return []string{store.ActiveDir(), store.IntakeDir(), store.LibraryDir()}, nil
	case "workspace":
		return []string{store.ActiveDir(), store.IntakeDir()}, nil
	case "library":
		return []string{store.LibraryDir()}, nil
	case "proposal":
		return []string{store.ProposalCardsDir(proposalID), store.ProposalCardDir()}, nil
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

func matchCardSearchFilters(card *core.Card, typeFilter map[core.CardType]bool, status, domain string, tagFilter []string) (string, bool) {
	if len(typeFilter) > 0 && !typeFilter[card.Type] {
		return "", false
	}
	if strings.TrimSpace(status) != "" && !strings.EqualFold(string(card.Status), strings.TrimSpace(status)) {
		return "", false
	}
	if strings.TrimSpace(domain) != "" && !strings.EqualFold(card.Domain, strings.TrimSpace(domain)) {
		return "", false
	}

	reasons := make([]string, 0, 4)
	if strings.TrimSpace(status) != "" {
		reasons = append(reasons, "status="+strings.TrimSpace(status))
	}
	if strings.TrimSpace(domain) != "" {
		reasons = append(reasons, "domain="+strings.TrimSpace(domain))
	}
	if len(tagFilter) > 0 {
		matchedTag, ok := matchCardTagFilter(card.Tags, tagFilter)
		if !ok {
			return "", false
		}
		reasons = append(reasons, "tag="+matchedTag)
	}

	return strings.Join(reasons, " | "), true
}

func matchCardTagFilter(cardTags, tagFilter []string) (string, bool) {
	if len(tagFilter) == 0 {
		return "", true
	}
	for _, want := range tagFilter {
		for _, tag := range cardTags {
			if strings.EqualFold(strings.TrimSpace(tag), want) {
				return want, true
			}
		}
	}
	return "", false
}

func splitNonEmptyCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	return filtered
}

func getOutputFormat() string {
	format := os.Getenv("FLOWFORGE_OUTPUT")
	if format == "" {
		format = "text"
	}
	return format
}

func readBody(body string) (string, error) {
	if body == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading body from stdin: %w", err)
		}
		return string(data), nil
	}
	return unescapeBody(body), nil
}

func unescapeBody(s string) string {
	if !strings.ContainsAny(s, "\\") {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case '\\':
				b.WriteByte('\\')
			default:
				b.WriteByte(s[i])
				i--
			}
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func upsertLinksSection(store *core.CardStore, card *core.Card) {
	if card.Type == core.CardTypeStructure {
		return
	}
	content := renderUnifiedLinksSection(store, card)
	card.Body = upsertMarkdownSection(card.Body, "Links", content)
}
