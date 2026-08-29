package subagent

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// CompileClaudeCode generates a Claude Code native agent definition file.
func CompileClaudeCode(def *Definition) ([]byte, error) {
	type claudeFrontmatter struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Model       string   `yaml:"model"`
		Skills      []string `yaml:"skills"`
	}

	fm := claudeFrontmatter{
		Name:        def.Name,
		Description: def.Description,
		Model:       def.ModelProfile.ClaudeModel(),
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
