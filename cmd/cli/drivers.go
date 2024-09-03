/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package cli

import (
	dockerDriver "crane.cloud.cranom.tech/cmd/docker-driver"
	nginxDriver "crane.cloud.cranom.tech/cmd/nginx-driver"
	"github.com/spf13/cobra"
)

var DockerDriverCmd = &cobra.Command{
	Use:   "docker-driver",
	Short: "Start the Crane Docker Driver",
	Long: `Start the Cranom API Server to interact with the Cranom Platform.

Learn More at: https://www.cranom.tech/plaform-tools/crane
`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerDriver.Listen()
	},
}

var NginxDriverCMD = &cobra.Command{
	Use:   "nginx-driver",
	Short: "Start the Crane Nginx Driver",
	Long: `Start the Cranom API Server to interact with the Cranom Platform.

Learn More at: https://www.cranom.tech/plaform-tools/crane
`,
	Run: func(cmd *cobra.Command, args []string) {
		nginxDriver.Listen()
	},
}
