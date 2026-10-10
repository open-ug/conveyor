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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/engine"
	"github.com/open-ug/conveyor/internal/utils"
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
	defer nc.Close()
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

	consumer, err := utils.CreateDurableConsumer(context.Background(), js, "messages", jetstream.ConsumerConfig{
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
	consumerErrors := make(chan error, 1)
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		data := msg.Data()
		var message types.DriverMessage
		err := json.Unmarshal([]byte(data), &message)
		if err != nil {
			color.Red("Error Occured while unmarshalling message: %v", err)
			return
		}

		logger := log.NewDriverLogger(d.Driver.Name, map[string]string{
			"event":  message.Event,
			"id":     message.ID,
			"run_id": message.RunID,
		}, nc)

		var pending types.Resource
		if json.Unmarshal([]byte(message.Payload), &pending) != nil {
			return
		}
		if pending.Pipeline != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			current, lookupErr := d.Client.GetResource(ctx, pending.Name, pending.Resource)
			cancel()
			if lookupErr != nil {
				var httpErr *HTTPStatusError
				if errors.As(lookupErr, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
					msg.Ack()
				} else {
					msg.NakWithDelay(time.Second)
				}
				return
			}
			// Queued stages must not recreate an application after deletion.
			if current.ID != pending.ID || (current.Metadata["deleting"] == true && (message.Event != "delete" || current.Metadata["deletion_runid"] != message.RunID)) {
				msg.Ack()
				return
			}
		}
		msg.Ack()
		result := d.Driver.Reconcile(message.Payload, message.Event, message.RunID, logger)

		driverevent := engine.DriverResultEvent{
			Success: result.Success,
			Message: result.Message,
			Driver:  d.Driver.Name,
			Data:    result.Data,
		}

		var resource types.Resource
		err = json.Unmarshal([]byte(message.Payload), &resource)
		if err != nil {
			color.Red("Error Occured while unmarshalling resource: %v", err)
			return
		}

		if err := driverevent.PublishEvent(message.RunID, resource, js, message.Event); err != nil {
			color.Red("Error publishing result for driver '%s', run '%s': %v", d.Driver.Name, message.RunID, err)
		}
	}, jetstream.PullMaxMessages(1), jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
		color.Red("Consumer error for driver '%s': %v", d.Driver.Name, err)
		if errors.Is(err, nats.ErrNoResponders) {
			select {
			case consumerErrors <- err:
			default:
			}
		}
	}))

	if err != nil {
		color.Red("Error Occured while consuming messages: %v", err)
		return err
	}
	defer cc.Stop()

	fmt.Println("Driver Manager is running for driver: ", d.Driver.Name)

	select {
	case <-cc.Closed():
		return fmt.Errorf("consumer stopped for driver '%s'", d.Driver.Name)
	case err := <-consumerErrors:
		return fmt.Errorf("consumer unavailable for driver '%s': %w", d.Driver.Name, err)
	}
}
