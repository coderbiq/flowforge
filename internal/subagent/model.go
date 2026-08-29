package subagent

// Definition represents a parsed subagent definition from a .md file.
type Definition struct {
	Name         string
	Description  string
	ModelProfile ModelProfile
	DefaultSkill string
	DetourSkills []string
	Permission   string
	After        []string
	Before       []string
	ReturnsTo    []string
	Body         string
}
