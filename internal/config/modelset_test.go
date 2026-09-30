package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyModelSetOverridesAllThreeLayers(t *testing.T) {
	base := AgentsConfig{
		Models:         map[string]ModelValue{"tool-capable": ModelValue{Model: "cpa/flash-a"}},
		ModelOverrides: map[string]ModelValue{"flowforge-reviewer": ModelValue{Model: "cpa/glm"}, "flowforge-planner": ModelValue{Model: "cpa/glm"}},
		ModelHostOverrides: map[string]map[string]ModelValue{
			"pi": {"flowforge-scribe": {Model: "cpa/flash-b"}},
		},
	}
	set := ModelSetConfig{
		ModelOverrides: map[string]ModelValue{"flowforge-reviewer": ModelValue{Model: "cpa/mimo"}},
		ModelHostOverrides: map[string]map[string]ModelValue{
			"pi": {"flowforge-scribe": {Model: "cpa/flash-c"}},
		},
	}

	got := ApplyModelSet(&base, set)

	if got.ModelOverrides["flowforge-reviewer"].Model != "cpa/mimo" {
		t.Errorf("set should override models_by_name: got %q", got.ModelOverrides["flowforge-reviewer"])
	}
	if got.ModelOverrides["flowforge-planner"].Model != "cpa/glm" {
		t.Errorf("base keys must survive: got %q", got.ModelOverrides["flowforge-planner"])
	}
	if got.Models["tool-capable"].Model != "cpa/flash-a" {
		t.Errorf("untouched base layer must survive: got %q", got.Models["tool-capable"])
	}
	if got.ModelHostOverrides["pi"]["flowforge-scribe"].Model != "cpa/flash-c" {
		t.Errorf("set should deep-override models_by_host: got %q", got.ModelHostOverrides["pi"]["flowforge-scribe"])
	}
	if base.ModelOverrides["flowforge-reviewer"].Model != "cpa/glm" {
		t.Errorf("ApplyModelSet must not mutate the base config")
	}
}

func TestActiveModelSetRoundTrip(t *testing.T) {
	dir := t.TempDir()

	if name, err := ReadActiveModelSet(dir); err != nil || name != "" {
		t.Fatalf("missing pointer should read as empty: got %q, err=%v", name, err)
	}
	if err := WriteActiveModelSet(dir, "offpeak"); err != nil {
		t.Fatalf("write: %v", err)
	}
	if name, err := ReadActiveModelSet(dir); err != nil || name != "offpeak" {
		t.Fatalf("round trip: got %q, err=%v", name, err)
	}
	if err := ClearActiveModelSet(dir); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if name, err := ReadActiveModelSet(dir); err != nil || name != "" {
		t.Fatalf("cleared pointer should read as empty: got %q, err=%v", name, err)
	}
}

func TestResolveModelSetByName(t *testing.T) {
	cfg := &AgentsConfig{
		ModelSets: map[string]ModelSetConfig{
			"offpeak": {ModelOverrides: map[string]ModelValue{"flowforge-reviewer": ModelValue{Model: "cpa/mimo"}}},
		},
	}
	if _, ok := ResolveModelSet(cfg, "offpeak"); !ok {
		t.Fatal("known set must resolve")
	}
	if _, ok := ResolveModelSet(cfg, "missing"); ok {
		t.Fatal("unknown set must not resolve")
	}
	if _, ok := ResolveModelSet(cfg, ""); ok {
		t.Fatal("empty name must not resolve")
	}
}

func TestModelSetNamesSorted(t *testing.T) {
	cfg := &AgentsConfig{
		ModelSets: map[string]ModelSetConfig{
			"zeta": {}, "alpha": {}, "mid": {},
		},
	}
	names := ModelSetNames(cfg)
	want := []string{"alpha", "mid", "zeta"}
	if len(names) != len(want) {
		t.Fatalf("got %v", names)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("want sorted %v, got %v", want, names)
		}
	}
}

// WriteActiveModelSet must trim whitespace so a stray newline in the state
// file never becomes part of the set name.
func TestReadActiveModelSetTrims(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ConfigDirName)
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p, activeModelSetFile), []byte("  offpeak\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	name, err := ReadActiveModelSet(dir)
	if err != nil || name != "offpeak" {
		t.Fatalf("trim: got %q, err=%v", name, err)
	}
}

func TestLoadConfigParsesModelSets(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ConfigDirName), 0o755); err != nil {
		t.Fatal(err)
	}
	yml := `version: 2.0.0
docs_dir: docs
agents:
  hosts: [pi]
  model_sets:
    offpeak:
      models_by_name:
        flowforge-reviewer: cpa/mimo-2.6-pro
      models:
        tool-capable-read-only: cpa/flash-x
`
	if err := os.WriteFile(filepath.Join(dir, ConfigDirName, ConfigFileName), []byte(yml), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	set, ok := ResolveModelSet(&cfg.Agents, "offpeak")
	if !ok {
		t.Fatal("model_sets.offpeak missing after load")
	}
	if set.ModelOverrides["flowforge-reviewer"].Model != "cpa/mimo-2.6-pro" {
		t.Fatalf("models_by_name in set: %+v", set)
	}
	if set.Models["tool-capable-read-only"].Model != "cpa/flash-x" {
		t.Fatalf("models in set: %+v", set)
	}
}

func TestApplyModelSetIndependentFieldsAndSparseHosts(t *testing.T) {
	base := AgentsConfig{Models: map[string]ModelValue{"tool-capable": {Model: "provider/base", ReasoningEffort: "high"}}, ModelOverrides: map[string]ModelValue{"flowforge-reviewer": {Model: "provider/reviewer", ReasoningEffort: "medium"}}, ModelHostOverrides: map[string]map[string]ModelValue{"codex": {"flowforge-investigator": {Model: "gpt-base", ReasoningEffort: "high"}, "flowforge-reviewer": {Model: "gpt-review"}}}}
	set := ModelSetConfig{Models: map[string]ModelValue{"tool-capable": {Model: "provider/new"}}, ModelOverrides: map[string]ModelValue{"flowforge-reviewer": {ReasoningEffort: "low"}}, ModelHostOverrides: map[string]map[string]ModelValue{"codex": {"flowforge-investigator": {ReasoningEffort: "inherit"}}}}
	got := ApplyModelSet(&base, set)
	if got.Models["tool-capable"] != (ModelValue{Model: "provider/new", ReasoningEffort: "high"}) {
		t.Fatal(got.Models)
	}
	if got.ModelOverrides["flowforge-reviewer"] != (ModelValue{Model: "provider/reviewer", ReasoningEffort: "low"}) {
		t.Fatal(got.ModelOverrides)
	}
	if got.ModelHostOverrides["codex"]["flowforge-investigator"] != (ModelValue{Model: "gpt-base", ReasoningEffort: "inherit"}) {
		t.Fatal(got.ModelHostOverrides)
	}
	if got.ModelHostOverrides["codex"]["flowforge-reviewer"].Model != "gpt-review" || len(got.ModelHostOverrides) != 1 {
		t.Fatal("sparse hosts lost")
	}
	got.ModelHostOverrides["codex"]["flowforge-reviewer"] = ModelValue{Model: "changed"}
	if base.ModelHostOverrides["codex"]["flowforge-reviewer"].Model != "gpt-review" || base.Models["tool-capable"].Model != "provider/base" || set.ModelHostOverrides["codex"]["flowforge-investigator"].Model != "" {
		t.Fatal("inputs mutated")
	}
}
