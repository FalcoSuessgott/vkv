package cmd

import (
	"fmt"
	"log"
	"path"
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

var errInvalidFlagCombination = fmt.Errorf("invalid flag combination specified")

// Options holds all available commandline options.
type Options struct {
	Path       string `env:"PATH" envDefault:"kv"`
	EnginePath string `env:"ENGINE_PATH" envDefault:""`

	OnlyKeys       bool `env:"ONLY_KEYS"`
	OnlyPaths      bool `env:"ONLY_PATHS"`
	ShowValues     bool `env:"SHOW_VALUES"`
	MaxValueLength int  `env:"MAX_VALUE_LENGTH" envDefault:"12"`

	TemplateFile   string `env:"TEMPLATE_FILE"`
	TemplateString string `env:"TEMPLATE_STRING"`

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
				fmt.Fprintf(cmd.OutOrStdout(), "vkv %s\n", version)

				return nil
			}

			if err := o.validateFlags(); err != nil {
				return err
			}

			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			m := map[string]interface{}{}
			s := &vault.Secrets{}

			rootPath, subPath := o.buildEnginePath()
			if err := s.ListRecursive(v, rootPath, subPath); err != nil {
				return fmt.Errorf("error reading secrets: %w", err)
			}

			m[path.Join(rootPath, subPath)] = (*s)

			if len(m) == 0 {
				return nil
			}

			printer := printer.NewPrinter(
				printer.OnlyKeys(o.OnlyKeys),
				printer.OnlyPaths(o.OnlyPaths),
				printer.CustomValueLength(o.MaxValueLength),
				printer.ShowValues(o.ShowValues),
				printer.WithTemplate(o.TemplateString, o.TemplateFile),
				printer.ToFormat(o.outputFormat),
			)

			if err := printer.Out(m); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KVv2 Engine path (env var: VKV_PATH)")
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath,
		"Specify the engine path using this flag in case your kv-engine contains special characters such as \"/\".\n"+
			"vkv will then append the values of the path-flag to the engine path, if specified (<engine-path>/<path>)"+
			"(env var: VKV_ENGINE_PATHS)")

	// Modify
	cmd.Flags().BoolVar(&o.OnlyKeys, "only-keys", o.OnlyKeys, "show only keys (env var: VKV_ONLY_KEYS)")
	cmd.Flags().BoolVar(&o.OnlyPaths, "only-paths", o.OnlyPaths, "show only paths (env var: VKV_ONLY_PATHS)")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "dont mask values (env var: VKV_SHOW_VALUES)")
	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values. Set to \"-1\" for disabling "+
		"(env var: VKV_MAX_VALUE_LENGTH)")

	// Template
	cmd.Flags().StringVar(&o.TemplateFile, "template-file", o.TemplateFile, "path to a file containing Go-template syntax to render the KV entries (env var: VKV_TEMPLATE_FILE)")
	cmd.Flags().StringVar(&o.TemplateString, "template-string", o.TemplateString, "template string containing Go-template syntax to render KV entries (env var: VKV_TEMPLATE_STRING)")

	// Output format
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "output format: \"base\", \"json\", \"yaml\", \"export\", \"markdown\", \"template\") "+
		"(env var: VKV_FORMAT)")

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

// nolint: cyclop
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
			o.OnlyKeys = false
			o.OnlyPaths = false
			o.ShowValues = true
		case "markdown":
			o.outputFormat = printer.Markdown
		case "base":
			o.outputFormat = printer.Base
		case "template", "tmpl":
			o.outputFormat = printer.Template
			o.OnlyKeys = false
			o.OnlyPaths = false

			if o.TemplateFile != "" && o.TemplateString != "" {
				return fmt.Errorf("%w: %s", errInvalidFlagCombination, "cannot specify both --template-file and --template-string")
			}

			if o.TemplateFile == "" && o.TemplateString == "" {
				return fmt.Errorf("%w: %s", errInvalidFlagCombination, "either --template-file or --template-string is required")
			}
		default:
			return printer.ErrInvalidFormat
		}
	case (o.EnginePath == "" && o.Path == ""):
		return fmt.Errorf("no paths KVv2 paths given. Either --engine-path (-e) or --path (-p) needs to be specified")
	}

	return nil
}

func (o *Options) buildEnginePath() (string, string) {
	// if engine path has been specified use that value as the root path and append the path
	if o.EnginePath != "" {
		return o.EnginePath, o.Path
	}

	return utils.SplitPath(o.Path)
}

func (o *Options) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarPrefix,
	}); err != nil {
		return err
	}

	return nil
}
