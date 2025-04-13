package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	models "github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/internal/utils"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
`MessageHandler` is a struct that holds the models for the message handler
*/
type MessageHandler struct {
	MessageModel models.DriverMessageModel
	RedisClient  *redis.Client
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(db *clientv3.Client, redisClient *redis.Client) *MessageHandler {
	return &MessageHandler{
		MessageModel: models.DriverMessageModel{
			Client: db,
			Prefix: "driver-messages/",
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
	mid, err := utils.GenerateRandomID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	message.ID = mid
	err = h.MessageModel.Insert(message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("Saved to database")
	jsonMsg, merr := json.Marshal(message)
	if merr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}
	// Publish to redis channel for driver to work on it
	errf := h.RedisClient.Publish(context.Background(), "application",
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
