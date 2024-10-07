package composer

import (
	"TinyKVStore/internal/application"
	"TinyKVStore/internal/application/client"
	"TinyKVStore/internal/args"
	"TinyKVStore/internal/config"
	"TinyKVStore/internal/logger/zap"
	"TinyKVStore/internal/network"
)

// ComposeNewClientApplication creates and returns a new client instance of IApplication.
// provides a flexible way to initialize the client application.
// Returns error if something went wrong during initialization.
func ComposeNewClientApplication() (application.IApplication, error) {
	flags, err := args.NewClientFlags()
	if err != nil {
		return nil, err
	}

	cfg, err := config.NewClientConfig(flags.ConfigDir())
	if err != nil {
		return nil, err
	}

	zapLogger, err := zap.NewLogger(&cfg.LoggerCfg)
	if err != nil {
		return nil, err
	}

	tcpClient, err := network.NewTCPClient(&cfg.NetworkCfg)
	if err != nil {
		return nil, err
	}

	_ = tcpClient

	return client.NewClientApplication(zapLogger)
}
