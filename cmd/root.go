package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

const defaultKVPath = "kv"

var defaultWriter = os.Stdout

// Options holds all available commandline options.
type Options struct {
	paths       []string
	writer      io.Writer
	onlyKeys    bool
	onlyPaths   bool
	showSecrets bool
	json        bool
	yaml        bool
	version     bool
}

func defaultOptions() *Options {
	return &Options{
		paths:       []string{defaultKVPath},
		showSecrets: false,
		writer:      defaultWriter,
	}
}

func newRootCmd(version string) *cobra.Command {
	o := defaultOptions()

	cmd := &cobra.Command{
		Use:           "vkv",
		Short:         "recursively list secrets from Vaults KV2 engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.version {
				fmt.Printf("vkv %s\n", version)

				return nil
			}

			if err := o.validateFlags(); err != nil {
				return err
			}

			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			for _, p := range o.paths {
				if err := v.ListRecursive(utils.SplitPath(p)); err != nil {
					return err
				}
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
	cmd.Flags().StringSliceVarP(&o.paths, "path", "p", o.paths, "engine paths")

	// Modify
	cmd.Flags().BoolVar(&o.onlyKeys, "only-keys", o.onlyKeys, "print only keys")
	cmd.Flags().BoolVar(&o.onlyPaths, "only-paths", o.onlyPaths, "print only paths")
	cmd.Flags().BoolVar(&o.showSecrets, "show-secrets", o.showSecrets, "print out secrets")

	// Output format
	cmd.Flags().BoolVarP(&o.json, "to-json", "j", o.json, "print secrets in json format")
	cmd.Flags().BoolVarP(&o.yaml, "to-yaml", "y", o.json, "print secrets in yaml format")

	cmd.Flags().BoolVarP(&o.version, "version", "v", o.version, "display version")

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}

func (o *Options) validateFlags() error {
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
