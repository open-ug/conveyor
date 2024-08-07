package dockerdriver

import (
	"fmt"

	runtime "crane.cloud.cranom.tech/cmd/driver-runtime"
)

func Listen() {
	driver := &runtime.Driver{
		Reconcile: func(resourceId string) error {
			fmt.Println("Reconciling resource: ", resourceId)
			return nil
		},
	}

	driverManager := runtime.NewDriverManager(driver)

	err := driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
