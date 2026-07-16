package core

import (
	"regexp"
	"strings"
)

var (
	frontmatterRe  = regexp.MustCompile(`(?s)^---\n.*?\n---\n`)
	headingRe      = regexp.MustCompile(`(?m)^#{1,6}\s+.*$`)
	autoNavSections = []string{"Links", "Outgoing", "FlowForge Navigation"}
)

func EffectiveContentLines(body string) int {
	cleaned := stripFrontmatter(body)
	cleaned = stripAutoNavSections(cleaned)
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return 0
	}

	lines := strings.Split(cleaned, "\n")
	count := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if headingRe.MatchString(trimmed) {
			continue
		}
		count++
	}
	return count
}

func stripFrontmatter(body string) string {
	return frontmatterRe.ReplaceAllString(body, "")
}

func stripAutoNavSections(body string) string {
	for _, section := range autoNavSections {
		body = stripMarkdownSection(body, section)
	}
	return body
}

func stripMarkdownSection(body string, section string) string {
	heading := "## " + section
	for {
		idx := strings.Index(body, heading)
		if idx < 0 {
			return body
		}
		end := strings.Index(body[idx+len(heading):], "\n## ")
		if end < 0 {
			return body[:idx]
		}
		body = body[:idx] + body[idx+len(heading)+end:]
	}
}
