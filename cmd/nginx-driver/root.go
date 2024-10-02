package nginxdriver

import (
	"encoding/json"
	"fmt"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	runtime "crane.cloud.cranom.tech/cmd/driver-runtime"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string) error {
	fmt.Println("NGINX_D: Reconcyling payload: " + payload)

	if event != "application" {
		// Unmarshal the payload
		var appMsg craneTypes.ApplicationMsg
		err := json.Unmarshal([]byte(payload), &appMsg)
		if err != nil {
			return fmt.Errorf("error unmarshalling application message: %v", err)
		}

		// Handle the action
		if appMsg.Action == "create" {
			CreateNginxConfig(appMsg.Payload)
		} else if appMsg.Action == "delete" {
			DeleteNginxConfig(appMsg.Payload)
		} else if appMsg.Action == "update" {
			UpdateNginxConfig(appMsg.Payload)
		} else if appMsg.Action == "start" {
			//
		} else if appMsg.Action == "stop" {
			//
		}
	}

	return nil
}

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
