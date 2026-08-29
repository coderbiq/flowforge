package subagent

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse reads and parses a single subagent definition file.
func Parse(path string) (*Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading subagent file %s: %w", path, err)
	}

	frontmatter, body, ok := splitFrontmatter(data)
	if !ok {
		return nil, fmt.Errorf("subagent file %s: missing or malformed frontmatter delimiters", path)
	}

	var envelope struct {
		FlowforgeAgent struct {
			Name         string   `yaml:"name"`
			Description  string   `yaml:"description"`
			ModelProfile string   `yaml:"model_profile"`
			DefaultSkill string   `yaml:"default_skill"`
			DetourSkills []string `yaml:"detour_skills"`
			Permission   string   `yaml:"permission"`
			After        []string `yaml:"after"`
			Before       []string `yaml:"before"`
			ReturnsTo    []string `yaml:"returns_to"`
		} `yaml:"flowforge_agent"`
	}

	if err := yaml.Unmarshal(frontmatter, &envelope); err != nil {
		return nil, fmt.Errorf("parsing frontmatter in %s: %w", path, err)
	}

	filename := filepath.Base(path)
	expectedName := strings.TrimSuffix(filename, filepath.Ext(filename))
	if envelope.FlowforgeAgent.Name != expectedName {
		return nil, fmt.Errorf("subagent file %s: frontmatter name %q does not match filename %q", path, envelope.FlowforgeAgent.Name, expectedName)
	}

	def := &Definition{
		Name:         envelope.FlowforgeAgent.Name,
		Description:  envelope.FlowforgeAgent.Description,
		ModelProfile: ModelProfile(envelope.FlowforgeAgent.ModelProfile),
		DefaultSkill: envelope.FlowforgeAgent.DefaultSkill,
		DetourSkills: envelope.FlowforgeAgent.DetourSkills,
		Permission:   envelope.FlowforgeAgent.Permission,
		After:        envelope.FlowforgeAgent.After,
		Before:       envelope.FlowforgeAgent.Before,
		ReturnsTo:    envelope.FlowforgeAgent.ReturnsTo,
		Body:         string(body),
	}

	return def, nil
}

// ParseDir reads all .md files in a directory and returns parsed definitions sorted by name.
func ParseDir(dir string) ([]*Definition, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	var definitions []*Definition
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		def, err := Parse(path)
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, def)
	}

	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].Name < definitions[j].Name
	})

	return definitions, nil
}

// splitFrontmatter extracts YAML frontmatter delimited by --- from markdown content.
func splitFrontmatter(data []byte) (frontmatter, body []byte, ok bool) {
	normalized := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(normalized, []byte("---\n")) {
		return nil, data, false
	}
	rest := normalized[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	delimiterLength := len("\n---\n")
	if end < 0 && bytes.HasSuffix(rest, []byte("\n---")) {
		end = len(rest) - len("\n---")
		delimiterLength = len("\n---")
	}
	if end < 0 {
		return nil, data, false
	}
	after := end + delimiterLength
	return rest[:end], rest[after:], true
}
