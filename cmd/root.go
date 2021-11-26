package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
)

type options struct {
	rootPath string
	list     bool
	backup   bool //
	restore  bool // import is a golang reserved key word
	format   string
}

func newRootCmd(version string) *cobra.Command {
	o := &options{}

	cmd := &cobra.Command{
		Use:   "vkv",
		Short: "vault kv engine exporter and importer",
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			keys, err := v.ListPath("secret")
			if err != nil {
				return err
			}

			fmt.Println("root", keys)

			fmt.Println(listRecursive(v, "secret", keys))
			return nil

		},
	}

	cmd.Flags().StringVarP(&o.rootPath, "path", "p", o.rootPath, "root path")
	cmd.Flags().BoolVarP(&o.list, "list", "l", o.list, "list all secrets from a path")
	cmd.Flags().BoolVarP(&o.backup, "export", "e", o.backup, "export all secrets in a specfied format")
	cmd.Flags().BoolVarP(&o.restore, "import", "i", o.restore, "import all secrets")

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}

func listRecursive(v *vault.Vault, rootPath string, rootKeys []string) map[string]interface{} {
	m := make(map[string]interface{})

	for _, k := range rootKeys {
		if strings.HasSuffix(k, "/") {
			subKeys, err := v.ListSubPath(rootPath, k)
			if err != nil {
				log.Fatalf("errored at %s/%s: %v", rootPath, k, err)
			}

			m[fmt.Sprintf("%s/%s", rootPath, k)] = subKeys
		}
	}

	return m
}
