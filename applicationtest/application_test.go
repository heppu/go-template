//go:build applicationtest

package applicationtest

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/go-tstr/tstr"
	"github.com/go-tstr/tstr/dep/cmd"
	"github.com/go-tstr/tstr/dep/compose"
	"github.com/stretchr/testify/assert"
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

	os.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.version=0.1.0")
	os.Setenv("API_ADDR", fmt.Sprintf("127.0.0.1:%d", apiPort))

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
	resp, err := http.Get(appURL + "/api/v1/healthz")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, `{"message":"OK"}`, string(data))
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
