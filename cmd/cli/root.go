/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package cli

import (
	"fmt"
	"os"

	config "github.com/open-ug/conveyor/internal/config"
	"github.com/spf13/cobra"
)

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
	cobra.OnInitialize(config.InitConfig)

	rootCmd.AddCommand(APIServerCmd)
	rootCmd.AddCommand(DockerDriverCmd)
	rootCmd.AddCommand(NginxDriverCMD)
	rootCmd.AddCommand(CertBotDriverCmd)
	rootCmd.AddCommand(WebhookDriverCmd)
	rootCmd.AddCommand(BuildPacksDriverCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is /etc/crane/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
