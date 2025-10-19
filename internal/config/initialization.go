package config

import (
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

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}
	log.Printf("Using config file: %s", viper.ConfigFileUsed())

}
