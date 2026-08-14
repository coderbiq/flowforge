package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"flowforge/internal/core"
)

func markdownLinkTarget(fromCard, toCard *core.Card) (string, error) {
	if fromCard.FilePath == "" || toCard.FilePath == "" {
		return "", fmt.Errorf("cannot render markdown link without file paths")
	}
	rel, err := filepath.Rel(filepath.Dir(fromCard.FilePath), toCard.FilePath)
	if err != nil {
		return "", fmt.Errorf("computing relative link: %w", err)
	}
	return filepath.ToSlash(rel), nil
}

func upsertMarkdownSection(body string, section string, content string) string {
	heading := "## " + section
	replacement := heading + "\n\n" + strings.TrimSpace(content)
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return replacement + "\n"
	}

	idx := strings.Index(trimmed, heading)
	if idx >= 0 {
		before := strings.TrimRight(trimmed[:idx], "\n")
		afterStart := idx + len(heading)
		after := trimmed[afterStart:]
		next := strings.Index(after, "\n## ")
		if next >= 0 {
			after = after[next:]
		} else {
			after = ""
		}
		parts := []string{}
		if before != "" {
			parts = append(parts, before)
		}
		parts = append(parts, replacement)
		if strings.TrimSpace(after) != "" {
			parts = append(parts, strings.TrimLeft(after, "\n"))
		}
		return strings.Join(parts, "\n\n") + "\n"
	}

	openQuestions := "\n## Open Questions"
	if idx := strings.Index(trimmed, openQuestions); idx >= 0 {
		before := strings.TrimRight(trimmed[:idx], "\n")
		after := strings.TrimLeft(trimmed[idx:], "\n")
		return before + "\n\n" + replacement + "\n\n" + after + "\n"
	}

	return trimmed + "\n\n" + replacement + "\n"
}
