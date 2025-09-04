package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/pkg/types"
)

type LogHandler struct {
	Model *models.LogModel
}

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
