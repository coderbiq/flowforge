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

func TestGenerateCardTimestampSupportsRapidCreation(t *testing.T) {
	first := GenerateCardTimestamp()
	second := GenerateCardTimestamp()
	if first == second {
		t.Fatalf("expected rapid timestamp generation to be unique, got %s", first)
	}
}

func TestGenerateCardID(t *testing.T) {
	tests := []struct {
		cardType   CardType
		proposalTs string
		prefix     string
	}{
		{CardTypeDecision, "def456", "DEC-def456-"},
		{CardTypeFeature, "abc123", "FEAT-abc123-"},
		{CardTypeProposal, "", "PROP-"},
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
		{"Title\u115Fwith\u3164invisible\u1160chars", "titlewithinvisiblechars"},
		{"\u200Bzero\u200Cwidth", "zerowidth"},
		{"\uFEFFbom\u00A0nbs\u0301combining", "bomnbscombining"},
		{"", ""},
		{"   ", ""},
	}
	for _, tt := range tests {
		if result := ToSlug(tt.input); result != tt.expected {
			t.Errorf("ToSlug(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestToSlugMaxLength(t *testing.T) {
	if result := ToSlug(strings.Repeat("a", 100)); len(result) > 50 {
		t.Errorf("expected slug length <= 50, got %d", len(result))
	}
	result := ToSlug(strings.Repeat("中", 60))
	if len([]rune(result)) > 50 || result == "" {
		t.Errorf("expected non-empty slug with <= 50 runes, got %q", result)
	}
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		id, title, expected string
	}{
		{"DEC-456", "Use PostgreSQL", "DEC-456_use-postgresql.md"},
		{"FEAT-789", "Feature", "FEAT-789_feature.md"},
	}
	for _, tt := range tests {
		if filename := GenerateFilename(tt.id, tt.title); filename != tt.expected {
			t.Errorf("GenerateFilename(%s, %s) = %s, expected %s", tt.id, tt.title, filename, tt.expected)
		}
	}
}

func TestParseFilename(t *testing.T) {
	tests := []struct {
		filename, expectedID, expectedSlug string
		hasError                           bool
	}{
		// Accepted: current-v3 global and proposal-scoped IDs.
		{"DEC-1ab_use-postgresql.md", "DEC-1ab", "use-postgresql", false},
		{"FEAT-CR26081401-003_runtime-boundary.md", "FEAT-CR26081401-003", "runtime-boundary", false},
		// Rejected: legacy TASK/subtask IDs, missing separators, and non-v3 IDs.
		{"TASK-CR26081401-001_old-task.md", "", "", true},
		{"SUBTASK-CR26081401-003-1_subtask.md", "", "", true},
		{"FEAT-CR26081401-003.md", "", "", true},
		{"OLD-456_use-postgresql.md", "", "", true},
		{"REQ-CR26081401-001_requirement.md", "", "", true},
		{"invalid-filename.md", "", "", true},
		{"FEAT-CR26081401-003_no-extension", "", "", true},
		{"FEAT-CR26081401-003_.md", "", "", true},
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
		if id != tt.expectedID || slug != tt.expectedSlug {
			t.Errorf("ParseFilename(%s) = (%s, %s), expected (%s, %s)", tt.filename, id, slug, tt.expectedID, tt.expectedSlug)
		}
	}
}
