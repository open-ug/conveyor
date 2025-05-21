/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	models "github.com/open-ug/conveyor/internal/models"
	craneTypes "github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationHandler struct {
	NatsCon          *nats.Conn
	ApplicationModel *models.ApplicationModel
}

func NewApplicationHandler(db *clientv3.Client, natsCon *nats.Conn) *ApplicationHandler {
	return &ApplicationHandler{
		NatsCon: natsCon,
		ApplicationModel: &models.ApplicationModel{
			Client: db,
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
	err := h.ApplicationModel.Insert(app)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("Saved to database")
	appMsg := craneTypes.ApplicationMsg{
		Action:  "create",
		Payload: app,
	}

	jsonMsg, merr := json.Marshal(appMsg)
	if merr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	driverMsg := craneTypes.DriverMessage{
		ID:      primitive.NewObjectID().Hex(),
		Payload: string(jsonMsg),
		Event:   "application",
		RunID:   uuid.New().String(),
	}

	jsonMsg, merr = json.Marshal(driverMsg)
	if merr != nil {
		fmt.Println(merr)
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
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"name":  app.Name,
		"runid": driverMsg.RunID,
	})
}

func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	app, err := h.ApplicationModel.FindOne(c.Params("name"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println(app)

	return c.Status(fiber.StatusOK).JSON(app)
}

func (h *ApplicationHandler) GetApplications(c *fiber.Ctx) error {
	apps, err := h.ApplicationModel.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(apps)
}

func (h *ApplicationHandler) UpdateApplication(c *fiber.Ctx) error {
	var appl craneTypes.Application

	if err := c.BodyParser(&appl); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}

	err := h.ApplicationModel.Update(appl)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newapp, err := h.ApplicationModel.FindOne(c.Params("name"))

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	appMsg := craneTypes.ApplicationMsg{
		Action:  "update",
		Payload: *newapp,
	}

	jsonMsg, merr := json.Marshal(appMsg)
	if merr != nil {
		fmt.Println(merr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	driverMsg := craneTypes.DriverMessage{
		ID:      primitive.NewObjectID().Hex(),
		Payload: string(jsonMsg),
		Event:   "application",
		RunID:   uuid.New().String(),
	}

	jsonMsg, merr = json.Marshal(driverMsg)
	if merr != nil {
		fmt.Println(merr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}
	errf := h.NatsCon.Publish("application", jsonMsg)

	if errf != nil {
		fmt.Println(errf)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "update successful",
		"runid":   driverMsg.RunID,
		"name":    newapp.Name,
	})
}

func (h *ApplicationHandler) DeleteApplication(c *fiber.Ctx) error {

	err := h.ApplicationModel.Delete(c.Params("name"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Publish to nats channel for driver to work on it
	errf := h.NatsCon.Publish("application", []byte(c.Params("name")))

	if errf != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete successful",
	})
}

func (h *ApplicationHandler) StartApplication(c *fiber.Ctx) error {

	app, err := h.ApplicationModel.FindOne(c.Params("name"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var appMsg craneTypes.ApplicationMsg
	appMsg.Payload = *app

	appMsg.Action = "start"
	appMsg.ID = appMsg.Payload.Name

	jsonMsg, merr := json.Marshal(appMsg)
	if merr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	driverMsg := craneTypes.DriverMessage{
		ID:      primitive.NewObjectID().Hex(),
		Payload: string(jsonMsg),
		Event:   "application",
		RunID:   uuid.New().String(),
	}

	jsonMsg, merr = json.Marshal(driverMsg)
	if merr != nil {
		fmt.Println(merr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	errf := h.NatsCon.Publish("application", jsonMsg)

	if errf != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Application started",
		"runid":   driverMsg.RunID,
		"name":    app.Name,
	})
}

func (h *ApplicationHandler) StopApplication(c *fiber.Ctx) error {

	app, err := h.ApplicationModel.FindOne(c.Params("name"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var appMsg craneTypes.ApplicationMsg
	appMsg.Payload = *app

	appMsg.Action = "stop"
	appMsg.ID = appMsg.Payload.Name

	jsonMsg, merr := json.Marshal(appMsg)
	if merr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	driverMsg := craneTypes.DriverMessage{
		ID:      primitive.NewObjectID().Hex(),
		Payload: string(jsonMsg),
		Event:   "application",
		RunID:   uuid.New().String(),
	}

	jsonMsg, merr = json.Marshal(driverMsg)
	if merr != nil {
		fmt.Println(merr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	errf := h.NatsCon.Publish("application", jsonMsg)

	if errf != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errf.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Application stopped",
		"runid":   driverMsg.RunID,
		"name":    app.Name,
	})
}
