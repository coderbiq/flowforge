package command

import (
	"os"
	"strings"
	"testing"

	"flowforge/internal/config"
	"flowforge/internal/core"
)

func TestTaskLifecycleCommands(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Task lifecycle proposal")

	store := testCardStore(t, tmpDir)
	design := core.NewCard(core.CardTypeDesign, "Lifecycle design")
	design.ID = "DES-abc123"
	if _, err := store.CreateCard(design, proposalID); err != nil {
		t.Fatalf("creating design card failed: %v", err)
	}

	createCmd := newTaskCreateCmd()
	createCmd.SetArgs([]string{"--title", "Implement task command", "--type", "i", "--links", "DES-abc123:implements"})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("task create failed: %v", err)
	}

	tasks, err := store.ListCardsByType(core.CardTypeTask)
	if err != nil {
		t.Fatalf("listing tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	taskID := tasks[0].ID
	if tasks[0].Status != core.CardStatusReady {
		t.Fatalf("expected created task to be ready, got %s", tasks[0].Status)
	}

	claimCmd := newTaskClaimCmd()
	claimCmd.SetArgs([]string{taskID})
	if err := claimCmd.Execute(); err != nil {
		t.Fatalf("task claim failed: %v", err)
	}

	claimed, err := store.ReadCard(taskID)
	if err != nil {
		t.Fatalf("reading claimed task failed: %v", err)
	}
	if claimed.Status != core.CardStatusInProgress {
		t.Fatalf("expected claimed task to be in_progress, got %s", claimed.Status)
	}

	doneCmd := newTaskDoneCmd()
	doneCmd.SetArgs([]string{taskID, "--summary", "Implemented lifecycle"})
	if err := doneCmd.Execute(); err != nil {
		t.Fatalf("task done failed: %v", err)
	}

	done, err := store.ReadCard(taskID)
	if err != nil {
		t.Fatalf("reading done task failed: %v", err)
	}
	if done.Status != core.CardStatusDone {
		t.Fatalf("expected done task to be done, got %s", done.Status)
	}
	if done.Body == "" {
		t.Fatal("expected done summary to be appended to body")
	}
}

func TestTaskCreateDefaultsToCurrentProposal(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Task default proposal")

	createCmd := newTaskCreateCmd()
	createCmd.SetArgs([]string{"--title", "Implement default proposal task", "--type", "i"})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("task create failed: %v", err)
	}

	store := testCardStore(t, tmpDir)
	tasks, err := store.ListCards(store.ProposalCardsDir(proposalID))
	if err != nil {
		t.Fatalf("listing proposal cards failed: %v", err)
	}

	var found *core.Card
	for _, task := range tasks {
		if task.Type == core.CardTypeTask && strings.Contains(task.Title, "Implement default proposal task") {
			found = task
			break
		}
	}
	if found == nil {
		t.Fatalf("expected task to be created under current proposal %s", proposalID)
	}
	if found.Source != proposalID {
		t.Fatalf("expected task source %s, got %q", proposalID, found.Source)
	}
	if !strings.Contains(found.ID, proposalID) {
		t.Fatalf("expected task ID %q to include proposal %s", found.ID, proposalID)
	}
}

func TestTaskCreateParsesCommaSeparatedLinks(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Task links proposal")

	store := testCardStore(t, tmpDir)
	req := core.NewCard(core.CardTypeRequirement, "Required card")
	req.ID = "REQ-abc"
	if _, err := store.CreateCard(req, proposalID); err != nil {
		t.Fatalf("creating requirement card failed: %v", err)
	}
	des := core.NewCard(core.CardTypeDesign, "Design card")
	des.ID = "DES-def"
	if _, err := store.CreateCard(des, proposalID); err != nil {
		t.Fatalf("creating design card failed: %v", err)
	}
	conv := core.NewCard(core.CardTypeConvention, "Convention card")
	conv.ID = "CONV-ghi"
	if _, err := store.CreateCard(conv, ""); err != nil {
		t.Fatalf("creating convention card failed: %v", err)
	}

	createCmd := newTaskCreateCmd()
	createCmd.SetArgs([]string{
		"--title", "Implement linked task",
		"--type", "i",
		"--links", "REQ-abc:requires,DES-def:implements",
		"--links", "CONV-ghi:constrains",
	})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("task create failed: %v", err)
	}

	tasks, err := store.ListCardsByType(core.CardTypeTask)
	if err != nil {
		t.Fatalf("listing tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	for _, want := range []core.Link{
		{Target: "REQ-abc", Relation: "requires"},
		{Target: "DES-def", Relation: "implements"},
		{Target: "CONV-ghi", Relation: "constrains"},
	} {
		if !hasLinkRelation(tasks[0], want.Target, want.Relation) {
			t.Fatalf("expected link %#v, got %#v", want, tasks[0].Links)
		}
	}
}

func TestTaskSubCreatesDecomposesParentLink(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Subtask proposal")

	store := testCardStore(t, tmpDir)
	parent := core.NewCard(core.CardTypeTask, "Parent task")
	parent.ID = "TASK-" + proposalID + "-i-parent"
	parent.Status = core.CardStatusReady
	parent.Source = proposalID
	parent.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(parent, proposalID); err != nil {
		t.Fatalf("creating parent task failed: %v", err)
	}

	cmd := newTaskSubCmd()
	cmd.SetArgs([]string{parent.ID, "--title", "Child task"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("task sub failed: %v", err)
	}

	child, err := store.ReadCard(parent.ID + "-a")
	if err != nil {
		t.Fatalf("reading child failed: %v", err)
	}
	if !hasLinkRelation(child, parent.ID, "decomposes") {
		t.Fatalf("expected child decomposes parent link, got %#v", child.Links)
	}
	if !hasLinkRelation(child, "PROP-"+proposalID, "belongs_to") {
		t.Fatalf("expected child belongs_to root link, got %#v", child.Links)
	}
}

func TestTaskLinkAddParsesRelationBeforeReadingTarget(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	proposalID := createProposalForTest(t, tmpDir, "Task link proposal")

	store := testCardStore(t, tmpDir)
	task := core.NewCard(core.CardTypeTask, "Link target task")
	task.ID = "TASK-" + proposalID + "-i-link"
	task.Status = core.CardStatusReady
	task.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(task, proposalID); err != nil {
		t.Fatalf("creating task failed: %v", err)
	}
	req := core.NewCard(core.CardTypeRequirement, "Requirement target")
	req.ID = "REQ-link-target"
	req.AddLink("PROP-"+proposalID, "belongs_to")
	if _, err := store.CreateCard(req, proposalID); err != nil {
		t.Fatalf("creating requirement failed: %v", err)
	}

	cmd := newTaskLinkAddCmd()
	cmd.SetArgs([]string{task.ID, "REQ-link-target:requires"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("task link-add failed: %v", err)
	}

	reloaded, err := store.ReadCard(task.ID)
	if err != nil {
		t.Fatalf("reading task failed: %v", err)
	}
	if !hasLinkRelation(reloaded, "REQ-link-target", "requires") {
		t.Fatalf("expected requires link, got %#v", reloaded.Links)
	}
}

func TestTaskCreateRendersSelfIDPlaceholder(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	createProposalForTest(t, tmpDir, "Task self id proposal")

	createCmd := newTaskCreateCmd()
	createCmd.SetArgs([]string{
		"--title", "Implement with self id",
		"--type", "i",
		"--body", "## Goal\n\nUpdate {{task.id}} and <self>.",
	})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("task create failed: %v", err)
	}

	store := testCardStore(t, tmpDir)
	tasks, err := store.ListCardsByType(core.CardTypeTask)
	if err != nil {
		t.Fatalf("listing tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if strings.Contains(tasks[0].Body, "{{task.id}}") || strings.Contains(tasks[0].Body, "<self>") {
		t.Fatalf("expected placeholders to be replaced, got body:\n%s", tasks[0].Body)
	}
	if !strings.Contains(tasks[0].Body, tasks[0].ID) {
		t.Fatalf("expected body to contain task id %s, got:\n%s", tasks[0].ID, tasks[0].Body)
	}
}

func TestTaskCreateAcceptsInitialStatus(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	createProposalForTest(t, tmpDir, "Task status proposal")

	createCmd := newTaskCreateCmd()
	createCmd.SetArgs([]string{"--title", "Clarify impact", "--type", "a", "--status", "not_ready"})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("task create failed: %v", err)
	}

	store := testCardStore(t, tmpDir)
	tasks, err := store.ListCardsByType(core.CardTypeTask)
	if err != nil {
		t.Fatalf("listing tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Status != core.CardStatusNotReady {
		t.Fatalf("expected created task to be not_ready, got %s", tasks[0].Status)
	}
	if !strings.Contains(tasks[0].ID, "-a-") {
		t.Fatalf("expected analysis task ID, got %s", tasks[0].ID)
	}
}

func TestTaskCreateRejectsInvalidInitialStatus(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")

	createCmd := newTaskCreateCmd()
	createCmd.SetArgs([]string{"--title", "Invalid status", "--status", "waiting"})
	err := createCmd.Execute()
	if err == nil {
		t.Fatal("expected invalid task status to fail")
	}
	if !strings.Contains(err.Error(), "invalid task status: waiting") {
		t.Fatalf("expected invalid status error, got %v", err)
	}
}

func TestTaskReadyRequiresAnalysisSections(t *testing.T) {
	tmpDir := t.TempDir()
	restoreWorkingDir(t)

	if err := runInit(tmpDir, true, "default"); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	createProjectForTest(t, "default")
	createProposalForTest(t, tmpDir, "Analysis ready proposal")

	incompleteCmd := newTaskCreateCmd()
	incompleteCmd.SetArgs([]string{"--title", "Incomplete analysis", "--type", "a", "--body", "## Goal\n\nInspect behavior."})
	if err := incompleteCmd.Execute(); err != nil {
		t.Fatalf("creating incomplete analysis task failed: %v", err)
	}

	completeBody := strings.Join([]string{
		"## Goal\n\nInspect behavior.",
		"## Inputs\n\n- Proposal context",
		"## Investigation Plan\n\n- Check CLI output",
		"## Expected Outputs\n\n- Design update",
		"## Done When\n\n- Output contract is clear",
	}, "\n\n")
	completeCmd := newTaskCreateCmd()
	completeCmd.SetArgs([]string{"--title", "Complete analysis", "--type", "a", "--body", completeBody})
	if err := completeCmd.Execute(); err != nil {
		t.Fatalf("creating complete analysis task failed: %v", err)
	}

	readyCmd := newTaskReadyCmd()
	var out strings.Builder
	readyCmd.SetOut(&out)
	readyCmd.SetArgs([]string{"--type", "a"})
	if err := readyCmd.Execute(); err != nil {
		t.Fatalf("task ready failed: %v", err)
	}

	text := out.String()
	if strings.Contains(text, "Incomplete analysis") {
		t.Fatalf("expected incomplete analysis task to be filtered out:\n%s", text)
	}
	if !strings.Contains(text, "Complete analysis") {
		t.Fatalf("expected complete analysis task to be listed:\n%s", text)
	}
}

func restoreWorkingDir(t *testing.T) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore working dir failed: %v", err)
		}
	})
}

func testCardStore(t *testing.T, projectRoot string) *core.CardStore {
	t.Helper()

	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatalf("loading config failed: %v", err)
	}
	return core.NewCardStore(cfg.WikiRoot(projectRoot))
}

func createProjectForTest(t *testing.T, projectID string) {
	t.Helper()

	cmd := newProjectCreateCmd()
	cmd.SetArgs([]string{projectID})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("project create failed: %v", err)
	}
}
