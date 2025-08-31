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

func ResourceRoutes(app *fiber.App, db *clientv3.Client, natsContext *utils.NatsContext) {

	// Initialize resource and resource definition handlers
	resourcePrefix := app.Group("/resources")
	resourceDefinitionPrefix := app.Group("/resource-definitions")
	resourceHandler := handlers.NewResourceHandler(db, natsContext)
	resourceDefinitionHandler := handlers.NewResourceDefinitionHandler(db, natsContext.NatsCon)

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
