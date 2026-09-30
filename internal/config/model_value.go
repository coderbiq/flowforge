package config

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ModelValue records independently declared model and reasoning fields. Empty
// fields are absent; explicit empty strings are rejected at the decoding seam.
type ModelValue struct {
	Model           string `yaml:"model,omitempty"`
	ReasoningEffort string `yaml:"reasoning_effort,omitempty"`
}

func parseModelValue(raw any) (ModelValue, error) {
	var out ModelValue
	switch value := raw.(type) {
	case string:
		if value == "" {
			return out, fmt.Errorf("model must be a non-empty string")
		}
		out.Model = value
	case map[string]any:
		if len(value) == 0 {
			return out, fmt.Errorf("model object must declare model or reasoning_effort")
		}
		keys := make([]string, 0, len(value))
		for k := range value {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if key != "model" && key != "reasoning_effort" {
				return out, fmt.Errorf("%s: unknown field", key)
			}
			s, ok := value[key].(string)
			if !ok || s == "" {
				return out, fmt.Errorf("%s: must be a non-empty string", key)
			}
			if key == "model" {
				out.Model = s
			} else {
				out.ReasoningEffort = s
			}
		}
	default:
		return out, fmt.Errorf("must be a model string or model/reasoning_effort object")
	}
	return out, nil
}
func (v *ModelValue) UnmarshalYAML(node *yaml.Node) error {
	var raw any
	if err := node.Decode(&raw); err != nil {
		return err
	}
	parsed, err := parseModelValue(raw)
	if err == nil {
		*v = parsed
	}
	return err
}
func (v ModelValue) MarshalYAML() (any, error) {
	if v.ReasoningEffort == "" {
		if v.Model == "" {
			return nil, fmt.Errorf("model must be non-empty")
		}
		return v.Model, nil
	}
	type object ModelValue
	return object(v), nil
}
func modelValueDecodeHook(from, to reflect.Type, raw any) (any, error) {
	if to != reflect.TypeOf(ModelValue{}) {
		return raw, nil
	}
	if from == to {
		return raw, nil
	}
	return parseModelValue(raw)
}

// Validate the original YAML before Viper discards null declarations or applies
// weak conversions. This also supplies precise declaration paths.
func validateRawModels(root string) error {
	data, err := os.ReadFile(ConfigPath(root))
	if err != nil {
		return err
	}
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}
	agents, ok := raw["agents"].(map[string]any)
	if !ok {
		return nil
	}
	return validateAgentModelDeclarations(agents)
}

func (a *AgentsConfig) UnmarshalYAML(node *yaml.Node) error {
	var raw map[string]any
	if err := node.Decode(&raw); err != nil {
		return err
	}
	if err := validateAgentModelDeclarations(raw); err != nil {
		return err
	}
	type plain AgentsConfig
	return node.Decode((*plain)(a))
}

func validateAgentModelDeclarations(agents map[string]any) error {
	var layer func(map[string]any, string) error
	values := func(raw any, path string) error {
		m, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: must be a mapping", path)
		}
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if _, err := parseModelValue(m[key]); err != nil {
				if _, isObject := m[key].(map[string]any); isObject {
					if field, message, ok := strings.Cut(err.Error(), ": "); ok {
						return fmt.Errorf("%s.%s.%s: %s", path, key, field, message)
					}
				}
				return fmt.Errorf("%s.%s: %w", path, key, err)
			}
		}
		return nil
	}
	layer = func(m map[string]any, path string) error {
		for _, key := range []string{"models", "models_by_name", "models_by_host"} {
			val, present := m[key]
			if !present {
				continue
			}
			if key != "models_by_host" {
				if err := values(val, path+"."+key); err != nil {
					return err
				}
				continue
			}
			hosts, ok := val.(map[string]any)
			if !ok {
				return fmt.Errorf("%s.%s: must be a mapping", path, key)
			}
			keys := make([]string, 0, len(hosts))
			for h := range hosts {
				keys = append(keys, h)
			}
			sort.Strings(keys)
			for _, h := range keys {
				if err := values(hosts[h], path+"."+key+"."+h); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := layer(agents, "agents"); err != nil {
		return err
	}
	if val, present := agents["model_sets"]; present {
		sets, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("agents.model_sets: must be a mapping")
		}
		keys := make([]string, 0, len(sets))
		for k := range sets {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			set, ok := sets[k].(map[string]any)
			if !ok {
				return fmt.Errorf("agents.model_sets.%s: must be a mapping", k)
			}
			if err := layer(set, "agents.model_sets."+k); err != nil {
				return err
			}
		}
	}
	return nil
}
