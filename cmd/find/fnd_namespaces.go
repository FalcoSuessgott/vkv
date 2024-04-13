package find

import (
	"io"
	"log"
	"strings"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/namespace"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

const envVarfindNamespacePrefix = "VKV_FIND_NAMESPACES_"

type findNamespaceOptions struct {
	Namespace string `env:"NS"`

	Regex        string `env:"REGEX"`
	All          bool   `env:"ALL"`
	FormatString string `env:"FORMAT" envDefault:"base"`

	outputFormat printer.OutputFormat
	writer       io.Writer
}

func newFindNamespacesCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
	var err error

	o := &findNamespaceOptions{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "namespaces",
		Short:         "find all namespaces",
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
				if namespaces, err = o.findNamespaces(vaultClient); err != nil {
					return err
				}
			} else {
				if namespaces, err = o.findAllNamespaces(vaultClient); err != nil {
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

	cmd.Flags().StringVarP(&o.Namespace, "ns", "n", o.Namespace, "specify the namespace (env: VKV_find_NAMESPACES_NS)")
	cmd.Flags().StringVarP(&o.Regex, "regex", "r", o.Regex, "filter namespaces by the specified regex pattern (env: VKV_find_NAMESPACES_REGEX)")
	cmd.Flags().BoolVarP(&o.All, "all", "a", o.All, "find all namespaces recursively from the specified namespace (env: VKV_find_NAMESPACES_ALL)")
	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\" (env: VKV_find_NAMESPACES_FORMAT")

	o.writer = writer

	return cmd
}

func (o *findNamespaceOptions) validateFlags() error {
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

func (o *findNamespaceOptions) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarfindNamespacePrefix,
	}); err != nil {
		return err
	}

	return nil
}

func (o *findNamespaceOptions) findNamespaces(v *vault.Vault) (vault.Namespaces, error) {
	ns, err := v.ListNamespaces(o.Namespace)
	if err != nil {
		return nil, err
	}

	m := make(vault.Namespaces)
	m[o.Namespace] = ns

	return m, nil
}

func (o *findNamespaceOptions) findAllNamespaces(v *vault.Vault) (vault.Namespaces, error) {
	ns, err := v.ListAllNamespaces(o.Namespace)
	if err != nil {
		return nil, err
	}

	return ns, nil
}
