package cli

import (
	apiServer "crane.cloud.cranom.tech/cmd/api"
	"github.com/spf13/cobra"
)

var APIServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Start the Cranom API Server",
	Long: `Start the Cranom API Server to interact with the Cranom Platform.

Learn More at: https://www.cranom.tech/plaform-tools/crane
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
