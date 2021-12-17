package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	defaultKVPath    = "kv"
	defaultNamespace = ""
)

var (
	defaultWriter = os.Stdout
	v             bool
)

func newRootCmd(version string) *cobra.Command {
	cmds := &cobra.Command{
		Use:           "vkv",
		Short:         "list, copy, move, remove paths from Vaults KV engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if v {
				fmt.Printf("vkv %s\n", version)
			}

			return nil
		},
	}

	cmds.Flags().BoolVarP(&v, "version", "v", v, "display version")

	cmds.AddCommand(newLSCmd())
	cmds.AddCommand(newCPCmd())

	return cmds
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}
