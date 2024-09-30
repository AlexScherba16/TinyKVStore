package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

const (
	configDirectory           = "config"
	configClientFile          = "client"
	configClientFileExtension = "yaml"
)

// ClientConfig client application config context.
type ClientConfig struct {
	LoggerCfg  LoggerConfig
	NetworkCfg NetworkConfig
}

// LoggerConfig logger configuration.
type LoggerConfig struct {
	Level  string
	Output string
}

// NetworkConfig client network configuration.
type NetworkConfig struct {
	Address        string
	MaxConnection  uint16
	MaxMessageSize string
	IdleTimeout    time.Duration
}

func getClientConfigPath() (string, error) {
	projectRootPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(projectRootPath, configDirectory)
	return configPath, nil
}

// newLoggerConfig logger config parser, returns LoggerConfig.
func newLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:  viper.GetString("logging.level"),
		Output: viper.GetString("logging.output"),
	}
}

// newLoggerConfig network config parser, returns NetworkConfig.
func newNetworkConfig() NetworkConfig {
	return NetworkConfig{
		Address:        viper.GetString("network.address"),
		MaxConnection:  viper.GetUint16("network.max_connections"),
		MaxMessageSize: viper.GetString("network.max_message_size"),
		IdleTimeout:    viper.GetDuration("network.idle_timeout"),
	}
}

// NewClientConfig client application config parser, returns ClientConfig
// or error if something went wrong during initialization.
func NewClientConfig() (*ClientConfig, error) {
	clientCfgPath, err := getClientConfigPath()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName(configClientFile)
	viper.SetConfigType(configClientFileExtension)
	viper.AddConfigPath(clientCfgPath)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &ClientConfig{
		LoggerCfg:  newLoggerConfig(),
		NetworkCfg: newNetworkConfig(),
	}, nil
}