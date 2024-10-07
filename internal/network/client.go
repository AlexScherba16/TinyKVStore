package network

import (
	"TinyKVStore/internal/config"
	"TinyKVStore/internal/helpers"
	"errors"
	"net"
	"time"
)

const (
	defaultClientAddress    = "127.0.0.1:3223"
	defaultClientTimeOut    = 0
	defaultClientBufferSize = 4096
)

// clientNetworkCfg private client network config, use it after validation config.ClientNetworkConfig.
type clientNetworkCfg struct {
	address    string
	bufferSize int
}

// validateConfig ensures that the correct settings will be used for the client network configuration.
func validateConfig(networkCfg *config.ClientNetworkConfig) clientNetworkCfg {
	clientAddress := networkCfg.Address
	if clientAddress == "" {
		clientAddress = defaultClientAddress
	}

	bufferSize := helpers.ParseBufferSize(networkCfg.MaxMessageSize)
	if bufferSize == -1 {
		bufferSize = defaultClientBufferSize
	}

	return clientNetworkCfg{
		address:    clientAddress,
		bufferSize: bufferSize,
	}
}

type TCPClient struct {
	connection  net.Conn
	idleTimeout time.Duration
	bufferSize  int
}

// NewTCPClient creates a new tcp client, or error if something went wrong.
func NewTCPClient(networkCfg *config.ClientNetworkConfig) (*TCPClient, error) {
	validatedCfg := validateConfig(networkCfg)
	_ = validatedCfg

	return nil, errors.New("Not yet implemented")
}
