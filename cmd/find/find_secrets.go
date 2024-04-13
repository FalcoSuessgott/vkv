package find

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/engine"
	"github.com/FalcoSuessgott/vkv/pkg/regex"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

// TODO: flag for no header, flag for including ns prefix, merge list and find in 1 subcommand, flag for printing the secret URL, print match kind (ns, engine, secret name, secret value)

const (
	envVarExportPrefix = "VKV_FIND_"

	Namespace match = iota
	Engine
	SecretName
	SecretValues
)

type match int

// findOptions holds all available commandline options.
type findOptions struct {
	Pattern      string `env:"PATTERN"`
	FormatString string `env:"FORMAT" envDefault:"base"`

	NoHeader    bool `env:"NO_HEADER"`
	NoMatchKind bool `env:"NO_MATCH_KIND"`
	PrintURL    bool `env:"PRINT_URL"`

	outputFormat printer.OutputFormat
}

// NewFindSecretsCmd find subcommand.
func NewFindSecretsCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
	var err error

	o := &findOptions{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "secrets",
		Short:         "search and print out KV secret engines that contain the specified regex pattern",
		Aliases:       []string{"s", "secret"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(); err != nil {
				return err
			}

			// vault auth
			if vaultClient == nil {
				if vaultClient, err = vault.NewDefaultClient(); err != nil {
					return err
				}
			}

			// find all visible engines in all NS
			engines, err := vaultClient.ListAllKVSecretEngines("")
			if err != nil {
				return err
			}

			var countedNS, countedEngines, countedSecrets int

			countedNS = len(engines)

			for _, v := range engines {
				countedEngines += len(v)
			}

			// TODO: count secrets
			fmt.Printf("searching for pattern \"%s\" in all visible namespaces (%d), KV engines (%d) & secrets (%d):\n", o.Pattern, countedNS, countedEngines, countedSecrets)

			// list secrets recursively for all found engines
			for ns, engine := range engines {
				for _, e := range engine {
					secrets, err := vaultClient.ListRecursive(e, "", true)
					if err != nil {
						return err
					}

					for k, v := range utils.ToMapStringInterface(secrets) {
						var found bool
						var kind match
						var url string

						// search for pattern in current NS
						if res, err := regex.MatchRegex(o.Pattern, ns); err == nil && res {
							kind = Namespace
							found = true
						}

						// search for pattern in current engine
						if res, err := regex.MatchRegex(o.Pattern, e); err == nil && res {
							kind = Engine
							found = true
							url = fmt.Sprintf("%s/ui/vault/secrets/%skv/list", os.Getenv("VAULT_ADDR"), e)

						}

						// search for pattern in current secret name
						if res, err := regex.MatchRegex(o.Pattern, k); err == nil && res {
							kind = SecretName
							found = true
							url = fmt.Sprintf("%s/ui/vault/secrets/%skv/list/%s", os.Getenv("VAULT_ADDR"), e, k)

						}

						// search for pattern in current secret values
						res := make(map[string]interface{})

						utils.TransformMap(e, utils.ToMapStringInterface(v), &res)

						for _, secrets := range res {
							for sk, sv := range utils.ToMapStringInterface(secrets) {
								if res, err := regex.MatchRegex(o.Pattern, sk); err == nil && res {
									kind = SecretValues
									found = true
									url = fmt.Sprintf("%s/ui/vault/secrets/%skv/%s", os.Getenv("VAULT_ADDR"), e, k)

								}

								if res, err := regex.MatchRegex(o.Pattern, sv.(string)); err == nil && res {
									kind = SecretValues
									found = true
								}
							}
						}

						// implement tab writer
						if found {
							// to a NS
							// to a secret: <VAULT_ADDR>/ui/vault/secrets/<ENGINE>/kv/list/<SECRET_NAME>
							// within secret: <VAULT_ADDR>/ui/vault/secrets/<ENGINE>/kv/<SECRET_NAME>

							fmt.Println(kind.String(), path.Join(ns, e, k), url)
						}
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVarP(&o.Pattern, "pattern", "p", o.Pattern, "search pattern to find secrets that match the pattern in its keys/values (env: VKV_FIND_PATTERN)")
	cmd.Flags().BoolVar(&o.NoHeader, "no-header", o.NoHeader, "do not print search header (env: VKV_FIND_NO_HEADER)")
	cmd.Flags().BoolVar(&o.PrintURL, "print-url", o.PrintURL, "print the url to the corresponding namespace, secret or engine (env: VKV_FIND_PRINT_URL)")
	cmd.Flags().BoolVar(&o.NoMatchKind, "no-match-kind", o.PrintURL, "do not print the kind of the match (env: VKV_FIND_NO_MATCH_KIND)")

	cmd.Flags().StringVarP(&o.FormatString, "format", "f", o.FormatString, "available output formats: \"base\", \"json\", \"yaml\" (env: VKV_LIST_ENGINES_FORMAT)")

	return cmd
}

func (o *findOptions) validateFlags() error {
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

func (o *findOptions) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarExportPrefix,
	}); err != nil {
		return err
	}

	return nil
}

func (m match) String() string {
	switch m {
	case Namespace:
		return "[Match in Namespace name] "
	case Engine:
		return "[Match in KV Engine Name] "
	case SecretName:
		return "[Match in Secret Name] "
	case SecretValues:
		return "[Match in Secret Values] "
	default:
		return ""
	}
}
