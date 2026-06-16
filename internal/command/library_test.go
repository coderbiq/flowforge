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

func TestLibraryFacetsClassifyAndSuggestByFacet(t *testing.T) {
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

	serviceRule := core.NewCard(core.CardTypeConvention, "Service page query rule")
	serviceRule.ID = "CONV-service-page-query"
	serviceRule.Status = core.CardStatusActive
	serviceRule.Importance = core.ImportanceMust
	serviceRule.Tags = []string{"layer:service", "scenario:page-query"}
	serviceRule.Body = "Service page-query implementations validate filters before repository calls."
	if _, err := store.CreateCard(serviceRule, ""); err != nil {
		t.Fatalf("creating service rule failed: %v", err)
	}

	controllerRule := core.NewCard(core.CardTypeConvention, "Controller create rule")
	controllerRule.ID = "CONV-controller-create"
	controllerRule.Status = core.CardStatusActive
	controllerRule.Tags = []string{"layer:controller", "scenario:create"}
	controllerRule.Body = "Controller create endpoints validate request payloads."
	if _, err := store.CreateCard(controllerRule, ""); err != nil {
		t.Fatalf("creating controller rule failed: %v", err)
	}

	task := core.NewCard(core.CardTypeTask, "Implement customer service pagination")
	task.ID = "TASK-focus-service-page"
	task.Status = core.CardStatusReady
	task.Tags = []string{"layer:service"}
	task.Body = "Implement customer page-query behavior in the service layer."
	if _, err := store.CreateCard(task, "CR26061401"); err != nil {
		t.Fatalf("creating task failed: %v", err)
	}

	facetsCmd := newLibraryFacetsCmd()
	var facetsOut bytes.Buffer
	facetsCmd.SetOut(&facetsOut)
	if err := facetsCmd.Execute(); err != nil {
		t.Fatalf("library facets failed: %v", err)
	}
	facetsText := facetsOut.String()
	for _, want := range []string{
		"## Library Facets",
		"| layer | service | 1 |",
		"| scenario | page-query | 1 |",
		"layer:service + scenario:page-query",
	} {
		if !strings.Contains(facetsText, want) {
			t.Fatalf("library facets output missing %q:\n%s", want, facetsText)
		}
	}

	classifyCmd := newLibraryClassifyCmd()
	var classifyOut bytes.Buffer
	classifyCmd.SetOut(&classifyOut)
	classifyCmd.SetArgs([]string{"--for", task.ID})
	if err := classifyCmd.Execute(); err != nil {
		t.Fatalf("library classify failed: %v", err)
	}
	classifyText := classifyOut.String()
	for _, want := range []string{
		"## Library Classification",
		"| layer:service | tag | layer:service | 1 |",
		"| scenario:page-query | text | page-query | 1 |",
		"flowforge library suggest --for TASK-focus-service-page",
		"--facet layer:service",
	} {
		if !strings.Contains(classifyText, want) {
			t.Fatalf("library classify output missing %q:\n%s", want, classifyText)
		}
	}

	suggestCmd := newLibrarySuggestCmd()
	var suggestOut bytes.Buffer
	suggestCmd.SetOut(&suggestOut)
	suggestCmd.SetArgs([]string{
		"--for", task.ID,
		"--facet", "layer:service",
		"--facet", "scenario:page-query",
		"--types", "convention",
	})
	if err := suggestCmd.Execute(); err != nil {
		t.Fatalf("library suggest failed: %v", err)
	}
	suggestText := suggestOut.String()
	for _, want := range []string{
		"CONV-service-page-query",
		"facets:layer:service,scenario:page-query",
		"constrains",
	} {
		if !strings.Contains(suggestText, want) {
			t.Fatalf("library suggest output missing %q:\n%s", want, suggestText)
		}
	}
	if strings.Contains(suggestText, "CONV-controller-create") {
		t.Fatalf("facet-filtered suggestion should omit controller rule:\n%s", suggestText)
	}
}

func TestLibraryImportCreatesValidatedLibraryCard(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Library import proposal")

	store := testCardStore(t, tmpDir)
	source := core.NewCard(core.CardTypeFinding, "Source finding")
	source.ID = "FIND-import-source"
	source.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(source, proposalID); err != nil {
		t.Fatalf("creating source finding failed: %v", err)
	}

	cmd := newLibraryImportCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{
		"--type", "convention",
		"--title", "Imported service rule",
		"--body", "## Rule\n\nValidate service filters before repository calls.",
		"--source-card", source.ID,
		"--tags", "layer:service,scenario:page-query",
		"--domain", "service",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("library import failed: %v", err)
	}
	if !strings.Contains(out.String(), "✓ Imported library card CONV-") {
		t.Fatalf("unexpected import output:\n%s", out.String())
	}

	cards, err := store.ListCards(store.LibraryTypeDir(core.CardTypeConvention))
	if err != nil {
		t.Fatalf("listing convention library failed: %v", err)
	}
	if len(cards) != 1 {
		t.Fatalf("expected 1 imported convention, got %d", len(cards))
	}
	if !hasLinkRelation(cards[0], source.ID, "references") {
		t.Fatalf("expected imported card to reference source card, got %#v", cards[0].Links)
	}

	validateCmd := newValidateAllCmd()
	if err := validateCmd.Execute(); err != nil {
		t.Fatalf("validate all failed after library import: %v", err)
	}
}

func TestLibraryPromoteCopiesProposalCardToLibrary(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Library promote proposal")

	store := testCardStore(t, tmpDir)
	finding := core.NewCard(core.CardTypeFinding, "Promotable finding")
	finding.ID = "FIND-promote-source"
	finding.Body = "## Finding\n\nStable reusable behavior."
	finding.Tags = []string{"behavior"}
	finding.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(finding, proposalID); err != nil {
		t.Fatalf("creating finding failed: %v", err)
	}

	cmd := newLibraryPromoteCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{finding.ID, "--type", "convention", "--title", "Promoted convention"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("library promote failed: %v", err)
	}
	if !strings.Contains(out.String(), "✓ Promoted FIND-promote-source to library card CONV-") {
		t.Fatalf("unexpected promote output:\n%s", out.String())
	}

	cards, err := store.ListCards(store.LibraryTypeDir(core.CardTypeConvention))
	if err != nil {
		t.Fatalf("listing convention library failed: %v", err)
	}
	if len(cards) != 1 {
		t.Fatalf("expected 1 promoted convention, got %d", len(cards))
	}
	promoted := cards[0]
	if promoted.Title != "Promoted convention" {
		t.Fatalf("expected overridden title, got %q", promoted.Title)
	}
	if promoted.Source != finding.ID {
		t.Fatalf("expected source %s, got %q", finding.ID, promoted.Source)
	}
	if !hasLinkRelation(promoted, finding.ID, "references") {
		t.Fatalf("expected promoted card to reference source card, got %#v", promoted.Links)
	}

	validateCmd := newValidateAllCmd()
	if err := validateCmd.Execute(); err != nil {
		t.Fatalf("validate all failed after library promote: %v", err)
	}
}

func TestLibrarySuggestRequiresFocusCard(t *testing.T) {
	cmd := newLibrarySuggestCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected missing --for to fail")
	}
}
