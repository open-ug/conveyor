package sampledriver

import (
	//"encoding/json"

	"fmt"
	"log"
	"time"

	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
	"github.com/open-ug/conveyor/pkg/types"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) types.DriverResult {

	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("Sample Driver Reconciling::: EVENT: %v PAYLOAD: %v", event, payload)
	time.Sleep(3 * time.Second)

	return types.DriverResult{
		Success: true,
		Message: "Sample Driver Reconciled Successfully",
	}
}

func Listen(
	name string,
	resources []string,
) {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
		Name:      name,
		Resources: resources,
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
