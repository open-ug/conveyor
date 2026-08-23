/*
Copyright © 2024 - Present Conveyor CI Contributors
*/
package routes

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/handlers"
	utils "github.com/open-ug/conveyor/internal/utils"
)

func PipelineRoutes(app *fiber.App, natsContext *utils.NatsContext, db *badger.DB) {

	// Initialize pipeline handler
	pipelinePrefix := app.Group("/pipelines")
	pipelineHandler := handlers.NewPipelineHandler(natsContext.NatsCon, db)
	// Define routes
	pipelinePrefix.Post("/", pipelineHandler.CreatePipeline)

}
