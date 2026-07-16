package state

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"flowforge/internal/core"
)

type CardSyncService struct {
	db *sql.DB
}

func NewCardSyncService(db *sql.DB) *CardSyncService {
	return &CardSyncService{db: db}
}

func (s *CardSyncService) SyncCard(card *core.Card) error {
	if s == nil || s.db == nil {
		return nil
	}
	if card == nil || card.ID == "" {
		return fmt.Errorf("card ID is required")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin sync transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	summary := computeCardSummary(card)

	if _, err := tx.Exec(`
INSERT INTO card_index(id, type, title, status, importance, source, domain, file_path, updated_at, created_at, body, summary)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
	type = excluded.type,
	title = excluded.title,
	status = excluded.status,
	importance = excluded.importance,
	source = excluded.source,
	domain = excluded.domain,
	file_path = excluded.file_path,
	updated_at = excluded.updated_at,
	created_at = excluded.created_at,
	body = excluded.body,
	summary = excluded.summary;`,
		card.ID, string(card.Type), card.Title, string(card.Status), string(card.Importance),
		card.Source, card.Domain, card.FilePath, card.Updated.UTC().Format(time.RFC3339Nano),
		card.Created.UTC().Format(time.RFC3339Nano), card.Body, summary,
	); err != nil {
		return fmt.Errorf("upserting card_index: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM card_link WHERE from_id = ?;`, card.ID); err != nil {
		return fmt.Errorf("clearing card_link: %w", err)
	}
	for _, link := range card.Links {
		if link.Target == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO card_link(from_id, to_id, relation) VALUES (?, ?, ?);`,
			card.ID, link.Target, link.Relation,
		); err != nil {
			return fmt.Errorf("inserting card_link: %w", err)
		}
	}

	if _, err := tx.Exec(`DELETE FROM card_tag WHERE card_id = ?;`, card.ID); err != nil {
		return fmt.Errorf("clearing card_tag: %w", err)
	}
	for _, tag := range card.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO card_tag(card_id, tag) VALUES (?, ?);`, card.ID, tag); err != nil {
			return fmt.Errorf("inserting card_tag: %w", err)
		}
	}

	if _, err := tx.Exec(`DELETE FROM card_term WHERE card_id = ?;`, card.ID); err != nil {
		return fmt.Errorf("clearing card_term: %w", err)
	}
	terms := buildCardTerms(card)
	for _, t := range terms {
		if _, err := tx.Exec(`INSERT INTO card_term(card_id, term, source) VALUES (?, ?, ?);`,
			card.ID, t.term, t.source,
		); err != nil {
			return fmt.Errorf("inserting card_term: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing sync: %w", err)
	}
	tx = nil
	return nil
}

func (s *CardSyncService) DeleteCard(cardID string) error {
	if s == nil || s.db == nil || cardID == "" {
		return nil
	}
	_, err := s.db.Exec(`DELETE FROM card_index WHERE id = ?;`, cardID)
	if err != nil {
		return fmt.Errorf("deleting from card_index: %w", err)
	}
	_, _ = s.db.Exec(`DELETE FROM card_link WHERE from_id = ? OR to_id = ?;`, cardID, cardID)
	_, _ = s.db.Exec(`DELETE FROM card_tag WHERE card_id = ?;`, cardID)
	_, _ = s.db.Exec(`DELETE FROM card_term WHERE card_id = ?;`, cardID)
	return nil
}

func (s *CardSyncService) RebuildAll(listCardsFromFiles func(dir string) ([]*core.Card, error), dirs []string) (int, int, error) {
	if s == nil || s.db == nil {
		return 0, 0, fmt.Errorf("sync service not initialized")
	}

	var cards []*core.Card
	for _, dir := range dirs {
		dirCards, err := listCardsFromFiles(dir)
		if err != nil {
			return 0, 0, fmt.Errorf("scanning %s: %w", dir, err)
		}
		cards = append(cards, dirCards...)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("begin rebuild transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	for _, table := range []string{"card_index", "card_link", "card_tag", "card_term"} {
		if _, err := tx.Exec("DELETE FROM " + table + ";"); err != nil {
			return 0, 0, fmt.Errorf("clearing %s: %w", table, err)
		}
	}

	cardCount, linkCount := 0, 0
	seen := map[string]bool{}
	for _, card := range cards {
		if card == nil || card.ID == "" || seen[card.ID] {
			continue
		}
		seen[card.ID] = true

		summary := computeCardSummary(card)
		if _, err := tx.Exec(`
INSERT INTO card_index(id, type, title, status, importance, source, domain, file_path, updated_at, created_at, body, summary)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
			card.ID, string(card.Type), card.Title, string(card.Status), string(card.Importance),
			card.Source, card.Domain, card.FilePath, card.Updated.UTC().Format(time.RFC3339Nano),
			card.Created.UTC().Format(time.RFC3339Nano), card.Body, summary,
		); err != nil {
			return 0, 0, fmt.Errorf("inserting card %s: %w", card.ID, err)
		}
		cardCount++

		for _, link := range card.Links {
			if link.Target == "" {
				continue
			}
			if _, err := tx.Exec(`INSERT INTO card_link(from_id, to_id, relation) VALUES (?, ?, ?);`,
				card.ID, link.Target, link.Relation,
			); err != nil {
				return 0, 0, fmt.Errorf("inserting link from %s to %s: %w", card.ID, link.Target, err)
			}
			linkCount++
		}

		for _, tag := range card.Tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			if _, err := tx.Exec(`INSERT INTO card_tag(card_id, tag) VALUES (?, ?);`, card.ID, tag); err != nil {
				return 0, 0, fmt.Errorf("inserting card_tag: %w", err)
			}
		}

		for _, t := range buildCardTerms(card) {
			if _, err := tx.Exec(`INSERT INTO card_term(card_id, term, source) VALUES (?, ?, ?);`,
				card.ID, t.term, t.source,
			); err != nil {
				return 0, 0, fmt.Errorf("inserting card_term: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("committing rebuild: %w", err)
	}
	tx = nil
	return cardCount, linkCount, nil
}

func (s *CardSyncService) ReadCard(id string) (*core.Card, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("sync service not initialized")
	}

	row := s.db.QueryRow(`
SELECT id, type, title, status, importance, source, domain, file_path, updated_at, created_at, body
FROM card_index WHERE id = ?;`, id)

	card, err := scanCard(row)
	if err != nil {
		return nil, err
	}

	card.Links = s.readLinks(id)
	card.Tags = s.readTags(id)
	return card, nil
}

func (s *CardSyncService) FindCardPath(id string) (string, error) {
	if s == nil || s.db == nil {
		return "", fmt.Errorf("sync service not initialized")
	}
	var path string
	err := s.db.QueryRow(`SELECT file_path FROM card_index WHERE id = ?;`, id).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("card not found: %s", id)
	}
	if err != nil {
		return "", fmt.Errorf("finding card path: %w", err)
	}
	return path, nil
}

func (s *CardSyncService) ListCards(dir string) ([]*core.Card, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("sync service not installed")
	}

	rows, err := s.db.Query(`
SELECT id, type, title, status, importance, source, domain, file_path, updated_at, created_at, body
FROM card_index WHERE file_path LIKE ? ORDER BY file_path;`, dir+"/%")
	if err != nil {
		return nil, fmt.Errorf("listing cards: %w", err)
	}
	defer rows.Close()

	cards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}
	if err := s.hydrateLinks(cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *CardSyncService) ListCardsByType(cardType core.CardType) ([]*core.Card, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("sync service not initialized")
	}

	rows, err := s.db.Query(`
SELECT id, type, title, status, importance, source, domain, file_path, updated_at, created_at, body
FROM card_index WHERE type = ? ORDER BY id;`, string(cardType))
	if err != nil {
		return nil, fmt.Errorf("listing cards by type: %w", err)
	}
	defer rows.Close()

	cards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}
	if err := s.hydrateLinks(cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *CardSyncService) GetDependents(cardID string) ([]*core.Card, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("sync service not initialized")
	}

	rows, err := s.db.Query(`
SELECT ci.id, ci.type, ci.title, ci.status, ci.importance, ci.source, ci.domain, ci.file_path, ci.updated_at, ci.created_at, ci.body
FROM card_index ci
INNER JOIN card_link cl ON cl.from_id = ci.id
WHERE cl.to_id = ?
ORDER BY ci.id;`, cardID)
	if err != nil {
		return nil, fmt.Errorf("querying dependents: %w", err)
	}
	defer rows.Close()

	cards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}
	if err := s.hydrateLinks(cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *CardSyncService) SearchCards(query string, typeFilter map[core.CardType]bool, statusFilter, domainFilter string, tagFilter []string, limit int) ([]core.CardSearchResult, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("sync service not initialized")
	}

	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return nil, fmt.Errorf("query is required")
	}
	if limit <= 0 {
		limit = 10
	}

	terms := tokenizeText(needle)

	var conditions []string
	var args []interface{}

	if len(terms) > 0 {
		placeholders := make([]string, len(terms))
		for i, term := range terms {
			placeholders[i] = "?"
			args = append(args, term)
		}
		conditions = append(conditions, fmt.Sprintf(`ci.id IN (
			SELECT card_id FROM card_term WHERE term IN (%s)
		)`, strings.Join(placeholders, ",")))
	}

	if len(typeFilter) > 0 {
		types := make([]string, 0, len(typeFilter))
		for ct := range typeFilter {
			types = append(types, "?")
			args = append(args, string(ct))
		}
		conditions = append(conditions, fmt.Sprintf("ci.type IN (%s)", strings.Join(types, ",")))
	}

	if strings.TrimSpace(statusFilter) != "" {
		conditions = append(conditions, "ci.status = ?")
		args = append(args, strings.TrimSpace(statusFilter))
	}

	if strings.TrimSpace(domainFilter) != "" {
		conditions = append(conditions, "ci.domain = ?")
		args = append(args, strings.TrimSpace(domainFilter))
	}

	if len(tagFilter) > 0 {
		placeholders := make([]string, len(tagFilter))
		for i, tag := range tagFilter {
			placeholders[i] = "?"
			args = append(args, strings.TrimSpace(tag))
		}
		conditions = append(conditions, fmt.Sprintf(`ci.id IN (
			SELECT card_id FROM card_tag WHERE tag IN (%s)
		)`, strings.Join(placeholders, ",")))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	sqlQuery := fmt.Sprintf(`
SELECT ci.id, ci.type, ci.title, ci.status, ci.importance, ci.source, ci.domain, ci.file_path, ci.updated_at, ci.created_at, ci.body
FROM card_index ci
%s
ORDER BY ci.id
LIMIT ?;`, whereClause)
	args = append(args, limit)

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("searching cards: %w", err)
	}
	defer rows.Close()

	cards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}

	results := make([]core.CardSearchResult, 0, len(cards))
	for _, card := range cards {
		reason := matchReason(card, needle, terms)
		results = append(results, core.CardSearchResult{
			Card:        card,
			MatchReason: reason,
		})
	}
	return results, nil
}

func (s *CardSyncService) SuggestLibraryCards(focus *core.Card, typeFilter map[core.CardType]bool, limit int) ([]core.LibrarySuggestion, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("sync service not initialized")
	}
	if limit <= 0 {
		limit = 10
	}

	terms := buildLibraryQueryTermsFromCard(focus)

	types := make([]string, 0, len(typeFilter))
	typeArgs := make([]interface{}, 0, len(typeFilter))
	for ct := range typeFilter {
		types = append(types, "?")
		typeArgs = append(typeArgs, string(ct))
	}

	termArgs := make([]interface{}, 0, len(terms))
	termPlaceholders := make([]string, 0, len(terms))
	for _, t := range terms {
		termPlaceholders = append(termPlaceholders, "?")
		termArgs = append(termArgs, t)
	}

	var args []interface{}
	args = append(args, typeArgs...)
	args = append(args, termArgs...)
	args = append(args, limit)

	sqlQuery := fmt.Sprintf(`
SELECT ci.id, ci.type, ci.title, ci.status, ci.importance, ci.source, ci.domain, ci.file_path, ci.updated_at, ci.created_at, ci.body,
	(SELECT COUNT(*) FROM card_term ct WHERE ct.card_id = ci.id AND ct.term IN (%s)) as term_count
FROM card_index ci
WHERE ci.type IN (%s)
  AND ci.status NOT IN ('deprecated', 'superseded')
  AND ci.id IN (SELECT card_id FROM card_term WHERE term IN (%s))
ORDER BY term_count DESC, ci.importance = 'must' DESC, ci.id
LIMIT ?;`,
		strings.Join(termPlaceholders, ","),
		strings.Join(types, ","),
		strings.Join(termPlaceholders, ","))

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("suggesting library cards: %w", err)
	}
	defer rows.Close()

	var suggestions []core.LibrarySuggestion
	for rows.Next() {
		card := &core.Card{}
		var updatedAt, createdAt string
		var termCount int
		if err := rows.Scan(
			&card.ID, &card.Type, &card.Title, &card.Status, &card.Importance,
			&card.Source, &card.Domain, &card.FilePath, &updatedAt, &createdAt, &card.Body,
			&termCount,
		); err != nil {
			return nil, fmt.Errorf("scanning suggestion: %w", err)
		}
		card.Updated, _ = time.Parse(time.RFC3339Nano, updatedAt)
		card.Created, _ = time.Parse(time.RFC3339Nano, createdAt)
		if card.Created.IsZero() {
			card.Created = card.Updated
		}
		suggestions = append(suggestions, core.LibrarySuggestion{
			Card:              card,
			Score:             termCount * 10,
			MatchReason:       fmt.Sprintf("matched %d terms", termCount),
			SuggestedRelation: suggestedRelationForType(card.Type),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating suggestions: %w", err)
	}

	return suggestions, nil
}

func (s *CardSyncService) CardSummary(id string) (string, error) {
	if s == nil || s.db == nil {
		return "", fmt.Errorf("sync service not initialized")
	}
	var summary string
	err := s.db.QueryRow(`SELECT summary FROM card_index WHERE id = ?;`, id).Scan(&summary)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return summary, err
}

func (s *CardSyncService) Backlinks(cardID string) ([]Backlink, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("store is not open")
	}
	if cardID == "" {
		return nil, fmt.Errorf("card ID is required")
	}

	rows, err := s.db.Query(`SELECT from_id, relation FROM card_link WHERE to_id = ? ORDER BY from_id, relation;`, cardID)
	if err != nil {
		return nil, fmt.Errorf("querying backlinks: %w", err)
	}
	defer rows.Close()

	var backlinks []Backlink
	for rows.Next() {
		var b Backlink
		if err := rows.Scan(&b.FromID, &b.Relation); err != nil {
			return nil, fmt.Errorf("scanning backlink: %w", err)
		}
		backlinks = append(backlinks, b)
	}
	return backlinks, rows.Err()
}

func (s *CardSyncService) DerivedIndexStatus() (DerivedIndexStatus, error) {
	if s == nil || s.db == nil {
		return DerivedIndexStatus{}, fmt.Errorf("store is not open")
	}

	status := DerivedIndexStatus{}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM card_index;`).Scan(&status.CardCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return status, nil
		}
		return DerivedIndexStatus{}, fmt.Errorf("counting card_index rows: %w", err)
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM card_link;`).Scan(&status.LinkCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return status, nil
		}
		return DerivedIndexStatus{}, fmt.Errorf("counting card_link rows: %w", err)
	}
	return status, nil
}

func (s *CardSyncService) readLinks(cardID string) []core.Link {
	rows, err := s.db.Query(`SELECT to_id, relation FROM card_link WHERE from_id = ? ORDER BY to_id;`, cardID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var links []core.Link
	for rows.Next() {
		var link core.Link
		if err := rows.Scan(&link.Target, &link.Relation); err != nil {
			continue
		}
		links = append(links, link)
	}
	return links
}

func (s *CardSyncService) readTags(cardID string) []string {
	rows, err := s.db.Query(`SELECT tag FROM card_tag WHERE card_id = ? ORDER BY tag;`, cardID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

func computeCardSummary(card *core.Card) string {
	if card == nil {
		return ""
	}
	body := strings.TrimSpace(card.Body)
	if body == "" {
		return ""
	}
	lines := strings.Split(body, "\n")
	var parts []string
	inParagraph := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			if inParagraph {
				break
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			if inParagraph {
				break
			}
			continue
		}
		parts = append(parts, line)
		inParagraph = true
	}
	result := strings.Join(parts, " ")
	if len(result) > 200 {
		result = result[:200]
	}
	return result
}

type cardTerm struct {
	term   string
	source string
}

func buildCardTerms(card *core.Card) []cardTerm {
	if card == nil {
		return nil
	}
	var terms []cardTerm
	seen := map[string]bool{}

	add := func(source, text string) {
		for _, t := range tokenizeText(text) {
			if t == "" || seen[t+"@"+source] {
				continue
			}
			seen[t+"@"+source] = true
			terms = append(terms, cardTerm{term: t, source: source})
		}
	}

	add("title", card.Title)
	for _, tag := range card.Tags {
		add("tag", tag)
	}
	add("domain", card.Domain)
	for _, w := range significantWords(card.Body) {
		add("body", w)
	}
	return terms
}

func tokenizeText(text string) []string {
	lower := strings.ToLower(text)
	runes := []rune(lower)
	tokens := make([]string, 0, len(runes)*2)
	seen := map[string]bool{}

	var latinBuf strings.Builder
	flushLatin := func() {
		if latinBuf.Len() == 0 {
			return
		}
		t := latinBuf.String()
		latinBuf.Reset()
		if t == "" || seen[t] {
			return
		}
		seen[t] = true
		tokens = append(tokens, t)
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if isCJK(r) {
			flushLatin()
			unigram := string(r)
			if !seen[unigram] {
				seen[unigram] = true
				tokens = append(tokens, unigram)
			}
			if i+1 < len(runes) && isCJK(runes[i+1]) {
				bigram := string(runes[i : i+2])
				if !seen[bigram] {
					seen[bigram] = true
					tokens = append(tokens, bigram)
				}
			}
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
			latinBuf.WriteRune(r)
		} else {
			flushLatin()
		}
	}
	flushLatin()
	return tokens
}

func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r)
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
		runes := []rune(term)
		if len(runes) == 1 && isCJK(runes[0]) {
			continue
		}
		if len(runes) == 1 && stop[term] {
			continue
		}
		if len(runes) < 2 {
			continue
		}
		if len(runes) >= 2 && len(runes) <= 3 && !hasCJK(term) && stop[term] {
			continue
		}
		filtered = append(filtered, term)
	}
	return filtered
}

func hasCJK(s string) bool {
	for _, r := range s {
		if isCJK(r) {
			return true
		}
	}
	return false
}

func buildLibraryQueryTermsFromCard(card *core.Card) []string {
	if card == nil {
		return nil
	}
	terms := make([]string, 0, 16)
	terms = append(terms, tokenizeText(card.Title)...)
	terms = append(terms, tokenizeText(card.Domain)...)
	for _, tag := range card.Tags {
		terms = append(terms, tokenizeText(tag)...)
	}
	terms = append(terms, significantWords(card.Body)...)
	seen := map[string]bool{}
	result := make([]string, 0, len(terms))
	for _, t := range terms {
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		result = append(result, t)
	}
	return result
}

func matchReason(card *core.Card, needle string, terms []string) string {
	lowerID := strings.ToLower(card.ID)
	lowerTitle := strings.ToLower(card.Title)
	lowerBody := strings.ToLower(card.Body)

	var reasons []string
	if strings.Contains(lowerID, needle) {
		reasons = append(reasons, "matched id")
	}
	if strings.Contains(lowerTitle, needle) {
		reasons = append(reasons, "matched title")
	}
	if strings.Contains(lowerBody, needle) {
		reasons = append(reasons, "matched body")
	}
	if len(reasons) == 0 && len(terms) > 0 {
		reasons = append(reasons, fmt.Sprintf("matched %d terms", len(terms)))
	}
	if len(reasons) == 0 {
		return "matched"
	}
	return strings.Join(reasons, " | ")
}

func suggestedRelationForType(cardType core.CardType) string {
	switch cardType {
	case core.CardTypeConvention, core.CardTypeModule:
		return "constrains"
	case core.CardTypeDecision, core.CardTypeDesign, core.CardTypeFinding:
		return "references"
	default:
		return "related"
	}
}

func scanCard(row *sql.Row) (*core.Card, error) {
	card := &core.Card{}
	var updatedAt, createdAt string
	err := row.Scan(
		&card.ID, &card.Type, &card.Title, &card.Status, &card.Importance,
		&card.Source, &card.Domain, &card.FilePath, &updatedAt, &createdAt, &card.Body,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("card not found")
	}
	if err != nil {
		return nil, fmt.Errorf("scanning card: %w", err)
	}
	card.Updated, _ = time.Parse(time.RFC3339Nano, updatedAt)
	card.Created, _ = time.Parse(time.RFC3339Nano, createdAt)
	if card.Created.IsZero() {
		card.Created = card.Updated
	}
	return card, nil
}

func scanCards(rows *sql.Rows) ([]*core.Card, error) {
	var cards []*core.Card
	for rows.Next() {
		card := &core.Card{}
		var updatedAt, createdAt string
		if err := rows.Scan(
			&card.ID, &card.Type, &card.Title, &card.Status, &card.Importance,
			&card.Source, &card.Domain, &card.FilePath, &updatedAt, &createdAt, &card.Body,
		); err != nil {
			return nil, fmt.Errorf("scanning card row: %w", err)
		}
		card.Updated, _ = time.Parse(time.RFC3339Nano, updatedAt)
		card.Created, _ = time.Parse(time.RFC3339Nano, createdAt)
		if card.Created.IsZero() {
			card.Created = card.Updated
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (s *CardSyncService) hydrateLinks(cards []*core.Card) error {
	if len(cards) == 0 || s == nil || s.db == nil {
		return nil
	}
	ids := make([]string, len(cards))
	idx := map[string]int{}
	for i, card := range cards {
		ids[i] = card.ID
		idx[card.ID] = i
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	rows, err := s.db.Query(
		fmt.Sprintf("SELECT from_id, to_id, relation FROM card_link WHERE from_id IN (%s) ORDER BY from_id, to_id;",
			strings.Join(placeholders, ",")),
		args...,
	)
	if err != nil {
		return fmt.Errorf("querying links: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fromID, toID, relation string
		if err := rows.Scan(&fromID, &toID, &relation); err != nil {
			continue
		}
		if i, ok := idx[fromID]; ok {
			cards[i].Links = append(cards[i].Links, core.Link{Target: toID, Relation: relation})
		}
	}
	return rows.Err()
}

func scanCardFromRows(rows *sql.Rows) (*core.Card, error) {
	card := &core.Card{}
	var updatedAt, createdAt string
	if err := rows.Scan(
		&card.ID, &card.Type, &card.Title, &card.Status, &card.Importance,
		&card.Source, &card.Domain, &card.FilePath, &updatedAt, &createdAt, &card.Body,
	); err != nil {
		return nil, err
	}
	card.Updated, _ = time.Parse(time.RFC3339Nano, updatedAt)
	card.Created, _ = time.Parse(time.RFC3339Nano, createdAt)
	if card.Created.IsZero() {
		card.Created = card.Updated
	}
	return card, nil
}
