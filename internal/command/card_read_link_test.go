package command

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestCardReadSummaryAndSection(t *testing.T) {
	tmpDir := prepareCardCommandFixture(t)

	summaryCmd := newCardReadCmd()
	var summaryOut bytes.Buffer
	summaryCmd.SetOut(&summaryOut)
	summaryCmd.SetArgs([]string{"DEC-260613-01", "--summary"})
	if err := summaryCmd.Execute(); err != nil {
		t.Fatalf("card read summary failed: %v", err)
	}

	summaryText := summaryOut.String()
	for _, want := range []string{
		"ID: DEC-260613-01",
		"Title: Link summary card",
		"Summary: Decision body paragraph.",
		"First section: Decision",
		"Sections:",
		"Links:",
		"Accept",
		"Decision",
	} {
		if !strings.Contains(summaryText, want) {
			t.Fatalf("summary output missing %q:\n%s", want, summaryText)
		}
	}

	sectionCmd := newCardReadCmd()
	var sectionOut bytes.Buffer
	sectionCmd.SetOut(&sectionOut)
	sectionCmd.SetArgs([]string{"DEC-260613-01", "--section", "decision"})
	if err := sectionCmd.Execute(); err != nil {
		t.Fatalf("card read section failed: %v", err)
	}

	sectionText := sectionOut.String()
	for _, want := range []string{
		"ID: DEC-260613-01",
		"Title: Link summary card",
		"Decision body paragraph.",
		"Second paragraph.",
	} {
		if !strings.Contains(sectionText, want) {
			t.Fatalf("section output missing %q:\n%s", want, sectionText)
		}
	}
	if strings.Contains(sectionText, "Acceptance") {
		t.Fatalf("section output unexpectedly included another section:\n%s", sectionText)
	}

	_ = tmpDir
}

func TestCardLinkAndUnlinkCommands(t *testing.T) {
	tmpDir := prepareCardCommandFixture(t)
	store := testCardStore(t, tmpDir)

	linkCmd := newCardLinkCmd()
	var linkOut bytes.Buffer
	linkCmd.SetOut(&linkOut)
	linkCmd.SetArgs([]string{"DEC-260613-01", "REQ-260613-01", "--relation", "references"})
	if err := linkCmd.Execute(); err != nil {
		t.Fatalf("card link failed: %v", err)
	}
	if !strings.Contains(linkOut.String(), "✓ Linked DEC-260613-01 -> REQ-260613-01 (references)") {
		t.Fatalf("unexpected link output:\n%s", linkOut.String())
	}

	card, err := store.ReadCard("DEC-260613-01")
	if err != nil {
		t.Fatalf("reading linked card failed: %v", err)
	}
	if !hasLinkRelation(card, "REQ-260613-01", "references") {
		t.Fatalf("expected link to be written to card")
	}

	unlinkCmd := newCardUnlinkCmd()
	var unlinkOut bytes.Buffer
	unlinkCmd.SetOut(&unlinkOut)
	unlinkCmd.SetArgs([]string{"DEC-260613-01", "REQ-260613-01", "--relation", "references"})
	if err := unlinkCmd.Execute(); err != nil {
		t.Fatalf("card unlink failed: %v", err)
	}
	if !strings.Contains(unlinkOut.String(), "✓ Unlinked DEC-260613-01 -> REQ-260613-01 (references)") {
		t.Fatalf("unexpected unlink output:\n%s", unlinkOut.String())
	}

	card, err = store.ReadCard("DEC-260613-01")
	if err != nil {
		t.Fatalf("reading unlinked card failed: %v", err)
	}
	if hasLinkRelation(card, "REQ-260613-01", "references") {
		t.Fatalf("expected link to be removed from card")
	}
}

func TestCardCreateAndUpdateParseCommaSeparatedLinks(t *testing.T) {
	tmpDir := prepareCardCommandFixture(t)
	store := testCardStore(t, tmpDir)

	createCmd := newCardCreateCmd()
	createCmd.SetArgs([]string{
		"--type", "design",
		"--title", "Comma linked design",
		"--links", "REQ-260613-01:requires,DEC-260613-01:references",
	})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("card create failed: %v", err)
	}

	cards, err := store.ListCardsByType(core.CardTypeDesign)
	if err != nil {
		t.Fatalf("listing design cards failed: %v", err)
	}
	var created *core.Card
	for _, card := range cards {
		if card.Title == "Comma linked design" {
			created = card
			break
		}
	}
	if created == nil {
		t.Fatal("expected created design card")
	}
	if !hasLinkRelation(created, "REQ-260613-01", "requires") || !hasLinkRelation(created, "DEC-260613-01", "references") {
		t.Fatalf("created card links not parsed correctly: %#v", created.Links)
	}

	updateCmd := newCardUpdateCmd()
	updateCmd.SetArgs([]string{created.ID, "--add-link", "TASK-260613-i-01:implements,LOG-260613-01:records"})
	if err := updateCmd.Execute(); err != nil {
		t.Fatalf("card update failed: %v", err)
	}
	updated, err := store.ReadCard(created.ID)
	if err != nil {
		t.Fatalf("reading updated card failed: %v", err)
	}
	if !hasLinkRelation(updated, "TASK-260613-i-01", "implements") || !hasLinkRelation(updated, "LOG-260613-01", "records") {
		t.Fatalf("updated card links not parsed correctly: %#v", updated.Links)
	}
}

func TestCardRelatedBacklinksDirection(t *testing.T) {
	prepareCardCommandFixture(t)

	relatedCmd := newCardRelatedCmd()
	var out bytes.Buffer
	relatedCmd.SetOut(&out)
	relatedCmd.SetArgs([]string{"REQ-260613-01", "--direction", "backlinks", "--relation", "references"})
	if err := relatedCmd.Execute(); err != nil {
		t.Fatalf("card related backlinks failed: %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "DEC-260613-01") || !strings.Contains(text, "Link summary card") {
		t.Fatalf("backlinks output missing dependent card:\n%s", text)
	}
	if strings.Contains(text, "TASK-260613-01") {
		t.Fatalf("backlinks output should respect relation filter:\n%s", text)
	}
}

func prepareCardCommandFixture(t *testing.T) string {
	t.Helper()

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

	req := core.NewCard(core.CardTypeRequirement, "Linked requirement")
	req.ID = "REQ-260613-01"
	req.Body = "# Requirement\n\nRequirement body paragraph.\n\n## Scope\n\n- First scope"
	if _, err := store.CreateCard(req, "CR26061301"); err != nil {
		t.Fatalf("creating requirement card failed: %v", err)
	}

	dec := core.NewCard(core.CardTypeDecision, "Link summary card")
	dec.ID = "DEC-260613-01"
	dec.Body = "# Decision\n\nDecision body paragraph.\n\nSecond paragraph.\n\n## Acceptance\n\n- A"
	dec.AddLink("REQ-260613-01", "references")
	if _, err := store.CreateCard(dec, "CR26061301"); err != nil {
		t.Fatalf("creating decision card failed: %v", err)
	}

	task := core.NewCard(core.CardTypeTask, "Other backlink")
	task.ID = "TASK-260613-01"
	task.Body = "# Task\n\nTask body."
	task.AddLink("REQ-260613-01", "implements")
	if _, err := store.CreateCard(task, "CR26061301"); err != nil {
		t.Fatalf("creating task card failed: %v", err)
	}

	return tmpDir
}
