/*
Copyright © 2024 Conveyor CI Contributors
*/
package cli

import (
	"fmt"

	apiServer "github.com/open-ug/conveyor/cmd/api"
	sampledriver "github.com/open-ug/conveyor/cmd/sample-driver"
	config "github.com/open-ug/conveyor/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var APIServerCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the Conveyor API Server",
	Long:  `Start the Conveyor API Server. This will start the server on the specified port (default: 8080) and connect to the embedded etcd instance.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Attempt to load the config
		if err := config.LoadConfig(); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("port").Value.String()
		if port == "" {
			if viper.IsSet("api.port") {
				port = viper.GetString("api.port")
			} else {
				port = "8080"
			}
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
