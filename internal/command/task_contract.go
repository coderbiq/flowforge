package command

import (
	"strings"

	"flowforge/internal/core"
)

func isAnalysisTask(card *core.Card) bool {
	return card != nil && taskKindFromID(card.ID) == "a"
}

func requiredTaskSections(card *core.Card) []string {
	if isAnalysisTask(card) {
		return []string{"Goal", "Inputs", "Investigation Plan", "Expected Outputs", "Done When"}
	}
	return []string{"Goal", "Inputs", "Deliverables", "Acceptance", "Out of Scope", "Read Before Work"}
}

func hasRequiredSections(body string, sections []string) bool {
	for _, section := range sections {
		if strings.TrimSpace(extractSection(body, section)) == "" {
			return false
		}
	}
	return true
}
