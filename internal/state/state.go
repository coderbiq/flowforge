package state

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const currentProjectKey = "currentProjectId"

type Store struct {
	db *sql.DB
}

func Open(dbPath string) (*Store, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("database path is required")
	}

	parentDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("pinging sqlite database: %w (closing db: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("pinging sqlite database: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) EnsureSchema() error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is not open")
	}

	const query = `
CREATE TABLE IF NOT EXISTS runtime_state (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS card_index (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	title TEXT NOT NULL,
	status TEXT NOT NULL,
	importance TEXT NOT NULL,
	source TEXT NOT NULL,
	domain TEXT NOT NULL,
	file_path TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS card_link (
	from_id TEXT NOT NULL,
	to_id TEXT NOT NULL,
	relation TEXT NOT NULL
);`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("ensuring runtime_state schema: %w", err)
	}

	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("closing sqlite database: %w", err)
	}

	s.db = nil
	return nil
}

func (s *Store) SetCurrentProjectID(id string) error {
	if id == "" {
		return fmt.Errorf("project ID is required")
	}

	return s.setValue(currentProjectKey, id)
}

func (s *Store) CurrentProjectID() (string, bool, error) {
	return s.value(currentProjectKey)
}

func (s *Store) ClearCurrentProjectID() error {
	return s.deleteValue(currentProjectKey)
}

func (s *Store) SetCurrentProposalID(projectID, proposalID string) error {
	if projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	if proposalID == "" {
		return fmt.Errorf("proposal ID is required")
	}

	return s.setValue(currentProposalKey(projectID), proposalID)
}

func (s *Store) CurrentProposalID(projectID string) (string, bool, error) {
	if projectID == "" {
		return "", false, fmt.Errorf("project ID is required")
	}

	return s.value(currentProposalKey(projectID))
}

func (s *Store) ClearCurrentProposalID(projectID string) error {
	if projectID == "" {
		return fmt.Errorf("project ID is required")
	}

	return s.deleteValue(currentProposalKey(projectID))
}

func currentProposalKey(projectID string) string {
	return fmt.Sprintf("project:%s:current-proposal", projectID)
}

func (s *Store) setValue(key string, value string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is not open")
	}
	if key == "" {
		return fmt.Errorf("state key is required")
	}

	const query = `
INSERT INTO runtime_state(key, value, updated_at)
VALUES (?, ?, ?)
ON CONFLICT(key) DO UPDATE SET
	value = excluded.value,
	updated_at = excluded.updated_at;`

	if _, err := s.db.Exec(query, key, value, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		return fmt.Errorf("setting runtime state %q: %w", key, err)
	}

	return nil
}

func (s *Store) value(key string) (string, bool, error) {
	if s == nil || s.db == nil {
		return "", false, fmt.Errorf("store is not open")
	}
	if key == "" {
		return "", false, fmt.Errorf("state key is required")
	}

	var value string
	err := s.db.QueryRow("SELECT value FROM runtime_state WHERE key = ?", key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("reading runtime state %q: %w", key, err)
	}

	return value, true, nil
}

func (s *Store) deleteValue(key string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is not open")
	}
	if key == "" {
		return fmt.Errorf("state key is required")
	}

	if _, err := s.db.Exec("DELETE FROM runtime_state WHERE key = ?", key); err != nil {
		return fmt.Errorf("clearing runtime state %q: %w", key, err)
	}

	return nil
}
