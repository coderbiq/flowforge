package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	initForce bool
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize FlowForge local tracker and deploy mattpocock skills",
		Long: `Initializes the local-first issue tracker environment in the target repository.
Deploys mattpocock engineering/productivity skills into .agents/skills/,
sets up docs/agents/ rules, and creates .scratch/ for issue tracking.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := "."
			if len(args) > 0 {
				targetDir = args[0]
			}

			absTarget, err := filepath.Abs(targetDir)
			if err != nil {
				return fmt.Errorf("resolving target path: %w", err)
			}

			// 1. Create .scratch directory
			scratchDir := filepath.Join(absTarget, ".scratch")
			if err := os.MkdirAll(scratchDir, 0755); err != nil {
				return fmt.Errorf("creating .scratch directory: %w", err)
			}

			// 2. Deploy managed assets (skills, agent docs, AGENTS.md)
			if err := deployManagedAssets(absTarget); err != nil {
				return fmt.Errorf("deploying assets: %w", err)
			}

			fmt.Println("✓ FlowForge local tracker initialized successfully.")
			fmt.Println("✓ Deployed mattpocock skills to .agents/skills/")
			fmt.Println("✓ Configured docs/agents/issue-tracker.md")
			fmt.Println("✓ Created .scratch/ directory")
			fmt.Println("\nReady to run /ask-matt, /grill-with-docs, or /to-spec in your agent.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force overwrite existing files")
	return cmd
}
