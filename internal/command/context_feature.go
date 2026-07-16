package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newContextFeatureCmd() *cobra.Command {
	var (
		featureID string
		stepN     int
	)

	cmd := &cobra.Command{
		Use:   "feature",
		Short: "Show FEATURE execution context",
		Long: `Show a FEATURE card's context, optionally scoped to a specific step.

Without --step: full feature context for design review.
With --step <n>: minimal context bundle for executing that step.

Examples:
  flowforge context feature --feature FEAT-001
  flowforge context feature --feature FEAT-001 --step 3
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if featureID == "" {
				return fmt.Errorf("--feature is required")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			card, err := store.ReadCard(featureID)
			if err != nil {
				return err
			}
			if card.Type != core.CardTypeFeature {
				return fmt.Errorf("card %s is not a feature card (type: %s)", featureID, card.Type)
			}

			out := cmd.OutOrStdout()

			if card.Role == "container" {
				if stepN > 0 {
					return fmt.Errorf("container feature %s has no Implementation Plan; use child feature IDs for --step", featureID)
				}
				return renderContainerFeatureContext(out, store, card)
			}

			if stepN > 0 {
				return renderStepContext(out, store, card, stepN)
			}
			return renderFullFeatureContext(out, store, card)
		},
	}

	cmd.Flags().StringVar(&featureID, "feature", "", "FEATURE card ID")
	cmd.Flags().IntVar(&stepN, "step", 0, "Step number for execution context")
	return cmd
}

func renderContainerFeatureContext(out interface{ Write([]byte) (int, error) }, store *core.CardStore, card *core.Card) error {
	w := out
	fmt.Fprintf(w, "## Container Feature: %s\n\n", card.ID)
	fmt.Fprintf(w, "- Title: %s\n", card.Title)
	fmt.Fprintf(w, "- Stage: %s\n", card.Status)

	summary := firstParagraph(card.Body)
	if summary != "" {
		fmt.Fprintf(w, "- Summary: %s\n", summary)
	}

	fmt.Fprintln(w, "\n### Sub-Features")
	for _, link := range card.Links {
		if link.Relation == "decomposes" {
			child, err := store.ReadCard(link.Target)
			if err != nil {
				fmt.Fprintf(w, "- %s (unreadable)\n", link.Target)
				continue
			}
			fmt.Fprintf(w, "- [%s] %s (stage: %s)\n", child.ID, child.Title, child.Status)
		}
	}

	fmt.Fprintln(w, "\n### Design Summary")
	designSection := extractSection(card.Body, "Design")
	if designSection != "" && !isPlaceholder(designSection) {
		kdSection := extractSubSection(designSection, "Key Decisions")
		if kdSection != "" {
			fmt.Fprintln(w, kdSection)
		}
	}

	fmt.Fprintln(w, "\n### Constraints")
	constraintsSection := extractSection(card.Body, "Constraints")
	if constraintsSection != "" {
		fmt.Fprintln(w, constraintsSection)
	}

	return nil
}

func renderFullFeatureContext(out interface{ Write([]byte) (int, error) }, store *core.CardStore, card *core.Card) error {
	w := out
	fmt.Fprintf(w, "## Feature Context: %s\n\n", card.ID)
	fmt.Fprintf(w, "- Title: %s\n", card.Title)
	fmt.Fprintf(w, "- Stage: %s\n", card.Status)

	summary := firstParagraph(card.Body)
	if summary != "" {
		fmt.Fprintf(w, "- Summary: %s\n", summary)
	}

	fmt.Fprintln(w, "\n### Linked Library Cards")
	linkedLib := linkedLibraryCards(store, card)
	if len(linkedLib) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		fmt.Fprintln(w, "| ID | Type | Title | Relation |")
		fmt.Fprintln(w, "|----|------|-------|----------|")
		for _, lc := range linkedLib {
			rel := linkRelation(card, lc.ID)
			fmt.Fprintf(w, "| %s | %s | %s | %s |\n", lc.ID, lc.Type, lc.Title, rel)
		}
	}

	fmt.Fprintln(w, "\n### Dependency Status")
	deps := featureDependencies(store, card)
	if len(deps) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		fmt.Fprintln(w, "| FEAT ID | Title | Stage | Blocks |")
		fmt.Fprintln(w, "|---------|-------|-------|--------|")
		for _, d := range deps {
			blocks := "no"
			if d.Status != core.CardStatusDone {
				blocks = "yes"
			}
			fmt.Fprintf(w, "| %s | %s | %s | %s |\n", d.ID, d.Title, d.Status, blocks)
		}
	}

	return nil
}

func renderStepContext(out interface{ Write([]byte) (int, error) }, store *core.CardStore, card *core.Card, stepN int) error {
	w := out

	ipSection := extractSection(card.Body, "Implementation Plan")
	stepBody := extractSubSection(ipSection, fmt.Sprintf("Step %d:", stepN))
	if stepBody == "" {
		return fmt.Errorf("step %d not found in Implementation Plan", stepN)
	}

	fmt.Fprintf(w, "## Step Context: %s Step %d\n\n", card.ID, stepN)

	fmt.Fprintln(w, "### Current Step")
	stepFields := parseStepFields(stepBody)
	for _, field := range []string{"Goal", "Files", "Approach", "Edge Cases", "Dependencies", "Parallel", "Verification"} {
		if val, ok := stepFields[field]; ok {
			fmt.Fprintf(w, "- **%s**: %s\n", field, val)
		}
	}

	fmt.Fprintln(w, "\n### Constraints (from FEATURE)")
	constraintsSection := extractSection(card.Body, "Constraints")
	if constraintsSection != "" && !isPlaceholder(constraintsSection) {
		fmt.Fprintln(w, constraintsSection)
	} else {
		fmt.Fprintln(w, "- None")
	}

	fmt.Fprintln(w, "\n### Relevant Library Cards")
	linkedLib := linkedLibraryCards(store, card)
	if len(linkedLib) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		fmt.Fprintln(w, "| ID | Type | Title | Relation |")
		fmt.Fprintln(w, "|----|------|-------|----------|")
		for _, lc := range linkedLib {
			rel := linkRelation(card, lc.ID)
			fmt.Fprintf(w, "| %s | %s | %s | %s |\n", lc.ID, lc.Type, lc.Title, rel)
		}
	}

	fmt.Fprintln(w, "\n### Dependency Status")
	deps := featureDependencies(store, card)
	if len(deps) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		fmt.Fprintln(w, "| FEAT ID | Title | Stage | Blocks |")
		fmt.Fprintln(w, "|---------|-------|-------|--------|")
		for _, d := range deps {
			blocks := "no"
			if d.Status != core.CardStatusDone && d.Status != core.CardStatusPlanned {
				blocks = "yes (wait strategy: check step Dependencies field)"
			}
			fmt.Fprintf(w, "| %s | %s | %s | %s |\n", d.ID, d.Title, d.Status, blocks)
		}
	}

	return nil
}

func linkedLibraryCards(store *core.CardStore, card *core.Card) []*core.Card {
	var result []*core.Card
	for _, link := range card.Links {
		target, err := store.ReadCard(link.Target)
		if err != nil {
			continue
		}
		switch target.Type {
		case core.CardTypeConvention, core.CardTypeDecision, core.CardTypeModule, core.CardTypeFinding:
			result = append(result, target)
		}
	}
	return result
}

func featureDependencies(store *core.CardStore, card *core.Card) []*core.Card {
	var result []*core.Card
	for _, link := range card.Links {
		if link.Relation != "depends_on" {
			continue
		}
		target, err := store.ReadCard(link.Target)
		if err != nil {
			continue
		}
		if target.Type == core.CardTypeFeature {
			result = append(result, target)
		}
	}
	return result
}

func linkRelation(card *core.Card, targetID string) string {
	for _, link := range card.Links {
		if link.Target == targetID {
			return link.Relation
		}
	}
	return "unknown"
}

var stepFieldLineRe = regexp.MustCompile(`^- \*\*(\w+(?:\s+\w+)*)\*\*: (.+)`)

func parseStepFields(stepBody string) map[string]string {
	fields := make(map[string]string)
	lines := strings.Split(stepBody, "\n")
	for _, line := range lines {
		matches := stepFieldLineRe.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) == 3 {
			fields[matches[1]] = strings.TrimSpace(matches[2])
		}
	}
	return fields
}
