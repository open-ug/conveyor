package streaming

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gofiber/websocket/v2"
)

type ContainerShellHandler struct {
	dockerClient *client.Client
}

func NewContainerShellHandler() (*ContainerShellHandler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &ContainerShellHandler{dockerClient: cli}, nil
}

func (h *ContainerShellHandler) HandleWebSocket(c *websocket.Conn) {
	appName := c.Params("name")

	// Create an exec instance for the shell
	execID, err := h.dockerClient.ContainerExecCreate(context.Background(), appName, container.ExecOptions{
		Cmd:          []string{"/bin/sh"},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		log.Printf("Error creating exec instance: %s", err)
		return
	}

	// Attach to the exec instance
	resp, err := h.dockerClient.ContainerExecAttach(context.Background(), execID.ID, container.ExecStartOptions{
		Tty: true,
	})
	if err != nil {
		log.Printf("Error attaching to exec instance: %s", err)
		return
	}
	defer resp.Close()

	// Bidirectional communication
	errChan := make(chan error, 2)

	// Read from WebSocket and send to container
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			if _, err := resp.Conn.Write(msg); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Read from container and send to WebSocket
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := resp.Reader.Read(buffer)
			if err != nil {
				errChan <- err
				return
			}
			if err := c.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Wait for an error or connection closure
	<-errChan
}
