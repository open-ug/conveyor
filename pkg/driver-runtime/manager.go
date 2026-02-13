/*
The Driver Runtime is responsible for managing the lifecycle of drivers and their interactions with the system.

Copyright © 2024 - Present Conveyor CI Contributors
*/
package driverruntime

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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

	// The API client to interact with the Conveyor API
	Client *Client
}

// NewDriverManager creates a new driver manager instance. It validates the driver and returns an error if the driver is invalid. The driver manager will listen to the specified events and reconcile the driver when those events are received.
func (c *Client) NewDriverManager(
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
		Client: c,
	}, nil
}

func (d *DriverManager) Run() error {
	// Setup NATS JetStream

	connectOptions := []nats.Option{
		nats.Name("Conveyor Driver Manager - " + d.Driver.Name),
		nats.MaxReconnects(-1), // infinite reconnects
	}

	if d.Client.Options.AuthEnabled {
		// the d.Client.Options.Cert and Key and RootCA are []byte of the files. we shall use nats.Secure()
		cert, err := tls.X509KeyPair(d.Client.Options.Cert, d.Client.Options.Key)
		if err != nil {
			color.Red("Error Occured while loading client certs: %v", err)
			return err
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(d.Client.Options.RootCA); !ok {
			color.Red("Error Occured while loading root CA cert: %v", err)
			return fmt.Errorf("failed to load root CA cert")
		}

		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
			InsecureSkipVerify: false,                          // Always verify the server's certificate
			ClientAuth:         tls.RequireAndVerifyClientCert, // Require client certs and verify them
		}

		connectOptions = append(connectOptions, nats.Secure(tlsConfig))

	}

	var natsConnectUrl string
	if d.Client.NatsURL != "" {
		natsConnectUrl = d.Client.NatsURL
	} else {
		natsConnectUrl = nats.DefaultURL
	}

	nc, err := nats.Connect(natsConnectUrl, connectOptions...)
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
		filterSubjects = append(filterSubjects, "drivers."+d.Driver.Name+".resources."+resource)
	}

	consumer, err := js.CreateOrUpdateConsumer(context.Background(), "messages", jetstream.ConsumerConfig{
		Name:           d.Driver.Name,
		FilterSubjects: filterSubjects,
		AckPolicy:      jetstream.AckExplicitPolicy,
		MaxAckPending:  1,
		// Deliver from last acknowledged message
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})

	if err != nil {
		color.Red("Error Occured while subscribing to NATS channel: %v", err)
		return err
	}

	// CONSUMER
	_, err = consumer.Consume(func(msg jetstream.Msg) {
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
