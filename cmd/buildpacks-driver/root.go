/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package buildpacksdriver

import (
	"encoding/json"
	"fmt"
	"log"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	runtime "crane.cloud.cranom.tech/cmd/driver-runtime"
)

func Reconcile(payload string, event string) error {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("Buiild Pcks Reconciling: %v", payload)

	if event == "application" {
		var appMsg craneTypes.ApplicationMsg
		err := json.Unmarshal([]byte(payload), &appMsg)
		if err != nil {
			return fmt.Errorf("error unmarshalling application message: %v", err)
		}

		if appMsg.Action == "create" {
			app := appMsg.Payload

			if app.Spec.Source.Type == "git" {
				err := CreateBuildpacksImage(&app)
				if err != nil {
					return fmt.Errorf("error creating buildpacks image: %v", err)
				}

				runtime.BroadCastMessage(
					craneTypes.DriverMessage{
						Event:   "buildpack-create-complete",
						Payload: payload,
					},
				)
			}
			return nil
		}
	}

	return nil
}

// Listen listens for events from the driver runtime
func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
	}

	driverManager := runtime.NewDriverManager(driver, []string{"application"})

	err := driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
