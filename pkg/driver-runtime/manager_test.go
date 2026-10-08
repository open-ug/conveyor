package driverruntime_test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	"github.com/open-ug/conveyor/pkg/driver-runtime/log"
	"github.com/open-ug/conveyor/pkg/types"
)

func TestDriverManagerReportsConsumerExit(t *testing.T) {
	s, err := server.NewServer(&server.Options{Host: "127.0.0.1", Port: -1, JetStream: true, StoreDir: t.TempDir(), NoLog: true, NoSigs: true})
	if err != nil {
		t.Fatal(err)
	}
	s.Start()
	t.Cleanup(s.Shutdown)
	if !s.ReadyForConnections(5 * time.Second) {
		t.Fatal("NATS did not start")
	}
	nc, err := nats.Connect(s.ClientURL())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := js.CreateStream(context.Background(), jetstream.StreamConfig{Name: "messages", Subjects: []string{"resources.>", "drivers.>"}}); err != nil {
		t.Fatal(err)
	}
	client, err := runtime.NewClient("http://unused", s.ClientURL(), runtime.ConfigOptions{})
	if err != nil {
		t.Fatal(err)
	}
	manager, err := client.NewDriverManager(&runtime.Driver{Name: "test", Resources: []string{"app"}, Reconcile: func(_, _, _ string, _ *log.DriverLogger) types.DriverResult { return types.DriverResult{} }}, nil)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- manager.Run() }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := js.Consumer(context.Background(), "messages", "test"); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("driver consumer did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := js.DeleteConsumer(context.Background(), "messages", "test"); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("stopped driver returned no error")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("driver silently remained running after its consumer stopped")
	}
}
