package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestCardCreateDefaultsToCurrentProposalForProposalScopedTypes(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Default proposal")

	cmd := newCardCreateCmd()
	cmd.SetArgs([]string{"--type", "requirement", "--title", "Default proposal requirement"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card create failed: %v", err)
	}

	store := testCardStore(t, tmpDir)
	cards, err := store.ListCards(store.ProposalCardsDir(proposalID))
	if err != nil {
		t.Fatalf("listing proposal cards failed: %v", err)
	}

	var requirement *core.Card
	for _, card := range cards {
		if card.Type == core.CardTypeRequirement && strings.Contains(card.Title, "Default proposal requirement") {
			requirement = card
			break
		}
	}
	if requirement == nil {
		t.Fatalf("expected requirement to be created under current proposal %s", proposalID)
	}
	if requirement.Source != proposalID {
		t.Fatalf("expected requirement source %s, got %q", proposalID, requirement.Source)
	}
	if !strings.Contains(requirement.ID, proposalID) {
		t.Fatalf("expected requirement ID %q to include proposal %s", requirement.ID, proposalID)
	}
}

func TestCardCreateKeepsLibraryTypesGlobalWithoutExplicitProposal(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Default proposal")

	cmd := newCardCreateCmd()
	cmd.SetArgs([]string{"--type", "convention", "--title", "Use explicit errors", "--links", "PROP-" + proposalID + ":references"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card create failed: %v", err)
	}

	store := testCardStore(t, tmpDir)
	proposalCards, err := store.ListCards(store.ProposalCardsDir(proposalID))
	if err != nil {
		t.Fatalf("listing proposal cards failed: %v", err)
	}
	for _, card := range proposalCards {
		if card.Type == core.CardTypeConvention && strings.Contains(card.Title, "Use explicit errors") {
			t.Fatalf("expected convention card to stay out of proposal cards")
		}
	}

	libraryCards, err := store.ListCards(filepath.Join(store.LibraryDir(), "60-conventions"))
	if err != nil {
		t.Fatalf("listing convention library failed: %v", err)
	}
	if len(libraryCards) != 1 {
		t.Fatalf("expected 1 convention library card, got %d", len(libraryCards))
	}
	if libraryCards[0].Source != "" {
		t.Fatalf("expected library convention to have empty source, got %q", libraryCards[0].Source)
	}
}

func TestCardCommandsUseCurrentProjectStore(t *testing.T) {
	tmpDir := createMultiProjectFixture(t)
	restoreWorkingDir(t)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	if err := useProjectForTest(t, "backend"); err != nil {
		t.Fatalf("project use failed: %v", err)
	}
	backendStore := core.NewCardStore(filepath.Join(tmpDir, "ff-wiki-be"))
	module := core.NewCard(core.CardTypeModule, "Backend module")
	module.ID = "MOD-BE"
	if _, err := backendStore.CreateCard(module, ""); err != nil {
		t.Fatalf("creating backend module failed: %v", err)
	}

	createCmd := newCardCreateCmd()
	createCmd.SetArgs([]string{"--type", "convention", "--title", "Backend only convention", "--links", "MOD-BE:references"})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("card create failed: %v", err)
	}

	listCmd := newCardListCmd()
	var out bytes.Buffer
	listCmd.SetOut(&out)
	listCmd.SetArgs([]string{"--type", "convention"})
	if err := listCmd.Execute(); err != nil {
		t.Fatalf("card list failed: %v", err)
	}
	if !strings.Contains(out.String(), "Backend only convention") {
		t.Fatalf("expected card list to use current backend project, got:\n%s", out.String())
	}
}
