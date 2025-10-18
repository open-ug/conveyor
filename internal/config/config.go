package config

type Config struct {
	API  APIConfig  `json:"api" yaml:"api" mapstructure:"api"`
	NATS NATSConfig `json:"nats" yaml:"nats" mapstructure:"nats"`
	TLS  TLSConfig  `json:"tls" yaml:"tls" mapstructure:"tls"`
}

type APIConfig struct {
	Port        int    `json:"port" yaml:"port" mapstructure:"port"`
	AuthEnabled bool   `json:"auth_enabled" yaml:"auth_enabled" mapstructure:"auth_enabled"`
	Data        string `json:"data" yaml:"data" mapstructure:"data"`
}

type NATSConfig struct {
	Port int `json:"port" yaml:"port" mapstructure:"port"`
}

type TLSConfig struct {
	CA   string `json:"ca" yaml:"ca" mapstructure:"ca"`
	Key  string `json:"key" yaml:"key" mapstructure:"key"`
	Cert string `json:"cert" yaml:"cert" mapstructure:"cert"`
}
