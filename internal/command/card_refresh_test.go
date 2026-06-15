package command

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"flowforge/internal/core"
)

func TestCardRefreshGeneratesRequirementAndDesignNavigation(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Navigation proposal")

	store := testCardStore(t, tmpDir)
	req := core.NewCard(core.CardTypeRequirement, "Navigation requirement")
	req.ID = "REQ-nav"
	req.AddLink("ROOT-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(req, proposalID); err != nil {
		t.Fatalf("creating requirement failed: %v", err)
	}

	analysis := core.NewCard(core.CardTypeTask, "Analyze requirement")
	analysis.ID = "TASK-" + proposalID + "-a-nav"
	analysis.Status = core.CardStatusReady
	analysis.AddLink("ROOT-"+proposalID, "belongs_to")
	analysis.AddLink(req.ID, "analyzes")
	if _, err := store.CreateCard(analysis, proposalID); err != nil {
		t.Fatalf("creating analysis task failed: %v", err)
	}

	design := core.NewCard(core.CardTypeDesign, "Navigation design")
	design.ID = "DES-nav"
	design.AddLink("ROOT-"+proposalID, "belongs_to")
	design.AddLink(req.ID, "designs")
	if _, err := store.CreateCard(design, proposalID); err != nil {
		t.Fatalf("creating design failed: %v", err)
	}

	task := core.NewCard(core.CardTypeTask, "Implement navigation")
	task.ID = "TASK-" + proposalID + "-i-nav"
	task.Status = core.CardStatusReady
	task.AddLink("ROOT-"+proposalID, "belongs_to")
	task.AddLink(design.ID, "implements")
	if _, err := store.CreateCard(task, proposalID); err != nil {
		t.Fatalf("creating implementation task failed: %v", err)
	}

	reqRefresh := newCardRefreshCmd()
	var reqOut bytes.Buffer
	reqRefresh.SetOut(&reqOut)
	reqRefresh.SetArgs([]string{req.ID})
	if err := reqRefresh.Execute(); err != nil {
		t.Fatalf("requirement refresh failed: %v", err)
	}
	if !strings.Contains(reqOut.String(), "✓ Refreshed REQ-nav") {
		t.Fatalf("unexpected requirement refresh output:\n%s", reqOut.String())
	}

	refreshedReq, err := store.ReadCard(req.ID)
	if err != nil {
		t.Fatalf("reading refreshed requirement failed: %v", err)
	}
	for _, want := range []string{
		"## FlowForge Navigation",
		"### Analysis Tasks",
		"[TASK-" + proposalID + "-a-nav]",
		"### Design Cards",
		"[DES-nav]",
		"### Implementation Tasks",
		"[TASK-" + proposalID + "-i-nav]",
	} {
		if !strings.Contains(refreshedReq.Body, want) {
			t.Fatalf("requirement navigation missing %q:\n%s", want, refreshedReq.Body)
		}
	}

	designRefresh := newCardRefreshCmd()
	designRefresh.SetArgs([]string{design.ID})
	if err := designRefresh.Execute(); err != nil {
		t.Fatalf("design refresh failed: %v", err)
	}

	refreshedDesign, err := store.ReadCard(design.ID)
	if err != nil {
		t.Fatalf("reading refreshed design failed: %v", err)
	}
	for _, want := range []string{
		"## FlowForge Navigation",
		"### Implementation Tasks",
		"[TASK-" + proposalID + "-i-nav]",
	} {
		if !strings.Contains(refreshedDesign.Body, want) {
			t.Fatalf("design navigation missing %q:\n%s", want, refreshedDesign.Body)
		}
	}

	validateCmd := newValidateAllCmd()
	if err := validateCmd.Execute(); err != nil {
		t.Fatalf("validate all failed after navigation refresh: %v", err)
	}
}
