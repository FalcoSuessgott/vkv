package cmd

import (
	"fmt"
	"io"
	"log"
	"path"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

// ImportOptions struct holding all import option flags.
type ImportOptions struct {
	Force          bool   `env:"IMPORT_FORCE"`
	DryRun         bool   `env:"IMPORT_DRY_RUN"`
	Silent         bool   `env:"IMPORT_SILENT"`
	ShowValues     bool   `env:"IMPORT_SHOW_VALUES"`
	MaxValueLength int    `env:"IMPORT_MAX_VALUE_LENGTH" envDefault:"12"`
	Path           string `env:"IMPORT_PATH"`
	File           string `env:"IMPORT_FILE"`

	writer io.Writer
}

//nolint: cyclop, gocognit
func newImportCmd() *cobra.Command {
	o := &ImportOptions{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "import -p <path> ",
		Short:         "import secrets from vkv's json or yaml output",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse envs
			if err := o.parseEnvs(); err != nil {
				return err
			}

			// validate flags
			if err := o.validateFlags(args); err != nil {
				return err
			}

			// vault auth
			v, err := vault.NewClient()
			if err != nil {
				return err
			}

			// get user input via -f or STDIN
			input, err := o.getInput(cmd)
			if err != nil {
				return err
			}

			// parse input
			secrets, err := o.parseInput(input)
			if err != nil {
				return err
			}

			// print preview during dryrun and exit
			if o.DryRun {
				return o.dryRun(v, secrets)
			}

			// enable kv engine, error if already enabled, unless force is used
			if err := v.EnableKV2EngineErrorIfNotForced(o.Force, o.Path); err != nil {
				return err
			}

			// write secrets
			if err := o.writeSecrets(secrets, v); err != nil {
				return err
			}

			// show result if not silence mode
			if !o.Silent {
				return o.printResult(v)
			}

			return nil
		},
	}

	// Input
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KVv2 Engine path (env: VKV_IMPORT_PATH)")
	cmd.Flags().StringVarP(&o.File, "file", "f", o.File, "path to a file containing vkv yaml or json output (env: VKV_IMPORT_FILE)")

	// Options
	cmd.Flags().BoolVar(&o.Force, "force", o.Force, "overwrite existing kv entries (env: VKV_IMPORT_FORCE)")
	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "d", o.DryRun, "print resulting KV engine (env: VKV_IMPORT_DRY_RUN)")
	cmd.Flags().BoolVarP(&o.Silent, "silent", "s", o.Silent, "do not output secrets (env: VKV_IMPORT_SILENT)")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "don't mask values (env: VKV_IMPORT_SHOW_VALUES)")
	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values. Set to \"-1\" for disabling "+
		"(env: VKV_IMPORT_MAX_VALUE_LENGTH)")

	o.writer = cmd.OutOrStdout()

	return cmd
}

func (o *ImportOptions) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarPrefix,
	}); err != nil {
		return err
	}

	return nil
}

//nolint: cyclop
func (o *ImportOptions) validateFlags(args []string) error {
	switch {
	case len(args) == 0 && o.Path == "":
		return fmt.Errorf("no KV-path given, -path / -p needs to be specified")
	case o.Force && o.DryRun:
		return fmt.Errorf("%w: %s", errInvalidFlagCombination, "cannot specify both --force and --dry-run")
	case o.Silent && o.DryRun:
		return fmt.Errorf("%w: %s", errInvalidFlagCombination, "cannot specify both --silent and --dry-run")
	case len(args) > 0:
		if o.File != "" && args[0] == "-" {
			return fmt.Errorf("%w: %s", errInvalidFlagCombination, "cannot specify both --file and read from STDIN")
		}
	}

	return nil
}

func (o *ImportOptions) getInput(cmd *cobra.Command) ([]byte, error) {
	if o.File != "" {
		out, err := utils.ReadFile(o.File)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(o.writer, "reading secrets from %s\n", o.File)

		return out, nil
	}

	out, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(o.writer, "reading secrets from STDIN")

	if len(out) == 0 {
		return nil, fmt.Errorf("no input found, perhaps the piped command failed or specified file is empty")
	}

	return out, nil
}

func (o *ImportOptions) parseInput(input []byte) (map[string]interface{}, error) {
	json, err := utils.FromJSON(input)
	if err != nil {
		yaml, err := utils.FromYAML(input)
		if err != nil {
			return nil, fmt.Errorf("cannot parse input, perhaps not a vkv output? Error: %w", err)
		}

		fmt.Fprintln(o.writer, "parsing secrets from YAML")

		return yaml, nil
	}

	fmt.Fprintln(o.writer, "parsing secrets from JSON")

	return json, nil
}

func (o *ImportOptions) writeSecrets(secrets map[string]interface{}, v *vault.Vault) error {
	transformedMap := make(map[string]interface{})
	utils.TransformMap("", secrets, &transformedMap)

	for p, m := range transformedMap {
		secrets, ok := m.(map[string]interface{})
		if !ok {
			log.Fatalf("cannot convert %T to map[string]interface", secrets)
		}

		rootPath, subPath := utils.SplitPath(o.Path)
		_, subPath2 := utils.SplitPath(p)

		nSubPath := path.Join(subPath, subPath2)
		if err := v.WriteSecrets(rootPath, nSubPath, secrets); err != nil {
			return fmt.Errorf("error writing secret \"%s\": %w", p, err)
		}

		fmt.Fprintf(o.writer, "writing secret \"%s\" \n", path.Join(rootPath, nSubPath))
	}

	fmt.Fprintln(o.writer, "successfully imported all secrets")

	return nil
}

func (o *ImportOptions) dryRun(v *vault.Vault, secrets map[string]interface{}) error {
	printer := printer.NewPrinter(
		printer.CustomValueLength(o.MaxValueLength),
		printer.ShowValues(o.ShowValues),
		printer.ToFormat(printer.Base),
		printer.WithVaultClient(v),
		printer.WithWriter(o.writer),
		printer.WithEnginePath(o.Path),
		printer.ShowVersion(true),
		printer.ShowMetadata(true),
	)

	fmt.Fprintln(o.writer, "")
	fmt.Fprintln(o.writer, "preview:")
	fmt.Fprintln(o.writer, "")

	// read existing ecrets from the rootPath
	rootPath, _ := utils.SplitPath(o.Path)
	existingSecrets := make(map[string]interface{})

	tmp, err := v.ListRecursive(rootPath, "")
	if err == nil {
		existingSecrets = utils.PathMap(rootPath, utils.ToMapStringInterface(tmp), false)
	}

	// add new secrets to it
	newSecrets := make(map[string]interface{})
	for _, v := range secrets {
		newSecrets = utils.PathMap(o.Path, utils.ToMapStringInterface(v), false)
	}

	// deep merge both secrets
	mergedSecrets := utils.DeepMergeMaps(newSecrets, existingSecrets)
	if err := printer.Out(mergedSecrets); err != nil {
		return err
	}

	fmt.Fprintln(o.writer, "")
	fmt.Fprintln(o.writer, "apply changes by using the --force flag")

	return nil
}

func (o *ImportOptions) printResult(v *vault.Vault) error {
	printer := printer.NewPrinter(
		printer.CustomValueLength(o.MaxValueLength),
		printer.ShowValues(o.ShowValues),
		printer.ToFormat(printer.Base),
		printer.WithVaultClient(v),
		printer.WithWriter(o.writer),
		printer.ShowVersion(true),
		printer.ShowMetadata(true),
		printer.WithEnginePath(o.Path),
	)

	fmt.Fprintln(o.writer, "")
	fmt.Fprintln(o.writer, "result:")
	fmt.Fprintln(o.writer, "")

	rootPath, _ := utils.SplitPath(o.Path)

	s, err := v.ListRecursive(rootPath, "")
	if err != nil {
		return err
	}

	secrets := utils.PathMap(rootPath, utils.ToMapStringInterface(s), false)

	if err := printer.Out(secrets); err != nil {
		return err
	}

	return nil
}
