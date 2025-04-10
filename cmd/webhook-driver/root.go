package webhookdriver

import (
	//"encoding/json"
	"fmt"
	"log"

	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string) error {

	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("Webhook Driver Reconciling::: EVENT: %v PAYLOAD: %v", event, payload)

	dm := craneTypes.DriverMessage{
		Payload: payload,
		Event:   event,
	}

	PostMessage(dm)

	return nil
}

func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
	}

	driverManager := runtime.NewDriverManager(driver, []string{"*"})

	err := driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
