package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
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

	configDir := createTmpConfigFile(t, func(t *testing.T, configFilePath string) {
		viper.Set("network.idle_timeout", "infinite")
		err := viper.WriteConfigAs(configFilePath)
		if err != nil {
			t.Fatal(err)
		}
	})
	defer os.RemoveAll(configDir)

	cfg, err := NewClientConfig(configDir)
	if err != nil {
		t.Fatal(err)
	}

	_ = cfg
}
