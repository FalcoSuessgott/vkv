package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	prt "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
)

var getRootPath = func(m map[string]interface{}) string {
	for k := range m {
		return k
	}

	return ""
}

type importOptions struct {
	EnginePath string `env:"ENGINE_PATH"`
	Path       string `env:"PATH"`

	File string `env:"FILE"`

	Force          bool `env:"FORCE"`
	DryRun         bool `env:"DRY_RUN"`
	Silent         bool `env:"SILENT"`
	ShowValues     bool `env:"SHOW_VALUES"`
	MaxValueLength int  `env:"MAX_VALUE_LENGTH" envDefault:"12"`

	input io.Reader
}

// NewImportCmd import subcommand.
// nolint: cyclop, gocognit, lll
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
			// read existing secrets from the rootPath
			rootPath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

			printer = prt.NewSecretPrinter(
				prt.CustomValueLength(o.MaxValueLength),
				prt.ShowValues(o.ShowValues),
				prt.ToFormat(prt.Base),
				prt.WithVaultClient(vaultClient),
				prt.WithWriter(writer),
				prt.ShowVersion(true),
				prt.ShowMetadata(true),
				prt.WithEnginePath(rootPath),
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
				return o.dryRun(rootPath, subPath, secrets)
			}

			// enable kv engine, error if already enabled, unless force is used
			if err := vaultClient.EnableKV2EngineErrorIfNotForced(o.Force, rootPath); err != nil {
				return err
			}

			// write secrets
			if err := o.writeSecrets(rootPath, subPath, secrets); err != nil {
				return err
			}

			// show result if not silence mode
			if !o.Silent {
				result, err := o.printResult(rootPath, subPath)
				if err != nil {
					return err
				}

				if err := printer.Out(result); err != nil {
					return err
				}
			}

			return nil
		},
	}

	// Input
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KVv2 Engine path (env: VKV_IMPORT_PATH)")
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath, "engine path in case your KV-engine contains special characters such as \"/\", the path value will then be appended if specified (\"<engine-path>/<path>\") (env: VKV_IMPORT_ENGINE_PATH)")
	cmd.Flags().StringVarP(&o.File, "file", "f", o.File, "path to a file containing vkv yaml or json output (env: VKV_IMPORT_FILE)")

	// Options
	cmd.Flags().BoolVar(&o.Force, "force", o.Force, "overwrite existing kv entries (env: VKV_IMPORT_FORCE)")
	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "d", o.DryRun, "print resulting KV engine (env: VKV_IMPORT_DRY_RUN)")
	cmd.Flags().BoolVarP(&o.Silent, "silent", "s", o.Silent, "do not output secrets (env: VKV_IMPORT_SILENT)")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "don't mask values (env: VKV_IMPORT_SHOW_VALUES)")
	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values. Set to \"-1\" for disabling "+
		"(env: VKV_IMPORT_MAX_VALUE_LENGTH)")

	o.input = cmd.InOrStdin()

	return cmd
}

func (o *importOptions) validateFlags(cmd *cobra.Command, args []string) error {
	switch {
	case o.EnginePath == "" && o.Path == "":
		return errors.New("no KV-paths given. Either --engine-path / -e or --path / -p needs to be specified")
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

func (o *importOptions) writeSecrets(rootPath, subPath string, secrets map[string]interface{}) error {
	transformedMap := make(map[string]interface{})
	utils.TransformMap(secrets, transformedMap, "")

	for p, m := range transformedMap {
		secret, ok := m.(map[string]interface{})
		if !ok {
			log.Fatalf("cannot convert %T to map[string]interface", secret)
		}

		// replace original path with the new engine path
		newSubPath := strings.TrimPrefix(p, getRootPath(secrets))

		// unless a subpath has been specified by the user
		if subPath != "" {
			newSubPath = path.Join(subPath, newSubPath)
		}

		if err := vaultClient.WriteSecrets(rootPath, newSubPath, secret); err != nil {
			return fmt.Errorf("error writing secret \"%s\": %w", p, err)
		}

		fmt.Fprintf(writer, "writing secret \"%s\" \n", path.Join(rootPath, newSubPath))
	}

	fmt.Fprintln(writer, "successfully imported all secrets")

	return nil
}

func (o *importOptions) dryRun(rootPath, subPath string, secrets map[string]interface{}) error {
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "preview:")
	fmt.Fprintln(writer, "")

	fmt.Println("root", rootPath, "sub", subPath)

	res, err := vaultClient.ListRecursive(rootPath, "", true)
	if err != nil {
		return fmt.Errorf("error listing secrets from \"%s/%s\": %w", rootPath, subPath, err)
	}

	pp.Println("res", res)

	existingSecrets := make(map[string]interface{})

	if o.EnginePath != "" {
		existingSecrets[rootPath] = utils.ToMapStringInterface(res)
	} else {
		existingSecrets = utils.PathMap(path.Join(rootPath, subPath), utils.ToMapStringInterface(res), false)
	}

	// add new secrets to it
	newSecrets := make(map[string]interface{})

	for _, v := range secrets {
		if o.EnginePath != "" {
			newSecrets[rootPath] = utils.ToMapStringInterface(v)
		} else {
			newSecrets = utils.PathMap(path.Join(rootPath, subPath), utils.ToMapStringInterface(v), false)
		}
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

func (o *importOptions) printResult(rootPath, subPath string) (map[string]interface{}, error) {
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "result:")
	fmt.Fprintln(writer, "")

	var isSecretPath bool

	// read recursive all secrets
	s, err := vaultClient.ListRecursive(rootPath, subPath, false)
	if err != nil {
		return nil, err
	}

	// check if path is a directory or secret path
	if _, isSecret := vaultClient.ReadSecrets(rootPath, subPath); isSecret == nil {
		isSecretPath = true
	}

	path := path.Join(rootPath, subPath)
	if o.EnginePath != "" {
		path = subPath
	}

	// prepare the output map
	pathMap := utils.PathMap(path, utils.ToMapStringInterface(s), isSecretPath)

	if o.EnginePath != "" {
		return map[string]interface{}{
			o.EnginePath: pathMap,
		}, nil
	}

	return pathMap, nil
}
