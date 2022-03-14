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

var (
	defaultWriter             = os.Stdout
	errInvalidFlagCombination = fmt.Errorf("invalid flag combination specified")
	errMultipleOutputFormats  = fmt.Errorf("specified multiple output formats, only one is allowed")
)

// Options holds all available commandline options.
type Options struct {
	paths          []string
	writer         io.Writer
	onlyKeys       bool
	onlyPaths      bool
	showValues     bool
	json           bool
	yaml           bool
	markdown       bool
	version        bool
	export         bool
	maxValueLength int
}

func defaultOptions() *Options {
	return &Options{
		paths:          []string{defaultKVPath},
		showValues:     false,
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
				printer.ShowSecrets(o.showValues),
				printer.ToExportFormat(o.export),
				printer.ToJSON(o.json),
				printer.ToYAML(o.yaml),
				printer.ToMarkdown(o.markdown),
			)

			if err := printer.Out(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringSliceVarP(&o.paths, "path", "p", o.paths, "kv engine paths (comma separated to define multiple paths)")

	// Modify
	cmd.Flags().BoolVar(&o.onlyKeys, "only-keys", o.onlyKeys, "show only keys")
	cmd.Flags().BoolVar(&o.onlyPaths, "only-paths", o.onlyPaths, "show only paths")
	cmd.Flags().BoolVar(&o.showValues, "show-values", o.showValues, "dont mask values")
	cmd.Flags().IntVar(&o.maxValueLength, "max-value-length", o.maxValueLength,
		"maximum char length of values (precedes VKV_MAX_PASSWORD_LENGTH). \"-1\" for disabling")

	// Output format
	cmd.Flags().BoolVarP(&o.markdown, "markdown", "m", o.markdown, "print entries in markdown table format")
	cmd.Flags().BoolVarP(&o.json, "json", "j", o.json, "print entries in json format")
	cmd.Flags().BoolVarP(&o.yaml, "yaml", "y", o.json, "print entries in yaml format")
	cmd.Flags().BoolVarP(&o.export, "export", "e", o.export,
		"print entries in export format (export \"key=value\")")

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
	case o.json && (o.markdown || o.yaml || o.export), o.yaml && (o.markdown || o.json || o.export),
		o.export && (o.markdown || o.json || o.yaml), o.markdown && (o.json || o.yaml || o.export):
		return errMultipleOutputFormats
	case (o.onlyKeys && o.showValues), (o.onlyPaths && o.showValues), (o.onlyKeys && o.onlyPaths):
		return errInvalidFlagCombination
	}

	// -m flag precedes VKV_MAX_PASSWORD_LENGTH, so we check if the flag has been provided
	if v, ok := os.LookupEnv(maxValueLengthEnvVar); ok && o.maxValueLength == printer.MaxValueLength {
		l, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid value \"%v\" for %s", v, maxValueLengthEnvVar)
		}

		o.maxValueLength = l
	}

	return nil
}
