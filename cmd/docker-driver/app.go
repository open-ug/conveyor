/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package dockerdriver

import (
	"encoding/json"
	"fmt"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	runtime "crane.cloud.cranom.tech/cmd/driver-runtime"
)

func Reconcile(payload string) error {
	dockerClient, err := GetDockerClient()
	if err != nil {
		return fmt.Errorf("error getting docker client: %v", err)
	}
	var appMsg craneTypes.ApplicationMsg
	err = json.Unmarshal([]byte(payload), &appMsg)
	if err != nil {
		return fmt.Errorf("error unmarshalling application message: %v", err)
	}

	if appMsg.Action == "create" {
		app := appMsg.Payload
		err := CreateContainer(dockerClient, &app)
		if err != nil {
			return fmt.Errorf("error creating container: %v", err)
		}
		// start the container
		serr := StartContainer(dockerClient, &app)
		if serr != nil {
			return fmt.Errorf("error starting container: %v", serr)
		}
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
	}

	return nil
}

// Listen listens for events from the driver runtime
func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
	}

	driverManager := runtime.NewDriverManager(driver)

	err := driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
