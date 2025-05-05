//go:build applicationtest

package applicationtest

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/go-tstr/golden"
	"github.com/go-tstr/tstr"
	"github.com/go-tstr/tstr/dep/cmd"
	"github.com/go-tstr/tstr/dep/compose"
	"github.com/stretchr/testify/require"
)

var (
	apiPort = mustFreePort()
	appURL  = fmt.Sprintf("http://127.0.0.1:%d", apiPort)
)

func TestMain(m *testing.M) {
	os.Setenv("POSTGRES_PORT", strconv.Itoa(mustFreePort()))
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")
	os.Setenv("POSTGRES_DB", "test")
	os.Setenv("POSTGRES_SSLMODE", "disable")
	os.Setenv("POSTGRES_HOST", "127.0.0.1")

	os.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.version=0.0.0-dev")
	os.Setenv("API_ADDR", fmt.Sprintf("127.0.0.1:%d", apiPort))

	golden.DefaultHandler.ProcessContent = golden.PrettyJSON
	tstr.RunMain(m, tstr.WithDeps(
		compose.New(
			compose.WithFile("../docker-compose.yaml"),
			compose.WithOsEnv(),
		),
		cmd.New(
			cmd.WithGoCode("../", "./cmd/demo"),
			cmd.WithReadyHTTP(appURL+"/api/v1/healthz"),
			cmd.WithEnvSet(os.Environ()...),
			cmd.WithGoCover(),
		),
	))
}

func TestHealthy(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, appURL+"/api/v1/healthz", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	golden.Request(t, http.DefaultClient, req, 200)
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
