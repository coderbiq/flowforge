package subagent

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// CompileOpenCode generates an OpenCode native agent definition file.
func CompileOpenCode(def *Definition) ([]byte, error) {
	type opencodeFrontmatter struct {
		Description string `yaml:"description"`
		Mode        string `yaml:"mode"`
	}

	fm := opencodeFrontmatter{
		Description: def.Description,
		Mode:        "subagent",
	}

	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return nil, fmt.Errorf("marshaling OpenCode frontmatter for %s: %w", def.Name, err)
	}

	// OpenCode omits model field to inherit from parent session
	// Body retains "invoke the Skill tool" instruction unchanged
	output := "---\n" + string(fmBytes) + "---\n" + def.Body
	return []byte(output), nil
}
