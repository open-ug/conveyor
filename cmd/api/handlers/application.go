package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	models "crane.cloud.cranom.tech/cmd/api/models"
	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationHandler struct {
	RedisClient      *redis.Client
	ApplicationModel *models.ApplicationModel
}

func NewApplicationHandler(db *mongo.Database, redisClient *redis.Client) *ApplicationHandler {
	return &ApplicationHandler{
		RedisClient: redisClient,
		ApplicationModel: &models.ApplicationModel{
			Collection: db.Collection("applications"),
		},
	}
}

func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	fmt.Println("Creating App")
	var app craneTypes.Application
	if err := c.BodyParser(&app); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}
	insertResult, err := h.ApplicationModel.Insert(app)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("Saved to database")
	appMsg := craneTypes.ApplicationMsg{
		Action:  "create",
		ID:      insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Payload: app,
	}

	jsonMsg, merr := json.Marshal(appMsg)
	if merr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}
	// Publish to redis channel for driver to work on it
	errf := h.RedisClient.Publish(context.Background(), "application",
		jsonMsg,
	).Err()
	fmt.Println("Sent Redis Pub")

	if errf != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":   insertResult.InsertedID,
		"name": app.Name,
	})
}

func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	filter := map[string]interface{}{
		"name": c.Params("name"),
	}
	app, err := h.ApplicationModel.FindOne(filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println(app)

	return c.Status(fiber.StatusOK).JSON(app)
}

func (h *ApplicationHandler) GetApplications(c *fiber.Ctx) error {
	apps, err := h.ApplicationModel.Find(nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(apps)
}

func (h *ApplicationHandler) UpdateApplication(c *fiber.Ctx) error {
	filter := map[string]interface{}{
		"name": c.Params("name"),
	}
	var appl craneTypes.Application

	if err := c.BodyParser(&appl); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}

	update := map[string]interface{}{
		"$set": appl,
	}
	updateResult, err := h.ApplicationModel.UpdateOne(filter, update)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Publish to redis channel for driver to work on it
	var app craneTypes.Application
	newapp := h.ApplicationModel.Collection.FindOne(context.Background(), filter)

	err = newapp.Decode(&app)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fmt.Println("Saved to database")

	rowapp, err := newapp.Raw()

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	appId := rowapp.Lookup("_id").ObjectID().Hex()

	appMsg := craneTypes.ApplicationMsg{
		Action:  "update",
		ID:      appId,
		Payload: app,
	}

	jsonMsg, merr := json.Marshal(appMsg)
	if merr != nil {
		fmt.Println(merr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	fmt.Println("Sent Redis Pub")
	errf := h.RedisClient.Publish(context.Background(), "application", jsonMsg).Err()

	if errf != nil {
		fmt.Println(errf)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updateResult)
}

func (h *ApplicationHandler) DeleteApplication(c *fiber.Ctx) error {
	filter := map[string]interface{}{
		"name": c.Params("name"),
	}
	deleteResult, err := h.ApplicationModel.DeleteOne(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Publish to redis channel for driver to work on it
	errf := h.RedisClient.Publish(context.Background(), "application", c.Params("name")).Err()

	if errf != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(deleteResult)
}
