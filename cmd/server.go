package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	srv "github.com/FalcoSuessgott/vkv/pkg/http"
	prt "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

type serverOptions struct {
	Port         string `env:"PORT" envDefault:"0.0.0.0:8080"`
	Path         string `env:"PATH"`
	EnginePath   string `env:"ENGINE_PATH"`
	SkipErrors   bool   `env:"SKIP_ERRORS" envDefault:"false"`
	LoginCommand string `env:"LoginCommand"`

	writer *bytes.Buffer
}

func defaultServerOptions() *serverOptions {
	return &serverOptions{
		writer: bytes.NewBufferString(""),
	}
}

// NewServerCmd export subcommand.
// nolint: lll
func NewServerCmd() *cobra.Command {
	o := defaultServerOptions()

	if err := utils.ParseEnvs(envVarServerPrefix, o); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "server",
		Short:         "expose a http server that returns the read secrets from Vault, useful during CI",
		SilenceUsage:  true,
		SilenceErrors: true,
		// PreRunE:       o.validateFlags,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(writer, "proxying secrets via \"%s/export?path=<PATH>?format=<FORMAT>\"\n", o.Port)

			ctx := context.Background()

			// otherwise create a new vault client
			vaultClient = &vault.Vault{}

			c, err := api.NewClient(api.DefaultConfig())
			if err != nil {
				return err
			}

			vaultClient.Client = c

			vaultClient.Client.SetAddress(os.Getenv("VAULT_ADDR"))
			// main handler for /export

			handler := srv.NewServer(vaultClient, prt.NewSecretPrinter())

			// server
			httpServer := &http.Server{
				Addr:    o.Port,
				Handler: handler,
			}

			// run in background
			go func() {
				log.Printf("listening on %s\n", o.Port)

				if err := httpServer.ListenAndServe(); err != nil {
					fmt.Fprintln(writer, "error listening and serving: %w", err)
				}
			}()

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()

				<-ctx.Done()

				shutdownCtx := context.Background()
				shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

				defer cancel()

				if err := httpServer.Shutdown(shutdownCtx); err != nil {
					fmt.Fprintf(writer, "error shutting down http server: %s\n", err)
				}
			}()

			wg.Wait()

			return nil
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringVarP(&o.Port, "port", "P", o.Port, "HTTP Server Port (env: VKV_SERVER_PORT)")
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KVv2 Engine path (env: VKV_SERVER_PATH)")
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath, "engine path in case your KV-engine contains special characters such as \"/\", the path value will then be appended if specified (\"<engine-path>/<path>\") (env: VKV_SERVER_ENGINE_PATH)")
	cmd.Flags().BoolVar(&o.SkipErrors, "skip-errors", o.SkipErrors, "dont exit on errors (permission denied, deleted secrets) (env: VKV_SERVER_SKIP_ERRORS)")

	return cmd
}

// func (o *serverOptions) validateFlags(cmd *cobra.Command, args []string) error {
// 	switch {
// 	case o.EnginePath == "" && o.Path == "":
// 		return errors.New("no KV-paths given. Either --engine-path / -e or --path / -p needs to be specified")
// 	case o.EnginePath != "" && o.Path != "":
// 		return errors.New("cannot specify both engine-path and path")
// 	}

// 	return nil
// }

func (o *serverOptions) buildMap() (map[string]interface{}, error) {
	rootPath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

	// read recursive all secrets
	s, err := vaultClient.ListRecursive(rootPath, subPath, o.SkipErrors)
	if err != nil {
		return nil, err
	}

	path := path.Join(rootPath, subPath)
	if o.EnginePath != "" {
		path = subPath
	}

	// prepare the output map
	pathMap := utils.UnflattenMap(path, utils.ToMapStringInterface(s))

	if o.EnginePath != "" {
		return map[string]interface{}{
			o.EnginePath: pathMap,
		}, nil
	}

	return pathMap, nil
}

func (o *serverOptions) readSecrets() []byte {
	o.writer.Reset()

	m, err := o.buildMap()
	if err != nil {
		log.Fatal(err)
	}

	if err := printer.Out(m); err != nil {
		log.Fatal(err)
	}

	return o.writer.Bytes()
}
