package command

import (
	"testing"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func TestComplexAnalysisBodyGateDesigned(t *testing.T) {
	body := complexGateBody("accepted: FIND-1", "None", true)
	if issues := validateComplexAnalysisBodyGate(body, core.CardStatusDesigned); len(issues) != 0 {
		t.Fatalf("expected designed gate to pass, got %#v", issues)
	}

	body = complexGateBody("FIND-1", "Investigate storage", true)
	issues := validateComplexAnalysisBodyGate(body, core.CardStatusDesigned)
	assertGateSections(t, issues, "Evidence", "Next Investigation")
}

func TestComplexAnalysisBodyGatePlannedRequiresCompletePlan(t *testing.T) {
	body := complexGateBody("accepted: FIND-1", "None", false)
	issues := validateComplexAnalysisBodyGate(body, core.CardStatusPlanned)
	assertGateSections(t, issues, "Verification", "Implementation Plan.Step 1")
}

func TestComplexAnalysisStateGate(t *testing.T) {
	for _, allowed := range []string{"no_plan", "completed"} {
		if issues := validateComplexAnalysisState(state.AnalysisView{State: allowed}, core.CardStatusDesigned); len(issues) != 0 {
			t.Fatalf("state %s should pass: %#v", allowed, issues)
		}
	}
	view := state.AnalysisView{State: "decision_required", NextAction: "request_user_decision", ActivePlan: &state.AnalysisPlan{Revision: 2}}
	if issues := validateComplexAnalysisState(view, core.CardStatusPlanned); len(issues) != 1 {
		t.Fatalf("active analysis should block: %#v", issues)
	}
}

func TestLegacyFeatureSkipsComplexAnalysisGate(t *testing.T) {
	card := &core.Card{Body: "## Design\n\n### Key Decisions\n\n- keep compatibility"}
	issues, err := complexAnalysisGateIssues(card, core.CardStatusDesigned)
	if err != nil || len(issues) != 0 {
		t.Fatalf("legacy feature must remain compatible: issues=%#v err=%v", issues, err)
	}
}

func complexGateBody(evidence, next string, completePlan bool) string {
	verification := "None"
	extraFields := ""
	if completePlan {
		verification = "- outcome -> go test ./internal/..."
		extraFields = "- **Dependencies**: None\n- **Parallel**: no\n- **Verification**: go test ./internal/...\n"
	}
	return "<!-- analysis-mode: complex -->\n\n" +
		"## Objective\n\nDeliver orchestration; non-goal: product code.\n\n" +
		"## Current Understanding\n\nFact: current routing is single analyst.\n\n" +
		"## Evidence\n\n" + evidence + "\n\n" +
		"## Working Design\n\nUse a recoverable loop.\n\n" +
		"## Rejected or Revised Assumptions\n\nNone\n\n" +
		"## Open Questions\n\nNone\n\n" +
		"## Next Investigation\n\n" + next + "\n\n" +
		"## Implementation Plan\n\n### Step 1: Work\n\n- **Files**: internal/example.go\n- **Approach**: implement\n- **Edge Cases**: empty input\n" + extraFields + "\n" +
		"## Verification\n\n" + verification + "\n"
}

func assertGateSections(t *testing.T, issues []gateIssue, sections ...string) {
	t.Helper()
	found := map[string]bool{}
	for _, issue := range issues {
		found[issue.Section] = true
	}
	for _, section := range sections {
		if !found[section] {
			t.Fatalf("missing issue for %s: %#v", section, issues)
		}
	}
}
