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
		"RootCard: ROOT-CR260613",
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
	completedDir := filepath.Join(store.CompletedDir(), proposalID)

	archiveCmd := newProposalArchiveCmd()
	if err := archiveCmd.Execute(); err == nil {
		t.Fatalf("expected archive to require proposal id")
	}
	archiveCmd.SetArgs([]string{proposalID})
	if err := archiveCmd.Execute(); err != nil {
		t.Fatalf("proposal archive failed: %v", err)
	}

	if _, err := os.Stat(activeDir); !os.IsNotExist(err) {
		t.Fatalf("expected active proposal dir to be moved, stat err=%v", err)
	}
	if _, err := os.Stat(completedDir); err != nil {
		t.Fatalf("expected completed proposal dir: %v", err)
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
	if _, err := os.Stat(completedDir); !os.IsNotExist(err) {
		t.Fatalf("expected completed proposal dir to be deleted, stat err=%v", err)
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
