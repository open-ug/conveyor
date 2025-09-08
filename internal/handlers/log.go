package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/pkg/types"
)

type LogHandler struct {
	Model *models.LogModel
}

// CreateLog persists a new log entry in BadgerDB
// @Summary Create a log entry
// @Description Accepts a JSON log entry and stores it in BadgerDB
// @Tags logs
// @Accept json
// @Produce json
// @Param log body types.Log true "Log entry object"
// @Success 201 {string} string "Log created"
// @Failure 400 {object} map[string]string "Bad request - Invalid payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs [post]
func (h *LogHandler) CreateLog(c *fiber.Ctx) error {
	var log types.Log
	if err := c.BodyParser(&log); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.Model.Insert(log); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

// GetLogs retrieves log entries from BadgerDB with optional filtering
// @Summary Get log entries
// @Description Returns logs filtered by pipeline, driver, and runid
// @Tags logs
// @Accept json
// @Produce json
// @Param pipeline query string false "Pipeline name"
// @Param driver query string false "Driver name"
// @Param runid query string false "Run ID"
// @Success 200 {array} types.Log "List of logs"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs [get]
func (h *LogHandler) GetLogs(c *fiber.Ctx) error {
	pipeline := c.Query("pipeline")
	driver := c.Query("driver")
	runid := c.Query("runid")

	logs, err := h.Model.Query(pipeline, driver, runid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(logs)
}
