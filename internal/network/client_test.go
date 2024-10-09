package network_test

import (
	"TinyKVStore/internal/config"
	"TinyKVStore/internal/helpers/syncutil"
	"TinyKVStore/internal/network"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

const (
	ip          = "127.0.0.1"
	port        = "3223"
	invalidPort = "7789"

	tcpClientWrite    = "Hi, from tcp client"
	tcpServerResponse = "Hello, from tcp server"
)

func runMockServer(t *testing.T, emitServerReadyState syncutil.Emit) {
	tcpAddress := fmt.Sprintf("%s:%s", ip, port)
	addr, err := net.ResolveTCPAddr("tcp", tcpAddress)
	require.NoErrorf(t, err, "failed to resolve TCP address %s", tcpAddress)

	listener, err := net.ListenTCP("tcp", addr)
	require.NoErrorf(t, err, "failed to create tcp listener %s", addr.String())

	defer func() {
		require.NoErrorf(t, listener.Close(), "failed to close tcp listener %s", listener.Addr())
	}()

	emitServerReadyState()
	for {
		conn, err := listener.AcceptTCP()
		require.NoErrorf(t, err, "failed to accept tcp connection %s", err)

		go func(connection *net.TCPConn) {
			defer func() {
				require.NoErrorf(t, conn.Close(), "failed to close tcp connection %s", listener.Addr())
			}()

			buffer := make([]byte, 4096)
			_, cerr := conn.Read(buffer)
			if cerr != nil && cerr == io.EOF {
				return
			}

			require.NoErrorf(t, cerr, "failed to read tcp data %s", cerr)

			_, cerr = conn.Write([]byte(tcpServerResponse))
			if cerr != nil && cerr == io.EOF {
				return
			}

			require.NoErrorf(t, cerr, "failed to write tcp data %s", cerr)
		}(conn)
	}
}

func TestTcpClient(t *testing.T) {
	t.Parallel()

	notifyServerReadyState, waitForRunningServer := syncutil.NewEmitterAwaiter()
	go runMockServer(t, notifyServerReadyState)
	waitForRunningServer()

	testCases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "Failed to create tcp client",
			run: func(t *testing.T) {
				tcpClient, err := network.NewTCPClient(&config.ClientNetworkConfig{
					Address: "something similar to tcp address",
				})

				assert.Nil(t, tcpClient, "tcp client should be nil")
				assert.Errorf(t, err, "expected an error due to provided invalid tcp address")

				var addrErr *net.AddrError
				assert.ErrorAsf(t, err, &addrErr, "expected net.AddrError, actual error is %#v", err)
			},
		}, {
			name: "Connect to incorrect server address",
			run: func(t *testing.T) {
				tcpClient, err := network.NewTCPClient(&config.ClientNetworkConfig{
					Address: fmt.Sprintf("%s:%s", ip, invalidPort),
				})

				assert.NotNilf(t, tcpClient, "tcp client shouldn't be nil")
				assert.NoErrorf(t, err, "expected nil error while creating tcp client, actual error is %#v", err)

				err = tcpClient.Open()
				switch runtime.GOOS {
				case "windows":
					assert.ErrorIsf(t, err, windows.WSAECONNREFUSED, "expected WSAECONNREFUSED, actual error is %#v", err)
				case "darwin", "linux":
					assert.ErrorIsf(t, err, syscall.ECONNREFUSED, "expected ECONNREFUSED, actual error is %#v", err)
				default:
					assert.FailNow(t, fmt.Sprintf("not tested for %s os", runtime.GOOS))
				}
			},
		}, {
			name: "Client with deadline timeout",
			run: func(t *testing.T) {
				tcpClient, err := network.NewTCPClient(&config.ClientNetworkConfig{
					Address:     fmt.Sprintf("%s:%s", ip, port),
					IdleTimeout: time.Duration(1 * time.Millisecond),
				})

				assert.NotNilf(t, tcpClient, "tcp client shouldn't be nil")
				assert.NoErrorf(t, err, "expected nil error while creating tcp client, actual error is %#v", err)

				err = tcpClient.Open()
				assert.NoErrorf(t, err, "expected nil error while open tcp client connection, actual error is %#v", err)

				defer func() {
					err := tcpClient.Close()
					require.NoErrorf(t, err, "expected nil error while closing tcp client connection, actual error is %#v", err)
				}()

				// wait for shutdown connection by timeout
				time.Sleep(200 * time.Millisecond)

				err = tcpClient.Write([]byte(tcpClientWrite))
				assert.Errorf(t, err, "expected an error due to timeout after Write operation")
				assert.ErrorIsf(t, err, os.ErrDeadlineExceeded, "expected os.ErrDeadlineExceeded after Write operation, actual error is %#v", err)

				response, err := tcpClient.Read()
				assert.Errorf(t, err, "expected an error due to timeout after Read operation")
				assert.ErrorIsf(t, err, os.ErrDeadlineExceeded, "expected os.ErrDeadlineExceeded after Read operation, actual error is %#v", err)
				assert.Nil(t, response, "response should be nil, actual response is %#v", response)
			},
		}, {
			name: "Small r/w buffers",
			run: func(t *testing.T) {
				tcpClient, err := network.NewTCPClient(&config.ClientNetworkConfig{
					Address:        fmt.Sprintf("%s:%s", ip, port),
					IdleTimeout:    time.Duration(0),
					MaxMessageSize: "1 kb",
				})

				assert.NotNilf(t, tcpClient, "tcp client shouldn't be nil")
				assert.NoErrorf(t, err, "expected nil error while creating tcp client, actual error is %#v", err)

				err = tcpClient.Open()
				assert.NoErrorf(t, err, "expected nil error while open tcp client connection, actual error is %#v", err)

				defer func() {
					err := tcpClient.Close()
					require.NoErrorf(t, err, "expected nil error while closing tcp client connection, actual error is %#v", err)
				}()

				randBuffer := make([]byte, 2048)
				n, err := rand.Read(randBuffer)
				assert.Equalf(t, n, len(randBuffer), "expected buffer len is 2048, actual len is %d", n)
				assert.NoErrorf(t, err, "expected nil error while generating random sequence, actual error is %#v", err)

				err = tcpClient.Write(randBuffer)
				assert.Errorf(t, err, "expected an error due to write operation")

				response, err := tcpClient.Read()
				assert.Nilf(t, response, "tcp response should be nil")
				assert.Errorf(t, err, "expected an error due to read operation")
			},
		}, {
			name: "Successful TCP communication",
			run: func(t *testing.T) {
				tcpClient, err := network.NewTCPClient(&config.ClientNetworkConfig{
					Address:        fmt.Sprintf("%s:%s", ip, port),
					IdleTimeout:    time.Duration(0),
					MaxMessageSize: "1 kb",
				})

				assert.NotNilf(t, tcpClient, "tcp client shouldn't be nil")
				assert.NoErrorf(t, err, "expected nil error while creating tcp client, actual error is %#v", err)

				err = tcpClient.Open()
				assert.NoErrorf(t, err, "expected nil error while open tcp client connection, actual error is %#v", err)

				defer func() {
					err := tcpClient.Close()
					require.NoErrorf(t, err, "expected nil error while closing tcp client connection, actual error is %#v", err)
				}()

				err = tcpClient.Write([]byte(tcpClientWrite))
				assert.NoErrorf(t, err, "expected nil error while tcp client write operation, actual error is %#v", err)

				response, err := tcpClient.Read()
				assert.NotNilf(t, response, "tcp response shouldn't be nil")
				assert.NoErrorf(t, err, "expected nil error while tcp client read operation, actual error is %#v", err)

				assert.Equalf(t, tcpServerResponse, string(response), "invalid response data, expected response \"%s\", actual response is \"%s\"", tcpServerResponse, response)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.run)
	}
}
