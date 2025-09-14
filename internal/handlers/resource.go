package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/open-ug/conveyor/internal/engine"
	models "github.com/open-ug/conveyor/internal/models"
	utils "github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ResourceHandler struct {
	PipelineModel           *models.PipelineModel
	ResourceModel           *models.ResourceModel
	ResourceDefinitionModel *models.ResourceDefinitionModel
	NatsContext             *utils.NatsContext
}

func NewResourceHandler(db *clientv3.Client, natsContext *utils.NatsContext) *ResourceHandler {
	return &ResourceHandler{
		NatsContext: natsContext,
		ResourceModel: &models.ResourceModel{
			Client: db,
		},
		ResourceDefinitionModel: &models.ResourceDefinitionModel{
			Client: db,
		},
		PipelineModel: &models.PipelineModel{
			Client: db,
		},
	}
}

// CreateResource creates a new resource
// @Summary Create a new resource
// @Description Create a new resource with the specified type and configuration
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body types.Resource true "Resource object"
// @Success 201 {object} map[string]interface{} "Resource created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid payload or validation failed"
// @Failure 404 {object} map[string]interface{} "Resource definition not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resources [post]
func (h *ResourceHandler) CreateResource(c *fiber.Ctx) error {
	var resource types.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	resource.ID = uuid.New().String()
	resource.Metadata = make(map[string]string)
	// set version to 1
	resource.Metadata["version"] = "1"

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
	isValid, err := utils.ValidateResource(resource, resourceDef)
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

	if resource.Pipeline != "" {
		// Check if resource is part of a pipeline
		pipeline, err := h.PipelineModel.GetPipeline(resource.Pipeline)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to get pipeline: %v", err),
			})
		}

		if pipeline == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Pipeline %s not found", resource.Pipeline),
			})
		}

		if pipeline.Resource != resource.Resource {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Resource type %s does not match pipeline resource type %s", resource.Resource, pipeline.Resource),
			})
		}
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

	// Publish resource creation event to NATS JetStream
	run_id, err := engine.PublishResourceEvent("create", resource, h.NatsContext.JetStream)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to publish resource event: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"name":    resource.Name,
		"runid":   run_id,
		"message": "Resource created successfully",
		"version": resource.Metadata["version"],
	})
}

// GetResource retrieves a specific resource by name and type
// @Summary Get a resource
// @Description Retrieve a specific resource by its name and type
// @Tags resources
// @Accept json
// @Produce json
// @Param type path string true "Resource type"
// @Param name path string true "Resource name"
// @Success 200 {object} types.Resource "Resource object"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resources/{type}/{name} [get]
func (h *ResourceHandler) GetResource(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	resourceType := c.Params("type")
	if resourceName == "" || resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name and type are required",
		})
	}

	resource, err := h.ResourceModel.FindOne(resourceName, resourceType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to find resource: %v", err),
		})
	}

	return c.JSON(resource)
}

// DeleteResource deletes a specific resource by name and type
// @Summary Delete a resource
// @Description Delete a specific resource by its name and type
// @Tags resources
// @Accept json
// @Produce json
// @Param type path string true "Resource type"
// @Param name path string true "Resource name"
// @Success 204 {string} string "Resource deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resources/{type}/{name} [delete]
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

// ListResources lists all resources of a specific type
// @Summary List resources
// @Description List all resources of a specific type
// @Tags resources
// @Accept json
// @Produce json
// @Param type path string true "Resource type"
// @Success 200 {array} types.Resource "List of resources"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resources/{type} [get]
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

// UpdateResource updates a specific resource by name and type
// @Summary Update a resource
// @Description Update a specific resource by its name and type
// @Tags resources
// @Accept json
// @Produce json
// @Param type path string true "Resource type"
// @Param name path string true "Resource name"
// @Param resource body types.Resource true "Resource object"
// @Success 200 {object} types.Resource "Resource updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid payload or missing parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resources/{type}/{name} [put]
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
	isValid, err := utils.ValidateResource(resource, resourceDef)
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

	r, err := h.ResourceModel.Update(resourceName, resourceType, resource)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update resource: %v", err),
		})
	}

	return c.JSON(r)
}

// GetResourceByVersion retrieves a specific resource by name, type, and version
// @Summary Get a resource by version
// @Description Retrieve a specific resource by its name, type, and version
// @Tags resources
// @Accept json
// @Produce json
// @Param type path string true "Resource type"
// @Param name path string true "Resource name"
// @Param version path string true "Resource version"
// @Success 200 {object} types.Resource "Resource object"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing parameters"
// @Failure 404 {object} map[string]interface{} "Resource not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /resources/{type}/{name}/{version} [get]
func (h *ResourceHandler) GetResourceByVersion(c *fiber.Ctx) error {
	resourceName := c.Params("name")
	resourceType := c.Params("type")
	resourceVersion := c.Params("version")
	if resourceName == "" || resourceType == "" || resourceVersion == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Resource name, type, and version are required",
		})
	}

	resource, err := h.ResourceModel.FindByVersion(resourceName, resourceType, resourceVersion)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to find resource: %v", err),
		})
	}

	return c.JSON(resource)
}
