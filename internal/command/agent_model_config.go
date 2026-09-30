package command

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"flowforge/internal/config"
	"flowforge/internal/subagent"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// resolveModelField resolves each field independently after same-layer overlay.
func resolveModelField(a *config.AgentsConfig, def *subagent.Definition, host string, effort bool) (string, string) {
	profile := string(def.ModelProfile)
	layers := []struct {
		value config.ModelValue
		path  string
	}{
		{a.ModelHostOverrides[host][def.Name], "models_by_host." + host + "." + def.Name},
		{a.ModelHostOverrides[host][profile], "models_by_host." + host + "." + profile},
		{a.ModelOverrides[def.Name], "models_by_name." + def.Name},
		{a.Models[profile], "models." + profile},
	}
	for _, layer := range layers {
		value := layer.value.Model
		field := "model"
		if effort {
			value = layer.value.ReasoningEffort
			field = "reasoning_effort"
		}
		if value != "" {
			return value, "agents." + layer.path + "." + field
		}
	}
	return "", ""
}
func validateModelConfig(cfg *config.Config, defs []*subagent.Definition, hosts []hostTarget) error {
	known := map[string]bool{}
	for _, def := range defs {
		known[def.Name] = true
	}
	knownHosts := map[string]bool{}
	for _, h := range allHostTargets() {
		knownHosts[h.key] = true
	}
	checkDeclarations := func(a *config.AgentsConfig, path string) error {
		check := func(value config.ModelValue, path string) error {
			if value.Model == "" && value.ReasoningEffort == "" {
				return fmt.Errorf("%s: model or reasoning_effort must be non-empty", path)
			}
			for _, field := range []struct{ name, value string }{{"model", value.Model}, {"reasoning_effort", value.ReasoningEffort}} {
				for _, r := range field.value {
					if unicode.IsSpace(r) || unicode.IsControl(r) {
						return fmt.Errorf("%s.%s: value %q must not contain whitespace or control characters", path, field.name, field.value)
					}
				}
			}
			return nil
		}
		for _, key := range slices.Sorted(maps.Keys(a.Models)) {
			if !validModelProfileKeys[key] {
				return fmt.Errorf("%s.models: unknown profile key %q (path: %s.models.%s)", path, key, path, key)
			}
			if err := check(a.Models[key], path+".models."+key); err != nil {
				return err
			}
		}
		for _, key := range slices.Sorted(maps.Keys(a.ModelOverrides)) {
			if !known[key] {
				return fmt.Errorf("%s.models_by_name: unknown agent %q (path: %s.models_by_name.%s)", path, key, path, key)
			}
			if err := check(a.ModelOverrides[key], path+".models_by_name."+key); err != nil {
				return err
			}
		}
		for _, host := range slices.Sorted(maps.Keys(a.ModelHostOverrides)) {
			if !knownHosts[host] {
				return fmt.Errorf("%s.models_by_host: unknown host %q (path: %s.models_by_host.%s)", path, host, path, host)
			}
			for _, key := range slices.Sorted(maps.Keys(a.ModelHostOverrides[host])) {
				if !known[key] && !validModelProfileKeys[key] {
					return fmt.Errorf("%s.models_by_host.%s: unknown key %q (path: %s.models_by_host.%s.%s)", path, host, key, path, host, key)
				}
				if err := check(a.ModelHostOverrides[host][key], path+".models_by_host."+host+"."+key); err != nil {
					return err
				}
			}
		}
		return nil
	}
	checkEffective := func(a *config.AgentsConfig, setName string) error {
		disabled := map[string]bool{}
		for _, n := range a.Disabled {
			disabled[n] = true
		}
		for _, host := range hosts {
			for _, def := range defs {
				if disabled[def.Name] {
					continue
				}
				model, modelPath := resolveModelField(a, def, host.key, false)
				effort, effortPath := resolveModelField(a, def, host.key, true)
				if setName != "" {
					set := cfg.Agents.ModelSets[setName]
					if declaredAt(set, modelPath, false) {
						modelPath = "agents.model_sets." + setName + strings.TrimPrefix(modelPath, "agents")
					}
					if declaredAt(set, effortPath, true) {
						effortPath = "agents.model_sets." + setName + strings.TrimPrefix(effortPath, "agents")
					}
				}
				if model != "" {
					if err := validateModelValueForHost(modelPath, model, host.key); err != nil {
						return fmt.Errorf("%w; configure agents.models_by_host for per-host models", err)
					}
				}
				if effort != "" && effort != "inherit" {
					var allowed []string
					switch host.key {
					case "claude":
						allowed = []string{"low", "medium", "high", "xhigh", "max"}
					case "pi":
						allowed = []string{"off", "minimal", "low", "medium", "high", "xhigh", "max"}
					}
					if allowed != nil && !slices.Contains(allowed, effort) {
						return fmt.Errorf("%s: reasoning_effort %q invalid for host %q (supported: %s)", effortPath, effort, host.key, strings.Join(allowed, ", "))
					}
				}
			}
		}
		return nil
	}
	if err := checkDeclarations(&cfg.Agents, "agents"); err != nil {
		return err
	}
	if err := checkEffective(&cfg.Agents, ""); err != nil {
		return err
	}
	for _, name := range config.ModelSetNames(&cfg.Agents) {
		set := cfg.Agents.ModelSets[name]
		fragment := config.AgentsConfig{Models: set.Models, ModelOverrides: set.ModelOverrides, ModelHostOverrides: set.ModelHostOverrides}
		if err := checkDeclarations(&fragment, "agents.model_sets."+name); err != nil {
			return fmt.Errorf("model set %q: %w", name, err)
		}
		merged := config.ApplyModelSet(&cfg.Agents, set)
		if err := checkEffective(&merged, name); err != nil {
			return err
		}
	}
	return nil
}
func declaredAt(set config.ModelSetConfig, path string, effort bool) bool {
	parts := strings.Split(strings.TrimPrefix(path, "agents."), ".")
	var value config.ModelValue
	if len(parts) < 3 {
		return false
	}
	switch parts[0] {
	case "models":
		value = set.Models[parts[1]]
	case "models_by_name":
		value = set.ModelOverrides[parts[1]]
	case "models_by_host":
		if len(parts) < 4 {
			return false
		}
		value = set.ModelHostOverrides[parts[1]][parts[2]]
	}
	if effort {
		return value.ReasoningEffort != ""
	}
	return value.Model != ""
}

// Read only native header fields. TOML instructions are decoded as text and
// can never be mistaken for top-level settings.
func readLocalModelFields(path, host string) (config.ModelValue, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return config.ModelValue{}, nil
	}
	if err != nil {
		return config.ModelValue{}, fmt.Errorf("reading %s: %w", path, err)
	}
	fields, err := parseNativeModelFields(data, host)
	if err != nil {
		return config.ModelValue{}, fmt.Errorf("reading %s: %w", path, err)
	}
	return config.ModelValue{Model: fields.Model.Value, ReasoningEffort: fields.Effort.Value}, nil
}

func parseNativeModelFields(data []byte, host string) (nativeModelFields, error) {
	var fields nativeModelFields
	var header map[string]any
	var err error
	if host == "codex" {
		err = toml.Unmarshal(data, &header)
	} else {
		normalized := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
		if !bytes.HasPrefix(normalized, []byte("---\n")) {
			return fields, nil
		}
		rest := normalized[4:]
		end := bytes.Index(rest, []byte("\n---\n"))
		if end < 0 {
			return fields, fmt.Errorf("unclosed YAML frontmatter")
		}
		err = yaml.Unmarshal(rest[:end], &header)
	}
	if err != nil {
		return fields, err
	}
	field := map[string]string{"codex": "model_reasoning_effort", "claude": "effort", "opencode": "reasoningEffort", "pi": "thinking"}[host]
	for _, entry := range []struct {
		key    string
		target *modelFieldValue
	}{{"model", &fields.Model}, {field, &fields.Effort}} {
		if value, present := header[entry.key]; present {
			str, ok := value.(string)
			if !ok {
				return fields, fmt.Errorf("%s must be a string", entry.key)
			}
			*entry.target = modelFieldValue{Present: true, Value: str}
		}
	}
	return fields, nil
}

type preparedOutput struct {
	path, host string
	content    []byte
	hints      []string
}
type preparedSubagents struct {
	hosts                []hostTarget
	allDefs, definitions []*subagent.Definition
	outputs              []preparedOutput
	stateBytes           []byte
}

func prepareSubagents(root string, cfg *config.Config, target string) (preparedSubagents, error) {
	var result preparedSubagents
	eff, _, err := config.EffectiveAgents(cfg, root)
	if err != nil {
		return result, err
	}
	hosts, err := resolveHostTargets(cfg)
	if err != nil {
		return result, err
	}
	result.hosts = hosts
	defs, err := discoverSubagentSources(root)
	if err != nil {
		return result, err
	}
	result.allDefs = defs
	validationConfig := *cfg
	if target != "" {
		validationConfig.Agents.Disabled = slices.DeleteFunc(slices.Clone(cfg.Agents.Disabled), func(name string) bool { return name == target })
	}
	if err := validateModelConfig(&validationConfig, defs, hosts); err != nil {
		return result, err
	}
	effective := *cfg
	effective.Agents = eff
	disabled := map[string]bool{}
	for _, n := range eff.Disabled {
		disabled[n] = true
	}
	for _, def := range defs {
		if target != "" {
			if def.Name == target {
				result.definitions = append(result.definitions, def)
			}
		} else if !disabled[def.Name] {
			result.definitions = append(result.definitions, def)
		}
	}
	if target != "" && len(result.definitions) == 0 {
		return result, fmt.Errorf("subagent %q not found in built-in or project sources", target)
	}
	state, err := readAgentModelState(root)
	if err != nil {
		return result, err
	}
	for _, def := range result.definitions {
		for _, h := range hosts {
			opts, err := resolveCompileOptions(&effective, def, h.key)
			if err != nil {
				return result, err
			}
			path := filepath.Join(root, h.relDir, def.Name+h.ext)
			actual, exists, err := readNativeModelFields(path, h.key)
			if err != nil {
				return result, err
			}
			relPath := filepath.ToSlash(filepath.Join(h.relDir, def.Name+h.ext))
			previous := state.Paths[relPath]
			if !exists {
				previous = nil
			}
			record := captureModelState(previous, actual)
			pin := config.ModelValue{}
			if record.Model.RetainedLocal != nil {
				pin.Model = record.Model.RetainedLocal.Value
			}
			if record.Effort.RetainedLocal != nil {
				pin.ReasoningEffort = record.Effort.RetainedLocal.Value
			}
			var hints []string
			defaultModel, defaultEffort := "", ""
			switch h.key {
			case "claude":
				defaultModel = def.ModelProfile.ClaudeModel()
			case "codex":
				defaultEffort = def.ModelProfile.CodexReasoningEffort()
			case "pi":
				defaultEffort = def.ModelProfile.PiThinking()
			}
			if opts.Model == "" && pin.Model != "" {
				opts.FallbackModel = pin.Model
				if pin.Model != defaultModel {
					hints = append(hints, fmt.Sprintf("  info: preserved local model %q for %s (set agents.models_by_name/models_by_host in .flowforge/config.yaml to pin explicitly)", pin.Model, filepath.Join(h.relDir, def.Name+h.ext)))
				}
			}
			if !opts.EffortConfigured && pin.ReasoningEffort != "" {
				opts.FallbackEffort = pin.ReasoningEffort
				if pin.ReasoningEffort != defaultEffort {
					hints = append(hints, fmt.Sprintf("  info: preserved local reasoning_effort %q for %s", pin.ReasoningEffort, filepath.Join(h.relDir, def.Name+h.ext)))
				}
			}
			content, err := h.compile(def, opts)
			if err != nil {
				return result, fmt.Errorf("compiling %s for %s: %w", def.Name, h.key, err)
			}
			generated, err := parseNativeModelFields(content, h.key)
			if err != nil {
				return result, fmt.Errorf("reading compiled %s: %w", path, err)
			}
			record.Model.LastGenerated = &generated.Model
			record.Effort.LastGenerated = &generated.Effort
			state.Paths[relPath] = record
			result.outputs = append(result.outputs, preparedOutput{path: path, host: h.key, content: content, hints: hints})
		}
	}
	if err := pruneCleanedModelState(root, hosts, defs, state); err != nil {
		return result, err
	}
	result.stateBytes, err = serializeAgentModelState(state)
	if err != nil {
		return result, fmt.Errorf("serializing %s: %w", agentModelStatePath(root), err)
	}
	return result, nil
}
