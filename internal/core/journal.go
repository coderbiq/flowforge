package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const proposalJournalFileName = "JOURNAL.md"

const journalHeader = "# Proposal Journal\n\n" +
	"Chronological collaboration notes for this proposal. Formal design, progress, and verification remain in their referenced artifacts.\n"

const (
	journalEntryStart = "<!-- flowforge:journal-entry -->"
	journalEntryEnd   = "<!-- /flowforge:journal-entry -->"
)

var journalLocks sync.Map

// JournalEntry is a concise, durable handoff note for a proposal. It never
// stores full agent output; referenced artifacts remain the source of truth.
type JournalEntry struct {
	Time          time.Time
	Actor         string
	Message       string
	References    []string
	Status        string
	Next          string
	BlockedReason string
}

func (s *CardStore) ProposalJournalPath(proposalID string) string {
	return filepath.Join(s.ProposalDir(proposalID), proposalJournalFileName)
}

func (s *CardStore) CreateProposalJournal(proposalID string) (string, error) {
	proposalDir := s.ProposalDir(proposalID)
	if err := os.MkdirAll(proposalDir, 0755); err != nil {
		return "", fmt.Errorf("creating proposal directory: %w", err)
	}

	path := s.ProposalJournalPath(proposalID)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return path, nil
		}
		return "", fmt.Errorf("creating proposal journal: %w", err)
	}
	if _, err := file.WriteString(journalHeader); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return "", fmt.Errorf("writing proposal journal header: %w (closing proposal journal: %v)", err, closeErr)
		}
		return "", fmt.Errorf("writing proposal journal header: %w", err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("closing proposal journal: %w", err)
	}
	return path, nil
}

func (s *CardStore) AppendProposalJournal(proposalID string, entry JournalEntry) error {
	if err := validateJournalEntry(entry); err != nil {
		return err
	}
	if entry.Time.IsZero() {
		entry.Time = time.Now().UTC()
	}

	lock := journalLock(s.ProposalJournalPath(proposalID))
	lock.Lock()
	defer lock.Unlock()

	path, err := s.CreateProposalJournal(proposalID)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("opening proposal journal: %w", err)
	}

	if _, err := file.WriteString(renderJournalEntry(entry)); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return fmt.Errorf("appending proposal journal: %w (closing proposal journal: %v)", err, closeErr)
		}
		return fmt.Errorf("appending proposal journal: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing proposal journal: %w", err)
	}
	return nil
}

func journalLock(path string) *sync.Mutex {
	lock, _ := journalLocks.LoadOrStore(path, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func (s *CardStore) RecentProposalJournal(proposalID string, limit int) ([]JournalEntry, error) {
	if limit < 0 {
		return nil, fmt.Errorf("journal entry limit cannot be negative")
	}

	data, err := os.ReadFile(s.ProposalJournalPath(proposalID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading proposal journal: %w", err)
	}

	entries, err := parseJournalEntries(string(data))
	if err != nil {
		return nil, err
	}
	if limit == 0 || len(entries) <= limit {
		return entries, nil
	}
	return entries[len(entries)-limit:], nil
}

func validateJournalEntry(entry JournalEntry) error {
	if strings.TrimSpace(entry.Actor) == "" {
		return fmt.Errorf("journal actor is required")
	}
	if strings.TrimSpace(entry.Message) == "" {
		return fmt.Errorf("journal message is required")
	}
	return nil
}

func renderJournalEntry(entry JournalEntry) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(journalEntryStart)
	b.WriteString("\n## ")
	b.WriteString(entry.Time.Format(time.RFC3339Nano))
	b.WriteString(" ")
	b.WriteString(journalValue(entry.Actor))
	b.WriteString("\n\n- Summary: ")
	b.WriteString(journalValue(entry.Message))
	b.WriteString("\n")
	if len(entry.References) > 0 {
		refs := make([]string, 0, len(entry.References))
		for _, reference := range entry.References {
			if value := journalValue(reference); value != "" {
				refs = append(refs, value)
			}
		}
		if len(refs) > 0 {
			b.WriteString("- References: ")
			b.WriteString(strings.Join(refs, ", "))
			b.WriteString("\n")
		}
	}
	if value := journalValue(entry.Status); value != "" {
		b.WriteString("- Status: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	if value := journalValue(entry.BlockedReason); value != "" {
		b.WriteString("- Blocked: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	if value := journalValue(entry.Next); value != "" {
		b.WriteString("- Next: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	b.WriteString(journalEntryEnd)
	b.WriteString("\n")
	return b.String()
}

func parseJournalEntries(content string) ([]JournalEntry, error) {
	var entries []JournalEntry
	remaining := content
	for {
		start := strings.Index(remaining, journalEntryStart)
		if start == -1 {
			return entries, nil
		}
		remaining = remaining[start+len(journalEntryStart):]
		end := strings.Index(remaining, journalEntryEnd)
		if end == -1 {
			return nil, fmt.Errorf("parsing proposal journal: unterminated entry")
		}
		entry, err := parseJournalEntry(strings.TrimSpace(remaining[:end]))
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		remaining = remaining[end+len(journalEntryEnd):]
	}
}

func parseJournalEntry(content string) (JournalEntry, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "## ") {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal: entry header is missing")
	}

	header := strings.TrimPrefix(lines[0], "## ")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal: entry actor is missing")
	}
	entryTime, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal timestamp: %w", err)
	}

	entry := JournalEntry{Time: entryTime, Actor: parts[1]}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "- Summary: "):
			entry.Message = strings.TrimPrefix(line, "- Summary: ")
		case strings.HasPrefix(line, "- References: "):
			entry.References = splitJournalReferences(strings.TrimPrefix(line, "- References: "))
		case strings.HasPrefix(line, "- Status: "):
			entry.Status = strings.TrimPrefix(line, "- Status: ")
		case strings.HasPrefix(line, "- Blocked: "):
			entry.BlockedReason = strings.TrimPrefix(line, "- Blocked: ")
		case strings.HasPrefix(line, "- Next: "):
			entry.Next = strings.TrimPrefix(line, "- Next: ")
		}
	}
	if err := validateJournalEntry(entry); err != nil {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal: %w", err)
	}
	return entry, nil
}

func splitJournalReferences(value string) []string {
	parts := strings.Split(value, ",")
	references := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			references = append(references, part)
		}
	}
	return references
}

func journalValue(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
