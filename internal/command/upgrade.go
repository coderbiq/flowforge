package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/update"
	"flowforge/internal/version"
)

func newUpgradeCmd() *cobra.Command {
	var targetVersion string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade FlowForge CLI to the latest version",
		Long: `Upgrade downloads and verifies the latest FlowForge binary,
then atomically replaces the current installation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRun {
				manifest, err := update.DryRunUpgrade(version.Version)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Current version: %s\n", version.Version)
				fmt.Fprintf(cmd.OutOrStdout(), "Latest  version: %s\n", manifest.Version)
				if update.CompareVersions(manifest.Version, version.Version) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "Upgrade available.\n")
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Already up to date.\n")
				}
				return nil
			}

			_, execErr := os.Executable()
			if execErr != nil {
				return fmt.Errorf("locating current executable: %w", execErr)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Current version: %s\n", version.Version)

			var result *update.UpgradeResult
			var err error

			if targetVersion != "" {
				murl := fmt.Sprintf("https://github.com/coderbiq/flowforge/releases/download/%s/manifest.json",
					targetVersion)
				manifest, mErr := update.FetchManifest(murl)
				if mErr != nil {
					return mErr
				}
				result, err = update.UpgradeToVersion(manifest, version.Version, targetVersion)
			} else {
				result, err = update.Upgrade(version.Version)
			}

			if err != nil {
				if errors.Is(err, update.ErrAlreadyUpToDate) {
					fmt.Fprintf(cmd.OutOrStdout(), "Already up to date (%s).\n", version.Version)
					syncProjectAssets(cmd, "✓ Project skills and agent rules synchronized.")
					return nil
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Upgraded from %s to %s\n",
				result.OldVersion, result.NewVersion)

			syncProjectAssets(cmd, "✓ Project skills and agent rules updated to latest version.")

			return nil
		},
	}

	cmd.Flags().StringVar(&targetVersion, "version", "",
		"upgrade to a specific version")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"show available upgrade without installing")

	return cmd
}

func syncProjectAssets(cmd *cobra.Command, successMessage string) {
	projectRoot, err := config.FindProjectRoot(".")
	if err != nil {
		if !isProjectDirectory(".") {
			return
		}
		projectRoot = "."
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to load project configuration: %v\n", err)
		return
	}
	if err := deployManagedAssets(projectRoot, cfg.DocsRoot(projectRoot)); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to synchronize project assets: %v\n", err)
		return
	}
	fmt.Fprintln(cmd.OutOrStdout(), successMessage)
}
