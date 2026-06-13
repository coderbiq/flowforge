package command

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newLibraryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "library",
		Short: "Work with library cards",
		Long:  "Inspect and recommend library cards for design and analysis work.",
	}

	cmd.AddCommand(newLibrarySuggestCmd())

	return cmd
}

func newLibrarySuggestCmd() *cobra.Command {
	var (
		forCardID string
		types     string
		relation  string
		limit     int
	)

	cmd := &cobra.Command{
		Use:   "suggest",
		Short: "Suggest related library cards",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(forCardID) == "" {
				return fmt.Errorf("--for is required")
			}
			if limit <= 0 {
				limit = 10
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			focus, err := store.ReadCard(forCardID)
			if err != nil {
				return err
			}

			typeFilter, err := parseCardTypeFilter(types, defaultLibrarySuggestionTypes())
			if err != nil {
				return err
			}

			candidates, err := suggestLibraryCards(store, focus, typeFilter, relation, limit)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			writeLibrarySuggestions(out, focus, candidates)
			return nil
		},
	}

	cmd.Flags().StringVar(&forCardID, "for", "", "Focus card ID")
	cmd.Flags().StringVar(&types, "types", "", "Comma-separated card types to include")
	cmd.Flags().StringVar(&types, "type", "", "Comma-separated card types to include")
	cmd.Flags().StringVar(&relation, "relation", "", "Suggested relation to prefer")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")

	return cmd
}

type librarySuggestion struct {
	Card              *core.Card
	Score             int
	MatchReason       string
	SuggestedRelation string
}

func defaultLibrarySuggestionTypes() []core.CardType {
	return []core.CardType{
		core.CardTypeConvention,
		core.CardTypeModule,
		core.CardTypeDesign,
		core.CardTypeFinding,
	}
}

func parseCardTypeFilter(raw string, defaults []core.CardType) (map[core.CardType]bool, error) {
	filter := map[core.CardType]bool{}
	if strings.TrimSpace(raw) == "" {
		for _, ct := range defaults {
			filter[ct] = true
		}
		return filter, nil
	}

	for _, part := range strings.Split(raw, ",") {
		ct := core.CardType(strings.TrimSpace(part))
		if !ct.Valid() {
			return nil, fmt.Errorf("invalid card type: %s", part)
		}
		filter[ct] = true
	}

	return filter, nil
}

func suggestLibraryCards(store *core.CardStore, focus *core.Card, typeFilter map[core.CardType]bool, relation string, limit int) ([]librarySuggestion, error) {
	cards, err := store.ListCards(store.LibraryDir())
	if err != nil {
		return nil, err
	}

	terms := buildLibraryQueryTerms(focus)
	suggestions := make([]librarySuggestion, 0, len(cards))
	for _, card := range cards {
		if len(typeFilter) > 0 && !typeFilter[card.Type] {
			continue
		}
		if card.Status == core.CardStatusDeprecated || card.Status == core.CardStatusSuperseded {
			continue
		}

		score, reason, ok := scoreLibraryCandidate(card, terms)
		if !ok {
			continue
		}

		suggestedRelation := "related"
		if strings.EqualFold(strings.TrimSpace(relation), "constrains") && constrainsPreferredType(card.Type) {
			suggestedRelation = "constrains"
		}

		suggestions = append(suggestions, librarySuggestion{
			Card:              card,
			Score:             score,
			MatchReason:       reason,
			SuggestedRelation: suggestedRelation,
		})
	}

	sort.SliceStable(suggestions, func(i, j int) bool {
		if suggestions[i].Score != suggestions[j].Score {
			return suggestions[i].Score > suggestions[j].Score
		}
		if suggestions[i].Card.Importance != suggestions[j].Card.Importance {
			return suggestions[i].Card.Importance == core.ImportanceMust
		}
		if suggestions[i].SuggestedRelation != suggestions[j].SuggestedRelation {
			return suggestions[i].SuggestedRelation == "constrains"
		}
		if suggestions[i].Card.ID != suggestions[j].Card.ID {
			return suggestions[i].Card.ID < suggestions[j].Card.ID
		}
		return suggestions[i].Card.Title < suggestions[j].Card.Title
	})

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions, nil
}

func constrainsPreferredType(cardType core.CardType) bool {
	switch cardType {
	case core.CardTypeConvention, core.CardTypeDecision, core.CardTypeModule:
		return true
	default:
		return false
	}
}

func buildLibraryQueryTerms(focus *core.Card) []string {
	terms := make([]string, 0, 16)
	terms = append(terms, tokenizeText(focus.Title)...)
	terms = append(terms, tokenizeText(focus.Domain)...)
	for _, tag := range focus.Tags {
		terms = append(terms, tokenizeText(tag)...)
	}
	terms = append(terms, significantWords(focus.Body)...)
	return uniqueStrings(terms)
}

func scoreLibraryCandidate(card *core.Card, terms []string) (int, string, bool) {
	score := 0
	reasons := make([]string, 0, 3)

	titleMatches := matchedTerms(card.Title, terms)
	if len(titleMatches) > 0 {
		score += 100
		reasons = append(reasons, "title:"+strings.Join(titleMatches, ","))
	}

	tagMatches := matchedTerms(strings.Join(card.Tags, " "), terms)
	domainMatches := matchedTerms(card.Domain, terms)
	if len(tagMatches) > 0 || len(domainMatches) > 0 {
		score += 50
		var parts []string
		if len(tagMatches) > 0 {
			parts = append(parts, "tags:"+strings.Join(tagMatches, ","))
		}
		if len(domainMatches) > 0 {
			parts = append(parts, "domain:"+strings.Join(domainMatches, ","))
		}
		reasons = append(reasons, strings.Join(parts, " "))
	}

	bodyMatches := matchedTerms(card.Body, terms)
	if len(bodyMatches) > 0 {
		score += 10
		reasons = append(reasons, "body:"+strings.Join(bodyMatches, ","))
	}

	if score == 0 {
		return 0, "", false
	}

	return score, strings.Join(reasons, "; "), true
}

func tokenizeText(text string) []string {
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return fields
}

func significantWords(text string) []string {
	stop := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true, "be": true,
		"by": true, "for": true, "from": true, "in": true, "into": true, "is": true, "it": true,
		"of": true, "on": true, "or": true, "the": true, "to": true, "with": true, "this": true,
		"that": true, "these": true, "those": true, "we": true, "you": true, "they": true,
		"should": true, "must": true, "may": true, "can": true, "will": true, "not": true,
	}

	terms := tokenizeText(text)
	filtered := make([]string, 0, len(terms))
	for _, term := range terms {
		if len(term) < 4 || stop[term] {
			continue
		}
		filtered = append(filtered, term)
	}
	return filtered
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func matchedTerms(text string, terms []string) []string {
	lower := strings.ToLower(text)
	matches := make([]string, 0, len(terms))
	seen := map[string]bool{}
	for _, term := range terms {
		if term == "" || seen[term] {
			continue
		}
		if strings.Contains(lower, term) {
			seen[term] = true
			matches = append(matches, term)
		}
	}
	return matches
}

func writeLibrarySuggestions(out io.Writer, focus *core.Card, candidates []librarySuggestion) {
	fmt.Fprintf(out, "## Library Suggestions\n\n")
	fmt.Fprintf(out, "For: %s\n\n", focus.ID)
	fmt.Fprintf(out, "| ID | Type | Title | Status | Importance | SuggestedRelation | MatchReason |\n")
	fmt.Fprintf(out, "|----|------|-------|--------|------------|-------------------|-------------|\n")
	for _, candidate := range candidates {
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s | %s | %s |\n",
			candidate.Card.ID,
			candidate.Card.Type,
			candidate.Card.Title,
			candidate.Card.Status,
			candidate.Card.Importance,
			candidate.SuggestedRelation,
			candidate.MatchReason,
		)
	}
}
