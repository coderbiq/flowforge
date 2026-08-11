package command

import (
	"testing"

	"flowforge/internal/core"
)

func TestProposalFeaturesComplete(t *testing.T) {
	done := &core.Card{Type: core.CardTypeFeature, Status: core.CardStatusDone}
	if !proposalFeaturesComplete(&proposalSnapshot{cards: []*core.Card{done}}) {
		t.Fatal("all done features should complete the proposal")
	}
	planned := &core.Card{Type: core.CardTypeFeature, Status: core.CardStatusPlanned}
	if proposalFeaturesComplete(&proposalSnapshot{cards: []*core.Card{done, planned}}) {
		t.Fatal("planned feature should keep a next action")
	}
	container := &core.Card{Type: core.CardTypeFeature, Status: core.CardStatusDraft, Role: "container"}
	if !proposalFeaturesComplete(&proposalSnapshot{cards: []*core.Card{done, container}}) {
		t.Fatal("container stage should not block child feature completion")
	}
}
