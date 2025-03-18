package main

import (
	// Timezone data for scratch image
	_ "time/tzdata"

	// Automated resource configuration for container envs
	_ "github.com/KimMachineGun/automemlimit"
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
		httpmod.New(httpmod.WithServerFn(server.New(s))),
	)
}
