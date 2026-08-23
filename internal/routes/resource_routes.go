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

func ResourceRoutes(app *fiber.App, natsContext *utils.NatsContext, db *badger.DB) {

	// Initialize resource and resource definition handlers
	resourcePrefix := app.Group("/resources")
	resourceDefinitionPrefix := app.Group("/resource-definitions")
	resourceHandler := handlers.NewResourceHandler(natsContext, db)
	resourceDefinitionHandler := handlers.NewResourceDefinitionHandler(natsContext.NatsCon, db)

	// Resource Routes
	resourcePrefix.Post("/", resourceHandler.CreateResource)
	resourcePrefix.Get("/:type/:name", resourceHandler.GetResource)
	resourcePrefix.Delete("/:type/:name", resourceHandler.DeleteResource)
	resourcePrefix.Put("/:type/:name", resourceHandler.UpdateResource)
	resourcePrefix.Get("/:type/:name/:version", resourceHandler.GetResourceByVersion)

	resourcePrefix.Get("/", resourceHandler.ListResources)

	// Resource Definition Routes
	resourceDefinitionPrefix.Post("/", resourceDefinitionHandler.CreateResourceDefinition)
	resourceDefinitionPrefix.Post("/apply", resourceDefinitionHandler.CreateOrUpdateResourceDefinition)
	resourceDefinitionPrefix.Get("/:name", resourceDefinitionHandler.GetResourceDefinition)
	resourceDefinitionPrefix.Delete("/:name", resourceDefinitionHandler.DeleteResourceDefinition)
	resourceDefinitionPrefix.Put("/:name", resourceDefinitionHandler.UpdateResourceDefinition)

}
