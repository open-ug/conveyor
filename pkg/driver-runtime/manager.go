/*
The Driver Runtime is responsible for managing the lifecycle of drivers and their interactions with the system.

Copyright © 2024 Conveyor CI Contributors
*/
package driverruntime

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/engine"
	"github.com/open-ug/conveyor/pkg/driver-runtime/log"
	types "github.com/open-ug/conveyor/pkg/types"
)

type DriverManager struct {
	// The driver manager is responsible for managing the drivers
	// and the driver lifecycle.

	Driver *Driver

	// an array of events that the driver manager will listen to
	// and reconcile
	Events []string
}

func NewDriverManager(
	driver *Driver,
	events []string,
) (*DriverManager, error) {
	// Validate the driver
	err := driver.Validate()
	if err != nil {
		color.Red("Error Occured while validating driver: %v", err)
		return nil, err
	}

	return &DriverManager{
		Driver: driver,
		Events: events,
	}, nil
}

func (d *DriverManager) Run() error {
	// Setup NATS JetStream
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		color.Red("Error Occured while connecting to NATS: %v", err)
		return err
	}
	js, err := jetstream.New(nc)
	if err != nil {
		color.Red("Error Occured while creating JetStream context: %v", err)
		return err
	}

	// Resources
	var filterSubjects []string
	for _, resource := range d.Driver.Resources {
		filterSubjects = append(filterSubjects, "resources."+resource)
	}

	consumer, err := js.CreateOrUpdateConsumer(context.Background(), "messages", jetstream.ConsumerConfig{
		Name:           d.Driver.Name,
		FilterSubjects: filterSubjects,
		AckPolicy:      jetstream.AckExplicitPolicy,
		MaxAckPending:  1,
	})

	if err != nil {
		color.Red("Error Occured while subscribing to NATS channel: %v", err)
		return err
	}

	// CONSUMER
	_, err = consumer.Consume(func(msg jetstream.Msg) {
		fmt.Printf("Received message on subject: %s\n", msg.Subject())
		msg.Ack()
		data := msg.Data()
		var message types.DriverMessage
		err := json.Unmarshal([]byte(data), &message)
		if err != nil {
			color.Red("Error Occured while unmarshalling message: %v", err)
			return
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
			}, nc)

			result := d.Driver.Reconcile(message.Payload, message.Event, message.RunID, logger)

			driverevent := engine.DriverResultEvent{
				Success: result.Success,
				Message: result.Message,
				Driver:  d.Driver.Name,
			}

			var resource types.Resource
			err := json.Unmarshal([]byte(message.Payload), &resource)
			if err != nil {
				color.Red("Error Occured while unmarshalling resource: %v", err)
				return
			}

			driverevent.PublishEvent(message.RunID, resource, js)

		}
	})

	if err != nil {
		color.Red("Error Occured while consuming messages: %v", err)
		return err
	}

	fmt.Println("Driver Manager is running for driver: ", d.Driver.Name)

	select {}
}
