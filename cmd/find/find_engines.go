package find

import (
	"io"
	"log"
	"strings"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/engine"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

const envVarFindEnginesPrefix = "VKV_FIND_ENGINES_"

type findEnginesOptions struct {
	Namespace string `env:"NS"`
	Prefix    bool   `env:"NS_PREFIX"`

	Regex string `env:"REGEX"`
	All   bool   `env:"ALL"`

	FormatString string `env:"FORMAT" envDefault:"base"`

	writer       io.Writer
	outputFormat printer.OutputFormat
}

func newFindEngineCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
	var err error

	o := &findEnginesOptions{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "engines",
		Short:         "find all KV engines",
		Aliases:       []string{"e", "eng"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var engines vault.Engines

			if err := o.validateFlags(); err != nil {
				return err
			}

			// vault auth
			if vaultClient == nil {
				if vaultClient, err = vault.NewDefaultClient(); err != nil {
					return err
				}
			}

			if !o.All {
				if engines, err = o.findEngines(vaultClient); err != nil {
					return err
				}
			} else {
				if engines, err = o.findAllEngines(vaultClient); err != nil {
					return err
				}
			}

			return printer.NewPrinter(
				printer.ToFormat(o.outputFormat),
				printer.WithWriter(o.writer),
				printer.WithRegex(o.Regex),
				printer.WithNSPrefix(o.Prefix),
			).Out(engines)
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "specify the namespace (env: VKV_find_ENGINES_NS)")
	cmd.Flags().BoolVarP(&o.Prefix, "include-ns-prefix", "p", o.Prefix, "prepend the namespaces (env: VKV_find_ENGINES_NS_PREFIX)")

	cmd.Flags().StringVarP(&o.Regex, "regex", "r", o.Regex, "filter engines by the specified regex pattern (env: VKV_find_ENGINES_REGEX")
	cmd.Flags().BoolVarP(&o.All, "all", "a", o.All, "find all KV engines recursively from the specified namespaces (env: VKV_find_ENGINES_ALL)")
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\" (env: VKV_find_ENGINES_FORMAT)")

	o.writer = writer

	return cmd
}

func (o *findEnginesOptions) validateFlags() error {
	switch strings.ToLower(o.FormatString) {
	case "yaml", "yml":
		o.outputFormat = printer.YAML
	case "json":
		o.outputFormat = printer.JSON
	case "base":
		o.outputFormat = printer.Base
	default:
		return printer.ErrInvalidFormat
	}

	return nil
}

func (o *findEnginesOptions) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarFindEnginesPrefix,
	}); err != nil {
		return err
	}

	return nil
}

func (o *findEnginesOptions) findEngines(v *vault.Vault) (vault.Engines, error) {
	engines, err := v.ListKVSecretEngines(o.Namespace)
	if err != nil {
		return nil, err
	}

	m := make(vault.Engines)
	m[o.Namespace] = engines

	return m, nil
}

func (o *findEnginesOptions) findAllEngines(v *vault.Vault) (vault.Engines, error) {
	engines, err := v.ListAllKVSecretEngines(o.Namespace)
	if err != nil {
		return nil, err
	}

	return engines, nil
}
