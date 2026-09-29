package command

import (
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/config"
)

// writeModelSetConfig appends a model_sets declaration to the project's
// config.yaml and returns the reloaded config.
func writeModelSetConfig(t *testing.T, projectRoot, extra string) *config.Config {
	t.Helper()
	p := filepath.Join(projectRoot, config.ConfigDirName, config.ConfigFileName)
	base, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, append(base, []byte(extra)...), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

// use <set> with a valid set must write the pointer, redeploy, and the
// deployed artifact must carry the set's model.
func TestModelSetUseSwitchesDeployedModel(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg := writeModelSetConfig(t, projectRoot, `agents:
  model_sets:
    offpeak:
      models_by_name:
        flowforge-reviewer: prov/mimo
        flowforge-planner: prov/mimo
`)

	if name, _ := config.ReadActiveModelSet(projectRoot); name != "" {
		t.Fatalf("fresh project should have no active set, got %q", name)
	}

	deployed, err := deploySubagents(projectRoot, cfg, "")
	if err != nil {
		t.Fatalf("deploy under set: %v", err)
	}
	if len(deployed) == 0 {
		t.Fatal("no subagents deployed")
	}

	// activate the set, then deploy again via the overlay path
	if err := config.WriteActiveModelSet(projectRoot, "offpeak"); err != nil {
		t.Fatal(err)
	}
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		t.Fatalf("redeploy under active set: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(projectRoot, ".pi", "agents", "flowforge-reviewer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(data), "prov/mimo") {
		t.Error("deployed flowforge-reviewer.md should carry the set model prov/mimo")
	}
	if name, _ := config.ReadActiveModelSet(projectRoot); name != "offpeak" {
		t.Errorf("pointer should read offpeak, got %q", name)
	}
}

// EffectiveAgents: active set that is not declared must fail loudly — the
// deploy refuses rather than silently using base layers.
func TestDeployFailsOnUndeclaredActiveSet(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := config.WriteActiveModelSet(projectRoot, "ghost"); err != nil {
		t.Fatal(err)
	}
	if _, err := deploySubagents(projectRoot, cfg, ""); err == nil {
		t.Fatal("deploy with undeclared active set must fail")
	}
}

// An invalid model set (unknown agent key) must be rejected by validation
// with an error path naming the set.
func TestValidateModelConfigRejectsInvalidSet(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	cfg := writeModelSetConfig(t, projectRoot, `agents:
  model_sets:
    broken:
      models_by_name:
        no-such-agent: prov/x
`)
	hosts, err := resolveHostTargets(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defs, err := discoverSubagentSources(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	err = validateModelConfig(cfg, defs, hosts)
	if err == nil {
		t.Fatal("invalid set must fail validation")
	}
	if !contains(err.Error(), "model set \"broken\"") || !contains(err.Error(), "no-such-agent") {
		t.Fatalf("error should name the set and the bad key, got: %v", err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}

// runModelSetUse executes the `model-set use` RunE with one arg.
func runModelSetUse(t *testing.T, arg string) error {
	t.Helper()
	cmd := newModelSetUseCmd()
	cmd.SetArgs([]string{arg})
	cmd.SilenceUsage = true
	return cmd.Execute()
}

// use with a set that fails deploy validation must roll the pointer back
// to its previous value (atomic switch, design switch-atomicity).
func TestModelSetUseRollsBackPointerOnDeployFailure(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	t.Chdir(projectRoot) // RunE roots at cwd via FindProjectRoot
	writeModelSetConfig(t, projectRoot, `agents:
  model_sets:
    broken:
      models_by_name:
        no-such-agent: prov/x
`)

	// prev is empty: failed use must leave no pointer behind
	err := runModelSetUse(t, "broken")
	if err == nil {
		t.Fatal("use with invalid set must fail")
	}
	if name, rerr := config.ReadActiveModelSet(projectRoot); rerr != nil || name != "" {
		t.Fatalf("pointer must roll back to empty, got %q (err=%v)", name, rerr)
	}
}

// use default on a fresh project is a no-op success (idempotent base state).
func TestModelSetUseDefaultIsIdempotent(t *testing.T) {
	projectRoot := t.TempDir()
	if err := initializeTestProject(projectRoot); err != nil {
		t.Fatal(err)
	}
	t.Chdir(projectRoot)
	if err := runModelSetUse(t, "default"); err != nil {
		t.Fatalf("use default: %v", err)
	}
	if name, _ := config.ReadActiveModelSet(projectRoot); name != "" {
		t.Fatalf("default must keep pointer empty, got %q", name)
	}
}
