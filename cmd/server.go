package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"path"
	"strings"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type serverOptions struct {
	Port         string `env:"PORT" envDefault:"0.0.0.0:8080"`
	Path         string `env:"PATH"`
	EnginePath   string `env:"ENGINE_PATH"`
	SkipErrors   bool   `env:"SKIP_ERRORS" envDefault:"false"`
	LoginCommand string `env:"LoginCommand"`

	writer  *bytes.Buffer
	printer *printer.Printer
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
		PreRunE:       o.validateFlags,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.printer = printer.NewPrinter(
				printer.ShowValues(true),
				printer.ToFormat(printer.Export),
				printer.WithVaultClient(vaultClient),
				printer.WithWriter(o.writer),
			)

			fmt.Fprintf(writer, "mirroring secrets from path: \"%s\" to \"%s/export\"\n", o.Path, o.Port)

			return o.serve()
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

func (o *serverOptions) validateFlags(cmd *cobra.Command, args []string) error {
	switch {
	case o.EnginePath == "" && o.Path == "":
		return errors.New("no KV-paths given. Either --engine-path / -e or --path / -p needs to be specified")
	case o.EnginePath != "" && o.Path != "":
		return errors.New("cannot specify both engine-path and path")
	}

	return nil
}

func (o *serverOptions) buildMap() (map[string]interface{}, error) {
	var isSecretPath bool

	rootPath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

	// read recursive all secrets
	s, err := vaultClient.ListRecursive(rootPath, subPath, o.SkipErrors)
	if err != nil {
		return nil, err
	}

	// check if path is a directory or secret path
	if _, isSecret := vaultClient.ReadSecrets(rootPath, subPath); isSecret == nil {
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

func (o *serverOptions) serve() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/export", func(c *gin.Context) {
		// get format specified per request via url query param
		format, ok := c.GetQuery("format")
		if ok {
			switch strings.ToLower(format) {
			case "yaml", "yml":
				o.printer.WithOption(printer.ToFormat(printer.YAML))
			case "json":
				o.printer.WithOption(printer.ToFormat(printer.JSON))
			case "export":
				o.printer.WithOption(printer.ToFormat(printer.Export))
			case "markdown":
				o.printer.WithOption(printer.ToFormat(printer.Markdown))
			case "base":
				o.printer.WithOption(printer.ToFormat(printer.Base))
			case "policy":
				o.printer.WithOption(printer.ToFormat(printer.Policy))
			case "template", "tmpl":
				o.printer.WithOption(printer.ToFormat(printer.Template))
			}
		}

		c.Data(200, "text/plain", o.readSecrets())
	})

	return r.Run(o.Port)
}

func (o *serverOptions) readSecrets() []byte {
	o.writer.Reset()

	m, err := o.buildMap()
	if err != nil {
		log.Fatal(err)
	}

	enginePath, _ := utils.HandleEnginePath(o.EnginePath, o.Path)

	if err := o.printer.Out(enginePath, m); err != nil {
		log.Fatal(err)
	}

	return o.writer.Bytes()
}
