/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crane",
	Short: "The Underlying Engine Behind Cranom Platform",
	Long: `The Underlying Engine Behind Cranom Platform
	
	Learn More at: https://www.cranom.tech/plaform-tools/crane-engine
	`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Crane Orchastrator Engine")
		fmt.Println("Use cranom --help to see available commands")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(KubeContollerCmd)
	rootCmd.AddCommand(APIServerCmd)
	rootCmd.AddCommand(DockerDriverCmd)
	rootCmd.AddCommand(NginxDriverCMD)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/crane/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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
