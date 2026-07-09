package command

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newCardStepsCmd() *cobra.Command {
	var (
		status string
		reason string
		start  bool
	)

	cmd := &cobra.Command{
		Use:   "steps <card-id>",
		Short: "Update Implementation Plan step status",
		Long: `Update the status of a step in a FEATURE card's Implementation Plan.
Step status is recorded as an HTML comment (<!-- step-status: ... -->) within the step.

Status values: not_started, in_progress, done, blocked
Use --start to mark first execution (auto-upgrades planned → in_progress).

Examples:
  flowforge card steps FEAT-xxx --status done 3
  flowforge card steps FEAT-xxx --status in_progress 1
  flowforge card steps FEAT-xxx --status blocked 2 --reason "Waiting for FEAT-001 API"
  flowforge card steps FEAT-xxx --start 1
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]
			stepNum, err := strconv.Atoi(args[1])
			if err != nil || stepNum < 1 {
				return fmt.Errorf("invalid step number: %s", args[1])
			}

			if status == "" && !start {
				return fmt.Errorf("--status is required (not_started|in_progress|done|blocked)")
			}
			if start {
				status = "in_progress"
			}
			if !isValidStepStatus(status) {
				return fmt.Errorf("invalid --status: %s (valid: not_started, in_progress, done, blocked)", status)
			}
			if status == "blocked" && reason == "" {
				return fmt.Errorf("--reason is required when status is 'blocked'")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			if err := store.UpdateCardWithLock(cardID, func(card *core.Card) error {
				if card.Type != core.CardTypeFeature {
					return fmt.Errorf("card %s is not a feature card", cardID)
				}

				oldBody := card.Body
				card.Body = setStepStatus(card.Body, stepNum, status, reason)

				if start && card.Status == core.CardStatusPlanned {
					card.Status = core.CardStatusInProgress
				}
				if card.Body == oldBody {
					return fmt.Errorf("step %d not found in Implementation Plan", stepNum)
				}
				return nil
			}); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			action := "updated"
			if start {
				action = "started"
			}
			fmt.Fprintf(out, "✓ step %d %s: %s\n", stepNum, action, status)
			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Step status (not_started|in_progress|done|blocked)")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for blocked status")
	cmd.Flags().BoolVar(&start, "start", false, "Mark step as started (auto-upgrades planned → in_progress)")
	return cmd
}

func isValidStepStatus(s string) bool {
	switch s {
	case "not_started", "in_progress", "done", "blocked":
		return true
	}
	return false
}

var stepStatusCommentRe = regexp.MustCompile(`<!-- step-status: [^>]* -->`)

func setStepStatus(body string, stepNum int, status, reason string) string {
	stepPattern := fmt.Sprintf("### Step %d:", stepNum)
	idx := strings.Index(body, stepPattern)
	if idx < 0 {
		return body
	}

	stepStart := idx
	rest := body[stepStart:]
	nextStep := regexp.MustCompile(`(?m)^### Step \d+:`)
	loc := nextStep.FindStringIndex(rest[len(stepPattern):])
	stepEnd := len(body)
	if loc != nil {
		stepEnd = stepStart + len(stepPattern) + loc[0]
	}

	stepBody := body[stepStart:stepEnd]
	newComment := fmt.Sprintf("<!-- step-status: %s -->", status)

	if strings.Contains(stepBody, "<!-- step-status:") {
		stepBody = stepStatusCommentRe.ReplaceAllString(stepBody, newComment)
	} else {
		firstLineEnd := strings.Index(stepBody, "\n")
		if firstLineEnd < 0 {
			stepBody += " " + newComment
		} else {
			stepBody = stepBody[:firstLineEnd+1] + newComment + "\n" + stepBody[firstLineEnd+1:]
		}
	}

	return body[:stepStart] + stepBody + body[stepEnd:]
}
