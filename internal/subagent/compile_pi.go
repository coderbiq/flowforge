package subagent

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// CompilePi generates a PI (pi-subagents) native agent definition file with
// zero options; behavior is identical to CompilePiWithOptions with empty
// options.
func CompilePi(def *Definition) ([]byte, error) {
	return CompilePiWithOptions(def, CompileOptions{})
}

// CompilePiWithOptions generates a PI native agent definition file
// (.pi/agents/<name>.md): YAML frontmatter followed by the original body as
// the system prompt. The body is never rewritten — the Default Skill
// section's "(or read `.agents/skills/...` directly if no Skill tool is
// available)" fallback is the native path under PI. Options PI cannot
// express (steps budget, edit-deny globs, question deny) are ignored per
// the "hosts that cannot express an option ignore it" convention; the
// model is carried via resolveModel (config pin or preserve-merge
// fallback, empty = inherit the parent session model); the edit-deny
// protection is carried by the project-level PI extension instead
// (pi-host-integration design, section three).
func CompilePiWithOptions(def *Definition, opts CompileOptions) ([]byte, error) {
	// Field order is fixed (name, description, model, thinking, tools,
	// skills, inheritSkills) so compiled output stays stable and diffable.
	type piFrontmatter struct {
		Name          string   `yaml:"name"`
		Description   string   `yaml:"description"`
		Model         string   `yaml:"model,omitempty"`
		Thinking      string   `yaml:"thinking"`
		Tools         []string `yaml:"tools,omitempty"`
		Skills        []string `yaml:"skills"`
		InheritSkills bool     `yaml:"inheritSkills"`
	}

	fm := piFrontmatter{
		Name:          def.Name,
		Description:   def.Description,
		Model:         resolveModel(opts),
		Thinking:      def.ModelProfile.PiThinking(),
		Skills:        piSkills(def),
		InheritSkills: false,
	}
	// read-only permission compiles to a strict read-tool allowlist, the
	// PI equivalent of Codex's sandbox_mode = "read-only".
	if def.Permission == "read-only" {
		fm.Tools = []string{"read", "grep", "find", "ls"}
	}

	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return nil, fmt.Errorf("marshaling PI frontmatter for %s: %w", def.Name, err)
	}

	output := "---\n" + string(fmBytes) + "---\n" + def.Body
	return []byte(output), nil
}

// piSkills orders one role's skill binding: the default skill first, then
// detour skills. PI children get exactly this set (inheritSkills: false).
func piSkills(def *Definition) []string {
	skills := make([]string, 0, 1+len(def.DetourSkills))
	if def.DefaultSkill != "" {
		skills = append(skills, def.DefaultSkill)
	}
	skills = append(skills, def.DetourSkills...)
	return skills
}
