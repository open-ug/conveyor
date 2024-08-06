package api

import (
	routes "crane.cloud.cranom.tech/cmd/api/routes"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func StartServer(port string) {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	redisClient := NewRedisClient()

	mongoClient := ConnectToMongoDB("mongodb+srv://jimjunior854:8bfKnA6cE2kq4kFW@cluster0.akews.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	db := GetMongoDBDatabase(mongoClient, "crane")

	routes.ApplicationRoutes(app, db, redisClient)

	app.Listen(":" + port)
}
