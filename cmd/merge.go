package cmd

import (
	"fmt"
	"io"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

type mergeOptions struct {
	srcNamespace, destNamespace string
	srcPath, destPath           string
	writer                      io.Writer
}

func defaultMergeptions() *mergeOptions {
	return &mergeOptions{
		writer: defaultWriter,
	}
}

func mergeCmd() *cobra.Command {
	o := defaultMergeptions()

	cmd := &cobra.Command{
		Use:           "merge",
		Short:         "merge the key-value pairs of a source and destination engine and print the result",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			srcC, err := vault.NewClient()
			if err != nil {
				return fmt.Errorf("error while starting source client: %w", err)
			}

			if err := srcC.ListRecursive(utils.SplitPath(o.srcPath)); err != nil {
				return fmt.Errorf("error while reading source path: %w", err)
			}

			destC, err := vault.NewClient()
			if err != nil {
				return fmt.Errorf("error while starting destination client: %w", err)
			}

			if err := destC.ListRecursive(utils.SplitPath(o.destPath)); err != nil {
				return fmt.Errorf("error while reading destination path: %w", err)
			}

			newDestination, _ := utils.SplitPath(o.destPath)
			result := utils.MergeMaps(srcC.Secrets, destC.Secrets, newDestination)

			if err := printer.NewPrinter(result, printer.ShowSecrets(true)).Out(); err != nil {
				return fmt.Errorf("error while printing result: %w", err)
			}

			return nil
		},
	}

	// Input
	cmd.Flags().StringVar(&o.srcNamespace, "src-namespace", o.srcNamespace, "source namespace")
	cmd.Flags().StringVar(&o.destNamespace, "dest-namespace", o.destNamespace, "destination namespace")

	cmd.Flags().StringVar(&o.srcPath, "src-path", o.srcPath, "source path")
	cmd.Flags().StringVar(&o.destPath, "dest-path", o.destPath, "destination path")

	return cmd
}

func (o *mergeOptions) validateFlags() error {
	if o.srcPath == "" || o.destPath == "" {
		return fmt.Errorf("both --src-path and --dest-path needs to be specified")
	}

	if o.srcPath == o.destPath {
		return fmt.Errorf("--src-path and --dest-path cannot be equal")
	}

	return nil
}
