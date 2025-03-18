package store_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/go-tstr/tstr"
	"github.com/go-tstr/tstr/dep/compose"
	"github.com/heppu/errgroup"
	"github.com/heppu/go-template/store"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("POSTGRES_PORT", strconv.Itoa(mustFreePort()))
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")
	os.Setenv("POSTGRES_DB", "test")
	os.Setenv("POSTGRES_SSLMODE", "disable")
	os.Setenv("POSTGRES_HOST", "127.0.0.1")
	fmt.Println("Opening DB in port", os.Getenv("POSTGRES_PORT"))

	tstr.RunMain(m, tstr.WithDeps(compose.New(
		compose.WithFile("../docker-compose.yaml"),
		compose.WithOsEnv(),
	)))
}

func TestStore(t *testing.T) {
	fn := func(t *testing.T) {
		s := store.New(store.WithDefaults()...)
		require.NoError(t, s.Init())
		eg := &errgroup.ErrGroup{}
		eg.Go(func() error { return s.Run() })
		require.NoError(t, s.Healthy(context.Background()))
		require.NoError(t, s.Stop())
		require.NoError(t, eg.Wait())
	}

	t.Run("Start on empty DB", fn)
	t.Run("Start on already initialized DB", fn)
}

func mustFreePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		log.Fatal(err)
	}

	return tcpAddr.Port
}
