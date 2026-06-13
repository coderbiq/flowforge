package command

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestStructureCommandsAddListRemove(t *testing.T) {
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
	structureCard := core.NewCard(core.CardTypeStructure, "Top-level structure")
	structureCard.ID = "STR-ROOT"
	if _, err := store.CreateCard(structureCard, ""); err != nil {
		t.Fatalf("creating structure card failed: %v", err)
	}

	linkedCard := core.NewCard(core.CardTypeRequirement, "Indexed requirement")
	linkedCard.ID = "REQ-001"
	if _, err := store.CreateCard(linkedCard, ""); err != nil {
		t.Fatalf("creating linked card failed: %v", err)
	}

	addCmd := newStructureAddCmd()
	var addOut bytes.Buffer
	addCmd.SetOut(&addOut)
	addCmd.SetArgs([]string{"STR-ROOT", "REQ-001"})
	if err := addCmd.Execute(); err != nil {
		t.Fatalf("structure add failed: %v", err)
	}
	if !strings.Contains(addOut.String(), "✓ Added REQ-001 to STR-ROOT") {
		t.Fatalf("unexpected add output:\n%s", addOut.String())
	}

	reloaded, err := store.ReadCard("STR-ROOT")
	if err != nil {
		t.Fatalf("reading structure card failed: %v", err)
	}
	if len(reloaded.Links) != 1 {
		t.Fatalf("expected 1 link after add, got %d", len(reloaded.Links))
	}
	if reloaded.Links[0].Target != "REQ-001" || reloaded.Links[0].Relation != "indexes" {
		t.Fatalf("unexpected structure link: %#v", reloaded.Links[0])
	}

	dupCmd := newStructureAddCmd()
	dupCmd.SetArgs([]string{"STR-ROOT", "REQ-001"})
	if err := dupCmd.Execute(); err != nil {
		t.Fatalf("duplicate structure add failed: %v", err)
	}
	reloaded, err = store.ReadCard("STR-ROOT")
	if err != nil {
		t.Fatalf("reading structure card after duplicate add failed: %v", err)
	}
	if len(reloaded.Links) != 1 {
		t.Fatalf("expected duplicate add to keep 1 link, got %d", len(reloaded.Links))
	}

	listCmd := newStructureListCmd()
	var listOut bytes.Buffer
	listCmd.SetOut(&listOut)
	listCmd.SetArgs([]string{"STR-ROOT"})
	if err := listCmd.Execute(); err != nil {
		t.Fatalf("structure list failed: %v", err)
	}
	if !strings.Contains(listOut.String(), "REQ-001 [requirement] Indexed requirement") {
		t.Fatalf("unexpected list output:\n%s", listOut.String())
	}

	removeCmd := newStructureRemoveCmd()
	var removeOut bytes.Buffer
	removeCmd.SetOut(&removeOut)
	removeCmd.SetArgs([]string{"STR-ROOT", "REQ-001"})
	if err := removeCmd.Execute(); err != nil {
		t.Fatalf("structure remove failed: %v", err)
	}
	if !strings.Contains(removeOut.String(), "✓ Removed REQ-001 from STR-ROOT") {
		t.Fatalf("unexpected remove output:\n%s", removeOut.String())
	}

	reloaded, err = store.ReadCard("STR-ROOT")
	if err != nil {
		t.Fatalf("reading structure card after remove failed: %v", err)
	}
	if len(reloaded.Links) != 0 {
		t.Fatalf("expected 0 links after remove, got %d", len(reloaded.Links))
	}
}

func TestStructureCommandsRejectNonStructureCard(t *testing.T) {
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
	card := core.NewCard(core.CardTypeRequirement, "Not a structure")
	card.ID = "REQ-100"
	if _, err := store.CreateCard(card, ""); err != nil {
		t.Fatalf("creating non-structure card failed: %v", err)
	}

	cmd := newStructureAddCmd()
	cmd.SetArgs([]string{"REQ-100", "REQ-100"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected add to reject non-structure card")
	}
}
