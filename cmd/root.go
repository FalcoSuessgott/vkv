package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

const (
	envVarVKVMode               = "VKV_MODE"
	envVarExportPrefix          = "VKV_EXPORT_"
	envVarImportPrefix          = "VKV_IMPORT_"
	envVarServerPrefix          = "VKV_SERVER_"
	envVarListEnginesPrefix     = "VKV_LIST_ENGINES_"
	envVarListNamespacePrefix   = "VKV_LIST_NAMESPACES_"
	envVarSnapshotRestorePrefix = "VKV_SNAPSHOT_RESTORE_"
	envVarSnapshotSavePrefix    = "VKV_SNAPSHOT_SAVE_"
)

var (
	Version string

	errInvalidFlagCombination = errors.New("invalid flag combination specified")
	vaultClient               *vault.Vault
	writer                    io.Writer
	rootContext               context.Context
	printer                   prt.Printer
)

// NewRootCmd vkv root command.
//
//nolint:cyclop
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vkv",
		Short:         "the swiss army knife when working with Vault KV engines",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// skip vault client creation for completion, version, help, docs and manpage generation
			if (cmd.HasParent() && cmd.Parent().Use == "completion") || cmd.Use == "docs" || cmd.Use == "help" || cmd.Use == "version" || cmd.Use == "man" {
				return nil
			}

			// required to inject the vault client for unit tests
			if vaultClient != nil {
				return nil
			}

			// otherwise create a new vault client
			vc, err := vault.NewDefaultClient(rootContext)
			if err != nil {
				return err
			}

			vaultClient = vc

			go func() {
				vaultClient.LeaseRefresher(rootContext)
			}()

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, ok := os.LookupEnv(envVarVKVMode)
			if !ok {
				return cmd.Help()
			}

			switch strings.ToUpper(mode) {
			case "EXPORT":
				return NewExportCmd().Execute()
			case "IMPORT":
				return NewImportCmd().Execute()
			case "SERVER":
				return NewServerCmd().Execute()
			case "LIST":
				return NewListCmd().Execute()
			case "SNAPSHOT_RESTORE":
				return NewSnapshotRestoreCmd().Execute()
			case "SNAPSHOT_SAVE":
				return NewSnapshotSaveCmd().Execute()
			default:
				return errors.New("invalid value for VKV_MODE")
			}
		},
	}

	// sub commands
	cmd.AddCommand(
		NewExportCmd(),
		NewListCmd(),
		NewSnapshotCmd(),
		NewVersionCmd(),
		NewImportCmd(),
		NewServerCmd(),
		NewManCmd(),
		NewDocCmd(),
	)

	cobra.OnInitialize(
		func() {
			// initialize writer if not already, used for injecting the writer in unit tests
			if writer == nil {
				writer = cmd.OutOrStdout()
			}
		},
	)

	return cmd
}

// Execute invokes the command.
func Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rootContext = ctx

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done

		log.Println("Received shutdown signal")
		cancel()
	}()

	if err := NewRootCmd().ExecuteContext(rootContext); err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}
