package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vkv/cmd/export"
	"github.com/FalcoSuessgott/vkv/cmd/find"
	imp "github.com/FalcoSuessgott/vkv/cmd/imp"
	"github.com/FalcoSuessgott/vkv/cmd/manpage"
	"github.com/FalcoSuessgott/vkv/cmd/server"
	"github.com/FalcoSuessgott/vkv/cmd/snapshot"
	"github.com/FalcoSuessgott/vkv/cmd/version"
	"github.com/spf13/cobra"
)

// NewRootCmd vkv root command.
//
//nolint:cyclop
func NewRootCmd(v string, writer io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vkv",
		Short:         "the swiss army knife when working with Vault KVv2 engines",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, ok := os.LookupEnv("VKV_MODE")
			if !ok {
				return cmd.Help()
			}

			switch strings.ToUpper(mode) {
			case "EXPORT":
				return export.NewExportCmd(writer, nil).Execute()
			case "IMPORT":
				return imp.NewImportCmd(writer, nil).Execute()
			case "FIND":
				return find.NewFindCmd(writer, nil).Execute()
			case "SERVER":
				return server.NewServerCmd(writer, nil).Execute()
			case "SNAPSHOT_RESTORE":
				cmd := snapshot.NewSnapshotCmd(writer, nil)

				for _, c := range cmd.Commands() {
					if c.Name() == "restore" {
						return c.Execute()
					}
				}
			case "SNAPSHOT_SAVE":
				cmd := snapshot.NewSnapshotCmd(writer, nil)

				for _, c := range cmd.Commands() {
					if c.Name() == "save" {
						return c.Execute()
					}
				}
			default:
				return errors.New("invalid value for VKV_MODE")
			}

			return cmd.Help()
		},
	}

	// sub commands
	cmd.AddCommand(
		export.NewExportCmd(writer, nil),
		find.NewFindCmd(writer, nil),
		snapshot.NewSnapshotCmd(writer, nil),
		version.NewVersionCmd(v),
		imp.NewImportCmd(writer, nil),
		server.NewServerCmd(writer, nil),
		manpage.NewManCmd().Cmd,
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
