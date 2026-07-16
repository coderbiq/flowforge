package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"flowforge/internal/uninstall"
)

func newUninstallCmd() *cobra.Command {
	var (
		yes        bool
		keepConfig bool
		project    string
	)

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall FlowForge CLI",
		Long: `Uninstall removes the FlowForge CLI binary, configuration,
and optionally project artifacts.

Without --yes, shows a preview of files to be removed and asks for confirmation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintln(cmd.OutOrStdout(), "Will remove:")
				fmt.Fprintln(cmd.OutOrStdout(), "  - FlowForge CLI binary")

				if !keepConfig {
					home, err := os.UserHomeDir()
					if err == nil {
						fmt.Fprintf(cmd.OutOrStdout(), "  - %s/.flowforge/\n", home)
					}
				}

				if project != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s/.flowforge/\n", project)
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s/.agents/skills/\n", project)
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s/AGENTS.md (FLOWFORGE block)\n", project)
				}

				fmt.Fprint(cmd.OutOrStdout(), "\nContinue? [y/N] ")
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}

			binaryResult, err := uninstall.CleanBinary()
			if err != nil {
				return err
			}

			totalRemoved := 0
			hasErrors := false

			for _, path := range binaryResult.Removed {
				fmt.Fprintf(cmd.OutOrStdout(), "✓ removed %s\n", path)
				totalRemoved++
			}
			if binaryResult.HasErrors() {
				fmt.Fprint(cmd.ErrOrStderr(), binaryResult.ErrorSummary())
				hasErrors = true
			}

			if !keepConfig {
				home, err := os.UserHomeDir()
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: finding home dir: %v\n", err)
				} else {
					configResult, err := uninstall.CleanConfig(home)
					if err != nil {
						return err
					}
					for _, path := range configResult.Removed {
						fmt.Fprintf(cmd.OutOrStdout(), "✓ removed %s\n", path)
						totalRemoved++
					}
					if configResult.HasErrors() {
						fmt.Fprint(cmd.ErrOrStderr(), configResult.ErrorSummary())
						hasErrors = true
					}
				}
			}

			if project != "" {
				projectResult, err := uninstall.CleanProject(project)
				if err != nil {
					return err
				}
				for _, path := range projectResult.Removed {
					fmt.Fprintf(cmd.OutOrStdout(), "✓ removed %s\n", path)
					totalRemoved++
				}
				if projectResult.HasErrors() {
					fmt.Fprint(cmd.ErrOrStderr(), projectResult.ErrorSummary())
					hasErrors = true
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nUninstall complete: %d items removed.\n", totalRemoved)

			if hasErrors {
				return fmt.Errorf("uninstall completed with warnings")
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&keepConfig, "keep-config", false, "Keep configuration, remove binary only")
	cmd.Flags().StringVar(&project, "project", "", "Also remove project managed files")

	return cmd
}
