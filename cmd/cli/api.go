/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package cli

import (
	apiServer "github.com/open-ug/conveyor/cmd/api"
	"github.com/open-ug/conveyor/cmd/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var APIServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Start the Conveyor API Server",
	Long: `Start the Conveyor API Server

`,
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("port").Value.String()
		if port == "" {
			port = viper.GetString("api.port")
		}
		apiServer.StartServer(port)
	},
}

func init() {
	APIServerCmd.Flags().StringP("port", "p", "", "Port to run the API Server on (default: 3000)")
}

var CertBotDriverCmd = &cobra.Command{
	Use:   "certbot-driver",
	Short: "Start the Crane CertBot Driver",
	Long: `Start the Cranom API Server to interact with the Cranom Platform.

Learn More at: https://www.cranom.tech/plaform-tools/crane
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Listen()
	},
}
