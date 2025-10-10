package utils_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/config"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestNatsContext_Integration(t *testing.T) {
	// Set test environment to enable random port allocation
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "test")
	defer os.Setenv("APP_ENV", originalEnv) // Restore original value after test

	config.InitConfig()

	// 1. Connect to NATS
	nc := utils.NewNatsConn()
	assert.NotNil(t, nc.NatsCon)
	assert.NotNil(t, nc.JetStream)

	// Ensure proper cleanup of NATS resources
	defer nc.Shutdown()

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
}

// helper to build consumer config
func sharedConsumerConfig(name string, subjects []string) jetstream.ConsumerConfig {
	return jetstream.ConsumerConfig{
		Name:           name,
		FilterSubjects: subjects,
	}
}
