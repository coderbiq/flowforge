package command

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type riskReport struct {
	Feature  string   `json:"feature"`
	Step     int      `json:"step"`
	Decision string   `json:"decision"`
	Signals  []string `json:"signals"`
	Next     string   `json:"next"`
}

func newContextRiskReviewCmd() *cobra.Command {
	var featureID, base string
	var step int
	cmd := &cobra.Command{Use: "risk-review", Short: "Route completed work to an independent review when needed", RunE: func(cmd *cobra.Command, args []string) error {
		if featureID == "" || step < 1 || base == "" {
			return fmt.Errorf("--feature, --step, and --base are required")
		}
		store, err := currentCardStore()
		if err != nil {
			return err
		}
		card, err := store.ReadCard(featureID)
		if err != nil {
			return err
		}
		r := riskReport{Feature: featureID, Step: step, Decision: "not_ready", Next: "complete the step and verification first"}
		body := extractSubSection(extractSection(card.Body, "Implementation Plan"), fmt.Sprintf("Step %d:", step))
		if !strings.Contains(body, "step-status: done") {
			return renderRisk(cmd, r)
		}
		out, err := exec.Command("git", "diff", "--name-only", base+"...HEAD").Output()
		if err != nil {
			return fmt.Errorf("reading git diff: %w", err)
		}
		changed := strings.Fields(string(out))
		fields := parseStepFields(body)["Files"]
		for _, path := range changed {
			if strings.Contains(path, "migration") || strings.Contains(path, "schema") || strings.Contains(path, "auth") || strings.Contains(path, "security") {
				r.Signals = append(r.Signals, "sensitive_path:"+path)
			}
			if fields != "" && !strings.Contains(fields, path) {
				r.Signals = append(r.Signals, "scope_drift:"+path)
			}
		}
		if len(changed) > 1 {
			r.Signals = append(r.Signals, "multi_file_change")
		}
		if len(r.Signals) == 0 {
			r.Decision = "review_not_required"
			r.Next = "continue FlowForge completion"
		} else {
			r.Decision = "review_required"
			r.Next = "run independent review or verification"
		}
		return renderRisk(cmd, r)
	}}
	cmd.Flags().StringVar(&featureID, "feature", "", "FEATURE card ID")
	cmd.Flags().IntVar(&step, "step", 0, "Step number")
	cmd.Flags().StringVar(&base, "base", "", "Git revision before implementation")
	return cmd
}
func renderRisk(cmd *cobra.Command, r riskReport) error {
	if isJSONOutput(cmd) {
		b, e := json.Marshal(r)
		if e != nil {
			return e
		}
		_, e = fmt.Fprintln(cmd.OutOrStdout(), string(b))
		return e
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Risk review %s Step %d: %s\n", r.Feature, r.Step, r.Decision)
	for _, s := range r.Signals {
		fmt.Fprintln(cmd.OutOrStdout(), "-", s)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Next:", r.Next)
	return nil
}
