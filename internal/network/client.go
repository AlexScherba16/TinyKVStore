package network

import (
	"TinyKVStore/internal/config"
	"errors"
	"net"
	"time"
)

const (
	defaultClientTimeOut = 0
)

type TCPClient struct {
	connection  net.Conn
	idleTimeout time.Duration
	bufferSize  int
}

// NewTCPClient creates a new tcp client, or error if something went wrong.
func NewTCPClient(networkCfg *config.ClientNetworkConfig) (*TCPClient, error) {
	connection, err := net.Dial("tcp", networkCfg.Address)
	if err != nil {
		return nil, err
	}

	connection.SetDeadline(time.Now().Add(defaultClientTimeOut * time.Second))

	return nil, errors.New("Not yet implemented")
}
