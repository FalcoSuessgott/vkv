package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

const (
	envVarPrefix = "VKV_"
)

var (
	errInvalidFlagCombination = fmt.Errorf("invalid flag combination specified")
	errInvalidFormat          = fmt.Errorf("invalid format (valid options: base, yaml, json, export, markdown)")
)

// Options holds all available commandline options.
type Options struct {
	Paths []string `env:"PATHS" envDefault:"kv" envSeparator:","`

	OnlyKeys       bool `env:"ONLY_KEYS"`
	OnlyPaths      bool `env:"ONLY_PATHS"`
	ShowValues     bool `env:"SHOW_VALUES"`
	MaxValueLength int  `env:"MAX_VALUE_LENGTH" envDefault:"12"`

	FormatString string `env:"FORMAT" envDefault:"base"`
	outputFormat printer.OutputFormat

	version bool
}

func newRootCmd(version string) *cobra.Command {
	o := &Options{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

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

			for _, p := range o.Paths {
				if err := v.ListRecursive(utils.SplitPath(p)); err != nil {
					return err
				}
			}

			printer := printer.NewPrinter(
				printer.OnlyKeys(o.OnlyKeys),
				printer.OnlyPaths(o.OnlyPaths),
				printer.CustomValueLength(o.MaxValueLength),
				printer.ShowSecrets(o.ShowValues),
				printer.ToFormat(o.outputFormat),
			)

			if err := printer.Out(v.Secrets); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringSliceVarP(&o.Paths, "path", "p", o.Paths, "kv engine mount paths (comma separated for specifying multiple paths)")

	// Modify
	cmd.Flags().BoolVar(&o.OnlyKeys, "only-keys", o.OnlyKeys, "show only keys")
	cmd.Flags().BoolVar(&o.OnlyPaths, "only-paths", o.OnlyPaths, "show only paths")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "dont mask values")
	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values (-1 for disable)")

	// Output format
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "output format (options: base, json, yaml, export, nmarkdown)")

	// version
	cmd.Flags().BoolVarP(&o.version, "version", "v", o.version, "display version")

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

//nolint: cyclop
func (o *Options) validateFlags() error {
	switch {
	case (o.OnlyKeys && o.ShowValues), (o.OnlyPaths && o.ShowValues), (o.OnlyKeys && o.OnlyPaths):
		return errInvalidFlagCombination
	case true:
		switch strings.ToLower(o.FormatString) {
		case "yaml", "yml":
			o.outputFormat = printer.YAML
		case "json":
			o.outputFormat = printer.JSON
		case "export":
			o.outputFormat = printer.Export
			o.ShowValues = true
		case "markdown":
			o.outputFormat = printer.Markdown
		case "base":
			o.outputFormat = printer.Base
		default:
			return errInvalidFormat
		}
	}

	return nil
}

func (o *Options) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarPrefix,
	}); err != nil {
		return err
	}

	return nil
}
