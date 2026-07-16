package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage FlowForge configuration",
	}
	cmd.AddCommand(newConfigListCmd())
	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigSetCmd())
	return cmd
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := config.New(".")
			if err != nil {
				return err
			}
			defer svc.Close()

			values, err := svc.List()
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Project Config:")
			for k, v := range values {
				if len(k) > 7 && k[:8] == "project." {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s = %s\n", k, v)
				}
			}
			fmt.Fprintln(cmd.OutOrStdout(), "\nRuntime State:")
			for k, v := range values {
				if len(k) > 7 && k[:8] == "runtime." {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s = %s\n", k, v)
				}
			}
			return nil
		},
	}
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := config.New(".")
			if err != nil {
				return err
			}
			defer svc.Close()

			value, err := svc.Get(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), value)
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would set %s = %s\n", key, value)
				return nil
			}

			svc, err := config.New(".")
			if err != nil {
				return err
			}
			defer svc.Close()

			if err := svc.Set(key, value); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ %s = %s\n", key, value)
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without applying")
	return cmd
}