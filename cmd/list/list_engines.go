package list

import (
	"io"
	"log"
	"strings"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/engine"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

const envVarListEnginesPrefix = "VKV_LIST_ENGINES_"

type listEnginesOptions struct {
	Namespace string `env:"NS"`
	Prefix    bool   `env:"NS_PREFIX"`

	Regex string `env:"REGEX"`
	All   bool   `env:"ALL"`

	FormatString string `env:"FORMAT" envDefault:"base"`

	writer       io.Writer
	outputFormat printer.OutputFormat
}

func newListEngineCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
	var err error

	o := &listEnginesOptions{}

	if err := utils.ParseEnvs(envVarListEnginesPrefix, o); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "engines",
		Short:         "list all KVv2 engines",
		Aliases:       []string{"e"},
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
				if engines, err = o.listEngines(vaultClient); err != nil {
					return err
				}
			} else {
				if engines, err = o.listAllEngines(vaultClient); err != nil {
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

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "specify the namespace (env: VKV_LIST_ENGINES_NS)")
	cmd.Flags().BoolVarP(&o.Prefix, "include-ns-prefix", "p", o.Prefix, "prepend the namespaces (env: VKV_LIST_ENGINES_NS_PREFIX)")

	cmd.Flags().StringVarP(&o.Regex, "regex", "r", o.Regex, "filter engines by the specified regex pattern (env: VKV_LIST_ENGINES_REGEX")
	cmd.Flags().BoolVarP(&o.All, "all", "a", o.All, "list all KV engines recursively from the specified namespaces (env: VKV_LIST_ENGINES_ALL)")
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\" (env: VKV_LIST_ENGINES_FORMAT)")

	o.writer = writer

	return cmd
}

func (o *listEnginesOptions) validateFlags() error {
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

func (o *listEnginesOptions) listEngines(v *vault.Vault) (vault.Engines, error) {
	engines, err := v.ListKVSecretEngines(o.Namespace)
	if err != nil {
		return nil, err
	}

	m := make(vault.Engines)
	m[o.Namespace] = engines

	return m, nil
}

func (o *listEnginesOptions) listAllEngines(v *vault.Vault) (vault.Engines, error) {
	engines, err := v.ListAllKVSecretEngines(o.Namespace)
	if err != nil {
		return nil, err
	}

	return engines, nil
}
