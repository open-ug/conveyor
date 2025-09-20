package engine

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
)

type DriverResultEvent struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Driver  string `json:"driver"`
}

func (dre *DriverResultEvent) PublishEvent(
	run_id string,
	resource types.Resource,
	js jetstream.JetStream) {
	pipelineEv := PipelineEvent{
		Event:             "driver.result",
		RunID:             run_id,
		Resource:          resource,
		DriverResultEvent: *dre,
	}

	resultJson, err := json.Marshal(pipelineEv)
	if err != nil {
		fmt.Printf("Error marshalling driver result event: %v", err)
		return
	}

	_, err = js.PublishAsync("pipelines.driver.result", resultJson)

	if err != nil {
		fmt.Printf("Error publishing driver result event: %v", err)
	}

}

func PublishResourceEvent(
	event string,
	resource types.Resource,
	js jetstream.JetStream) (string, error) {

	run_id := uuid.New().String()
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
			Event:   "create",
			RunID:   run_id,
		}

		jsonMsg, merr := json.Marshal(driverMsg)
		if merr != nil {
			return "", merr
		}

		// Publish message to jetstream
		subjectName := "resources." + resource.Resource
		_, err = js.PublishAsync(subjectName, jsonMsg)

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

		_, err = js.PublishAsync("pipelines.pipeline.init", eventJson)
		if err != nil {
			return "", err
		}
	}

	return run_id, nil
}
