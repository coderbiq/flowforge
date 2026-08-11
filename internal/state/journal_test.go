package state

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"flowforge/internal/core"
)

func TestBuildAnalysisViewFoldsDependenciesAndReentry(t *testing.T) {
	events := []core.JournalEvent{
		testJournalEvent("plan", "analysis.plan_published", map[string]any{
			"cycleId": "cycle-1", "revision": 1,
			"work": []any{
				map[string]any{"workId": "w1", "question": "first", "required": true},
				map[string]any{"workId": "w2", "question": "second", "required": true, "dependsOn": []any{"w1"}},
			},
		}),
		testJournalEvent("dispatch-1", "work.dispatched", map[string]any{"revision": 1, "workId": "w1"}),
		testJournalEvent("complete-1", "work.completed", map[string]any{"revision": 1, "workId": "w1"}),
	}
	view, err := BuildAnalysisView("CR1", "hash", events)
	if err != nil {
		t.Fatalf("building view: %v", err)
	}
	if len(view.ReadyWork) != 1 || view.ReadyWork[0].WorkID != "w2" {
		t.Fatalf("unexpected ready work: %#v", view.ReadyWork)
	}
	if view.NextAction != "dispatch_ready_work" {
		t.Fatalf("unexpected next action: %s", view.NextAction)
	}
}

func TestRebuildAnalysisIsEquivalentAndDetectsStale(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "state.sqlite")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening state: %v", err)
	}
	defer store.Close()
	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("ensuring schema: %v", err)
	}
	events := []core.JournalEvent{testJournalEvent("plan", "analysis.plan_published", map[string]any{"cycleId": "c", "revision": 1, "work": []any{}})}
	first, err := store.RebuildAnalysis("CR1", "revision-a", events)
	if err != nil {
		t.Fatalf("rebuilding analysis: %v", err)
	}
	loaded, ok, err := store.AnalysisView("CR1")
	if err != nil || !ok {
		t.Fatalf("loading analysis: ok=%v err=%v", ok, err)
	}
	want, _ := json.Marshal(first)
	got, _ := json.Marshal(loaded)
	if string(want) != string(got) {
		t.Fatalf("rebuilt view mismatch:\nwant %s\ngot  %s", want, got)
	}
	status, err := store.AnalysisIndexStatus("CR1", "revision-b")
	if err != nil {
		t.Fatalf("checking stale status: %v", err)
	}
	if !status.Present || !status.Stale {
		t.Fatalf("expected stale index: %#v", status)
	}
}

func TestBuildAnalysisViewRejectsIllegalTransition(t *testing.T) {
	events := []core.JournalEvent{
		testJournalEvent("plan", "analysis.plan_published", map[string]any{"cycleId": "c", "revision": 1, "work": []any{map[string]any{"workId": "w1"}}}),
		testJournalEvent("done", "work.completed", map[string]any{"revision": 1, "workId": "w1"}),
	}
	if _, err := BuildAnalysisView("CR1", "hash", events); err == nil {
		t.Fatal("expected illegal transition error")
	}
}

func testJournalEvent(id, kind string, data map[string]any) core.JournalEvent {
	return core.JournalEvent{ID: id, Version: core.JournalEventSchemaVersion, State: "sealed", Kind: kind, Time: time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC), Data: data}
}
