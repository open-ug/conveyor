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
	internals "github.com/open-ug/conveyor/internal/shared"
	utils "github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/driver-runtime/log"
	types "github.com/open-ug/conveyor/pkg/types"
)

type DriverManager struct {
	// The driver manager is responsible for managing the drivers
	// and the driver lifecycle.

	NatsContext *internals.NatsContext

	Driver *Driver

	// an array of events that the driver manager will listen to
	// and reconcile
	Events []string
}

type Driver struct {
	// The driver is responsible for managing the driver
	Reconcile func(message string, event string, runID string, logger *log.DriverLogger) error

	Name string

	Resources []string
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
		NatsContext: natsCon,
		Driver:      driver,
		Events:      events,
	}, nil
}

func (d *DriverManager) Run() error {

	// The driver manager will run the driver's reconcile function
	// in a loop

	// Get the resource from the message queue
	randomStr, serr := utils.GenerateRandomShortStr()
	if serr != nil {
		color.Red("Error Occured while generating random string: %v", serr)
	}
	fmt.Println(randomStr)

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
	}

	// CONSUMER
	consumer.Consume(func(msg jetstream.Msg) {
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

	fmt.Println("Driver Manager is running for driver: ", d.Driver.Name)

	select {}
}
