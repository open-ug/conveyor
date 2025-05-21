/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/handlers"
	stream "github.com/open-ug/conveyor/internal/streaming"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func DriverRoutes(app *fiber.App, db *clientv3.Client, natsCon *nats.Conn) {
	applicationPrefix := app.Group("/drivers")
	applicationHandler := handlers.NewMessageHandler(db, natsCon)

	applicationPrefix.Post("/broadcast-message", applicationHandler.PublishMessage)

	// Streams
	applicationPrefix.Get("/streams/logs/:name/:runid", websocket.New(stream.NewDriverLogsStreamer(natsCon).StreamLogs))
}
