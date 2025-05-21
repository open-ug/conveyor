/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package driverruntime

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/nats-io/nats.go"
	config "github.com/open-ug/conveyor/internal/config"
	internals "github.com/open-ug/conveyor/internal/shared"
	"github.com/open-ug/conveyor/pkg/driver-runtime/log"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
)

type DriverManager struct {
	// The driver manager is responsible for managing the drivers
	// and the driver lifecycle.

	NatsCon *nats.Conn

	Driver *Driver

	// an array of events that the driver manager will listen to
	// and reconcile
	Events []string
}

type Driver struct {
	// The driver is responsible for managing the driver
	Reconcile func(message string, event string, runID string, logger *log.DriverLogger) error

	Name string
}

// validate the driver
func (d *Driver) Validate() error {
	if d.Reconcile == nil {
		return fmt.Errorf("driver reconcile function is not set")
	}
	if d.Name == "" {
		return fmt.Errorf("driver name is not set")
	}
	return nil
}

func NewDriverManager(
	driver *Driver,
	events []string,
) (*DriverManager, error) {
	// Load the configuration
	config.InitConfig()

	// Validate the driver
	err := driver.Validate()
	if err != nil {
		color.Red("Error Occured while validating driver: %v", err)
		return nil, err
	}

	natsCon := internals.NewNatsConn()

	return &DriverManager{
		NatsCon: natsCon,
		Driver:  driver,
		Events:  events,
	}, nil
}

func (d *DriverManager) Run() error {
	// The driver manager will run the driver's reconcile function
	// in a loop

	// Get the resource from the message queue
	_, err := d.NatsCon.Subscribe("application", func(msg *nats.Msg) {

		var message craneTypes.DriverMessage
		err := json.Unmarshal([]byte(msg.Data), &message)
		if err != nil {
			fmt.Println("Error unmarshalling message: ", err)
			//return err
		}

		// check if the event is in the list of events
		// that the driver manager is listening to
		var eventFound bool = false
		for _, event := range d.Events {
			if event == message.Event {
				eventFound = true
				break
			} else if event == "*" {
				eventFound = true
				break
			}
		}

		if eventFound {
			logger := log.NewDriverLogger(d.Driver.Name, map[string]string{
				"event":  message.Event,
				"id":     message.ID,
				"run_id": message.RunID,
			}, d.NatsCon)

			err = d.Driver.Reconcile(message.Payload, message.Event, message.RunID, logger)
			if err != nil {
				fmt.Println("Error reconciling resource: ", err)
				//return err
			}
		}

	})
	if err != nil {
		color.Red("Error Occured while subscribing to NATS channel: %v", err)
	}
	fmt.Println("Subscribed to NATS channel")

	select {}
}

func GetDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error creating docker client: %v", err)
	}
	return cli, nil
}

// a function that watches for docker container to start. it keeps checking the container status until it is running
func WatchContainerStart(containerID string) error {
	cli, err := GetDockerClient()
	if err != nil {
		return err
	}

	for {
		fmt.Println("Waiting for System Component to start")
		inspect, err := cli.ContainerInspect(context.Background(), containerID)
		if err != nil {
			color.Red("System Component failed to start")
			return fmt.Errorf("error inspecting container: %v", err)
		}

		if inspect.State.Status == "running" {
			color.Green("System Component is running")
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
