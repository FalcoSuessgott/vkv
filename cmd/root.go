package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

const (
	DELIMITER       = "/"
	SEPERATOR       = " "
	NEW_LINE        = "\n"
	MASK_CHAR       = "*"
	DEFAULT_KV_PATH = "kv2"
)

var defaultWriter = os.Stdout

// Secrets map containing all readable secrets.
type (
	Secrets map[string]interface{}
	Keys    []string
)

// Options holds all available commandline options.
type Options struct {
	rootPath    string
	subPath     string
	writer      io.Writer
	onlyKeys    bool
	onlyPaths   bool
	showSecrets bool
	json        bool
	yaml        bool
}

func defaultOptions() *Options {
	return &Options{
		rootPath:    DEFAULT_KV_PATH,
		showSecrets: false,
		writer:      defaultWriter,
	}
}

func newRootCmd(version string) *cobra.Command {
	o := defaultOptions()
	s := Secrets{}

	cmd := &cobra.Command{
		Use:           "vkv",
		Short:         "recursively list secrets from Vaults KV2 engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(); err != nil {
				return err
			}

			// READ
			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			s.listRecursive(v, o.rootPath, o.subPath)

			// MODIFY
			o.evalModifyFlags(s)

			// PRINT
			if !o.json && !o.yaml {
				o.print(s)
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

	// Input
	cmd.Flags().StringVarP(&o.rootPath, "root-path", "p", o.rootPath, "root path")
	cmd.Flags().StringVarP(&o.subPath, "sub-path", "s", o.subPath, "sub path")

	// Modify
	cmd.Flags().BoolVar(&o.onlyKeys, "only-keys", o.onlyKeys, "print only keys")
	cmd.Flags().BoolVar(&o.onlyPaths, "only-paths", o.onlyPaths, "print only paths")
	cmd.Flags().BoolVar(&o.showSecrets, "show-secrets", o.showSecrets, "print out secrets")

	// Output format
	cmd.Flags().BoolVarP(&o.json, "to-json", "j", o.json, "print secrets in json format")
	cmd.Flags().BoolVarP(&o.yaml, "to-yaml", "y", o.json, "print secrets in yaml format")

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}

func (o *Options) validateFlags() error {
	if o.json && o.yaml {
		return fmt.Errorf("cannot specify both --to-json and --to-yaml\n")
	}

	if o.onlyKeys && o.showSecrets {
		return fmt.Errorf("cannot specify both --only-keys and --show-secrets\n")
	}

	if o.onlyPaths && o.showSecrets {
		return fmt.Errorf("cannot specify both --only-paths and --show-secrets\n")
	}

	if o.onlyKeys && o.onlyPaths {
		return fmt.Errorf("cannot specify both --only-keys and --only-paths\n")
	}
	return nil
}

func (o *Options) evalModifyFlags(s Secrets) {
	if o.onlyKeys {
		s.onlyKeys()
	}

	if o.onlyPaths {
		s.onlyPaths()
	}

	if !o.showSecrets {
		s.maskSecrets()
	}
}

func (s *Secrets) listRecursive(v *vault.Vault, rootPath, subPath string) {
	keys, err := v.ListPath(rootPath, subPath)
	if err != nil {
		log.Fatalf("error listing secrets at %s/%s: %v.\n", rootPath, subPath, err)
	}

	for _, k := range keys {
		if strings.HasSuffix(k, DELIMITER) {
			s.listRecursive(v, rootPath, buildPath(subPath, k))
		} else {
			secrets, err := v.ReadSecrets(rootPath, buildPath(subPath, k))
			if err != nil {
				log.Fatalf("error reading secret at %s/%s/%s: %v.\n", rootPath, subPath, k, err)
			}

			(*s)[buildPath(rootPath, subPath, k)] = secrets
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

func sortMapKeys(m map[string]interface{}) []string {
	keys := make(Keys, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	return keys
}

func (k Keys) Len() int {
	return len(k)
}

func (k Keys) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (k Keys) Less(i, j int) bool {
	k1 := strings.Replace(k[i], "/", "\x00", -1)
	k2 := strings.Replace(k[j], "/", "\x00", -1)

	return k1 < k2
}

func (o *Options) print(s Secrets) {
	w := tabwriter.NewWriter(o.writer, 0, 8, 1, '\t', tabwriter.AlignRight)

	fmt.Fprintf(w, "%s%s%s", o.rootPath, DELIMITER, NEW_LINE)

	for _, k := range sortMapKeys(s) {
		if o.onlyKeys {
			fmt.Fprintf(w, "%s\t%v\n", k, printMap(s[k]))
		} else if o.onlyPaths {
			fmt.Fprintln(w, k)
		} else {
			fmt.Fprintf(w, "%s\t%v\n", k, printMap(s[k]))
		}
	}

	w.Flush()
}

func (s Secrets) onlyKeys() {
	for k := range s {
		m, ok := s[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			m[k] = ""
		}
	}
}

func (s Secrets) onlyPaths() {
	for k := range s {
		s[k] = nil
	}
}

func (s Secrets) maskSecrets() {
	for k := range s {
		m, ok := s[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			secret := fmt.Sprintf("%v", m[k])
			m[k] = strings.Repeat(MASK_CHAR, len(secret))
		}
	}
}

func buildPath(elements ...string) string {
	p := ""

	for i, e := range elements {
		e = strings.TrimSuffix(e, DELIMITER)

		if e == "" {
			continue
		}

		p += e

		if i < len(elements) {
			p += DELIMITER
		}
	}

	return strings.TrimSuffix(p, DELIMITER)
}

func printMap(m interface{}) string {
	out := ""

	secrets, ok := m.(map[string]interface{})
	if !ok {
		return ""
	}

	for _, k := range sortMapKeys(secrets) {
		out += k

		if secrets[k] == "" {
			out += SEPERATOR
		} else {
			out += fmt.Sprintf("=%v ", secrets[k])
		}

	}

	return strings.TrimSuffix(out, SEPERATOR)
}
