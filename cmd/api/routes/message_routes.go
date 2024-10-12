/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package routes

import (
	"crane.cloud.cranom.tech/cmd/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func DriverRoutes(app *fiber.App, db *mongo.Database, redisClient *redis.Client) {
	applicationPrefix := app.Group("/drivers")
	applicationHandler := handlers.NewMessageHandler(db, redisClient)

	applicationPrefix.Post("/broadcast-message", applicationHandler.PublishMessage)
}
