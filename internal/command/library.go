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
	cmd.AddCommand(newLibraryFacetsCmd())
	cmd.AddCommand(newLibraryClassifyCmd())
	cmd.AddCommand(newLibraryImportCmd())
	cmd.AddCommand(newLibraryPromoteCmd())

	return cmd
}

func newLibraryImportCmd() *cobra.Command {
	var (
		cardType     string
		title        string
		body         string
		status       string
		importance   string
		domain       string
		source       string
		sourceCardID string
		links        []string
		tags         []string
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a structured candidate into library",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(title) == "" {
				return fmt.Errorf("--title is required")
			}
			ct := core.CardType(cardType)
			if err := validateLibraryImportType(ct); err != nil {
				return err
			}
			cardStatus := core.CardStatus(status)
			if !cardStatus.Valid() {
				return fmt.Errorf("invalid status: %s", status)
			}
			cardImportance := core.Importance(importance)
			if !cardImportance.Valid() {
				return fmt.Errorf("invalid importance: %s", importance)
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
			card.ID = core.GenerateCardID(ct, "")
			card.Status = cardStatus
			card.Importance = cardImportance
			card.Body = body
			card.Tags = tags
			card.Domain = domain
			card.Source = source

			if sourceCardID != "" {
				if _, err := store.ReadCard(sourceCardID); err != nil {
					return fmt.Errorf("reading source card %s: %w", sourceCardID, err)
				}
				card.AddLink(sourceCardID, "references")
				if card.Source == "" {
					card.Source = sourceCardID
				}
			}

			parsedLinks, err := parseLinkArgs(links)
			if err != nil {
				return err
			}
			if err := ensureLinkTargetsExist(store, parsedLinks); err != nil {
				return err
			}
			for _, link := range parsedLinks {
				card.AddLink(link.target, link.relation)
			}

			upsertLinksSection(store, card)

			if len(card.Links) == 0 {
				return fmt.Errorf("library import requires at least one outbound link; pass --source-card or --links")
			}

			_, err = store.CreateCard(card, "")
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

	cmd.Flags().StringVar(&cardType, "type", "", "Library card type")
	cmd.Flags().StringVar(&title, "title", "", "Card title")
	cmd.Flags().StringVar(&body, "body", "", "Card body content; use '-' to read from stdin")
	cmd.Flags().StringVar(&status, "status", string(core.CardStatusActive), "Card status")
	cmd.Flags().StringVar(&importance, "importance", string(core.ImportanceShould), "Card importance")
	cmd.Flags().StringVar(&domain, "domain", "", "Card domain")
	cmd.Flags().StringVar(&source, "source", "", "Source label")
	cmd.Flags().StringVar(&sourceCardID, "source-card", "", "Source card ID to reference")
	cmd.Flags().StringSliceVar(&links, "links", nil, "Links to cards (format: CARD_ID or CARD_ID:relation)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags for the card")

	return cmd
}

func newLibraryPromoteCmd() *cobra.Command {
	var (
		cardType   string
		title      string
		status     string
		importance string
		domain     string
		links      []string
		tags       []string
	)

	cmd := &cobra.Command{
		Use:   "promote <card-id>",
		Short: "Promote a proposal card copy into library",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceCardID := args[0]

			store, err := currentCardStore()
			if err != nil {
				return err
			}
			sourceCard, err := store.ReadCard(sourceCardID)
			if err != nil {
				return err
			}

			ct := sourceCard.Type
			if strings.TrimSpace(cardType) != "" {
				ct = core.CardType(cardType)
			}
			if err := validateLibraryImportType(ct); err != nil {
				return err
			}

			cardStatus := core.CardStatus(status)
			if !cardStatus.Valid() {
				return fmt.Errorf("invalid status: %s", status)
			}
			cardImportance := core.Importance(importance)
			if !cardImportance.Valid() {
				return fmt.Errorf("invalid importance: %s", importance)
			}

			cardTitle := title
			if strings.TrimSpace(cardTitle) == "" {
				cardTitle = sourceCard.Title
			}
			cardDomain := domain
			if cardDomain == "" {
				cardDomain = sourceCard.Domain
			}
			cardTags := tags
			if len(cardTags) == 0 {
				cardTags = append([]string{}, sourceCard.Tags...)
			}

			card := core.NewCard(ct, cardTitle)
			card.ID = core.GenerateCardID(ct, "")
			card.Status = cardStatus
			card.Importance = cardImportance
			card.Body = sourceCard.Body
			card.Tags = cardTags
			card.Domain = cardDomain
			card.Source = sourceCard.ID
			card.AddLink(sourceCard.ID, "references")

			parsedLinks, err := parseLinkArgs(links)
			if err != nil {
				return err
			}
			if err := ensureLinkTargetsExist(store, parsedLinks); err != nil {
				return err
			}
			for _, link := range parsedLinks {
				card.AddLink(link.target, link.relation)
			}

			upsertLinksSection(store, card)

			_, err = store.CreateCard(card, "")
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Promoted %s to library card %s\n", sourceCard.ID, card.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Type: %s\n", card.Type)
			return nil
		},
	}

	cmd.Flags().StringVar(&cardType, "type", "", "Override library card type")
	cmd.Flags().StringVar(&title, "title", "", "Override card title")
	cmd.Flags().StringVar(&status, "status", string(core.CardStatusActive), "Card status")
	cmd.Flags().StringVar(&importance, "importance", string(core.ImportanceShould), "Card importance")
	cmd.Flags().StringVar(&domain, "domain", "", "Override card domain")
	cmd.Flags().StringSliceVar(&links, "links", nil, "Additional links to cards (format: CARD_ID or CARD_ID:relation)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Override tags for the card")

	return cmd
}

func validateLibraryImportType(cardType core.CardType) error {
	if !cardType.Valid() {
		return fmt.Errorf("invalid library card type: %s", cardType)
	}
	switch cardType {
	case core.CardTypeRequirement,
		core.CardTypeDecision,
		core.CardTypeDesign,
		core.CardTypeConvention,
		core.CardTypeFinding,
		core.CardTypeModule,
		core.CardTypeStructure:
		return nil
	default:
		return fmt.Errorf("card type %s cannot be imported into library through this command", cardType)
	}
}

func newLibrarySuggestCmd() *cobra.Command {
	var (
		forCardID string
		types     string
		relation  string
		facets    []string
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

			facetFilter, err := parseFacetArgs(facets)
			if err != nil {
				return err
			}

			candidates, err := suggestLibraryCards(store, focus, typeFilter, relation, facetFilter, limit)
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
	cmd.Flags().StringSliceVar(&facets, "facet", nil, "Facet filter (format: key:value, repeatable or comma-separated)")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")

	return cmd
}

func newLibraryFacetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facets",
		Short: "List discovered library facets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := currentCardStore()
			if err != nil {
				return err
			}

			cards, err := store.ListCards(store.LibraryDir())
			if err != nil {
				return err
			}

			writeLibraryFacets(cmd.OutOrStdout(), buildLibraryFacetIndex(cards))
			return nil
		},
	}

	return cmd
}

func newLibraryClassifyCmd() *cobra.Command {
	var forCardID string

	cmd := &cobra.Command{
		Use:   "classify",
		Short: "Classify a focus card against discovered library facets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(forCardID) == "" {
				return fmt.Errorf("--for is required")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			focus, err := store.ReadCard(forCardID)
			if err != nil {
				return err
			}

			cards, err := store.ListCards(store.LibraryDir())
			if err != nil {
				return err
			}

			index := buildLibraryFacetIndex(cards)
			matches := classifyLibraryFacets(focus, index)
			writeLibraryClassification(cmd.OutOrStdout(), focus, index, matches)
			return nil
		},
	}

	cmd.Flags().StringVar(&forCardID, "for", "", "Focus card ID")
	return cmd
}

type librarySuggestion struct {
	Card              *core.Card
	Score             int
	MatchReason       string
	SuggestedRelation string
}

type libraryFacet struct {
	Key   string
	Value string
}

type libraryFacetValue struct {
	Facet libraryFacet
	Count int
}

type libraryFacetCombination struct {
	Left  libraryFacet
	Right libraryFacet
	Count int
}

type libraryFacetIndex struct {
	Values       map[string]map[string]int
	Combinations map[string]int
	CardCount    int
}

type libraryFacetMatch struct {
	Facet    libraryFacet
	Source   string
	Evidence string
	CardHits int
}

func defaultLibrarySuggestionTypes() []core.CardType {
	return []core.CardType{
		core.CardTypeConvention,
		core.CardTypeDecision,
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

func suggestLibraryCards(store *core.CardStore, focus *core.Card, typeFilter map[core.CardType]bool, relation string, facetFilter []libraryFacet, limit int) ([]librarySuggestion, error) {
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
		facetScore, facetReason, ok := matchFacetFilter(card, facetFilter)
		if !ok {
			continue
		}

		score, reason, ok := scoreLibraryCandidate(card, terms)
		if !ok && facetScore == 0 {
			continue
		}
		if facetScore > 0 {
			score += facetScore
			if reason == "" {
				reason = facetReason
			} else {
				reason = reason + "; " + facetReason
			}
		}

		suggestedRelation := suggestedLibraryRelation(card.Type, relation)

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

func parseFacetArgs(values []string) ([]libraryFacet, error) {
	var facets []libraryFacet
	for _, value := range values {
		for _, raw := range strings.Split(value, ",") {
			raw = strings.TrimSpace(raw)
			if raw == "" {
				continue
			}
			facet, ok := parseFacetTag(raw)
			if !ok {
				return nil, fmt.Errorf("invalid facet %q (expected key:value)", raw)
			}
			facets = append(facets, facet)
		}
	}
	return facets, nil
}

func parseFacetTag(tag string) (libraryFacet, bool) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return libraryFacet{}, false
	}

	if strings.HasPrefix(tag, "facet:") {
		tag = strings.TrimPrefix(tag, "facet:")
	}

	parts := strings.SplitN(tag, ":", 2)
	if len(parts) != 2 {
		return libraryFacet{}, false
	}

	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.ToLower(strings.TrimSpace(parts[1]))
	if key == "" || value == "" {
		return libraryFacet{}, false
	}
	return libraryFacet{Key: key, Value: value}, true
}

func facetString(facet libraryFacet) string {
	return facet.Key + ":" + facet.Value
}

func cardFacets(card *core.Card) []libraryFacet {
	if card == nil {
		return nil
	}
	var facets []libraryFacet
	seen := map[string]bool{}
	for _, tag := range card.Tags {
		facet, ok := parseFacetTag(tag)
		if !ok {
			continue
		}
		key := facetString(facet)
		if seen[key] {
			continue
		}
		seen[key] = true
		facets = append(facets, facet)
	}
	sort.SliceStable(facets, func(i, j int) bool {
		if facets[i].Key == facets[j].Key {
			return facets[i].Value < facets[j].Value
		}
		return facets[i].Key < facets[j].Key
	})
	return facets
}

func buildLibraryFacetIndex(cards []*core.Card) libraryFacetIndex {
	index := libraryFacetIndex{
		Values:       map[string]map[string]int{},
		Combinations: map[string]int{},
	}

	for _, card := range cards {
		if card.Status == core.CardStatusDeprecated || card.Status == core.CardStatusSuperseded {
			continue
		}
		facets := cardFacets(card)
		if len(facets) == 0 {
			continue
		}
		index.CardCount++
		for _, facet := range facets {
			if index.Values[facet.Key] == nil {
				index.Values[facet.Key] = map[string]int{}
			}
			index.Values[facet.Key][facet.Value]++
		}
		for i := 0; i < len(facets); i++ {
			for j := i + 1; j < len(facets); j++ {
				key := facetString(facets[i]) + " + " + facetString(facets[j])
				index.Combinations[key]++
			}
		}
	}

	return index
}

func sortedFacetValues(index libraryFacetIndex) []libraryFacetValue {
	var values []libraryFacetValue
	for key, byValue := range index.Values {
		for value, count := range byValue {
			values = append(values, libraryFacetValue{
				Facet: libraryFacet{Key: key, Value: value},
				Count: count,
			})
		}
	}
	sort.SliceStable(values, func(i, j int) bool {
		if values[i].Facet.Key == values[j].Facet.Key {
			if values[i].Count == values[j].Count {
				return values[i].Facet.Value < values[j].Facet.Value
			}
			return values[i].Count > values[j].Count
		}
		return values[i].Facet.Key < values[j].Facet.Key
	})
	return values
}

func sortedFacetCombinations(index libraryFacetIndex, limit int) []libraryFacetCombination {
	var combinations []libraryFacetCombination
	for raw, count := range index.Combinations {
		parts := strings.Split(raw, " + ")
		if len(parts) != 2 {
			continue
		}
		left, leftOK := parseFacetTag(parts[0])
		right, rightOK := parseFacetTag(parts[1])
		if !leftOK || !rightOK {
			continue
		}
		combinations = append(combinations, libraryFacetCombination{Left: left, Right: right, Count: count})
	}
	sort.SliceStable(combinations, func(i, j int) bool {
		if combinations[i].Count == combinations[j].Count {
			leftI := facetString(combinations[i].Left) + " + " + facetString(combinations[i].Right)
			leftJ := facetString(combinations[j].Left) + " + " + facetString(combinations[j].Right)
			return leftI < leftJ
		}
		return combinations[i].Count > combinations[j].Count
	})
	if limit > 0 && len(combinations) > limit {
		combinations = combinations[:limit]
	}
	return combinations
}

func matchFacetFilter(card *core.Card, filter []libraryFacet) (int, string, bool) {
	if len(filter) == 0 {
		return 0, "", true
	}

	cardFacetSet := map[string]bool{}
	for _, facet := range cardFacets(card) {
		cardFacetSet[facetString(facet)] = true
	}

	var matched []string
	for _, facet := range filter {
		key := facetString(facet)
		if !cardFacetSet[key] {
			return 0, "", false
		}
		matched = append(matched, key)
	}

	return 75 * len(matched), "facets:" + strings.Join(matched, ","), true
}

func classifyLibraryFacets(focus *core.Card, index libraryFacetIndex) []libraryFacetMatch {
	text := strings.ToLower(strings.Join([]string{focus.Title, focus.Domain, strings.Join(focus.Tags, " "), focus.Body}, " "))
	focusFacetSet := map[string]bool{}
	for _, facet := range cardFacets(focus) {
		focusFacetSet[facetString(facet)] = true
	}

	var matches []libraryFacetMatch
	for _, value := range sortedFacetValues(index) {
		facetKey := facetString(value.Facet)
		if focusFacetSet[facetKey] {
			matches = append(matches, libraryFacetMatch{
				Facet:    value.Facet,
				Source:   "tag",
				Evidence: facetKey,
				CardHits: value.Count,
			})
			continue
		}

		if containsFacetValue(text, value.Facet.Value) {
			matches = append(matches, libraryFacetMatch{
				Facet:    value.Facet,
				Source:   "text",
				Evidence: value.Facet.Value,
				CardHits: value.Count,
			})
		}
	}

	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].Source != matches[j].Source {
			return matches[i].Source == "tag"
		}
		if matches[i].CardHits != matches[j].CardHits {
			return matches[i].CardHits > matches[j].CardHits
		}
		return facetString(matches[i].Facet) < facetString(matches[j].Facet)
	})
	return matches
}

func containsFacetValue(text string, value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return false
	}
	if strings.Contains(text, value) {
		return true
	}
	return strings.Contains(text, strings.ReplaceAll(value, "-", " "))
}

func suggestedLibraryRelation(cardType core.CardType, requested string) string {
	requested = strings.ToLower(strings.TrimSpace(requested))
	if requested != "" {
		if requested == "constrains" && constrainsPreferredType(cardType) {
			return "constrains"
		}
		return requested
	}

	switch cardType {
	case core.CardTypeConvention, core.CardTypeModule:
		return "constrains"
	case core.CardTypeDecision, core.CardTypeDesign, core.CardTypeFinding:
		return "references"
	default:
		return "related"
	}
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
	fmt.Fprintf(out, "| ID | Type | Title | Status | Importance | Domain | Score | SuggestedRelation |\n")
	fmt.Fprintf(out, "|----|------|-------|--------|------------|--------|-------|-------------------|\n")
	for _, candidate := range candidates {
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s | %s | %d | %s |\n",
			candidate.Card.ID,
			candidate.Card.Type,
			escapeTableCell(candidate.Card.Title),
			candidate.Card.Status,
			candidate.Card.Importance,
			escapeTableCell(candidate.Card.Domain),
			candidate.Score,
			candidate.SuggestedRelation,
		)
	}

	fmt.Fprintln(out, "\n## Match Reasons")
	fmt.Fprintln(out)
	if len(candidates) == 0 {
		fmt.Fprintln(out, "- No candidates matched the focus card.")
	} else {
		fmt.Fprintln(out, "| ID | matchedBy | Reason |")
		fmt.Fprintln(out, "|----|-----------|--------|")
		for _, candidate := range candidates {
			fmt.Fprintf(out, "| %s | keyword | %s |\n",
				candidate.Card.ID,
				escapeTableCell(candidate.MatchReason),
			)
		}
	}

	fmt.Fprintln(out, "\n## Recommended Reads")
	fmt.Fprintln(out)
	if len(candidates) == 0 {
		fmt.Fprintln(out, "- None")
	} else {
		fmt.Fprintln(out, "| ID | Section | Reason |")
		fmt.Fprintln(out, "|----|---------|--------|")
		for _, candidate := range candidates {
			fmt.Fprintf(out, "| %s | summary | Validate before linking with `flowforge card read %s --summary` |\n",
				candidate.Card.ID,
				candidate.Card.ID,
			)
		}
	}

	fmt.Fprintln(out, "\n## Not Included")
	fmt.Fprintln(out, "- Deprecated and superseded cards are omitted by default.")
}

func writeLibraryFacets(out io.Writer, index libraryFacetIndex) {
	fmt.Fprintln(out, "## Library Facets")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "- CardsWithFacets: %d\n", index.CardCount)
	fmt.Fprintln(out)

	values := sortedFacetValues(index)
	if len(values) == 0 {
		fmt.Fprintln(out, "- None")
	} else {
		fmt.Fprintln(out, "| Facet | Value | Cards |")
		fmt.Fprintln(out, "|-------|-------|-------|")
		for _, value := range values {
			fmt.Fprintf(out, "| %s | %s | %d |\n", value.Facet.Key, value.Facet.Value, value.Count)
		}
	}

	fmt.Fprintln(out, "\n## Common Combinations")
	combinations := sortedFacetCombinations(index, 10)
	if len(combinations) == 0 {
		fmt.Fprintln(out, "- None")
		return
	}
	fmt.Fprintln(out, "| Facets | Cards |")
	fmt.Fprintln(out, "|--------|-------|")
	for _, combination := range combinations {
		fmt.Fprintf(out, "| %s + %s | %d |\n", facetString(combination.Left), facetString(combination.Right), combination.Count)
	}
}

func writeLibraryClassification(out io.Writer, focus *core.Card, index libraryFacetIndex, matches []libraryFacetMatch) {
	fmt.Fprintln(out, "## Library Classification")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "For: %s\n", focus.ID)
	fmt.Fprintf(out, "KnownFacetCards: %d\n", index.CardCount)
	fmt.Fprintln(out)

	fmt.Fprintln(out, "## Extracted Facets")
	if len(matches) == 0 {
		fmt.Fprintln(out, "- None")
	} else {
		fmt.Fprintln(out, "| Facet | Source | Evidence | LibraryCards |")
		fmt.Fprintln(out, "|-------|--------|----------|--------------|")
		for _, match := range matches {
			fmt.Fprintf(out, "| %s | %s | %s | %d |\n",
				facetString(match.Facet),
				match.Source,
				escapeTableCell(match.Evidence),
				match.CardHits,
			)
		}
	}

	fmt.Fprintln(out, "\n## Suggested Commands")
	if len(matches) == 0 {
		fmt.Fprintf(out, "- flowforge library facets\n")
		fmt.Fprintf(out, "- flowforge library suggest --for %s\n", focus.ID)
		return
	}

	var args []string
	limit := len(matches)
	if limit > 4 {
		limit = 4
	}
	for i := 0; i < limit; i++ {
		args = append(args, "--facet "+facetString(matches[i].Facet))
	}
	fmt.Fprintf(out, "- flowforge library suggest --for %s %s\n", focus.ID, strings.Join(args, " "))
}
