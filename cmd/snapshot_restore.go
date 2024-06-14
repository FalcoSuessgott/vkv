package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/spf13/cobra"
)

type snapshotRestoreOptions struct {
	Source string `env:"SOURCE" envDefault:"./vkv-snapshot-export"`
}

func NewSnapshotRestoreCmd() *cobra.Command {
	o := &snapshotRestoreOptions{}

	if err := utils.ParseEnvs(envVarSnapshotRestorePrefix, o); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "restore",
		Short:         "restore the KV engines defined in the specified snapshot",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// make sure source is absolute
			absolutePath, err := filepath.Abs(o.Source)
			if err != nil {
				return err
			}

			return o.restoreSecrets(absolutePath)
		},
	}

	cmd.Flags().StringVarP(&o.Source, "source", "s", o.Source, "source of a vkv snapshot export (env :VKV_SNAPSHOT_RESTORE_SOURCE)")

	return cmd
}

// nolint: cyclop
func (o *snapshotRestoreOptions) restoreSecrets(source string) error {
	return filepath.Walk(source, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		absPath := path.Join(filepath.Dir(p), info.Name())

		// directory == namespace
		//nolint: nestif
		if info.IsDir() {
			ns := strings.TrimPrefix(strings.ReplaceAll(absPath, source, ""), "/")

			if ns == "" {
				return nil
			}

			// get the namespace name
			nsParts := strings.Split(ns, "/")
			nsName := nsParts[len(nsParts)-1]
			nsParent := strings.Join(nsParts[:len(nsParts)-1], "/")

			if ns != "" {
				if nsParent == "" {
					nsParent = "root"
				}

				fmt.Fprintf(writer, "[%s] restore namespace: \"%s\"\n", nsParent, nsName)
			}

			if err := vaultClient.CreateNamespaceErrorIfNotForced(nsParent, nsName, true); err != nil {
				return err
			}
		} else { // file == engine
			// get namespace and engine name
			engine := utils.RemoveExtension(filepath.Base(absPath))
			ns := strings.Trim(strings.Trim(strings.ReplaceAll(absPath, source, ""), info.Name()), utils.Delimiter)

			if ns == "" {
				ns = "root"
			}

			fmt.Fprintf(writer, "[%s] restore engine: %s\n", ns, engine)

			// create engine
			vaultClient.Client.SetNamespace(ns)

			if err := vaultClient.EnableKV2EngineErrorIfNotForced(true, engine); err != nil {
				return err
			}

			// read file
			input, err := fs.ReadFile(absPath)
			if err != nil {
				return err
			}

			// parse input
			json, err := utils.FromJSON(input)
			if err != nil {
				return err
			}

			// write secret
			if err := o.writeSecrets(json, vaultClient, ns, engine); err != nil {
				return err
			}
		}

		return nil
	})
}

func (o *snapshotRestoreOptions) writeSecrets(secrets map[string]interface{}, v *vault.Vault, ns, rootPath string) error {
	transformedMap := make(map[string]interface{})
	utils.FlattenMap(secrets, transformedMap, "")

	for p, m := range transformedMap {
		secrets, ok := m.(map[string]interface{})
		if !ok {
			log.Fatalf("cannot convert %T to map[string]interface", secrets)
		}

		if err := v.WriteSecrets(rootPath, p, secrets); err != nil {
			return fmt.Errorf("[%s] error writing secret \"%s\": %w", ns, p, err)
		}

		fmt.Fprintf(writer, "[%s] writing secret \"%s\" \n", ns, path.Join(rootPath, p))
	}

	return nil
}
