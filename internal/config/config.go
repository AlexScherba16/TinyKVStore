package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	configDirectory           = "config"
	configClientFile          = "client"
	configClientFileExtension = "yaml"
)

// ClientConfig client application config context.
type ClientConfig struct {
	LoggerCfg  LoggerConfig
	NetworkCfg ClientNetworkConfig
}

// LoggerConfig logger configuration.
type LoggerConfig struct {
	Level  string
	Output string
}

// ClientNetworkConfig client network configuration.
type ClientNetworkConfig struct {
	Address        string
	MaxMessageSize string
	IdleTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

// newLoggerConfig logger config parser, returns LoggerConfig.
func newLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:  viper.GetString("logging.level"),
		Output: viper.GetString("logging.output"),
	}
}

// newLoggerConfig client application network config, returns ClientNetworkConfig.
func newNetworkConfig() ClientNetworkConfig {
	return ClientNetworkConfig{
		Address:        viper.GetString("network.address"),
		MaxMessageSize: viper.GetString("network.max_message_size"),
		IdleTimeout:    viper.GetDuration("network.idle_timeout"),
		ReadTimeout:    viper.GetDuration("network.read_timeout"),
		WriteTimeout:   viper.GetDuration("network.write_timeout"),
	}
}

// NewClientConfig client application config parser, returns ClientConfig
// or error if something went wrong during initialization.
func NewClientConfig(configPath string) (*ClientConfig, error) {
	viper.SetConfigName(configClientFile)
	viper.SetConfigType(configClientFileExtension)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &ClientConfig{
		LoggerCfg:  newLoggerConfig(),
		NetworkCfg: newNetworkConfig(),
	}, nil
}
