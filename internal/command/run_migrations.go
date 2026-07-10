package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
)

func newRunMigrationsCmd() *cobra.Command {
	var fromVersion string

	cmd := &cobra.Command{
		Use:    "_run-migrations",
		Short:  "Run pending data migrations after upgrade",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromVersion == "" {
				return fmt.Errorf("--from is required")
			}

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
				fmt.Println("No wiki root configured — skipping migrations.")
				return nil
			}

			return runPendingMigrations(fromVersion, "", wikiRoot)
		},
	}

	cmd.Flags().StringVar(&fromVersion, "from", "", "Version being upgraded from")
	return cmd
}
