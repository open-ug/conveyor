/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package driverruntime

import (
	"context"
	"fmt"
	"time"

	apiServer "crane.cloud.cranom.tech/cmd/api"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
)

type DriverManager struct {
	// The driver manager is responsible for managing the drivers
	// and the driver lifecycle.

	RedisClient *redis.Client

	Driver *Driver
}

type Driver struct {
	// The driver is responsible for managing the driver
	Reconcile func(message string) error
}

func NewDriverManager(
	driver *Driver,
) *DriverManager {
	rdb := apiServer.NewRedisClient()

	return &DriverManager{
		RedisClient: rdb,
		Driver:      driver,
	}
}

func (d *DriverManager) Run() error {
	err := WatchContainerStart("cranom-redis")
	if err != nil {
		return err
	}
	// The driver manager will run the driver's reconcile function
	// in a loop
	for {
		// Get the resource from the message queue
		pubsub := d.RedisClient.Subscribe(context.Background(), "application")

		ch := pubsub.Channel()

		for msg := range ch {
			err := d.Driver.Reconcile(msg.Payload)
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
			//color.Green("System Component is running")
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
