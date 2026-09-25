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
	// Resolve enabled hosts first (same order as deploySubagents, so both
	// paths surface identical config errors); deselected hosts are out of
	// scope entirely.
	hosts, err := resolveHostTargets(cfg)
	if err != nil {
		return subagentStatusResult{}, err
	}

	// Discover all non-disabled subagent definitions
	definitions, err := discoverSubagentSources(projectRoot)
	if err != nil {
		return subagentStatusResult{}, err
	}

	// Same config validation as deploySubagents: a bad model config must
	// fail status with the deploy error, before any expected-content work.
	if err := validateModelConfig(cfg, definitions, hosts); err != nil {
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
		for i, h := range hosts {
			opts, err := resolveCompileOptions(cfg, def, h.key)
			if err != nil {
				return subagentStatusResult{}, err
			}
			path := filepath.Join(expectations[i].dir, def.Name+h.ext)
			content, err := h.compile(def, opts)
			if err != nil {
				return subagentStatusResult{}, fmt.Errorf("compiling %s for %s: %w", def.Name, h.key, err)
			}
			// Preserve-merge parity with deploySubagents: the expected
			// content is what deploy would write, so when config does
			// not pin a model and the deployed file carries a local
			// `model:` the fresh compile would not reproduce, recompile
			// with it as the fallback. This keeps preserved-model files
			// current instead of falsely drifted; hosts whose files
			// carry no yaml frontmatter `model` (codex TOML) are
			// naturally unaffected.
			if opts.Model == "" {
				if existing := deployedModel(path); existing != "" && existing != frontmatterModel(content) {
					fallbackOpts := opts
					fallbackOpts.FallbackModel = existing
					content, err = h.compile(def, fallbackOpts)
					if err != nil {
						return subagentStatusResult{}, fmt.Errorf("compiling %s for %s: %w", def.Name, h.key, err)
					}
				}
			}
			expectations[i].expected[path] = content
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
