package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
)

func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run data migrations manually",
		Long: `Run all registered data migrations unconditionally.
Each migration checks its own prerequisites and is safe to run multiple times.

Use this when upgrading was not triggered by 'flowforge upgrade'
(e.g. reinstalled the binary manually, or upgrade had bugs).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return fmt.Errorf("finding project root: %w", err)
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			wikiRoot := cfg.WikiRoot(projectRoot)
			if wikiRoot == "" {
				return fmt.Errorf("no wiki root configured")
			}

			return forceRunMigrations(wikiRoot)
		},
	}
	return cmd
}
