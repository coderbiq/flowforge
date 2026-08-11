package command

import (
	"testing"

	"flowforge/internal/core"
)

func TestBuildPreflightReportAllowsPlannedCompleteStep(t *testing.T) {
	store, root := setupContextTestStore(t)
	card := testFeatureCard("FEAT-ready", core.CardStatusPlanned, completeStepBody(1))
	if _, err := store.CreateCard(card, "CR260810"); err != nil {
		t.Fatal(err)
	}
	r := buildPreflightReport(store, card, 1, "implement")
	if r.Decision != "allow" {
		t.Fatalf("expected allow, got %#v", r)
	}
	if r.Owner != "flowforge-executor" || !r.HandoffRequired {
		t.Fatalf("expected executor handoff, got %#v", r)
	}
	if r.Context != "flowforge context feature --feature FEAT-ready --step 1" {
		t.Fatalf("unexpected context command: %q", r.Context)
	}
	if r.Next != "delegate flowforge-executor with the exact Step context; the primary thread must not implement locally" {
		t.Fatalf("unexpected next action: %q", r.Next)
	}
	_ = root
}

func TestBuildPreflightReportBlocksMissingIntentAndDependency(t *testing.T) {
	store, _ := setupContextTestStore(t)
	dependency := testFeatureCard("FEAT-dependency", core.CardStatusPlanned, completeStepBody(1))
	if _, err := store.CreateCard(dependency, "CR260810"); err != nil {
		t.Fatal(err)
	}
	card := testFeatureCard("FEAT-dependent", core.CardStatusPlanned, completeStepBody(1))
	card.AddLink("FEAT-dependency", "requires")
	if _, err := store.CreateCard(card, "CR260810"); err != nil {
		t.Fatal(err)
	}
	if r := buildPreflightReport(store, card, 1, ""); r.Decision != "blocked" || r.HandoffRequired || r.Owner != "coordinator" || r.Context != "" {
		t.Fatalf("missing intent should block without handoff: %#v", r)
	}
	if r := buildPreflightReport(store, card, 1, "implement"); r.Decision != "blocked" {
		t.Fatalf("unfinished dependency should block: %#v", r)
	}
}

func setupContextTestStore(t *testing.T) (*core.CardStore, string) {
	t.Helper()
	root := t.TempDir()
	store := core.NewCardStore(root)
	if _, _, err := store.CreateProposal("CR260810", "Test"); err != nil {
		t.Fatal(err)
	}
	return store, root
}
func testFeatureCard(id string, status core.CardStatus, body string) *core.Card {
	c := core.NewCard(core.CardTypeFeature, id)
	c.ID = id
	c.Status = status
	c.Body = body
	return c
}
func completeStepBody(n int) string {
	return "## Constraints\n\n- preserve behavior\n\n## Implementation Plan\n\n### Step 1: Work\n\n<!-- step-status: not_started -->\n\n- **Goal**: Complete work\n- **Files**: internal/example.go\n- **Approach**: Change implementation\n- **Edge Cases**: Empty input\n- **Verification**: go test\n\n## Open Questions\n\nNone\n"
}
