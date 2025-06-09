package logger

import (
	//"encoding/json"

	"fmt"
	"log"

	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) error {

	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("Webhook Driver Reconciling::: EVENT: %v PAYLOAD: %v", event, payload)

	return nil
}

func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
		Name:      "logger",
		Resources: []string{"job"},
	}

	driverManager, err := runtime.NewDriverManager(driver, []string{"*"})

	if err != nil {
		fmt.Println("Error creating driver manager: ", err)
		return
	}

	err = driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
