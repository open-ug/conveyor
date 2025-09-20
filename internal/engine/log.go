package engine

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/open-ug/conveyor/pkg/types"
)

func (ec *EngineContext) consumeLogEvents(msg jetstream.Msg) {
	// Logic to consume log-related events from NATS
	msg.Ack()

	data := msg.Data()

	var logEvent types.Log
	err := json.Unmarshal(data, &logEvent)
	if err != nil {
		log.Println("Error unmarshaling log event: ", err)
		return
	}

	// Process the log event
	err = ec.LogModel.Insert(logEvent)
	if err != nil {
		log.Println("Error inserting log event: ", err)
		return
	}
}
