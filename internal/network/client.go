package network

import (
	"TinyKVStore/internal/config"
	"TinyKVStore/internal/helpers/parser"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	defaultClientAddress    = "127.0.0.1:3223"
	defaultClientTimeOut    = 0
	defaultClientBufferSize = 4096

	readTimeOut  = time.Duration(5 * time.Second)
	writeTimeOut = time.Duration(5 * time.Second)
)

// clientNetworkCfg private client network config, use it after validation config.ClientNetworkConfig.
type clientNetworkCfg struct {
	address     string
	bufferSize  int
	idleTimeout time.Duration
}

// validateConfig ensures that the correct settings will be used for the client network configuration.
func validateConfig(networkCfg *config.ClientNetworkConfig) clientNetworkCfg {
	clientAddress := networkCfg.Address
	if clientAddress == "" {
		clientAddress = defaultClientAddress
	}

	bufferSize := parser.ParseBufferSize(networkCfg.MaxMessageSize)
	if bufferSize == -1 {
		bufferSize = defaultClientBufferSize
	}

	return clientNetworkCfg{
		address:     clientAddress,
		bufferSize:  bufferSize,
		idleTimeout: networkCfg.IdleTimeout,
	}
}

type TCPClient struct {
	address     *net.TCPAddr
	connection  *net.TCPConn
	idleTimeout time.Duration
	bufferSize  int
}

// NewTCPClient creates a new tcp client, or error if something went wrong.
func NewTCPClient(networkCfg *config.ClientNetworkConfig) (*TCPClient, error) {
	validatedCfg := validateConfig(networkCfg)

	tcpAddress, err := net.ResolveTCPAddr("tcp", validatedCfg.address)
	if err != nil {
		return nil, err
	}

	return &TCPClient{
		address:     tcpAddress,
		connection:  nil,
		idleTimeout: validatedCfg.idleTimeout,
		bufferSize:  validatedCfg.bufferSize,
	}, nil
}

// Open TCP connection, or error if something went wrong.
func (c *TCPClient) Open() error {
	if c.connection != nil {
		return errors.New("connection already open")
	}

	connection, err := net.DialTCP("tcp", nil, c.address)
	if err != nil {
		return err
	}

	deadlineTime := time.Time{}
	if c.idleTimeout != 0 {
		deadlineTime = time.Now().Add(c.idleTimeout)
	}
	err = connection.SetDeadline(deadlineTime)
	if err != nil {
		return err
	}

	err = connection.SetReadBuffer(c.bufferSize)
	if err != nil {
		return err
	}

	err = connection.SetWriteBuffer(c.bufferSize)
	if err != nil {
		return err
	}

	c.connection = connection
	return nil
}

// Close TCP connection, or error if something went wrong.
func (c *TCPClient) Close() error {
	if c.connection == nil {
		return errors.New("connection already closed")
	}

	err := c.connection.Close()
	c.connection = nil
	return err
}

// Write data to server, returns error if something went wrong.
func (c *TCPClient) Write(b []byte) error {
	if c.connection == nil {
		return errors.New("invalid tcp connection for writing")
	}

	type writeResult struct {
		err error
	}
	writeCh := make(chan writeResult)

	go func() {
		defer close(writeCh)

		writeLen := len(b)
		if writeLen > c.bufferSize {
			writeCh <- writeResult{err: fmt.Errorf("payload size %d is overflow buffer size %d", writeLen, c.bufferSize)}
			return
		}

		n, err := c.connection.Write(b)
		if err != nil {
			writeCh <- writeResult{err: err}
			return
		}

		if n != writeLen {
			writeCh <- writeResult{err: fmt.Errorf("invalid transfered data len, expected : %d, actual : %d", len(b), n)}
			return
		}

		writeCh <- writeResult{err: nil}
	}()

	select {
	case result := <-writeCh:
		return result.err

	case <-time.After(writeTimeOut):
		return errors.New("expired timeout during tcp write operation")
	}
}

// Read server response, returns error if something went wrong.
func (c *TCPClient) Read() ([]byte, error) {
	if c.connection == nil {
		return nil, errors.New("invalid tcp connection for reading")
	}

	type readResult struct {
		response []byte
		err      error
	}
	readCh := make(chan readResult)

	go func() {
		defer close(readCh)

		response, err := io.ReadAll(c.connection)
		if err != nil {
			readCh <- readResult{nil, err}
			return
		}

		respSize := len(response)
		if respSize > c.bufferSize {
			readCh <- readResult{nil, fmt.Errorf("response invalid size, response payload size : %d, buffer size : %d", respSize, c.bufferSize)}
			return
		}

		readCh <- readResult{response, nil}
	}()

	select {
	case result := <-readCh:
		return result.response, result.err

	case <-time.After(readTimeOut):
		return nil, errors.New("expired timeout during tcp read operation")
	}
}
