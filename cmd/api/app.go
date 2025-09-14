/*
Copyright © 2024 Conveyor CI Contributors
*/

// @title Conveyor CI API
// @version 0.1.31
// @description Conveyor is a lightweight, distributed CI/CD engine built for platform developers who demand simplicity without compromise.

// @contact.name Conveyor Support
// @contact.url https://conveyor.open.ug/
// @contact.email conveyor@open.ug

// @license.name Apache 2.0
// @license.url https://opensource.org/license/apache-2-0

// @host localhost:8080
// @BasePath /
// @schemes http https

package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/open-ug/conveyor/internal/engine"
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

	appCtx, err := Setup()
	if err != nil {
		color.Red("Error setting up the server: %v", err)
		return
	}

	// Run server in a goroutine
	go func() {
		if err := appCtx.App.Listen(":" + port); err != nil {
			fmt.Printf("Server stopped: %v\n", err)
		}
	}()

	engineCtx := engine.NewEngineContext(appCtx.ETCD.Client, *appCtx.NatsContext)

	go func() {
		err := engineCtx.Start()
		if err != nil {
			color.Red("Error starting the engine: %v", err)
			return
		}
	}()

	// Setup channel to listen for interrupt/terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Wait until a signal is received

	fmt.Println("Shutting down server...")

	// Create a context with timeout for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := appCtx.App.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	appCtx.NatsContext.Shutdown()
	fmt.Println("Sracefully shutting down Datastore")
	appCtx.ETCD.ServerStop()

	fmt.Println("Server gracefully stopped")

}

type APIServerContext struct {
	App         *fiber.App
	NatsContext *utils.NatsContext
	ETCD        *utils.EtcdClient
}

func (c *APIServerContext) ShutDown() {
	c.App.Shutdown()
	c.NatsContext.Shutdown()
	c.ETCD.ServerStop()
}

// Setup Server
func Setup() (APIServerContext, error) {

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
		return APIServerContext{}, err
	}

	routes.DriverRoutes(app, etcd.Client, natsContext.NatsCon)
	routes.ResourceRoutes(app, etcd.Client, natsContext)

	return APIServerContext{
		NatsContext: natsContext,
		App:         app,
		ETCD:        etcd,
	}, nil
}
