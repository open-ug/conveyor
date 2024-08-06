package handlers

import (
	"fmt"

	models "crane.cloud.cranom.tech/cmd/api/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationHandler struct {
	ApplicationModel *models.ApplicationModel
}

func NewApplicationHandler(db *mongo.Database) *ApplicationHandler {
	return &ApplicationHandler{
		ApplicationModel: &models.ApplicationModel{
			Collection: db.Collection("applications"),
		},
	}
}

func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	var app models.Application
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
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"insertedId": insertResult.InsertedID,
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
	var update map[string]interface{}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}
	updateResult, err := h.ApplicationModel.UpdateOne(filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
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
	return c.Status(fiber.StatusOK).JSON(deleteResult)
}
