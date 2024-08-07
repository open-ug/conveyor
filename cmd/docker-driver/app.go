package dockerdriver

import (
	"fmt"

	runtime "crane.cloud.cranom.tech/cmd/driver-runtime"
)

func Reconcile(payload string) error {
	/* dockerClient, err := GetDockerClient()
	if err != nil {
		return fmt.Errorf("Error getting docker client: %v", err)
	} */
	fmt.Println("Reconciling resource: ", payload)

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
