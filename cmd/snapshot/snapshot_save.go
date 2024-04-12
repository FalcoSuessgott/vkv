package snapshot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	printer "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

const snapshotSavePrefix = "VKV_SNAPSHOT_SAVE_"

type snapshotSaveOptions struct {
	Namespace   string `env:"NS"`
	Destination string `env:"DESTINATION" envDefault:"./vkv-snapshot-export"`
	SkipErrors  bool   `env:"SKIP_ERRORS" envDefault:"false"`

	writer io.Writer
}

func newSnapshotSaveCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
	var err error

	o := &snapshotSaveOptions{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "save",
		Short:         "create a snapshot of all visible KV engines recursively for all namespaces",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if vaultClient == nil {
				if vaultClient, err = vault.NewDefaultClient(); err != nil {
					return err
				}
			}

			engines, err := vaultClient.ListAllKVSecretEngines(o.Namespace)
			if err != nil {
				return err
			}

			return o.backupKVEngines(vaultClient, engines)
		},
	}

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "namespaces from which to save recursively all visible KV engines (env: VKV_SNAPSHOT_SAVE_NS)")
	cmd.Flags().StringVarP(&o.Destination, "destination", "d", o.Destination, "vkv snapshot destination path (env: VKV_SNAPSHOT_SAVE_DESTINATION)")
	cmd.Flags().BoolVar(&o.SkipErrors, "skip-errors", o.SkipErrors, "dont exit on errors (permission denied, deleted secrets) (env: VKV_SNAPSHOT_SAVE_SKIP_ERRORS)")

	o.writer = writer

	return cmd
}

// nolint: cyclop
func (o *snapshotSaveOptions) backupKVEngines(v *vault.Vault, engines map[string][]string) error {
	for _, ns := range utils.SortMapKeys(utils.ToMapStringInterface(engines)) {
		nsDir := path.Join(o.Destination, ns)

		if err := fs.CreateDirectory(nsDir); err != nil {
			return err
		}

		fmt.Fprintf(o.writer, "created %s\n", nsDir)

		for _, e := range engines[ns] {
			enginePath := path.Join(ns, e)

			out, err := v.ListRecursive(enginePath, "", o.SkipErrors)
			if err != nil {
				return err
			}

			b := bytes.NewBufferString("")

			p := printer.NewPrinter(
				printer.CustomValueLength(-1),
				printer.ShowValues(true),
				printer.ToFormat(printer.JSON),
				printer.WithVaultClient(v),
				printer.WithWriter(b),
				printer.ShowVersion(false),
				printer.ShowMetadata(false),
			)

			if err := p.Out(strings.TrimSuffix(e, utils.Delimiter), utils.ToMapStringInterface(out)); err != nil {
				return err
			}

			c, err := io.ReadAll(b)
			if err != nil {
				return err
			}

			engineFile := path.Join(o.Destination, ns, e) + ".yaml"

			if err := os.WriteFile(engineFile, c, 0o600); err != nil {
				return err
			}

			fmt.Fprintf(o.writer, "created %s\n", engineFile)

			b.Reset()
		}
	}

	return nil
}

func (o *snapshotSaveOptions) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: snapshotSavePrefix,
	}); err != nil {
		return err
	}

	return nil
}
