package core

import (
	"strings"
	"sync"
	"testing"
)

func TestNextProposalCardIDUsesPaddedSequence(t *testing.T) {
	store, _ := setupTestStore(t)
	createTestDirs(t, store.wikiRoot)
	createTestProposal(t, store, "CR26081102", "Sequence proposal")

	first, err := store.NextCardID(CardTypeFeature, "CR26081102")
	if err != nil {
		t.Fatal(err)
	}
	if first != "FEAT-CR26081102-001" {
		t.Fatalf("unexpected first ID: %s", first)
	}

	card := NewCard(CardTypeFeature, "First")
	card.ID = first
	if _, err := store.CreateCard(card, "CR26081102"); err != nil {
		t.Fatal(err)
	}

	second, err := store.NextTaskID("CR26081102", "i")
	if err != nil {
		t.Fatal(err)
	}
	if second != "TASK-CR26081102-i-002" {
		t.Fatalf("unexpected second ID: %s", second)
	}
}

func TestNextProposalCardIDRemainsUniqueAcrossConcurrentAllocations(t *testing.T) {
	store, _ := setupTestStore(t)
	createTestDirs(t, store.wikiRoot)
	createTestProposal(t, store, "CR26081103", "Concurrent sequence proposal")

	const count = 20
	ids := make(chan string, count)
	errs := make(chan error, count)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := store.NextCardID(CardTypeFeature, "CR26081103")
			if err != nil {
				errs <- err
				return
			}
			ids <- id
		}()
	}
	wg.Wait()
	close(ids)
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}

	seen := map[string]bool{}
	for id := range ids {
		if seen[id] {
			t.Fatalf("duplicate allocated ID: %s", id)
		}
		seen[id] = true
	}
	if len(seen) != count {
		t.Fatalf("expected %d allocated IDs, got %d", count, len(seen))
	}
}

func TestProposalCardSequenceExpandsBeyondThreeDigits(t *testing.T) {
	store, _ := setupTestStore(t)
	createTestDirs(t, store.wikiRoot)
	createTestProposal(t, store, "CR26081104", "Large sequence proposal")

	sequencePath := store.ProposalCardsDir("CR26081104") + "/.flowforge-card-sequence"
	if err := writeSequenceCounter(sequencePath, 1000); err != nil {
		t.Fatal(err)
	}
	id, err := store.NextCardID(CardTypeRequirement, "CR26081104")
	if err != nil {
		t.Fatal(err)
	}
	if id != "REQ-CR26081104-1000" {
		t.Fatalf("expected four-digit sequence after 999, got %s", id)
	}
}

func TestProposalCardSequenceRecognizesOnlyNumericSequenceParts(t *testing.T) {
	tests := []struct {
		id      string
		want    int
		matched bool
	}{
		{"FEAT-CR26081102-001", 1, true},
		{"TASK-CR26081102-i-1000", 1000, true},
		{"FEAT-CR26081102-dkm31krtv79s", 0, false},
		{"FEAT-OTHER-001", 0, false},
	}
	for _, tt := range tests {
		got, matched := proposalCardSequence(tt.id, "CR26081102")
		if got != tt.want || matched != tt.matched {
			t.Errorf("proposalCardSequence(%q) = (%d, %t), want (%d, %t)", tt.id, got, matched, tt.want, tt.matched)
		}
	}
}

func TestNextCardIDUsesLegacyFormatWithoutProposal(t *testing.T) {
	store, _ := setupTestStore(t)
	id, err := store.NextCardID(CardTypeConvention, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(id, "CONV-") || strings.Count(id, "-") != 1 {
		t.Fatalf("unexpected global card ID: %s", id)
	}
}
