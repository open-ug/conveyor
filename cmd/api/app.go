/*
Copyright © 2024 Beingana Jim Junior and Contributors
*/

// @title Conveyor API
// @version 1.0
// @description Conveyor is a platform for managing and orchestrating resources, drivers, and workflows.
// @termsOfService https://conveyor.open.ug/terms

// @contact.name Conveyor Support
// @contact.url https://conveyor.open.ug/support
// @contact.email support@conveyor.open.ug

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

package api

import (
	"time"

	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	metrics "github.com/open-ug/conveyor/internal/metrics"
	routes "github.com/open-ug/conveyor/internal/routes"
	_ "github.com/open-ug/conveyor/internal/swagger"
	utils "github.com/open-ug/conveyor/internal/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @Summary Health check
// @Description Get the health status of the API server
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health status"
// @Router /health [get]
func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}

func StartServer(port string) {

	app, err := Setup()
	if err != nil {
		color.Red("Error setting up the server: %v", err)
		return
	}
	app.Listen(":" + port)

}

// Setup Server
func Setup() (*fiber.App, error) {

	app := fiber.New(fiber.Config{
		AppName:     "Conveyor API Server",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New())

	// Start metrics updater
	metrics.StartAPIMetricsUpdater()

	// Add Prometheus middleware
	app.Use(metrics.PrometheusMiddleware())

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("CONVEYOR API SERVER. Visit https://conveyor.open.ug for Documentation")
	})

	app.Get("/health", healthHandler)

	// Metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	natsContext := utils.NewNatsConn()
	natsContext.InitiateStreams()

	etcd, err := utils.NewEtcdClient()

	if err != nil {
		color.Red("Error Occured while creating etcd client: %v", err)
		return nil, err
	}

	routes.DriverRoutes(app, etcd.Client, natsContext.NatsCon)
	routes.ResourceRoutes(app, etcd.Client, natsContext)

	return app, nil
}
