package list

import (
	"io"
	"log"
	"strings"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/namespace"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

const envVarListNamespacePrefix = "VKV_LIST_NAMESPACES_"

type listNamespaceOptions struct {
	Namespace string `env:"NS"`

	Regex        string `env:"REGEX"`
	All          bool   `env:"ALL"`
	FormatString string `env:"FORMAT" envDefault:"base"`

	outputFormat printer.OutputFormat
	writer       io.Writer
}

func newListNamespacesCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var namespaces vault.Namespaces

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
				printer.WithWriter(o.writer),
				printer.WithRegex(o.Regex),
			).Out(namespaces)
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVarP(&o.Namespace, "ns", "n", o.Namespace, "specify the namespace (env: VKV_LIST_NAMESPACES_NS)")
	cmd.Flags().StringVarP(&o.Regex, "regex", "r", o.Regex, "filter namespaces by the specified regex pattern (env: VKV_LIST_NAMESPACES_REGEX)")
	cmd.Flags().BoolVarP(&o.All, "all", "a", o.All, "list all namespaces recursively from the specified namespace (env: VKV_LIST_NAMESPACES_ALL)")
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\" (env: VKV_LIST_NAMESPACES_FORMAT")

	o.writer = writer

	return cmd
}

func (o *listNamespaceOptions) validateFlags() error {
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
