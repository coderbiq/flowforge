package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

type preflightIssue struct {
	Code   string `json:"code"`
	Level  string `json:"level"`
	Detail string `json:"detail"`
}

type preflightReport struct {
	Feature         string           `json:"feature"`
	Step            int              `json:"step"`
	Decision        string           `json:"decision"`
	Owner           string           `json:"owner"`
	HandoffRequired bool             `json:"handoffRequired"`
	Context         string           `json:"context,omitempty"`
	Issues          []preflightIssue `json:"issues"`
	Next            string           `json:"next"`
}

func newContextPreflightCmd() *cobra.Command {
	var featureID, intent string
	var step int
	cmd := &cobra.Command{
		Use:   "preflight",
		Short: "Check whether a planned FEATURE step may start",
		RunE: func(cmd *cobra.Command, args []string) error {
			if featureID == "" || step < 1 {
				return fmt.Errorf("--feature and --step are required")
			}
			store, err := currentCardStore()
			if err != nil {
				return err
			}
			card, err := store.ReadCard(featureID)
			if err != nil {
				return err
			}
			report := buildPreflightReport(store, card, step, intent)
			if isJSONOutput(cmd) {
				data, err := json.Marshal(report)
				if err != nil {
					return fmt.Errorf("encoding preflight: %w", err)
				}
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Preflight %s Step %d: %s\n", featureID, step, report.Decision)
			for _, issue := range report.Issues {
				fmt.Fprintf(cmd.OutOrStdout(), "- [%s] %s: %s\n", issue.Level, issue.Code, issue.Detail)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Owner: %s\n", report.Owner)
			if report.HandoffRequired {
				fmt.Fprintln(cmd.OutOrStdout(), "Handoff: required")
				fmt.Fprintf(cmd.OutOrStdout(), "Context: %s\n", report.Context)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "Handoff: not required")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Next: %s\n", report.Next)
			return nil
		},
	}
	cmd.Flags().StringVar(&featureID, "feature", "", "FEATURE card ID")
	cmd.Flags().IntVar(&step, "step", 0, "Step number")
	cmd.Flags().StringVar(&intent, "intent", "", "Confirmed intent (must be implement)")
	return cmd
}

func buildPreflightReport(store *core.CardStore, card *core.Card, step int, intent string) preflightReport {
	contextCommand := fmt.Sprintf("flowforge context feature --feature %s --step %d", card.ID, step)
	r := preflightReport{
		Feature:         card.ID,
		Step:            step,
		Decision:        "allow",
		Owner:           "flowforge-executor",
		HandoffRequired: true,
		Context:         contextCommand,
		Next:            "delegate flowforge-executor with the exact Step context; the primary thread must not implement locally",
	}
	block := func(code, detail string) {
		r.Issues = append(r.Issues, preflightIssue{Code: code, Level: "blocked", Detail: detail})
		r.Decision = "blocked"
		r.Owner = "coordinator"
		r.HandoffRequired = false
		r.Context = ""
		r.Next = "return to flowforge-design or flowforge-feedback"
	}
	if card.Type != core.CardTypeFeature {
		block("not_feature", "card is not a FEATURE")
		return r
	}
	if card.Role == "container" {
		block("container_feature", "container FEATURE has no executable steps")
		return r
	}
	if intent != "implement" {
		block("implementation_intent_missing", "--intent implement is required")
		return r
	}
	if card.Status != core.CardStatusPlanned {
		block("feature_not_planned", fmt.Sprintf("current stage is %s", card.Status))
		return r
	}
	stepBody := extractSubSection(extractSection(card.Body, "Implementation Plan"), fmt.Sprintf("Step %d:", step))
	if stepBody == "" {
		block("step_not_found", "Step is absent from the Implementation Plan")
		return r
	}
	if strings.Contains(stepBody, "step-status: done") || strings.Contains(stepBody, "step-status: blocked") {
		block("step_terminal", "Step is already done or blocked")
	}
	for _, field := range []string{"Goal", "Files", "Symbols", "Actions", "Constraints", "Done When", "Verification"} {
		if !hasStepField(stepBody, field) {
			block("step_missing_field", field+" is required")
		}
	}
	if vaguePlanRe.MatchString(stepBody) {
		block("step_vague_language", "step contains unresolved execution choices")
	}
	if countOpenQuestions(extractSection(card.Body, "Open Questions")) > 0 {
		block("unresolved_open_questions", "FEATURE contains unresolved questions")
	}
	for _, dep := range featureDependencies(store, card) {
		if dep.Status != core.CardStatusDone {
			block("dependency_not_done", fmt.Sprintf("%s is %s", dep.ID, dep.Status))
		}
	}
	return r
}
