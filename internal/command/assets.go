package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/orchestration"
	"flowforge/internal/version"
)

func newAssetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Manage FlowForge managed assets in the current project",
	}
	cmd.AddCommand(newAssetsUpdateCmd())
	cmd.AddCommand(newAssetsAdapterCmd())
	return cmd
}

func newAssetsAdapterCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "adapter", Short: "Manage optional host adapters"}
	cmd.AddCommand(newAssetsAdapterInstallCmd(), newAssetsAdapterUpdateCmd(), newAssetsAdapterUninstallCmd())
	return cmd
}
func newAssetsAdapterUpdateCmd() *cobra.Command {
	var adapter string
	cmd := &cobra.Command{Use: "update", Short: "Update an installed optional adapter", RunE: func(cmd *cobra.Command, args []string) error {
		if adapter != "opencode" {
			return fmt.Errorf("unsupported adapter %q", adapter)
		}
		svc, err := currentConfigService()
		if err != nil {
			return err
		}
		defer svc.Close()
		return applyOpenCodeAdapter(cmd, svc.ProjectRoot(), true)
	}}
	cmd.Flags().StringVar(&adapter, "adapter", "", "Adapter name (opencode)")
	return cmd
}

func newAssetsAdapterInstallCmd() *cobra.Command {
	var adapter string
	cmd := &cobra.Command{Use: "install", Short: "Install an optional adapter", RunE: func(cmd *cobra.Command, args []string) error {
		if adapter != "opencode" {
			return fmt.Errorf("unsupported adapter %q", adapter)
		}
		svc, err := currentConfigService()
		if err != nil {
			return err
		}
		defer svc.Close()
		return applyOpenCodeAdapter(cmd, svc.ProjectRoot(), false)
	}}
	cmd.Flags().StringVar(&adapter, "adapter", "", "Adapter name (opencode)")
	return cmd
}
func newAssetsAdapterUninstallCmd() *cobra.Command {
	var adapter string
	cmd := &cobra.Command{Use: "uninstall", Short: "Remove an optional adapter", RunE: func(cmd *cobra.Command, args []string) error {
		if adapter != "opencode" {
			return fmt.Errorf("unsupported adapter %q", adapter)
		}
		svc, err := currentConfigService()
		if err != nil {
			return err
		}
		defer svc.Close()
		manifest, err := core.LoadProjectManifest(svc.ProjectRoot())
		if err != nil {
			return err
		}
		kept := manifest.Files[:0]
		for _, entry := range manifest.Files {
			if entry.Type != "opencode_agent" {
				kept = append(kept, entry)
				continue
			}
			path := filepath.Join(svc.ProjectRoot(), entry.Target)
			data, readErr := os.ReadFile(path)
			if readErr == nil && core.SHA256Hex(data) != entry.SHA256 {
				fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s (preserved)\n", entry.Target)
				kept = append(kept, entry)
				continue
			}
			if readErr != nil && !os.IsNotExist(readErr) {
				return fmt.Errorf("reading adapter file: %w", readErr)
			}
			if readErr == nil {
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("removing adapter file: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", entry.Target)
			}
		}
		manifest.Files = kept
		return manifest.Save(svc.ProjectRoot())
	}}
	cmd.Flags().StringVar(&adapter, "adapter", "", "Adapter name (opencode)")
	return cmd
}
func applyOpenCodeAdapter(cmd *cobra.Command, root string, update bool) error {
	policy := orchestration.DefaultPolicy()
	files, err := orchestration.RenderOpenCode(policy)
	if err != nil {
		return err
	}
	manifest, err := core.LoadProjectManifest(root)
	if err != nil {
		manifest = &core.ProjectManifest{}
	}
	byTarget := map[string]core.FileEntry{}
	for _, e := range manifest.Files {
		byTarget[e.Target] = e
	}
	for name, content := range files {
		target := filepath.Join(".opencode", "agents", name)
		path := filepath.Join(root, target)
		old, managed := byTarget[target]
		existing, readErr := os.ReadFile(path)
		if readErr != nil && !os.IsNotExist(readErr) {
			return fmt.Errorf("reading adapter target: %w", readErr)
		}
		if !managed && !os.IsNotExist(readErr) {
			fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s (not managed)\n", target)
			continue
		}
		if managed && !os.IsNotExist(readErr) && core.SHA256Hex(existing) != old.SHA256 {
			fmt.Fprintf(cmd.ErrOrStderr(), "! conflict: %s (preserved)\n", target)
			continue
		}
		if managed && !os.IsNotExist(readErr) && core.SHA256Hex(existing) == core.SHA256Hex(content) {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
		entry := core.FileEntry{Source: "generated/opencode/agents/" + name, Target: target, SHA256: core.SHA256Hex(content), Type: "opencode_agent"}
		replaced := false
		for i, e := range manifest.Files {
			if e.Target == target {
				manifest.Files[i] = entry
				replaced = true
			}
		}
		if !replaced {
			manifest.Files = append(manifest.Files, entry)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", map[bool]string{true: "~", false: "+"}[managed], target)
	}
	return manifest.Save(root)
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
