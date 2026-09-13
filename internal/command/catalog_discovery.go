package command

import (
	"fmt"
	"os"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

// discoverProposalCatalog discovers proposal artifacts for a scan directory,
// applying project-level evidence exemptions when a FlowForge project config
// is discoverable. The project root is resolved from the scanned directory
// first (so `check --dir <other-project>` applies that project's exemptions)
// and falls back to the working directory. A missing config keeps legacy
// zero-option semantics; a corrupt config warns on stderr and degrades to
// zero options rather than failing the scan.
func discoverProposalCatalog(dir string) (*tracker.Catalog, error) {
	opts := tracker.Options{}
	projectRoot, err := config.FindProjectRoot(dir)
	if err != nil {
		projectRoot, err = config.FindProjectRoot(".")
		if err != nil {
			projectRoot = ""
		}
	}
	if projectRoot != "" {
		if cfg, err := config.Load(projectRoot); err == nil {
			opts.ExemptProposals = cfg.Evidence.ExemptProposals
		} else if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: loading %s: %v (evidence exemptions ignored)\n", config.ConfigPath(projectRoot), err)
		}
	}
	return tracker.DiscoverArtifactsWithConfig(dir, opts)
}
