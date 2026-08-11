package core

import (
	"strings"
	"testing"
)

func TestHandoffJournalBindsIdempotentlyAndBecomesReadOnly(t *testing.T) {
	projectRoot := t.TempDir()
	handoffs := NewHandoffJournalStore(projectRoot)
	journal, err := handoffs.Create("Investigate before proposal")
	if err != nil {
		t.Fatalf("create handoff journal: %v", err)
	}
	first, err := handoffs.Append(journal.ID, HandoffJournalEntry{Actor: "coordinator", Kind: "delegation", WorkID: "W1", Message: "Inspect the failing path", References: []string{"GIIS-2219"}})
	if err != nil {
		t.Fatalf("append delegation: %v", err)
	}
	if _, err := handoffs.Append(journal.ID, HandoffJournalEntry{Actor: "investigator", Kind: "result", WorkID: "W1", Message: "Root cause found"}); err != nil {
		t.Fatalf("append result: %v", err)
	}

	cardStore := NewCardStore(t.TempDir())
	if _, _, err := cardStore.CreateProposal("CR26081101", "Journal binding"); err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	imported, err := handoffs.Bind(journal.ID, "CR26081101", cardStore)
	if err != nil {
		t.Fatalf("bind handoff journal: %v", err)
	}
	if imported != 2 {
		t.Fatalf("imported = %d, want 2", imported)
	}
	entries, err := cardStore.RecentProposalJournal("CR26081101", 0)
	if err != nil {
		t.Fatalf("read proposal journal: %v", err)
	}
	if len(entries) != 2 || entries[0].EventID != first.ID || entries[0].Source != journal.ID {
		t.Fatalf("unexpected imported entries: %#v", entries)
	}
	if imported, err = handoffs.Bind(journal.ID, "CR26081101", cardStore); err != nil || imported != 0 {
		t.Fatalf("idempotent bind = (%d, %v), want (0, nil)", imported, err)
	}
	if _, err := handoffs.Append(journal.ID, HandoffJournalEntry{Actor: "worker", Kind: "result", Message: "late write"}); err == nil || !strings.Contains(err.Error(), "append to proposal") {
		t.Fatalf("expected bound journal append rejection, got %v", err)
	}
}

func TestHandoffJournalRequiresExplicitSupportedKind(t *testing.T) {
	store := NewHandoffJournalStore(t.TempDir())
	journal, err := store.Create("Bounded handoff")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Append(journal.ID, HandoffJournalEntry{Actor: "worker", Kind: "analysis", Message: "too broad"}); err == nil {
		t.Fatal("expected unsupported kind error")
	}
}

func TestHandoffJournalRejectsTraversalAndOversizedMessages(t *testing.T) {
	store := NewHandoffJournalStore(t.TempDir())
	if _, err := store.Read("../../outside"); err == nil {
		t.Fatal("expected invalid journal ID error")
	}
	journal, err := store.Create("Bounded handoff")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Append(journal.ID, HandoffJournalEntry{Actor: "worker", Kind: "result", Message: strings.Repeat("x", 1001)}); err == nil {
		t.Fatal("expected oversized message error")
	}
}
