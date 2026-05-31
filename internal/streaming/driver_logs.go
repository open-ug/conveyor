package streaming

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/valyala/fasthttp"
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

// StreamDriverLogsByRunID streams logs for a specific driver and run ID using Server-Sent Events (SSE)
// @Summary Stream logs by driver name and run ID
// @Description Streams logs for a specific driver and run ID using Server-Sent Events (SSE)
// @Tags logs
// @Accept json
// @Produce json
// @Param drivername path string true "Driver Name"
// @Param runid path string true "Run ID"
// @Success 200 {string} string "Stream of log entries"
// @Failure 400 {object} map[string]string "Bad request - Invalid driver name or run ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs/streams/{drivername}/{runid} [get]
func (s *DriverLogsStreamer) StreamDriverLogsByRunID(c *fiber.Ctx) error {
	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	driverName := c.Params("drivername")
	runID := c.Params("runid")

	// Get the logs from the database
	logs, err := s.LogModel.Query("", driverName, runID)
	if err != nil {
		fmt.Println("Error getting logs:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting logs")
	}

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(
		fasthttp.StreamWriter(func(w *bufio.Writer) {
			// send dummy data
			w.Flush()

			for _, logEntry := range logs {
				jsonData, err := json.Marshal(logEntry)
				if err != nil {
					fmt.Println("Error: failed to marshal json")
					continue
				}

				fmt.Fprintf(w, "data: %s\n\n", jsonData)
				if err := w.Flush(); err != nil {
					return
				}
			}

			// Channel for incoming NATS messages
			msgCh := make(chan *nats.Msg, 256)

			// Async NATS subscription
			sub, err := s.NatsCon.Subscribe(
				"live.logs."+runID+"."+driverName,
				func(msg *nats.Msg) {
					select {
					case msgCh <- msg:
					default:
						// Drop message if client is slow (important!)
					}
				},
			)
			if err != nil {
				return
			}
			defer sub.Unsubscribe()

			// 4. Heartbeat ticker
			heartbeat := time.NewTicker(15 * time.Second)
			defer heartbeat.Stop()

			// 5. Streaming loop
			for {
				select {
				case <-heartbeat.C:
					// SSE heartbeat (comment line)
					fmt.Fprintf(w, ": heartbeat\n\n")
					if err := w.Flush(); err != nil {
						return
					}

				case msg := <-msgCh:
					fmt.Fprintf(w, "data: %s\n\n", msg.Data)
					if err := w.Flush(); err != nil {
						return
					}

				}
			}
		}),
	)

	return nil
}
