package command

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestCardSearchScopesAndFilters(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	store := testCardStore(t, tmpDir)

	workspaceCard := core.NewCard(core.CardTypeDesign, "Keyword from workspace")
	workspaceCard.ID = "DES-work-1"
	workspaceCard.Status = core.CardStatusActive
	workspaceCard.Tags = []string{"searchable"}
	workspaceCard.Domain = "workspace-domain"
	workspaceCard.Body = "This body mentions UniqueKeyword only here."
	if _, err := store.CreateCard(workspaceCard, "CR26061301"); err != nil {
		t.Fatalf("creating workspace card failed: %v", err)
	}

	libraryCard := core.NewCard(core.CardTypeConvention, "Keyword from library")
	libraryCard.ID = "CONV-lib-1"
	libraryCard.Status = core.CardStatusActive
	libraryCard.Tags = []string{"library-tag"}
	libraryCard.Domain = "library-domain"
	libraryCard.Body = "Library body with UniqueKeyword and search text."
	if _, err := store.CreateCard(libraryCard, ""); err != nil {
		t.Fatalf("creating library card failed: %v", err)
	}

	otherCard := core.NewCard(core.CardTypeTask, "Different topic")
	otherCard.ID = "TASK-work-2"
	otherCard.Status = core.CardStatusReady
	otherCard.Body = "No match here."
	if _, err := store.CreateCard(otherCard, "CR26061301"); err != nil {
		t.Fatalf("creating other card failed: %v", err)
	}

	cmd := newCardSearchCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"UniqueKeyword", "--scope", "workspace", "--type", "design", "--limit", "5"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card search failed: %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "DES-work-1") {
		t.Fatalf("expected workspace design card in output:\n%s", text)
	}
	if strings.Contains(text, "CONV-lib-1") {
		t.Fatalf("expected library card to be excluded from workspace search:\n%s", text)
	}
	if !strings.Contains(text, "Match: matched body") {
		t.Fatalf("expected match reason in output:\n%s", text)
	}
	if strings.Contains(text, "UniqueKeyword only here") || strings.Contains(text, "Library body with UniqueKeyword") {
		t.Fatalf("expected search output to omit full body text:\n%s", text)
	}
}

func TestCardSearchAllScopeIncludesLibraryCards(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	store := testCardStore(t, tmpDir)

	card := core.NewCard(core.CardTypeConvention, "All scope match")
	card.ID = "CONV-all-1"
	card.Status = core.CardStatusActive
	card.Domain = "all-scope"
	card.Tags = []string{"alpha"}
	card.Body = "Body text with mixedcasequery."
	if _, err := store.CreateCard(card, ""); err != nil {
		t.Fatalf("creating card failed: %v", err)
	}

	cmd := newCardSearchCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"mixedcasequery", "--scope", "all"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card search failed: %v", err)
	}

	if !strings.Contains(out.String(), "CONV-all-1") {
		t.Fatalf("expected all-scope search to find library card:\n%s", out.String())
	}
}
