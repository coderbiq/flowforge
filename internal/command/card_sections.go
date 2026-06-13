package command

import "strings"

func extractSection(body, section string) string {
	marker := "## " + section
	idx := strings.Index(body, marker)
	if idx < 0 {
		return ""
	}
	sectionBody := body[idx+len(marker):]
	next := strings.Index(sectionBody, "\n## ")
	if next >= 0 {
		sectionBody = sectionBody[:next]
	}
	return strings.TrimSpace(sectionBody)
}

func splitBulletLines(section string) []string {
	var lines []string
	for _, raw := range strings.Split(section, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
