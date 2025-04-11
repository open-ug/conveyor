package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/utils"
)

type MetricsFilters struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (h *ApplicationHandler) GetApplicationCPUUsage(c *fiber.Ctx) error {
	filter := map[string]interface{}{
		"name": c.Params("name"),
	}
	app, err := h.ApplicationModel.FindOne(filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var data MetricsFilters
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}

	res, err := utils.GetCPUUsage(app.Name, data.Start, data.End)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *ApplicationHandler) GetApplicationMemoryUsage(c *fiber.Ctx) error {
	filter := map[string]interface{}{
		"name": c.Params("name"),
	}
	app, err := h.ApplicationModel.FindOne(filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var data MetricsFilters
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}

	res, err := utils.GetMemoryUsage(app.Name, data.Start, data.End)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse request body",
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
