/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package shared

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/spf13/viper"
)

type NatsContext struct {
	NatsCon   *nats.Conn
	JetStream jetstream.JetStream
}

func NewNatsConn() *NatsContext {
	nc, err := nats.Connect(viper.GetString("nats.url"))
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}
	return &NatsContext{
		NatsCon:   nc,
		JetStream: js,
	}
}

func (n *NatsContext) InitiateStreams() error {

	// Create a stream for resource messages
	_, err := n.JetStream.CreateOrUpdateStream(context.Background(),
		jetstream.StreamConfig{
			Name:      "messages",
			Subjects:  []string{"resources.>", "events.>"},
			Retention: jetstream.WorkQueuePolicy,
		})
	if err != nil {
		return err
	}

	return nil
}
