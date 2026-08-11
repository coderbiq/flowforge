package core

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const handoffJournalVersion = 1

type HandoffJournal struct {
	ID              string    `json:"id"`
	Version         int       `json:"version"`
	Title           string    `json:"title"`
	State           string    `json:"state"`
	CreatedAt       time.Time `json:"createdAt"`
	BoundProposalID string    `json:"boundProposalId,omitempty"`
	BoundAt         time.Time `json:"boundAt,omitempty"`
}

type HandoffJournalEntry struct {
	ID            string    `json:"id"`
	Time          time.Time `json:"time"`
	Actor         string    `json:"actor"`
	Kind          string    `json:"kind"`
	WorkID        string    `json:"workId,omitempty"`
	Message       string    `json:"message"`
	References    []string  `json:"references,omitempty"`
	Status        string    `json:"status,omitempty"`
	Next          string    `json:"next,omitempty"`
	BlockedReason string    `json:"blocked,omitempty"`
}

type HandoffJournalStore struct {
	root string
}

func NewHandoffJournalStore(projectRoot string) *HandoffJournalStore {
	return &HandoffJournalStore{root: filepath.Join(projectRoot, ".flowforge", "handoff-journals")}
}

func (s *HandoffJournalStore) Create(title string) (HandoffJournal, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return HandoffJournal{}, fmt.Errorf("handoff journal title is required")
	}
	id, err := newHandoffID("JRN")
	if err != nil {
		return HandoffJournal{}, err
	}
	journal := HandoffJournal{ID: id, Version: handoffJournalVersion, Title: title, State: "open", CreatedAt: time.Now().UTC()}
	if err := os.MkdirAll(s.journalDir(id), 0755); err != nil {
		return HandoffJournal{}, fmt.Errorf("creating handoff journal directory: %w", err)
	}
	if err := s.writeMetadata(journal); err != nil {
		return HandoffJournal{}, err
	}
	return journal, nil
}

func (s *HandoffJournalStore) Read(id string) (HandoffJournal, error) {
	if !validHandoffID(id, "JRN") {
		return HandoffJournal{}, fmt.Errorf("invalid handoff journal ID %q", id)
	}
	data, err := os.ReadFile(s.metadataPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return HandoffJournal{}, fmt.Errorf("handoff journal %q does not exist", id)
		}
		return HandoffJournal{}, fmt.Errorf("reading handoff journal metadata: %w", err)
	}
	var journal HandoffJournal
	if err := json.Unmarshal(data, &journal); err != nil {
		return HandoffJournal{}, fmt.Errorf("parsing handoff journal metadata: %w", err)
	}
	return journal, nil
}

func (s *HandoffJournalStore) Append(id string, entry HandoffJournalEntry) (HandoffJournalEntry, error) {
	journal, err := s.Read(id)
	if err != nil {
		return HandoffJournalEntry{}, err
	}
	if journal.State != "open" {
		return HandoffJournalEntry{}, fmt.Errorf("handoff journal %s is %s; append to proposal %s instead", id, journal.State, journal.BoundProposalID)
	}
	if strings.TrimSpace(entry.Actor) == "" {
		return HandoffJournalEntry{}, fmt.Errorf("handoff journal actor is required")
	}
	if strings.TrimSpace(entry.Message) == "" {
		return HandoffJournalEntry{}, fmt.Errorf("handoff journal message is required")
	}
	if len([]rune(entry.Message)) > 1000 {
		return HandoffJournalEntry{}, fmt.Errorf("handoff journal message exceeds 1000 characters; store detailed evidence in an artifact and reference it")
	}
	if len(entry.References) > 20 {
		return HandoffJournalEntry{}, fmt.Errorf("handoff journal entry has %d references (maximum: 20)", len(entry.References))
	}
	entry.Kind = strings.TrimSpace(entry.Kind)
	if entry.Kind == "" {
		return HandoffJournalEntry{}, fmt.Errorf("handoff journal kind is required")
	}
	if !map[string]bool{"delegation": true, "result": true, "blocked": true, "synthesis": true}[entry.Kind] {
		return HandoffJournalEntry{}, fmt.Errorf("unsupported handoff journal kind %q", entry.Kind)
	}
	if entry.ID == "" {
		entry.ID, err = newHandoffID("JEV")
		if err != nil {
			return HandoffJournalEntry{}, err
		}
	}
	if entry.Time.IsZero() {
		entry.Time = time.Now().UTC()
	}
	encoded, err := json.Marshal(entry)
	if err != nil {
		return HandoffJournalEntry{}, fmt.Errorf("encoding handoff journal entry: %w", err)
	}
	file, err := os.OpenFile(s.entriesPath(id), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return HandoffJournalEntry{}, fmt.Errorf("opening handoff journal entries: %w", err)
	}
	if _, err := file.Write(append(encoded, '\n')); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return HandoffJournalEntry{}, fmt.Errorf("appending handoff journal entry: %w (closing: %v)", err, closeErr)
		}
		return HandoffJournalEntry{}, fmt.Errorf("appending handoff journal entry: %w", err)
	}
	if err := file.Close(); err != nil {
		return HandoffJournalEntry{}, fmt.Errorf("closing handoff journal entries: %w", err)
	}
	return entry, nil
}

func (s *HandoffJournalStore) Entries(id string, limit int) ([]HandoffJournalEntry, error) {
	if limit < 0 {
		return nil, fmt.Errorf("handoff journal entry limit cannot be negative")
	}
	if _, err := s.Read(id); err != nil {
		return nil, err
	}
	file, err := os.Open(s.entriesPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("opening handoff journal entries: %w", err)
	}
	var entries []HandoffJournalEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry HandoffJournalEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, fmt.Errorf("parsing handoff journal entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return nil, fmt.Errorf("reading handoff journal entries: %w (closing: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("reading handoff journal entries: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("closing handoff journal entries: %w", err)
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Time.Before(entries[j].Time) })
	if limit > 0 && len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}
	return entries, nil
}

func (s *HandoffJournalStore) Bind(id, proposalID string, cardStore *CardStore) (int, error) {
	journal, err := s.Read(id)
	if err != nil {
		return 0, err
	}
	if journal.State == "bound" {
		if journal.BoundProposalID != proposalID {
			return 0, fmt.Errorf("handoff journal %s is already bound to proposal %s", id, journal.BoundProposalID)
		}
		return 0, nil
	}
	entries, err := s.Entries(id, 0)
	if err != nil {
		return 0, err
	}
	existing, err := cardStore.RecentProposalJournal(proposalID, 0)
	if err != nil {
		return 0, err
	}
	seen := map[string]bool{}
	for _, entry := range existing {
		seen[entry.EventID] = true
	}
	imported := 0
	for _, entry := range entries {
		if seen[entry.ID] {
			continue
		}
		proposalEntry := JournalEntry{EventID: entry.ID, Source: id, Time: entry.Time, Actor: entry.Actor, Message: entry.Message, References: entry.References, Status: entry.Status, Next: entry.Next, BlockedReason: entry.BlockedReason}
		if err := cardStore.AppendProposalJournal(proposalID, proposalEntry); err != nil {
			return imported, fmt.Errorf("importing handoff journal entry %s: %w", entry.ID, err)
		}
		imported++
	}
	journal.State = "bound"
	journal.BoundProposalID = proposalID
	journal.BoundAt = time.Now().UTC()
	if err := s.writeMetadata(journal); err != nil {
		return imported, err
	}
	return imported, nil
}

func (s *HandoffJournalStore) journalDir(id string) string { return filepath.Join(s.root, id) }
func (s *HandoffJournalStore) metadataPath(id string) string {
	return filepath.Join(s.journalDir(id), "journal.json")
}
func (s *HandoffJournalStore) entriesPath(id string) string {
	return filepath.Join(s.journalDir(id), "entries.jsonl")
}

func (s *HandoffJournalStore) writeMetadata(journal HandoffJournal) error {
	data, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding handoff journal metadata: %w", err)
	}
	if err := os.WriteFile(s.metadataPath(journal.ID), append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("writing handoff journal metadata: %w", err)
	}
	return nil
}

func newHandoffID(prefix string) (string, error) {
	var random [6]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("generating handoff journal ID: %w", err)
	}
	return prefix + "-" + strings.ToLower(hex.EncodeToString(random[:])), nil
}

func validHandoffID(value, prefix string) bool {
	if !strings.HasPrefix(value, prefix+"-") || len(value) != len(prefix)+1+12 {
		return false
	}
	for _, r := range value[len(prefix)+1:] {
		if !((r >= 'a' && r <= 'f') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}
