package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateCard(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Title:      "Test Card",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	result := ValidateCard(card)
	if result.HasErrors() {
		t.Errorf("expected valid card, got errors: %s", result.String())
	}
}

func TestValidateCardMissingID(t *testing.T) {
	card := &Card{
		Title:      "Test Card",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected error for missing ID")
	}

	found := false
	for _, e := range result.Errors {
		if e.Field == "id" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error for id field")
	}
}

func TestValidateCardMissingTitle(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected error for missing title")
	}
}

func TestValidateCardInvalidType(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Title:      "Test",
		Type:       CardType("invalid"),
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected error for invalid type")
	}

	found := false
	for _, e := range result.Errors {
		if e.Field == "type" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error for type field")
	}
}

func TestValidateCardInvalidStatus(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Title:      "Test",
		Type:       CardTypeRequirement,
		Status:     CardStatus("invalid"),
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected error for invalid status")
	}
}

func TestValidateCardInvalidImportance(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Title:      "Test",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: Importance("invalid"),
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected error for invalid importance")
	}
}

func TestValidateCardMissingTimestamps(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Title:      "Test",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected error for missing timestamps")
	}

	createdFound := false
	updatedFound := false
	for _, e := range result.Errors {
		if e.Field == "created" {
			createdFound = true
		}
		if e.Field == "updated" {
			updatedFound = true
		}
	}
	if !createdFound || !updatedFound {
		t.Error("expected errors for created and updated fields")
	}
}

func TestValidateCardIDPrefixMismatch(t *testing.T) {
	card := &Card{
		ID:         "DEC-abc-123",
		Title:      "Test",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected error for ID prefix mismatch")
	}

	found := false
	for _, e := range result.Errors {
		if e.Field == "id" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error for id field")
	}
}

func TestValidateCardInvalidLink(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Title:      "Test",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Links: []Link{
			{Target: "", Relation: "references"},
			{Target: "DEC-456", Relation: ""},
			{Target: "DEC-789", Relation: "invalid-relation"},
		},
		Created: time.Now(),
		Updated: time.Now(),
	}

	result := ValidateCard(card)
	if !result.HasErrors() {
		t.Error("expected errors for invalid links")
	}

	if len(result.Errors) < 3 {
		t.Errorf("expected at least 3 link errors, got %d", len(result.Errors))
	}
}

func TestValidateCardValidLinks(t *testing.T) {
	card := &Card{
		ID:         "REQ-abc-123",
		Title:      "Test",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Links: []Link{
			{Target: "DEC-456", Relation: "references"},
			{Target: "DES-789", Relation: "implements"},
			{Target: "TASK-abc", Relation: "blocks"},
		},
		Created: time.Now(),
		Updated: time.Now(),
	}

	result := ValidateCard(card)
	if result.HasErrors() {
		t.Errorf("expected valid links, got errors: %s", result.String())
	}
}

func TestValidateCardComplexAnalysisFeatureRequiresRecoverableSections(t *testing.T) {
	card := &Card{
		ID:         "FEAT-abc-analysis",
		Title:      "Complex analysis",
		Type:       CardTypeFeature,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body: `<!-- analysis-mode: complex -->

## Objective

Preserve the analysis state.

## Current Understanding

TBD
`,
	}

	result := ValidateCard(card)
	assertValidationField(t, result, "body.analysis.current-understanding")
	assertValidationField(t, result, "body.analysis.evidence")
	assertValidationField(t, result, "body.analysis.working-design")
}

func TestValidateCardSimpleFeatureDoesNotRequireAnalysisSections(t *testing.T) {
	card := &Card{
		ID:         "FEAT-abc-simple",
		Title:      "Simple feature",
		Type:       CardTypeFeature,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body:       "## Summary\n\nSmall, self-contained change.",
	}

	result := ValidateCard(card)
	if result.HasErrors() {
		t.Fatalf("simple feature should remain backward compatible: %s", result.String())
	}
}

func TestValidateCardComplexAnalysisFindingRequiresWorkIDAndSource(t *testing.T) {
	card := &Card{
		ID:         "FIND-abc-analysis",
		Title:      "Investigation evidence",
		Type:       CardTypeFinding,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body: `<!-- analysis-mode: complex -->

## Evidence

Observed behavior.

## Source

TBD

## Impact

Changes the working design.

## Open Questions

None
`,
	}

	result := ValidateCard(card)
	assertValidationField(t, result, "body.analysis-work-id")
	assertValidationField(t, result, "body.analysis.source")
}

func TestValidateCardComplexAnalysisFindingAcceptsStableWorkID(t *testing.T) {
	card := &Card{
		ID:         "FIND-abc-valid",
		Title:      "Investigation evidence",
		Type:       CardTypeFinding,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body: `<!-- analysis-mode: complex -->
<!-- analysis-work-id: r2.cache-recovery -->

## Evidence

Observed behavior.

## Source

internal/core/store.go:42

## Impact

Changes the working design.

## Open Questions

None
`,
	}

	result := ValidateCard(card)
	if result.HasErrors() {
		t.Fatalf("expected recoverable finding to validate: %s", result.String())
	}
}

func TestValidateCardFileInStoreRejectsMissingAnalysisFindingReference(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)
	if err := os.MkdirAll(store.ProposalCardsDir("CR260811"), 0755); err != nil {
		t.Fatalf("creating proposal cards dir failed: %v", err)
	}

	card := &Card{
		ID:         "FEAT-abc-reference",
		Title:      "Complex analysis reference",
		Type:       CardTypeFeature,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body: `<!-- analysis-mode: complex -->

## Objective
Recover analysis.
## Current Understanding
Current facts.
## Evidence
Accepted FIND-missing-evidence.
## Working Design
Current design.
## Rejected or Revised Assumptions
None
## Open Questions
None
## Next Investigation
None
`,
	}
	filePath := filepath.Join(store.ProposalCardsDir("CR260811"), GenerateFilename(card.ID, card.Title))
	if err := card.Save(filePath); err != nil {
		t.Fatalf("saving feature failed: %v", err)
	}

	result := ValidateCardFileInStore(filePath, store)
	assertValidationField(t, result, "body.analysis.evidence")
}

func assertValidationField(t *testing.T, result *ValidationResult, field string) {
	t.Helper()
	for _, validationErr := range result.Errors {
		if validationErr.Field == field {
			return
		}
	}
	t.Fatalf("expected validation field %q, got %s", field, result.String())
}

func TestValidateCardFile(t *testing.T) {
	tmpDir := t.TempDir()

	card := &Card{
		ID:         "REQ-abc-file",
		Title:      "Test Card",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body:       "Card body content",
	}

	filename := GenerateFilename(card.ID, card.Title)
	filePath := filepath.Join(tmpDir, filename)

	if err := card.Save(filePath); err != nil {
		t.Fatalf("failed to save card: %v", err)
	}

	result := ValidateCardFile(filePath)
	if result.HasErrors() {
		t.Errorf("expected valid card file, got errors: %s", result.String())
	}
}

func TestValidateCardFileFilenameMismatch(t *testing.T) {
	tmpDir := t.TempDir()

	card := &Card{
		ID:         "REQ-abc-mis",
		Title:      "Test Card",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	wrongFilename := "REQ-different_wrong-name.md"
	filePath := filepath.Join(tmpDir, wrongFilename)

	if err := card.Save(filePath); err != nil {
		t.Fatalf("failed to save card: %v", err)
	}

	result := ValidateCardFile(filePath)
	if !result.HasErrors() {
		t.Error("expected error for filename mismatch")
	}

	found := false
	for _, e := range result.Errors {
		if e.Field == "filename" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected error for filename field")
	}
}

func TestValidateCardFileParseError(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.md")

	if err := os.WriteFile(filePath, []byte("no frontmatter"), 0644); err != nil {
		t.Fatal(err)
	}

	result := ValidateCardFile(filePath)
	if !result.HasErrors() {
		t.Error("expected error for parse failure")
	}
}

func TestValidateCardFileInStoreRejectsMissingLinkTargets(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)

	if err := os.MkdirAll(store.ProposalCardsDir("CR260612"), 0755); err != nil {
		t.Fatalf("creating proposal cards dir failed: %v", err)
	}

	card := &Card{
		ID:         "DES-abc-link",
		Title:      "Broken link card",
		Type:       CardTypeDesign,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Links: []Link{
			{Target: "REQ-MISSING", Relation: "references"},
		},
		Created: time.Now(),
		Updated: time.Now(),
	}

	filePath := filepath.Join(store.ProposalCardsDir("CR260612"), GenerateFilename(card.ID, card.Title))
	if err := card.Save(filePath); err != nil {
		t.Fatalf("saving card failed: %v", err)
	}

	result := ValidateCardFileInStore(filePath, store)
	if !result.HasErrors() {
		t.Fatal("expected missing target validation error")
	}
}

func TestValidateCardFileInStoreRejectsMissingMarkdownLinkTargets(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)

	if err := os.MkdirAll(store.ProposalCardsDir("CR260612"), 0755); err != nil {
		t.Fatalf("creating proposal cards dir failed: %v", err)
	}

	card := &Card{
		ID:         "STR-abc-req",
		Title:      "Broken structure",
		Type:       CardTypeStructure,
		Status:     CardStatusActive,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body:       "## Entries\n\n- [REQ-MISSING](REQ-MISSING.md) - missing",
	}

	filePath := filepath.Join(store.ProposalCardsDir("CR260612"), GenerateFilename(card.ID, card.Title))
	if err := card.Save(filePath); err != nil {
		t.Fatalf("saving card failed: %v", err)
	}

	result := ValidateCardFileInStore(filePath, store)
	if !result.HasErrors() {
		t.Fatal("expected missing markdown link validation error")
	}
}

func TestValidateCardFileInStoreRejectsWikiLinks(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)

	if err := os.MkdirAll(store.ProposalCardsDir("CR260612"), 0755); err != nil {
		t.Fatalf("creating proposal cards dir failed: %v", err)
	}

	card := &Card{
		ID:         "STR-abc-wiki",
		Title:      "Wiki structure",
		Type:       CardTypeStructure,
		Status:     CardStatusActive,
		Importance: ImportanceShould,
		Links:      []Link{{Target: "PROP-CR260612", Relation: "belongs_to"}},
		Created:    time.Now(),
		Updated:    time.Now(),
		Body:       "## Entries\n\n- [[REQ-MISSING]] - missing",
	}

	root := NewCard(CardTypeProposal, "Root")
	root.ID = "PROP-CR260612"
	root.Status = CardStatusActive
	proposalDir := store.ProposalDir("CR260612")
	if err := os.MkdirAll(proposalDir, 0755); err != nil {
		t.Fatalf("creating proposal dir failed: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(proposalDir, "90-cards"), 0755); err != nil {
		t.Fatalf("creating cards dir failed: %v", err)
	}
	if err := root.Save(filepath.Join(proposalDir, "ROOT-CR260612.md")); err != nil {
		t.Fatalf("saving root failed: %v", err)
	}

	filePath := filepath.Join(store.ProposalCardsDir("CR260612"), GenerateFilename(card.ID, card.Title))
	if err := card.Save(filePath); err != nil {
		t.Fatalf("saving card failed: %v", err)
	}

	result := ValidateCardFileInStore(filePath, store)
	if !result.HasErrors() {
		t.Fatal("expected wikilink validation error")
	}
}

func TestValidateCardFileInStoreRequiresFrontmatterOutboundLink(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)

	if err := os.MkdirAll(store.ProposalCardsDir("CR260612"), 0755); err != nil {
		t.Fatalf("creating proposal cards dir failed: %v", err)
	}

	target := NewCard(CardTypeProposal, "Root")
	target.ID = "PROP-CR260612"
	target.Status = CardStatusActive
	if err := target.Save(store.ProposalRootCardPath("CR260612")); err != nil {
		t.Fatalf("saving target failed: %v", err)
	}

	card := &Card{
		ID:         "REQ-abc-orphan",
		Title:      "Orphan with body link",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    time.Now(),
		Updated:    time.Now(),
		Body:       "[Root](../PROP-CR260612.md)",
	}
	filePath := filepath.Join(store.ProposalCardsDir("CR260612"), GenerateFilename(card.ID, card.Title))
	if err := card.Save(filePath); err != nil {
		t.Fatalf("saving card failed: %v", err)
	}

	result := ValidateCardFileInStore(filePath, store)
	if !result.HasErrors() {
		t.Fatal("expected missing outbound frontmatter link validation error")
	}
}

func TestValidateCardFileInStoreSkipsExternalAndAnchorMarkdownLinks(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewCardStore(tmpDir)

	if err := os.MkdirAll(store.ProposalCardsDir("CR260612"), 0755); err != nil {
		t.Fatalf("creating proposal cards dir failed: %v", err)
	}

	root := NewCard(CardTypeProposal, "Root")
	root.ID = "PROP-CR260612"
	root.Status = CardStatusActive
	if err := root.Save(store.ProposalRootCardPath("CR260612")); err != nil {
		t.Fatalf("saving root failed: %v", err)
	}

	card := &Card{
		ID:         "REQ-abc-external",
		Title:      "External link card",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Links:      []Link{{Target: "PROP-CR260612", Relation: "belongs_to"}},
		Created:    time.Now(),
		Updated:    time.Now(),
		Body:       "[site](https://example.com)\n[mail](mailto:test@example.com)\n[anchor](#local)",
	}
	filePath := filepath.Join(store.ProposalCardsDir("CR260612"), GenerateFilename(card.ID, card.Title))
	if err := card.Save(filePath); err != nil {
		t.Fatalf("saving card failed: %v", err)
	}

	result := ValidateCardFileInStore(filePath, store)
	if result.HasErrors() {
		t.Fatalf("expected external/anchor links to be skipped, got: %s", result.String())
	}
}

func TestValidationResultString(t *testing.T) {
	result := &ValidationResult{}
	if result.String() != "valid" {
		t.Errorf("expected 'valid' for empty result, got %s", result.String())
	}

	result.AddError("field1", "error1")
	result.AddError("field2", "error2")

	str := result.String()
	if str == "valid" {
		t.Error("expected error string, got 'valid'")
	}

	if !contains(str, "field1") || !contains(str, "error1") {
		t.Error("expected field1 and error1 in string")
	}
}

func TestIsValidRelation(t *testing.T) {
	validRelations := []string{
		"references", "extends", "refines", "contradicts",
		"supersedes", "supports", "questions", "related",
		"implements", "satisfies", "blocks", "produced",
		"indexes", "decomposes", "analyzes", "designs",
		"constrains", "records", "discovers", "belongs_to",
		"requires",
	}

	for _, rel := range validRelations {
		if !isValidRelation(rel) {
			t.Errorf("expected %s to be valid relation", rel)
		}
	}

	if isValidRelation("invalid") {
		t.Error("expected 'invalid' to be invalid relation")
	}
}

func TestIsValidTaskType(t *testing.T) {
	validTypes := []string{"a", "i", "t", "d", "f", "r", "c"}

	for _, tt := range validTypes {
		if !isValidTaskType(tt) {
			t.Errorf("expected %s to be valid task type", tt)
		}
	}

	if isValidTaskType("x") {
		t.Error("expected 'x' to be invalid task type")
	}
}
