package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCardTypeValid(t *testing.T) {
	tests := []struct {
		cardType CardType
		valid    bool
	}{
		{CardTypeRequirement, true},
		{CardTypeDecision, true},
		{CardTypeDesign, true},
		{CardTypeTask, true},
		{CardTypeLog, true},
		{CardTypeConvention, true},
		{CardTypeFinding, true},
		{CardTypeModule, true},
		{CardTypeStructure, true},
		{CardType("invalid"), false},
	}

	for _, tt := range tests {
		if got := tt.cardType.Valid(); got != tt.valid {
			t.Errorf("CardType(%q).Valid() = %v, want %v", tt.cardType, got, tt.valid)
		}
	}
}

func TestCardTypePrefix(t *testing.T) {
	tests := []struct {
		cardType CardType
		prefix   string
	}{
		{CardTypeRequirement, "REQ"},
		{CardTypeDecision, "DEC"},
		{CardTypeDesign, "DES"},
		{CardTypeTask, "TASK"},
		{CardTypeLog, "LOG"},
		{CardTypeConvention, "CONV"},
		{CardTypeFinding, "FIND"},
		{CardTypeModule, "MOD"},
		{CardTypeStructure, "STR"},
	}

	for _, tt := range tests {
		if got := tt.cardType.Prefix(); got != tt.prefix {
			t.Errorf("CardType(%q).Prefix() = %q, want %q", tt.cardType, got, tt.prefix)
		}
	}
}

func TestCardStatusValid(t *testing.T) {
	tests := []struct {
		status CardStatus
		valid  bool
	}{
		{CardStatusDraft, true},
		{CardStatusActive, true},
		{CardStatusAccepted, true},
		{CardStatusDeprecated, true},
		{CardStatusSuperseded, true},
		{CardStatusBacklog, true},
		{CardStatusNotReady, true},
		{CardStatusReady, true},
		{CardStatusInProgress, true},
		{CardStatusDone, true},
		{CardStatusBlocked, true},
		{CardStatusCancelled, true},
		{CardStatus("invalid"), false},
	}

	for _, tt := range tests {
		if got := tt.status.Valid(); got != tt.valid {
			t.Errorf("CardStatus(%q).Valid() = %v, want %v", tt.status, got, tt.valid)
		}
	}
}

func TestImportanceValid(t *testing.T) {
	tests := []struct {
		importance Importance
		valid      bool
	}{
		{ImportanceMust, true},
		{ImportanceShould, true},
		{ImportanceMay, true},
		{Importance(""), true},
		{Importance("invalid"), false},
	}

	for _, tt := range tests {
		if got := tt.importance.Valid(); got != tt.valid {
			t.Errorf("Importance(%q).Valid() = %v, want %v", tt.importance, got, tt.valid)
		}
	}
}

func TestNewCard(t *testing.T) {
	card := NewCard(CardTypeRequirement, "Test Card")

	if card.Type != CardTypeRequirement {
		t.Errorf("expected type %s, got %s", CardTypeRequirement, card.Type)
	}

	if card.Title != "Test Card" {
		t.Errorf("expected title 'Test Card', got %s", card.Title)
	}

	if card.Status != CardStatusDraft {
		t.Errorf("expected status draft, got %s", card.Status)
	}

	if card.Importance != ImportanceShould {
		t.Errorf("expected importance should, got %s", card.Importance)
	}

	if card.Created.IsZero() {
		t.Error("expected created time to be set")
	}
}

func TestCardAddLink(t *testing.T) {
	card := NewCard(CardTypeRequirement, "Test")

	card.AddLink("DEC-123", "references")

	if len(card.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(card.Links))
	}

	if card.Links[0].Target != "DEC-123" {
		t.Errorf("expected target DEC-123, got %s", card.Links[0].Target)
	}

	if card.Links[0].Relation != "references" {
		t.Errorf("expected relation references, got %s", card.Links[0].Relation)
	}

	card.AddLink("DEC-123", "references")
	if len(card.Links) != 1 {
		t.Error("expected duplicate link to be ignored")
	}
}

func TestCardRemoveLink(t *testing.T) {
	card := NewCard(CardTypeRequirement, "Test")
	card.AddLink("DEC-123", "references")
	card.AddLink("DEC-456", "implements")

	removed := card.RemoveLink("DEC-123", "references")
	if !removed {
		t.Error("expected link to be removed")
	}

	if len(card.Links) != 1 {
		t.Errorf("expected 1 link, got %d", len(card.Links))
	}

	removed = card.RemoveLink("DEC-999", "references")
	if removed {
		t.Error("expected non-existent link removal to return false")
	}
}

func TestParseCard(t *testing.T) {
	content := `---
id: REQ-123
title: Test Requirement
type: requirement
status: draft
importance: should
tags:
  - cli
  - init
links:
  - target: DEC-456
    relation: references
created: 2026-06-12T10:00:00Z
updated: 2026-06-12T10:00:00Z
source: CR260612
domain: cli
---

# Test Requirement

This is the body content.
`

	card, err := ParseCard([]byte(content), "/test/path.md")
	if err != nil {
		t.Fatalf("failed to parse card: %v", err)
	}

	if card.ID != "REQ-123" {
		t.Errorf("expected ID REQ-123, got %s", card.ID)
	}

	if card.Title != "Test Requirement" {
		t.Errorf("expected title 'Test Requirement', got %s", card.Title)
	}

	if card.Type != CardTypeRequirement {
		t.Errorf("expected type requirement, got %s", card.Type)
	}

	if card.Status != CardStatusDraft {
		t.Errorf("expected status draft, got %s", card.Status)
	}

	if len(card.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(card.Tags))
	}

	if len(card.Links) != 1 {
		t.Errorf("expected 1 link, got %d", len(card.Links))
	}

	if card.Source != "CR260612" {
		t.Errorf("expected source CR260612, got %s", card.Source)
	}

	if card.Body != "# Test Requirement\n\nThis is the body content." {
		t.Errorf("unexpected body: %q", card.Body)
	}
}

func TestParseCardMissingFrontmatter(t *testing.T) {
	content := `# No Frontmatter

Just content.
`

	_, err := ParseCard([]byte(content), "/test/path.md")
	if err == nil {
		t.Error("expected error for missing frontmatter")
	}
}

func TestCardToMarkdown(t *testing.T) {
	card := &Card{
		ID:         "REQ-123",
		Title:      "Test Card",
		Type:       CardTypeRequirement,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Tags:       []string{"test", "cli"},
		Links: []Link{
			{Target: "DEC-456", Relation: "references"},
		},
		Created: time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC),
		Updated: time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC),
		Source:  "CR260612",
		Domain:  "cli",
		Body:    "# Test\n\nBody content.",
	}

	data, err := card.ToMarkdown()
	if err != nil {
		t.Fatalf("failed to convert to markdown: %v", err)
	}

	content := string(data)

	if !contains(content, "id: REQ-123") {
		t.Error("expected ID in frontmatter")
	}

	if !contains(content, "title: Test Card") {
		t.Error("expected title in frontmatter")
	}

	if !contains(content, "type: requirement") {
		t.Error("expected type in frontmatter")
	}

	if !contains(content, "# Test") {
		t.Error("expected body content")
	}
}

func TestCardSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test-card.md")

	card := NewCard(CardTypeDecision, "Test Decision")
	card.ID = "DEC-123"
	card.Body = "# Decision\n\nWe chose this approach."

	if err := card.Save(filePath); err != nil {
		t.Fatalf("failed to save card: %v", err)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("card file not created: %v", err)
	}

	loaded, err := ParseCardFile(filePath)
	if err != nil {
		t.Fatalf("failed to load card: %v", err)
	}

	if loaded.ID != card.ID {
		t.Errorf("expected ID %s, got %s", card.ID, loaded.ID)
	}

	if loaded.Title != card.Title {
		t.Errorf("expected title %s, got %s", card.Title, loaded.Title)
	}

	if loaded.Body != card.Body {
		t.Errorf("expected body %s, got %s", card.Body, loaded.Body)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
