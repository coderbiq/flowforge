package command

import (
	"bytes"
	"strings"
	"testing"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func TestRenderInvestigatorContextIsBounded(t *testing.T) {
	card := &core.Card{ID: "FEAT-1"}
	work := state.AnalysisWork{WorkID: "W-1", Question: "What happens?", Scope: "internal/core only", State: "ready", EvidenceTarget: "FIND-1", DoneWhen: "tests cited", Sources: []string{"code", "library"}, Inputs: []string{"FEAT-1"}, Budget: 3}
	var out bytes.Buffer
	if err := renderInvestigatorContext(&out, card, work); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"What happens?", "FIND-1", "internal/core only", "edit only the designated FIND"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q in:\n%s", want, out.String())
		}
	}
}

func TestRenderAnalysisContextSeparatesCoordinatorAndAnalystDetail(t *testing.T) {
	view := state.AnalysisView{State: "investigating", NextAction: "dispatch_ready_work", ReadyWork: []state.AnalysisWork{{WorkID: "W-1", Question: "Inspect code"}}, History: []core.JournalEvent{{ID: "EV-1"}}}
	var coordinator bytes.Buffer
	if err := renderAnalysisContext(&coordinator, view, false); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(coordinator.String(), "Sealed Events") {
		t.Fatalf("coordinator should receive scheduling state only:\n%s", coordinator.String())
	}
	var analyst bytes.Buffer
	if err := renderAnalysisContext(&analyst, view, true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(analyst.String(), "Sealed Events: 1") {
		t.Fatalf("analyst context should include synthesis history count:\n%s", analyst.String())
	}
}

func TestFindAnalysisWorkAcrossStates(t *testing.T) {
	view := state.AnalysisView{BlockedWork: []state.AnalysisWork{{WorkID: "W-blocked"}}}
	if _, ok := findAnalysisWork(view, "W-blocked"); !ok {
		t.Fatal("expected blocked work to remain recoverable")
	}
}
