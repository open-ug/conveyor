/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/open-ug/conveyor/cmd/api/handlers"
	streams "github.com/open-ug/conveyor/cmd/api/streaming"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApplicationRoutes(app *fiber.App, db *mongo.Database, redisClient *redis.Client) {
	applicationPrefix := app.Group("/applications")
	// Create Application Handler
	applicationHandler := handlers.NewApplicationHandler(db, redisClient)
	streamHandler := streams.NewApplicationStreamer(redisClient, applicationHandler.ApplicationModel)
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
