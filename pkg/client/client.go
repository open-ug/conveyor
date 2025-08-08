package client

import (
	"github.com/go-resty/resty/v2"
	config "github.com/open-ug/conveyor/internal/config"
	"github.com/spf13/viper"
)

type Client struct {
	HTTPClient *resty.Client
}

func NewClient() *Client {
	// Initialize the config
	config.InitConfig()

	// Initialize HTTP client
	client := resty.New()
	client.SetBaseURL(viper.GetString("api.host"))
	client.SetHeader("Content-Type", "application/json")

	// Create a new client instance
	return &Client{
		HTTPClient: client,
	}
}

func (c *Client) GetAPIURL() string {
	// Return the API URL
	return viper.GetString("api.host")
}
