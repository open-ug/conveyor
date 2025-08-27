package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsContext struct {
	NatsCon    *nats.Conn
	JetStream  jetstream.JetStream
	natsServer *server.Server
}

func NewNatsConn() *NatsContext {

	// Create a temporary directory for JetStream's data.
	// This is critical, as JetStream needs a place to persist messages.
	dataDir, err := os.MkdirTemp("", "nats-jetstream-data-")
	if err != nil {
		log.Fatalf("failed to create temp dir for JetStream: %v", err)
	}

	opts := &server.Options{
		Port:      -1,
		JetStream: true,
		StoreDir:  dataDir,
	}

	log.Println("Starting embedded NATS server with JetStream...")
	natsServer, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("failed to create NATS server: %v", err)
	}

	// Start the server in a separate goroutine so it doesn't block the main thread.
	go natsServer.Start()
	// Wait for the server to be ready.
	if !natsServer.ReadyForConnections(5 * time.Second) {
		log.Fatal("NATS server failed to start within timeout.")
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
			Subjects:  []string{"resources.>", "events.>"},
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
