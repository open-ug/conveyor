package certbotdriver

import (
	"encoding/json"
	"fmt"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	runtime "crane.cloud.cranom.tech/cmd/driver-runtime"
)

// Listen for messages from the runtime
func Reconcile(payload string) error {
	fmt.Println("CERTBOT_D: Reconcyling payload: " + payload)

	// Unmarshal the payload
	var appMsg craneTypes.ApplicationMsg
	err := json.Unmarshal([]byte(payload), &appMsg)
	if err != nil {
		return fmt.Errorf("error unmarshalling application message: %v", err)
	}

	// Handle the action
	if appMsg.Action == "create" {
		CreateCertBotConfig(appMsg.Payload)
	} else if appMsg.Action == "delete" {
		DeleteCertBotConfig(appMsg.Payload)
	} else if appMsg.Action == "update" {
		//
	} else if appMsg.Action == "start" {
		//
	} else if appMsg.Action == "stop" {
		//
	}

	return nil
}

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
