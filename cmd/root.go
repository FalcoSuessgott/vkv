package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/FalcoSuessgott/vkv/cmd/export"
	imp "github.com/FalcoSuessgott/vkv/cmd/imp"
	"github.com/FalcoSuessgott/vkv/cmd/list"
	"github.com/FalcoSuessgott/vkv/cmd/snapshot"
	"github.com/FalcoSuessgott/vkv/cmd/version"
	"github.com/spf13/cobra"
)

// NewRootCmd vkv root command.
func NewRootCmd(v string, writer io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vkv",
		Short:         "the swiss army knife when working with Vault KVv2 engines",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// sub commands
	cmd.AddCommand(
		export.NewExportCmd(writer, nil),
		list.NewListCmd(writer, nil),
		snapshot.NewSnapshotCmd(writer, nil),
		version.NewVersionCmd(v),
		imp.NewImportCmd(writer, nil),
	)

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := NewRootCmd(version, os.Stdout).Execute(); err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}
