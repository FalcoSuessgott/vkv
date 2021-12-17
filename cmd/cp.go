package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"

	"github.com/spf13/cobra"
)

type cpOptions struct {
	srcPath  string
	destPath string

	srcNamespace  string
	destNamespace string
}

func defaultCpOptions() *cpOptions {
	return &cpOptions{
		srcNamespace:  defaultNamespace,
		destNamespace: defaultNamespace,
	}
}

func newCPCmd() *cobra.Command {
	o := defaultCpOptions()

	cmd := &cobra.Command{
		Use:           "cp",
		Short:         "copy secrets within Vaults KV2 engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// create vault clients incl. namespaces
			srcC, err := newClient(o.srcNamespace)
			if err != nil {
				return err
			}

			destC, err := newClient(o.destNamespace)
			if err != nil {
				return err
			}

			// check if src path exist
			// read secrets from source path
			if err := srcC.ListRecursive(utils.SplitPath(o.srcPath)); err != nil {
				return err
			}

			// create dest path if necessary
			if err := destC.EnableKV2Engine(o.destPath); err != nil {
				log.Println(err)
			}

			// fetch secrets on dest path and inform user if secrets are going to be overwritten
			// evaluate if merging to secrets is possible
			// prompt user to perform action with showing the result

			// test this
			fmt.Printf("copying from Path: %s and Namespace: %s to Path: %s Namespace: %s\n", o.srcPath, o.srcNamespace, o.destPath, o.destNamespace)
			srcPrinter := printer.NewPrinter(srcC.Secrets)

			if err := srcPrinter.Out(); err != nil {
				return err
			}

			// create dest path if necessary
			// inform user in case of overwriting
			// write src to dest

			fmt.Println()
			// s, _ := utils.SplitPath(o.srcPath)

			for k, v := range srcC.Secrets {
				if k == o.srcPath {
					continue
				}

				newPath := strings.Join(strings.Split(k, "/")[1:], "/")

				if err := destC.WriteSecrets(o.destPath, newPath, v.(map[string]interface{})); err != nil {
					return err
				}
			}

			if err := destC.ListRecursive(utils.SplitPath(o.destPath)); err != nil {
				return err
			}

			destPrinter := printer.NewPrinter(destC.Secrets)

			fmt.Println()
			if err := destPrinter.Out(); err != nil {
				return err
			}

			destC.DisableKV2Engine(o.destPath)

			return nil
		},
	}

	// Input
	cmd.Flags().StringVarP(&o.srcPath, "src-path", "s", o.srcPath, "source path")
	cmd.Flags().StringVarP(&o.destPath, "dest-path", "d", o.destPath, "destination path")

	cmd.Flags().StringVarP(&o.srcNamespace, "src-namespace", "n", o.srcNamespace, "source namespace")
	cmd.Flags().StringVarP(&o.destNamespace, "dest-namespace", "m", o.destNamespace, "destination namespace")

	return cmd
}

func (o *cpOptions) validateFlags() error {
	if o.srcPath == "" || o.destPath == "" {
		return fmt.Errorf("--src-path and --dest-path required")
	}

	return nil
}

func newClient(namespace string) (*vault.Vault, error) {
	var e error

	vC, err := vault.NewClient()
	if err != nil {
		e = err
	}

	if namespace != "" {
		if vC, err = vault.NewNamespacedClient(namespace); err != nil {
			return nil, err
		}
	}

	if e != nil {
		return nil, e
	}

	return vC, nil
}
