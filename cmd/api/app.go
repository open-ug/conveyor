/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package api

import (
	"fmt"

	"context"
	"time"

	helpers "crane.cloud.cranom.tech/cmd/api/helpers"
	routes "crane.cloud.cranom.tech/cmd/api/routes"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func StartServer(port string) {
	err := WatchContainerStart("crane-redis")
	if err != nil {
		panic(err)
	}

	err = WatchContainerStart("crane-mongo-db")
	if err != nil {
		panic(err)
	}
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

	privateKey, err := helpers.LoadPrivateKey()
	if err != nil {
		panic(err)
	}

	encryptedDbPass := viper.GetString("db.pass")
	//fmt.Println("Encrypted DB Pass: ", encryptedDbPass)
	decryptedDbPass, err := helpers.DecryptData(encryptedDbPass, privateKey)
	if err != nil {
		fmt.Println("Error decrypting DB Pass: ", err)
		panic(err)
	}

	uri := "mongodb://" + viper.GetString("db.user") + ":" + string(string(decryptedDbPass)) + "@" + viper.GetString("db.host") + ":" + viper.GetString("db.port")

	mongoClient := ConnectToMongoDB(uri)
	db := GetMongoDBDatabase(mongoClient, "crane")

	routes.ApplicationRoutes(app, db, redisClient)

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
