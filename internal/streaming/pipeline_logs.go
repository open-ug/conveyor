package streaming

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/valyala/fasthttp"
)

type PipelineLogsStreamer struct {
	NatsContext *utils.NatsContext
	LogModal    *models.LogModel
}

func NewPipelineLogsStreamer(natsContext *utils.NatsContext, logModel *models.LogModel) *PipelineLogsStreamer {

	return &PipelineLogsStreamer{
		NatsContext: natsContext,
		LogModal:    logModel,
	}
}

// StreamLogsByRunID streams logs for a specific pipeline run ID using Server-Sent Events (SSE)
// @Summary Stream logs by pipeline run ID
// @Description Streams logs for a specific pipeline run ID using Server-Sent Events (SSE)
// @Tags logs
// @Accept json
// @Produce json
// @Param runid path string true "Run ID"
// @Success 200 {string} string "Stream of log entries"
// @Failure 400 {object} map[string]string "Bad request - Invalid run ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs/pipeline/{runid} [get]
func (s *PipelineLogsStreamer) StreamLogsByRunID(c *fiber.Ctx) error {
	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	runID := c.Params("runid")

	logs, err := s.LogModal.Query("", "", runID)
	if err != nil {
		fmt.Println("Error getting logs:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting logs")
	}

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(
		fasthttp.StreamWriter(func(w *bufio.Writer) {

			//Send existing logs first (replay)

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
			sub, err := s.NatsContext.NatsCon.Subscribe(
				"live.logs."+runID+".*",
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

			//Heartbeat ticker
			heartbeat := time.NewTicker(15 * time.Second)
			defer heartbeat.Stop()

			// Streaming loop
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
