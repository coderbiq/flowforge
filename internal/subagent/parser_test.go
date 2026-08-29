package subagent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseValidDefinition(t *testing.T) {
	dir := testAssetsDir(t)
	path := filepath.Join(dir, "flowforge-analyst.md")
	def, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if def.Name != "flowforge-analyst" {
		t.Errorf("expected name 'flowforge-analyst', got %q", def.Name)
	}
	if def.DefaultSkill == "" {
		t.Error("expected non-empty DefaultSkill")
	}
	if def.Body == "" {
		t.Error("expected non-empty Body")
	}
}

func TestParseMissingFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "no-frontmatter.md")
	content := `# This file has no frontmatter
Just body content.`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for missing frontmatter, got nil")
	}
}
