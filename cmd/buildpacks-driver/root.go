/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package buildpacksdriver

import (
	"fmt"

	runtime "crane.cloud.cranom.tech/cmd/driver-runtime"
)

func Reconcile(payload string, event string) error {

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
