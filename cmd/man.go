package cmd

import (
	"fmt"
	"os"

	mcoral "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

// NewManCmd manpage cmd.
func NewManCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates GoReleaser's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcoral.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))

			return err
		},
	}

	return cmd
}
