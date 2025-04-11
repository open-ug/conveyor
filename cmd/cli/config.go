/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package cli

type Config struct {
	Api   ApiConfig   `mapstructure:"api" yaml:"api"`
	DB    DBConfig    `mapstructure:"db" yaml:"db"`
	Redis RedisConfig `mapstructure:"redis" yaml:"redis"`
}

type ApiConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port string `mapstructure:"port" yaml:"port"`
}

type DBConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port string `mapstructure:"port" yaml:"port"`
	User string `mapstructure:"user" yaml:"user"`
	Pass string `mapstructure:"pass" yaml:"pass"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port string `mapstructure:"port" yaml:"port"`
}
