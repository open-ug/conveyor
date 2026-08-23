package engine

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
)

type EngineContext struct {
	NatsContext   utils.NatsContext
	PipelineModel *models.PipelineModel
	ResourceModel *models.ResourceModel
	LogModel      *models.LogModel
}

type PipelineEvent struct {
	Event             string            `json:"event"` // e.g "create", "update", "delete"
	RunID             string            `json:"run_id"`
	Resource          types.Resource    `json:"resource"`
	DriverResultEvent DriverResultEvent `json:"driverresult"`
}

func NewEngineContext(logmodel *models.LogModel, natsContext utils.NatsContext, db *badger.DB) *EngineContext {

	return &EngineContext{
		NatsContext:   natsContext,
		PipelineModel: models.NewPipelineModel(db),
		ResourceModel: models.NewResourceModel(db),
		LogModel:      logmodel,
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

	// Create log consumer
	logconsumer, err := ec.NatsContext.JetStream.CreateOrUpdateConsumer(context.Background(), "logs-engine", jetstream.ConsumerConfig{
		Name:          "logs-engine",
		FilterSubject: "logs.>",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Println("Error creating consumer: ", err)
		return err
	}

	log.Println("Log consumer started...")
	lc, err := logconsumer.Consume(ec.consumeLogEvents)
	if err != nil {
		log.Println("Error consuming log events: ", err)
		return err
	}
	defer lc.Stop()

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
		log.Println("Error unmarshaling pipeline event: ", err)
		return
	}

	if event.Resource.Pipeline == "" {
		// No pipeline associated, ignore
		return
	}

	// Get pipeline details
	pipeline, err := ec.PipelineModel.GetPipeline(event.Resource.Pipeline)
	if err != nil {
		log.Println("Error getting pipeline details: ", err)
		return
	}

	mID, _ := utils.GenerateRandomID()

	switch subject {
	case "pipelines.driver.result":

		// Process driver result and move to next step
		ec.handleProcessDriverResult(event, pipeline)

	case "pipelines.pipeline.init":

		resourceJson, err := json.Marshal(event.Resource)
		if err != nil {
			log.Println("Error marshaling resource: ", err)
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
			// set current step in resource metadata
			err = ec.ResourceModel.SetCurrentPipelineStep(event.Resource.Name, event.Resource.Resource, firstStep.ID)
			if err != nil {
				log.Println("Error setting current pipeline step: ", err)
				return
			}
			subject := "drivers." + firstStep.Driver + ".resources." + event.Resource.Resource
			err = ec.publishEvent(subject, driverMessage)
			if err != nil {
				log.Println("Error publishing event to driver: ", err)
				return
			}
		}
	default:
		log.Println("Unhandled pipeline event subject: ", subject)
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
		log.Println("Error publishing event to driver: ", err)
		return err
	}
	return nil
}

func (ec *EngineContext) handleProcessDriverResult(event PipelineEvent, pipeline *types.Pipeline) {

	// save driver result to resource metadata
	err := ec.ResourceModel.SaveDriverResult(event.Resource.Name, event.Resource.Resource, event.DriverResultEvent.Driver, event.DriverResultEvent)
	if err != nil {
		log.Println("Error saving driver result: ", err)
		return
	}

	// update resource to include latest driver result
	updatedResource, err := ec.ResourceModel.FindOne(event.Resource.Name, event.Resource.Resource)
	if err != nil {
		log.Println("Error retrieving updated resource: ", err)
		return
	}
	event.Resource = updatedResource

	currentStepID, err := ec.ResourceModel.GetCurrentPipelineStep(event.Resource.Name, event.Resource.Resource)
	if err != nil {
		log.Println("Error getting current pipeline step: ", err)
		return
	}

	var nextStep *types.Step

	for i, step := range pipeline.Steps {
		if step.ID == currentStepID {
			if i+1 < len(pipeline.Steps) {
				nextStep = &pipeline.Steps[i+1]
			}
			break
		}
	}

	if nextStep == nil {
		log.Println("Pipeline completed or no next step found for resource: ", event.Resource.Name)
		return
	}

	log.Printf("Moving to next step '%s' for resource '%s'\n", nextStep.Name, event.Resource.Name)

	// If the driver result indicates success, move to the next step
	if event.DriverResultEvent.Success {

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

		subject := "drivers." + nextStep.Driver + ".resources." + event.Resource.Resource
		ec.publishEvent(subject, driverMessage)

		// Update the current step in resource metadata
		err = ec.ResourceModel.SetCurrentPipelineStep(event.Resource.Name, event.Resource.Resource, nextStep.ID)
		if err != nil {
			log.Println("Error setting current pipeline step: ", err)
			return
		}

	} else {
		// TODO: Handle failure case
	}
}

func (ec *EngineContext) Stop() error {
	// Logic to stop the engine
	return nil
}
