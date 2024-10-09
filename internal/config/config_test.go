package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func createTmpConfigFile(t *testing.T, setConfigFn func(t *testing.T, configFilePath string)) string {
	confDir, err := os.MkdirTemp("", "*-client_conf")
	if err != nil {
		t.Fatal(err)
	}

	fileName := fmt.Sprintf("%s.%s", configClientFile, configClientFileExtension)
	configFilePath := filepath.Join(confDir, fileName)
	f, err := os.Create(configFilePath)
	if err != nil {
		os.RemoveAll(confDir)
		t.Fatal(err)
	}
	f.Close()

	setConfigFn(t, configFilePath)

	return confDir
}

func TestNewClientConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "Network infinite idle timeout",
			run: func(t *testing.T) {
				configDir := createTmpConfigFile(t, func(t *testing.T, configFilePath string) {
					viper.Set("network.idle_timeout", "infinite")
					err := viper.WriteConfigAs(configFilePath)
					require.NoErrorf(t, err, "failed to write config %s", err)
				})
				defer os.RemoveAll(configDir)

				cfg, err := NewClientConfig(configDir)
				require.NoErrorf(t, err, "failed to create client config %s", err)
				require.Equal(t, time.Duration(0), cfg.NetworkCfg.IdleTimeout)
			},
		}, {
			name: "Network read/write timeouts",
			run: func(t *testing.T) {
				configDir := createTmpConfigFile(t, func(t *testing.T, configFilePath string) {
					viper.Set("network.read_timeout", "infinite")
					viper.Set("network.write_timeout", "5s")
					err := viper.WriteConfigAs(configFilePath)
					require.NoErrorf(t, err, "failed to write config %s", err)
				})
				defer os.RemoveAll(configDir)

				cfg, err := NewClientConfig(configDir)
				require.NoErrorf(t, err, "failed to create client config %s", err)
				require.Equal(t, time.Duration(0), cfg.NetworkCfg.ReadTimeout)
				require.Equal(t, time.Duration(5*time.Second), cfg.NetworkCfg.WriteTimeout)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.run)
	}
}
