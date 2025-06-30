/*
Copyright © 2024 Beingana Jim Junior and Contributors
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
	Use:   "conveyor",
	Short: "Conveyor CLI",
	Long: `Conveyor CLI is a command line interface for managing and interacting with the Conveyor platform.
	`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Conveyor CLI")
		fmt.Println("Use conveyor --help to see available commands")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.AddCommand(APIServerCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
