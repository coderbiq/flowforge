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

	// Resolve enabled hosts; deselected hosts are out of scope entirely
	hosts, err := resolveHostTargets(cfg)
	if err != nil {
		return subagentStatusResult{}, err
	}

	// Build expected content maps per enabled host
	type hostExpectation struct {
		dir      string
		expected map[string][]byte
	}
	expectations := make([]hostExpectation, 0, len(hosts))
	for _, h := range hosts {
		expectations = append(expectations, hostExpectation{
			dir:      filepath.Join(projectRoot, h.relDir),
			expected: make(map[string][]byte),
		})
	}

	for _, def := range active {
		opts, err := resolveCompileOptions(cfg, def)
		if err != nil {
			return subagentStatusResult{}, err
		}
		for i, h := range hosts {
			content, err := h.compile(def, opts)
			if err != nil {
				return subagentStatusResult{}, fmt.Errorf("compiling %s for %s: %w", def.Name, h.key, err)
			}
			expectations[i].expected[filepath.Join(expectations[i].dir, def.Name+h.ext)] = content
		}
	}

	result := subagentStatusResult{Current: true}

	// Compare each enabled host directory
	for _, tc := range expectations {
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
