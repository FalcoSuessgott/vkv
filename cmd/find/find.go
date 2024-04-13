package find

import (
	"io"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

// NewFindCmd holds the list subcommands.
func NewFindCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "find",
		Short:         "find and list namespaces, KV engines or secrets ",
		Aliases:       []string{"f", "fd"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newFindNamespacesCmd(writer, vaultClient),
		newFindEngineCmd(writer, vaultClient),
		NewFindSecretsCmd(writer, vaultClient),
	)

	return cmd
}
