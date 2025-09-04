package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/handlers"
)

func LogRoutes(app *fiber.App, logHandler *handlers.LogHandler) {

	app.Post("/logs", logHandler.CreateLog)

	// GET /logs?pipeline=...&driver=...&runid=...
	app.Get("/logs", logHandler.GetLogs)
}
