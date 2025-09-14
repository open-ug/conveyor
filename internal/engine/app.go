package engine

import (
	"context"
	"encoding/json"
	"log"

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
	Event             string            `json:"event"` // e.g "create", "update", "delete"
	RunID             string            `json:"run_id"`
	Resource          types.Resource    `json:"resource"`
	DriverResultEvent DriverResultEvent `json:"driverresult"`
}

func NewEngineContext(db *clientv3.Client, natsContext utils.NatsContext) *EngineContext {

	return &EngineContext{
		NatsContext:   natsContext,
		PipelineModel: models.NewPipelineModel(db),
		ResourceModel: models.NewResourceModel(db),
	}
}

func (ec *EngineContext) Start() error {
	log.Println("Starting the engine...")
	consumer, err := ec.NatsContext.JetStream.CreateOrUpdateConsumer(context.Background(), "pipeline-engine", jetstream.ConsumerConfig{
		Name:          "pipeline-engine",
		FilterSubject: "pipelines.>",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Println("Error creating consumer: ", err)
		return err
	}

	log.Println("Engine started and listening for pipeline events...")
	cc, err := consumer.Consume(ec.consumePipelineEvents)
	if err != nil {
		log.Println("Error consuming pipeline events: ", err)
		return err
	}
	defer cc.Stop()

	select {}
}

func (ec *EngineContext) consumePipelineEvents(msg jetstream.Msg) {
	// Logic to consume pipeline-related events from NATS
	msg.Ack()

	data := msg.Data()
	subject := msg.Subject()
	var event PipelineEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		// Handle error
		return
	}

	// Get pipeline details
	pipeline, err := ec.PipelineModel.GetPipeline(event.Resource.Pipeline)
	if err != nil {
		// Handle error
		return
	}

	mID, _ := utils.GenerateRandomID()

	if subject == "pipelines.driver.result" {
		if event.Resource.Pipeline == "" {
			// No pipeline associated, ignore
			return
		}
		// Process driver result and move to next step
		ec.handleProcessDriverResult(event, pipeline)

	} else if subject == "pipelines.pipeline.init" {

		resourceJson, err := json.Marshal(event.Resource)
		if err != nil {
			// Handle error
			return
		}
		// Driver message
		driverMessage := types.DriverMessage{
			Event:   event.Event,
			RunID:   event.RunID,
			Payload: string(resourceJson),
			ID:      mID,
		}

		// Publish to the first step's driver
		if len(pipeline.Steps) > 0 {
			firstStep := pipeline.Steps[0]
			subject := firstStep.Driver + ".resources." + event.Resource.Resource
			ec.publishEvent(subject, driverMessage)
		}
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

func (ec *EngineContext) handleProcessDriverResult(event PipelineEvent, pipeline *types.Pipeline) {

	// Find the current step based on the driver name
	var currentStepIndex int = -1
	for i, step := range pipeline.Steps {
		if step.Driver == event.DriverResultEvent.Driver {
			currentStepIndex = i
			break
		}
	}

	if currentStepIndex == -1 {
		// Current step not found, handle error
		return
	}

	// If the driver result indicates success, move to the next step
	if event.DriverResultEvent.Success {
		nextStepIndex := currentStepIndex + 1
		if nextStepIndex < len(pipeline.Steps) {
			nextStep := pipeline.Steps[nextStepIndex]

			mID, _ := utils.GenerateRandomID()

			resourceJson, err := json.Marshal(event.Resource)
			if err != nil {
				// Handle error
				return
			}

			// Driver message for the next step
			driverMessage := types.DriverMessage{
				Event:   "process",
				RunID:   event.RunID,
				Payload: string(resourceJson),
				ID:      mID,
			}

			subject := nextStep.Driver + ".resources." + event.Resource.Resource
			ec.publishEvent(subject, driverMessage)
		} else {
			// Pipeline completed successfully
			// You can add logic here to mark the pipeline as completed in your system
		}
	} else {
		// Handle failure case, e.g., log the error or notify someone
	}
}

func (ec *EngineContext) Stop() error {
	// Logic to stop the engine
	return nil
}
