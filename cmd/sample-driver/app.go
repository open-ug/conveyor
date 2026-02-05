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

	/// simulate doing some work by looping 20 times
	for i := 0; i < 20; i++ {
		log.Printf("Sample Driver Working... %d/20", i+1)
		logger.Log(map[string]string{"step": fmt.Sprintf("%d", i+1)}, fmt.Sprintf("Sample Driver Working... %d/20", i+1))
		time.Sleep(2 * time.Second)
	}

	return types.DriverResult{
		Success: true,
		Message: "Sample Driver Reconciled Successfully",
		Data:    nil,
	}
}

func Listen(
	name string,
	resources []string,
) {

	/**
	 * UNCOMMENT TO ENABLE AUTHENTICATION

	// LOAD TLS CERTS in  ./client-cert.pem and ./client-key.pem and ./ca.pem
	// and pass them to the runtime client options
	cert, err := os.ReadFile("./client-cert.pem")
	if err != nil {
		fmt.Println("Error reading client cert: ", err)
		return
	}

	key, err := os.ReadFile("./client-key.pem")
	if err != nil {
		fmt.Println("Error reading client key: ", err)
		return
	}

	rootCA, err := os.ReadFile("./ca.pem")
	if err != nil {
		fmt.Println("Error reading root CA: ", err)
		return
	}
	*/

	client, err := runtime.NewClient("http://localhost:8080", "tls://localhost:4222", runtime.ConfigOptions{
		//AuthEnabled: true,
		//Cert:        cert,
		//Key:         key,
		//RootCA:      rootCA,
	})

	driver := &runtime.Driver{
		Reconcile: Reconcile,
		Name:      name,
		Resources: resources,
	}

	driverManager, err := client.NewDriverManager(driver, []string{"*"})

	if err != nil {
		fmt.Println("Error creating driver manager: ", err)
		return
	}

	err = driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
