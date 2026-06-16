package core

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func setupTestStore(t *testing.T) (*CardStore, string) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)
	return store, tmpDir
}

func createTestProposal(t *testing.T, store *CardStore, proposalID, title string) (string, string) {
	t.Helper()

	rootPath, indexPath, err := store.CreateProposal(proposalID, title)
	if err != nil {
		t.Fatalf("failed to create proposal: %v", err)
	}

	return rootPath, indexPath
}

func createTestDirs(t *testing.T, wikiRoot string) {
	dirs := []string{
		filepath.Join(wikiRoot, "01-workspace", "01-active"),
		filepath.Join(wikiRoot, "01-workspace", "02-intake"),
		filepath.Join(wikiRoot, "01-workspace", "03-completed"),
		filepath.Join(wikiRoot, "02-library", "10-requirements"),
		filepath.Join(wikiRoot, "02-library", "20-decisions"),
		filepath.Join(wikiRoot, "02-library", "30-designs"),
		filepath.Join(wikiRoot, "02-library", "40-tasks"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
}

func cardHasLink(card *Card, target string, relation string) bool {
	for _, link := range card.Links {
		if link.Target == target && link.Relation == relation {
			return true
		}
	}
	return false
}

func TestNewCardStore(t *testing.T) {
	store := NewCardStore("/test/wiki")
	if store.wikiRoot != "/test/wiki" {
		t.Errorf("expected wikiRoot /test/wiki, got %s", store.wikiRoot)
	}
}

func TestCardStoreDirectories(t *testing.T) {
	store := NewCardStore("/test/wiki")

	if store.WorkspaceDir() != "/test/wiki/01-workspace" {
		t.Errorf("unexpected WorkspaceDir: %s", store.WorkspaceDir())
	}

	if store.ActiveDir() != "/test/wiki/01-workspace/01-active" {
		t.Errorf("unexpected ActiveDir: %s", store.ActiveDir())
	}

	if store.IntakeDir() != "/test/wiki/01-workspace/02-intake" {
		t.Errorf("unexpected IntakeDir: %s", store.IntakeDir())
	}

	if store.CompletedDir() != "/test/wiki/01-workspace/03-completed" {
		t.Errorf("unexpected CompletedDir: %s", store.CompletedDir())
	}

	if store.LibraryDir() != "/test/wiki/02-library" {
		t.Errorf("unexpected LibraryDir: %s", store.LibraryDir())
	}
}

func TestCardStoreLibraryTypeDir(t *testing.T) {
	store := NewCardStore("/test/wiki")

	tests := []struct {
		cardType CardType
		expected string
	}{
		{CardTypeRequirement, "/test/wiki/02-library/10-requirements"},
		{CardTypeDecision, "/test/wiki/02-library/20-decisions"},
		{CardTypeDesign, "/test/wiki/02-library/30-designs"},
		{CardTypeTask, "/test/wiki/02-library/40-tasks"},
		{CardTypeLog, "/test/wiki/02-library/50-logs"},
		{CardTypeConvention, "/test/wiki/02-library/60-conventions"},
		{CardTypeFinding, "/test/wiki/02-library/70-findings"},
		{CardTypeModule, "/test/wiki/02-library/80-modules"},
		{CardTypeStructure, "/test/wiki/02-library/structures"},
	}

	for _, tt := range tests {
		result := store.LibraryTypeDir(tt.cardType)
		if result != tt.expected {
			t.Errorf("LibraryTypeDir(%s) = %s, expected %s", tt.cardType, result, tt.expected)
		}
	}
}

func TestCreateProposal(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	rootPath, indexPath := createTestProposal(t, store, "CR260612", "Test Proposal")
	proposalDir := store.ProposalDir("CR260612")

	if _, err := os.Stat(proposalDir); err != nil {
		t.Errorf("proposal directory not created: %v", err)
	}

	cardsDir := filepath.Join(proposalDir, "90-cards")
	if _, err := os.Stat(cardsDir); err != nil {
		t.Errorf("cards directory not created: %v", err)
	}

	expectedRootPath := filepath.Join(store.ProposalCardDir(), "CR260612_test-proposal.md")
	if rootPath != expectedRootPath {
		t.Fatalf("expected rootPath %s, got %s", expectedRootPath, rootPath)
	}
	if _, err := os.Stat(rootPath); err != nil {
		t.Errorf("root card not created: %v", err)
	}
	rootCard, err := ParseCardFile(rootPath)
	if err != nil {
		t.Fatalf("parsing root card failed: %v", err)
	}
	if !cardHasLink(rootCard, "STR-CR260612-REQ", "indexes") {
		t.Fatalf("root card missing indexes link to requirement index: %#v", rootCard.Links)
	}
	for _, want := range []string{
		"## Purpose",
		"## Entries",
		"[STR-CR260612-REQ](STR-CR260612-REQ.md) (structure, active) - Requirement index",
	} {
		if !strings.Contains(rootCard.Body, want) {
			t.Fatalf("root body missing %q:\n%s", want, rootCard.Body)
		}
	}

	expectedIndexPath := filepath.Join(proposalDir, "STR-CR260612-REQ.md")
	if indexPath != expectedIndexPath {
		t.Fatalf("expected indexPath %s, got %s", expectedIndexPath, indexPath)
	}
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("requirement index card not created: %v", err)
	}
	indexCard, err := ParseCardFile(indexPath)
	if err != nil {
		t.Fatalf("parsing index card failed: %v", err)
	}
	if !cardHasLink(indexCard, "PROP-CR260612", "belongs_to") {
		t.Fatalf("requirement index missing belongs_to link to root: %#v", indexCard.Links)
	}
	for _, want := range []string{
		"## Purpose",
		"Top-level requirement index for Test Proposal.",
		"## Entries\n\n- None",
		"## Open Questions\n\n- None",
	} {
		if !strings.Contains(indexCard.Body, want) {
			t.Fatalf("index body missing %q:\n%s", want, indexCard.Body)
		}
	}

	_, _, err = store.CreateProposal("CR260612", "Test Proposal")
	if err == nil {
		t.Error("expected error for duplicate proposal")
	}
}

func TestReadCardFindsProposalRootAndIndexFiles(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)

	if _, _, err := store.CreateProposal("CR260612", "Test Proposal"); err != nil {
		t.Fatalf("CreateProposal failed: %v", err)
	}

	root, err := store.ReadCard("PROP-CR260612")
	if err != nil {
		t.Fatalf("ReadCard root failed: %v", err)
	}
	if root.Type != CardTypeProposal {
		t.Fatalf("expected root type proposal, got %s", root.Type)
	}

	index, err := store.ReadCard("STR-CR260612-REQ")
	if err != nil {
		t.Fatalf("ReadCard requirement index failed: %v", err)
	}
	if index.Type != CardTypeStructure {
		t.Fatalf("expected index type structure, got %s", index.Type)
	}
}

func TestCreateCard(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	card := NewCard(CardTypeRequirement, "Test Requirement")
	card.ID = "REQ-abc-123"

	filePath, err := store.CreateCard(card, "CR260612")
	if err != nil {
		t.Fatalf("failed to create card: %v", err)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Errorf("card file not created: %v", err)
	}

	if card.Source != "CR260612" {
		t.Errorf("expected source CR260612, got %s", card.Source)
	}
}

func TestReadCard(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	card := NewCard(CardTypeDecision, "Test Decision")
	card.ID = "DEC-abc-456"
	card.Body = "Decision body"

	filePath, _ := store.CreateCard(card, "CR260612")

	loaded, err := store.ReadCard("DEC-abc-456")
	if err != nil {
		t.Fatalf("failed to read card: %v", err)
	}

	if loaded.ID != card.ID {
		t.Errorf("expected ID %s, got %s", card.ID, loaded.ID)
	}

	if loaded.Title != card.Title {
		t.Errorf("expected title %s, got %s", card.Title, loaded.Title)
	}

	if loaded.FilePath != filePath {
		t.Errorf("expected filePath %s, got %s", filePath, loaded.FilePath)
	}
}

func TestReadCardNotFound(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	_, err := store.ReadCard("NONEXISTENT")
	if err == nil {
		t.Error("expected error for non-existent card")
	}
}

func TestUpdateCard(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	card := NewCard(CardTypeRequirement, "Original Title")
	card.ID = "REQ-abc-789"
	filePath, err := store.CreateCard(card, "CR260612")
	if err != nil {
		t.Fatalf("failed to create card: %v", err)
	}

	loaded, _ := store.ReadCard("REQ-abc-789")
	loaded.Title = "Updated Title"

	if err := store.UpdateCard(loaded); err != nil {
		t.Fatalf("failed to update card: %v", err)
	}

	reloaded, _ := store.ReadCard("REQ-abc-789")
	if reloaded.Title != "Updated Title" {
		t.Errorf("expected updated title, got %s", reloaded.Title)
	}
	if reloaded.FilePath != filePath {
		t.Fatalf("expected filePath to remain stable, got %s want %s", reloaded.FilePath, filePath)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected original file path to remain after update: %v", err)
	}
}

func TestUpdateCardWithLockPreservesConcurrentLinkAdds(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	structure := NewCard(CardTypeStructure, "Concurrent structure")
	structure.ID = "STR-atomic"
	if _, err := store.CreateCard(structure, "CR260612"); err != nil {
		t.Fatalf("failed to create structure card: %v", err)
	}

	linkIDs := []string{"REQ-atomic-1", "REQ-atomic-2", "REQ-atomic-3"}
	for i, linkID := range linkIDs {
		req := NewCard(CardTypeRequirement, "Requirement")
		req.ID = linkID
		req.Body = filepath.Base(linkID) + " body"
		if _, err := store.CreateCard(req, "CR260612"); err != nil {
			t.Fatalf("failed to create linked card %d: %v", i, err)
		}
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	errs := make(chan error, len(linkIDs))
	for _, linkID := range linkIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			<-start
			errs <- store.UpdateCardWithLock("STR-atomic", func(card *Card) error {
				card.AddLink(id, "indexes")
				time.Sleep(20 * time.Millisecond)
				return nil
			})
		}(linkID)
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent update failed: %v", err)
		}
	}

	reloaded, err := store.ReadCard("STR-atomic")
	if err != nil {
		t.Fatalf("failed to reload structure card: %v", err)
	}
	if len(reloaded.Links) != len(linkIDs) {
		t.Fatalf("expected %d links after concurrent updates, got %d", len(linkIDs), len(reloaded.Links))
	}

	found := map[string]bool{}
	for _, link := range reloaded.Links {
		if link.Relation != "indexes" {
			t.Fatalf("unexpected relation on concurrent link: %#v", link)
		}
		found[link.Target] = true
	}
	for _, linkID := range linkIDs {
		if !found[linkID] {
			t.Fatalf("missing concurrent link %s", linkID)
		}
	}
}

func TestCardLockReleaseDoesNotRemoveForeignLock(t *testing.T) {
	lockPath := filepath.Join(t.TempDir(), "card.lock")

	release, err := acquireCardLock(lockPath)
	if err != nil {
		t.Fatalf("acquireCardLock failed: %v", err)
	}
	if err := os.WriteFile(lockPath, []byte("foreign-owner"), 0600); err != nil {
		t.Fatalf("replacing lock owner failed: %v", err)
	}
	if err := release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	data, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("expected foreign lock to remain: %v", err)
	}
	if string(data) != "foreign-owner" {
		t.Fatalf("expected foreign lock owner to remain, got %q", string(data))
	}
}

func TestDeleteCard(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	card := NewCard(CardTypeRequirement, "To Delete")
	card.ID = "REQ-abc-del"
	filePath, _ := store.CreateCard(card, "CR260612")

	if err := store.DeleteCard("REQ-abc-del"); err != nil {
		t.Fatalf("failed to delete card: %v", err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected card file to be deleted")
	}
}

func TestDeleteCardNonDraft(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	card := NewCard(CardTypeRequirement, "Active Card")
	card.ID = "REQ-abc-act"
	card.Status = CardStatusActive
	store.CreateCard(card, "CR260612")

	err := store.DeleteCard("REQ-abc-act")
	if err == nil {
		t.Error("expected error when deleting non-draft card")
	}
}

func TestListCards(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	for i := 0; i < 3; i++ {
		card := NewCard(CardTypeRequirement, "Req "+string(rune('A'+i)))
		card.ID = "REQ-abc-" + string(rune('0'+i))
		store.CreateCard(card, "CR260612")
	}

	cardsDir := store.ProposalCardsDir("CR260612")
	cards, err := store.ListCards(cardsDir)
	if err != nil {
		t.Fatalf("failed to list cards: %v", err)
	}

	if len(cards) != 3 {
		t.Errorf("expected 3 cards, got %d", len(cards))
	}
}

func TestListCardsByType(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	req1 := NewCard(CardTypeRequirement, "Req 1")
	req1.ID = "REQ-abc-r1"
	store.CreateCard(req1, "CR260612")

	dec1 := NewCard(CardTypeDecision, "Dec 1")
	dec1.ID = "DEC-abc-d1"
	store.CreateCard(dec1, "CR260612")

	req2 := NewCard(CardTypeRequirement, "Req 2")
	req2.ID = "REQ-abc-r2"
	store.CreateCard(req2, "CR260612")

	reqCards, _ := store.ListCardsByType(CardTypeRequirement)
	if len(reqCards) != 2 {
		t.Errorf("expected 2 requirements, got %d", len(reqCards))
	}

	decCards, _ := store.ListCardsByType(CardTypeDecision)
	if len(decCards) != 1 {
		t.Errorf("expected 1 decision, got %d", len(decCards))
	}
}

func TestGetDependents(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	req := NewCard(CardTypeRequirement, "Requirement")
	req.ID = "REQ-abc-main"
	store.CreateCard(req, "CR260612")

	dec := NewCard(CardTypeDecision, "Decision")
	dec.ID = "DEC-abc-dep"
	dec.AddLink("REQ-abc-main", "references")
	store.CreateCard(dec, "CR260612")

	dependents, err := store.GetDependents("REQ-abc-main")
	if err != nil {
		t.Fatalf("failed to get dependents: %v", err)
	}

	if len(dependents) != 1 {
		t.Errorf("expected 1 dependent, got %d", len(dependents))
	}

	if dependents[0].ID != "DEC-abc-dep" {
		t.Errorf("expected dependent DEC-abc-dep, got %s", dependents[0].ID)
	}
}

func TestGetRelated(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	req := NewCard(CardTypeRequirement, "Requirement")
	req.ID = "REQ-abc-rel"
	store.CreateCard(req, "CR260612")

	dec := NewCard(CardTypeDecision, "Decision")
	dec.ID = "DEC-abc-rel"
	dec.AddLink("REQ-abc-rel", "references")
	store.CreateCard(dec, "CR260612")

	related, err := store.GetRelated("DEC-abc-rel", "", 1)
	if err != nil {
		t.Fatalf("failed to get related: %v", err)
	}

	if len(related) != 1 {
		t.Errorf("expected 1 related card, got %d", len(related))
	}

	if related[0].ID != "REQ-abc-rel" {
		t.Errorf("expected related REQ-abc-rel, got %s", related[0].ID)
	}
}

func TestGetRelatedWithFilter(t *testing.T) {
	store, wikiRoot := setupTestStore(t)
	createTestDirs(t, wikiRoot)

	createTestProposal(t, store, "CR260612", "Test")

	req := NewCard(CardTypeRequirement, "Requirement")
	req.ID = "REQ-abc-fil"
	store.CreateCard(req, "CR260612")

	dec := NewCard(CardTypeDecision, "Decision")
	dec.ID = "DEC-abc-fil"
	dec.AddLink("REQ-abc-fil", "references")
	store.CreateCard(dec, "CR260612")

	related, _ := store.GetRelated("DEC-abc-fil", "implements", 1)
	if len(related) != 0 {
		t.Errorf("expected 0 related cards with implements filter, got %d", len(related))
	}

	related, _ = store.GetRelated("DEC-abc-fil", "references", 1)
	if len(related) != 1 {
		t.Errorf("expected 1 related card with references filter, got %d", len(related))
	}
}
