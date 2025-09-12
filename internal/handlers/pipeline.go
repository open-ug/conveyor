package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type PipelineHandler struct {
	Model                   *models.PipelineModel
	NatsCon                 *nats.Conn
	ResourceDefinitionModel *models.ResourceDefinitionModel
}

func NewPipelineHandler(db *clientv3.Client, natsCon *nats.Conn) *PipelineHandler {
	return &PipelineHandler{
		Model:                   models.NewPipelineModel(db),
		NatsCon:                 natsCon,
		ResourceDefinitionModel: models.NewResourceDefinitionModel(db),
	}
}

// CreatePipeline creates a new pipeline
// @Summary Create a new pipeline
// @Description Create a new pipeline with the specified configuration
// @Tags pipelines
// @Accept json
// @Produce json
// @Param pipeline body types.Pipeline true "Pipeline object"
// @Success 201 {object} types.Pipeline "Pipeline created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid payload"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pipelines [post]
func (h *PipelineHandler) CreatePipeline(c *fiber.Ctx) error {
	var pipeline types.Pipeline
	if err := c.BodyParser(&pipeline); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	// Validate that the resource definition exists
	resourceDefinition, err := h.ResourceDefinitionModel.FindOne(pipeline.Resource)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to find resource definition: %v", err),
		})
	}
	if resourceDefinition == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Resource definition for type %s not found", pipeline.Resource),
		})
	}

	err = h.Model.CreatePipeline(&pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pipeline",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(pipeline)
}

// GetPipeline retrieves a pipeline by name
// @Summary Get a pipeline by name
// @Description Retrieve a pipeline by its name
// @Tags pipelines
// @Accept json
// @Produce json
// @Param name path string true "Pipeline name"
// @Success 200 {object} types.Pipeline "Pipeline retrieved successfully"
// @Failure 404 {object} map[string]interface{} "Not found - Pipeline does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pipelines/{name} [get]
func (h *PipelineHandler) GetPipeline(c *fiber.Ctx) error {
	name := c.Params("name")
	pipeline, err := h.Model.GetPipeline(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pipeline not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(pipeline)
}

// UpdatePipeline updates an existing pipeline
// @Summary Update an existing pipeline
// @Description Update an existing pipeline with new configuration
// @Tags pipelines
// @Accept json
// @Produce json
// @Param name path string true "Pipeline name"
// @Param pipeline body types.Pipeline true "Updated pipeline object"
// @Success 200 {object} types.Pipeline "Pipeline updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid payload"
// @Failure 404 {object} map[string]interface{} "Not found - Pipeline does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pipelines/{name} [put]
func (h *PipelineHandler) UpdatePipeline(c *fiber.Ctx) error {
	name := c.Params("name")
	var pipeline types.Pipeline
	if err := c.BodyParser(&pipeline); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if name != pipeline.Name {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Pipeline name in URL and body do not match",
		})
	}
	err := h.Model.UpdatePipeline(&pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update pipeline",
		})
	}
	return c.Status(fiber.StatusOK).JSON(pipeline)
}

// DeletePipeline deletes a pipeline by name
// @Summary Delete a pipeline by name
// @Description Delete a pipeline by its name
// @Tags pipelines
// @Accept json
// @Produce json
// @Param name path string true "Pipeline name"
// @Success 204 "Pipeline deleted successfully"
// @Failure 404 {object} map[string]interface{} "Not found - Pipeline does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pipelines/{name} [delete]
func (h *PipelineHandler) DeletePipeline(c *fiber.Ctx) error {
	name := c.Params("name")
	err := h.Model.DeletePipeline(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pipeline not found",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
