/*
Copyright © 2024 Cranom Technologies Limited, Beingana Jim Junior and Contributors
*/
package api

import (
	"fmt"

	"context"
	"time"

	"github.com/docker/docker/client"
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

	fmt.Printf("LOKI API SERVER: %s\n", viper.GetString("loki.host"))
	fmt.Printf("ETCD API SERVER: %s\n", viper.GetString("etcd.host"))
	fmt.Printf("NATS API SERVER: %s\n", viper.GetString("nats.url"))

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("CONVEYOR API SERVER contact info@cranom.tech for Documentation")
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	natsCon := internals.NewNatsConn()

	etcd, err := utils.NewEtcdClient(
		viper.GetString("etcd.host"),
	)

	if err != nil {
		color.Red("Error Occured while creating etcd client: %v", err)
		return
	}

	routes.ApplicationRoutes(app, etcd.Client, natsCon)
	routes.DriverRoutes(app, etcd.Client, natsCon)

	app.Listen(":" + port)

}

func GetDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error creating docker client: %v", err)
	}
	return cli, nil
}

// a function that watches for docker container to start. it keeps checking the container status until it is running
func WatchContainerStart(containerID string) error {
	cli, err := GetDockerClient()
	if err != nil {
		return err
	}

	for {
		fmt.Println("Waiting for System Component to start")
		inspect, err := cli.ContainerInspect(context.Background(), containerID)
		if err != nil {
			color.Red("System Component failed to start")
			return fmt.Errorf("error inspecting container: %v", err)
		}

		if inspect.State.Status == "running" {
			//color.Green("System Component is running")
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
