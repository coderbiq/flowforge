package core

import (
	"os"
	"path/filepath"
	"testing"
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

	expectedRootPath := filepath.Join(proposalDir, "ROOT-CR260612.md")
	if rootPath != expectedRootPath {
		t.Fatalf("expected rootPath %s, got %s", expectedRootPath, rootPath)
	}
	if _, err := os.Stat(rootPath); err != nil {
		t.Errorf("root card not created: %v", err)
	}

	expectedIndexPath := filepath.Join(proposalDir, "STR-CR260612-REQ.md")
	if indexPath != expectedIndexPath {
		t.Fatalf("expected indexPath %s, got %s", expectedIndexPath, indexPath)
	}
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("requirement index card not created: %v", err)
	}

	_, _, err := store.CreateProposal("CR260612", "Duplicate")
	if err == nil {
		t.Error("expected error for duplicate proposal")
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
	store.CreateCard(card, "CR260612")

	loaded, _ := store.ReadCard("REQ-abc-789")
	loaded.Title = "Updated Title"

	if err := store.UpdateCard(loaded); err != nil {
		t.Fatalf("failed to update card: %v", err)
	}

	reloaded, _ := store.ReadCard("REQ-abc-789")
	if reloaded.Title != "Updated Title" {
		t.Errorf("expected updated title, got %s", reloaded.Title)
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
