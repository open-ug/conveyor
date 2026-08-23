/*
Copyright © 2024 - Present Conveyor CI Contributors
*/
package routes

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/handlers"
)

func DriverRoutes(app *fiber.App, natsCon *nats.Conn, db *badger.DB) {
	applicationPrefix := app.Group("/drivers")
	applicationHandler := handlers.NewMessageHandler(natsCon, db)

	applicationPrefix.Post("/broadcast-message", applicationHandler.PublishMessage)

	// Streams

}
