package cmd

import (
	"fmt"
	"io"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

type lsOptions struct {
	path        string
	writer      io.Writer
	onlyKeys    bool
	onlyPaths   bool
	showSecrets bool
	json        bool
	yaml        bool
}

func defaultLsOptions() *lsOptions {
	return &lsOptions{
		path:        defaultKVPath,
		showSecrets: false,
		writer:      defaultWriter,
	}
}

func newLSCmd() *cobra.Command {
	o := defaultLsOptions()

	cmd := &cobra.Command{
		Use:           "ls",
		Short:         "recursively list secrets from Vaults KV2 engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			if err := v.ListRecursive(utils.SplitPath(o.path)); err != nil {
				return err
			}

			printer := printer.NewPrinter(v.Secrets,
				printer.OnlyKeys(o.onlyKeys),
				printer.OnlyPaths(o.onlyPaths),
				printer.ShowSecrets(o.showSecrets),
				printer.ToJSON(o.json),
				printer.ToYAML(o.yaml),
			)

			if err := printer.Out(); err != nil {
				return err
			}

			return nil
		},
	}

	// Input
	cmd.Flags().StringVarP(&o.path, "path", "p", o.path, "path")

	// Modify
	cmd.Flags().BoolVar(&o.onlyKeys, "only-keys", o.onlyKeys, "print only keys")
	cmd.Flags().BoolVar(&o.onlyPaths, "only-paths", o.onlyPaths, "print only paths")
	cmd.Flags().BoolVar(&o.showSecrets, "show-secrets", o.showSecrets, "print out secrets")

	// Output format
	cmd.Flags().BoolVarP(&o.json, "to-json", "j", o.json, "print secrets in json format")
	cmd.Flags().BoolVarP(&o.yaml, "to-yaml", "y", o.json, "print secrets in yaml format")

	return cmd
}

func (o *lsOptions) validateFlags() error {
	if o.json && o.yaml {
		return fmt.Errorf("cannot specify both --to-json and --to-yaml")
	}

	if o.onlyKeys && o.showSecrets {
		return fmt.Errorf("cannot specify both --only-keys and --show-secrets")
	}

	if o.onlyPaths && o.showSecrets {
		return fmt.Errorf("cannot specify both --only-paths and --show-secrets")
	}

	if o.onlyKeys && o.onlyPaths {
		return fmt.Errorf("cannot specify both --only-keys and --only-paths")
	}

	return nil
}
