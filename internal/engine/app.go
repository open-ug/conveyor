package engine

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EngineContext struct {
	NatsContext   utils.NatsContext
	PipelineModel *models.PipelineModel
	ResourceModel *models.ResourceModel
}

type PipelineEvent struct {
	Event    string `json:"event"` // e.g "create", "update", "delete"
	RunID    string `json:"run_id"`
	Resource string `json:"resource"`
}

func NewEngineContext(db *clientv3.Client, natsContext utils.NatsContext) *EngineContext {

	return &EngineContext{
		NatsContext:   natsContext,
		PipelineModel: models.NewPipelineModel(db),
		ResourceModel: models.NewResourceModel(db),
	}
}

func (ec *EngineContext) Start() error {
	// Start the engine as a goroutine
	go ec.startEngine()
	return nil
}

func (ec *EngineContext) startEngine() error {
	consumer, err := ec.NatsContext.JetStream.CreateOrUpdateConsumer(context.Background(), "pipeline-engine", jetstream.ConsumerConfig{
		Name:          "pipeline-engine",
		FilterSubject: "pipelines.*",
		AckPolicy:     jetstream.AckAllPolicy,
	})
	if err != nil {
		return err
	}

	consumer.Consume(ec.consumePipelineEvents)
	select {}
}

func (ec *EngineContext) consumePipelineEvents(msg jetstream.Msg) {
	// Logic to consume pipeline-related events from NATS
	msg.Ack()

	data := msg.Data()
	var event PipelineEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		// Handle error
		return
	}

	var resource types.Resource
	err = json.Unmarshal([]byte(event.Resource), &resource)
	if err != nil {
		// Handle error
		return
	}

	// Get pipeline details
	pipeline, err := ec.PipelineModel.GetPipeline(resource.Pipeline)
	if err != nil {
		// Handle error
		return
	}

	mID, _ := utils.GenerateRandomID()

	// Driver message
	driverMessage := types.DriverMessage{
		Event:   event.Event,
		RunID:   event.RunID,
		Payload: event.Resource,
		ID:      mID,
	}

}

// publishEvent publishes an event to driver
func (ec *EngineContext) publishEvent(subject string, message types.DriverMessage) error {

	jsonMsg, merr := json.Marshal(message)
	if merr != nil {
		return merr
	}

	_, err := ec.NatsContext.JetStream.PublishAsync(subject, jsonMsg)

	if err != nil {
		return err
	}
	return nil
}
