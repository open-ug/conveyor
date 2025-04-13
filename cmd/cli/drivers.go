/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package cli

import (
	"fmt"

	buildpacksDriver "github.com/open-ug/conveyor/cmd/buildpacks-driver"
	dockerDriver "github.com/open-ug/conveyor/cmd/docker-driver"
	"github.com/open-ug/conveyor/cmd/logger"
	nginxDriver "github.com/open-ug/conveyor/cmd/nginx-driver"
	webhookDriver "github.com/open-ug/conveyor/cmd/webhook-driver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var BuildPacksDriverCmd = &cobra.Command{
	Use:   "bp-driver",
	Short: "Start the Crane Buildpacks Driver",
	Long: `Start the Cranom API Server to interact with the Cranom Platform.

Learn More at: https://www.cranom.tech/plaform-tools/crane
`,
	Run: func(cmd *cobra.Command, args []string) {
		buildpacksDriver.Listen()
	},
}

var WebhookDriverCmd = &cobra.Command{
	Use:   "webhook-driver",
	Short: "Start the Crane Webhook Driver",
	Long: `Start the Cranom API Server to interact with the Cranom Platform.

Learn More at: https://www.cranom.tech/plaform-tools/crane
`,
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("endpoint").Value.String()
		if port == "" {
			fmt.Println("Please Specify --endpoint flag")
			return
		}
		viper.Set("webhook.endpoint", port)
		webhookDriver.Listen()
	},
}

func init() {
	WebhookDriverCmd.Flags().StringP("endpoint", "e", "", "Webhook endpoint where events will be published")
}
