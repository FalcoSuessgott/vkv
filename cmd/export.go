package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
)

type exportOptions struct {
	path string
}

func defaultExportOptions() *exportOptions {
	return &exportOptions{}
}

func newExportCmd() *cobra.Command {
	o := defaultExportOptions()

	cmd := &cobra.Command{
		Use:           "export",
		Short:         "export secrets as env vars from Vaults KV2 engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			secrets, err := v.ReadSecrets(utils.SplitPath(o.path))
			if err != nil {
				return err
			}

			for k, v := range secrets {
				fmt.Printf("export %s=\"%v\"\n", k, v)

				if err := os.Setenv(k, fmt.Sprintf("%v", v)); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.path, "path", "p", o.path, "engine paths")

	return cmd
}
