package command

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root Cobra command for FlowForge.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flowforge",
		Short: "FlowForge: Local-first issue tracker and DAG engine for mattpocock skills",
		Long: `FlowForge provides deterministic DAG dependency calculation,
frontier task queue extraction, and local-first issue tracking
for the mattpocock agile skills methodology.`,
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newInitCmd(),
		newFrontierCmd(),
		newCheckCmd(),
		newStatusCmd(),
		newVersionCmd(),
		newUpgradeCmd(),
	)

	return cmd
}
