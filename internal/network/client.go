package network

import (
	"TinyKVStore/internal/config"
	"errors"
	"net"
	"time"
)

type TCPClient struct {
	connection  net.Conn
	idleTimeout time.Duration
	bufferSize  int
}

// NewTCPClient creates a new tcp client, or error if something went wrong.
func NewTCPClient(networkCfg *config.NetworkConfig) (*TCPClient, error) {
	return nil, errors.New("Not yet implemented")
}
