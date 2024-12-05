/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package driverruntime

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	apiServer "crane.cloud.cranom.tech/cmd/api"
	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
)

type DriverManager struct {
	// The driver manager is responsible for managing the drivers
	// and the driver lifecycle.

	RedisClient *redis.Client

	Driver *Driver

	// an array of events that the driver manager will listen to
	// and reconcile
	Events []string
}

type Driver struct {
	// The driver is responsible for managing the driver
	Reconcile func(message string, event string) error
}

func NewDriverManager(
	driver *Driver,
	events []string,
) *DriverManager {
	rdb := apiServer.NewRedisClient()

	return &DriverManager{
		RedisClient: rdb,
		Driver:      driver,
		Events:      events,
	}
}

func (d *DriverManager) Run() error {
	err := WatchContainerStart("crane-redis")
	if err != nil {
		color.Red("Error Occured while waiting for Redis to start: %v", err)
		return err
	}
	// The driver manager will run the driver's reconcile function
	// in a loop
	for {
		// Get the resource from the message queue
		pubsub := d.RedisClient.Subscribe(context.Background(), "application")

		ch := pubsub.Channel()

		for msg := range ch {

			var message craneTypes.DriverMessage
			err := json.Unmarshal([]byte(msg.Payload), &message)
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
				}
			}

			if !eventFound {
				continue
			}

			err = d.Driver.Reconcile(message.Payload, message.Event)
			if err != nil {
				fmt.Println("Error reconciling resource: ", err)
				//return err
			}
		}
	}
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
