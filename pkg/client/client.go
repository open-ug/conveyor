package client

import (
	config "conveyor.cloud.cranom.tech/internal/config"
	"github.com/spf13/viper"
)

type Client struct {
	APIHost string
	APIPort string
}

func NewClient() *Client {
	// Initialize the config
	config.InitConfig()
	// Create a new client instance
	return &Client{
		APIHost: viper.GetString("api.host"),
		APIPort: viper.GetString("api.port"),
	}
}
