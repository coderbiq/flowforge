package state

import (
	"database/sql"
	"fmt"
	"time"

	"flowforge/internal/core"
)

type DerivedIndexStatus struct {
	CardCount int
	LinkCount int
}

type Backlink struct {
	FromID   string
	Relation string
}

func (s *Store) RebuildDerivedIndex(cards []*core.Card) (int, int, error) {
	if s == nil || s.db == nil {
		return 0, 0, fmt.Errorf("store is not open")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("starting index rebuild transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.Exec(`DELETE FROM card_index;`); err != nil {
		return 0, 0, fmt.Errorf("clearing card_index: %w", err)
	}
	if _, err := tx.Exec(`DELETE FROM card_link;`); err != nil {
		return 0, 0, fmt.Errorf("clearing card_link: %w", err)
	}

	cardStmt, err := tx.Prepare(`
INSERT INTO card_index(id, type, title, status, importance, source, domain, file_path, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return 0, 0, fmt.Errorf("preparing card index insert: %w", err)
	}
	defer cardStmt.Close()

	linkStmt, err := tx.Prepare(`
INSERT INTO card_link(from_id, to_id, relation)
VALUES (?, ?, ?);`)
	if err != nil {
		return 0, 0, fmt.Errorf("preparing card link insert: %w", err)
	}
	defer linkStmt.Close()

	seen := map[string]bool{}
	cardCount := 0
	linkCount := 0

	for _, card := range cards {
		if card == nil {
			continue
		}
		if seen[card.ID] {
			return 0, 0, fmt.Errorf("duplicate card ID in rebuild input: %s", card.ID)
		}
		seen[card.ID] = true

		if _, err := cardStmt.Exec(card.ID, string(card.Type), card.Title, string(card.Status), string(card.Importance), card.Source, card.Domain, card.FilePath, card.Updated.UTC().Format(time.RFC3339Nano)); err != nil {
			return 0, 0, fmt.Errorf("inserting card %s: %w", card.ID, err)
		}
		cardCount++

		for _, link := range card.Links {
			if link.Target == "" {
				continue
			}
			if _, err := linkStmt.Exec(card.ID, link.Target, link.Relation); err != nil {
				return 0, 0, fmt.Errorf("inserting link from %s to %s: %w", card.ID, link.Target, err)
			}
			linkCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("committing index rebuild: %w", err)
	}
	tx = nil

	return cardCount, linkCount, nil
}

func (s *Store) DerivedIndexStatus() (DerivedIndexStatus, error) {
	if s == nil || s.db == nil {
		return DerivedIndexStatus{}, fmt.Errorf("store is not open")
	}

	status := DerivedIndexStatus{}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM card_index;`).Scan(&status.CardCount); err != nil {
		if err == sql.ErrNoRows {
			return status, nil
		}
		return DerivedIndexStatus{}, fmt.Errorf("counting card_index rows: %w", err)
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM card_link;`).Scan(&status.LinkCount); err != nil {
		if err == sql.ErrNoRows {
			return status, nil
		}
		return DerivedIndexStatus{}, fmt.Errorf("counting card_link rows: %w", err)
	}

	return status, nil
}

func (s *Store) Backlinks(cardID string) ([]Backlink, error) {
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
		var backlink Backlink
		if err := rows.Scan(&backlink.FromID, &backlink.Relation); err != nil {
			return nil, fmt.Errorf("scanning backlink: %w", err)
		}
		backlinks = append(backlinks, backlink)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating backlinks: %w", err)
	}

	return backlinks, nil
}
