package list

import (
	"io"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

// NewListCmd holds the list subcommands.
func NewListCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
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
		newListNamespacesCmd(writer, vaultClient),
		newListEngineCmd(writer, vaultClient),
	)

	return cmd
}
