package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
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

	if len(issues) > 0 {
		fmt.Fprintf(out, "Evolve to '%s' rejected — %d issues:\n\n", target, len(issues))
		for i, issue := range issues {
			fmt.Fprintf(out, "  [%d] %s: %s\n", i+1, issue.Section, issue.Detail)
			if issue.Fix != "" {
				fmt.Fprintf(out, "      → %s\n", issue.Fix)
			}
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "Commands:\n")
		fmt.Fprintf(out, "  flowforge card read %s --summary\n", card.ID)
		return nil
	}

	card.Status = target
	if err := store.UpdateCard(card); err != nil {
		return fmt.Errorf("updating card: %w", err)
	}

	fmt.Fprintf(out, "✓ %s evolved to '%s'\n", card.ID, target)
	return nil
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
	placeholderRe = regexp.MustCompile(`^(None|TBD|N/A|<!-- TBD -->)\s*$`)
	crossRefRe    = regexp.MustCompile(`参考\s*(DES|REQ|TASK|STR)-|参见.*卡片|see\s+(DES|REQ|TASK|STR)-`)
	stepHeaderRe  = regexp.MustCompile(`(?m)^### Step (\d+):`)
	stepFieldRe   = regexp.MustCompile(`(?m)^- \*\*(\w+)\*\*: (.+)`)
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

	if crossRefRe.MatchString(ipSection) {
		issues = append(issues, gateIssue{
			Section: "Implementation Plan",
			Detail:  "cross-card references found (参考 DES/REQ/TASK)",
			Fix:     "replace with information inline in the step",
		})
	}

	for _, match := range steps {
		stepNum := match[1]
		stepBody := extractSubSection(ipSection, "Step "+stepNum+":")
		var missing []string
		if !hasStepField(stepBody, "Files") {
			missing = append(missing, "Files")
		}
		if !hasStepField(stepBody, "Approach") {
			missing = append(missing, "Approach")
		}
		if !hasStepField(stepBody, "Edge Cases") {
			missing = append(missing, "Edge Cases")
		}
		if len(missing) > 0 {
			issues = append(issues, gateIssue{
				Section: fmt.Sprintf("Implementation Plan.Step %s", stepNum),
				Detail:  fmt.Sprintf("missing required fields: %s", strings.Join(missing, ", ")),
				Fix:     fmt.Sprintf("add %s with substantive content", strings.Join(missing, ", ")),
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
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^- \*\*%s\*\*: (.+)`, field))
	matches := re.FindAllStringSubmatch(stepBody, -1)
	for _, m := range matches {
		val := strings.TrimSpace(m[1])
		if val != "" && !placeholderRe.MatchString(val) && val != "TBD" {
			return true
		}
	}
	return false
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
