package command

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"flowforge/internal/config"
)

func modelStateProject(t *testing.T, hosts string) (string, *config.Config) {
	t.Helper()
	root := t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	cfg := loadReasoningConfig(t, root, "agents:\n  hosts: ["+hosts+"]\n")
	if _, err := deploySubagents(root, cfg, ""); err != nil {
		t.Fatal(err)
	}
	return root, cfg
}
func deployModelState(t *testing.T, root string, cfg *config.Config, target string) {
	t.Helper()
	if _, err := deploySubagents(root, cfg, target); err != nil {
		t.Fatal(err)
	}
}
func assertNativeModel(t *testing.T, root, host, name, model, effort string) {
	t.Helper()
	var path string
	for _, h := range allHostTargets() {
		if h.key == host {
			path = filepath.Join(root, h.relDir, name+h.ext)
		}
	}
	value, err := readLocalModelFields(path, host)
	if err != nil {
		t.Fatal(err)
	}
	if value.Model != model || value.ReasoningEffort != effort {
		t.Fatalf("%s/%s got %#v, want model=%q effort=%q", host, name, value, model, effort)
	}
}
func editNativeModel(t *testing.T, root, host, name, model, effort string) {
	t.Helper()
	for _, h := range allHostTargets() {
		if h.key == host {
			var content string
			switch host {
			case "codex":
				if model != "" {
					content += "model = \"" + model + "\"\n"
				}
				if effort != "" {
					content += "model_reasoning_effort = \"" + effort + "\"\n"
				}
			default:
				content = "---\n"
				if model != "" {
					content += "model: " + model + "\n"
				}
				key := map[string]string{"claude": "effort", "opencode": "reasoningEffort", "pi": "thinking"}[host]
				if effort != "" {
					content += key + ": " + effort + "\n"
				}
				content += "---\nlocal body\n"
			}
			if err := os.WriteFile(filepath.Join(root, h.relDir, name+h.ext), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
}
func snapshotModelStateProject(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		result[rel] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func assertModelStateProjectUnchanged(t *testing.T, before map[string]string, root string) {
	t.Helper()
	if after := snapshotModelStateProject(t, root); !reflect.DeepEqual(before, after) {
		t.Fatal("preparation/status changed project files")
	}
}

func TestModelSetRestoresRetainedLocalPins(t *testing.T) {
	for _, host := range []string{"codex", "claude", "pi", "opencode"} {
		t.Run(host, func(t *testing.T) {
			root, cfg := modelStateProject(t, host)
			editNativeModel(t, root, host, "flowforge-reviewer", "provider/local", "medium")
			cfg.Agents.ModelSets = map[string]config.ModelSetConfig{"quick": {ModelOverrides: map[string]config.ModelValue{"flowforge-reviewer": {Model: "provider/quick", ReasoningEffort: "low"}}}}
			if err := config.WriteActiveModelSet(root, "quick"); err != nil {
				t.Fatal(err)
			}
			deployModelState(t, root, cfg, "")
			assertNativeModel(t, root, host, "flowforge-reviewer", "provider/quick", "low")
			deployModelState(t, root, cfg, "")
			if err := config.WriteActiveModelSet(root, ""); err != nil {
				t.Fatal(err)
			}
			deployModelState(t, root, cfg, "")
			assertNativeModel(t, root, host, "flowforge-reviewer", "provider/local", "medium")
			deployModelState(t, root, cfg, "")
			assertNativeModel(t, root, host, "flowforge-reviewer", "provider/local", "medium")
		})
	}
}
func TestModelSetRestoresFourHostDefaults(t *testing.T) {
	root, cfg := modelStateProject(t, "codex,claude,pi,opencode")
	cfg.Agents.ModelSets = map[string]config.ModelSetConfig{
		"a": {ModelOverrides: map[string]config.ModelValue{"flowforge-reviewer": {Model: "provider/quick", ReasoningEffort: "low"}}},
		"b": {ModelOverrides: map[string]config.ModelValue{"flowforge-reviewer": {Model: "provider/next"}}},
	}
	if err := config.WriteActiveModelSet(root, "a"); err != nil {
		t.Fatal(err)
	}
	deployModelState(t, root, cfg, "")
	if err := config.WriteActiveModelSet(root, "b"); err != nil {
		t.Fatal(err)
	}
	deployModelState(t, root, cfg, "")
	for _, host := range []string{"codex", "claude", "pi", "opencode"} {
		effort := ""
		if host == "codex" || host == "pi" {
			effort = "high"
		}
		assertNativeModel(t, root, host, "flowforge-reviewer", "provider/next", effort)
	}
	if err := config.WriteActiveModelSet(root, ""); err != nil {
		t.Fatal(err)
	}
	deployModelState(t, root, cfg, "")
	for _, host := range []string{"codex", "claude", "pi", "opencode"} {
		model, effort := "", ""
		if host == "claude" {
			model = "opus"
		}
		if host == "codex" || host == "pi" {
			effort = "high"
		}
		assertNativeModel(t, root, host, "flowforge-reviewer", model, effort)
	}
}
func TestDeployModelStateAfterConfigEdits(t *testing.T) {
	root, cfg := modelStateProject(t, "codex")
	cfg.Agents.ModelOverrides = map[string]config.ModelValue{"flowforge-reviewer": {Model: "gpt-old", ReasoningEffort: "low"}}
	deployModelState(t, root, cfg, "")
	delete(cfg.Agents.ModelOverrides, "flowforge-reviewer")
	before := snapshotModelStateProject(t, root)
	status, err := computeSubagentStatus(root, cfg)
	if err != nil || status.Current {
		t.Fatalf("expected drift: %#v %v", status, err)
	}
	assertModelStateProjectUnchanged(t, before, root)
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "", "high")
	// Changing the definition's profile updates defaults rather than preserving
	// the prior profile's generated high effort.
	assets, cleanup, err := locateAssetsDir()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	data, err := os.ReadFile(filepath.Join(assets, "subagents/flowforge-reviewer.md"))
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.ReplaceAll(data, []byte("model_profile: high-capability"), []byte("model_profile: tool-capable"))
	dir := filepath.Join(root, ".flowforge/subagents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "flowforge-reviewer.md"), data, 0644); err != nil {
		t.Fatal(err)
	}
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "", "medium")
	cfg.Agents.ModelSets = map[string]config.ModelSetConfig{"quick": {ModelOverrides: map[string]config.ModelValue{"flowforge-reviewer": {ReasoningEffort: "low"}}}}
	if err := config.WriteActiveModelSet(root, "quick"); err != nil {
		t.Fatal(err)
	}
	deployModelState(t, root, cfg, "")
	cfg.Agents.ModelSets["quick"] = config.ModelSetConfig{}
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "", "medium")
}
func TestDeployModelStateLocalFieldChanges(t *testing.T) {
	root, cfg := modelStateProject(t, "codex")
	editNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-local", "ultra")
	deployModelState(t, root, cfg, "")
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-local", "ultra")
	editNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-changed", "ultra")
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-changed", "ultra")
	editNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-changed", "")
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-changed", "high")
	editNativeModel(t, root, "codex", "flowforge-reviewer", "", "high")
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "", "high")
	editNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-local", "ultra")
	deployModelState(t, root, cfg, "")
	if err := os.Remove(filepath.Join(root, ".codex/agents/flowforge-reviewer.toml")); err != nil {
		t.Fatal(err)
	}
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "", "high")
}
func TestDeployModelStateInherit(t *testing.T) {
	for _, local := range []bool{false, true} {
		t.Run(map[bool]string{false: "default", true: "local"}[local], func(t *testing.T) {
			root, cfg := modelStateProject(t, "codex")
			if local {
				editNativeModel(t, root, "codex", "flowforge-reviewer", "", "ultra")
			}
			cfg.Agents.ModelOverrides = map[string]config.ModelValue{"flowforge-reviewer": {ReasoningEffort: "inherit"}}
			deployModelState(t, root, cfg, "")
			deployModelState(t, root, cfg, "")
			assertNativeModel(t, root, "codex", "flowforge-reviewer", "", "")
			delete(cfg.Agents.ModelOverrides, "flowforge-reviewer")
			deployModelState(t, root, cfg, "")
			want := "high"
			if local {
				want = "ultra"
			}
			assertNativeModel(t, root, "codex", "flowforge-reviewer", "", want)
		})
	}
}
func TestDeployModelStateLegacyMigration(t *testing.T) {
	root, cfg := modelStateProject(t, "codex")
	if err := os.Remove(agentModelStatePath(root)); err != nil {
		t.Fatal(err)
	}
	editNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-legacy", "ultra")
	cfg.Agents.ModelOverrides = map[string]config.ModelValue{"flowforge-reviewer": {Model: "gpt-explicit", ReasoningEffort: "low"}}
	deployModelState(t, root, cfg, "")
	delete(cfg.Agents.ModelOverrides, "flowforge-reviewer")
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "gpt-legacy", "ultra")
}
func TestDeployModelStateTargetIsolation(t *testing.T) {
	root, cfg := modelStateProject(t, "codex,pi")
	before, err := readAgentModelState(root)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Agents.ModelOverrides = map[string]config.ModelValue{"flowforge-reviewer": {ReasoningEffort: "low"}}
	deployModelState(t, root, cfg, "flowforge-reviewer")
	after, err := readAgentModelState(root)
	if err != nil {
		t.Fatal(err)
	}
	for path, record := range before.Paths {
		if !strings.Contains(path, "flowforge-reviewer.") && !reflect.DeepEqual(record, after.Paths[path]) {
			t.Fatalf("target deploy changed %s", path)
		}
	}
	cfg.Agents.Hosts = []string{"codex"}
	deployModelState(t, root, cfg, "flowforge-reviewer")
	after, err = readAgentModelState(root)
	if err != nil {
		t.Fatal(err)
	}
	for path := range after.Paths {
		if strings.HasPrefix(path, ".pi/") {
			t.Fatalf("cleaned path retained: %s", path)
		}
	}
}
func TestModelStatePreparationFailureIsReadOnly(t *testing.T) {
	failures := map[string]string{
		"broken": "{", "version": `{"version":2,"paths":{}}`, "unknown": `{"version":1,"paths":{},"extra":true}`,
		"null":             `{"version":1,"paths":{".codex/agents/a.toml":null}}`,
		"invalid_field":    `{"version":1,"paths":{".codex/agents/a.toml":{"model":{"last_generated":{"present":"yes","value":"x"}},"effort":{"last_generated":{"present":false,"value":""}}}}}`,
		"missing_presence": `{"version":1,"paths":{".codex/agents/a.toml":{"model":{"last_generated":{"value":"x"}},"effort":{"last_generated":{"present":false,"value":""}}}}}`,
		"config":           "", "late_pin": "",
	}
	for name, bad := range failures {
		t.Run(name, func(t *testing.T) {
			root, cfg := modelStateProject(t, "codex")
			cfg.Agents.ModelSets = map[string]config.ModelSetConfig{"quick": {ModelOverrides: map[string]config.ModelValue{"flowforge-reviewer": {ReasoningEffort: "low"}}}}
			if err := cfg.Save(root); err != nil {
				t.Fatal(err)
			}
			switch name {
			case "config":
				cfg.Agents.ModelOverrides = map[string]config.ModelValue{"flowforge-reviewer": {ReasoningEffort: "bad token"}}
				if err := cfg.Save(root); err != nil {
					t.Fatal(err)
				}
			case "late_pin":
				if err := os.WriteFile(filepath.Join(root, ".codex/agents/flowforge-scribe.toml"), []byte("model = [\n"), 0644); err != nil {
					t.Fatal(err)
				}
			default:
				if err := os.WriteFile(agentModelStatePath(root), []byte(bad), 0644); err != nil {
					t.Fatal(err)
				}
			}
			before := snapshotModelStateProject(t, root)
			_, depErr := deploySubagents(root, cfg, "")
			_, statusErr := computeSubagentStatus(root, cfg)
			if depErr == nil || statusErr == nil || depErr.Error() != statusErr.Error() {
				t.Fatalf("diagnostic mismatch: %v / %v", depErr, statusErr)
			}
			if bad != "" && !strings.Contains(depErr.Error(), agentModelStatePath(root)) {
				t.Fatalf("missing path: %v", depErr)
			}
			assertModelStateProjectUnchanged(t, before, root)
			t.Chdir(root)
			if err := runModelSetUse(t, "quick"); err == nil {
				t.Fatal("use should fail")
			}
			if active, err := config.ReadActiveModelSet(root); err != nil || active != "" {
				t.Fatalf("pointer changed: %q %v", active, err)
			}
			assertModelStateProjectUnchanged(t, before, root)
		})
	}
}
func TestModelStateStatusIsReadOnly(t *testing.T) {
	root, cfg := modelStateProject(t, "codex")
	for _, stage := range []string{"normal", "missing", "drift"} {
		t.Run(stage, func(t *testing.T) {
			if stage == "missing" {
				if err := os.Remove(agentModelStatePath(root)); err != nil {
					t.Fatal(err)
				}
			}
			if stage == "drift" {
				cfg.Agents.ModelOverrides = map[string]config.ModelValue{"flowforge-reviewer": {ReasoningEffort: "low"}}
				deployModelState(t, root, cfg, "")
				delete(cfg.Agents.ModelOverrides, "flowforge-reviewer")
			}
			before := snapshotModelStateProject(t, root)
			status, err := computeSubagentStatus(root, cfg)
			if err != nil {
				t.Fatal(err)
			}
			if stage == "drift" && status.Current {
				t.Fatal("removed config should drift")
			}
			assertModelStateProjectUnchanged(t, before, root)
		})
	}
}
func TestModelStateAtomicWriteFailure(t *testing.T) {
	root := t.TempDir()
	path := agentModelStatePath(root)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(path, "old.json")
	if err := os.WriteFile(old, []byte("old snapshot"), 0644); err != nil {
		t.Fatal(err)
	}
	err := writeAgentModelState(root, []byte("new snapshot"))
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("expected path error: %v", err)
	}
	data, err := os.ReadFile(old)
	if err != nil || string(data) != "old snapshot" {
		t.Fatal("old snapshot changed")
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".agent-model-state-") {
			t.Fatal("temporary state file leaked")
		}
	}
}

func TestDeployModelStateDeletedFileDuringInherit(t *testing.T) {
	root, cfg := modelStateProject(t, "codex")
	editNativeModel(t, root, "codex", "flowforge-reviewer", "", "ultra")
	cfg.Agents.ModelOverrides = map[string]config.ModelValue{"flowforge-reviewer": {ReasoningEffort: "inherit"}}
	deployModelState(t, root, cfg, "")
	if err := os.Remove(filepath.Join(root, ".codex/agents/flowforge-reviewer.toml")); err != nil {
		t.Fatal(err)
	}
	deployModelState(t, root, cfg, "")
	delete(cfg.Agents.ModelOverrides, "flowforge-reviewer")
	deployModelState(t, root, cfg, "")
	assertNativeModel(t, root, "codex", "flowforge-reviewer", "", "high")
}

func TestModelStateInvalidRecordDiagnosticIsStable(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".flowforge"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(agentModelStatePath(root), []byte(`{"version":1,"paths":{"../a":null,"../b":null}}`), 0644); err != nil {
		t.Fatal(err)
	}
	for attempt := 0; attempt < 100; attempt++ {
		_, err := readAgentModelState(root)
		if err == nil || !strings.Contains(err.Error(), `invalid record path "../a"`) {
			t.Fatalf("same invalid snapshot must report ../a first, attempt %d: %v", attempt, err)
		}
	}
}
