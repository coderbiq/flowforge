package command

import (
	"fmt"

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

After the CLI binary is upgraded, managed project assets are
also updated (equivalent to running flowforge assets update).`,
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
				fmt.Fprintf(cmd.ErrOrStderr(), "Skipping project assets update: %v\n", pErr)
				return nil
			}

			report, aErr := applyAssetUpdates(projectRoot)
			if aErr != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Project assets update: %v\n", aErr)
				return nil
			}

			if report == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "Project assets are up to date.")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Project assets updated: %s\n", report.Summary())
			if report.BlockUpdated {
				fmt.Fprintln(cmd.OutOrStdout(), "  AGENTS.md: block updated")
			}
			for _, f := range report.Added {
				fmt.Fprintf(cmd.OutOrStdout(), "  + %s\n", f.Target)
			}
			for _, f := range report.Updated {
				fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s\n", f.Target)
			}
			for _, f := range report.Conflict {
				fmt.Fprintf(cmd.ErrOrStderr(), "  ! conflict: %s (manual merge needed)\n", f.Target)
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
