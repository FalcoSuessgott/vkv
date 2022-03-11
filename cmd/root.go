package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

const (
	defaultKVPath        = "kv"
	maxValueLengthEnvVar = "VKV_MAX_VALUE_LENGTH"
)

var defaultWriter = os.Stdout

// Options holds all available commandline options.
type Options struct {
	paths          []string
	writer         io.Writer
	onlyKeys       bool
	onlyPaths      bool
	showSecrets    bool
	json           bool
	yaml           bool
	version        bool
	maxValueLength int
}

func defaultOptions() *Options {
	return &Options{
		paths:          []string{defaultKVPath},
		showSecrets:    false,
		writer:         defaultWriter,
		maxValueLength: printer.MaxValueLength,
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
				printer.CustomValueLength(o.maxValueLength),
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
	cmd.Flags().StringSliceVarP(&o.paths, "path", "p", o.paths, "kv engine paths (comma separated to define multiple paths)")

	// Modify
	cmd.Flags().BoolVar(&o.onlyKeys, "only-keys", o.onlyKeys, "print only keys")
	cmd.Flags().BoolVar(&o.onlyPaths, "only-paths", o.onlyPaths, "print only paths")
	cmd.Flags().BoolVar(&o.showSecrets, "show-secrets", o.showSecrets, "print out values")

	// Output format
	cmd.Flags().BoolVarP(&o.json, "to-json", "j", o.json, "print entries in json format")
	cmd.Flags().BoolVarP(&o.yaml, "to-yaml", "y", o.json, "print entries in yaml format")
	cmd.Flags().IntVarP(&o.maxValueLength, "max-value-length", "m",
		o.maxValueLength, "maximum char length of values (precedes VKV_MAX_PASSWORD_LENGTH)")

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

	// -m flag precedes VKV_MAX_PASSWORD_LENGTH, so we check if the flag has been provided
	if v, ok := os.LookupEnv(maxValueLengthEnvVar); ok && o.maxValueLength == printer.MaxValueLength {
		o.maxValueLength, _ = strconv.Atoi(v)
	}

	return nil
}
