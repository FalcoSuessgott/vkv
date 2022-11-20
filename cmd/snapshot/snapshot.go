package snapshot

import (
	"io"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

// NewSnapshotCmd snaphot command.
func NewSnapshotCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
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
		newSnapshotSaveCmd(writer, vaultClient),
		newSnapshotRestoreCmd(writer, nil),
	)

	return cmd
}
