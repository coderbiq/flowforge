package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func TestProposalLifecycleCommandsUseCurrentProposalPointer(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")

	firstProposalID := createProposalForTest(t, tmpDir, "First proposal")
	store := runtimeStateStore(t, tmpDir)
	currentProposalID, ok, err := store.CurrentProposalID("default")
	if err != nil {
		t.Fatalf("CurrentProposalID failed: %v", err)
	}
	if !ok || currentProposalID != firstProposalID {
		t.Fatalf("expected current proposal %s, got %q ok=%v", firstProposalID, currentProposalID, ok)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}

	secondProposalID := "CR260613-2"
	cardStore := testCardStore(t, tmpDir)
	if _, _, err := cardStore.CreateProposal(secondProposalID, "Second proposal"); err != nil {
		t.Fatalf("manual CreateProposal failed: %v", err)
	}
	runtimeStore := runtimeStateStore(t, tmpDir)
	currentProposalID, ok, err = runtimeStore.CurrentProposalID("default")
	if err != nil {
		t.Fatalf("CurrentProposalID after manual create failed: %v", err)
	}
	if !ok || currentProposalID != firstProposalID {
		t.Fatalf("expected current proposal to remain %s, got %q ok=%v", firstProposalID, currentProposalID, ok)
	}
	if err := runtimeStore.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}

	listCmd := newProposalListCmd()
	var listOut bytes.Buffer
	listCmd.SetOut(&listOut)
	if err := listCmd.Execute(); err != nil {
		t.Fatalf("proposal list failed: %v", err)
	}
	listText := listOut.String()
	for _, want := range []string{
		"Active proposals:",
		"- " + firstProposalID,
		"- " + secondProposalID,
	} {
		if !strings.Contains(listText, want) {
			t.Fatalf("proposal list output missing %q:\n%s", want, listText)
		}
	}

	useCmd := newProposalUseCmd()
	var useOut bytes.Buffer
	useCmd.SetOut(&useOut)
	useCmd.SetArgs([]string{secondProposalID})
	if err := useCmd.Execute(); err != nil {
		t.Fatalf("proposal use failed: %v", err)
	}
	if !strings.Contains(useOut.String(), "✓ Current proposal: "+secondProposalID) {
		t.Fatalf("proposal use output missing current proposal:\n%s", useOut.String())
	}

	store = runtimeStateStore(t, tmpDir)
	currentProposalID, ok, err = store.CurrentProposalID("default")
	if err != nil {
		t.Fatalf("CurrentProposalID after use failed: %v", err)
	}
	if !ok || currentProposalID != secondProposalID {
		t.Fatalf("expected current proposal %s after use, got %q ok=%v", secondProposalID, currentProposalID, ok)
	}

	currentCmd := newProposalCurrentCmd()
	var currentOut bytes.Buffer
	currentCmd.SetOut(&currentOut)
	if err := currentCmd.Execute(); err != nil {
		t.Fatalf("proposal current failed: %v", err)
	}
	if !strings.Contains(currentOut.String(), "Proposal: "+secondProposalID) {
		t.Fatalf("proposal current output missing current proposal:\n%s", currentOut.String())
	}

	contextCmd := newContextProposalCmd()
	var contextOut bytes.Buffer
	contextCmd.SetOut(&contextOut)
	if err := contextCmd.Execute(); err != nil {
		t.Fatalf("context proposal failed: %v", err)
	}
	contextText := contextOut.String()
	for _, want := range []string{
		"## Context",
		"- Proposal: " + secondProposalID,
		"- Focus: STR-" + secondProposalID + "-REQ",
		"## Root Summary",
		"## Requirement Map",
		"## Library Discovery Suggestions",
	} {
		if !strings.Contains(contextText, want) {
			t.Fatalf("context proposal output missing %q:\n%s", want, contextText)
		}
	}

	if err := store.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}
}

func TestProposalCreateOutputPassesValidateAll(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	createProposalForTest(t, tmpDir, "Valid empty proposal")

	cmd := newValidateAllCmd()
	if err := cmd.Execute(); err != nil {
		t.Fatalf("validate all failed for newly created proposal: %v", err)
	}
}

func TestProposalCurrentPointerIsolatedPerProject(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "frontend")
	createProjectForTest(t, "backend")

	if err := useProjectForTest(t, "frontend"); err != nil {
		t.Fatalf("using frontend project failed: %v", err)
	}
	frontendProposalID := createProposalForTest(t, tmpDir, "Frontend proposal")

	if err := useProjectForTest(t, "backend"); err != nil {
		t.Fatalf("using backend project failed: %v", err)
	}
	backendProposalID := createProposalForTest(t, tmpDir, "Backend proposal")

	store := runtimeStateStore(t, tmpDir)
	frontendCurrent, ok, err := store.CurrentProposalID("frontend")
	if err != nil {
		t.Fatalf("CurrentProposalID frontend failed: %v", err)
	}
	if !ok || frontendCurrent != frontendProposalID {
		t.Fatalf("expected frontend current proposal %s, got %q ok=%v", frontendProposalID, frontendCurrent, ok)
	}
	backendCurrent, ok, err := store.CurrentProposalID("backend")
	if err != nil {
		t.Fatalf("CurrentProposalID backend failed: %v", err)
	}
	if !ok || backendCurrent != backendProposalID {
		t.Fatalf("expected backend current proposal %s, got %q ok=%v", backendProposalID, backendCurrent, ok)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}

	if err := useProjectForTest(t, "frontend"); err != nil {
		t.Fatalf("switching to frontend failed: %v", err)
	}
	if got := runProposalCurrentForTest(t); !strings.Contains(got, "Proposal: "+frontendProposalID) {
		t.Fatalf("expected frontend proposal current output to mention %s, got %s", frontendProposalID, got)
	}
	if got := runProposalContextForTest(t); !strings.Contains(got, "- Proposal: "+frontendProposalID) {
		t.Fatalf("expected frontend context output to mention %s, got %s", frontendProposalID, got)
	}
	if got := runProposalListForTest(t); !strings.Contains(got, "- "+frontendProposalID) {
		t.Fatalf("expected frontend proposal list output to mention %s, got %s", frontendProposalID, got)
	}

	if err := useProjectForTest(t, "backend"); err != nil {
		t.Fatalf("switching to backend failed: %v", err)
	}
	if got := runProposalCurrentForTest(t); !strings.Contains(got, "Proposal: "+backendProposalID) {
		t.Fatalf("expected backend proposal current output to mention %s, got %s", backendProposalID, got)
	}
	if got := runProposalContextForTest(t); !strings.Contains(got, "- Proposal: "+backendProposalID) {
		t.Fatalf("expected backend context output to mention %s, got %s", backendProposalID, got)
	}
	if got := runProposalListForTest(t); !strings.Contains(got, "- "+backendProposalID) {
		t.Fatalf("expected backend proposal list output to mention %s, got %s", backendProposalID, got)
	}
}

func TestNextProposalIDForPrefixUsesDailySequence(t *testing.T) {
	tmpDir := t.TempDir()
	store := core.NewCardStore(filepath.Join(tmpDir, "ff-wiki"))

	nextID, err := nextProposalIDForPrefix(store, "CR260613")
	if err != nil {
		t.Fatalf("nextProposalIDForPrefix on empty store failed: %v", err)
	}
	if nextID != "CR26061301" {
		t.Fatalf("expected first proposal ID CR26061301, got %s", nextID)
	}

	if _, _, err := store.CreateProposal("CR260613", "Old format proposal"); err != nil {
		t.Fatalf("creating old format proposal failed: %v", err)
	}
	nextID, err = nextProposalIDForPrefix(store, "CR260613")
	if err != nil {
		t.Fatalf("nextProposalIDForPrefix with old format proposal failed: %v", err)
	}
	if nextID != "CR26061301" {
		t.Fatalf("expected old format proposal to reserve only base ID, got %s", nextID)
	}

	for _, proposalID := range []string{"CR26061301", "CR26061303"} {
		if _, _, err := store.CreateProposal(proposalID, proposalID); err != nil {
			t.Fatalf("creating active proposal %s failed: %v", proposalID, err)
		}
	}
	completedProposalDir := filepath.Join(store.CompletedDir(), "CR26061302")
	if err := os.MkdirAll(completedProposalDir, 0755); err != nil {
		t.Fatalf("creating completed proposal dir failed: %v", err)
	}

	nextID, err = nextProposalIDForPrefix(store, "CR260613")
	if err != nil {
		t.Fatalf("nextProposalIDForPrefix with existing sequence failed: %v", err)
	}
	if nextID != "CR26061304" {
		t.Fatalf("expected next proposal ID CR26061304, got %s", nextID)
	}
}

func TestProposalInspectAndContextCommands(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")

	store := testCardStore(t, tmpDir)
	rootPath, indexPath, err := store.CreateProposal("CR260613", "Command Proposal")
	if err != nil {
		t.Fatalf("CreateProposal failed: %v", err)
	}
	if _, err := os.Stat(rootPath); err != nil {
		t.Fatalf("root card missing: %v", err)
	}
	if _, err := os.Stat(indexPath); err != nil {
		t.Fatalf("requirement index missing: %v", err)
	}

	task := core.NewCard(core.CardTypeTask, "Analyze command output")
	task.ID = core.GenerateTaskID("260613", "a")
	task.Status = core.CardStatusReady
	task.Body = strings.Join([]string{
		"## Goal\n\nInspect command output.",
		"## Inputs\n\n- Proposal context",
		"## Investigation Plan\n\n- Compare current output to contract",
		"## Expected Outputs\n\n- Stable output fields",
		"## Done When\n\n- Output fields are stable",
		"## Open Questions\n\n- Which output fields are stable?",
	}, "\n\n")
	task.AddLink("STR-CR260613-REQ", "analyzes")
	if _, err := store.CreateCard(task, "CR260613"); err != nil {
		t.Fatalf("creating analysis task failed: %v", err)
	}

	notReadyTask := core.NewCard(core.CardTypeTask, "Implement command output")
	notReadyTask.ID = core.GenerateTaskID("260613", "i")
	notReadyTask.Status = core.CardStatusNotReady
	notReadyTask.Body = "## Goal\n\nImplement command output."
	if _, err := store.CreateCard(notReadyTask, "CR260613"); err != nil {
		t.Fatalf("creating not-ready task failed: %v", err)
	}

	notReadyAnalysis := core.NewCard(core.CardTypeTask, "Clarify command output")
	notReadyAnalysis.ID = "TASK-260613-a-notready"
	notReadyAnalysis.Status = core.CardStatusNotReady
	notReadyAnalysis.Body = "## Goal\n\nClarify command output."
	if _, err := store.CreateCard(notReadyAnalysis, "CR260613"); err != nil {
		t.Fatalf("creating not-ready analysis task failed: %v", err)
	}

	blockedTask := core.NewCard(core.CardTypeTask, "Implement after analysis")
	blockedTask.ID = "TASK-260613-i-blocked"
	blockedTask.Status = core.CardStatusBlocked
	blockedTask.Body = strings.Join([]string{
		"## Goal\n\nImplement after analysis.",
		"## Deliverables\n\n- Command implementation",
		"## Acceptance\n\n- Tests pass",
		"## Blocked\n\n- Waiting for analysis output",
	}, "\n\n")
	blockedTask.AddLink(task.ID, "depends")
	if _, err := store.CreateCard(blockedTask, "CR260613"); err != nil {
		t.Fatalf("creating blocked task failed: %v", err)
	}

	design := core.NewCard(core.CardTypeDesign, "Command output design")
	design.ID = "DES-260613-output"
	design.Body = "## Goal\n\nDesign command output."
	if _, err := store.CreateCard(design, "CR260613"); err != nil {
		t.Fatalf("creating design card failed: %v", err)
	}

	inspectCmd := newProposalInspectCmd()
	var inspectOut bytes.Buffer
	inspectCmd.SetOut(&inspectOut)
	inspectCmd.SetArgs([]string{"CR260613"})
	if err := inspectCmd.Execute(); err != nil {
		t.Fatalf("proposal inspect failed: %v", err)
	}
	inspectText := inspectOut.String()
	for _, want := range []string{
		"## Proposal",
		"RootCard: PROP-CR260613",
		"RequirementIndex: STR-CR260613-REQ",
		"## Task Summary",
		"## Open Questions",
		"Analyze command output",
		"| ID | Title | Status | Analyzes | Done When |",
		"STR-CR260613-REQ",
		"Output fields are stable",
		"| ID | Title | Status | Missing |",
		"Implement command output",
		"Deliverables",
		"Clarify command output",
		"Implement after analysis",
		"blocked: Waiting for analysis output",
		"dependency " + task.ID + " is ready",
	} {
		if !strings.Contains(inspectText, want) {
			t.Fatalf("proposal inspect output missing %q:\n%s", want, inspectText)
		}
	}
	if strings.Contains(inspectText, "Clarify command output | not_ready | Inputs, Investigation Plan, Expected Outputs, Done When, links") {
		t.Fatalf("not-ready analysis task should not require links:\n%s", inspectText)
	}

	contextCmd := newContextProposalCmd()
	var contextOut bytes.Buffer
	contextCmd.SetOut(&contextOut)
	contextCmd.SetArgs([]string{"--proposal", "CR260613"})
	if err := contextCmd.Execute(); err != nil {
		t.Fatalf("context proposal failed: %v", err)
	}
	contextText := contextOut.String()
	for _, want := range []string{
		"## Context",
		"Focus: STR-CR260613-REQ",
		"## Root Summary",
		"## Requirement Map",
		"## Library Discovery Suggestions",
	} {
		if !strings.Contains(contextText, want) {
			t.Fatalf("context proposal output missing %q:\n%s", want, contextText)
		}
	}
	if strings.Contains(contextText, "Omitted: 1 non-focused requirement map cards") {
		t.Fatalf("context proposal should not count design cards as omitted requirement map cards:\n%s", contextText)
	}
}

func TestProposalInspectReportsStructureAndNavigationHealth(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Health proposal")

	store := testCardStore(t, tmpDir)
	req := core.NewCard(core.CardTypeRequirement, "Unindexed requirement")
	req.ID = "REQ-health"
	req.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(req, proposalID); err != nil {
		t.Fatalf("creating requirement failed: %v", err)
	}

	design := core.NewCard(core.CardTypeDesign, "Health design")
	design.ID = "DES-health"
	design.AddLink("PROP-"+proposalID, "belongs_to")
	design.AddLink(req.ID, "designs")
	if _, err := store.CreateCard(design, proposalID); err != nil {
		t.Fatalf("creating design failed: %v", err)
	}

	task := core.NewCard(core.CardTypeTask, "Implement health design")
	task.ID = "TASK-" + proposalID + "-i-health"
	task.Status = core.CardStatusReady
	task.Body = strings.Join([]string{
		"## Goal\n\nImplement health design.",
		"## Deliverables\n\n- Code",
		"## Acceptance\n\n- Tests pass",
		"## Out of Scope\n\n- None",
		"## Read Before Work\n\n- Design",
	}, "\n\n")
	task.AddLink("PROP-"+proposalID, "belongs_to")
	task.AddLink(design.ID, "implements")
	if _, err := store.CreateCard(task, proposalID); err != nil {
		t.Fatalf("creating task failed: %v", err)
	}

	// P6: orphan design without requirement link should trigger health warning
	orphanDesign := core.NewCard(core.CardTypeDesign, "Orphan design")
	orphanDesign.ID = "DES-orphan"
	orphanDesign.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(orphanDesign, proposalID); err != nil {
		t.Fatalf("creating orphan design failed: %v", err)
	}

	// P1: finding card should appear in context
	finding := core.NewCard(core.CardTypeFinding, "Important restriction")
	finding.ID = "FIND-health"
	finding.Body = "transType determines clearance ruleset."
	finding.AddLink("PROP-"+proposalID, "belongs_to")
	finding.AddLink(design.ID, "references")
	if _, err := store.CreateCard(finding, proposalID); err != nil {
		t.Fatalf("creating finding failed: %v", err)
	}

	// P2: create a card linked to the design to test extended context
	refCard := core.NewCard(core.CardTypeRequirement, "Referenced requirement")
	refCard.ID = "REQ-ref"
	refCard.AddLink("PROP-"+proposalID, "belongs_to")
	refCard.AddLink(finding.ID, "references")
	if _, err := store.CreateCard(refCard, proposalID); err != nil {
		t.Fatalf("creating ref requirement failed: %v", err)
	}

	inspectCmd := newProposalInspectCmd()
	var inspectOut bytes.Buffer
	inspectCmd.SetOut(&inspectOut)
	inspectCmd.SetArgs([]string{proposalID})
	if err := inspectCmd.Execute(); err != nil {
		t.Fatalf("proposal inspect failed: %v", err)
	}
	inspectText := inspectOut.String()
	for _, want := range []string{
		"## Health Issues",
		"PROP-" + proposalID,
		"proposal card has no meaningful summary",
		"STR-" + proposalID + "-REQ",
		"structure card has no meaningful purpose description",
		"REQ-health",
		"requirement is not reachable from a requirement index",
		"requirement navigation is stale or missing",
		"DES-health",
		"design navigation is stale or missing",
		"DES-orphan",
		"design card does not link to a requirement",
		"TASK-" + proposalID + "-i-health",
		"ready implementation task has no linked convention constraints",
		"flowforge card refresh REQ-health",
		"flowforge card refresh DES-health",
	} {
		if !strings.Contains(inspectText, want) {
			t.Fatalf("proposal inspect health output missing %q:\n%s", want, inspectText)
		}
	}

	contextCmd := newContextProposalCmd()
	var contextOut bytes.Buffer
	contextCmd.SetOut(&contextOut)
	contextCmd.SetArgs([]string{"--proposal", proposalID, "--cards", req.ID})
	if err := contextCmd.Execute(); err != nil {
		t.Fatalf("context proposal failed: %v", err)
	}
	contextText := contextOut.String()
	for _, want := range []string{
		"## Health Summary",
		"## Proposal Findings",
		"FIND-health",
		"transType determines",
		"## Extended Context",
		"REQ-health",
	} {
		if !strings.Contains(contextText, want) {
			t.Fatalf("context health output missing %q:\n%s", want, contextText)
		}
	}
}

func TestProposalArchiveAndDeleteCommands(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Archive proposal")

	store := testCardStore(t, tmpDir)
	activeDir := store.ProposalDir(proposalID)

	archiveCmd := newProposalArchiveCmd()
	if err := archiveCmd.Execute(); err == nil {
		t.Fatalf("expected archive to require proposal id")
	}
	archiveCmd.SetArgs([]string{proposalID})
	if err := archiveCmd.Execute(); err != nil {
		t.Fatalf("proposal archive failed: %v", err)
	}

	if _, err := os.Stat(activeDir); err != nil {
		t.Fatalf("expected proposal dir to still exist after archive, stat err=%v", err)
	}

	propCard, err := store.ReadCard("PROP-" + proposalID)
	if err != nil {
		t.Fatalf("reading PROP card after archive: %v", err)
	}
	if propCard.Status != core.CardStatusCompleted {
		t.Fatalf("expected PROP status 'completed', got '%s'", propCard.Status)
	}

	runtimeStore := runtimeStateStore(t, tmpDir)
	currentID, ok, err := runtimeStore.CurrentProposalID("default")
	if err != nil {
		t.Fatalf("CurrentProposalID failed: %v", err)
	}
	if ok {
		t.Fatalf("expected current proposal to be cleared, got %s", currentID)
	}
	if err := runtimeStore.Close(); err != nil {
		t.Fatalf("closing runtime store failed: %v", err)
	}

	rejectDeleteCmd := newProposalDeleteCmd()
	rejectDeleteCmd.SetArgs([]string{proposalID})
	if err := rejectDeleteCmd.Execute(); err == nil {
		t.Fatalf("expected proposal delete without --force to fail")
	}

	deleteCmd := newProposalDeleteCmd()
	deleteCmd.SetArgs([]string{proposalID, "--force"})
	if err := deleteCmd.Execute(); err != nil {
		t.Fatalf("proposal delete failed: %v", err)
	}
	if _, err := os.Stat(activeDir); !os.IsNotExist(err) {
		t.Fatalf("expected proposal dir to be deleted, stat err=%v", err)
	}
}

func createProposalForTest(t *testing.T, projectRoot string, title string) string {
	t.Helper()

	cmd := newProposalCreateCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{title})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("proposal create failed: %v", err)
	}

	store := runtimeStateStore(t, projectRoot)
	defer func() {
		if err := store.Close(); err != nil {
			t.Fatalf("closing runtime store failed: %v", err)
		}
	}()

	projectID, ok, err := store.CurrentProjectID()
	if err != nil {
		t.Fatalf("CurrentProjectID failed: %v", err)
	}
	if !ok {
		t.Fatal("expected current project to be set")
	}

	currentID, ok, err := store.CurrentProposalID(projectID)
	if err != nil {
		t.Fatalf("CurrentProposalID failed: %v", err)
	}
	if !ok {
		t.Fatal("expected current proposal to be set")
	}
	return currentID
}

func runtimeStateStore(t *testing.T, projectRoot string) *state.Store {
	t.Helper()

	store, err := state.Open(filepath.Join(projectRoot, ".flowforge", "cache", "flowforge.sqlite"))
	if err != nil {
		t.Fatalf("opening runtime state store failed: %v", err)
	}
	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("ensuring runtime state schema failed: %v", err)
	}

	return store
}

func useProjectForTest(t *testing.T, projectID string) error {
	t.Helper()

	cmd := newProjectUseCmd()
	cmd.SetArgs([]string{projectID})
	return cmd.Execute()
}

func runProposalCurrentForTest(t *testing.T) string {
	t.Helper()

	cmd := newProposalCurrentCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("proposal current failed: %v", err)
	}
	return out.String()
}

func runProposalContextForTest(t *testing.T) string {
	t.Helper()

	cmd := newContextProposalCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("context proposal failed: %v", err)
	}
	return out.String()
}

func runProposalListForTest(t *testing.T) string {
	t.Helper()

	cmd := newProposalListCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("proposal list failed: %v", err)
	}
	return out.String()
}

func TestReadyTaskWithEmptyBodyIsError(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Empty task proposal")

	store := testCardStore(t, tmpDir)
	task := core.NewCard(core.CardTypeTask, "Implement empty check")
	task.ID = "TASK-" + proposalID + "-i-empty"
	task.Status = core.CardStatusReady
	task.Body = "## Links\n\n### Outgoing\n\n- [PROP-" + proposalID + "](../../../03-proposal/" + proposalID + "_empty-task-proposal.md) [proposal] - Empty task proposal"
	task.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(task, proposalID); err != nil {
		t.Fatalf("creating empty ready task failed: %v", err)
	}

	inspectCmd := newProposalInspectCmd()
	var inspectOut bytes.Buffer
	inspectCmd.SetOut(&inspectOut)
	inspectCmd.SetArgs([]string{proposalID})
	if err := inspectCmd.Execute(); err != nil {
		t.Fatalf("proposal inspect failed: %v", err)
	}
	inspectText := inspectOut.String()
	for _, want := range []string{
		"TASK-" + proposalID + "-i-empty",
		"ready task has no body content",
		"error",
	} {
		if !strings.Contains(inspectText, want) {
			t.Fatalf("inspect should flag empty ready task as error, missing %q:\n%s", want, inspectText)
		}
	}
}

func TestReadyTaskWithBodyIsNotFlagged(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Full task proposal")

	store := testCardStore(t, tmpDir)
	task := core.NewCard(core.CardTypeTask, "Implement full task")
	task.ID = "TASK-" + proposalID + "-i-full"
	task.Status = core.CardStatusReady
	task.Body = "## Goal\n\nImplement the feature.\n\n## Deliverables\n\n- Code\n\n## Acceptance\n\n- Tests pass\n\n## Links\n\n### Outgoing\n\n- [PROP-" + proposalID + "](../../../03-proposal/" + proposalID + "_full-task-proposal.md) [proposal]"
	task.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(task, proposalID); err != nil {
		t.Fatalf("creating full ready task failed: %v", err)
	}

	inspectCmd := newProposalInspectCmd()
	var inspectOut bytes.Buffer
	inspectCmd.SetOut(&inspectOut)
	inspectCmd.SetArgs([]string{proposalID})
	if err := inspectCmd.Execute(); err != nil {
		t.Fatalf("proposal inspect failed: %v", err)
	}
	inspectText := inspectOut.String()
	if strings.Contains(inspectText, "ready task has no body content") {
		t.Fatalf("inspect should NOT flag task with body content:\n%s", inspectText)
	}
}

func TestProposalAndStructureCardsWithMeaningfulContentAreNotFlagged(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Meaningful proposal")

	store := testCardStore(t, tmpDir)

	rootCard, err := store.ReadCard("PROP-" + proposalID)
	if err != nil {
		t.Fatalf("reading root card failed: %v", err)
	}
	rootCard.Body = "# Meaningful proposal\n\n## Purpose\n\nStable entry.\n\n## Summary\n\nThis proposal implements user authentication with JWT tokens, role-based access control, and session management."
	if err := store.UpdateCard(rootCard); err != nil {
		t.Fatalf("updating root card failed: %v", err)
	}

	strCard, err := store.ReadCard("STR-" + proposalID + "-REQ")
	if err != nil {
		t.Fatalf("reading STR card failed: %v", err)
	}
	strCard.Body = "# Meaningful proposal Requirements\n\n## Purpose\n\nThis index organizes requirements for the authentication module.\n\n## Entries\n\n- None\n\n## Open Questions\n\n- None"
	if err := store.UpdateCard(strCard); err != nil {
		t.Fatalf("updating STR card failed: %v", err)
	}

	inspectCmd := newProposalInspectCmd()
	var inspectOut bytes.Buffer
	inspectCmd.SetOut(&inspectOut)
	inspectCmd.SetArgs([]string{proposalID})
	if err := inspectCmd.Execute(); err != nil {
		t.Fatalf("proposal inspect failed: %v", err)
	}
	inspectText := inspectOut.String()
	for _, want := range []string{
		"proposal card has no meaningful summary",
		"structure card has no meaningful purpose description",
	} {
		if strings.Contains(inspectText, want) {
			t.Fatalf("inspect should NOT flag card with meaningful content for %q:\n%s", want, inspectText)
		}
	}
}
