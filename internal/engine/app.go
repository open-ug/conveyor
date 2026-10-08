package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go"
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
	ctx           context.Context
	cancel        context.CancelFunc
}

type PipelineEvent struct {
	Event             string            `json:"event"` // e.g "create", "update", "delete"
	RunID             string            `json:"run_id"`
	Resource          types.Resource    `json:"resource"`
	DriverResultEvent DriverResultEvent `json:"driverresult"`
}

func NewEngineContext(logmodel *models.LogModel, natsContext utils.NatsContext, db *badger.DB) *EngineContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &EngineContext{
		NatsContext:   natsContext,
		PipelineModel: models.NewPipelineModel(db),
		ResourceModel: models.NewResourceModel(db),
		LogModel:      logmodel,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (ec *EngineContext) Start() error {
	log.Println("Starting the engine...")

	consumer, err := utils.CreateDurableConsumer(context.Background(), ec.NatsContext.JetStream, "pipeline-engine", jetstream.ConsumerConfig{
		Name:          "pipeline-engine",
		FilterSubject: "pipelines.>",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Println("Error creating consumer: ", err)
		return err
	}

	log.Println("Engine started and listening for pipeline events...")
	consumerErrors := make(chan error, 1)
	reportUnavailable := func(err error) {
		if errors.Is(err, nats.ErrNoResponders) {
			select {
			case consumerErrors <- err:
			default:
			}
		}
	}
	cc, err := consumer.Consume(ec.consumePipelineEvents, jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
		log.Printf("Pipeline consumer error: %v", err)
		reportUnavailable(err)
	}))
	if err != nil {
		log.Println("Error consuming pipeline events: ", err)
		return err
	}
	defer cc.Stop()

	// Create log consumer
	logconsumer, err := utils.CreateDurableConsumer(context.Background(), ec.NatsContext.JetStream, "logs-engine", jetstream.ConsumerConfig{
		Name:          "logs-engine",
		FilterSubject: "logs.>",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Println("Error creating consumer: ", err)
		return err
	}

	log.Println("Log consumer started...")
	lc, err := logconsumer.Consume(ec.consumeLogEvents, jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
		log.Printf("Log consumer error: %v", err)
		reportUnavailable(err)
	}))
	if err != nil {
		log.Println("Error consuming log events: ", err)
		return err
	}
	defer lc.Stop()

	select {
	case <-ec.ctx.Done():
		return nil
	case <-cc.Closed():
		return fmt.Errorf("pipeline consumer stopped unexpectedly")
	case <-lc.Closed():
		return fmt.Errorf("log consumer stopped unexpectedly")
	case err := <-consumerErrors:
		return fmt.Errorf("engine consumer unavailable: %w", err)
	}
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
		if len(pipeline.Steps) == 0 {
			log.Printf("Pipeline '%s' has no steps for run '%s'", pipeline.Name, event.RunID)
			return
		}
		firstStep := pipeline.Steps[0]
		event.Resource, err = ec.ResourceModel.StartPipelineRun(event.Resource.Name, event.Resource.Resource, event.RunID, event.Event, firstStep.ID)
		if err != nil {
			log.Println("Error starting pipeline run: ", err)
			return
		}

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
		driverSubject := "drivers." + firstStep.Driver + ".resources." + event.Resource.Resource
		if err := ec.publishEvent(driverSubject, driverMessage); err != nil {
			log.Printf("Failed to dispatch first step '%s' for run '%s': %v", firstStep.Name, event.RunID, err)
			return
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := ec.NatsContext.JetStream.Publish(ctx, subject, jsonMsg)

	if err != nil {
		log.Println("Error publishing event to driver: ", err)
		return err
	}
	return nil
}

func (ec *EngineContext) handleProcessDriverResult(event PipelineEvent, pipeline *types.Pipeline) {
	resource, err := ec.ResourceModel.FindOne(event.Resource.Name, event.Resource.Resource)
	if err != nil {
		log.Println("Error retrieving resource: ", err)
		return
	}
	if runID, ok := resource.Metadata["current_pipeline_run"].(string); ok && runID != event.RunID {
		log.Printf("Ignoring result from superseded run '%s' for resource '%s'", event.RunID, resource.Name)
		return
	}
	currentStepID, err := ec.ResourceModel.GetCurrentPipelineStep(resource.Name, resource.Resource)
	if err != nil {
		log.Println("Error getting current pipeline step: ", err)
		return
	}
	if stepID, ok := event.Resource.Metadata["current_pipeline_step"].(string); ok && stepID != currentStepID {
		log.Printf("Ignoring result from previous step '%s' for run '%s'", stepID, event.RunID)
		return
	}
	currentStepIndex := -1
	for i, step := range pipeline.Steps {
		if step.ID == currentStepID && step.Driver == event.DriverResultEvent.Driver {
			currentStepIndex = i
			break
		}
	}
	if currentStepIndex < 0 {
		log.Printf("Ignoring result from driver '%s' outside current step '%s' for run '%s'", event.DriverResultEvent.Driver, currentStepID, event.RunID)
		return
	}

	// save driver result to resource metadata
	err = ec.ResourceModel.SaveDriverResult(event.Resource.Name, event.Resource.Resource, event.DriverResultEvent.Driver, event.DriverResultEvent)
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

	if !event.DriverResultEvent.Success {
		log.Printf("Pipeline run '%s' stopped at step '%s' for resource '%s': %s", event.RunID, pipeline.Steps[currentStepIndex].Name, event.Resource.Name, event.DriverResultEvent.Message)
		return
	}

	if currentStepIndex+1 == len(pipeline.Steps) {
		log.Printf("Pipeline run '%s' completed for resource '%s'", event.RunID, event.Resource.Name)
		return
	}
	nextStep := pipeline.Steps[currentStepIndex+1]

	mID, err := utils.GenerateRandomID()
	if err != nil {
		log.Println("Error generating driver message ID: ", err)
		return
	}

	// Preserve the original operation for every driver in the run.
	operation, ok := event.Resource.Metadata["current_pipeline_event"].(string)
	if !ok || operation == "" {
		log.Printf("Missing pipeline operation for run '%s'; start a new run", event.RunID)
		return
	}
	event.Resource.Metadata["current_pipeline_step"] = nextStep.ID
	resourceJson, err := json.Marshal(event.Resource)
	if err != nil {
		log.Println("Error marshaling resource: ", err)
		return
	}

	driverMessage := types.DriverMessage{
		Event:   operation,
		RunID:   event.RunID,
		Payload: string(resourceJson),
		ID:      mID,
	}
	subject := "drivers." + nextStep.Driver + ".resources." + event.Resource.Resource
	if err := ec.publishEvent(subject, driverMessage); err != nil {
		log.Printf("Failed to dispatch step '%s' for run '%s': %v", nextStep.Name, event.RunID, err)
		return
	}

	if err := ec.ResourceModel.SetCurrentPipelineStep(event.Resource.Name, event.Resource.Resource, nextStep.ID); err != nil {
		log.Println("Error setting current pipeline step: ", err)
		return
	}
	log.Printf("Moving to next step '%s' for resource '%s' (run '%s')", nextStep.Name, event.Resource.Name, event.RunID)
}

func (ec *EngineContext) Stop() error {
	ec.cancel()
	return nil
}
