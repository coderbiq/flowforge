package subagent

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// CompileClaudeCode generates a Claude Code native agent definition file with
// zero options: the model field carries the profile default, matching
// pre-option behavior byte for byte.
func CompileClaudeCode(def *Definition) ([]byte, error) {
	return CompileClaudeCodeWithOptions(def, CompileOptions{})
}

// CompileClaudeCodeWithOptions generates a Claude Code native agent
// definition file with explicit compile options applied. The model field
// precedence: config-pinned Model, then preserved FallbackModel, then the
// profile default.
func CompileClaudeCodeWithOptions(def *Definition, opts CompileOptions) ([]byte, error) {
	type claudeFrontmatter struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Model       string   `yaml:"model"`
		Skills      []string `yaml:"skills"`
	}

	model := resolveModel(opts)
	if model == "" {
		model = def.ModelProfile.ClaudeModel()
	}

	fm := claudeFrontmatter{
		Name:        def.Name,
		Description: def.Description,
		Model:       model,
		Skills:      []string{def.DefaultSkill},
	}

	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return nil, fmt.Errorf("marshaling Claude Code frontmatter for %s: %w", def.Name, err)
	}

	fallbackNote := fmt.Sprintf("\n\n_Note: If the skill is not preloaded, explicitly invoke the Skill tool with `%s` before proceeding._\n", def.DefaultSkill)
	body := def.Body + fallbackNote

	output := "---\n" + string(fmBytes) + "---\n" + body
	return []byte(output), nil
}
