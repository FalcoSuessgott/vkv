package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const envVarExportPrefix = "VKV_SERVER_"

type serverOptions struct {
	Port         string `env:"PORT" envDefault:"8080"`
	Path         string `env:"PATH"`
	EnginePath   string `env:"ENGINE_PATH"`
	SkipErrors   bool   `env:"SKIP_ERRORS" envDefault:"false"`
	LoginCommand string `env:"LoginCommand"`
}

// NewServerCmd export subcommand.
//nolint:lll
func NewServerCmd(writer io.Writer, vaultClient *vault.Vault) *cobra.Command {
	var err error

	o := &serverOptions{}

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "server",
		Short:         "expose a http server that returns the read secrets from Vault, useful during CI",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// vault auth
			if vaultClient == nil {
				if vaultClient, err = vault.NewDefaultClient(); err != nil {
					return err
				}
			}

			b := bytes.NewBufferString("")

			// prepare printer
			printer := printer.NewPrinter(
				printer.ShowValues(true),
				printer.ToFormat(printer.Export),
				printer.WithVaultClient(vaultClient),
				printer.WithWriter(b),
			)

			// prepare map
			m, err := o.buildMap(vaultClient)
			if err != nil {
				return err
			}

			// print secrets
			enginePath, _ := utils.HandleEnginePath(o.EnginePath, o.Path)

			if err := printer.Out(enginePath, m); err != nil {
				return err
			}

			return o.serve(b.Bytes())
		},
	}

	cmd.Flags().SortFlags = false

	// Input
	cmd.Flags().StringVarP(&o.Port, "port", "P", o.Port, "HTTP Server Port (env: VKV_SERVER_PATH)")
	cmd.Flags().StringVarP(&o.Path, "path", "p", o.Path, "KVv2 Engine path (env: VKV_SERVER_PATH)")
	cmd.Flags().StringVarP(&o.EnginePath, "engine-path", "e", o.EnginePath, "engine path in case your KV-engine contains special characters such as \"/\", the path value will then be appended if specified (\"<engine-path>/<path>\") (env: VKV_SERVER_ENGINE_PATH)")
	cmd.Flags().BoolVar(&o.SkipErrors, "skip-errors", o.SkipErrors, "dont exit on errors (permission denied, deleted secrets) (env: VKV_SERVER_SKIP_ERRORS)")

	return cmd
}

func (o *serverOptions) buildMap(v *vault.Vault) (map[string]interface{}, error) {
	var isSecretPath bool

	rootPath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

	// read recursive all secrets
	s, err := v.ListRecursive(rootPath, subPath, o.SkipErrors)
	if err != nil {
		return nil, err
	}

	// check if path is a directory or secret path
	if _, isSecret := v.ReadSecrets(rootPath, subPath); isSecret == nil {
		isSecretPath = true
	}

	path := path.Join(rootPath, subPath)
	if o.EnginePath != "" {
		path = subPath
	}

	// prepare the output map
	pathMap := utils.PathMap(path, utils.ToMapStringInterface(s), isSecretPath)

	if o.EnginePath != "" {
		return map[string]interface{}{
			o.EnginePath: pathMap,
		}, nil
	}

	return pathMap, nil
}

func (o *serverOptions) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarExportPrefix,
	}); err != nil {
		return err
	}

	return nil
}

func (o *serverOptions) serve(secrets []byte) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/export", func(c *gin.Context) {
		c.Data(200, "text/plain", secrets)
	})

	return r.Run(fmt.Sprintf(":%s", o.Port))
}
