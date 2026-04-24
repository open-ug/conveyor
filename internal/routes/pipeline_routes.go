/*
Copyright © 2024 - Present Conveyor CI Contributors
*/
package routes

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/open-ug/conveyor/internal/handlers"
	utils "github.com/open-ug/conveyor/internal/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func PipelineRoutes(app *fiber.App, cli *clientv3.Client, natsContext *utils.NatsContext, db *badger.DB) {

	// Initialize pipeline handler
	pipelinePrefix := app.Group("/pipelines")
	pipelineHandler := handlers.NewPipelineHandler(cli, natsContext.NatsCon, db)
	// Define routes
	pipelinePrefix.Post("/", pipelineHandler.CreatePipeline)

}
