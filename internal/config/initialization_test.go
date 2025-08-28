package config_test

import (
	"os"
	"testing"

	"github.com/open-ug/conveyor/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	// Backup original env vars
	originalEnv := map[string]string{
		"CONVEYOR_SERVER_HOST": os.Getenv("CONVEYOR_SERVER_HOST"),
		"LOKI_ENDPOINT":        os.Getenv("LOKI_ENDPOINT"),
	}

	// Set test environment variables
	os.Setenv("CONVEYOR_SERVER_HOST", "127.0.0.1")
	os.Setenv("LOKI_ENDPOINT", "http://loki:3100")

	// Clear viper keys in case another test already ran
	viper.Reset()

	config.InitConfig()

	assert.Equal(t, "127.0.0.1", viper.GetString("api.host"))
	assert.Equal(t, "http://loki:3100", viper.GetString("loki.host"))

	// Unset test env vars and restore original
	for k, v := range originalEnv {
		if v == "" {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}

	viper.Reset()
}
