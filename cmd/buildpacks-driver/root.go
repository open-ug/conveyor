/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package buildpacksdriver

import (
	"encoding/json"
	"fmt"
	"log"

	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
)

func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) error {
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
						RunID:   runID,
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
