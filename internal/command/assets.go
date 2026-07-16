package command

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/version"
)

func newAssetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Manage FlowForge managed assets in the current project",
	}
	cmd.AddCommand(newAssetsUpdateCmd())
	return cmd
}

func newAssetsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Apply managed asset updates to the current project",
		Long: `Reconcile managed assets (skills, templates, AGENTS.md) for the current project.

Compares the embedded asset manifest with the deployed manifest and applies
only the differences: adds new files, updates changed files, and reports
conflicts for files that were modified outside of FlowForge.

Managed assets are:
  .agents/skills/      SKILL definitions
  .flowforge/templates/  Card templates
  AGENTS.md            FlowForge directive block`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := currentConfigService()
			if err != nil {
				return fmt.Errorf("finding project root: %w (run flowforge init first)", err)
			}
			defer svc.Close()

			report, err := applyAssetUpdates(svc.ProjectRoot())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if report == nil {
				fmt.Fprintln(out, "Project assets are up to date.")
				return nil
			}

			fmt.Fprintf(out, "Project assets updated: %s\n", report.Summary())
			if report.BlockUpdated {
				fmt.Fprintln(out, "  AGENTS.md: block updated")
			}
			for _, f := range report.Added {
				fmt.Fprintf(out, "  + %s\n", f.Target)
			}
			for _, f := range report.Updated {
				fmt.Fprintf(out, "  ~ %s\n", f.Target)
			}
			for _, f := range report.Conflict {
				fmt.Fprintf(cmd.ErrOrStderr(), "  ! conflict: %s (manual merge needed)\n", f.Target)
			}
			return nil
		},
	}
	return cmd
}

func applyAssetUpdates(projectRoot string) (*AssetUpdateReport, error) {
	oldManifest, err := core.LoadProjectManifest(projectRoot)
	if err != nil {
		oldManifest = &core.ProjectManifest{}
	}

	newManifest, err := core.GenerateManifest(embeddedAssets, version.Version)
	if err != nil {
		return nil, fmt.Errorf("generating manifest: %w", err)
	}

	diff := core.CompareManifests(oldManifest, newManifest)
	if !diff.HasChanges() {
		return nil, nil
	}

	backupDir := ""
	if oldManifest.CLIVersion != "" {
		backupDir = filepath.Join(projectRoot, ".flowforge", "backup", oldManifest.CLIVersion)
	}
	report := core.ApplyUpgrade(diff, newManifest, projectRoot, embeddedAssets, backupDir)

	if report.Error != nil {
		return nil, report.Error
	}

	if err := newManifest.Save(projectRoot); err != nil {
		return nil, fmt.Errorf("saving updated manifest: %w", err)
	}

	return &AssetUpdateReport{
		SummaryLine:  diff.Summary(),
		BlockUpdated: report.BlockUpdated,
		Added:        report.Added,
		Updated:      report.Updated,
		Conflict:     report.Conflict,
	}, nil
}

type AssetUpdateReport struct {
	SummaryLine  string
	BlockUpdated bool
	Added        []core.FileEntry
	Updated      []core.FileEntry
	Conflict     []core.FileEntry
}

func (r *AssetUpdateReport) Summary() string { return r.SummaryLine }
