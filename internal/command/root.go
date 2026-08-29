package command

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root Cobra command for FlowForge.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flowforge",
		Short: "FlowForge: Local-first issue tracker and DAG engine for FlowForge engineering skills",
		Long: `FlowForge provides deterministic DAG dependency calculation,
frontier task queue extraction, and local-first issue tracking
for the FlowForge engineering skills methodology.`,
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newInitCmd(),
		newConfigCmd(),
		newFrontierCmd(),
		newCheckCmd(),
		newStatusCmd(),
		newVersionCmd(),
		newUpgradeCmd(),
		newAssetsCmd(),
		newAgentsCmd(),
	)

	return cmd
}
