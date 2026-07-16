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
	workspaceCard.Tags = []string{"searchable", "shared"}
	workspaceCard.Domain = "workspace-domain"
	workspaceCard.Body = "This body mentions UniqueKeyword only here."
	if _, err := store.CreateCard(workspaceCard, "CR26061301"); err != nil {
		t.Fatalf("creating workspace card failed: %v", err)
	}

	libraryCard := core.NewCard(core.CardTypeConvention, "Keyword from library")
	libraryCard.ID = "CONV-lib-1"
	libraryCard.Status = core.CardStatusActive
	libraryCard.Tags = []string{"library-tag", "shared"}
	libraryCard.Domain = "library-domain"
	libraryCard.Body = "Library body with UniqueKeyword and search text."
	if _, err := store.CreateCard(libraryCard, ""); err != nil {
		t.Fatalf("creating library card failed: %v", err)
	}

	otherCard := core.NewCard(core.CardTypeTask, "Different topic")
	otherCard.ID = "TASK-work-2"
	otherCard.Status = core.CardStatusReady
	otherCard.Tags = []string{"other-tag"}
	otherCard.Domain = "other-domain"
	otherCard.Body = "UniqueKeyword also appears here."
	if _, err := store.CreateCard(otherCard, "CR26061301"); err != nil {
		t.Fatalf("creating other card failed: %v", err)
	}

	cases := []struct {
		name       string
		args       []string
		wantIDs    []string
		wantMatch  string
		wantNotIDs []string
	}{
		{
			name:       "status filter",
			args:       []string{"UniqueKeyword", "--scope", "all", "--status", "active", "--limit", "5"},
			wantIDs:    []string{"CONV-lib-1", "DES-work-1"},
			wantMatch:  "Match: matched body | status=active",
			wantNotIDs: []string{"TASK-work-2"},
		},
		{
			name:       "domain filter",
			args:       []string{"UniqueKeyword", "--scope", "all", "--domain", "workspace-domain", "--limit", "5"},
			wantIDs:    []string{"DES-work-1"},
			wantMatch:  "Match: matched body | domain=workspace-domain",
			wantNotIDs: []string{"CONV-lib-1", "TASK-work-2"},
		},
		{
			name:       "tag filter",
			args:       []string{"UniqueKeyword", "--scope", "all", "--tag", "missing, searchable", "--limit", "5"},
			wantIDs:    []string{"DES-work-1"},
			wantMatch:  "Match: matched body | tag=searchable",
			wantNotIDs: []string{"CONV-lib-1", "TASK-work-2"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runCardSearchCommand(t, tc.args...)
			for _, id := range tc.wantIDs {
				if !strings.Contains(out, id) {
					t.Fatalf("expected output to include %s:\n%s", id, out)
				}
			}
			for _, id := range tc.wantNotIDs {
				if strings.Contains(out, id) {
					t.Fatalf("expected output to exclude %s:\n%s", id, out)
				}
			}
			if !strings.Contains(out, tc.wantMatch) {
				t.Fatalf("expected output to include match reason %q:\n%s", tc.wantMatch, out)
			}
			if strings.Contains(out, "UniqueKeyword only here.") || strings.Contains(out, "Library body with UniqueKeyword") {
				t.Fatalf("expected search output to omit full body text:\n%s", out)
			}
		})
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

func TestCardSearchInvalidTypeRejected(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")

	cmd := newCardSearchCmd()
	cmd.SetArgs([]string{"UniqueKeyword", "--type", "not-a-type"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected invalid type to fail")
	}
	if !strings.Contains(err.Error(), "invalid card type: not-a-type") {
		t.Fatalf("expected invalid type error, got: %v", err)
	}
}

func runCardSearchCommand(t *testing.T, args ...string) string {
	t.Helper()

	cmd := newCardSearchCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card search failed: %v", err)
	}
	return out.String()
}

func TestCardSearchProposalScope(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Search proposal")

	store := testCardStore(t, tmpDir)
	req := core.NewCard(core.CardTypeRequirement, "Searchable requirement")
	req.ID = "REQ-search-" + proposalID[len(proposalID)-2:]
	req.AddLink("PROP-"+proposalID, "belongs_to")
	req.Body = "UniqueProposalKeyword for scoped search test."
	if _, err := store.CreateCard(req, proposalID); err != nil {
		t.Fatalf("creating proposal card failed: %v", err)
	}

	libraryCard := core.NewCard(core.CardTypeDesign, "Library card with keyword")
	libraryCard.ID = "DES-lib-search"
	libraryCard.Body = "UniqueProposalKeyword also in library."
	if _, err := store.CreateCard(libraryCard, ""); err != nil {
		t.Fatalf("creating library card failed: %v", err)
	}

	out := runCardSearchCommand(t, "UniqueProposalKeyword", "--scope", "proposal", "--proposal", proposalID, "--limit", "5")
	if !strings.Contains(out, req.ID) {
		t.Fatalf("proposal scope search should find proposal card:\n%s", out)
	}
	if strings.Contains(out, libraryCard.ID) {
		t.Fatalf("proposal scope search should NOT find library card:\n%s", out)
	}

	outAll := runCardSearchCommand(t, "UniqueProposalKeyword", "--scope", "all", "--limit", "5")
	if !strings.Contains(outAll, req.ID) || !strings.Contains(outAll, libraryCard.ID) {
		t.Fatalf("all scope search should find both cards:\n%s", outAll)
	}
}

func TestCardSearchProposalScopeRequiresProposalFlag(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")

	cmd := newCardSearchCmd()
	cmd.SetArgs([]string{"keyword", "--scope", "proposal"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--proposal is required") {
		t.Fatalf("expected --proposal required error, got: %v", err)
	}
}
