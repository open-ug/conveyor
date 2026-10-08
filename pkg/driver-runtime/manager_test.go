package driverruntime_test

import (
	"context"
	"encoding/json"
	"errors"
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

func blockedManager(t *testing.T, existingEphemeral bool) (jetstream.JetStream, chan string, chan struct{}, <-chan error) {
	t.Helper()
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
	for _, config := range []jetstream.StreamConfig{
		{Name: "messages", Subjects: []string{"resources.>", "drivers.>"}},
		{Name: "pipeline-engine", Subjects: []string{"pipelines.>"}},
	} {
		if _, err := js.CreateStream(context.Background(), config); err != nil {
			t.Fatal(err)
		}
	}
	if existingEphemeral {
		_, err := js.CreateConsumer(context.Background(), "messages", jetstream.ConsumerConfig{Name: "test", FilterSubjects: []string{"resources.app", "drivers.test.resources.app"}, AckPolicy: jetstream.AckExplicitPolicy, MaxAckPending: 1})
		if err != nil {
			t.Fatal(err)
		}
	}
	client, err := runtime.NewClient("http://unused", s.ClientURL(), runtime.ConfigOptions{})
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan string, 2)
	release := make(chan struct{})
	manager, err := client.NewDriverManager(&runtime.Driver{Name: "test", Resources: []string{"app"}, Reconcile: func(_, _, runID string, _ *log.DriverLogger) types.DriverResult {
		started <- runID
		if runID == "first" {
			<-release
		}
		return types.DriverResult{Success: runID != "first"}
	}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- manager.Run() }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		select {
		case err := <-done:
			t.Fatalf("manager stopped: %v", err)
		default:
		}
		consumer, err := js.Consumer(context.Background(), "messages", "test")
		if err == nil {
			// Migration of an existing consumer also runs before the first delivery.
			info, err := consumer.Info(context.Background())
			if err == nil && (!existingEphemeral || info.Config.Durable != "") {
				break
			}
		}
		if time.Now().After(deadline) {
			t.Fatal("manager did not create its durable consumer")
		}
		time.Sleep(10 * time.Millisecond)
	}
	return js, started, release, done
}

func sendDriverMessage(t *testing.T, js jetstream.JetStream, runID string) {
	t.Helper()
	payload, _ := json.Marshal(types.Resource{Name: "app", Resource: "app"})
	message, _ := json.Marshal(types.DriverMessage{Payload: string(payload), Event: "update", RunID: runID})
	if _, err := js.Publish(context.Background(), "drivers.test.resources.app", message); err != nil {
		t.Fatal(err)
	}
}

func TestDriverManagerSurvivesLongFailedReconcile(t *testing.T) {
	for _, existing := range []bool{false, true} {
		name := "new consumer"
		if existing {
			name = "existing ephemeral consumer"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			js, started, release, done := blockedManager(t, existing)
			released := false
			defer func() {
				if !released {
					close(release)
				}
			}()
			sendDriverMessage(t, js, "first")
			select {
			case <-started:
			case <-time.After(2 * time.Second):
				t.Fatal("first reconcile did not start")
			}
			// No pull is pending while the callback runs. Ephemerals expire in 5–6s.
			time.Sleep(server.JsDeleteWaitTimeDefault + 2*time.Second)
			consumer, err := js.Consumer(context.Background(), "messages", "test")
			if err != nil {
				t.Fatalf("consumer disappeared during the build: %v", err)
			}
			info, err := consumer.Info(context.Background())
			if err != nil || info.Config.Durable != "test" || info.Config.InactiveThreshold != 0 {
				t.Fatalf("consumer is not durable: %+v err=%v", info, err)
			}
			close(release)
			released = true
			sendDriverMessage(t, js, "second")
			select {
			case runID := <-started:
				if runID != "second" {
					t.Fatalf("unexpected run %s", runID)
				}
			case err := <-done:
				t.Fatalf("manager stopped after failed reconcile: %v", err)
			case <-time.After(2 * time.Second):
				t.Fatal("second run never reached the driver")
			}
		})
	}
}

func TestDriverManagerReportsMissingConsumerAfterReconcile(t *testing.T) {
	js, started, release, done := blockedManager(t, false)
	sendDriverMessage(t, js, "first")
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		close(release)
		t.Fatal("reconcile did not start")
	}
	if err := js.DeleteConsumer(context.Background(), "messages", "test"); err != nil {
		close(release)
		t.Fatal(err)
	}
	close(release)
	select {
	case err := <-done:
		if !errors.Is(err, nats.ErrNoResponders) {
			t.Fatalf("missing consumer not reported: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("manager kept running without a consumer")
	}
}
