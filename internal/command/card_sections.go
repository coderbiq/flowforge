package command

import "strings"

func extractSection(body, section string) string {
	if strings.Contains(section, ".") {
		return extractHierarchicalSection(body, section)
	}

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

func extractHierarchicalSection(body, path string) string {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) != 2 {
		return ""
	}
	parentSection := extractSection(body, parts[0])
	if parentSection == "" {
		return ""
	}

	childHeading := "### " + parts[1]
	idx := strings.Index(parentSection, childHeading)
	if idx < 0 {
		return ""
	}
	childBody := parentSection[idx+len(childHeading):]
	nextH3 := strings.Index(childBody, "\n### ")
	nextH2 := strings.Index(childBody, "\n## ")
	cutAt := len(childBody)
	if nextH3 >= 0 && nextH3 < cutAt {
		cutAt = nextH3
	}
	if nextH2 >= 0 && nextH2 < cutAt {
		cutAt = nextH2
	}
	return strings.TrimSpace(childBody[:cutAt])
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
