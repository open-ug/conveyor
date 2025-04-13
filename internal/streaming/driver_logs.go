package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/websocket/v2"
	logging "github.com/open-ug/conveyor/internal/logging"
	"github.com/redis/go-redis/v9"
)

type DriverLogsStreamer struct {
	RedisClient *redis.Client
	Logger      *logging.LokiClient
}

func NewDriverLogsStreamer(redisClient *redis.Client) *DriverLogsStreamer {
	lokiClient := logging.New("http://localhost:3100")
	return &DriverLogsStreamer{
		RedisClient: redisClient,
		Logger:      lokiClient,
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

	// Subscribe to the Redis channel for driver logs
	pubsub := s.RedisClient.Subscribe(context.Background(), fmt.Sprintf("driver:%s:logs:%s", driverName, runID))
	defer pubsub.Close()

	// Wait for the subscription to be ready
	_, err = pubsub.Receive(context.Background())
	if err != nil {
		ws.Close()
		return
	}

	// Listen for messages on the Redis channel
	ch := pubsub.Channel()
	for msg := range ch {
		// Unmarshal the message
		var logMessage []string
		err := json.Unmarshal([]byte(msg.Payload), &logMessage)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		err = ws.WriteJSON(logMessage)
		if err != nil {
			ws.Close()
			return
		}
	}
	// Close the WebSocket connection when done
	ws.Close()

}
