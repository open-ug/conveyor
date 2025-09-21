package utils

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/spf13/viper"
)

type NatsContext struct {
	NatsCon    *nats.Conn
	JetStream  jetstream.JetStream
	natsServer *server.Server
}

func NewNatsConn() *NatsContext {

	conveyorDataDir := viper.GetString("api.data")
	dataDir := conveyorDataDir + "/nats"

	var port int
	if IsTestMode() {
		port = -1
	} else {
		port = 4222
	}
	opts := &server.Options{
		Port:      port,
		JetStream: true,
		StoreDir:  dataDir,
		NoLog:     true,
	}

	log.Println("Starting embedded NATS server with JetStream...")
	natsServer, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("failed to create NATS server: %v", err)
	}

	natsServer.ConfigureLogger()

	// Start the server in a separate goroutine so it doesn't block the main thread.
	go natsServer.Start()
	// Wait for the server to be ready.
	if !natsServer.ReadyForConnections(5 * time.Second) {
		log.Println("NATS server failed to start within timeout. retrying...")
		if !natsServer.ReadyForConnections(20 * time.Second) {
			log.Fatal("NATS server failed to start within timeout.")
		}
	}

	log.Printf("NATS server started on %s\n", natsServer.ClientURL())

	// --- Connect a NATS client to the embedded server ---
	nc, err := nats.Connect(natsServer.ClientURL())
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}
	return &NatsContext{
		NatsCon:    nc,
		JetStream:  js,
		natsServer: natsServer,
	}
}

func (n *NatsContext) InitiateStreams() error {

	// Create a stream for resource messages
	_, err := n.JetStream.CreateOrUpdateStream(context.Background(),
		jetstream.StreamConfig{
			Name:      "messages",
			Subjects:  []string{"resources.>", "events.>", "drivers.>"},
			Retention: jetstream.InterestPolicy,
		})
	if err != nil {
		return err
	}

	// Create a stream for pipeline events
	_, err = n.JetStream.CreateOrUpdateStream(context.Background(),
		jetstream.StreamConfig{
			Name:      "pipeline-engine",
			Subjects:  []string{"pipelines.>"},
			Retention: jetstream.WorkQueuePolicy,
		})
	if err != nil {
		return err
	}

	// Create a stream for logging
	_, err = n.JetStream.CreateOrUpdateStream(context.Background(),
		jetstream.StreamConfig{
			Name:      "logs-engine",
			Subjects:  []string{"logs.>"},
			Retention: jetstream.WorkQueuePolicy,
		})
	if err != nil {
		return err
	}

	return nil
}

// Shutdown gracefully stops the NATS client, the embedded NATS server,
// and removes the temporary JetStream data directory.
func (n *NatsContext) Shutdown() {
	log.Println("Initiating graceful shutdown of NATS resources...")

	// 1. Close the client connection first.
	if n.NatsCon != nil {
		n.NatsCon.Close()
	}

	// 2. Shutdown the embedded server.
	if n.natsServer != nil {
		n.natsServer.Shutdown()
	}

	log.Println("NATS shutdown complete.")
}
