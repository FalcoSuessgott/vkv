package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
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
	Path       string `env:"PATH"`
	EnginePath string `env:"ENGINE_PATH"`

	OnlyKeys       bool `env:"ONLY_KEYS"`
	OnlyPaths      bool `env:"ONLY_PATHS"`
	ShowValues     bool `env:"SHOW_VALUES"`
	MaxValueLength int  `env:"MAX_VALUE_LENGTH" envDefault:"12"`

	TemplateFile   string `env:"TEMPLATE_FILE"`
	TemplateString string `env:"TEMPLATE_STRING"`

	FormatString string `env:"FORMAT" envDefault:"base"`

	outputFormat printer.OutputFormat
}

func newRootCmd(version string, writer io.Writer) *cobra.Command {
	o := &Options{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "vkv -p <kv-path>",
		Short:         "recursively list secrets from Vaults KV2 engine in various formats",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// validate flags
			if err := o.validateFlags(); err != nil {
				return err
			}

			// vault auth
			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			// prepare printer
			printer := printer.NewPrinter(
				printer.OnlyKeys(o.OnlyKeys),
				printer.OnlyPaths(o.OnlyPaths),
				printer.CustomValueLength(o.MaxValueLength),
				printer.ShowValues(o.ShowValues),
				printer.WithTemplate(o.TemplateString, o.TemplateFile),
				printer.ToFormat(o.outputFormat),
				printer.WithVaultClient(v),
				printer.WithWriter(writer),
			)

			// prepare map
			m, err := o.buildMap(v)
			if err != nil {
				return err
			}

			// print secrets
			if err := printer.Out(m); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KVv2 Engine path (env: VKV_PATH)")
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath,
		"Specify the engine path using this flag in case your kv-engine contains special characters such as \"/\".\n"+
			"vkv will then append the values of the path-flag to the engine path, if specified (<engine-path>/<path>)"+
			"(env: VKV_ENGINE_PATH)")

	// Modify
	cmd.Flags().BoolVar(&o.OnlyKeys, "only-keys", o.OnlyKeys, "show only keys (env: VKV_ONLY_KEYS)")
	cmd.Flags().BoolVar(&o.OnlyPaths, "only-paths", o.OnlyPaths, "show only paths (env: VKV_ONLY_PATHS)")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "don't mask values (env: VKV_SHOW_VALUES)")
	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values. Set to \"-1\" for disabling "+
		"(env: VKV_MAX_VALUE_LENGTH)")

	// Template
	cmd.Flags().StringVar(&o.TemplateFile, "template-file", o.TemplateFile, "path to a file containing Go-template syntax to render the KV entries (env: VKV_TEMPLATE_FILE)")
	cmd.Flags().StringVar(&o.TemplateString, "template-string", o.TemplateString, "template string containing Go-template syntax to render KV entries (env: VKV_TEMPLATE_STRING)")

	// Output format
	//nolint: lll
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\", \"export\", \"policy\", \"markdown\", \"template\" "+
		"(env: VKV_FORMAT)")

	// sub commands
	cmd.AddCommand(
		newVersionCmd(version),
		newImportCmd(),
	)

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version, os.Stdout).Execute(); err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}

// nolint: cyclop
func (o *Options) validateFlags() error {
	switch {
	case (o.OnlyKeys && o.ShowValues), (o.OnlyPaths && o.ShowValues), (o.OnlyKeys && o.OnlyPaths):
		return errInvalidFlagCombination
	case o.EnginePath == "" && o.Path == "":
		return fmt.Errorf("no KV-paths given. Either --engine-path / -e or --path / -p needs to be specified")
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
		case "policy":
			o.outputFormat = printer.Policy
			o.OnlyKeys = false
			o.OnlyPaths = false
			o.ShowValues = true
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
	}

	return nil
}

func (o *Options) buildMap(v *vault.Vault) (map[string]interface{}, error) {
	var isSecretPath bool

	rootPath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

	// read recursive all secrets
	s, err := v.ListRecursive(rootPath, subPath)
	if err != nil {
		return nil, err
	}

	// check if path is a directory or secret path
	if _, isSecret := v.ReadSecrets(rootPath, subPath); isSecret == nil {
		isSecretPath = true
	}

	path := path.Join(rootPath, subPath)
	if o.EnginePath != "" {
		path = subPath
	}

	// prepare the output map
	pathMap := utils.PathMap(path, utils.ToMapStringInterface(s), isSecretPath)
	m := pathMap

	if o.EnginePath != "" {
		m[o.EnginePath] = pathMap
	}

	return m, nil
}

func (o *Options) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarPrefix,
	}); err != nil {
		return err
	}

	return nil
}
