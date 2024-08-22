/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package api

import (
	"os"

	routes "crane.cloud.cranom.tech/cmd/api/routes"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func StartServer(port string) {
	app := fiber.New(fiber.Config{
		AppName:     "Crane API Server",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("CRANE API SERVER contact info@cranom.tech for Documentation")
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	redisClient := NewRedisClient()

	uri := os.Getenv("CRANE_MONGO_URI")

	mongoClient := ConnectToMongoDB(uri)
	db := GetMongoDBDatabase(mongoClient, "crane")

	routes.ApplicationRoutes(app, db, redisClient)

	app.Listen(":" + port)
}
