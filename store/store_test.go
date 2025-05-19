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
	mustSetenv("POSTGRES_PORT", strconv.Itoa(mustFreePort()))
	mustSetenv("POSTGRES_USER", "test")
	mustSetenv("POSTGRES_PASSWORD", "test")
	mustSetenv("POSTGRES_DB", "test")
	mustSetenv("POSTGRES_SSLMODE", "disable")
	mustSetenv("POSTGRES_HOST", "127.0.0.1")
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

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		log.Fatal(err)
	}

	if err := listener.Close(); err != nil {
		log.Fatal("Failed to close listener:", err)
	}

	return tcpAddr.Port
}

func mustSetenv(k, v string) {
	if err := os.Setenv(k, v); err != nil {
		log.Fatal(err)
	}
}
