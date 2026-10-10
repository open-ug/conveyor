package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
)

type DriverResultEvent struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Driver  string      `json:"driver"`
}

func (dre *DriverResultEvent) PublishEvent(
	run_id string,
	resource types.Resource,
	js jetstream.JetStream, operation ...string) error {
	pipelineEv := PipelineEvent{
		Event:             "driver.result",
		RunID:             run_id,
		Resource:          resource,
		DriverResultEvent: *dre,
	}

	if len(operation) > 0 {
		pipelineEv.Event = operation[0]
	}

	resultJson, err := json.Marshal(pipelineEv)
	if err != nil {
		return fmt.Errorf("marshal driver result event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = js.Publish(ctx, "pipelines.driver.result", resultJson)

	if err != nil {
		return fmt.Errorf("publish driver result event: %w", err)
	}
	return nil
}

func PublishResourceEvent(
	event string,
	resource types.Resource,
	js jetstream.JetStream) (string, error) {

	// Generate a unique run ID for this event
	run_id := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if event == "delete" {
		if id, ok := resource.Metadata["deletion_runid"].(string); ok {
			run_id = id
		}
	}

	// If the resource is not part of a pipeline, publish it directly to the resource subject
	if resource.Pipeline == "" {

		mID, err := utils.GenerateRandomID()
		if err != nil {
			return "", err
		}
		resourceData, err := json.Marshal(resource)
		if err != nil {
			return "", err
		}
		driverMsg := types.DriverMessage{
			ID:      mID,
			Payload: string(resourceData),
			Event:   event,
			RunID:   run_id,
		}

		jsonMsg, merr := json.Marshal(driverMsg)
		if merr != nil {
			return "", merr
		}

		// Publish message to jetstream
		subjectName := "resources." + resource.Resource
		_, err = js.Publish(ctx, subjectName, jsonMsg)

		if err != nil {
			return "", err
		}
		return run_id, nil
	} else {
		resourceEvent := PipelineEvent{
			Event:    event,
			RunID:    run_id,
			Resource: resource,
		}

		eventJson, err := json.Marshal(resourceEvent)
		if err != nil {
			return "", err
		}

		_, err = js.Publish(ctx, "pipelines.pipeline.init", eventJson)
		if err != nil {
			return "", err
		}
	}

	return run_id, nil
}
