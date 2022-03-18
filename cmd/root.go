package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const (
	defaultKVPath        = "kv"
	maxValueLengthEnvVar = "VKV_MAX_VALUE_LENGTH"
)

var (
	defaultWriter             = os.Stdout
	errInvalidFlagCombination = fmt.Errorf("invalid flag combination specified")
	errMultipleOutputFormats  = fmt.Errorf("specified multiple output formats, only one is allowed")
)

type options struct {
	version bool
	writer  io.Writer
}

func defaultOptions() *options {
	return &options{
		writer: defaultWriter,
	}
}

func rootCmd(version string) *cobra.Command {
	o := defaultOptions()

	cmd := &cobra.Command{
		Use:           "vkv",
		Short:         "utility for interacting with Vaults KV engine",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.version {
				fmt.Fprintf(cmd.OutOrStdout(), "vkv: %s\n", version)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&o.version, "version", "v", o.version, "display version")

	// commands
	cmd.AddCommand(
		mergeCmd(),
		lsCmd(),
		importCmd(),
	)

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := rootCmd(version).Execute(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (o *options) validateFlags() error {
	return nil
}
