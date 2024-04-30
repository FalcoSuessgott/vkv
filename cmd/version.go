package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd version subcommand.
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "print vkv version",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(writer, "vkv %s\n", Version)

			return nil
		},
	}

	return cmd
}
