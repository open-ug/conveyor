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
		app, ok := appMsg.Payload.(craneTypes.Application)
		if !ok {
			return fmt.Errorf("error converting payload to Application type")
		}
		err = CreateContainer(dockerClient, &app)
		if err != nil {
			return fmt.Errorf("error creating container: %v", err)
		}
		return nil
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
