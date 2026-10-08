package utils

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

// CreateDurableConsumer keeps worker subscriptions alive while their callbacks
// run. A consumer Name alone is ephemeral and can expire during a long build.
func CreateDurableConsumer(ctx context.Context, js jetstream.JetStream, stream string, config jetstream.ConsumerConfig) (jetstream.Consumer, error) {
	config.Durable = config.Name
	config.InactiveThreshold = 0
	consumer, err := js.CreateOrUpdateConsumer(ctx, stream, config)
	if err != nil {
		return nil, err
	}
	if consumer.CachedInfo().Config.InactiveThreshold != 0 {
		// NATS 2.11 updates the expiry threshold using the previous durability
		// during an ephemeral-to-durable conversion. Apply the configuration
		// again once the consumer is durable to remove the inherited 5s expiry.
		consumer, err = js.CreateOrUpdateConsumer(ctx, stream, config)
		if err != nil {
			return nil, err
		}
		if consumer.CachedInfo().Config.InactiveThreshold != 0 {
			return nil, fmt.Errorf("consumer '%s' still has an inactivity timeout; check stream '%s' consumer limits", config.Name, stream)
		}
	}
	return consumer, nil
}
