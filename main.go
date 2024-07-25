package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create new sample GET routes
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("CRANE APPLICATION ORCHESTRATION RUNNING")
	})

	app.Post("/api/kube/metacontroller/apps/sync", func(c *fiber.Ctx) error {
		// print the request body
		fmt.Println(string(c.Body()))
		return c.SendString("CRANE APPLICATION ORCHESTRATION RUNNING")
	})

	app.Listen(":3000")
}
