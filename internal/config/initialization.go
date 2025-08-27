package config

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var CfgFile string

// initConfig reads in config file and ENV variables if set.
func InitConfig() {

	// Load .env if present
	_ = godotenv.Load() // silently ignore if not found

	// Bind environment variables
	viper.AutomaticEnv()
	viper.BindEnv("api.host", "CONVEYOR_SERVER_HOST")
	viper.BindEnv("api.data", "CONVEYOR_DATA_DIR")
	viper.BindEnv("loki.host", "LOKI_ENDPOINT")

	// Set defaults
	viper.SetDefault("api.host", "http://localhost:8080")
	viper.SetDefault("loki.host", "http://localhost:3100")

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

	viper.SetDefault("api.data", dataDir)

}
