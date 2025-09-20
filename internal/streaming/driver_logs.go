package streaming

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
)

type DriverLogsStreamer struct {
	NatsCon *nats.Conn
}

func NewDriverLogsStreamer(NatsCon *nats.Conn) *DriverLogsStreamer {

	return &DriverLogsStreamer{
		NatsCon: NatsCon,
	}
}

func (s *DriverLogsStreamer) StreamLogs(ws *websocket.Conn) {
	fmt.Println("Streaming driver logs")
	driverName := ws.Params("name")
	runID := ws.Params("runid")

	// Subscribe to the NATS channel for driver logs
	sub, errf := s.NatsCon.Subscribe(fmt.Sprintf("driver:%s:logs:%s", driverName, runID), func(msg *nats.Msg) {

		// Unmarshal the message
		var logMessage []string
		err := json.Unmarshal([]byte(msg.Data), &logMessage)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			ws.Close()
			return
		}

		err = ws.WriteJSON(logMessage)
		if err != nil {
			ws.Close()
			return
		}
	})
	if errf != nil {
		fmt.Println("Error subscribing to NATS channel:", errf)
		ws.Close()
		return
	}

	defer sub.Unsubscribe()

	// Keep the connection open until the client closes it
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// Client disconnected
			break
		}
	}

	ws.Close()
}
