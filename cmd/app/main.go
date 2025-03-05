package main

import (
	"cmp"
	"fmt"
	"os"

	"github.com/elisasre/go-common/v2/service"
	"github.com/elisasre/go-common/v2/service/module/httpserver"
	"github.com/elisasre/go-common/v2/service/module/httpserver/pprof"
	"github.com/elisasre/go-common/v2/service/module/siglistener"
	"github.com/heppu/go-template/api"
	"github.com/heppu/go-template/app"
	"github.com/heppu/go-template/store"
)

func main() {
	str := store.New(store.FromEnv())
	service.RunAndExit(service.Modules{
		siglistener.New(os.Interrupt),
		httpserver.New(
			pprof.WithProfiling(),
			httpserver.WithAddr(GetEnv("PPROF_ADDR", ":8081")),
		),
		str,
		httpserver.New(
			httpserver.WithAddr(GetEnv("API_ADDR", ":8080")),
			func(s *httpserver.Server) error {
				srv, err := api.NewServer(app.New(str))
				if err != nil {
					return fmt.Errorf("failed to create api server: %w", err)
				}
				return httpserver.WithHandler(srv)(s)
			},
		),
	})
}

func GetEnv(key, def string) string {
	return cmp.Or(os.Getenv(key), def)
}
