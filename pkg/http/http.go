package http

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Authenticator interface {
	Authenticate(string) error
}

type Printer interface {
	Print(v any, format string) ([]byte, error)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func NewServer(a Authenticator, p Printer) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/export", Chain(PrintSecrets(p).ServeHTTP, Logging(), AuthHandler(a)))

	return mux
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

// use slog

func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				log.Println(r.URL.Path, r.URL.Query().Encode(), time.Since(start))
			}()

			f(w, r)
		}
	}
}

// print secrets here
func PrintSecrets(p Printer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get path
		_, ok := r.URL.Query()["path"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error no path found but required.\n"))

			return
		}

		// get secrets
		// get format
		format, _ := r.URL.Query()["format"]

		// print
		out, err := p.Print(nil, format[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error printing secrets: %s\n", err.Error())))

			return
		}

		// return
		w.Write(out)
	}
}

func AuthHandler(a Authenticator) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header

			// check for token
			token, ok := headers["X-Vault-Token"]
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Vault Token required in \"X-Vault-Token\" request header.\n"))

				return
			}

			if len(token) > 1 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("only one token is allowed\n"))

				return
			}

			// login
			if err := a.Authenticate(token[0]); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(fmt.Sprintf("error authenticating to Vault: %s\n", err.Error())))

				return
			}

			// call next middleware
			f.ServeHTTP(w, r)
		})
	}
}
