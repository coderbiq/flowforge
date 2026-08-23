package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WorkspaceProposal represents a proposal in Tier 2 working memory (01-workspace/<proposal_id>).
type WorkspaceProposal struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Path      string `json:"path"`      // Absolute path to proposal directory
	Readme    string `json:"readme"`    // Absolute path to README.md
	Mode      string `json:"mode"`      // flat or hierarchical
	Status    string `json:"status"`    // active or completed
	CreatedAt string `json:"createdAt"`
}

// WorkspaceSlice represents a single tracer bullet or batch slice parsed from Scratchpad.
type WorkspaceSlice struct {
	Index       string   `json:"index"`       // e.g. "1", "2.1", "Slice 1"
	Title       string   `json:"title"`       // Title or summary
	Objective   string   `json:"objective"`   // Objective / description
	TouchPoints []string `json:"touchPoints"` // Files / modules to modify
	TestCommand string   `json:"testCommand"` // Test verification command
	Completed   bool     `json:"completed"`   // Checkbox status
	Raw         string   `json:"raw"`         // Full raw section text
}

// ListWorkspaceProposals returns all proposals found under 01-workspace.
func (s *CardStore) ListWorkspaceProposals() ([]WorkspaceProposal, error) {
	workspaceDir := s.WorkspaceDir()
	entries, err := os.ReadDir(workspaceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading workspace directory %s: %w", workspaceDir, err)
	}

	var proposals []WorkspaceProposal
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirName := entry.Name()
		if strings.HasPrefix(dirName, ".") {
			continue
		}

		propDir := filepath.Join(workspaceDir, dirName)
		readmePath := filepath.Join(propDir, "README.md")
		if _, err := os.Stat(readmePath); err != nil {
			continue
		}

		proposalID := dirName

		title, status := parseReadmeMeta(readmePath, dirName)
		mode := "flat"
		modulesDir := filepath.Join(propDir, "modules")
		if info, err := os.Stat(modulesDir); err == nil && info.IsDir() {
			mode = "hierarchical"
		}

		proposals = append(proposals, WorkspaceProposal{
			ID:     proposalID,
			Title:  title,
			Path:   propDir,
			Readme: readmePath,
			Mode:   mode,
			Status: status,
		})
	}

	return proposals, nil
}

// FindWorkspaceProposal locates a proposal by ID or directory name.
func (s *CardStore) FindWorkspaceProposal(proposalID string) (*WorkspaceProposal, error) {
	proposals, err := s.ListWorkspaceProposals()
	if err != nil {
		return nil, err
	}
	for _, p := range proposals {
		if p.ID == proposalID || filepath.Base(p.Path) == proposalID || strings.HasPrefix(filepath.Base(p.Path), proposalID+"_") || strings.HasPrefix(filepath.Base(p.Path), proposalID+"-") {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("proposal %q not found in workspace", proposalID)
}

// ParseProposalSlices extracts slice definitions from 01-workspace/<proposal_id>/README.md.
func ParseProposalSlices(readmePath string) ([]WorkspaceSlice, error) {
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, fmt.Errorf("reading proposal readme: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var slices []WorkspaceSlice

	inSlicesSection := false
	var currentSlice *WorkspaceSlice
	var currentSectionLines []string

	flushCurrent := func() {
		if currentSlice != nil {
			currentSlice.Raw = strings.TrimSpace(strings.Join(currentSectionLines, "\n"))
			slices = append(slices, *currentSlice)
			currentSlice = nil
			currentSectionLines = nil
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if strings.Contains(strings.ToLower(trimmed), "slice") || strings.Contains(strings.ToLower(trimmed), "batch") || strings.Contains(strings.ToLower(trimmed), "execution plan") || strings.Contains(strings.ToLower(trimmed), "plan") {
				inSlicesSection = true
			} else if inSlicesSection {
				// Reached another major section
				flushCurrent()
				inSlicesSection = false
			}
		}

		if !inSlicesSection {
			continue
		}

		// Detect slice/batch header: "### Slice 1: ...", "### [Batch 1: Expand] ...", "### 1. ...", "### [ ] Slice 1: ...", "#### Slice 1: ...", "- [ ] **Slice 1: ...**"
		if strings.HasPrefix(trimmed, "### ") || strings.HasPrefix(trimmed, "#### ") || (strings.HasPrefix(trimmed, "- [") && (strings.Contains(strings.ToLower(trimmed), "slice") || strings.Contains(strings.ToLower(trimmed), "batch"))) {
			flushCurrent()
			headerText := strings.TrimSpace(strings.TrimLeft(trimmed, "#- "))
			completed := false
			if strings.HasPrefix(headerText, "[x]") || strings.HasPrefix(headerText, "[X]") {
				completed = true
				headerText = strings.TrimSpace(headerText[3:])
			} else if strings.HasPrefix(headerText, "[ ]") {
				headerText = strings.TrimSpace(headerText[3:])
			}
			headerText = strings.Trim(headerText, "*_` ")

			// Extract index and title
			index := headerText
			title := headerText
			if idx := strings.Index(headerText, "]"); idx > 0 && strings.HasPrefix(headerText, "[") {
				index = strings.TrimSpace(headerText[1:idx])
				title = strings.TrimSpace(headerText[idx+1:])
			} else if idx := strings.Index(headerText, ":"); idx > 0 {
				index = strings.TrimSpace(headerText[:idx])
				title = strings.TrimSpace(headerText[idx+1:])
			} else if idx := strings.Index(headerText, "."); idx > 0 {
				index = strings.TrimSpace(headerText[:idx])
				title = strings.TrimSpace(headerText[idx+1:])
			}

			currentSlice = &WorkspaceSlice{
				Index:     index,
				Title:     title,
				Completed: completed,
			}
			currentSectionLines = append(currentSectionLines, line)
			continue
		}

		if currentSlice != nil {
			currentSectionLines = append(currentSectionLines, line)
			lower := strings.ToLower(trimmed)
			if strings.Contains(lower, "test") || strings.Contains(lower, "verification") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					currentSlice.TestCommand = strings.Trim(strings.TrimSpace(parts[1]), "`* ")
				}
			} else if strings.Contains(lower, "touchpoint") || strings.Contains(lower, "seam") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					points := strings.Split(parts[1], ",")
					for _, pt := range points {
						cleaned := strings.Trim(strings.TrimSpace(pt), "`* ")
						if cleaned != "" {
							currentSlice.TouchPoints = append(currentSlice.TouchPoints, cleaned)
						}
					}
				}
			}
		}
	}

	flushCurrent()
	return slices, nil
}

func parseReadmeMeta(readmePath, defaultID string) (title string, status string) {
	status = "active"
	title = defaultID

	file, err := os.Open(readmePath)
	if err != nil {
		return title, status
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") && title == defaultID {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
		if strings.HasPrefix(strings.ToLower(line), "status:") || strings.HasPrefix(strings.ToLower(line), "- **status**:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				val := strings.ToLower(strings.TrimSpace(parts[1]))
				if strings.Contains(val, "completed") || strings.Contains(val, "archived") || strings.Contains(val, "done") {
					status = "completed"
				}
			}
		}
	}

	return title, status
}
