package core

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAppendProposalJournalCreatesMissingJournal(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	proposalID := "CR260810"
	proposalDir := filepath.Join(wikiRoot, "01-workspace", proposalID)
	if err := os.MkdirAll(proposalDir, 0755); err != nil {
		t.Fatalf("creating proposal directory: %v", err)
	}

	entry := JournalEntry{
		Time:    time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC),
		Actor:   "design-analyst",
		Message: "完成设计",
		References: []string{
			"FEAT-001",
			"DEC-001",
		},
		Status: "designed",
		Next:   "等待确认",
	}
	if err := store.AppendProposalJournal(proposalID, entry); err != nil {
		t.Fatalf("appending proposal journal: %v", err)
	}

	content, err := os.ReadFile(store.ProposalJournalPath(proposalID))
	if err != nil {
		t.Fatalf("reading proposal journal: %v", err)
	}
	for _, want := range []string{
		journalHeader,
		"## 2026-08-10T12:00:00Z design-analyst",
		"- Summary: 完成设计",
		"- References: FEAT-001, DEC-001",
		"- Status: designed",
		"- Next: 等待确认",
	} {
		if !strings.Contains(string(content), want) {
			t.Fatalf("journal missing %q:\n%s", want, content)
		}
	}
}

func TestRecentProposalJournalReturnsTailInOrder(t *testing.T) {
	store, _ := setupTestStore(t)
	proposalID := "CR260810"
	times := []time.Time{
		time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC),
		time.Date(2026, time.August, 10, 12, 1, 0, 0, time.UTC),
		time.Date(2026, time.August, 10, 12, 2, 0, 0, time.UTC),
	}
	for i, entryTime := range times {
		if err := store.AppendProposalJournal(proposalID, JournalEntry{
			Time:    entryTime,
			Actor:   "executor",
			Message: "step " + string(rune('1'+i)),
		}); err != nil {
			t.Fatalf("appending entry %d: %v", i, err)
		}
	}

	entries, err := store.RecentProposalJournal(proposalID, 2)
	if err != nil {
		t.Fatalf("reading recent entries: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Time != times[1] || entries[1].Time != times[2] {
		t.Fatalf("expected tail entries in order, got %#v", entries)
	}
}

func TestRecentProposalJournalHandlesMissingJournal(t *testing.T) {
	store, _ := setupTestStore(t)
	entries, err := store.RecentProposalJournal("CR260810", 10)
	if err != nil {
		t.Fatalf("reading missing journal: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no entries, got %#v", entries)
	}
}

func TestAppendProposalJournalNormalizesMarkdownValues(t *testing.T) {
	store, _ := setupTestStore(t)
	if err := store.AppendProposalJournal("CR260810", JournalEntry{
		Time:          time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC),
		Actor:         " design-analyst ",
		Message:       "line one\nline two",
		References:    []string{" FEAT-001 ", "", "DEC-001"},
		BlockedReason: "waiting\nfor dependency",
	}); err != nil {
		t.Fatalf("appending journal: %v", err)
	}

	entries, err := store.RecentProposalJournal("CR260810", 1)
	if err != nil {
		t.Fatalf("reading journal: %v", err)
	}
	want := JournalEntry{
		Time:          time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC),
		Actor:         "design-analyst",
		Message:       "line one line two",
		References:    []string{"FEAT-001", "DEC-001"},
		BlockedReason: "waiting for dependency",
	}
	if !reflect.DeepEqual(entries[0], want) {
		t.Fatalf("unexpected entry:\n got: %#v\nwant: %#v", entries[0], want)
	}
}

func TestAppendProposalJournalRejectsMissingRequiredValues(t *testing.T) {
	store, _ := setupTestStore(t)
	for _, entry := range []JournalEntry{
		{Message: "missing actor"},
		{Actor: "executor"},
	} {
		if err := store.AppendProposalJournal("CR260810", entry); err == nil {
			t.Fatalf("expected validation error for %#v", entry)
		}
	}
}

func TestRecentProposalJournalRejectsMalformedManagedEntry(t *testing.T) {
	store, _ := setupTestStore(t)
	path, err := store.CreateProposalJournal("CR260810")
	if err != nil {
		t.Fatalf("creating journal: %v", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("opening journal: %v", err)
	}
	if _, err := file.WriteString("\n" + journalEntryStart + "\n## invalid\n" + journalEntryEnd + "\n"); err != nil {
		t.Fatalf("writing malformed entry: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("closing journal: %v", err)
	}

	if _, err := store.RecentProposalJournal("CR260810", 1); err == nil {
		t.Fatal("expected malformed journal error")
	}
}

func TestAppendProposalJournalPreservesConcurrentEntries(t *testing.T) {
	store, _ := setupTestStore(t)
	const entryCount = 20

	var group sync.WaitGroup
	errs := make(chan error, entryCount)
	for i := 0; i < entryCount; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			errs <- store.AppendProposalJournal("CR260810", JournalEntry{
				Actor:   "executor",
				Message: "concurrent entry " + string(rune('a'+i)),
			})
		}(i)
	}
	group.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent append failed: %v", err)
		}
	}

	entries, err := store.RecentProposalJournal("CR260810", 0)
	if err != nil {
		t.Fatalf("reading journal: %v", err)
	}
	if len(entries) != entryCount {
		t.Fatalf("expected %d entries, got %d", entryCount, len(entries))
	}
}
