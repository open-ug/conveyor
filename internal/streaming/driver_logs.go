package streaming

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
	logging "github.com/open-ug/conveyor/internal/logging"
)

type DriverLogsStreamer struct {
	NatsCon *nats.Conn
	Logger  *logging.LokiClient
}

func NewDriverLogsStreamer(NatsCon *nats.Conn) *DriverLogsStreamer {
	lokiClient := logging.New("http://localhost:3100")
	return &DriverLogsStreamer{
		NatsCon: NatsCon,
		Logger:  lokiClient,
	}
}

func (s *DriverLogsStreamer) StreamLogs(ws *websocket.Conn) {
	fmt.Println("Streaming driver logs")
	driverName := ws.Params("name")
	runID := ws.Params("runid")

	lokiQuery := map[string]string{
		"driver": driverName,
		"run_id": runID,
	}

	// First fetch previous logs in loki. time is zero
	logs, err := s.Logger.QueryLoki(lokiQuery, time.Time{}, time.Time{})
	if err != nil {
		fmt.Println("Error fetching logs from Loki:", err)
		ws.Close()
		return
	}
	for _, log := range logs {
		for _, line := range log.Values {
			// Send the log line to the WebSocket
			err = ws.WriteJSON(line)
			if err != nil {
				ws.Close()
				return
			}
		}
	}

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
		fmt.Println("Error subscribing to NATS channel:", err)
		ws.Close()
		return
	}

	defer sub.Unsubscribe()

	// 👇 Keep the connection open until the client closes it
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// Client disconnected
			break
		}
	}

	ws.Close()
}
