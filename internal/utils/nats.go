package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/pkg/types"
)

type NatsContext struct {
	NatsCon    *nats.Conn
	JetStream  jetstream.JetStream
	natsServer *server.Server
}

func NewNatsConn(config *types.ServerConfig) *NatsContext {

	conveyorDataDir := config.API.Data
	dataDir := conveyorDataDir + "/nats"
	fmt.Println("Using NATS data directory:", dataDir)

	var port int
	if IsTestMode() {
		port = -1
		dataDir = "" // Use in-memory store for tests
	} else {
		port = config.NATS.Port
	}
	opts := &server.Options{
		Port:      port,
		JetStream: true,
		StoreDir:  dataDir,
		NoLog:     true,
		//NoSigs:    true,
	}

	authEnabled := config.API.AuthEnabled
	connectionOptions := []nats.Option{} // NATS CLIENT CONNNECTION OPTIONS
	if authEnabled {
		// enable TLS Authentication
		caFilePath := config.TLS.CA
		certFilePath := config.TLS.Cert
		keyFilePath := config.TLS.Key

		cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
		if err != nil {
			log.Fatalf("Error loading server certs: %v", err)
		}

		// 2. Load Root CA (For verifying clients)
		// We must manually build the CertPool since we want 'verify: true'
		caCert, err := os.ReadFile(caFilePath)
		if err != nil {
			log.Fatalf("Error loading CA cert: %v", err)
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)

		// 3. Create the TLS Config Object
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caPool,
			// This is the equivalent of 'opts.TLSVerify = true'
			ClientAuth: tls.RequireAndVerifyClientCert,
			MinVersion: tls.VersionTLS12,
		}

		opts.TLSConfig = tlsConfig

		// For the NATS client connection, we need to use the same certs to authenticate to the server

		connectionOptions = append(connectionOptions, nats.ClientCert(certFilePath, keyFilePath))
		connectionOptions = append(connectionOptions, nats.RootCAs(caFilePath))

		log.Println("NATS server TLS authentication enabled.")
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
		if !natsServer.ReadyForConnections(50 * time.Second) {
			log.Fatal("NATS server failed to start within timeout.")
		}
	}

	log.Printf("NATS server started on %s\n", natsServer.ClientURL())

	// NATS CLIENT CONNNECTION SETUP

	nc, err := nats.Connect(natsServer.ClientURL(), connectionOptions...)
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
	if n.NatsCon != nil && !n.NatsCon.IsClosed() {
		n.NatsCon.Close()
	}

	// 2. Shutdown the embedded server.
	if n.natsServer != nil && n.natsServer.Running() {
		n.natsServer.Shutdown()
	}

	log.Println("NATS shutdown complete.")
}
