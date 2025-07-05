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
	"github.com/spf13/cobra"
)

type importOptions struct {
	EnginePath string `env:"ENGINE_PATH"`
	Path       string `env:"PATH"`

	File string `env:"FILE"`

	Force          bool `env:"FORCE"`
	DryRun         bool `env:"DRY_RUN"`
	Silent         bool `env:"SILENT"`
	ShowValues     bool `env:"SHOW_VALUES"`
	MaxValueLength int  `env:"MAX_VALUE_LENGTH" envDefault:"12"`

	SkipErrors bool `env:"SKIP_ERRORS" envDefault:"false"`

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
		Short:         "import secrets from vkv's export json or yaml output",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE:       o.validateFlags,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// if no path specified, use the path from the secrets to be imported
			if o.EnginePath == "" && o.Path == "" {
				fmt.Fprintln(writer, "no path specified, trying to determine root path from the provided input")

				rootPath, err := utils.GetRootElement(secrets)
				if err != nil {
					return fmt.Errorf("try specifying a destination path using -p/-e. %w", err)
				}

				fmt.Fprintf(writer, "using \"%s\" as KV engine path\n", rootPath)

				// detect whether its an engine path or a normal root path
				if len(strings.Split(rootPath, utils.Delimiter)) > 1 {
					o.EnginePath = rootPath
				} else {
					o.Path = rootPath
				}
			}

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
				prt.WithEnginePath(utils.NormalizePath(rootPath)),
				prt.WithContext(rootContext),
			)

			// print preview during dry run and exit
			if o.DryRun {
				// replace the root path in the secrets with the specified path
				secretsWithNewPath := make(map[string]interface{})
				for _, v := range secrets {
					secretsWithNewPath = utils.UnflattenMap(utils.NormalizePath(path.Join(rootPath, subPath)), utils.ToMapStringInterface(v), o.EnginePath)
				}

				return o.dryRun(rootPath, secretsWithNewPath)
			}

			// enable kv engine, error if already enabled, unless force is used
			if err := vaultClient.EnableKV2EngineErrorIfNotForced(rootContext, o.Force, rootPath); !o.SkipErrors && err != nil {
				return err
			}

			// write secrets
			if err := o.writeSecrets(rootPath, subPath, secrets); err != nil {
				return err
			}

			// show result if not silence mode
			if !o.Silent {
				result, err := o.printResult(rootPath)
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
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KV engine path (env: VKV_IMPORT_PATH)")
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath, "engine path in case your KV-engine contains special characters such as \"/\", the path (-p) flag will then be appended if specified (\"<engine-path>/<path>\") (env: VKV_IMPORT_PATH)")
	cmd.Flags().StringVarP(&o.File, "file", "f", o.File, "path to a file containing vkv export json or yaml output (env: VKV_IMPORT_FILE)")
	cmd.Flags().BoolVar(&o.SkipErrors, "skip-errors", o.SkipErrors, "don't exit on errors (permission denied, ...) (env: VKV_EXPORT_SKIP_ERRORS)")

	// Options
	cmd.Flags().BoolVar(&o.Force, "force", o.Force, "overwrite existing kv secrets (env: VKV_IMPORT_FORCE)")
	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "d", o.DryRun, "print resulting KV secrets (env: VKV_IMPORT_DRY_RUN)")
	cmd.Flags().BoolVarP(&o.Silent, "silent", "s", o.Silent, "do not output secrets (env: VKV_IMPORT_SILENT)")
	cmd.Flags().BoolVar(&o.ShowValues, "show-values", o.ShowValues, "don't mask values (env: VKV_IMPORT_SHOW_VALUES)")
	cmd.Flags().IntVar(&o.MaxValueLength, "max-value-length", o.MaxValueLength, "maximum char length of values. Set to \"-1\" for disabling "+
		"(env: VKV_IMPORT_MAX_VALUE_LENGTH)")

	o.input = cmd.InOrStdin()

	return cmd
}

func (o *importOptions) validateFlags(cmd *cobra.Command, args []string) error {
	switch {
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
	utils.FlattenMap(secrets, transformedMap, "")

	for p, m := range transformedMap {
		secret, ok := m.(map[string]interface{})
		if !ok {
			log.Fatalf("cannot convert %T to map[string]interface", secret)
		}

		// replace original path with the new engine path
		t, _ := utils.GetRootElement(secrets)
		newSubPath := strings.TrimPrefix(p, t)

		// unless a subpath has been specified by the user
		if subPath != "" {
			newSubPath = path.Join(subPath, newSubPath)
		}

		if err := vaultClient.WriteSecrets(rootContext, rootPath, newSubPath, secret); !o.SkipErrors && err != nil {
			return fmt.Errorf("error writing secret \"%s\": %w", p, err)
		}

		fmt.Fprintf(writer, "writing secret \"%s\" \n", path.Join(rootPath, newSubPath))
	}

	fmt.Fprintln(writer, "successfully imported all secrets")

	return nil
}

func (o *importOptions) dryRun(rootPath string, secrets map[string]interface{}) error {
	fmt.Printf("fetching any existing KV secrets from \"%s\" (if any)\n", utils.NormalizePath(rootPath))

	tmp, err := vaultClient.ListRecursive(rootContext, rootPath, "", true)
	if err != nil {
		return fmt.Errorf("error listing secrets from \"%s/\": %w", rootPath, err)
	}

	if len(utils.ToMapStringInterface(tmp)) == 0 {
		fmt.Println("no secrets found - nothing to compare with")
	}

	existingSecrets := utils.UnflattenMap(utils.NormalizePath(rootPath), utils.ToMapStringInterface(tmp), o.EnginePath)

	// check whether new and existing secrets are equal
	if fmt.Sprint(secrets) == fmt.Sprint(existingSecrets) {
		fmt.Fprintln(writer, "")
		fmt.Fprintln(writer, "input matches secrets - no changes needed:")
		fmt.Fprintln(writer, "")

		if err := printer.Out(existingSecrets); err != nil {
			return err
		}

		return nil
	}

	fmt.Fprintf(writer, "deep merging provided secrets with existing secrets read from \"%s\"\n", utils.NormalizePath(rootPath))
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "preview:")
	fmt.Fprintln(writer, "")

	if err := printer.Out(utils.DeepMergeMaps(secrets, existingSecrets)); err != nil {
		return err
	}

	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "apply changes by using the --force flag")

	return nil
}

func (o *importOptions) printResult(rootPath string) (map[string]interface{}, error) {
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "result:")
	fmt.Fprintln(writer, "")

	printer = prt.NewSecretPrinter(
		prt.CustomValueLength(o.MaxValueLength),
		prt.ShowValues(o.ShowValues),
		prt.ToFormat(prt.Base),
		prt.WithVaultClient(vaultClient),
		prt.WithWriter(writer),
		prt.ShowVersion(true),
		prt.ShowMetadata(true),
		prt.ShowVersion(true),
		prt.WithEnginePath(utils.NormalizePath(rootPath)),
		prt.WithContext(rootContext),
	)

	secrets, err := vaultClient.ListRecursive(rootContext, rootPath, "", false)
	if err != nil {
		return nil, err
	}

	return utils.UnflattenMap(utils.NormalizePath(rootPath), utils.ToMapStringInterface(secrets), o.EnginePath), nil
}
