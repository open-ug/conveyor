package shared_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/shared"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNatsContext_Integration(t *testing.T) {
	// Ensure NATS URL is set (assume container is on localhost)
	viper.Set("nats.url", "nats://localhost:4222")
	if os.Getenv("CI") != "" {
		// If running in GitHub Actions, ensure service name works
		viper.Set("nats.url", "nats://nats:4222")
	}

	// 1. Connect to NATS
	nc := shared.NewNatsConn()
	assert.NotNil(t, nc.NatsCon)
	assert.NotNil(t, nc.JetStream)

	// 2. Initiate streams
	err := nc.InitiateStreams()
	assert.NoError(t, err, "Expected to create or update stream without error")

	// 3. Publish a test message to resources stream
	subject := "resources.test"
	msgData := []byte("hello-nats")

	_, err = nc.JetStream.Publish(context.Background(), subject, msgData)
	assert.NoError(t, err, "Expected to publish message without error")

	// 4. Subscribe and consume message
	consumer, err := nc.JetStream.CreateOrUpdateConsumer(context.Background(), "messages",
		sharedConsumerConfig("test-consumer", []string{subject}),
	)
	assert.NoError(t, err)

	// Use a channel to receive the message
	msgCh := make(chan string, 1)
	_, err = consumer.Consume(func(msg jetstream.Msg) {
		msgCh <- string(msg.Data())
		msg.Ack()
	})
	assert.NoError(t, err, "Expected to start consuming messages")

	// 5. Wait for the message or timeout
	select {
	case received := <-msgCh:
		assert.Equal(t, "hello-nats", received, "Expected to receive the same message")
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for NATS message")
	}

	// 6. Close connection
	nc.NatsCon.Close()
}

// helper to build consumer config
func sharedConsumerConfig(name string, subjects []string) jetstream.ConsumerConfig {
	return jetstream.ConsumerConfig{
		Name:           name,
		FilterSubjects: subjects,
		AckPolicy:      jetstream.AckExplicitPolicy,
		MaxAckPending:  1,
	}
}
