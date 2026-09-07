package command

import (
	"errors"
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

			// Re-exec the new binary to sync assets with the correct embedded version.
			// The current process still holds the old binary's embed.FS in memory,
			// so calling deployManagedAssets here would deploy stale assets.
			syncErr := syncAssetsViaReExec(cmd)
			if syncErr != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Warning: post-upgrade asset sync failed: %v\n", syncErr)
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
	comparison, err := verifyManagedAssets(projectRoot)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to verify synchronized project assets: %v\n", err)
		return
	}
	if !comparison.IsCurrent() {
		fmt.Fprintf(cmd.ErrOrStderr(), "Warning: project assets remain missing or drifted after synchronization: %s\n", comparison.DivergenceMessage())
		return
	}

	// Synchronize subagents (non-fatal warning on failure)
	if _, err := deploySubagents(projectRoot, cfg, ""); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to synchronize subagents: %v\n", err)
		return
	}

	fmt.Fprintln(cmd.OutOrStdout(), successMessage)
}

// syncAssetsViaReExec runs `flowforge init` using the newly replaced binary
// so that the correct embedded assets (from the new version) are deployed.
// The current process still holds the old binary's embed.FS in memory.
func syncAssetsViaReExec(cmd *cobra.Command) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locating executable: %w", err)
	}

	projectRoot, err := config.FindProjectRoot(".")
	if err != nil {
		if !isProjectDirectory(".") {
			return nil
		}
		projectRoot = "."
	}

	initCmd := exec.Command(exePath, "init", projectRoot, "--force")
	initCmd.Stdout = cmd.OutOrStdout()
	initCmd.Stderr = cmd.ErrOrStderr()
	initCmd.Dir = projectRoot
	if err := initCmd.Run(); err != nil {
		return fmt.Errorf("running %s init: %w", exePath, err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "✓ Project skills and agent rules updated to latest version.")
	return nil
}
