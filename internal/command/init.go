package command

import (
	"fmt"
	"os"
	"path/filepath"

	"flowforge/internal/config"

	"github.com/spf13/cobra"
)

var (
	initForce bool
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [path]",
		Aliases: []string{"sync"},
		Short:   "Initialize or sync FlowForge local tracker, Wiki structure, and skills",
		Long: `Initializes or synchronizes the local-first issue tracker environment in the target repository.
Deploys or updates flowforge engineering/productivity skills into .agents/skills/,
sets up <docs_dir>/agents/ rules, creates .flowforge/ configuration, and sets up the Wiki hierarchy.`,
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

			// 1. Create .flowforge/ configuration directory and default config if not present
			configDir := filepath.Join(absTarget, config.ConfigDirName)
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("creating .flowforge directory: %w", err)
			}

			configFile := filepath.Join(configDir, config.ConfigFileName)
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				defaultYAML := "version: 5.0.0\nversion_check: true\ndocs_dir: docs\n"
				if err := os.WriteFile(configFile, []byte(defaultYAML), 0644); err != nil {
					return fmt.Errorf("writing config file: %w", err)
				}
			} else if err != nil {
				return fmt.Errorf("checking config file: %w", err)
			}

			// 2. Resolve docs directory and create proposals/ and adr/ directories
			cfg, err := config.Load(absTarget)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			docsRoot := cfg.DocsRoot(absTarget)

			proposalsDir := filepath.Join(docsRoot, "proposals")
			if err := os.MkdirAll(proposalsDir, 0755); err != nil {
				return fmt.Errorf("creating proposals directory: %w", err)
			}

			adrDir := filepath.Join(docsRoot, "adr")
			if err := os.MkdirAll(adrDir, 0755); err != nil {
				return fmt.Errorf("creating adr directory: %w", err)
			}

			// 3. Create initial CONTEXT.md template if missing
			contextFile := filepath.Join(docsRoot, "CONTEXT.md")
			if _, err := os.Stat(contextFile); os.IsNotExist(err) {
				initialContext := "# System Context & Domain Glossary\n\nThis document maintains the core domain glossary, bounded contexts, and system architecture invariants.\n"
				if err := os.WriteFile(contextFile, []byte(initialContext), 0644); err != nil {
					return fmt.Errorf("writing context file: %w", err)
				}
			} else if err != nil {
				return fmt.Errorf("checking context file: %w", err)
			}

			// 4. Deploy managed assets (skills, agent docs, AGENTS.md)
			if err := deployManagedAssets(absTarget, docsRoot); err != nil {
				return fmt.Errorf("deploying assets: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "✓ FlowForge local tracker initialized successfully.")
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Created %s/ and %s hierarchy\n", config.ConfigDirName, proposalsDir)
			fmt.Fprintln(cmd.OutOrStdout(), "✓ Deployed flowforge skills to .agents/skills/")
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Configured %s/agents/ rules\n", docsRoot)
			fmt.Fprintln(cmd.OutOrStdout(), "\nReady to run /flowforge-route, /flowforge-align, or /flowforge-to-spec in your agent.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force overwrite existing files")
	return cmd
}
