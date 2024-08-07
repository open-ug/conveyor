package driverruntime

import (
	"context"
	"fmt"

	apiServer "crane.cloud.cranom.tech/cmd/api"
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
	Reconcile func(resourceId string) error
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
	// The driver manager will run the driver's reconcile function
	// in a loop
	for {
		// Get the resource id from the message queue
		pubsub := d.RedisClient.Subscribe(context.Background(), "application")

		ch := pubsub.Channel()

		for msg := range ch {
			err := d.Driver.Reconcile(msg.Payload)
			if err != nil {
				fmt.Println("Error reconciling resource: ", err)
				return err
			}
		}
	}
}
