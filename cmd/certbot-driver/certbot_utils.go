package certbotdriver

import (
	"fmt"
	"os/exec"

	craneTypes "conveyor.cloud.cranom.tech/pkg/types"
)

func CreateCertBotConfig(app craneTypes.Application) {
	fmt.Println("Creating CertBot Config")

	// loop through app.Ports and create a certbot config for each port that has a SSL toTrue
	for _, port := range app.Spec.Ports {
		if port.SSL {
			// create the certbot config
			cmd := exec.Command("certbot", "--nginx", "-d", port.Domain)
			err := cmd.Run()
			if err != nil {
				fmt.Println("Failed to create certbot config")
			}
		}
	}
}

func DeleteCertBotConfig(app craneTypes.Application) {
	fmt.Println("Deleting CertBot Config")

	// loop through app.Ports and delete a certbot config for each port that has a SSL toTrue
	for _, port := range app.Spec.Ports {
		if port.SSL {
			// delete the certbot config
			cmd := exec.Command("certbot", "delete", "--nginx", "-d", port.Domain)
			err := cmd.Run()
			if err != nil {
				fmt.Println("Failed to delete certbot config")
			}
		}
	}
}
