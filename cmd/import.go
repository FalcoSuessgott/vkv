package cmd

import (
	"fmt"
	"io"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/ui"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

type importOptions struct {
	namespace string
	path      string
	file      string
	force     bool
	writer    io.Writer
}

func defaultImportOptions() *importOptions {
	return &importOptions{
		writer: defaultWriter,
	}
}

//nolint: gocognit, cyclop
func importCmd() *cobra.Command {
	o := defaultImportOptions()

	cmd := &cobra.Command{
		Use:           "import",
		Short:         "import key-values from file or STDIN to a KV engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(args); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vc, err := vault.NewClient()
			if err != nil {
				return fmt.Errorf("error while starting client: %w", err)
			}

			destC, err := vault.NewClient()
			if err != nil {
				return fmt.Errorf("error while starting destination client: %w", err)
			}

			// read file or stdin
			content, err := utils.ReadFile(o.file)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", o.file, err)
			}

			// determine input
			data, _, err := utils.JSONorYAML(content)
			if err != nil {
				return fmt.Errorf("error while determining %s format: %w", o.file, err)
			}

			// check if path is already enabled, if not enable
			rootPath, _ := utils.SplitPath(o.path)

			if err := vc.EnableKV2Engine(rootPath); err != nil {
				// this fails if no secrets are written on the dest path
				if err := destC.ListRecursive(utils.SplitPath(o.path)); err != nil {
					return fmt.Errorf("error while reading destination path: %w", err)
				}

				if !o.force {
					return fmt.Errorf("error path %s is already enabled and might contain key-values. Use --force for enforcing the import: %w", o.path, err)
				}
			}

			srcMap, ok := data.(map[string]interface{})
			if !ok {
				return fmt.Errorf("cannot parse data to map")
			}

			// merge src and dest and print result
			result := utils.MergeMaps(srcMap, destC.Secrets, o.path)

			// print result
			ui.FInfoMsg(o.writer, "Merging the content of \"%s\" to kv engine \"%s\" the result will look like:", o.file, o.path)
			if err := printer.NewPrinter(result, printer.ShowSecrets(true)).Out(); err != nil {
				return fmt.Errorf("error while printing result: %w", err)
			}

			// prompt for confirmation
			if ui.PromptYN("Are you sure?") {
				for k, v := range result {
					// check here if the k already exist on the existing path
					rootPath, subPath := utils.SplitPath(k)
					if err := vc.WriteSecrets(rootPath, subPath, v.(map[string]interface{})); err != nil {
						return fmt.Errorf("error while importing %s: %w", k, io.ErrClosedPipe)
					}

					ui.FSuccessMsg(o.writer, "successfully uploaded %s", k)
				}
			} else {
				return fmt.Errorf("aborted")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", o.namespace, "namespace of the kv engine")
	cmd.Flags().StringVarP(&o.path, "path", "p", o.path, "path to kv engine to import data")
	cmd.Flags().StringVarP(&o.file, "file", "f", o.file, "path to file to import")
	cmd.Flags().BoolVar(&o.force, "force", o.force, "dont ask for approving")

	return cmd
}

func (o *importOptions) validateFlags(args []string) error {
	if len(args) > 0 {
		if o.file != "" && args[0] == "-" {
			return fmt.Errorf("cannot specify both --path and reading from STDIN")
		}
	}

	if o.file == "" && len(args) == 0 {
		return fmt.Errorf("input either via --path or via STDIN (-) is required")
	}

	if o.file != "" && !utils.FileExists(o.file) {
		return fmt.Errorf("file %s does not exist", o.file)
	}

	if o.path == "" {
		return fmt.Errorf("--path is required")
	}

	return nil
}
