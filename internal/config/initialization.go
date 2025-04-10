package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var CfgFile string

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {

		// /etc/crane/config.yaml
		viper.AddConfigPath("/etc/crane")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file: ", err)
	}

}
