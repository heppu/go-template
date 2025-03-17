package server

import (
	"cmp"
	"net/http"
	"os"
	"time"

	"github.com/heppu/go-template/api"
	"github.com/heppu/go-template/app"
)

// New returns a function that creates a new http.Server.
// This functions binds composes the application and binds it to the server.
// Here you can configure custom middlewares, routes, etc.
func New(s app.Store) func() (*http.Server, error) {
	return func() (*http.Server, error) {
		srv, err := api.NewServer(app.New(s))
		if err != nil {
			return nil, err
		}
		return &http.Server{
			Addr:              cmp.Or(os.Getenv("API_ADDR"), ":8080"),
			Handler:           srv,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       65 * time.Second,
		}, nil
	}
}
