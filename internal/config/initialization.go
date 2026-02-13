package config

import (
	"fmt"
	"log"
	"os"

	"github.com/open-ug/conveyor/pkg/types"
	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func InitConfig() {

	configDir := "/etc/conveyor"

	// Load configuration from file
	viper.SetConfigName("conveyor")
	viper.AddConfigPath(configDir)
	viper.SetConfigType("yaml")

	configFileEnv := os.Getenv("CONVEYOR_CONFIG_FILE")
	if configFileEnv != "" {
		viper.SetConfigFile(configFileEnv)
	}

}

// LoadConfig reads in config file and ENV variables if set.
// This should be called in PreRunE of commands that need config.
func LoadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; provide a helpful error
			return fmt.Errorf("config file not found. Please run 'conveyor init' first")
		}
		// Config file was found but another error was produced
		return fmt.Errorf("error reading config file: %w", err)
	}

	log.Printf("Using config file: %s", viper.ConfigFileUsed())
	return nil
}

func LoadTestEnvConfig(testConfigPath string) error {
	viper.SetConfigFile(testConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading test config file: %w", err)
	}

	log.Printf("Using test config file: %s", viper.ConfigFileUsed())
	return nil
}

func GetTestConfig() (types.ServerConfig, error) {
	var config types.ServerConfig
	err := viper.Unmarshal(&config)
	if err != nil {
		return types.ServerConfig{}, fmt.Errorf("unable to decode into struct: %w", err)
	}
	return config, nil
}
