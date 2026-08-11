package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func newContextFeatureCmd() *cobra.Command {
	var (
		featureID string
		stepN     int
		role      string
		workID    string
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

			if role != "" {
				return renderRoleFeatureContext(out, store, card, stepN, role, workID)
			}

			if stepN > 0 {
				return renderStepContext(out, store, card, stepN)
			}
			return renderFullFeatureContext(out, store, card)
		},
	}

	cmd.Flags().StringVar(&featureID, "feature", "", "FEATURE card ID")
	cmd.Flags().IntVar(&stepN, "step", 0, "Step number for execution context")
	cmd.Flags().StringVar(&role, "role", "", "Role-specific context (design-analyst|coordinator|investigator|executor)")
	cmd.Flags().StringVar(&workID, "work", "", "Analysis work ID for investigator context")
	return cmd
}

func renderRoleFeatureContext(out interface{ Write([]byte) (int, error) }, store *core.CardStore, card *core.Card, stepN int, role, workID string) error {
	switch role {
	case "executor":
		if stepN <= 0 {
			return fmt.Errorf("--step is required for executor context")
		}
		return renderStepContext(out, store, card, stepN)
	case "design-analyst":
		if err := renderFullFeatureContext(out, store, card); err != nil {
			return err
		}
		view, closeFn, err := currentAnalysisView(card.Source, true)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return fmt.Errorf("loading analyst context: %w", err)
		}
		return renderAnalysisContext(out, view, true)
	case "coordinator":
		view, closeFn, err := currentAnalysisView(card.Source, true)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return fmt.Errorf("loading coordinator context: %w", err)
		}
		fmt.Fprintf(out, "## Coordinator Context: %s\n", card.ID)
		return renderAnalysisContext(out, view, false)
	case "investigator":
		if workID == "" {
			return fmt.Errorf("--work is required for investigator context")
		}
		view, closeFn, err := currentAnalysisView(card.Source, true)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return fmt.Errorf("loading investigator context: %w", err)
		}
		work, ok := findAnalysisWork(view, workID)
		if !ok {
			return fmt.Errorf("analysis work %s is not in the active plan", workID)
		}
		return renderInvestigatorContext(out, card, work)
	default:
		return fmt.Errorf("invalid --role %q (valid: design-analyst, coordinator, investigator, executor)", role)
	}
}

func renderAnalysisContext(out interface{ Write([]byte) (int, error) }, view state.AnalysisView, includeHistory bool) error {
	fmt.Fprintln(out, "\n### Analysis State")
	fmt.Fprintf(out, "- State: %s\n- Next Action: %s\n", view.State, view.NextAction)
	if view.ActivePlan != nil {
		fmt.Fprintf(out, "- Active Plan: %s revision %d\n", view.ActivePlan.CycleID, view.ActivePlan.Revision)
	}
	fmt.Fprintf(out, "- Ready: %d\n- Running: %d\n- Returned: %d\n- Blocked: %d\n", len(view.ReadyWork), len(view.RunningWork), len(view.ReturnedWork), len(view.BlockedWork))
	if includeHistory {
		fmt.Fprintf(out, "- Sealed Events: %d\n", len(view.History))
	}
	for _, work := range view.ReadyWork {
		fmt.Fprintf(out, "- Dispatchable: %s — %s\n", work.WorkID, work.Question)
	}
	return nil
}

func findAnalysisWork(view state.AnalysisView, workID string) (state.AnalysisWork, bool) {
	groups := [][]state.AnalysisWork{view.ReadyWork, view.RunningWork, view.ReturnedWork, view.BlockedWork}
	for _, group := range groups {
		for _, work := range group {
			if work.WorkID == workID {
				return work, true
			}
		}
	}
	return state.AnalysisWork{}, false
}

func renderInvestigatorContext(out interface{ Write([]byte) (int, error) }, card *core.Card, work state.AnalysisWork) error {
	fmt.Fprintf(out, "## Investigator Context: %s\n\n", work.WorkID)
	fmt.Fprintf(out, "- FEATURE: %s\n- Question: %s\n- Scope: %s\n- State: %s\n- Evidence Target: %s\n- Done When: %s\n", card.ID, work.Question, work.Scope, work.State, work.EvidenceTarget, work.DoneWhen)
	if len(work.Sources) > 0 {
		fmt.Fprintf(out, "- Allowed Sources: %s\n", strings.Join(work.Sources, ", "))
	}
	if len(work.Inputs) > 0 {
		fmt.Fprintf(out, "- Persisted Inputs: %s\n", strings.Join(work.Inputs, ", "))
	}
	if work.Skill != "" {
		fmt.Fprintf(out, "- Skill: %s\n", work.Skill)
	}
	if work.Budget != nil {
		fmt.Fprintf(out, "- Budget: %v\n", work.Budget)
	}
	fmt.Fprintln(out, "- Write Boundary: edit only the designated FIND Evidence, Source, Impact, and Open Questions")
	return nil
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
		if link.Relation != "requires" {
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
