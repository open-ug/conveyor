/*
Copyright © 2024 Conveyor CI Contributors
*/
package cli

import (
	"fmt"
	"os"

	config "github.com/open-ug/conveyor/internal/config"
	"github.com/spf13/cobra"
)

var Version string = "development"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "conveyor",
	Short: "Conveyor CI",
	Long: `The lightweight, distributed CI/CD engine built for platform developers who demand simplicity without compromise.
	`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Conveyor CI " + Version)
		fmt.Println("Run 'conveyor --help' to see available commands")
		fmt.Println("Visit https://conveyor.open.ug for more information")
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
	rootCmd.AddCommand(SampleDriverCmd)
	rootCmd.AddCommand(initCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}