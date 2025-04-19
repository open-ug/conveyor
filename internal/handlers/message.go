package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	models "github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/internal/utils"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
`MessageHandler` is a struct that holds the models for the message handler
*/
type MessageHandler struct {
	MessageModel models.DriverMessageModel
	NatsCon      *nats.Conn
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(db *clientv3.Client, natsCon *nats.Conn) *MessageHandler {
	return &MessageHandler{
		MessageModel: models.DriverMessageModel{
			Client: db,
			Prefix: "driver-messages/",
		},
		NatsCon: natsCon,
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
	// Publish to nats channel for driver to work on it
	errf := h.NatsCon.Publish("application",
		jsonMsg,
	)
	if errf != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "message published",
	})
}
