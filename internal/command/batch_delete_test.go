package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"flowforge/internal/core"
)

func TestCardBatchStdin(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Batch stdin proposal")

	yaml := fmt.Sprintf(`proposal: %s
cards:
  - ref: "stdin-req"
    type: feature
    title: "Stdin Batch Requirement"
    status: draft
    body: |
      # Stdin Batch Requirement
      A requirement created via stdin.
    tags: [batch, stdin]
    domain: flowforge
  - type: feature
    title: "Stdin Batch Design"
    status: draft
    body: |
      # Stdin Batch Design
      A design created via stdin.
    tags: [batch, stdin]
    domain: flowforge
    links:
      - "@stdin-req:references"
`, proposalID)

	cmd := newCardCreateBatchCmd()
	stdin := strings.NewReader(yaml)
	cmd.SetIn(stdin)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"-"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card batch - (stdin) failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "✓ Created 2 card(s)") {
		t.Fatalf("expected 2 cards created, got:\n%s", result)
	}

	store := testCardStore(t, tmpDir)
	cards, err := store.ListCards(store.ProposalCardsDir(proposalID))
	if err != nil {
		t.Fatalf("listing proposal cards: %v", err)
	}

	var reqFound, desFound bool
	for _, c := range cards {
		if c.Type == core.CardTypeFeature && c.Title == "Stdin Batch Requirement" {
			reqFound = true
		}
		if c.Type == core.CardTypeFeature && c.Title == "Stdin Batch Design" {
			desFound = true
		}
	}
	if !reqFound {
		t.Fatalf("requirement card not found from stdin batch")
	}
	if !desFound {
		t.Fatalf("design card not found from stdin batch")
	}
}

func TestCardBatchForwardRefResolution(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Batch forward ref proposal")

	yaml := fmt.Sprintf(`proposal: %s
cards:
  - ref: "design-ref"
    type: feature
    title: "Batch Forward Ref Design"
    status: draft
    body: |
      # Batch Forward Ref Design

      ## Summary
      Design card referenced by a later card via @ref.

      ## Source
      test

      ## Acceptance
      - @task-ref can resolve to this card's ID
    tags: [batch, ref-test]
    domain: flowforge

  - type: feature
    title: "Batch Forward Ref Task (implements design)"
    status: not_ready
    body: |
      # Batch Forward Ref Task

      ## Goal
      Verify @ref forward-reference works across batch cards.

      ## Inputs
      - @design-ref (design card created earlier in same batch)
    tags: [batch, ref-test]
    domain: flowforge
    links:
      - "@design-ref:implements"
`, proposalID)

	manifestPath := filepath.Join(tmpDir, "batch-test.yaml")
	if err := os.WriteFile(manifestPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("writing batch yaml: %v", err)
	}

	cmd := newCardCreateBatchCmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{manifestPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("card batch failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "✓ Created 2 card(s)") {
		t.Fatalf("expected 2 cards created, got:\n%s", result)
	}

	store := testCardStore(t, tmpDir)

	// Find the task card and verify it has implements link to the design card
	cards, err := store.ListCards(store.ProposalCardsDir(proposalID))
	if err != nil {
		t.Fatalf("listing proposal cards: %v", err)
	}

	var taskCard *core.Card
	for _, c := range cards {
		if c.Type == core.CardTypeFeature && strings.Contains(c.Title, "Batch Forward Ref Task") {
			taskCard = c
			break
		}
	}
	if taskCard == nil {
		t.Fatalf("task card not found in proposal cards")
	}

	found := false
	for _, link := range taskCard.Links {
		if link.Relation == "implements" {
			targetCard, terr := store.ReadCard(link.Target)
			if terr != nil {
				t.Fatalf("reading linked design card %s: %v", link.Target, terr)
			}
			if targetCard.Type == core.CardTypeFeature && targetCard.Title == "Batch Forward Ref Design" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatalf("task card missing implements link to design card; links: %#v", taskCard.Links)
	}
}

func TestCardBatchPhaseOneDuplicateIDRetainsSuccessfulPeers(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Batch duplicate ID proposal")
	store := testCardStore(t, tmpDir)
	existing := core.NewCard(core.CardTypeFeature, "Existing card")
	existing.ID = "FEAT-" + proposalID + "-001"
	existing.Body = "existing body"
	if _, err := store.CreateCard(existing, proposalID); err != nil {
		t.Fatalf("creating existing card: %v", err)
	}
	sequencePath := filepath.Join(store.ProposalCardsDir(proposalID), ".flowforge-card-sequence")
	if err := os.WriteFile(sequencePath, []byte("1\n"), 0600); err != nil {
		t.Fatalf("writing sequence counter: %v", err)
	}

	manifest := fmt.Sprintf(`proposal: %s
cards:
  - type: feature
    title: "Duplicate card"
  - type: feature
    title: "Successful peer"
`, proposalID)
	cmd := newCardCreateBatchCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--manifest", manifest})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected batch to report the duplicate ID")
	}
	text := out.String()
	if !strings.Contains(text, "phase 1") || !strings.Contains(text, "manifest[0]") || !strings.Contains(text, "FEAT-"+proposalID+"-001") {
		t.Fatalf("duplicate error is not auditable: %s", text)
	}
	if !strings.Contains(text, "Successful peer") {
		t.Fatalf("successful peer missing from report: %s", text)
	}
	peer, err := store.ReadCard("FEAT-" + proposalID + "-002")
	if err != nil {
		t.Fatalf("successful peer was not retained: %v", err)
	}
	if peer.Title != "Successful peer" {
		t.Fatalf("unexpected retained peer: %#v", peer)
	}
	reloaded, err := store.ReadCard(existing.ID)
	if err != nil {
		t.Fatalf("reading existing card: %v", err)
	}
	if reloaded.Body != existing.Body {
		t.Fatalf("duplicate ID overwrote existing body: %q", reloaded.Body)
	}
	next, err := store.NextCardID(core.CardTypeFeature, proposalID)
	if err != nil {
		t.Fatalf("allocating subsequent ID: %v", err)
	}
	if next != "FEAT-"+proposalID+"-003" {
		t.Fatalf("sequence counter rolled back or collided: %s", next)
	}
}

func TestCardBatchJSONPartialFailureReportIsStable(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Batch JSON proposal")
	store := testCardStore(t, tmpDir)
	existing := core.NewCard(core.CardTypeFeature, "Existing JSON card")
	existing.ID = "FEAT-" + proposalID + "-001"
	if _, err := store.CreateCard(existing, proposalID); err != nil {
		t.Fatalf("creating existing card: %v", err)
	}
	if err := os.WriteFile(filepath.Join(store.ProposalCardsDir(proposalID), ".flowforge-card-sequence"), []byte("1\n"), 0600); err != nil {
		t.Fatalf("writing sequence counter: %v", err)
	}
	viper.Set("output", "json")
	defer viper.Set("output", "text")
	cmd := newCardCreateBatchCmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--manifest", fmt.Sprintf("proposal: %s\ncards:\n  - type: feature\n    title: Duplicate JSON\n", proposalID)})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected JSON batch partial failure")
	}
	var result batchResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON batch report: %v; output=%s", err, out.String())
	}
	if len(result.Errors) != 1 || result.Errors[0].Index != 0 || !strings.Contains(result.Errors[0].Error, "phase 1") || !strings.Contains(result.Errors[0].Error, existing.ID) {
		t.Fatalf("unstable or incomplete JSON error report: %#v", result)
	}
}

func TestCardBatchPhaseTwoBadTargetRetainsCreatedLinksAndNavigation(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Batch bad target proposal")
	store := testCardStore(t, tmpDir)
	target := core.NewCard(core.CardTypeDecision, "Existing target")
	target.ID = "DEC-" + proposalID + "-target"
	if _, err := store.CreateCard(target, proposalID); err != nil {
		t.Fatalf("creating target: %v", err)
	}
	manifest := fmt.Sprintf(`proposal: %s
cards:
  - type: feature
    title: "Partial links"
    links:
      - "%s:references"
      - "DEC-missing-target:references"
`, proposalID, target.ID)
	cmd := newCardCreateBatchCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--manifest", manifest})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected batch to report the bad target")
	}
	text := out.String()
	if !strings.Contains(text, "phase 2") || !strings.Contains(text, "target DEC-missing-target") {
		t.Fatalf("bad target error is not auditable: %s", text)
	}
	created, err := store.ReadCard("FEAT-" + proposalID + "-001")
	if err != nil {
		t.Fatalf("created card was not retained: %v", err)
	}
	if !hasLinkRelation(created, target.ID, "references") {
		t.Fatalf("valid link was not retained: %#v", created.Links)
	}
	if strings.Contains(created.Body, "DEC-missing-target") {
		t.Fatalf("bad target leaked into generated navigation: %s", created.Body)
	}
	updatedTarget, err := store.ReadCard(target.ID)
	if err != nil {
		t.Fatalf("reading navigation target: %v", err)
	}
	if !strings.Contains(updatedTarget.Body, created.ID) {
		t.Fatalf("existing target navigation was not retained: %s", updatedTarget.Body)
	}
}

func TestCardBatchDuplicateRefRejectedBeforeWrites(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)
	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Batch duplicate ref proposal")
	manifest := fmt.Sprintf(`proposal: %s
cards:
  - ref: same
    type: feature
    title: "First"
  - ref: same
    type: decision
    title: "Second"
`, proposalID)
	cmd := newCardCreateBatchCmd()
	cmd.SetArgs([]string{"--manifest", manifest})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected duplicate ref preflight rejection")
	}
	store := testCardStore(t, tmpDir)
	cards, err := store.ListCards(store.ProposalCardsDir(proposalID))
	if err != nil {
		t.Fatalf("listing proposal cards: %v", err)
	}
	if len(cards) != 0 {
		t.Fatalf("duplicate ref preflight wrote cards: %#v", cards)
	}
}

func TestForceDeleteCardRemovesBacklinks(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Force delete proposal")

	store := testCardStore(t, tmpDir)

	// Create target card (the one to be force-deleted)
	targetCard := core.NewCard(core.CardTypeFinding, "Target card for force-delete test")
	targetCard.Status = core.CardStatusActive
	targetCard.ID = core.GenerateCardID(core.CardTypeFinding, proposalTimestamp(proposalID))
	targetCard.AddLink("PROP-"+proposalID, "belongs_to")
	_, err := store.CreateCard(targetCard, proposalID)
	if err != nil {
		t.Fatalf("creating target card: %v", err)
	}

	// Create source card that references the target card
	sourceCard := core.NewCard(core.CardTypeFinding, "Source card with backlink")
	sourceCard.Status = core.CardStatusActive
	sourceCard.ID = core.GenerateCardID(core.CardTypeFinding, proposalTimestamp(proposalID))
	sourceCard.AddLink(targetCard.ID, "references")
	sourceCard.AddLink("PROP-"+proposalID, "belongs_to")
	_, err = store.CreateCard(sourceCard, proposalID)
	if err != nil {
		t.Fatalf("creating source card: %v", err)
	}

	// Verify the backlink exists before deletion
	sourceAfter, err := store.ReadCard(sourceCard.ID)
	if err != nil {
		t.Fatalf("reading source card before delete: %v", err)
	}
	preDeleteHasLink := false
	for _, link := range sourceAfter.Links {
		if link.Target == targetCard.ID {
			preDeleteHasLink = true
			break
		}
	}
	if !preDeleteHasLink {
		t.Fatalf("expected source card to have backlink to target before delete")
	}

	// Force-delete the target card
	if err := store.ForceDeleteCard(targetCard.ID); err != nil {
		t.Fatalf("ForceDeleteCard failed: %v", err)
	}

	// Verify the target card file is gone
	if _, err := store.ReadCard(targetCard.ID); err == nil {
		t.Fatalf("expected target card %s to be deleted, but it still exists", targetCard.ID)
	}

	// Verify backlinks were cleaned up from source card
	sourceReloaded, err := store.ReadCard(sourceCard.ID)
	if err != nil {
		t.Fatalf("reading source card after force delete: %v", err)
	}
	for _, link := range sourceReloaded.Links {
		if link.Target == targetCard.ID {
			t.Fatalf("backlink to deleted card %s still present in source card %s", targetCard.ID, sourceCard.ID)
		}
	}
}

func TestForceDeleteCardNonDraft(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Force delete non-draft proposal")

	store := testCardStore(t, tmpDir)

	// Create an active (non-draft) card
	activeCard := core.NewCard(core.CardTypeFinding, "Active card to force delete")
	activeCard.Status = core.CardStatusActive
	activeCard.ID = core.GenerateCardID(core.CardTypeFinding, proposalTimestamp(proposalID))
	activeCard.AddLink("PROP-"+proposalID, "belongs_to")
	_, err := store.CreateCard(activeCard, proposalID)
	if err != nil {
		t.Fatalf("creating active card: %v", err)
	}

	// Normal DeleteCard should reject non-draft
	if err := store.DeleteCard(activeCard.ID); err == nil {
		t.Fatalf("expected DeleteCard to reject non-draft card %s", activeCard.ID)
	}

	// ForceDeleteCard should succeed regardless of status
	if err := store.ForceDeleteCard(activeCard.ID); err != nil {
		t.Fatalf("ForceDeleteCard should delete non-draft card: %v", err)
	}

	if _, err := store.ReadCard(activeCard.ID); err == nil {
		t.Fatalf("card %s should have been deleted by ForceDeleteCard", activeCard.ID)
	}
}
