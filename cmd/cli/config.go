/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package cli

type Config struct {
	Api ApiConfig `mapstructure:"api"`
}

type ApiConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
