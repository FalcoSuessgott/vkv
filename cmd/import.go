package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	prt "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/spf13/cobra"
)

type importOptions struct {
	Force          bool   `env:"FORCE"`
	DryRun         bool   `env:"DRY_RUN"`
	Silent         bool   `env:"SILENT"`
	ShowValues     bool   `env:"SHOW_VALUES"`
	MaxValueLength int    `env:"MAX_VALUE_LENGTH" envDefault:"12"`
	Path           string `env:"PATH"`
	File           string `env:"FILE"`

	input io.Reader
}

// NewImportCmd import subcommand.
// nolint: cyclop, gocognit
func NewImportCmd() *cobra.Command {
	o := &importOptions{}

	if err := utils.ParseEnvs(envVarImportPrefix, o); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "import",
		Short:         "import secrets from vkv's json or yaml output",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE:       o.validateFlags,
		RunE: func(cmd *cobra.Command, args []string) error {
			printer = prt.NewSecretPrinter(
				prt.CustomValueLength(o.MaxValueLength),
				prt.ShowValues(o.ShowValues),
				prt.ToFormat(prt.Base),
				prt.WithVaultClient(vaultClient),
				prt.WithWriter(writer),
				prt.ShowVersion(true),
				prt.ShowMetadata(true),
				prt.WithEnginePath(o.Path),
			)

			// get user input via -f or STDIN
			input, err := o.getInput()
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
				return o.dryRun(secrets)
			}

			// enable kv engine, error if already enabled, unless force is used
			if err := vaultClient.EnableKV2EngineErrorIfNotForced(o.Force, o.Path); err != nil {
				return err
			}

			// write secrets
			if err := o.writeSecrets(secrets); err != nil {
				return err
			}

			// show result if not silence mode
			if !o.Silent {
				return o.printResult()
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

	// flag constraints
	if err := cmd.MarkFlagRequired("path"); err != nil {
		log.Fatal(err)
	}

	o.input = cmd.InOrStdin()

	return cmd
}

func (o *importOptions) validateFlags(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0 && o.Path == "":
		return errors.New("no KV-path given, -path / -p needs to be specified")
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

func (o *importOptions) getInput() ([]byte, error) {
	if o.File != "" {
		out, err := fs.ReadFile(o.File)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(writer, "reading secrets from %s\n", o.File)

		return out, nil
	}

	out, err := io.ReadAll(o.input)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(writer, "reading secrets from STDIN")

	if len(out) == 0 {
		return nil, errors.New("no input found, perhaps the piped command failed or specified file is empty")
	}

	return out, nil
}

func (o *importOptions) parseInput(input []byte) (map[string]interface{}, error) {
	json, err := utils.FromJSON(input)
	if err != nil {
		yaml, err := utils.FromYAML(input)
		if err != nil {
			return nil, fmt.Errorf("cannot parse input, perhaps not a vkv output? Error: %w", err)
		}

		fmt.Fprintln(writer, "parsing secrets from YAML")

		return yaml, nil
	}

	fmt.Fprintln(writer, "parsing secrets from JSON")

	return json, nil
}

func (o *importOptions) writeSecrets(secrets map[string]interface{}) error {
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
		if err := vaultClient.WriteSecrets(rootPath, nSubPath, secrets); err != nil {
			return fmt.Errorf("error writing secret \"%s\": %w", p, err)
		}

		fmt.Fprintf(writer, "writing secret \"%s\" \n", path.Join(rootPath, nSubPath))
	}

	fmt.Fprintln(writer, "successfully imported all secrets")

	return nil
}

func (o *importOptions) dryRun(secrets map[string]interface{}) error {
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "preview:")
	fmt.Fprintln(writer, "")

	// read existing ecrets from the rootPath
	rootPath, _ := utils.SplitPath(o.Path)
	existingSecrets := make(map[string]interface{})

	tmp, err := vaultClient.ListRecursive(rootPath, "", false)
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

	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "apply changes by using the --force flag")

	return nil
}

func (o *importOptions) printResult() error {
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "result:")
	fmt.Fprintln(writer, "")

	rootPath, _ := utils.SplitPath(o.Path)

	s, err := vaultClient.ListRecursive(rootPath, "", false)
	if err != nil {
		return err
	}

	secrets := utils.PathMap(rootPath, utils.ToMapStringInterface(s), false)

	if err := printer.Out(secrets); err != nil {
		return err
	}

	return nil
}
