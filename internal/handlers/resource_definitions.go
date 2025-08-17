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

// CreateResourceDefinition creates a new resource definition
// @Summary Create a resource definition
// @Description Create a new resource definition with the specified configuration
// @Tags resource-definitions
// @Accept json
// @Produce json
// @Param resourceDefinition body types.ResourceDefinition true "Resource definition object"
// @Success 201 {object} types.ResourceDefinition "Resource definition created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid payload"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resource-definitions [post]
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

// GetResourceDefinition retrieves a specific resource definition by name
// @Summary Get a resource definition
// @Description Retrieve a specific resource definition by its name
// @Tags resource-definitions
// @Accept json
// @Produce json
// @Param name path string true "Resource definition name"
// @Success 200 {object} types.ResourceDefinition "Resource definition object"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resource-definitions/{name} [get]
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

// DeleteResourceDefinition deletes a specific resource definition by name
// @Summary Delete a resource definition
// @Description Delete a specific resource definition by its name
// @Tags resource-definitions
// @Accept json
// @Produce json
// @Param name path string true "Resource definition name"
// @Success 204 {string} string "Resource definition deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resource-definitions/{name} [delete]
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

// UpdateResourceDefinition updates a specific resource definition by name
// @Summary Update a resource definition
// @Description Update a specific resource definition by its name
// @Tags resource-definitions
// @Accept json
// @Produce json
// @Param name path string true "Resource definition name"
// @Param resourceDefinition body types.ResourceDefinition true "Resource definition object"
// @Success 200 {object} types.ResourceDefinition "Resource definition updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid payload or missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resource-definitions/{name} [put]
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

// CreateOrUpdateResourceDefinition creates or updates a resource definition
// @Summary Create or update a resource definition
// @Description Create a new resource definition or update existing one by name
// @Tags resource-definitions
// @Accept json
// @Produce json
// @Param resourceDefinition body types.ResourceDefinition true "Resource definition object"
// @Success 201 {object} types.ResourceDefinition "Resource definition created or updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid payload or missing name"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resource-definitions/apply [post]
func (h *ResourceDefinitionHandler) CreateOrUpdateResourceDefinition(c *fiber.Ctx) error {
	fmt.Println("CreateOrUpdateResourceDefinition called")
	var resourceDefinition types.ResourceDefinition
	if err := c.BodyParser(&resourceDefinition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	resourceName := resourceDefinition.Name
	if resourceName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name is required",
		})
	}
	resourceDef, err := json.Marshal(resourceDefinition)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal resource definition",
		})
	}

	// Check if the resource definition already exists
	existingDef, _ := h.ResourceDefinitionModel.FindOne(resourceName)
	fmt.Println("resource exists")

	if existingDef != nil {
		// If it exists, update it
		err = h.ResourceDefinitionModel.Update(resourceName, resourceDef)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to update resource definition: %v", err),
			})
		}
	} else {
		// If it doesn't exist, create it
		resourceDefinition.ID = uuid.New().String()
		err = h.ResourceDefinitionModel.Insert(resourceName, resourceDef)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to insert resource definition: %v", err),
			})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(resourceDefinition)
}
