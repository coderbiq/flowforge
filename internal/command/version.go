package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"flowforge/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of FlowForge CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("flowforge %s\n", version.Version)
		},
	}
}
