package cmd

import (
	"github.com/spf13/cobra"
)

// NewListCmd holds the list subcommands.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "list namespaces or KV engines",
		Aliases:       []string{"ls"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newListNamespacesCmd(),
		newListEngineCmd(),
	)

	return cmd
}
