package certbotdriver

import (
	"encoding/json"
	"fmt"

	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string) error {
	fmt.Println("CERTBOT_D: Reconcyling payload: " + payload)

	if event != "application" {
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
