package command

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestLibrarySuggestRanksAndFiltersCards(t *testing.T) {
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

	focus := core.NewCard(core.CardTypeDesign, "Library search flow")
	focus.ID = "DES-focus-1"
	focus.Status = core.CardStatusActive
	focus.Tags = []string{"indexing", "search"}
	focus.Domain = "library"
	focus.Body = "Need a search index for library cards and matching conventions."
	if _, err := store.CreateCard(focus, "CR26061301"); err != nil {
		t.Fatalf("creating focus card failed: %v", err)
	}

	mustConvention := core.NewCard(core.CardTypeConvention, "Library search flow")
	mustConvention.ID = "CONV-lib-1"
	mustConvention.Status = core.CardStatusActive
	mustConvention.Importance = core.ImportanceMust
	mustConvention.Domain = "library"
	mustConvention.Tags = []string{"search"}
	mustConvention.Body = "Shared rule about search flow."
	if _, err := store.CreateCard(mustConvention, ""); err != nil {
		t.Fatalf("creating convention failed: %v", err)
	}

	moduleCard := core.NewCard(core.CardTypeModule, "Search indexing module")
	moduleCard.ID = "MOD-lib-1"
	moduleCard.Status = core.CardStatusActive
	moduleCard.Domain = "library"
	moduleCard.Body = "Index library entries for lookup."
	if _, err := store.CreateCard(moduleCard, ""); err != nil {
		t.Fatalf("creating module failed: %v", err)
	}

	bodyOnly := core.NewCard(core.CardTypeFinding, "Unrelated note")
	bodyOnly.ID = "FIND-lib-1"
	bodyOnly.Status = core.CardStatusActive
	bodyOnly.Body = "This finding mentions search flow in the body only."
	if _, err := store.CreateCard(bodyOnly, ""); err != nil {
		t.Fatalf("creating finding failed: %v", err)
	}

	deprecated := core.NewCard(core.CardTypeDesign, "Deprecated design")
	deprecated.ID = "DES-lib-2"
	deprecated.Status = core.CardStatusDeprecated
	deprecated.Body = "search flow"
	if _, err := store.CreateCard(deprecated, ""); err != nil {
		t.Fatalf("creating deprecated card failed: %v", err)
	}

	cmd := newLibrarySuggestCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--for", "DES-focus-1", "--relation", "constrains", "--limit", "3"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("library suggest failed: %v", err)
	}

	text := out.String()
	for _, want := range []string{
		"## Library Suggestions",
		"| ID | Type | Title | Status | Importance | Domain | Score | SuggestedRelation |",
		"CONV-lib-1",
		"MOD-lib-1",
		"constrains",
		"## Match Reasons",
		"## Recommended Reads",
		"flowforge card read CONV-lib-1 --summary",
		"## Not Included",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("library suggest output missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "DES-lib-2") {
		t.Fatalf("expected deprecated card to be omitted:\n%s", text)
	}
	if strings.Contains(text, "This finding mentions search flow in the body only.") {
		t.Fatalf("expected output to omit full card body:\n%s", text)
	}

	conventionPos := strings.Index(text, "CONV-lib-1")
	modulePos := strings.Index(text, "MOD-lib-1")
	if conventionPos < 0 || modulePos < 0 || conventionPos > modulePos {
		t.Fatalf("expected must convention to rank before module:\n%s", text)
	}
}

func TestLibrarySuggestRequiresFocusCard(t *testing.T) {
	cmd := newLibrarySuggestCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected missing --for to fail")
	}
}
