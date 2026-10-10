package engine

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/pkg/types"
)

func deletionEngine(t *testing.T) (*EngineContext, types.Resource, *types.Pipeline) {
	t.Helper()
	db, err := badger.Open(badger.DefaultOptions(t.TempDir()).WithLogger(nil))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	model := models.NewResourceModel(db)
	resource := types.Resource{ID: "resource-id", Name: "app", Resource: "dockerdep", Pipeline: "deployment", Metadata: map[string]interface{}{"version": "1"}}
	data, _ := json.Marshal(resource)
	if err := model.Insert(resource.Name, resource.Resource, data); err != nil {
		t.Fatal(err)
	}
	resource, err = model.BeginDeletion(resource.Name, resource.Resource)
	if err != nil {
		t.Fatal(err)
	}
	if err := model.SetCurrentPipelineStep(resource.Name, resource.Resource, "cleanup"); err != nil {
		t.Fatal(err)
	}
	return &EngineContext{ResourceModel: model}, resource, &types.Pipeline{Steps: []types.Step{{ID: "cleanup", Driver: "status"}}}
}

func TestDeletionKeepsRecordAfterOperationalFailure(t *testing.T) {
	ec, resource, pipeline := deletionEngine(t)
	pipeline.Steps[0].Driver = "docker"
	ec.handleProcessDriverResult(PipelineEvent{Event: "delete", RunID: resource.Metadata["deletion_runid"].(string), Resource: resource, DriverResultEvent: DriverResultEvent{
		Driver: "docker", Success: true, Data: map[string]interface{}{"success": false},
	}}, pipeline)
	current, err := ec.ResourceModel.FindOne(resource.Name, resource.Resource)
	if err != nil {
		t.Fatalf("failed cleanup removed retry context: %v", err)
	}
	if current.Metadata["deletion_failed"] != true {
		t.Fatal("failure was not made visible")
	}
}

func TestDeletionRemovesRecordAndSecretsOnlyAfterFinalSuccess(t *testing.T) {
	ec, resource, pipeline := deletionEngine(t)
	ec.handleProcessDriverResult(PipelineEvent{Event: "delete", RunID: resource.Metadata["deletion_runid"].(string), Resource: resource, DriverResultEvent: DriverResultEvent{Driver: "status", Success: true}}, pipeline)
	_, err := ec.ResourceModel.FindOne(resource.Name, resource.Resource)
	if !errors.Is(err, badger.ErrKeyNotFound) {
		t.Fatalf("record retained after success: %v", err)
	}
	if _, err := ec.ResourceModel.FindByVersion(resource.Name, resource.Resource, "1"); err == nil {
		t.Fatal("secret-bearing snapshot survived deletion")
	}
}

func TestDeletionCannotBeReplacedWithUpdate(t *testing.T) {
	ec, resource, _ := deletionEngine(t)
	if _, err := ec.ResourceModel.Update(resource.Name, resource.Resource, resource); err == nil {
		t.Fatal("update allowed while deleting")
	}
}

func TestDeletionOperationIsPreservedAcrossStages(t *testing.T) {
	if nextOperation("delete") != "delete" {
		t.Fatal("cleanup became deployment processing")
	}
	if deletionStageSucceeded(DriverResultEvent{Success: false}) {
		t.Fatal("failed driver treated as successful")
	}
}

// Run the real engine through each cleanup stage, including the actual
// DriverResultEvent serialization, without Docker or live infrastructure.
func TestDeletionPipelineCompletesWithoutRecreatingResource(t *testing.T) {
	t.Run("explicit operation", func(t *testing.T) { runDeletionPipeline(t, true) })
	t.Run("legacy result event", func(t *testing.T) { runDeletionPipeline(t, false) })
}

func runDeletionPipeline(t *testing.T, includeOperation bool) {
	t.Helper()
	ec, resource, _ := deletionEngine(t)
	pipeline := &types.Pipeline{Name: resource.Pipeline, Resource: resource.Resource, Steps: []types.Step{
		{ID: "build", Driver: "buildpacks"}, {ID: "containers", Driver: "docker"},
		{ID: "routes", Driver: "ingress"}, {ID: "done", Driver: "status"},
	}}
	ec.PipelineModel = models.NewPipelineModel(ec.ResourceModel.DB)
	if err := ec.PipelineModel.CreatePipeline(pipeline); err != nil {
		t.Fatal(err)
	}
	publisher := &recordingPublisher{}
	ec.NatsContext.JetStream = publisher
	runID := resource.Metadata["deletion_runid"].(string)
	init, _ := json.Marshal(PipelineEvent{Event: "delete", RunID: runID, Resource: resource})
	ec.consumePipelineEvents(testMessage{subject: "pipelines.pipeline.init", data: init})
	for index, step := range pipeline.Steps {
		if publisher.subject != "drivers."+step.Driver+".resources."+resource.Resource {
			t.Fatalf("wrong stage %d: %s", index, publisher.subject)
		}
		var message types.DriverMessage
		if err := json.Unmarshal(publisher.data, &message); err != nil {
			t.Fatal(err)
		}
		if message.Event != "delete" {
			t.Fatalf("stage %d became %s", index, message.Event)
		}
		if err := json.Unmarshal([]byte(message.Payload), &resource); err != nil {
			t.Fatal(err)
		}
		result := DriverResultEvent{Driver: step.Driver, Success: true, Data: map[string]interface{}{"success": true}}
		var operations []string
		if includeOperation {
			operations = append(operations, message.Event)
		}
		if err := result.PublishEvent(runID, resource, publisher, operations...); err != nil {
			t.Fatal(err)
		}
		ec.consumePipelineEvents(testMessage{subject: "pipelines.driver.result", data: publisher.data})
		if index < len(pipeline.Steps)-1 {
			if _, err := ec.ResourceModel.FindOne(resource.Name, resource.Resource); err != nil {
				t.Fatal("resource vanished before cleanup finished")
			}
		}
	}
	if _, err := ec.ResourceModel.FindOne(resource.Name, resource.Resource); !errors.Is(err, badger.ErrKeyNotFound) {
		t.Fatalf("cleanup did not finish: %v", err)
	}
}

type recordingPublisher struct {
	jetstream.JetStream
	subject string
	data    []byte
}

func (p *recordingPublisher) Publish(_ context.Context, subject string, data []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	p.subject, p.data = subject, append([]byte{}, data...)
	return nil, nil
}

type testMessage struct {
	jetstream.Msg
	subject string
	data    []byte
}

func (m testMessage) Data() []byte    { return m.data }
func (m testMessage) Subject() string { return m.subject }
func (m testMessage) Ack() error      { return nil }

func TestOutdatedDeletionCannotFinishNewRun(t *testing.T) {
	ec, resource, _ := deletionEngine(t)
	oldRun := resource.Metadata["deletion_runid"].(string)
	if err := ec.ResourceModel.MarkDeletionFailed(resource.Name, resource.Resource); err != nil {
		t.Fatal(err)
	}
	retry, err := ec.ResourceModel.BeginDeletion(resource.Name, resource.Resource)
	if err != nil {
		t.Fatal(err)
	}
	if retry.Metadata["deletion_runid"] == oldRun {
		t.Fatal("retry reused failed operation")
	}
	if err := ec.ResourceModel.FinishDeletion(resource.Name, resource.Resource, oldRun); err == nil {
		t.Fatal("outdated deletion erased retry context")
	}
}
