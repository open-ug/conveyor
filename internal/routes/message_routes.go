/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/open-ug/conveyor/internal/handlers"
	stream "github.com/open-ug/conveyor/internal/streaming"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func DriverRoutes(app *fiber.App, db *clientv3.Client, redisClient *redis.Client) {
	applicationPrefix := app.Group("/drivers")
	applicationHandler := handlers.NewMessageHandler(db, redisClient)

	applicationPrefix.Post("/broadcast-message", applicationHandler.PublishMessage)

	// Streams
	applicationPrefix.Get("/streams/logs/:name/:runid", websocket.New(stream.NewDriverLogsStreamer(redisClient).StreamLogs))
}
