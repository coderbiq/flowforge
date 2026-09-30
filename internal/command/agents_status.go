package command

import (
	"fmt"
	"path/filepath"

	"flowforge/internal/config"
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
	prepared, err := prepareSubagents(projectRoot, cfg, "")
	if err != nil {
		return subagentStatusResult{}, err
	}
	type hostExpectation struct {
		dir      string
		expected map[string][]byte
	}
	expectations := make([]hostExpectation, 0, len(prepared.hosts))
	for _, h := range prepared.hosts {
		e := hostExpectation{dir: filepath.Join(projectRoot, h.relDir), expected: map[string][]byte{}}
		for _, o := range prepared.outputs {
			if o.host == h.key {
				e.expected[o.path] = o.content
			}
		}
		expectations = append(expectations, e)
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
