package subagent

import (
	"fmt"
	"strings"
)

// CompileCodex generates a Codex native agent definition file in TOML format.
func CompileCodex(def *Definition) ([]byte, error) {
	sandboxMode := "workspace-write"
	if def.Permission == "read-only" {
		sandboxMode = "read-only"
	}

	// Build TOML frontmatter manually (no TOML library dependency for now)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("name = %q\n", def.Name))
	sb.WriteString(fmt.Sprintf("description = %q\n", def.Description))
	sb.WriteString(fmt.Sprintf("sandbox_mode = %q\n", sandboxMode))
	sb.WriteString(fmt.Sprintf("model_reasoning_effort = %q\n", def.ModelProfile.CodexReasoningEffort()))
	sb.WriteString("developer_instructions = \"\"\"\n")

	// Replace "invoke the Skill tool" instruction with file read directive
	body := replaceSkillInvocationWithFileRead(def.Body, def.DefaultSkill)
	sb.WriteString(body)

	if !strings.HasSuffix(body, "\n") {
		sb.WriteString("\n")
	}
	sb.WriteString("\"\"\"\n")

	return []byte(sb.String()), nil
}

// replaceSkillInvocationWithFileRead replaces the "invoke the Skill tool" section
// in the Default Skill paragraph with a file read directive.
func replaceSkillInvocationWithFileRead(body, defaultSkill string) string {
	// Look for the Default Skill section and replace the invocation instruction
	lines := strings.Split(body, "\n")
	var result []string
	inDefaultSkillSection := false

	for _, line := range lines {
		if strings.HasPrefix(line, "## Default Skill") {
			inDefaultSkillSection = true
			result = append(result, line)
			continue
		}
		if inDefaultSkillSection && strings.HasPrefix(line, "## ") {
			inDefaultSkillSection = false
		}

		if inDefaultSkillSection && strings.Contains(line, "invoke the Skill tool") {
			// Replace the entire instruction
			replacement := fmt.Sprintf("Read and follow `.agents/skills/%s/SKILL.md` completely before taking any other action.", defaultSkill)
			// Preserve leading whitespace
			leading := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			result = append(result, leading+replacement)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
