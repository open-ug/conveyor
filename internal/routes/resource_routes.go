/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/open-ug/conveyor/internal/handlers"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func ResourceRoutes(app *fiber.App, db *clientv3.Client, natsCon *nats.Conn) {
	resourcePrefix := app.Group("/resources")
	resourceDefinitionPrefix := app.Group("/resource-definitions")
	resourceHandler := handlers.NewResourceHandler(db, natsCon)
	resourceDefinitionHandler := handlers.NewResourceDefinitionHandler(db, natsCon)

	// Resource Routes
	resourcePrefix.Post("/", resourceHandler.CreateResource)
	resourcePrefix.Get("/:name/:type", resourceHandler.GetResource)
	resourcePrefix.Delete("/:name/:type", resourceHandler.DeleteResource)
	resourcePrefix.Put("/:name/:type", resourceHandler.UpdateResource)

	resourcePrefix.Get("/", resourceHandler.ListResources)

	// Resource Definition Routes
	resourceDefinitionPrefix.Post("/", resourceDefinitionHandler.CreateResourceDefinition)
	resourceDefinitionPrefix.Get("/:name", resourceDefinitionHandler.GetResourceDefinition)
	resourceDefinitionPrefix.Delete("/:name", resourceDefinitionHandler.DeleteResourceDefinition)
	resourceDefinitionPrefix.Put("/:name", resourceDefinitionHandler.UpdateResourceDefinition)

}
