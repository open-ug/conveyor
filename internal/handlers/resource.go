package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/open-ug/conveyor/internal/helpers"
	models "github.com/open-ug/conveyor/internal/models"
	internals "github.com/open-ug/conveyor/internal/shared"
	utils "github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ResourceHandler struct {
	ResourceModel           *models.ResourceModel
	ResourceDefinitionModel *models.ResourceDefinitionModel
	NatsContext             *internals.NatsContext
}

func NewResourceHandler(db *clientv3.Client, natsContext *internals.NatsContext) *ResourceHandler {
	return &ResourceHandler{
		NatsContext: natsContext,
		ResourceModel: &models.ResourceModel{
			Client: db,
		},
		ResourceDefinitionModel: &models.ResourceDefinitionModel{
			Client: db,
		},
	}
}

func (h *ResourceHandler) CreateResource(c *fiber.Ctx) error {
	var resource types.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	resource.ID = uuid.New().String()

	resourceType := resource.Resource
	if resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource type is required",
		})
	}

	resourceDefinition, err := h.ResourceDefinitionModel.FindOne(resourceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to find resource definition: %v", err),
		})
	}
	if resourceDefinition == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Resource definition for type %s not found", resourceType),
		})
	}

	var resourceDef types.ResourceDefinition
	err = json.Unmarshal(resourceDefinition, &resourceDef)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unmarshal resource definition",
		})
	}
	if resourceDef.Schema == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource definition schema is required",
		})
	}

	// Validate the resource against the schema
	isValid, err := helpers.ValidateResource(resource, resourceDef)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Resource validation failed: %v", err),
		})
	}
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource does not conform to the schema",
		})
	}

	resourceData, err := json.Marshal(resource)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal resource",
		})
	}

	err = h.ResourceModel.Insert(resource.Name, resource.Resource, resourceData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to insert resource: %v", err),
		})
	}

	mID, err := utils.GenerateRandomID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generating id:",
		})
	}
	driverMsg := types.DriverMessage{
		ID:      mID,
		Payload: string(resourceData),
		Event:   "create",
		RunID:   uuid.New().String(),
	}

	jsonMsg, merr := json.Marshal(driverMsg)
	if merr != nil {
		fmt.Println(merr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": merr.Error(),
		})
	}

	// Publish message to jetstream
	subjectName := "resources." + resourceDef.Name
	_, err = h.NatsContext.JetStream.PublishAsync(subjectName, jsonMsg)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resource)
}

func (h *ResourceHandler) GetResource(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	resourceType := c.Params("type")
	if resourceName == "" || resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name and type are required",
		})
	}

	resourceData, err := h.ResourceModel.FindOne(resourceName, resourceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to find resource: %v", err),
		})
	}

	var resource types.Resource
	err = json.Unmarshal(resourceData, &resource)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unmarshal resource data",
		})
	}

	return c.JSON(resource)
}

func (h *ResourceHandler) DeleteResource(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	resourceType := c.Params("type")
	if resourceName == "" || resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name and type are required",
		})
	}

	err := h.ResourceModel.Delete(resourceName, resourceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete resource: %v", err),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ResourceHandler) ListResources(c *fiber.Ctx) error {
	resourceType := c.Params("type")
	if resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource type is required",
		})
	}

	resources, err := h.ResourceModel.List(resourceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to list resources: %v", err),
		})
	}

	return c.JSON(resources)
}

func (h *ResourceHandler) UpdateResource(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	resourceType := c.Params("type")
	if resourceName == "" || resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name and type are required",
		})
	}

	var resource types.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	resourceData, err := json.Marshal(resource)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal resource",
		})
	}
	resourceDefinition, err := h.ResourceDefinitionModel.FindOne(resourceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to find resource definition: %v", err),
		})
	}

	if resourceDefinition == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Resource definition for type %s not found", resourceType),
		})
	}

	var resourceDef types.ResourceDefinition
	err = json.Unmarshal(resourceDefinition, &resourceDef)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unmarshal resource definition",
		})
	}

	if resourceDef.Schema == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource definition schema is required",
		})
	}

	// Validate the resource against the schema
	isValid, err := helpers.ValidateResource(resource, resourceDef)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Resource validation failed: %v", err),
		})
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource does not conform to the schema",
		})
	}

	err = h.ResourceModel.Update(resourceName, resourceType, resourceData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update resource: %v", err),
		})
	}

	return c.JSON(resource)
}
