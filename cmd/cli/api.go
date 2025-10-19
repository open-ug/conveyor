/*
Copyright © 2024 Conveyor CI Contributors
*/
package cli

import (
	"fmt"
	"log"

	apiServer "github.com/open-ug/conveyor/cmd/api"
	sampledriver "github.com/open-ug/conveyor/cmd/sample-driver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var APIServerCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the Conveyor Service",
	Long: `Start the Conveyor Service

`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("api.data") {
			log.Fatal("ERROR: API data directory is not set in configuration")
		}
		port := cmd.Flag("port").Value.String()
		if port == "" {
			port = "8080"
		}
		apiServer.StartServer(port)
	},
}

var SampleDriverCmd = &cobra.Command{
	Use:   "sampledriver",
	Short: "Start the Sample Driver",
	Long:  `Start the Sample Driver for testing purposes, You can specify the name and resources the driver will manage.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flag("name").Value.String()
		resources, err := cmd.Flags().GetStringSlice("resources")
		if err != nil {
			fmt.Println("Error getting resources flag: ", err)
			fmt.Println("Defaulting to 'pipe' resource")
			resources = []string{"pipe"}
		}
		if name == "" {
			name = "sampledriver"
		}
		if len(resources) == 0 {
			resources = []string{"pipe"}
		}
		sampledriver.Listen(name, resources)
	},
}

func init() {
	APIServerCmd.Flags().StringP("port", "p", "", "Port to run the API Server on (default: 3000)")
	SampleDriverCmd.Flags().StringP("name", "n", "sampledriver", "Name of the driver")
	SampleDriverCmd.Flags().StringSliceP("resources", "r", []string{"pipe"}, "Resources the driver will manage")
}
