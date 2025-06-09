/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package api

import (
	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	routes "github.com/open-ug/conveyor/internal/routes"
	internals "github.com/open-ug/conveyor/internal/shared"
	utils "github.com/open-ug/conveyor/internal/utils"
	"github.com/spf13/viper"
)

func StartServer(port string) {

	app := fiber.New(fiber.Config{
		AppName:     "Conveyor API Server",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("CONVEYOR API SERVER. Visit https://conveyor.open.ug for Documentation")
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	natsContext := internals.NewNatsConn()
	natsContext.InitiateStreams()

	etcd, err := utils.NewEtcdClient(
		viper.GetString("etcd.host"),
	)

	if err != nil {
		color.Red("Error Occured while creating etcd client: %v", err)
		return
	}

	routes.ApplicationRoutes(app, etcd.Client, natsContext.NatsCon)
	routes.DriverRoutes(app, etcd.Client, natsContext.NatsCon)
	routes.ResourceRoutes(app, etcd.Client, natsContext)

	app.Listen(":" + port)

}
