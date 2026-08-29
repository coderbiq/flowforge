package subagent

// ModelProfile represents the capability tier of a subagent's model.
type ModelProfile string

const (
	ModelProfileHighCapability      ModelProfile = "high-capability"
	ModelProfileToolCapable         ModelProfile = "tool-capable"
	ModelProfileToolCapableReadOnly ModelProfile = "tool-capable-read-only"
)

// ClaudeModel maps the profile to Claude Code's model field.
func (m ModelProfile) ClaudeModel() string {
	switch m {
	case ModelProfileHighCapability:
		return "opus"
	case ModelProfileToolCapable, ModelProfileToolCapableReadOnly:
		return "sonnet"
	default:
		return "sonnet"
	}
}

// CodexReasoningEffort maps the profile to Codex's model_reasoning_effort field.
func (m ModelProfile) CodexReasoningEffort() string {
	switch m {
	case ModelProfileHighCapability:
		return "high"
	case ModelProfileToolCapable, ModelProfileToolCapableReadOnly:
		return "medium"
	default:
		return "medium"
	}
}
