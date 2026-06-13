package command

import (
	"os"
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

	createCmd := newTaskCreateCmd()
	createCmd.SetArgs([]string{"--title", "Implement task command", "--type", "i", "--links", "DES-abc123:implements"})
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
