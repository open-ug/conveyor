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

// MessageHandler is a struct that holds the models for the message handler
type MessageHandler struct {
	MessageModel models.DriverMessageModel
	RedisClient  *redis.Client
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(db *mongo.Database, redisClient *redis.Client) *MessageHandler {
	return &MessageHandler{
		MessageModel: models.DriverMessageModel{
			Collection: db.Collection("messages"),
		},
		RedisClient: redisClient,
	}
}

// PublishMessage creates a new message
func (h *MessageHandler) PublishMessage(c *fiber.Ctx) error {
	var message craneTypes.DriverMessage
	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}
	insertResult, err := h.MessageModel.Insert(message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("Saved to database")
	message.ID = insertResult.InsertedID.(primitive.ObjectID).Hex()
	jsonMsg, merr := json.Marshal(message)
	if merr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}
	// Publish to redis channel for driver to work on it
	errf := h.RedisClient.Publish(context.Background(), "message",
		jsonMsg,
	)
	if errf.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Err().Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "message published",
	})
}
