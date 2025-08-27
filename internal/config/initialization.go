package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var CfgFile string

// initConfig reads in config file and ENV variables if set.
func InitConfig() {

	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Error loading .env file")
	}
	viper.AutomaticEnv()

	viper.BindEnv("api.host", "CONVEYOR_SERVER_HOST")
	viper.BindEnv("etcd.data", "ETCD_ENDPOINT")
	viper.BindEnv("loki.host", "LOKI_ENDPOINT")
	viper.BindEnv("nats.url", "NATS_URL")

	viper.SetDefault("api.host", "http://localhost:8080")
	viper.SetDefault("loki.host", "http://localhost:3100")
	viper.SetDefault("nats.url", "nats://localhost:4222")

}
