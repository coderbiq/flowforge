package command

import (
	"fmt"
	"os"
	"os/exec"

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
then atomically replaces the current installation.

If a newer version is available, the binary is downloaded,
verified with Ed25519 signature and SHA256 checksum, and
installed atomatically. On failure, the previous version
is automatically restored.

After the CLI binary is upgraded, project facilities are synchronized
(equivalent to running flowforge sync).`,
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

			execPath, execErr := os.Executable()
			if execErr != nil {
				execPath = ""
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
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Upgraded from %s to %s\n",
				result.OldVersion, result.NewVersion)

			projectRoot, pErr := config.FindProjectRoot(".")
			if pErr != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Skipping project synchronization: %v\n", pErr)
				return nil
			}

			if execPath == "" {
				return fmt.Errorf("synchronizing project: cannot locate upgraded executable")
			}

			// The child command is the project lifecycle boundary. Disable the
			// global version checker here so an upgrade cannot start a second
			// asynchronous check while reconciling manifest assets.
			assetCmd := exec.Command(execPath, "--no-version-check", "sync")
			assetCmd.Dir = projectRoot
			assetCmd.Stdout = cmd.OutOrStdout()
			assetCmd.Stderr = cmd.ErrOrStderr()
			assetCmd.Stdin = cmd.InOrStdin()
			if aErr := assetCmd.Run(); aErr != nil {
				return fmt.Errorf("synchronizing project: %w", aErr)
			}

			migrateCmd := exec.Command(execPath, "--no-version-check", "_run-migrations", "--from", result.OldVersion)
			migrateCmd.Dir = projectRoot
			migrateCmd.Stdout = cmd.OutOrStdout()
			migrateCmd.Stderr = cmd.ErrOrStderr()
			if mErr := migrateCmd.Run(); mErr != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Migration: %v\n", mErr)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&targetVersion, "version", "",
		"upgrade to a specific version")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"show available upgrade without installing")

	return cmd
}
