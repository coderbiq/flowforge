package config

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ModelSetConfig is one named overlay over the base agent model layers
// (agents.models / models_by_name / models_by_host). Only the keys present
// here override the base; everything else is inherited from the base.
type ModelSetConfig struct {
	Models             map[string]string            `yaml:"models,omitempty" mapstructure:"models"`
	ModelOverrides     map[string]string            `yaml:"models_by_name,omitempty" mapstructure:"models_by_name"`
	ModelHostOverrides map[string]map[string]string `yaml:"models_by_host,omitempty" mapstructure:"models_by_host"`
}

// activeModelSetFile holds the per-machine active model-set pointer inside
// the .flowforge directory. Absence (or empty content) means: base layers.
const activeModelSetFile = "model-set.active"

// ActiveModelSetPath returns the state-file path for a project root.
func ActiveModelSetPath(projectRoot string) string {
	return filepath.Join(projectRoot, ConfigDirName, activeModelSetFile)
}

// ReadActiveModelSet returns the active model-set name, or "" when no set
// is active (base layers in effect). A missing file is not an error.
func ReadActiveModelSet(projectRoot string) (string, error) {
	data, err := os.ReadFile(ActiveModelSetPath(projectRoot))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	name := strings.TrimSpace(string(data))
	if name == "" {
		return "", nil
	}
	return name, nil
}

// WriteActiveModelSet records the active model-set pointer.
func WriteActiveModelSet(projectRoot, name string) error {
	p := ActiveModelSetPath(projectRoot)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(name+"\n"), 0o644)
}

// ClearActiveModelSet removes the pointer file (back to base layers).
// A missing file is not an error.
func ClearActiveModelSet(projectRoot string) error {
	err := os.Remove(ActiveModelSetPath(projectRoot))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ResolveModelSet looks up a model set by name. The empty name never
// resolves; callers treat that as "base layers".
func ResolveModelSet(cfg *AgentsConfig, name string) (ModelSetConfig, bool) {
	if name == "" {
		return ModelSetConfig{}, false
	}
	set, ok := cfg.ModelSets[name]
	return set, ok
}

// ModelSetNames returns the declared model-set names in sorted order.
func ModelSetNames(cfg *AgentsConfig) []string {
	return slices.Sorted(maps.Keys(cfg.ModelSets))
}

// mergeSS overlays set over base (string→string maps).
func mergeSS(base, set map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(set))
	maps.Copy(merged, base)
	maps.Copy(merged, set)
	return merged
}

// ApplyModelSet returns a copy of base with the set overlaid. base is not
// mutated. Layer semantics:
//   - models / models_by_name: set key overrides the same base key
//   - models_by_host: merged per host; set keys override base keys per host
func ApplyModelSet(base *AgentsConfig, set ModelSetConfig) AgentsConfig {
	merged := *base
	merged.Models = mergeSS(base.Models, set.Models)
	merged.ModelOverrides = mergeSS(base.ModelOverrides, set.ModelOverrides)
	hosts := make(map[string]map[string]string, len(base.ModelHostOverrides)+len(set.ModelHostOverrides))
	for host, inner := range base.ModelHostOverrides {
		hosts[host] = mergeSS(inner, nil)
	}
	for host, inner := range set.ModelHostOverrides {
		hosts[host] = mergeSS(hosts[host], inner)
	}
	merged.ModelHostOverrides = hosts
	return merged
}

// EffectiveAgents returns the AgentsConfig with the active model set
// overlaid (base layers when no set is active). The returned struct is a
// shallow overlay copy; non-model fields share storage with cfg.
func EffectiveAgents(cfg *Config, projectRoot string) (AgentsConfig, string, error) {
	name, err := ReadActiveModelSet(projectRoot)
	if err != nil {
		return AgentsConfig{}, "", fmt.Errorf("reading active model set: %w", err)
	}
	if name == "" {
		return cfg.Agents, "", nil
	}
	set, ok := ResolveModelSet(&cfg.Agents, name)
	if !ok {
		return cfg.Agents, "", fmt.Errorf("active model set %q is not declared in agents.model_sets", name)
	}
	return ApplyModelSet(&cfg.Agents, set), name, nil
}
