package state

import (
	"path/filepath"
	"testing"
)

func TestEnsureSchemaCreatesRuntimeStateTable(t *testing.T) {
	store := openTestStore(t)
	defer closeTestStore(t, store)

	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("EnsureSchema failed: %v", err)
	}

	var name string
	err := store.db.QueryRow(
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'runtime_state'",
	).Scan(&name)
	if err != nil {
		t.Fatalf("querying sqlite_master failed: %v", err)
	}
	if name != "runtime_state" {
		t.Fatalf("expected runtime_state table, got %q", name)
	}
}

func TestCurrentProjectIDLifecycle(t *testing.T) {
	store := openTestStore(t)
	defer closeTestStore(t, store)

	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("EnsureSchema failed: %v", err)
	}

	value, ok, err := store.CurrentProjectID()
	if err != nil {
		t.Fatalf("CurrentProjectID before set failed: %v", err)
	}
	if ok {
		t.Fatalf("expected no current project, got %q", value)
	}

	if err := store.SetCurrentProjectID("frontend"); err != nil {
		t.Fatalf("SetCurrentProjectID failed: %v", err)
	}

	value, ok, err = store.CurrentProjectID()
	if err != nil {
		t.Fatalf("CurrentProjectID after set failed: %v", err)
	}
	if !ok || value != "frontend" {
		t.Fatalf("expected current project frontend, got %q, ok=%v", value, ok)
	}

	if err := store.ClearCurrentProjectID(); err != nil {
		t.Fatalf("ClearCurrentProjectID failed: %v", err)
	}

	value, ok, err = store.CurrentProjectID()
	if err != nil {
		t.Fatalf("CurrentProjectID after clear failed: %v", err)
	}
	if ok {
		t.Fatalf("expected cleared current project, got %q", value)
	}
}

func TestCurrentProjectIDRejectsEmptyValue(t *testing.T) {
	store := openTestStore(t)
	defer closeTestStore(t, store)

	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("EnsureSchema failed: %v", err)
	}

	if err := store.SetCurrentProjectID(""); err == nil {
		t.Fatalf("expected empty project ID to be rejected")
	}
}

func TestCurrentProposalIDLifecyclePerProject(t *testing.T) {
	store := openTestStore(t)
	defer closeTestStore(t, store)

	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("EnsureSchema failed: %v", err)
	}

	if err := store.SetCurrentProposalID("frontend", "CR26061301"); err != nil {
		t.Fatalf("SetCurrentProposalID frontend failed: %v", err)
	}
	if err := store.SetCurrentProposalID("backend", "CR26061302"); err != nil {
		t.Fatalf("SetCurrentProposalID backend failed: %v", err)
	}

	frontendValue, frontendOK, err := store.CurrentProposalID("frontend")
	if err != nil {
		t.Fatalf("CurrentProposalID frontend failed: %v", err)
	}
	if !frontendOK || frontendValue != "CR26061301" {
		t.Fatalf("expected frontend proposal CR26061301, got %q, ok=%v", frontendValue, frontendOK)
	}

	backendValue, backendOK, err := store.CurrentProposalID("backend")
	if err != nil {
		t.Fatalf("CurrentProposalID backend failed: %v", err)
	}
	if !backendOK || backendValue != "CR26061302" {
		t.Fatalf("expected backend proposal CR26061302, got %q, ok=%v", backendValue, backendOK)
	}

	if err := store.ClearCurrentProposalID("frontend"); err != nil {
		t.Fatalf("ClearCurrentProposalID frontend failed: %v", err)
	}

	frontendValue, frontendOK, err = store.CurrentProposalID("frontend")
	if err != nil {
		t.Fatalf("CurrentProposalID frontend after clear failed: %v", err)
	}
	if frontendOK {
		t.Fatalf("expected frontend proposal to be cleared, got %q", frontendValue)
	}

	backendValue, backendOK, err = store.CurrentProposalID("backend")
	if err != nil {
		t.Fatalf("CurrentProposalID backend after frontend clear failed: %v", err)
	}
	if !backendOK || backendValue != "CR26061302" {
		t.Fatalf("expected backend proposal to remain CR26061302, got %q, ok=%v", backendValue, backendOK)
	}
}

func TestCurrentProposalIDRejectsEmptyValues(t *testing.T) {
	store := openTestStore(t)
	defer closeTestStore(t, store)

	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("EnsureSchema failed: %v", err)
	}

	if err := store.SetCurrentProposalID("", "CR26061301"); err == nil {
		t.Fatalf("expected empty project ID to be rejected")
	}
	if err := store.SetCurrentProposalID("frontend", ""); err == nil {
		t.Fatalf("expected empty proposal ID to be rejected")
	}
	if _, _, err := store.CurrentProposalID(""); err == nil {
		t.Fatalf("expected empty project ID read to be rejected")
	}
	if err := store.ClearCurrentProposalID(""); err == nil {
		t.Fatalf("expected empty project ID clear to be rejected")
	}
}

func TestStatePersistsAfterReopen(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "flowforge.sqlite")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if err := store.EnsureSchema(); err != nil {
		t.Fatalf("EnsureSchema failed: %v", err)
	}
	if err := store.SetCurrentProjectID("frontend"); err != nil {
		t.Fatalf("SetCurrentProjectID failed: %v", err)
	}
	if err := store.SetCurrentProposalID("frontend", "CR26061301"); err != nil {
		t.Fatalf("SetCurrentProposalID failed: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	reopened, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open reopen failed: %v", err)
	}
	defer closeTestStore(t, reopened)

	projectID, ok, err := reopened.CurrentProjectID()
	if err != nil {
		t.Fatalf("CurrentProjectID after reopen failed: %v", err)
	}
	if !ok || projectID != "frontend" {
		t.Fatalf("expected persisted current project frontend, got %q, ok=%v", projectID, ok)
	}

	proposalID, ok, err := reopened.CurrentProposalID("frontend")
	if err != nil {
		t.Fatalf("CurrentProposalID after reopen failed: %v", err)
	}
	if !ok || proposalID != "CR26061301" {
		t.Fatalf("expected persisted proposal CR26061301, got %q, ok=%v", proposalID, ok)
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "flowforge.sqlite")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	return store
}

func closeTestStore(t *testing.T, store *Store) {
	t.Helper()

	if err := store.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}
