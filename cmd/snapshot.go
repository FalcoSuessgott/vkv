package cmd

import (
	"github.com/spf13/cobra"
)

// NewSnapshotCmd snaphot command.
func NewSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "snapshot",
		Short:         "save or restore a snapshot of all KVv2 engines",
		Aliases:       []string{"ss"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		NewSnapshotSaveCmd(),
		NewSnapshotRestoreCmd(),
	)

	return cmd
}
