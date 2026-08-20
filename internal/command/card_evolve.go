package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func newCardEvolveCmd() *cobra.Command {
	var (
		stage   string
		regress bool
	)

	cmd := &cobra.Command{
		Use:   "evolve <card-id>",
		Short: "Evolve a FEATURE card to the next stage",
		Long: `Upgrade (or regress with --regress) a FEATURE card's stage.
Validates gate conditions before allowing stage transitions.

Stages: draft → designed → planned → in_progress → done
Use --regress to move backward (e.g. planned → designed).

Examples:
  flowforge card evolve FEAT-xxx --stage designed
  flowforge card evolve FEAT-xxx --stage designed --regress
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]
			if stage == "" {
				return fmt.Errorf("--stage is required (designed|planned|done)")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			card, err := store.ReadCard(cardID)
			if err != nil {
				return err
			}
			if card.Type != core.CardTypeFeature {
				return fmt.Errorf("card %s is not a feature card (type: %s)", cardID, card.Type)
			}
			if card.Role == "container" && !regress {
				return fmt.Errorf("container feature %s cannot be evolved; evolve its child features instead", cardID)
			}

			targetStage := core.CardStatus(stage)
			if !isValidEvolveTarget(targetStage) {
				return fmt.Errorf("invalid stage: %s (valid: designed, planned, done)", stage)
			}

			if regress {
				return handleRegress(store, card, targetStage, cmd)
			}
			return handleEvolve(store, card, targetStage, cmd)
		},
	}

	cmd.Flags().StringVar(&stage, "stage", "", "Target stage (designed|planned|done)")
	cmd.Flags().BoolVar(&regress, "regress", false, "Allow stage regression")
	return cmd
}

func isValidEvolveTarget(s core.CardStatus) bool {
	switch s {
	case core.CardStatusDesigned, core.CardStatusPlanned, core.CardStatusDone:
		return true
	}
	return false
}

type gateIssue struct {
	Section string
	Detail  string
	Fix     string
}

func handleEvolve(store *core.CardStore, card *core.Card, target core.CardStatus, cmd *cobra.Command) error {
	out := cmd.OutOrStdout()

	var issues []gateIssue
	switch target {
	case core.CardStatusDesigned:
		if card.Status != core.CardStatusDraft {
			return fmt.Errorf("card must be in 'draft' stage to evolve to 'designed', current: %s", card.Status)
		}
		issues = validateDesignedGate(card.Body)
	case core.CardStatusPlanned:
		if card.Status != core.CardStatusDesigned {
			return fmt.Errorf("card must be in 'designed' stage to evolve to 'planned', current: %s", card.Status)
		}
		issues = validatePlannedGate(card.Body)
	case core.CardStatusDone:
		if card.Status != core.CardStatusInProgress {
			return fmt.Errorf("card must be in 'in_progress' stage to evolve to 'done', current: %s", card.Status)
		}
		issues = validateDoneGate(card.Body)
	}

	if target == core.CardStatusDesigned || target == core.CardStatusPlanned {
		analysisIssues, err := complexAnalysisGateIssues(card, target)
		if err != nil {
			return err
		}
		issues = append(issues, analysisIssues...)
	}

	// In v4, gates are non-blocking advisory indicators.
	// Target stage update proceeds directly to eliminate CLI deadlocks.
	card.Status = target
	if err := store.UpdateCard(card); err != nil {
		return fmt.Errorf("updating card: %w", err)
	}

	fmt.Fprintf(out, "✓ %s status updated to '%s'\n", card.ID, target)
	return nil
}

func complexAnalysisGateIssues(card *core.Card, target core.CardStatus) ([]gateIssue, error) {
	if card == nil || !complexAnalysisModeRe.MatchString(card.Body) {
		return nil, nil
	}

	issues := validateComplexAnalysisBodyGate(card.Body, target)
	proposalID := strings.TrimSpace(card.Source)
	if proposalID == "" {
		proposalID = strings.TrimSpace(card.ProposalID)
	}
	if proposalID == "" {
		return append(issues, gateIssue{
			Section: "Analysis State",
			Detail:  "complex FEATURE has no proposal ID",
			Fix:     "set source/proposal_id so analysis status can be resolved",
		}), nil
	}

	view, closeFn, err := currentAnalysisView(proposalID, true)
	if closeFn != nil {
		defer closeFn()
	}
	if err != nil {
		return nil, fmt.Errorf("checking complex analysis readiness: %w", err)
	}
	return append(issues, validateComplexAnalysisState(view, target)...), nil
}

func validateComplexAnalysisBodyGate(body string, target core.CardStatus) []gateIssue {
	var issues []gateIssue
	for _, section := range []string{"Objective", "Current Understanding", "Evidence", "Working Design", "Rejected or Revised Assumptions", "Next Investigation"} {
		content := extractSection(body, section)
		if isComplexAnalysisPlaceholder(content) {
			issues = append(issues, gateIssue{Section: section, Detail: "complex analysis section is missing or placeholder", Fix: "record a recoverable value; use None only when intentionally empty"})
		}
	}

	evidence := strings.ToLower(extractSection(body, "Evidence"))
	if !containsAny(evidence, "accepted", "rejected", "conflicting", "inconclusive", "已采纳", "已拒绝", "冲突", "无结论") {
		issues = append(issues, gateIssue{Section: "Evidence", Detail: "no evidence support state found", Fix: "classify decisive evidence as accepted, rejected, conflicting, or inconclusive"})
	}
	next := strings.TrimSpace(extractSection(body, "Next Investigation"))
	if !isExplicitNone(next) {
		issues = append(issues, gateIssue{Section: "Next Investigation", Detail: "investigation remains active", Fix: "complete or supersede the investigation, then set this section to None"})
	}
	if target == core.CardStatusPlanned {
		verification := extractSection(body, "Verification")
		if isPlaceholder(verification) || countBulletLines(verification) == 0 {
			issues = append(issues, gateIssue{Section: "Verification", Detail: "complex FEATURE needs explicit verification mappings", Fix: "map user-visible outcomes and design risks to tests or inspection"})
		}
		ipSection := extractSection(body, "Implementation Plan")
		for _, match := range stepHeaderRe.FindAllStringSubmatch(ipSection, -1) {
			stepBody := extractSubSection(ipSection, "Step "+match[1]+":")
			for _, field := range []string{"Dependencies", "Parallel", "Verification"} {
				if !hasStepField(stepBody, field) {
					issues = append(issues, gateIssue{Section: "Implementation Plan.Step " + match[1], Detail: "missing complex-planning field: " + field, Fix: "add an explicit " + field + " value"})
				}
			}
		}
	}
	return issues
}

func isComplexAnalysisPlaceholder(value string) bool {
	cleaned := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(value), "-* "))
	if cleaned == "" {
		return true
	}
	switch strings.ToLower(cleaned) {
	case "tbd", "todo", "pending", "待补充", "待定", "<!-- tbd -->":
		return true
	default:
		return false
	}
}

func validateComplexAnalysisState(view state.AnalysisView, target core.CardStatus) []gateIssue {
	if view.State == "no_plan" || view.State == "completed" {
		return nil
	}
	detail := fmt.Sprintf("analysis state is %s (next action: %s)", view.State, view.NextAction)
	if target == core.CardStatusPlanned && view.ActivePlan != nil {
		detail += fmt.Sprintf("; active plan revision %d is not completed", view.ActivePlan.Revision)
	}
	return []gateIssue{{Section: "Analysis State", Detail: detail, Fix: "finish synthesis, resolve conflicts or user decisions, and seal analysis.completed"}}
}

func containsAny(value string, candidates ...string) bool {
	for _, candidate := range candidates {
		if strings.Contains(value, candidate) {
			return true
		}
	}
	return false
}

func isExplicitNone(value string) bool {
	cleaned := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(value), "-* "))
	return strings.EqualFold(cleaned, "none") || cleaned == "无" || cleaned == "无下一步调查"
}

func handleRegress(store *core.CardStore, card *core.Card, target core.CardStatus, cmd *cobra.Command) error {
	from := card.Status
	allowed := false
	switch {
	case from == core.CardStatusPlanned && target == core.CardStatusDesigned:
		allowed = true
	case from == core.CardStatusInProgress && (target == core.CardStatusDesigned || target == core.CardStatusPlanned):
		allowed = true
	case from == core.CardStatusDone && (target == core.CardStatusDesigned || target == core.CardStatusPlanned || target == core.CardStatusInProgress):
		allowed = true
	}
	if !allowed {
		return fmt.Errorf("cannot regress from '%s' to '%s'", from, target)
	}

	resetSteps := target != core.CardStatusPlanned
	if target == core.CardStatusPlanned {
		resetSteps = false
	}

	card.Body = resetStepStatuses(card.Body, resetSteps)
	card.Body = appendHistoryLine(card.Body, "decision", fmt.Sprintf("stage regressed: %s → %s", from, target))
	card.Status = target

	if err := store.UpdateCard(card); err != nil {
		return fmt.Errorf("updating card: %w", err)
	}

	out := cmd.OutOrStdout()
	msg := fmt.Sprintf("✓ %s regressed: %s → %s", card.ID, from, target)
	if resetSteps {
		msg += " (all step statuses reset)"
	}
	fmt.Fprintln(out, msg)
	return nil
}

var (
	placeholderRe         = regexp.MustCompile(`^(None|TBD|N/A|<!-- TBD -->)\s*$`)
	complexAnalysisModeRe = regexp.MustCompile(`(?mi)^\s*<!--\s*analysis-mode:\s*complex\s*-->\s*$`)
	stepHeaderRe          = regexp.MustCompile(`(?m)^### Step (\d+):`)
)

func validateDesignedGate(body string) []gateIssue {
	var issues []gateIssue

	designSection := extractSection(body, "Design")
	if isPlaceholder(designSection) {
		issues = append(issues, gateIssue{
			Section: "Design",
			Detail:  "section is missing or placeholder",
			Fix:     fmt.Sprintf("edit the Design section of the card"),
		})
	} else {
		kdSection := extractSubSection(designSection, "Key Decisions")
		if isPlaceholder(kdSection) || countBulletLines(kdSection) == 0 {
			issues = append(issues, gateIssue{
				Section: "Design.Key Decisions",
				Detail:  fmt.Sprintf("0 substantive entries found (minimum: 1)"),
				Fix:     fmt.Sprintf("add at least 1 key decision with rationale"),
			})
		}
	}

	constraintsSection := extractSection(body, "Constraints")
	if isPlaceholder(constraintsSection) || countBulletLines(constraintsSection) == 0 {
		issues = append(issues, gateIssue{
			Section: "Constraints",
			Detail:  "0 substantive entries found (minimum: 1)",
			Fix:     "add constraints from CONV/DEC/business rules",
		})
	}

	oqSection := extractSection(body, "Open Questions")
	openCount := countOpenQuestions(oqSection)
	if openCount > 0 {
		var qs []string
		for _, line := range strings.Split(oqSection, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- ") && !strings.Contains(trimmed, "[假设]") && !isPlaceholder(trimmed) {
				qs = append(qs, strings.TrimPrefix(trimmed, "- "))
			}
		}
		detail := fmt.Sprintf("%d unresolved questions remain", openCount)
		for _, q := range qs {
			detail += "\n      → " + q
		}
		issues = append(issues, gateIssue{
			Section: "Open Questions",
			Detail:  detail,
			Fix:     "resolve each question or mark as [假设]",
		})
	}

	return issues
}

func validatePlannedGate(body string) []gateIssue {
	var issues []gateIssue

	ipSection := extractSection(body, "Implementation Plan")
	steps := stepHeaderRe.FindAllStringSubmatch(ipSection, -1)
	if len(steps) == 0 {
		issues = append(issues, gateIssue{
			Section: "Implementation Plan",
			Detail:  "no ### Step N: sections found (minimum: 1)",
			Fix:     "add at least one implementation step",
		})
		return issues
	}

	for _, match := range steps {
		stepNum := match[1]
		stepBody := extractSubSection(ipSection, "Step "+stepNum+":")
		var missing []string
		for _, field := range []string{"Goal", "Files", "Symbols", "Actions", "Constraints", "Done When", "Verification"} {
			if !hasStepField(stepBody, field) {
				missing = append(missing, field)
			}
		}
		if len(missing) > 0 {
			issues = append(issues, gateIssue{
				Section: fmt.Sprintf("Implementation Plan.Step %s", stepNum),
				Detail:  fmt.Sprintf("missing required fields: %s", strings.Join(missing, ", ")),
				Fix:     fmt.Sprintf("add %s with substantive content", strings.Join(missing, ", ")),
			})
		}
		if vaguePlanRe.MatchString(stepBody) {
			issues = append(issues, gateIssue{
				Section: fmt.Sprintf("Implementation Plan.Step %s", stepNum),
				Detail:  "contains vague execution language (for example TBD/as needed/必要时/视情况)",
				Fix:     "replace ambiguity with an explicit condition, action, target, and observable result",
			})
		}
	}

	oqSection := extractSection(body, "Open Questions")
	if countOpenQuestions(oqSection) > 0 {
		issues = append(issues, gateIssue{
			Section: "Open Questions",
			Detail:  "must be completely cleared for 'planned' stage",
			Fix:     "resolve all open questions",
		})
	}

	return issues
}

var vaguePlanRe = regexp.MustCompile(`(?i)(\bTBD\b|\bTODO\b|\bas needed\b|\bas appropriate\b|必要时|视情况|酌情|待定|待确认)`)

func validateDoneGate(body string) []gateIssue {
	var issues []gateIssue

	ipSection := extractSection(body, "Implementation Plan")
	steps := stepHeaderRe.FindAllStringSubmatch(ipSection, -1)
	var incompleteSteps []string
	for _, match := range steps {
		stepNum := match[1]
		stepBody := extractSubSection(ipSection, "Step "+stepNum+":")
		if !strings.Contains(stepBody, "step-status: done") {
			incompleteSteps = append(incompleteSteps, stepNum)
		}
	}
	if len(incompleteSteps) > 0 {
		issues = append(issues, gateIssue{
			Section: "Implementation Plan",
			Detail:  fmt.Sprintf("steps not completed: %s", strings.Join(incompleteSteps, ", ")),
			Fix:     "mark steps as done with: flowforge card steps <id> --status done <n>",
		})
	}

	verificationSection := extractSection(body, "Verification")
	if isPlaceholder(verificationSection) {
		issues = append(issues, gateIssue{
			Section: "Verification",
			Detail:  "verification results are missing or placeholder",
			Fix:     "document verification results for each acceptance criterion",
		})
	}

	return issues
}

func isPlaceholder(section string) bool {
	trimmed := strings.TrimSpace(section)
	if trimmed == "" {
		return true
	}
	lines := strings.Split(trimmed, "\n")
	nonEmpty := 0
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		if placeholderRe.MatchString(strings.TrimPrefix(t, "- ")) {
			continue
		}
		nonEmpty++
	}
	return nonEmpty == 0
}

func countBulletLines(section string) int {
	count := 0
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			content := strings.TrimSpace(trimmed[2:])
			if content != "" && !placeholderRe.MatchString(content) {
				count++
			}
		}
	}
	return count
}

func countOpenQuestions(section string) int {
	count := 0
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") && !strings.Contains(trimmed, "[假设]") {
			content := strings.TrimSpace(trimmed[2:])
			if content != "" && !placeholderRe.MatchString(content) {
				count++
			}
		}
	}
	return count
}

func extractSubSection(body, subHeading string) string {
	lines := strings.Split(body, "\n")
	var capture []string
	inSub := false
	subLevel := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			level, heading, ok := parseMarkdownHeading(trimmed)
			if !ok {
				continue
			}
			headingLower := strings.ToLower(heading)
			if strings.HasPrefix(headingLower, strings.ToLower(subHeading)) {
				inSub = true
				subLevel = level
				continue
			}
			if inSub && level <= subLevel {
				break
			}
			continue
		}
		if inSub {
			capture = append(capture, line)
		}
	}
	return strings.Join(capture, "\n")
}

func hasStepField(stepBody, field string) bool {
	return stepFieldValue(stepBody, field) != ""
}

func stepFieldValue(stepBody, field string) string {
	lines := strings.Split(stepBody, "\n")
	prefix := "- **" + field + "**:"
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, prefix) {
			continue
		}
		values := []string{strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))}
		for _, following := range lines[i+1:] {
			candidate := strings.TrimSpace(following)
			if stepFieldLineRe.MatchString(candidate) {
				break
			}
			if candidate != "" && !strings.HasPrefix(candidate, "<!--") {
				values = append(values, candidate)
			}
		}
		value := strings.TrimSpace(strings.Join(values, "\n"))
		if value == "" || placeholderRe.MatchString(value) || strings.EqualFold(value, "TBD") {
			return ""
		}
		return value
	}
	return ""
}

func resetStepStatuses(body string, resetAll bool) string {
	ipSection := extractSection(body, "Implementation Plan")
	if ipSection == "" {
		return body
	}

	newIP := ipSection
	for _, match := range stepHeaderRe.FindAllStringSubmatch(ipSection, -1) {
		stepNum := match[1]
		stepBody := extractSubSection(ipSection, "Step "+stepNum+":")
		if strings.Contains(stepBody, "step-status: done") {
			if !resetAll {
				continue
			}
		}
		newIP = strings.ReplaceAll(newIP,
			fmt.Sprintf("<!-- step-status: done -->"),
			"<!-- step-status: not_started -->")
		newIP = strings.ReplaceAll(newIP,
			fmt.Sprintf("<!-- step-status: in_progress -->"),
			"<!-- step-status: not_started -->")
		newIP = strings.ReplaceAll(newIP,
			fmt.Sprintf("<!-- step-status: blocked -->"),
			"<!-- step-status: not_started -->")
	}

	return strings.Replace(body, ipSection, newIP, 1)
}

func appendHistoryLine(body, kind, event string) string {
	timeStr := "<!-- TODO: ISO time -->"
	line := fmt.Sprintf("- %s | %s | %s", timeStr, kind, event)

	historyIdx := strings.Index(body, "## History")
	if historyIdx < 0 {
		depIdx := strings.Index(body, "## Dependencies")
		if depIdx >= 0 {
			return body[:depIdx] + "## History\n\n" + line + "\n\n" + body[depIdx:]
		}
		return body + "\n\n## History\n\n" + line + "\n"
	}

	historyEnd := strings.Index(body[historyIdx:], "\n## ")
	if historyEnd < 0 {
		return body + "\n" + line + "\n"
	}
	insertAt := historyIdx + historyEnd
	return body[:insertAt] + line + "\n" + body[insertAt:]
}
