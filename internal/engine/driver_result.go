package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
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
