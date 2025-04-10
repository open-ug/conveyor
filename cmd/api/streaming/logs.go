package streaming

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	models "conveyor.cloud.cranom.tech/cmd/api/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
)

type ApplicationStreamer struct {
	RedisClient      *redis.Client
	ApplicationModel *models.ApplicationModel
}

func NewApplicationStreamer(redisClient *redis.Client, appModel *models.ApplicationModel) *ApplicationStreamer {
	return &ApplicationStreamer{
		RedisClient:      redisClient,
		ApplicationModel: appModel,
	}
}

func (s *ApplicationStreamer) StreamLogs(ws *websocket.Conn) {
	fmt.Println("Streaming logs")
	appName := ws.Params("name")

	filter := map[string]interface{}{
		"name": appName,
	}
	app, err := s.ApplicationModel.FindOne(filter)
	if err != nil {
		fmt.Println(err)
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	fmt.Println("Fetched model")

	err = StreamDockerLogs(app.Name, ws)
	if err != nil {
		fmt.Println(err)
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}
}

func StreamDockerLogs(containerID string, ws *websocket.Conn) error {
	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to get docker client")
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

	buf := make([]byte, 1024)

	// Stream logs to the WebSocket connection
	for {
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
			logOutput := buf[:n]

			encodedData := base64.StdEncoding.EncodeToString(logOutput)

			if err := ws.WriteMessage(websocket.TextMessage, []byte(encodedData)); err != nil {
				return err
			}
		}
	}
	return nil
}
