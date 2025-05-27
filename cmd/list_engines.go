package cmd

import (
	"log"
	"strings"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer/engine"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

var engines vault.Engines

type listEnginesOptions struct {
	Namespace string `env:"NS"`
	Prefix    bool   `env:"NS_PREFIX"`

	Regex string `env:"REGEX"`
	All   bool   `env:"ALL"`

	FormatString string `env:"FORMAT" envDefault:"base"`

	outputFormat prt.OutputFormat
}

func newListEngineCmd() *cobra.Command {
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
		PreRunE:       o.Validate,
		RunE: func(cmd *cobra.Command, args []string) error {
			printer = prt.NewEnginePrinter(
				prt.ToFormat(o.outputFormat),
				prt.WithWriter(writer),
				prt.WithRegex(o.Regex),
				prt.WithNSPrefix(o.Prefix),
			)
			if !o.All {
				if engines, err = o.listEngines(); err != nil {
					return err
				}
			} else {
				if engines, err = o.listAllEngines(); err != nil {
					return err
				}
			}

			return printer.Out(engines)
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "specify the namespace (env: VKV_LIST_ENGINES_NS)")
	cmd.Flags().BoolVarP(&o.Prefix, "include-ns-prefix", "p", o.Prefix, "prepend the namespaces (env: VKV_LIST_ENGINES_NS_PREFIX)")

	cmd.Flags().StringVarP(&o.Regex, "regex", "r", o.Regex, "filter engines by the specified regex pattern (env: VKV_LIST_ENGINES_REGEX")
	cmd.Flags().BoolVarP(&o.All, "all", "a", o.All, "list all KV engines recursively from the specified namespaces (env: VKV_LIST_ENGINES_ALL)")
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\" (env: VKV_LIST_ENGINES_FORMAT)")

	return cmd
}

func (o *listEnginesOptions) Validate(cmd *cobra.Command, args []string) error {
	switch strings.ToLower(o.FormatString) {
	case "yaml", "yml":
		o.outputFormat = prt.YAML
	case "json":
		o.outputFormat = prt.JSON
	case "base":
		o.outputFormat = prt.Base
	default:
		return prt.ErrInvalidFormat
	}

	return nil
}

func (o *listEnginesOptions) listEngines() (vault.Engines, error) {
	engines, err := vaultClient.ListKVSecretEngines(rootContext, o.Namespace)
	if err != nil {
		return nil, err
	}

	m := make(vault.Engines)
	m[o.Namespace] = engines

	return m, nil
}

func (o *listEnginesOptions) listAllEngines() (vault.Engines, error) {
	engines, err := vaultClient.ListAllKVSecretEngines(rootContext, o.Namespace)
	if err != nil {
		return nil, err
	}

	return engines, nil
}
