package cmd

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"path"
	"strings"
	"time"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/spf13/cobra"
)

// exportOptions holds all available commandline options.
type exportOptions struct {
	Path       string `env:"PATH"`
	EnginePath string `env:"ENGINE_PATH"`

	OnlyKeys       bool `env:"ONLY_KEYS"`
	OnlyPaths      bool `env:"ONLY_PATHS"`
	ShowValues     bool `env:"SHOW_VALUES"`
	ShowVersion    bool `env:"SHOW_VERSION" envDefault:"true"`
	ShowMetadata   bool `env:"SHOW_METADATA" envDefault:"true"`
	WithHyperLink  bool `env:"WITH_HYPERLINK" envDefault:"true"`
	MaxValueLength int  `env:"MAX_VALUE_LENGTH" envDefault:"12"`

	SkipErrors bool `env:"SKIP_ERRORS" envDefault:"false"`

	TemplateFile   string `env:"TEMPLATE_FILE"`
	TemplateString string `env:"TEMPLATE_STRING"`

	FormatString string `env:"FORMAT" envDefault:"base"`

	outputFormat prt.OutputFormat
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
			enginePath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

			printer = prt.NewSecretPrinter(
				prt.OnlyKeys(o.OnlyKeys),
				prt.OnlyPaths(o.OnlyPaths),
				prt.CustomValueLength(o.MaxValueLength),
				prt.ShowValues(o.ShowValues),
				prt.WithTemplate(o.TemplateString, o.TemplateFile),
				prt.ToFormat(o.outputFormat),
				prt.WithVaultClient(vaultClient),
				prt.WithWriter(writer),
				prt.ShowVersion(o.ShowVersion),
				prt.ShowMetadata(o.ShowMetadata),
				prt.WithHyperLinks(o.WithHyperLink),
				prt.WithEnginePath(utils.NormalizePath(enginePath)),
			)

			secrets, err := vaultClient.ListRecursive(enginePath, subPath, o.SkipErrors)
			if err != nil {
				return err
			}

			go func() {
				slog.Info("Starting token refresher")
				vaultClient.LeaseRefresher(cmd.Context())
			}()

			time.Sleep(60 * time.Second) // wait for the refresher to start

			p := path.Join(enginePath, subPath)
			if subPath == "" {
				p = utils.NormalizePath(p)
			}

			result := utils.UnflattenMap(p, utils.ToMapStringInterface(secrets), o.EnginePath)

			if err := printer.Out(result); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KV Engine path (env: VKV_EXPORT_PATH")
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath, "engine path in case your KV-engine contains special characters such as \"/\", the path (-p) flag will then be appended if specified (\"<engine-path>/<path>\") (env: VKV_EXPORT_ENGINE_PATH)")
	cmd.Flags().BoolVar(&o.SkipErrors, "skip-errors", o.SkipErrors, "don't exit on errors (permission denied, deleted secrets) (env: VKV_EXPORT_SKIP_ERRORS)")

	// Modify
	cmd.Flags().BoolVar(&o.OnlyKeys, "only-keys", o.OnlyKeys, "show only keys (env: VKV_EXPORT_ONLY_KEYS)")
	cmd.Flags().BoolVar(&o.OnlyPaths, "only-paths", o.OnlyPaths, "show only paths (env: VKV_EXPORT_ONLY_PATHS)")
	cmd.Flags().BoolVar(&o.ShowVersion, "show-version", o.ShowVersion, "show the secret version (env: VKV_EXPORT_VERSION)")
	cmd.Flags().BoolVar(&o.ShowMetadata, "show-metadata", o.ShowMetadata, "show the secrets metadata (env: VKV_EXPORT_METADATA)")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "don't mask values (env: VKV_EXPORT_SHOW_VALUES)")
	cmd.Flags().BoolVar(&o.WithHyperLink, "with-hyperlink", o.WithHyperLink, "don't link to the Vault UI (env: VKV_EXPORT_WITH_HYPERLINK)")

	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values. Set to \"-1\" for disabling "+
		"(env: VKV_EXPORT_MAX_VALUE_LENGTH)")

	// Template
	cmd.Flags().StringVar(&o.TemplateFile, "template-file", o.TemplateFile, "path to a file containing Go-template syntax to render the KV entries (env: VKV_EXPORT_TEMPLATE_FILE)")
	cmd.Flags().StringVar(&o.TemplateString, "template-string", o.TemplateString, "template string containing Go-template syntax to render KV entries (env: VKV_EXPORT_TEMPLATE_STRING)")

	// Output format
	//nolint: lll
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\", \"export\", \"policy\", \"markdown\", \"template\" "+
		"(env: VKV_EXPORT_FORMAT)")

	return cmd
}

// nolint: cyclop, goconst
func (o *exportOptions) validateFlags(cmd *cobra.Command, args []string) error {
	switch {
	case (o.OnlyKeys && o.ShowValues), (o.OnlyPaths && o.ShowValues), (o.OnlyKeys && o.OnlyPaths):
		return errInvalidFlagCombination
	case o.EnginePath == "" && o.Path == "":
		return errors.New("no KV-paths given. Either --engine-path/-e or --path/-p needs to be specified")
	case true:
		switch strings.ToLower(o.FormatString) {
		case "yaml", "yml":
			o.outputFormat = prt.YAML
			o.OnlyKeys = false
			o.OnlyPaths = false
			o.MaxValueLength = -1
			o.ShowValues = true
		case "json":
			o.outputFormat = prt.JSON
			o.OnlyKeys = false
			o.OnlyPaths = false
			o.MaxValueLength = -1
			o.ShowValues = true
		case "export":
			o.outputFormat = prt.Export
			o.OnlyKeys = false
			o.OnlyPaths = false
			o.ShowValues = true
			o.MaxValueLength = -1
		case "markdown":
			o.outputFormat = prt.Markdown
		case "base":
			o.outputFormat = prt.Base
		case "policy":
			o.outputFormat = prt.Policy
			o.OnlyKeys = false
			o.OnlyPaths = false
			o.ShowValues = true
		case "template", "tmpl":
			o.outputFormat = prt.Template
			o.OnlyKeys = false
			o.OnlyPaths = false
			o.ShowValues = true
			o.MaxValueLength = -1

			if o.TemplateFile != "" && o.TemplateString != "" {
				return fmt.Errorf("%w: %s", errInvalidFlagCombination, "cannot specify both --template-file and --template-string")
			}

			if o.TemplateFile == "" && o.TemplateString == "" {
				return fmt.Errorf("%w: %s", errInvalidFlagCombination, "either --template-file or --template-string is required")
			}
		default:
			return prt.ErrInvalidFormat
		}
	}

	return nil
}
