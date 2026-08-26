package command

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
)

func newAssetsCmd() *cobra.Command {
	assets := &cobra.Command{Use: "assets", Short: "Inspect FlowForge managed project assets"}
	var jsonOutput bool
	verify := &cobra.Command{
		Use:   "verify [project]",
		Short: "Compare managed Skills and agent rules without modifying the project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := "."
			if len(args) == 1 {
				project = args[0]
			}
			root, err := filepath.Abs(project)
			if err != nil {
				return fmt.Errorf("resolving project path: %w", err)
			}
			comparison, err := verifyManagedAssets(root)
			if err != nil {
				return err
			}
			if jsonOutput {
				payload := make([]map[string]string, 0, len(comparison.Entries))
				for _, entry := range comparison.Entries {
					payload = append(payload, map[string]string{"state": string(entry.State), "target": entry.TargetPath})
				}
				data, err := json.MarshalIndent(map[string]any{"current": comparison.IsCurrent(), "entries": payload}, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(data))
			} else {
				for _, entry := range comparison.Entries {
					cmd.Printf("%s %s\n", entry.State, entry.TargetPath)
				}
			}
			if !comparison.IsCurrent() {
				return errPolicyViolation
			}
			return nil
		},
	}
	verify.Flags().BoolVar(&jsonOutput, "json", false, "Output JSON")
	assets.AddCommand(verify)
	return assets
}

func verifyManagedAssets(projectRoot string) (managedAssetComparison, error) {
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return managedAssetComparison{}, fmt.Errorf("loading project configuration: %w", err)
	}
	assetsDir, cleanup, err := locateAssetsDir()
	if err != nil {
		return managedAssetComparison{}, err
	}
	defer cleanup()
	return compareManagedAssets(assetsDir, projectRoot, cfg.DocsRoot(projectRoot))
}
