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

	//check for the following environment variables

	viper.BindEnv("api.host", "CONVEYOR_SERVER_HOST")
	viper.BindEnv("api.port", "CONVEYOR_SERVER_PORT")
	viper.BindEnv("etcd.host", "ETCD_ENDPOINT")
	viper.BindEnv("loki.host", "LOKI_ENDPOINT")
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")

	viper.SetDefault("api.host", "localhost")
	viper.SetDefault("api.port", "8080")
	viper.SetDefault("etcd.host", "localhost:2379")
	viper.SetDefault("loki.host", "localhost:3100")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")

}
