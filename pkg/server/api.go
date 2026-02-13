package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/open-ug/conveyor/internal/config/auth"
	"github.com/open-ug/conveyor/internal/engine"
	"github.com/open-ug/conveyor/internal/handlers"
	"github.com/open-ug/conveyor/internal/metrics"
	"github.com/open-ug/conveyor/internal/models"
	"github.com/open-ug/conveyor/internal/routes"
	"github.com/open-ug/conveyor/internal/utils"
	"github.com/open-ug/conveyor/pkg/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

/* APIServerContext holds references to the API server and its dependencies. You can call Start() on the APIServerContext to run the server, and ShutDown() to gracefully stop it.
 */
type APIServerContext struct {
	// The API server instance. This is a Fiber app that you extend with your routes and middleware. You can also access it to add custom routes or middleware if needed.
	App *fiber.App
	// NATS context holds the NATS connection and JetStream context. You can use this to publish messages to NATS or manage streams.
	NatsContext *utils.NatsContext
	// ETCD client for interacting with the ETCD datastore. You can use this to read/write data to ETCD directly if needed.
	ETCD *utils.EtcdClient
	// LogModel is a wrapper around BadgerDB for storing logs. You can use this to write logs to the database.
	LogModel *models.LogModel
	// BadgerDB instance for direct access to the database if needed. You can use this for advanced queries or operations that are not covered by the LogModel.
	BadgerDB *badger.DB
	// Config holds the server configuration. You can access this to read any configuration values that were used to set up the server. This can be useful if you want to make decisions based on the configuration at runtime.
	Config *types.ServerConfig
}

/*
Graceful shutdown of server and dependencies
*/
func (c *APIServerContext) ShutDown() {
	c.App.Shutdown()
	c.NatsContext.Shutdown()
	c.ETCD.ServerStop()
	if c.BadgerDB != nil {
		c.BadgerDB.Close()
	}
}

/*
Setup the API server with all routes and dependencies. It returns an APIServerContext which holds references to the server. You can later call Start() on the APIServerContext to run the server.
*/
func Setup(config *types.ServerConfig) (APIServerContext, error) {

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

	if config.API.AuthEnabled {
		// Auth middleware
		rootPool, err := auth.LoadRootCAs(config)
		if err != nil {
			color.Red("Error loading root CAs: %v", err)
			return APIServerContext{}, err
		}
		app.Use(auth.JWTCertMiddleware(rootPool))
	}

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("CONVEYOR API SERVER. Visit https://conveyor.open.ug for Documentation")
	})

	// Metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	natsContext := utils.NewNatsConn(config)
	natsContext.InitiateStreams()

	etcd, err := utils.NewEtcdClient(config)
	if err != nil {
		color.Red("Error Occured while creating etcd client: %v", err)
		return APIServerContext{}, err
	}

	// Initialize BadgerDB
	conveyorDataDir := config.API.Data
	badgerOpts := badger.DefaultOptions(conveyorDataDir + "/badger")
	badgerDB, err := badger.Open(badgerOpts)
	if err != nil {
		color.Red("Error opening BadgerDB: %v", err)
		return APIServerContext{}, err
	}

	// Register routes
	logModel := &models.LogModel{DB: badgerDB}
	routes.LogRoutes(app, &handlers.LogHandler{Model: logModel}, natsContext)

	routes.DriverRoutes(app, etcd.Client, natsContext.NatsCon)
	routes.ResourceRoutes(app, etcd.Client, natsContext)
	routes.PipelineRoutes(app, etcd.Client, natsContext)

	return APIServerContext{
		NatsContext: natsContext,
		App:         app,
		ETCD:        etcd,
		LogModel:    logModel,
		BadgerDB:    badgerDB,
		Config:      config,
	}, nil
}

/*
Start the API server and all dependencies. This function will block until an interrupt signal is received, at which point it will attempt to gracefully shut down the server and dependencies.
*/
func (appCtx *APIServerContext) Start() {

	// Run server in a goroutine
	go func() {
		if err := appCtx.App.Listen(":" + fmt.Sprintf("%d", appCtx.Config.API.Port)); err != nil {
			fmt.Printf("Server stopped: %v\n", err)
		}
	}()

	engineCtx := engine.NewEngineContext(appCtx.ETCD.Client, appCtx.LogModel, *appCtx.NatsContext)

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
	fmt.Println("Gracefully shutting down Datastore")
	appCtx.ETCD.ServerStop()

	fmt.Println("Server gracefully stopped")

}
