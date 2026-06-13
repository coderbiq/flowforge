package core

import (
	"strings"
	"testing"
)

func TestGenerateCardTimestamp(t *testing.T) {
	ts := GenerateCardTimestamp()

	if ts == "" {
		t.Error("expected non-empty timestamp")
	}

	if len(ts) < 6 {
		t.Errorf("expected timestamp length >= 6, got %d", len(ts))
	}
}

func TestGenerateCardID(t *testing.T) {
	tests := []struct {
		cardType   CardType
		proposalTs string
		prefix     string
	}{
		{CardTypeRequirement, "abc123", "REQ-abc123-"},
		{CardTypeDecision, "def456", "DEC-def456-"},
		{CardTypeDesign, "", "DES-"},
		{CardTypeTask, "ghi789", "TASK-ghi789-"},
	}

	for _, tt := range tests {
		id := GenerateCardID(tt.cardType, tt.proposalTs)

		if !strings.HasPrefix(id, tt.prefix) {
			t.Errorf("GenerateCardID(%s, %s) = %s, expected prefix %s", tt.cardType, tt.proposalTs, id, tt.prefix)
		}

		if tt.proposalTs == "" && strings.Count(id, "-") < 1 {
			t.Errorf("expected at least 1 dash in ID without proposal, got %s", id)
		}
	}
}

func TestGenerateTaskID(t *testing.T) {
	tests := []struct {
		proposalTs string
		taskType   string
		prefix     string
	}{
		{"abc123", "i", "TASK-abc123-i-"},
		{"def456", "t", "TASK-def456-t-"},
		{"ghi789", "", "TASK-ghi789-i-"},
	}

	for _, tt := range tests {
		id := GenerateTaskID(tt.proposalTs, tt.taskType)

		if !strings.HasPrefix(id, tt.prefix) {
			t.Errorf("GenerateTaskID(%s, %s) = %s, expected prefix %s", tt.proposalTs, tt.taskType, id, tt.prefix)
		}
	}
}

func TestGenerateSubTaskID(t *testing.T) {
	tests := []struct {
		parentID string
		expected string
		hasError bool
	}{
		{"TASK-abc-i-123", "TASK-abc-i-123-a", false},
		{"TASK-abc-i-123-a", "TASK-abc-i-123-b", false},
		{"TASK-abc-i-123-b", "TASK-abc-i-123-c", false},
		{"INVALID", "", true},
	}

	for _, tt := range tests {
		id, err := GenerateSubTaskID(tt.parentID)

		if tt.hasError {
			if err == nil {
				t.Errorf("GenerateSubTaskID(%s) expected error, got %s", tt.parentID, id)
			}
			continue
		}

		if err != nil {
			t.Errorf("GenerateSubTaskID(%s) unexpected error: %v", tt.parentID, err)
			continue
		}

		if id != tt.expected {
			t.Errorf("GenerateSubTaskID(%s) = %s, expected %s", tt.parentID, id, tt.expected)
		}
	}
}

func TestGenerateProposalID(t *testing.T) {
	id := GenerateProposalID()

	if !strings.HasPrefix(id, "CR") {
		t.Errorf("expected proposal ID to start with CR, got %s", id)
	}

	if len(id) < 8 {
		t.Errorf("expected proposal ID length >= 8, got %d", len(id))
	}
}

func TestToSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"  Multiple   Spaces  ", "multiple-spaces"},
		{"Special!@#Characters", "specialcharacters"},
		{"CamelCaseText", "camel-case-text"},
		{"already-slugged", "already-slugged"},
		{"under_score_text", "under-score-text"},
		{"中文标题", "中文标题"},
		{"Mixed 中文 and English", "mixed-中文-and-english"},
		{"", ""},
		{"   ", ""},
	}

	for _, tt := range tests {
		result := ToSlug(tt.input)
		if result != tt.expected {
			t.Errorf("ToSlug(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestToSlugMaxLength(t *testing.T) {
	longInput := strings.Repeat("a", 100)
	result := ToSlug(longInput)

	if len(result) > 50 {
		t.Errorf("expected slug length <= 50, got %d", len(result))
	}
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		id       string
		title    string
		expected string
	}{
		{"REQ-123", "Test Requirement", "REQ-123_test-requirement.md"},
		{"DEC-456", "Use PostgreSQL", "DEC-456_use-postgresql.md"},
		{"TASK-789", "Implement API", "TASK-789_implement-api.md"},
		{"STR-001", "CLI Architecture", "STR-001_cli-architecture.md"},
	}

	for _, tt := range tests {
		filename := GenerateFilename(tt.id, tt.title)
		if filename != tt.expected {
			t.Errorf("GenerateFilename(%s, %s) = %s, expected %s", tt.id, tt.title, filename, tt.expected)
		}
	}
}

func TestParseFilename(t *testing.T) {
	tests := []struct {
		filename     string
		expectedID   string
		expectedSlug string
		hasError     bool
	}{
		{"REQ-123_test-requirement.md", "REQ-123", "test-requirement", false},
		{"DEC-456_use-postgresql.md", "DEC-456", "use-postgresql", false},
		{"TASK-789-i-abc_implement-api.md", "TASK-789-i-abc", "implement-api", false},
		{"invalid-filename.md", "", "", true},
		{"no-extension", "", "", true},
	}

	for _, tt := range tests {
		id, slug, err := ParseFilename(tt.filename)

		if tt.hasError {
			if err == nil {
				t.Errorf("ParseFilename(%s) expected error, got id=%s, slug=%s", tt.filename, id, slug)
			}
			continue
		}

		if err != nil {
			t.Errorf("ParseFilename(%s) unexpected error: %v", tt.filename, err)
			continue
		}

		if id != tt.expectedID {
			t.Errorf("ParseFilename(%s) id = %s, expected %s", tt.filename, id, tt.expectedID)
		}

		if slug != tt.expectedSlug {
			t.Errorf("ParseFilename(%s) slug = %s, expected %s", tt.filename, slug, tt.expectedSlug)
		}
	}
}

func TestParseCardID(t *testing.T) {
	tests := []struct {
		cardID         string
		expectedType   CardType
		expectedPropTs string
		expectedCardTs string
		hasError       bool
	}{
		{"REQ-abc-123", CardTypeRequirement, "abc", "123", false},
		{"DEC-def-456", CardTypeDecision, "def", "456", false},
		{"TASK-ghi-i-789", CardTypeTask, "ghi", "789", false},
		{"TASK-ghi-i-789-a", CardTypeTask, "ghi", "789-a", false},
		{"CONV-001", CardTypeConvention, "", "001", false},
		{"INVALID-123", "", "", "", true},
		{"REQ", "", "", "", true},
	}

	for _, tt := range tests {
		cardType, proposalTs, cardTs, err := ParseCardID(tt.cardID)

		if tt.hasError {
			if err == nil {
				t.Errorf("ParseCardID(%s) expected error", tt.cardID)
			}
			continue
		}

		if err != nil {
			t.Errorf("ParseCardID(%s) unexpected error: %v", tt.cardID, err)
			continue
		}

		if cardType != tt.expectedType {
			t.Errorf("ParseCardID(%s) type = %s, expected %s", tt.cardID, cardType, tt.expectedType)
		}

		if proposalTs != tt.expectedPropTs {
			t.Errorf("ParseCardID(%s) proposalTs = %s, expected %s", tt.cardID, proposalTs, tt.expectedPropTs)
		}

		if cardTs != tt.expectedCardTs {
			t.Errorf("ParseCardID(%s) cardTs = %s, expected %s", tt.cardID, cardTs, tt.expectedCardTs)
		}
	}
}

func TestIsSubTaskID(t *testing.T) {
	tests := []struct {
		cardID   string
		expected bool
	}{
		{"TASK-abc-i-123-a", true},
		{"TASK-abc-i-123-b", true},
		{"TASK-abc-i-123", false},
		{"REQ-abc-123", false},
		{"TASK-abc-i-123-ab", false},
	}

	for _, tt := range tests {
		result := IsSubTaskID(tt.cardID)
		if result != tt.expected {
			t.Errorf("IsSubTaskID(%s) = %v, expected %v", tt.cardID, result, tt.expected)
		}
	}
}

func TestGetParentTaskID(t *testing.T) {
	tests := []struct {
		subTaskID string
		expected  string
		hasError  bool
	}{
		{"TASK-abc-i-123-a", "TASK-abc-i-123", false},
		{"TASK-abc-i-123-b", "TASK-abc-i-123", false},
		{"TASK-abc-i-123", "", true},
		{"REQ-abc-123", "", true},
	}

	for _, tt := range tests {
		parentID, err := GetParentTaskID(tt.subTaskID)

		if tt.hasError {
			if err == nil {
				t.Errorf("GetParentTaskID(%s) expected error, got %s", tt.subTaskID, parentID)
			}
			continue
		}

		if err != nil {
			t.Errorf("GetParentTaskID(%s) unexpected error: %v", tt.subTaskID, err)
			continue
		}

		if parentID != tt.expected {
			t.Errorf("GetParentTaskID(%s) = %s, expected %s", tt.subTaskID, parentID, tt.expected)
		}
	}
}
