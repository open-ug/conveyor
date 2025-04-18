/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package dockerdriver

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
	log.Printf("Docker Driver Reconciling: %v", payload)

	dockerClient, err := GetDockerClient()
	if err != nil {
		return fmt.Errorf("error getting docker client: %v", err)
	}
	if event == "application" || event == "buildpack-create-complete" {
		var appMsg craneTypes.ApplicationMsg
		err = json.Unmarshal([]byte(payload), &appMsg)
		if err != nil {
			return fmt.Errorf("error unmarshalling application message: %v", err)
		}

		if event == "application" && appMsg.Payload.Spec.Source.Type == "git" {
			return nil
		}

		if appMsg.Action == "create" {
			// Create Action
			// Creating a new Container and Run It
			app := appMsg.Payload
			err := CreateContainer(dockerClient, &app)
			if err != nil {
				return fmt.Errorf("error creating container: %v", err)
			}
			// Broadcast Complete Message
			runtime.BroadCastMessage(
				craneTypes.DriverMessage{
					Event:   "docker-create-complete",
					Payload: payload,
					RunID:   runID,
				},
			)
			// start the container
			serr := StartContainer(dockerClient, &app)
			if serr != nil {
				return fmt.Errorf("error starting container: %v", serr)
			}
			// Broadcast Complete Message
			runtime.BroadCastMessage(
				craneTypes.DriverMessage{
					Event:   "docker-start-complete",
					Payload: payload,
					RunID:   runID,
				},
			)

			// Broadcast Complete Message
			runtime.BroadCastMessage(
				craneTypes.DriverMessage{
					Event:   "docker-build-complete",
					Payload: payload,
				},
			)
			return nil
		} else if appMsg.Action == "delete" {
			app := appMsg.Payload
			err = DeleteContainer(dockerClient, &app)
			if err != nil {
				return fmt.Errorf("error deleting container: %v", err)
			}
			return nil
		} else if appMsg.Action == "update" {
			app := appMsg.Payload
			err = UpdateContainer(dockerClient, &app)
			if err != nil {
				return fmt.Errorf("error updating container: %v", err)
			}
		} else if appMsg.Action == "start" {
			app := appMsg.Payload
			err = StartContainer(dockerClient, &app)
			if err != nil {
				return fmt.Errorf("error starting container: %v", err)
			}
		} else if appMsg.Action == "stop" {
			app := appMsg.Payload
			err = StopContainer(dockerClient, &app)
			if err != nil {
				return fmt.Errorf("error stopping container: %v", err)
			}

			// Broadcast Stop Message
			runtime.BroadCastMessage(
				craneTypes.DriverMessage{
					Event:   "docker-stop-complete",
					Payload: payload,
				},
			)
		}
	}

	return nil
}

// Listen listens for events from the driver runtime
func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
	}

	driverManager, err := runtime.NewDriverManager(driver, []string{"application", "buildpack-create-complete"})

	if err != nil {
		fmt.Println("Error creating driver manager: ", err)
		return
	}

	err = driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
