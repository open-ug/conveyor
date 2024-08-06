package routes

import (
	"crane.cloud.cranom.tech/cmd/api/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApplicationRoutes(app *fiber.App, db *mongo.Database) {
	applicationPrefix := app.Group("/applications")
	applicationHandler := handlers.NewApplicationHandler(db)

	applicationPrefix.Post("/", applicationHandler.CreateApplication)
	applicationPrefix.Get("/:name", applicationHandler.GetApplication)
	applicationPrefix.Get("/", applicationHandler.GetApplications)
	applicationPrefix.Put("/:name", applicationHandler.UpdateApplication)
}
