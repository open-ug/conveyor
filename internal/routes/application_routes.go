/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/handlers"
	streams "github.com/open-ug/conveyor/internal/streaming"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func ApplicationRoutes(app *fiber.App, db *clientv3.Client, natsCon *nats.Conn) {
	applicationPrefix := app.Group("/applications")
	// Create Application Handler
	applicationHandler := handlers.NewApplicationHandler(db, natsCon)
	streamHandler := streams.NewApplicationStreamer(natsCon, applicationHandler.ApplicationModel)
	execHandler, err := streams.NewContainerShellHandler()

	if err != nil {
		fmt.Println("Error occ")
	}

	applicationPrefix.Post("/", applicationHandler.CreateApplication)
	applicationPrefix.Get("/:name", applicationHandler.GetApplication)
	applicationPrefix.Get("/", applicationHandler.GetApplications)
	applicationPrefix.Put("/:name", applicationHandler.UpdateApplication)
	applicationPrefix.Delete("/:name", applicationHandler.DeleteApplication)
	applicationPrefix.Post("/:name/start", applicationHandler.StartApplication)
	applicationPrefix.Post("/:name/stop", applicationHandler.StopApplication)

	// Streams
	applicationPrefix.Get("/streams/logs/:name", websocket.New(streamHandler.StreamLogs))
	applicationPrefix.Get("/streams/exec/:name", websocket.New(execHandler.HandleWebSocket))

	//Metrics

	applicationPrefix.Post("/metrics/cpu/:name", applicationHandler.GetApplicationCPUUsage)
	applicationPrefix.Post("/metrics/memory/:name", applicationHandler.GetApplicationMemoryUsage)
}
