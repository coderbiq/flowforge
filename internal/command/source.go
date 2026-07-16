package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
)

func newSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Manage external knowledge sources",
	}
	cmd.AddCommand(newSourceListCmd())
	cmd.AddCommand(newSourceAddCmd())
	cmd.AddCommand(newSourceRemoveCmd())
	return cmd
}

func newSourceListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List registered external knowledge sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			if len(cfg.KnowledgeSources) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No external knowledge sources configured.")
				fmt.Fprintln(cmd.OutOrStdout(), "Use 'flowforge source add <name> --path <path>' to add one.")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-10s %-12s %-10s %s\n", "NAME", "TYPE", "CATEGORY", "TRUST", "PATH")
			fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 90))
			for _, src := range cfg.KnowledgeSources {
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-10s %-12s %-10s %s\n",
					src.Name, src.Type, src.Category, src.Trust, src.Path)
			}
			return nil
		},
	}
}

func newSourceAddCmd() *cobra.Command {
	var srcPath, srcType, category, trust, description string

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add an external knowledge source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if srcPath == "" {
				return fmt.Errorf("--path is required")
			}

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			for _, src := range cfg.KnowledgeSources {
				if src.Name == name {
					return fmt.Errorf("source %q already exists; use 'source remove' first or pick a different name", name)
				}
			}

			if srcType == "" {
				srcType = "file"
			}
			if trust == "" {
				trust = "medium"
			}

			cfg.KnowledgeSources = append(cfg.KnowledgeSources, config.KnowledgeSourceConfig{
				Name:        name,
				Path:        srcPath,
				Type:        srcType,
				Category:    category,
				Trust:       trust,
				Description: description,
			})

			if err := cfg.Save(projectRoot); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Added external knowledge source %q\n", name)
			fmt.Fprintf(cmd.OutOrStdout(), "  Path:  %s\n", srcPath)
			fmt.Fprintf(cmd.OutOrStdout(), "  Type:  %s\n", srcType)
			fmt.Fprintf(cmd.OutOrStdout(), "  Trust: %s\n", trust)
			return nil
		},
	}

	cmd.Flags().StringVar(&srcPath, "path", "", "Path to the knowledge source (directory or file)")
	cmd.Flags().StringVar(&srcType, "type", "file", "Access mechanism type (file, jira, confluence, url)")
	cmd.Flags().StringVar(&category, "category", "", "Content category (official_docs, team_knowledge, community, experimental, legacy)")
	cmd.Flags().StringVar(&trust, "trust", "medium", "Trust level (high, medium, low, unknown)")
	cmd.Flags().StringVar(&description, "description", "", "Human-readable description of the source")
	return cmd
}

func newSourceRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an external knowledge source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			projectRoot, err := config.FindProjectRoot(".")
			if err != nil {
				return err
			}

			cfg, err := config.Load(projectRoot)
			if err != nil {
				return err
			}

			found := false
			newSources := make([]config.KnowledgeSourceConfig, 0, len(cfg.KnowledgeSources))
			for _, src := range cfg.KnowledgeSources {
				if src.Name == name {
					found = true
					continue
				}
				newSources = append(newSources, src)
			}

			if !found {
				return fmt.Errorf("source %q not found", name)
			}

			cfg.KnowledgeSources = newSources
			if err := cfg.Save(projectRoot); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Removed external knowledge source %q\n", name)
			return nil
		},
	}
}
