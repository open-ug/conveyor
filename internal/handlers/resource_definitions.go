package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	models "github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ResourceDefinitionHandler struct {
	NatsCon                 *nats.Conn
	ResourceDefinitionModel *models.ResourceDefinitionModel
}

func NewResourceDefinitionHandler(db *clientv3.Client, natsCon *nats.Conn) *ResourceDefinitionHandler {
	return &ResourceDefinitionHandler{
		NatsCon: natsCon,
		ResourceDefinitionModel: &models.ResourceDefinitionModel{
			Client: db,
		},
	}
}

func (h *ResourceDefinitionHandler) CreateResourceDefinition(c *fiber.Ctx) error {
	var resourceDefinition types.ResourceDefinition
	if err := c.BodyParser(&resourceDefinition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	resourceDefinition.ID = uuid.New().String()

	resourceDef, err := json.Marshal(resourceDefinition)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal resource definition",
		})
	}

	err = h.ResourceDefinitionModel.Insert(resourceDefinition.Name, resourceDef)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to insert resource definition: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resourceDefinition)
}

func (h *ResourceDefinitionHandler) GetResourceDefinition(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	if resourceName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name is required",
		})
	}

	resourceDef, err := h.ResourceDefinitionModel.FindOne(resourceName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to find resource definition: %v", err),
		})
	}

	var resourceDefinition types.ResourceDefinition
	err = json.Unmarshal(resourceDef, &resourceDefinition)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unmarshal resource definition",
		})
	}

	return c.JSON(resourceDefinition)
}

func (h *ResourceDefinitionHandler) DeleteResourceDefinition(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	if resourceName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name is required",
		})
	}

	err := h.ResourceDefinitionModel.Delete(resourceName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete resource definition: %v", err),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ResourceDefinitionHandler) UpdateResourceDefinition(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	if resourceName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name is required",
		})
	}

	var resourceDefinition types.ResourceDefinition
	if err := c.BodyParser(&resourceDefinition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	resourceDef, err := json.Marshal(resourceDefinition)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal resource definition",
		})
	}

	err = h.ResourceDefinitionModel.Update(resourceName, resourceDef)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update resource definition: %v", err),
		})
	}

	return c.JSON(resourceDefinition)
}
