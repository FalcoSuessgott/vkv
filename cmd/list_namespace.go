package cmd

import (
	"log"
	"strings"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/namespace"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

var namespaces vault.Namespaces

type listNamespaceOptions struct {
	Namespace string `env:"NS"`

	Regex        string `env:"REGEX"`
	All          bool   `env:"ALL"`
	FormatString string `env:"FORMAT" envDefault:"base"`

	outputFormat printer.OutputFormat
}

func newListNamespacesCmd() *cobra.Command {
	var err error

	o := &listNamespaceOptions{}

	if err := utils.ParseEnvs(envVarListNamespacePrefix, o); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "namespaces",
		Short:         "list all namespaces",
		Aliases:       []string{"ns"},
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !o.All {
				if namespaces, err = o.listNamespaces(vaultClient); err != nil {
					return err
				}
			} else {
				if namespaces, err = o.listAllNamespaces(vaultClient); err != nil {
					return err
				}
			}

			return printer.NewPrinter(
				printer.ToFormat(o.outputFormat),
				printer.WithWriter(writer),
				printer.WithRegex(o.Regex),
			).Out(namespaces)
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVarP(&o.Namespace, "ns", "n", o.Namespace, "specify the namespace (env: VKV_LIST_NAMESPACES_NS)")
	cmd.Flags().StringVarP(&o.Regex, "regex", "r", o.Regex, "filter namespaces by the specified regex pattern (env: VKV_LIST_NAMESPACES_REGEX)")
	cmd.Flags().BoolVarP(&o.All, "all", "a", o.All, "list all namespaces recursively from the specified namespace (env: VKV_LIST_NAMESPACES_ALL)")
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\" (env: VKV_LIST_NAMESPACES_FORMAT")

	return cmd
}

func (o *listNamespaceOptions) Validate() error {
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

func (o *listNamespaceOptions) listNamespaces(v *vault.Vault) (vault.Namespaces, error) {
	ns, err := v.ListNamespaces(o.Namespace)
	if err != nil {
		return nil, err
	}

	m := make(vault.Namespaces)
	m[o.Namespace] = ns

	return m, nil
}

func (o *listNamespaceOptions) listAllNamespaces(v *vault.Vault) (vault.Namespaces, error) {
	ns, err := v.ListAllNamespaces(o.Namespace)
	if err != nil {
		return nil, err
	}

	return ns, nil
}
