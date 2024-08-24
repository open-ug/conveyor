/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package api

import (
	helpers "crane.cloud.cranom.tech/cmd/api/helpers"
	routes "crane.cloud.cranom.tech/cmd/api/routes"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
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

	privateKey, err := helpers.LoadPrivateKey()
	if err != nil {
		panic(err)
	}

	encryptedDbPass := viper.GetString("mongo.pass")
	decryptedDbPass, err := helpers.DecryptData([]byte(encryptedDbPass), privateKey)
	if err != nil {
		panic(err)
	}

	uri := "mongodb://" + viper.GetString("mongo.user") + ":" + string(decryptedDbPass) + "@" + viper.GetString("mongo.host") + ":" + viper.GetString("mongo.port")

	mongoClient := ConnectToMongoDB(uri)
	db := GetMongoDBDatabase(mongoClient, "crane")

	routes.ApplicationRoutes(app, db, redisClient)

	app.Listen(":" + port)
}
