package config

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func InitConfig() {

	// Determine proper data directory based on user context
	dataDir := "/data" // fallback

	currentUser, err := user.Current()
	if err != nil {
		log.Printf("Warning: unable to get current user: %v, using default data dir", err)
	} else if currentUser.Uid == "0" {
		// root/system app
		dataDir = "/var/lib/conveyor"
	} else {
		// normal user
		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			xdgData = filepath.Join(currentUser.HomeDir, ".local", "share")
		}
		dataDir = filepath.Join(xdgData, "conveyor")
	}

	// Load configuration from file
	viper.SetConfigName("conveyor")
	viper.AddConfigPath(dataDir)
	viper.SetConfigType("yaml")

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
