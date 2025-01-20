/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package cli

import (
	"fmt"

	certbotDriver "crane.cloud.cranom.tech/cmd/certbot-driver"
	dockerDriver "crane.cloud.cranom.tech/cmd/docker-driver"
	nginxDriver "crane.cloud.cranom.tech/cmd/nginx-driver"
	webhookDriver "crane.cloud.cranom.tech/cmd/webhook-driver"
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
		certbotDriver.Listen()
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
