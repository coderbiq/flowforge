package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newSkillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage FlowForge skills",
	}

	cmd.AddCommand(newSkillUpdateCmd())

	return cmd
}

func newSkillUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update deployed skills to the latest version",
		Long: `Redeploy FlowForge skills from the current binary to an already-initialized project.

This is useful after upgrading the flowforge binary to pick up SKILL improvements.
Only skills and templates are updated; project config, wiki content, and runtime state are not modified.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := currentConfigService()
			if err != nil {
				return fmt.Errorf("finding project root: %w (run flowforge init first)", err)
			}
			defer svc.Close()
			projectRoot := svc.ProjectRoot()

			assetsDir, cleanup, err := locateAssetsDir()
			if err != nil {
				return err
			}
			defer cleanup()

			skillsSrc := filepath.Join(assetsDir, "skills")
			skillsDst := filepath.Join(projectRoot, ".agents", "skills")
			if _, err := os.Stat(skillsSrc); err != nil {
				return fmt.Errorf("skills directory not found in assets: %w", err)
			}

			if err := copyDir(skillsSrc, skillsDst, true); err != nil {
				return fmt.Errorf("deploying skills: %w", err)
			}

			templatesSrc := filepath.Join(assetsDir, "templates")
			templatesDst := filepath.Join(projectRoot, ".flowforge", "templates")
			if _, err := os.Stat(templatesSrc); err == nil {
				if err := copyDir(templatesSrc, templatesDst, true); err != nil {
					return fmt.Errorf("deploying templates: %w", err)
				}
			}

			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "✓ Skills updated")
			fmt.Fprintf(out, "  Skills: %s\n", skillsDst)
			fmt.Fprintf(out, "  Templates: %s\n", templatesDst)
			return nil
		},
	}

	return cmd
}
