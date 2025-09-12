/*
Copyright © 2024 Conveyor CI Contributors
*/
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/handlers"
	utils "github.com/open-ug/conveyor/internal/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func PipelineRoutes(app *fiber.App, db *clientv3.Client, natsContext *utils.NatsContext) {

	// Initialize pipeline handler
	pipelinePrefix := app.Group("/pipelines")
	pipelineHandler := handlers.NewPipelineHandler(db, natsContext.NatsCon)

	// Define routes
	pipelinePrefix.Post("/", pipelineHandler.CreatePipeline)

}
