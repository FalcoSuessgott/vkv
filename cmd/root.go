package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/ghodss/yaml"
)

const directorySuffix = "/"

// Secrets map containing all readable secrets
type Secrets map[string]interface{}

// Options holds all available commandline options
type Options struct {
	rootPath string
	subPath  string
	list     bool
	restore  bool // import is a golang reserved key word
	json     bool
	yaml     bool
}

func newRootCmd(version string) *cobra.Command {
	o := &Options{}
	s := Secrets{}

	cmd := &cobra.Command{
		Use:           "vkv",
		Short:         "vault kv engine exporter and importer",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			s.listRecursive(v, o.rootPath, o.subPath)

			if o.json && o.yaml {
				return fmt.Errorf("cannot specify both --to-json and --to-yaml")
			}

			if !o.json && !o.yaml {
				fmt.Println(s.print())
			}

			if o.json {
				fmt.Println(s.toJSON())
			}

			if o.yaml {
				fmt.Println(s.toYAML())
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.rootPath, "root-path", "p", o.rootPath, "root path")
	cmd.Flags().StringVarP(&o.subPath, "sub-path", "s", o.subPath, "sub path")

	cmd.Flags().BoolVarP(&o.json, "to-json", "j", o.json, "print secrets in json format")
	cmd.Flags().BoolVarP(&o.yaml, "to-yaml", "y", o.json, "print secrets in yaml format")

	cmd.Flags().BoolVarP(&o.list, "list", "l", o.list, "list all secrets from a path")
	cmd.Flags().BoolVarP(&o.restore, "import", "i", o.restore, "import all secrets")

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}

func (s *Secrets) listRecursive(v *vault.Vault, rootPath, subPath string) {
	keys, err := v.ListPath(rootPath, subPath)
	if err != nil {
		log.Fatalf("error listing secrets at %s/%s: %v.\n", rootPath, subPath, err)
	}

	for _, k := range keys {
		if strings.HasSuffix(k, directorySuffix) {
			s.listRecursive(v, rootPath, subPath+directorySuffix+strings.TrimSuffix(k, directorySuffix))
		} else {
			secrets, err := v.ReadSecrets(rootPath, subPath+directorySuffix+k)
			if err != nil {
				log.Fatalf("error reading secret at %s/%s/%s: %v.\n", rootPath, subPath, k, err)
			}

			path := rootPath + directorySuffix + subPath + directorySuffix + k
			if subPath == "" {
				path = rootPath + directorySuffix + k
			}

			(*s)[path] = secrets
		}
	}
}

func (s Secrets) toJSON() string {
	out, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("error while marshalling map: %v\n", err)
	}

	return string(out)
}

func (s *Secrets) toYAML() string {
	out, err := yaml.JSONToYAML([]byte(s.toJSON()))
	if err != nil {
		log.Fatalf("error while marshalling from json: %v\n", err)
	}

	return string(out)
}

func (s Secrets) sortKeys() []string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (s Secrets) print() string {
	out := ""
	for _, k := range s.sortKeys() {
		out += fmt.Sprintf("%s:\t%v\n", k, s[k])
	}

	return out
}
