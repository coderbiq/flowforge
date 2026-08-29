package command

import (
	"fmt"
	"path/filepath"

	"flowforge/internal/config"
	"flowforge/internal/subagent"
)

// subagentStatusEntry represents the deployment state of one subagent across hosts.
type subagentStatusEntry struct {
	State  string `json:"state"`
	Target string `json:"target"`
}

// subagentStatusResult aggregates all subagent status entries with an overall current flag.
type subagentStatusResult struct {
	Current bool                  `json:"current"`
	Entries []subagentStatusEntry `json:"entries"`
}

// computeSubagentStatus compares compiled-expected content against deployed files on disk.
func computeSubagentStatus(projectRoot string, cfg *config.Config) (subagentStatusResult, error) {
	// Discover all non-disabled subagent definitions
	definitions, err := discoverSubagentSources(projectRoot)
	if err != nil {
		return subagentStatusResult{}, err
	}

	// Filter out disabled
	disabled := make(map[string]bool)
	for _, name := range cfg.Agents.Disabled {
		disabled[name] = true
	}
	var active []*subagent.Definition
	for _, def := range definitions {
		if !disabled[def.Name] {
			active = append(active, def)
		}
	}

	// Build expected content maps for each host directory
	claudeDir := filepath.Join(projectRoot, ".claude", "agents")
	opencodeDir := filepath.Join(projectRoot, ".opencode", "agent")
	codexDir := filepath.Join(projectRoot, ".codex", "agents")

	claudeExpected := make(map[string][]byte)
	opencodeExpected := make(map[string][]byte)
	codexExpected := make(map[string][]byte)

	for _, def := range active {
		// Claude Code
		cc, err := subagent.CompileClaudeCode(def)
		if err != nil {
			return subagentStatusResult{}, fmt.Errorf("compiling %s for Claude Code: %w", def.Name, err)
		}
		claudeExpected[filepath.Join(claudeDir, def.Name+".md")] = cc

		// OpenCode
		oc, err := subagent.CompileOpenCode(def)
		if err != nil {
			return subagentStatusResult{}, fmt.Errorf("compiling %s for OpenCode: %w", def.Name, err)
		}
		opencodeExpected[filepath.Join(opencodeDir, def.Name+".md")] = oc

		// Codex
		cx, err := subagent.CompileCodex(def)
		if err != nil {
			return subagentStatusResult{}, fmt.Errorf("compiling %s for Codex: %w", def.Name, err)
		}
		codexExpected[filepath.Join(codexDir, def.Name+".toml")] = cx
	}

	result := subagentStatusResult{Current: true}

	// Compare each host directory
	for _, tc := range []struct {
		dir      string
		expected map[string][]byte
	}{
		{claudeDir, claudeExpected},
		{opencodeDir, opencodeExpected},
		{codexDir, codexExpected},
	} {
		entries, err := compareExpectedContent(tc.expected, tc.dir)
		if err != nil {
			return subagentStatusResult{}, fmt.Errorf("comparing %s: %w", tc.dir, err)
		}
		for _, e := range entries {
			result.Entries = append(result.Entries, subagentStatusEntry{
				State:  string(e.State),
				Target: e.TargetPath,
			})
			if e.State == managedAssetMissing || e.State == managedAssetDrifted {
				result.Current = false
			}
		}
	}

	return result, nil
}
