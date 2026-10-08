package engine

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
)

func testEngine(t *testing.T) (*EngineContext, jetstream.Consumer, map[string]jetstream.Consumer) {
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
	natsCtx := utils.NatsContext{NatsCon: nc, JetStream: js}
	if err := natsCtx.InitiateStreams(); err != nil {
		t.Fatal(err)
	}
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	ec := NewEngineContext(&models.LogModel{DB: db}, natsCtx, db)
	pipeline := &types.Pipeline{Name: "gitops", Resource: "gitop", Steps: []types.Step{
		{ID: "1", Name: "Build Application", Driver: "buildpacks"},
		{ID: "2", Name: "Update Deployment Status", Driver: "status"},
		{ID: "3", Name: "Deploy Application", Driver: "docker"},
		{ID: "4", Name: "Update Deployment Status", Driver: "status"},
	}}
	if err := ec.PipelineModel.CreatePipeline(pipeline); err != nil {
		t.Fatal(err)
	}
	consumer, err := js.CreateOrUpdateConsumer(context.Background(), "pipeline-engine", jetstream.ConsumerConfig{Name: "test-engine", AckPolicy: jetstream.AckExplicitPolicy})
	if err != nil {
		t.Fatal(err)
	}
	drivers := make(map[string]jetstream.Consumer)
	for _, name := range []string{"buildpacks", "status", "docker"} {
		consumer, err := js.CreateOrUpdateConsumer(context.Background(), "messages", jetstream.ConsumerConfig{Name: name, FilterSubject: "drivers." + name + ".resources.gitop", AckPolicy: jetstream.AckExplicitPolicy})
		if err != nil {
			t.Fatal(err)
		}
		drivers[name] = consumer
	}
	return ec, consumer, drivers
}

func processEvent(t *testing.T, ec *EngineContext, consumer jetstream.Consumer) {
	t.Helper()
	msg, err := consumer.Next(jetstream.FetchMaxWait(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	ec.consumePipelineEvents(msg)
}

func driverMessage(t *testing.T, consumer jetstream.Consumer, operation, runID, stepID string) (types.DriverMessage, types.Resource) {
	t.Helper()
	msg, err := consumer.Next(jetstream.FetchMaxWait(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if err := msg.Ack(); err != nil {
		t.Fatal(err)
	}
	var message types.DriverMessage
	if err := json.Unmarshal(msg.Data(), &message); err != nil {
		t.Fatal(err)
	}
	if message.Event != operation || message.RunID != runID {
		t.Fatalf("wrong operation or run: %+v", message)
	}
	var resource types.Resource
	if err := json.Unmarshal([]byte(message.Payload), &resource); err != nil {
		t.Fatal(err)
	}
	if resource.Metadata["current_pipeline_step"] != stepID {
		t.Fatalf("wrong step: %+v", resource.Metadata)
	}
	return message, resource
}

func TestPipelineFailureRetryAndNewResource(t *testing.T) {
	ec, consumer, drivers := testEngine(t)
	resource := types.Resource{Name: "roots-edu", Resource: "gitop", Pipeline: "gitops", Metadata: map[string]interface{}{"label": "keep"}}
	data, _ := json.Marshal(resource)
	if err := ec.ResourceModel.Insert(resource.Name, resource.Resource, data); err != nil {
		t.Fatal(err)
	}
	start := func(operation string, resource types.Resource) string {
		t.Helper()
		runID, err := PublishResourceEvent(operation, resource, ec.NatsContext.JetStream)
		if err != nil {
			t.Fatal(err)
		}
		processEvent(t, ec, consumer)
		return runID
	}
	result := func(driver string, success bool, runID string, resource types.Resource) {
		t.Helper()
		event := DriverResultEvent{Driver: driver, Success: success, Message: "test result"}
		if err := event.PublishEvent(runID, resource, ec.NatsContext.JetStream); err != nil {
			t.Fatal(err)
		}
		processEvent(t, ec, consumer)
	}
	firstRun := start("update", resource)
	_, build := driverMessage(t, drivers["buildpacks"], "update", firstRun, "1")
	result("buildpacks", true, firstRun, build)
	_, status := driverMessage(t, drivers["status"], "update", firstRun, "2")
	result("status", false, firstRun, status)
	step, err := ec.ResourceModel.GetCurrentPipelineStep(resource.Name, resource.Resource)
	if err != nil || step != "2" {
		t.Fatalf("failed run advanced: step=%s err=%v", step, err)
	}
	batch, err := drivers["docker"].FetchNoWait(1)
	if err != nil {
		t.Fatal(err)
	}
	for msg := range batch.Messages() {
		t.Fatalf("failed run dispatched Docker: %s", msg.Data())
	}
	if err := batch.Error(); err != nil {
		t.Fatal(err)
	}

	retryRun := start("update", resource)
	_, retryBuild := driverMessage(t, drivers["buildpacks"], "update", retryRun, "1")
	if results, ok := retryBuild.Metadata["driverresults"].(map[string]interface{}); !ok || len(results) != 0 {
		t.Fatalf("retry inherited previous results: %+v", retryBuild.Metadata)
	}
	if retryBuild.Metadata["label"] != "keep" {
		t.Fatal("resource metadata was lost")
	}
	// A late result from the failed run must not stop or advance the retry.
	result("status", false, firstRun, status)
	result("buildpacks", true, retryRun, retryBuild)
	_, retryStatus := driverMessage(t, drivers["status"], "update", retryRun, "2")
	result("status", true, retryRun, retryStatus)
	_, deployment := driverMessage(t, drivers["docker"], "update", retryRun, "3")
	// Repeated status drivers make a step check necessary as well as a driver check.
	result("docker", true, retryRun, deployment)
	_, finalStatus := driverMessage(t, drivers["status"], "update", retryRun, "4")
	result("status", true, retryRun, retryStatus)
	step, _ = ec.ResourceModel.GetCurrentPipelineStep(resource.Name, resource.Resource)
	if step != "4" {
		t.Fatalf("duplicate status result changed step to %s", step)
	}
	result("status", true, retryRun, finalStatus)

	resource.Name = "new-app"
	data, _ = json.Marshal(resource)
	if err := ec.ResourceModel.Insert(resource.Name, resource.Resource, data); err != nil {
		t.Fatal(err)
	}
	newRun := start("create", resource)
	_, newBuild := driverMessage(t, drivers["buildpacks"], "create", newRun, "1")
	result("buildpacks", false, newRun, newBuild)
	// A first-step failure also allows a fresh run.
	newRun = start("create", resource)
	driverMessage(t, drivers["buildpacks"], "create", newRun, "1")
}

func TestPublishResourceEventReportsServerRejection(t *testing.T) {
	ec, _, _ := testEngine(t)
	if err := ec.NatsContext.JetStream.DeleteStream(context.Background(), "pipeline-engine"); err != nil {
		t.Fatal(err)
	}
	runID, err := PublishResourceEvent("create", types.Resource{Name: "app", Resource: "gitop", Pipeline: "gitops"}, ec.NatsContext.JetStream)
	if err == nil || runID != "" {
		t.Fatalf("returned success for an event NATS rejected: run=%s err=%v", runID, err)
	}
}

func TestEngineReportsConsumerExit(t *testing.T) {
	ec, _, _ := testEngine(t)
	if err := ec.NatsContext.JetStream.DeleteConsumer(context.Background(), "pipeline-engine", "test-engine"); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- ec.Start() }()
	t.Cleanup(func() { ec.Stop() })
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := ec.NatsContext.JetStream.Consumer(context.Background(), "pipeline-engine", "pipeline-engine"); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("engine consumer did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := ec.NatsContext.JetStream.DeleteConsumer(context.Background(), "pipeline-engine", "pipeline-engine"); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("engine did not report its stopped consumer")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("engine remained running with a stopped consumer")
	}
}

func TestFailedDispatchDoesNotAdvanceStep(t *testing.T) {
	ec, consumer, drivers := testEngine(t)
	resource := types.Resource{Name: "app", Resource: "gitop", Pipeline: "gitops"}
	data, _ := json.Marshal(resource)
	if err := ec.ResourceModel.Insert(resource.Name, resource.Resource, data); err != nil {
		t.Fatal(err)
	}
	runID, err := PublishResourceEvent("create", resource, ec.NatsContext.JetStream)
	if err != nil {
		t.Fatal(err)
	}
	processEvent(t, ec, consumer)
	_, build := driverMessage(t, drivers["buildpacks"], "create", runID, "1")
	if err := ec.NatsContext.JetStream.DeleteStream(context.Background(), "messages"); err != nil {
		t.Fatal(err)
	}
	result := DriverResultEvent{Driver: "buildpacks", Success: true}
	if err := result.PublishEvent(runID, build, ec.NatsContext.JetStream); err != nil {
		t.Fatal(err)
	}
	processEvent(t, ec, consumer)
	step, err := ec.ResourceModel.GetCurrentPipelineStep(resource.Name, resource.Resource)
	if err != nil || step != "1" {
		t.Fatalf("advanced after dispatch failed: step=%s err=%v", step, err)
	}
}
