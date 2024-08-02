package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

// exportOptions holds all available commandline options.
type exportOptions struct {
	Path       string `env:"PATH"`
	EnginePath string `env:"ENGINE_PATH"`

	SkipErrors bool `env:"SKIP_ERRORS" envDefault:"false"`

	ShowValues bool `env:"SHOW_VALUES"`

	AllVersions bool `env:"ALL_VERSIONS"`

	TemplateFile   string `env:"TEMPLATE_FILE"`
	TemplateString string `env:"TEMPLATE_STRING"`

	FormatString  string `env:"FORMAT" envDefault:"default"`
	FormatOptions []vault.Option
}

// NewExportCmd export subcommand.
//
//nolint:lll
func NewExportCmd() *cobra.Command {
	o := &exportOptions{}

	if err := utils.ParseEnvs(envVarExportPrefix, o); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "export",
		Short:         "recursively list secrets from Vaults KV2 engine in various formats",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE:       o.validateFlags,
		RunE: func(cmd *cobra.Command, args []string) error {
			rootPath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

			// TODO: consider passing the context through

			kv, err := vaultClient.NewKVSecrets(rootPath, subPath, o.SkipErrors, o.AllVersions)
			if err != nil {
				return err
			}

			opts := prt.DefaultPrinterOptions()
			opts.Format = o.FormatString

			if err := prt.Print(kv.PrinterFuncs(vault.NewFormatOptions(o.FormatOptions...)), opts); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, fmt.Sprintf("KV Engine path (env: %s)", envVarExportPrefix+"_PATH"))
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath, "engine path in case your KV-engine contains special characters such as \"/\", the path value will then be appended if specified (\"<engine-path>/<path>\") (env: VKV_EXPORT_ENGINE_PATH)")
	cmd.Flags().BoolVar(&o.SkipErrors, "skip-errors", o.SkipErrors, "don't exit on errors (permission denied, deleted secrets) (env: VKV_EXPORT_SKIP_ERRORS)")
	cmd.Flags().BoolVarP(&o.AllVersions, "all-versions", "a", o.AllVersions, "wether to fetch all KVv2 secret versions (env: VKV_EXPORT_ALL_VERSIONS)")

	// Modify
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "don't mask values (env: VKV_EXPORT_SHOW_VALUES)")

	// Template
	cmd.Flags().StringVar(&o.TemplateFile, "template-file", o.TemplateFile, "path to a file containing Go-template syntax to render the KV entries (env: VKV_EXPORT_TEMPLATE_FILE)")
	cmd.Flags().StringVar(&o.TemplateString, "template-string", o.TemplateString, "template string containing Go-template syntax to render KV entries (env: VKV_EXPORT_TEMPLATE_STRING)")

	// Output format
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "one of the allowed output formats env: VKV_EXPORT_FORMAT)")

	return cmd
}

// nolint: cyclop, goconst
func (o *exportOptions) validateFlags(cmd *cobra.Command, args []string) error {
	switch {
	case o.EnginePath == "" && o.Path == "":
		return errors.New("no KV-paths given. Either --engine-path / -e or --path / -p needs to be specified")
	case o.EnginePath != "" && o.Path != "":
		return errors.New("cannot specify both engine-path and path")
	}

	switch strings.ToLower(o.FormatString) {
	// for certain output formats we dont mask secrets
	case "yaml", "yml", "json", "export", "policy":
		o.ShowValues = true
	// for some we make it configurable
	case "markdown", "default":
		if !o.ShowValues {
			o.FormatOptions = append(o.FormatOptions, vault.MaskSecrets())
		}
	case "full":
		o.AllVersions = true

		if !o.ShowValues {
			o.FormatOptions = append(o.FormatOptions, vault.MaskSecrets())
		}

		o.FormatOptions = append(o.FormatOptions, vault.ShowDiff())
	case "template", "tmpl":
		o.ShowValues = true

		if o.TemplateFile != "" && o.TemplateString != "" {
			return fmt.Errorf("%w: %s", errInvalidFlagCombination, "cannot specify both --template-file and --template-string")
		}

		if o.TemplateFile == "" && o.TemplateString == "" {
			return fmt.Errorf("%w: %s", errInvalidFlagCombination, "either --template-file or --template-string is required")
		}

		o.FormatOptions = append(o.FormatOptions, vault.WithTemplate(o.TemplateString, o.TemplateFile))
	default:
		return prt.ErrInvalidPrinterFormat
	}

	return nil
}
