package dockerdriver

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gofiber/websocket/v2"
)

func StreamDockerLogs(containerID string, ws *websocket.Conn) error {
	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// Set options to stream logs
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	}

	// Open logs stream
	logStream, err := cli.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		return err
	}
	defer logStream.Close()

	// Stream logs to the WebSocket connection
	for {
		buf := make([]byte, 1024)
		n, err := logStream.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Log stream closed")
				break
			}
			return err
		}

		// Send log data to WebSocket client
		if n > 0 {
			if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return err
			}
		}
	}
	return nil
}
