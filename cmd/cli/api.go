/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package cli

import (
	apiServer "github.com/open-ug/conveyor/cmd/api"
	"github.com/spf13/cobra"
)

var APIServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Start the Conveyor API Server",
	Long: `Start the Conveyor API Server

`,
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("port").Value.String()
		if port == "" {
			port = "3000"
		}
		apiServer.StartServer(port)
	},
}

func init() {
	APIServerCmd.Flags().StringP("port", "p", "", "Port to run the API Server on (default: 3000)")
}
