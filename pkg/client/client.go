package client

import (
	config "github.com/open-ug/conveyor/internal/config"
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

func (c *Client) GetAPIURL() string {
	// Return the API URL
	return c.APIHost + ":" + c.APIPort
}
