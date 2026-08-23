package tracker

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	titleRegex     = regexp.MustCompile(`^#\s*(?:(\d+)\s*[:：\-]\s*)?(.*)$`)
	statusRegex    = regexp.MustCompile(`(?i)^\*\*Status:\*\*\s*(.+)$`)
	typeRegex      = regexp.MustCompile(`(?i)^\*\*Type:\*\*\s*(.+)$`)
	blockedByRegex = regexp.MustCompile(`(?i)^\*\*Blocked\s+by:\*\*\s*(.+)$`)
	assigneeRegex  = regexp.MustCompile(`(?i)^\*\*Assignee:\*\*\s*(.+)$`)
	labelsRegex    = regexp.MustCompile(`(?i)^\*\*Labels:\*\*\s*(.+)$`)
	idPrefixRegex  = regexp.MustCompile(`^(\d+)(?:-(.*))?\.md$`)
)

// ParseIssueFile reads and parses a single issue markdown file.
func ParseIssueFile(filePath string) (*Issue, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read issue file: %w", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("stat issue file: %w", err)
	}

	baseName := filepath.Base(filePath)
	featureDir := filepath.Base(filepath.Dir(filepath.Dir(filePath)))
	if featureDir == "." || featureDir == "/" {
		featureDir = ""
	}

	id := ""
	slug := strings.TrimSuffix(baseName, ".md")
	if m := idPrefixRegex.FindStringSubmatch(baseName); len(m) > 1 {
		id = m[1]
		if len(m) > 2 && m[2] != "" {
			slug = m[2]
		}
	}

	issue := &Issue{
		ID:        id,
		Slug:      slug,
		FilePath:  filePath,
		Feature:   featureDir,
		Status:    StatusOpen,
		Type:      "task",
		BlockedBy: []string{},
		Labels:    []string{},
		ModTime:   info.ModTime(),
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var bodyLines []string
	inHeader := true

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if inHeader {
			if issue.Title == "" && strings.HasPrefix(trimmed, "#") {
				if m := titleRegex.FindStringSubmatch(trimmed); len(m) > 0 {
					if m[1] != "" && issue.ID == "" {
						issue.ID = m[1]
					}
					issue.Title = strings.TrimSpace(m[2])
				}
				continue
			}

			if m := statusRegex.FindStringSubmatch(trimmed); len(m) > 1 {
				issue.Status = Status(strings.ToLower(strings.TrimSpace(m[1])))
				continue
			}

			if m := typeRegex.FindStringSubmatch(trimmed); len(m) > 1 {
				issue.Type = strings.ToLower(strings.TrimSpace(m[1]))
				continue
			}

			if m := blockedByRegex.FindStringSubmatch(trimmed); len(m) > 1 {
				raw := strings.TrimSpace(m[1])
				if raw != "" && raw != "none" && raw != "None" {
					parts := strings.Split(raw, ",")
					for _, p := range parts {
						p = strings.TrimSpace(p)
						if p != "" {
							issue.BlockedBy = append(issue.BlockedBy, p)
						}
					}
				}
				continue
			}

			if m := assigneeRegex.FindStringSubmatch(trimmed); len(m) > 1 {
				issue.Assignee = strings.TrimSpace(m[1])
				continue
			}

			if m := labelsRegex.FindStringSubmatch(trimmed); len(m) > 1 {
				raw := strings.TrimSpace(m[1])
				parts := strings.Split(raw, ",")
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if p != "" {
						issue.Labels = append(issue.Labels, p)
					}
				}
				continue
			}

			// If we hit an empty line or a non-metadata line, transition to body
			if trimmed != "" && !strings.HasPrefix(trimmed, "**") {
				inHeader = false
				bodyLines = append(bodyLines, line)
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	issue.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
	return issue, nil
}

// DiscoverIssues scans a root directory (e.g. .scratch/) for all issues.
func DiscoverIssues(root string) ([]*Issue, error) {
	var issues []*Issue

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return issues, nil
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Match any .md inside an issues/ directory or named spec.md
		if strings.HasSuffix(path, ".md") && (strings.Contains(path, "issues/") || strings.HasSuffix(path, "spec.md")) {
			issue, parseErr := ParseIssueFile(path)
			if parseErr == nil && issue != nil {
				issues = append(issues, issue)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk issues: %w", err)
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Feature != issues[j].Feature {
			return issues[i].Feature < issues[j].Feature
		}
		return issues[i].ID < issues[j].ID
	})

	return issues, nil
}
