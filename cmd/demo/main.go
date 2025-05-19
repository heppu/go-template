package main

import (
	"cmp"
	"net/http"
	"os"

	// Instrument http.DefaultServeMux with pprof
	_ "net/http/pprof" //nolint: gosec
	// Import time zone data
	_ "time/tzdata"

	// Set GOMEMLIMIT automatically
	_ "github.com/KimMachineGun/automemlimit"
	// Set GOMAXPROCS automatically
	_ "go.uber.org/automaxprocs"

	"github.com/go-srvc/mods/httpmod"
	"github.com/go-srvc/mods/logmod"
	"github.com/go-srvc/mods/metermod"
	"github.com/go-srvc/mods/sigmod"
	"github.com/go-srvc/mods/tracemod"
	"github.com/go-srvc/srvc"
	"github.com/heppu/go-template/server"
	"github.com/heppu/go-template/store"
)

func main() {
	s := store.New(store.WithDefaults()...)

	srvc.RunAndExit(
		logmod.New(),
		sigmod.New(),
		tracemod.New(),
		metermod.New(),
		s,
		httpmod.New(
			httpmod.WithAddr(cmp.Or(os.Getenv("PPROF_ADDR"), ":6060")),
			httpmod.WithHandler(http.DefaultServeMux),
		),
		httpmod.New(httpmod.WithServerFn(server.New(s))),
	)
}
