package streaming

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/models"
)

type DriverLogsStreamer struct {
	NatsCon  *nats.Conn
	LogModel *models.LogModel
}

func NewDriverLogsStreamer(NatsCon *nats.Conn, LogModel *models.LogModel) *DriverLogsStreamer {

	return &DriverLogsStreamer{
		NatsCon:  NatsCon,
		LogModel: LogModel,
	}
}

func (s *DriverLogsStreamer) StreamLogs(ws *websocket.Conn) {
	driverName := ws.Params("name")
	runID := ws.Params("runid")

	// Get the logs from the database
	logs, err := s.LogModel.Query("", driverName, runID)
	if err != nil {
		fmt.Println("Error getting logs:", err)
		ws.Close()
		return
	}
	// Send the existing logs to the client
	for _, logEntry := range logs {
		err = ws.WriteJSON([]string{logEntry.Timestamp, logEntry.Message})
		if err != nil {
			ws.Close()
			return
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
