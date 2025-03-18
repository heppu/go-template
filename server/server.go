// Package server wraps the application logic inside HTTP server.
package server

//go:generate go tool github.com/ogen-go/ogen/cmd/ogen --clean --config ../.ogen.yaml --target ../api -package api ./openapi.yaml
import (
	"cmp"
	"embed"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/heppu/go-template/api"
	"github.com/heppu/go-template/app"
)

//go:embed openapi.yaml
//go:embed swaggerui/*
var swaggerUI embed.FS

// New returns a function that creates a new http.Server.
// This functions composes the application and binds it to the server.
// Here you can configure custom middlewares, routes, etc.
//
// The binding chain goes as follows:
//
// http.Server -> http.ServeMux -> api.Server -> app.App -> store.Store
func New(s app.Store) func() (*http.Server, error) {
	return func() (*http.Server, error) {
		srv, err := api.NewServer(app.New(s))
		if err != nil {
			return nil, err
		}

		mux := http.NewServeMux()
		mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))
		mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.FS(swaggerUI))))

		addr := cmp.Or(os.Getenv("API_ADDR"), ":8080")
		slog.Info("Starting server at", slog.String("addr", addr))
		return &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       65 * time.Second,
		}, nil
	}
}
