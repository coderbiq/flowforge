package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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

func TestModelSetIndependentSourcesAndActiveStatus(t *testing.T) {
	root := t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	cfg := loadReasoningConfig(t, root, `agents:
  hosts: [codex,pi,opencode]
  models_by_name:
    flowforge-investigator: {model: provider/shared, reasoning_effort: high}
    flowforge-reviewer: {model: provider/reviewer, reasoning_effort: medium}
  models_by_host:
    codex:
      flowforge-investigator: gpt-base
  model_sets:
    quick:
      models_by_host:
        codex:
          flowforge-investigator: gpt-quick
      models_by_name:
        flowforge-reviewer: {reasoning_effort: low}
`)
	if err := runModelSetUse(t, "quick"); err != nil {
		t.Fatal(err)
	}
	status, err := computeSubagentStatus(root, cfg)
	if err != nil || !status.Current {
		t.Fatalf("active-set status: %#v %v", status, err)
	}
	var out bytes.Buffer
	cmd := newModelSetShowCmd()
	cmd.SetOut(&out)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var codexLine, piLine, reviewerLine string
	for _, line := range strings.Split(out.String(), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		if fields[0] == "codex" && fields[1] == "flowforge-investigator" {
			codexLine = line
		}
		if fields[0] == "pi" && fields[1] == "flowforge-investigator" {
			piLine = line
		}
		if fields[0] == "codex" && fields[1] == "flowforge-reviewer" {
			reviewerLine = line
		}
	}
	for _, want := range []string{"gpt-quick", "set:agents.model_sets.quick.models_by_host.codex.flowforge-investigator.model", "high", "base:agents.models_by_name.flowforge-investigator.reasoning_effort"} {
		if !strings.Contains(codexLine, want) {
			t.Fatalf("codex independent source %s missing: %s", want, codexLine)
		}
	}
	for _, want := range []string{"provider/shared", "base:agents.models_by_name.flowforge-investigator.model", "high"} {
		if !strings.Contains(piLine, want) {
			t.Fatalf("pi unified source %s missing: %s", want, piLine)
		}
	}
	for _, want := range []string{"provider/reviewer", "base:agents.models_by_name.flowforge-reviewer.model", "low", "set:agents.model_sets.quick.models_by_name.flowforge-reviewer.reasoning_effort"} {
		if !strings.Contains(reviewerLine, want) {
			t.Fatalf("reviewer effort-only set %s missing: %s", want, reviewerLine)
		}
	}
	if err := runModelSetUse(t, "default"); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(root, ".codex/agents/flowforge-investigator.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `model = "gpt-base"`) {
		t.Fatalf("default did not restore base: %s", body)
	}
	body, err = os.ReadFile(filepath.Join(root, ".pi/agents/flowforge-reviewer.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "thinking: medium") || !strings.Contains(string(body), "model: provider/reviewer") {
		t.Fatalf("default did not restore base effort/model: %s", body)
	}
}
func TestModelSetRollbackKeepsArtifactsOnPreparationFailure(t *testing.T) {
	for _, failure := range []string{"config", "pin"} {
		t.Run(failure, func(t *testing.T) {
			root := t.TempDir()
			if err := initializeTestProject(root); err != nil {
				t.Fatal(err)
			}
			t.Chdir(root)
			cfg := loadReasoningConfig(t, root, `agents:
  hosts: [codex]
  models_by_name:
    flowforge-investigator: {model: gpt-base, reasoning_effort: high}
  model_sets:
    previous:
      models_by_name:
        flowforge-investigator: gpt-previous
    next:
      models_by_name:
        flowforge-investigator: {model: gpt-next, reasoning_effort: low}
`)
			if err := runModelSetUse(t, "previous"); err != nil {
				t.Fatal(err)
			}
			if failure == "config" {
				cfg.Agents.ModelSets["next"] = config.ModelSetConfig{ModelOverrides: map[string]config.ModelValue{"flowforge-investigator": {Model: "bad token"}}}
				if err := cfg.Save(root); err != nil {
					t.Fatal(err)
				}
			} else {
				path := filepath.Join(root, ".codex/agents/flowforge-scribe.toml")
				if err := os.WriteFile(path, []byte("model = [\n"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			before := snapshotAgentArtifacts(t, root)
			if err := runModelSetUse(t, "next"); err == nil {
				t.Fatal("failed preparation must reject switch")
			}
			if name, err := config.ReadActiveModelSet(root); err != nil || name != "previous" {
				t.Fatalf("pointer not restored: %s %v", name, err)
			}
			assertArtifactsEqual(t, before, snapshotAgentArtifacts(t, root))
		})
	}
}
func TestUndeclaredSetDiagnosticParityAndShowDoesNotReadPins(t *testing.T) {
	root := t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	cfg := loadReasoningConfig(t, root, "agents:\n  hosts: [codex]\n")
	if _, err := deploySubagents(root, cfg, ""); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".codex/agents/flowforge-investigator.toml")
	if err := os.WriteFile(path, []byte("model = \"local-only\"\nmodel_reasoning_effort = \"ultra\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd := newModelSetShowCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"default"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "local-only") || strings.Contains(out.String(), "ultra") {
		t.Fatal("show mislabeled local pin as config")
	}
	if err := config.WriteActiveModelSet(root, "missing"); err != nil {
		t.Fatal(err)
	}
	before := snapshotAgentArtifacts(t, root)
	_, depErr := deploySubagents(root, cfg, "")
	_, statusErr := computeSubagentStatus(root, cfg)
	if depErr == nil || statusErr == nil || depErr.Error() != statusErr.Error() {
		t.Fatalf("undeclared active set parity: %v / %v", depErr, statusErr)
	}
	assertArtifactsEqual(t, before, snapshotAgentArtifacts(t, root))
}

func TestModelSetRestoresUnconfiguredDefaults(t *testing.T) {
	root := t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	cfg := loadReasoningConfig(t, root, `agents:
  hosts: [codex]
  model_sets:
    quick:
      models_by_name:
        flowforge-reviewer: {reasoning_effort: low}
`)
	if _, err := deploySubagents(root, cfg, ""); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".codex/agents/flowforge-reviewer.toml")
	for _, step := range []struct{ set, want string }{{"quick", "low"}, {"default", "high"}} {
		if err := runModelSetUse(t, step.set); err != nil {
			t.Fatal(err)
		}
		value, err := readLocalModelFields(path, "codex")
		if err != nil {
			t.Fatal(err)
		}
		if value.ReasoningEffort != step.want {
			t.Fatalf("%s effort = %q, want %q", step.set, value.ReasoningEffort, step.want)
		}
	}
}
