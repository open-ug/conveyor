package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/handlers"
	"github.com/open-ug/conveyor/internal/streaming"
	utils "github.com/open-ug/conveyor/internal/utils"
)

func LogRoutes(app *fiber.App, logHandler *handlers.LogHandler, natsContext *utils.NatsContext) {

	logStreamHandler := streaming.NewPipelineLogsStreamer(natsContext, logHandler.Model)

	app.Post("/logs", logHandler.CreateLog)

	// GET /logs?pipeline=...&driver=...&runid=...
	app.Get("/logs", logHandler.GetLogs)

	// Streaming logs
	app.Get("/logs/streams/:drivername/:runid", streaming.NewDriverLogsStreamer(natsContext.NatsCon, logHandler.Model).StreamDriverLogsByRunID)

	app.Get("/logs/pipeline/:runid", logStreamHandler.StreamLogsByRunID)
}
