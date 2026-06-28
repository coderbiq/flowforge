package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSkillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage FlowForge skills",
		Long:  `Manage FlowForge skills. Use "flowforge assets update" for all managed assets.`,
	}
	cmd.AddCommand(newSkillUpdateCmd())
	return cmd
}

func newSkillUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "update",
		Short:      "Alias for flowforge assets update",
		Long:       `Alias for "flowforge assets update". Kept for compatibility.`,
		Args:       cobra.NoArgs,
		Deprecated: "use \"flowforge assets update\" instead",
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
