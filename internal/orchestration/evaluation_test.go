package orchestration

import (
	"strings"
	"testing"
	"time"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func TestAnalysisOrchestrationScenarios(t *testing.T) {
	tests := []struct {
		name       string
		events     []core.JournalEvent
		wantState  string
		wantAction string
		analystMin int
		investMin  int
	}{
		{name: "simple direct", wantState: "no_plan", wantAction: "publish_analysis_plan", analystMin: 1},
		{name: "complex multi round", events: []core.JournalEvent{
			evalEvent("plan", "analysis.plan_published", map[string]any{"cycleId": "c1", "revision": 1, "work": []any{map[string]any{"workId": "w1", "question": "inspect code", "required": true}}}),
			evalEvent("dispatch", "work.dispatched", map[string]any{"revision": 1, "workId": "w1"}),
			evalEvent("complete", "work.completed", map[string]any{"revision": 1, "workId": "w1"}),
		}, wantState: "synthesize", wantAction: "invoke_design_analyst", analystMin: 2, investMin: 1},
		{name: "investigator failure", events: []core.JournalEvent{
			evalEvent("plan", "analysis.plan_published", map[string]any{"cycleId": "c1", "revision": 1, "work": []any{map[string]any{"workId": "w1", "required": true}}}),
			evalEvent("dispatch", "work.dispatched", map[string]any{"revision": 1, "workId": "w1"}),
			evalEvent("blocked", "work.blocked", map[string]any{"revision": 1, "workId": "w1"}),
		}, wantState: "synthesize", wantAction: "invoke_design_analyst", analystMin: 2, investMin: 1},
		{name: "user decision", events: []core.JournalEvent{
			evalEvent("decision", "user.decision_required", map[string]any{"reason": "compatibility"}),
		}, wantState: "decision_required", wantAction: "request_user_decision", analystMin: 1},
		{name: "completed recovery", events: []core.JournalEvent{
			evalEvent("done", "analysis.completed", map[string]any{}),
		}, wantState: "completed", wantAction: "none", analystMin: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			view, err := state.BuildAnalysisView("CR-eval", "journal-hash", test.events)
			if err != nil {
				t.Fatalf("folding persisted events: %v", err)
			}
			if view.State != test.wantState || view.NextAction != test.wantAction {
				t.Fatalf("unexpected route: state=%s action=%s", view.State, view.NextAction)
			}
			analystCalls, investigatorCalls := scenarioCallProxy(test.events)
			if analystCalls < test.analystMin || investigatorCalls < test.investMin {
				t.Fatalf("unexpected cost proxy: analyst=%d investigator=%d", analystCalls, investigatorCalls)
			}
		})
	}
}

func TestHostAdaptersRecordProfilesAndEnforcement(t *testing.T) {
	policy := DefaultPolicy()
	codex, err := RenderCodex(policy)
	if err != nil {
		t.Fatal(err)
	}
	opencode, err := RenderOpenCode(policy)
	if err != nil {
		t.Fatal(err)
	}
	for host, files := range map[string]map[string][]byte{"codex": codex, "opencode": opencode} {
		for role, profile := range map[string]ModelProfile{"coordinator": ModelProfileLowCostGeneral, "design-analyst": ModelProfileHighCapability, "investigator": ModelProfileToolCapableReadOnly} {
			ext := ".toml"
			if host == "opencode" {
				ext = ".md"
			}
			content := string(files["flowforge-"+role+ext])
			if !strings.Contains(content, "Model Profile: "+string(profile)) {
				t.Fatalf("%s %s did not record profile %s", host, role, profile)
			}
		}
		if EnforcementSummary(host) == "unsupported" {
			t.Fatalf("missing enforcement mapping for %s", host)
		}
	}
}

func scenarioCallProxy(events []core.JournalEvent) (analyst, investigator int) {
	analyst = 1
	for _, event := range events {
		switch event.Kind {
		case "work.dispatched":
			investigator++
		case "work.completed", "work.blocked", "work.inconclusive", "analysis.reentry_requested":
			analyst = 2
		}
	}
	return analyst, investigator
}

func evalEvent(id, kind string, data map[string]any) core.JournalEvent {
	return core.JournalEvent{ID: id, Version: core.JournalEventSchemaVersion, State: "sealed", Kind: kind, Time: time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC), Data: data}
}
