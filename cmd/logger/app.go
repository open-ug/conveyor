package logger

import (
	//"encoding/json"

	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) error {

	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("Webhook Driver Reconciling::: EVENT: %v PAYLOAD: %v", event, payload)

	captureAndStreamLogs(PrintLog, func(line string) {
		logger.Log(map[string]string{
			"event": event,
		}, line)
	})

	return nil
}

func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
		Name:      "logger",
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

// A function that prints out 10 lines of the log using fmt
func PrintLog() {
	for i := 0; i < 10; i++ {
		fmt.Println("This is line", i+1)
		time.Sleep(1 * time.Second)
	}
}

func captureAndStreamLogs(f func(), streamFn func(string)) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Start a goroutine to read and stream each line
	done := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			streamFn(line)
		}
		close(done)
	}()

	f()

	// Close writer and restore stdout
	w.Close()
	os.Stdout = origStdout
	<-done
}
