/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package routes

import (
	"crane.cloud.cranom.tech/cmd/api/handlers"
	streams "crane.cloud.cranom.tech/cmd/api/streaming"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApplicationRoutes(app *fiber.App, db *mongo.Database, redisClient *redis.Client) {
	applicationPrefix := app.Group("/applications")
	applicationHandler := handlers.NewApplicationHandler(db, redisClient)
	streamHandler := streams.NewApplicationStreamer(redisClient, applicationHandler.ApplicationModel)

	applicationPrefix.Post("/", applicationHandler.CreateApplication)
	applicationPrefix.Get("/:name", applicationHandler.GetApplication)
	applicationPrefix.Get("/", applicationHandler.GetApplications)
	applicationPrefix.Put("/:name", applicationHandler.UpdateApplication)
	applicationPrefix.Delete("/:name", applicationHandler.DeleteApplication)
	applicationPrefix.Post("/:name/start", applicationHandler.StartApplication)
	applicationPrefix.Post("/:name/stop", applicationHandler.StopApplication)

	// Streams
	applicationPrefix.Get("/streams/logs/:name", websocket.New(streamHandler.StreamLogs))
}
