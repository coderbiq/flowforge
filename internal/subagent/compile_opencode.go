package subagent

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// CompileOptions carries per-deployment compilation parameters resolved from
// project config and existing deployed files: an explicit model to pin
// (otherwise inherit), a fallback model preserved from a prior deployment
// (used only when Model is empty), edit-deny globs for host-level file
// protection, and an optional host-level execution budget (OpenCode `steps`,
// nil = inherit host behavior).
type CompileOptions struct {
	Model         string
	FallbackModel string
	EditDeny      []string
	MaxSteps      *int
}

// resolveModel applies the model precedence: config-pinned Model first, then
// a preserved FallbackModel, then empty (host default: the field is omitted).
func resolveModel(opts CompileOptions) string {
	if opts.Model != "" {
		return opts.Model
	}
	return opts.FallbackModel
}

// CompileOpenCode generates an OpenCode native agent definition file with
// zero options: the model field is omitted so the agent inherits the parent
// session model, matching pre-option behavior byte for byte.
func CompileOpenCode(def *Definition) ([]byte, error) {
	return CompileOpenCodeWithOptions(def, CompileOptions{})
}

// CompileOpenCodeWithOptions generates an OpenCode native agent definition
// file with explicit compile options applied.
func CompileOpenCodeWithOptions(def *Definition, opts CompileOptions) ([]byte, error) {
	type opencodeFrontmatter struct {
		Description string                       `yaml:"description"`
		Mode        string                       `yaml:"mode"`
		Model       string                       `yaml:"model,omitempty"`
		Steps       *int                         `yaml:"steps,omitempty"`
		Permission  map[string]map[string]string `yaml:"permission,omitempty"`
	}

	fm := opencodeFrontmatter{
		Description: def.Description,
		Mode:        "subagent",
		Model:       resolveModel(opts),
		Steps:       opts.MaxSteps,
	}
	if len(opts.EditDeny) > 0 {
		deny := make(map[string]string, len(opts.EditDeny))
		for _, glob := range opts.EditDeny {
			deny[glob] = "deny"
		}
		fm.Permission = map[string]map[string]string{"edit": deny}
	}

	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return nil, fmt.Errorf("marshaling OpenCode frontmatter for %s: %w", def.Name, err)
	}

	// Body retains "invoke the Skill tool" instruction unchanged
	output := "---\n" + string(fmBytes) + "---\n" + def.Body
	return []byte(output), nil
}
