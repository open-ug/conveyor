/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package cli

import (
	driver "crane.cloud.cranom.tech/cmd/docker-driver"
	"github.com/spf13/cobra"
)

var DockerDriverCmd = &cobra.Command{
	Use:   "docker-driver",
	Short: "Start the Cranom API Server",
	Long: `Start the Cranom API Server to interact with the Cranom Platform.

Learn More at: https://www.cranom.tech/plaform-tools/crane
`,
	Run: func(cmd *cobra.Command, args []string) {
		driver.Listen()
	},
}
