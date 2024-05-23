package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
)

var fOpts = []vault.Option{}

// exportOptions holds all available commandline options.
type exportOptions struct {
	Path       string `env:"PATH"`
	EnginePath string `env:"ENGINE_PATH"`

	OnlyKeys       bool `env:"ONLY_KEYS"`
	OnlyPaths      bool `env:"ONLY_PATHS"`
	ShowValues     bool `env:"SHOW_VALUES"`
	WithHyperLink  bool `env:"WITH_HYPERLINK" envDefault:"true"`
	MaxValueLength int  `env:"MAX_VALUE_LENGTH" envDefault:"12"`
	ShowDiff       bool `env:"SHOW_DIFF"  envDefault:"true"`

	SkipErrors bool `env:"SKIP_ERRORS" envDefault:"false"`

	PrintLegend bool `env:"PRINT_LEGEND" envDefault:"true"`

	TemplateFile   string `env:"TEMPLATE_FILE"`
	TemplateString string `env:"TEMPLATE_STRING"`

	FormatString string `env:"FORMAT" envDefault:"default"`
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

			// TODO flag for all versions
			kv, err := vaultClient.NewKVSecrets(rootPath, subPath, o.SkipErrors, true)
			if err != nil {
				return err
			}

			opts := prt.DefaultPrinterOptions()
			opts.Format = o.FormatString

			pp.Println(kv.Secrets)
			pp.Println("")

			formatOptions := vault.NewFormatOptions(fOpts...)

			fmt.Printf("%#v\n", formatOptions)

			if err := prt.Print(kv.PrinterFuncs(formatOptions), opts); err != nil {
				return err
			}

			if o.FormatString == "default" && o.PrintLegend {
				fmt.Printf("[ ] = no changes\n%s = added\n%s = changed\n%s = removed\n",
					utils.ColorGreen("[+]"),
					utils.ColorYellow("[~]"),
					utils.ColorRed("[-]"),
				)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, fmt.Sprintf("KV Engine path (env: %s)", envVarExportPrefix+"_PATH"))
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath, "engine path in case your KV-engine contains special characters such as \"/\", the path value will then be appended if specified (\"<engine-path>/<path>\") (env: VKV_EXPORT_ENGINE_PATH)")
	cmd.Flags().BoolVar(&o.SkipErrors, "skip-errors", o.SkipErrors, "don't exit on errors (permission denied, deleted secrets) (env: VKV_EXPORT_SKIP_ERRORS)")

	// Modify
	cmd.Flags().BoolVar(&o.OnlyKeys, "only-keys", o.OnlyKeys, "show only keys (env: VKV_EXPORT_ONLY_KEYS)")
	cmd.Flags().BoolVar(&o.OnlyPaths, "only-paths", o.OnlyPaths, "show only paths (env: VKV_EXPORT_ONLY_PATHS)")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "don't mask values (env: VKV_EXPORT_SHOW_VALUES)")
	cmd.Flags().BoolVar(&o.WithHyperLink, "with-hyperlink", o.WithHyperLink, "don't link to the Vault UI (env: VKV_EXPORT_WITH_HYPERLINK)")
	cmd.Flags().BoolVar(&o.ShowDiff, "show-diff", o.WithHyperLink, "when enabled highlights the diff for each secret version (env: VKV_EXPORT_SHOW_DIFF)")

	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values. Set to \"-1\" for disabling "+
		"(env: VKV_EXPORT_MAX_VALUE_LENGTH)")

	// Template
	cmd.Flags().StringVar(&o.TemplateFile, "template-file", o.TemplateFile, "path to a file containing Go-template syntax to render the KV entries (env: VKV_EXPORT_TEMPLATE_FILE)")
	cmd.Flags().StringVar(&o.TemplateString, "template-string", o.TemplateString, "template string containing Go-template syntax to render KV entries (env: VKV_EXPORT_TEMPLATE_STRING)")

	cmd.Flags().BoolVarP(&o.PrintLegend, "legend", "l", o.PrintLegend, "wether to print a legend (env: VKV_EXPORT_PRINT_LEGEND)")

	// Output format
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "one of the allowed output formats env: VKV_EXPORT_FORMAT)")

	return cmd
}

// nolint: cyclop, goconst
func (o *exportOptions) validateFlags(cmd *cobra.Command, args []string) error {
	switch {
	case (o.OnlyKeys && o.ShowValues), (o.OnlyPaths && o.ShowValues), (o.OnlyKeys && o.OnlyPaths):
		return errInvalidFlagCombination
	case o.EnginePath == "" && o.Path == "":
		return errors.New("no KV-paths given. Either --engine-path / -e or --path / -p needs to be specified")
	case o.EnginePath != "" && o.Path != "":
		return errors.New("cannot specify both engine-path and path")
	}

	switch strings.ToLower(o.FormatString) {
	case "yaml", "yml":
		o.OnlyKeys = false
		o.OnlyPaths = false
		o.MaxValueLength = -1
		o.ShowValues = true
	case "json":
		o.OnlyKeys = false
		o.OnlyPaths = false
		o.MaxValueLength = -1
		o.ShowValues = true
	case "export":
		o.OnlyKeys = false
		o.OnlyPaths = false
		o.ShowValues = true
		o.MaxValueLength = -1
	case "markdown":
	case "default":
		if o.ShowDiff {
			fOpts = append(fOpts, vault.ShowDiff())
		}

		if o.OnlyKeys {
			fOpts = append(fOpts, vault.OnlyKeys())
		}

		if !o.ShowValues {
			fOpts = append(fOpts, vault.MaskSecrets())
		}

	case "policy":
		o.OnlyKeys = false
		o.OnlyPaths = false
		o.ShowValues = true
	case "template", "tmpl":
		o.OnlyKeys = false
		o.OnlyPaths = false
		o.MaxValueLength = -1

		if o.TemplateFile != "" && o.TemplateString != "" {
			return fmt.Errorf("%w: %s", errInvalidFlagCombination, "cannot specify both --template-file and --template-string")
		}

		if o.TemplateFile == "" && o.TemplateString == "" {
			return fmt.Errorf("%w: %s", errInvalidFlagCombination, "either --template-file or --template-string is required")
		}
	default:
		return prt.ErrInvalidPrinterFormat
	}

	return nil
}
