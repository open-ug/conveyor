/*
Copyright © 2024 - Present Conveyor CI Contributors
*/
package routes

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/handlers"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func DriverRoutes(app *fiber.App, cli *clientv3.Client, natsCon *nats.Conn, db *badger.DB) {
	applicationPrefix := app.Group("/drivers")
	applicationHandler := handlers.NewMessageHandler(cli, natsCon, db)

	applicationPrefix.Post("/broadcast-message", applicationHandler.PublishMessage)

	// Streams

}
