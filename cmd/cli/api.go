/*
Copyright © 2024 Conveyor CI Contributors
*/
package cli

import (
	apiServer "github.com/open-ug/conveyor/cmd/api"
	sampledriver "github.com/open-ug/conveyor/cmd/sample-driver"
	"github.com/spf13/cobra"
)

var APIServerCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the Conveyor Service",
	Long: `Start the Conveyor Service

`,
	Run: func(cmd *cobra.Command, args []string) {
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
	Long:  `Start the Sample Driver`,
	Run: func(cmd *cobra.Command, args []string) {
		sampledriver.Listen()
	},
}

func init() {
	APIServerCmd.Flags().StringP("port", "p", "", "Port to run the API Server on (default: 3000)")
}
