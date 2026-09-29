package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyModelSetOverridesAllThreeLayers(t *testing.T) {
	base := AgentsConfig{
		Models:         map[string]string{"tool-capable": "cpa/flash-a"},
		ModelOverrides: map[string]string{"flowforge-reviewer": "cpa/glm", "flowforge-planner": "cpa/glm"},
		ModelHostOverrides: map[string]map[string]string{
			"pi": {"flowforge-scribe": "cpa/flash-b"},
		},
	}
	set := ModelSetConfig{
		ModelOverrides: map[string]string{"flowforge-reviewer": "cpa/mimo"},
		ModelHostOverrides: map[string]map[string]string{
			"pi": {"flowforge-scribe": "cpa/flash-c"},
		},
	}

	got := ApplyModelSet(&base, set)

	if got.ModelOverrides["flowforge-reviewer"] != "cpa/mimo" {
		t.Errorf("set should override models_by_name: got %q", got.ModelOverrides["flowforge-reviewer"])
	}
	if got.ModelOverrides["flowforge-planner"] != "cpa/glm" {
		t.Errorf("base keys must survive: got %q", got.ModelOverrides["flowforge-planner"])
	}
	if got.Models["tool-capable"] != "cpa/flash-a" {
		t.Errorf("untouched base layer must survive: got %q", got.Models["tool-capable"])
	}
	if got.ModelHostOverrides["pi"]["flowforge-scribe"] != "cpa/flash-c" {
		t.Errorf("set should deep-override models_by_host: got %q", got.ModelHostOverrides["pi"]["flowforge-scribe"])
	}
	if base.ModelOverrides["flowforge-reviewer"] != "cpa/glm" {
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
			"offpeak": {ModelOverrides: map[string]string{"flowforge-reviewer": "cpa/mimo"}},
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
	if set.ModelOverrides["flowforge-reviewer"] != "cpa/mimo-2.6-pro" {
		t.Fatalf("models_by_name in set: %+v", set)
	}
	if set.Models["tool-capable-read-only"] != "cpa/flash-x" {
		t.Fatalf("models in set: %+v", set)
	}
}
