package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/spf13/cobra"
)

type serverOptions struct {
	Port       string `env:"PORT" envDefault:"0.0.0.0:8080"`
	Path       string `env:"PATH"`
	EnginePath string `env:"ENGINE_PATH"`
	SkipErrors bool   `env:"SKIP_ERRORS" envDefault:"false"`

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
		PreRunE:       o.validateFlags,
		RunE: func(cmd *cobra.Command, args []string) error {
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
	rootPath, subPath := utils.HandleEnginePath(o.EnginePath, o.Path)

	// read recursive all secrets
	s, err := vaultClient.ListRecursive(rootContext, rootPath, subPath, o.SkipErrors)
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

func (o *serverOptions) serve() error {
	server := &http.Server{
		Addr:              o.Port,
		Handler:           loggingMiddleware(o.export()),
		ReadHeaderTimeout: 5 * time.Second,
		BaseContext:       func(l net.Listener) context.Context { return rootContext },
	}

	slog.Info("Server started",
		slog.String("address", server.Addr),
		slog.String("port", o.Port),
		slog.String("path", o.Path),
	)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("server failed to start: %w", err)
		}
	}()

	<-rootContext.Done()

	slog.Info("Shutdown signal received")
	slog.Info("Server stopped")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(rootContext, 3*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "err", err)
	} else {
		slog.Info("Server shut down gracefully")
	}

	return nil
}

// nolint: cyclop
func (o *serverOptions) export() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		format := r.URL.Query().Get("format")
		enginePath, _ := utils.HandleEnginePath(o.EnginePath, o.Path)

		if format == "" {
			format = "export"
		}

		opts := []prt.Option{
			prt.ShowValues(true),
			prt.WithVaultClient(vaultClient),
			prt.WithWriter(o.writer),
			prt.WithEnginePath(enginePath),
			prt.ToFormat(prt.Export),
			//nolint: contextcheck
			prt.WithContext(rootContext),
		}

		switch strings.ToLower(format) {
		case "yaml", "yml":
			opts = append(opts, prt.ToFormat(prt.YAML))
		case "json":
			opts = append(opts, prt.ToFormat(prt.JSON))
		case "export":
			opts = append(opts, prt.ToFormat(prt.Export))
		case "markdown":
			opts = append(opts, prt.ToFormat(prt.Markdown))
		case "base":
			opts = append(opts, prt.ToFormat(prt.Base))
		case "policy":
			opts = append(opts, prt.ToFormat(prt.Policy))
		case "template", "tmpl":
			opts = append(opts, prt.ToFormat(prt.Template))
		default:
			//nolint: perfsprint
			http.Error(w, fmt.Sprintf("unsupported format: %s", format), http.StatusBadRequest)

			return
		}

		printer = prt.NewSecretPrinter(opts...)

		defer o.writer.Reset()

		w.WriteHeader(http.StatusOK)
		//nolint: errcheck, contextcheck
		w.Write(func() []byte {
			m, err := o.buildMap()
			if err != nil {
				log.Fatal(err)
			}

			if err := printer.Out(m); err != nil {
				log.Fatal(err)
			}

			return o.writer.Bytes()
		}())
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", duration),
			slog.String("remote", r.RemoteAddr),
		)
	})
}
