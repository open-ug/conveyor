package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/handlers"
	"github.com/open-ug/conveyor/internal/streaming"
)

func LogRoutes(app *fiber.App, logHandler *handlers.LogHandler, natsCon *nats.Conn) {

	app.Post("/logs", logHandler.CreateLog)

	// GET /logs?pipeline=...&driver=...&runid=...
	app.Get("/logs", logHandler.GetLogs)

	// Streaming logs
	app.Get("/logs/streams/:name/:runid", websocket.New(streaming.NewDriverLogsStreamer(natsCon, logHandler.Model).StreamLogs))
}
