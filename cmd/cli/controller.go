/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package cli

import (
	"fmt"

	craneController "crane.cloud.cranom.tech/cmd/controller"

	"github.com/spf13/cobra"
)

var KubeContollerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Start the Cranom Controller",
	Long: `Start the Cranom Controller to manage Cranom Applications in a Kubernetes Cluster.

Learn More at: https://www.cranom.tech/plaform-tools/controller
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Cranom Controller")
		craneController.RunController()
	},
}
