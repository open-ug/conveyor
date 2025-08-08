/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package driverruntime

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/nats-io/nats.go/jetstream"
	config "github.com/open-ug/conveyor/internal/config"
	utils "github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/driver-runtime/log"
	types "github.com/open-ug/conveyor/pkg/types"
)

type DriverManager struct {
	// The driver manager is responsible for managing the drivers
	// and the driver lifecycle.

	NatsContext *utils.NatsContext

	Driver *Driver

	// an array of events that the driver manager will listen to
	// and reconcile
	Events []string
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

	natsCon := utils.NewNatsConn()

	return &DriverManager{
		NatsContext: natsCon,
		Driver:      driver,
		Events:      events,
	}, nil
}

func (d *DriverManager) Run() error {

	// Resources
	var filterSubjects []string
	for _, resource := range d.Driver.Resources {
		filterSubjects = append(filterSubjects, "resources."+resource)
	}

	consumer, err := d.NatsContext.JetStream.CreateOrUpdateConsumer(context.Background(), "messages", jetstream.ConsumerConfig{
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
			}, d.NatsContext.NatsCon)

			err = d.Driver.Reconcile(message.Payload, message.Event, message.RunID, logger)
			if err != nil {
				fmt.Println("Error reconciling resource: ", err)
				//return err
			}
		}
	})

	if err != nil {
		color.Red("Error Occured while consuming messages: %v", err)
		return err
	}

	fmt.Println("Driver Manager is running for driver: ", d.Driver.Name)

	select {}
}
