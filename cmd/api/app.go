/*
Copyright © 2024 - Present Conveyor CI Contributors
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
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/fatih/color"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/open-ug/conveyor/internal/auth"
	"github.com/open-ug/conveyor/internal/engine"
	"github.com/open-ug/conveyor/internal/handlers"
	"github.com/spf13/viper"

	metrics "github.com/open-ug/conveyor/internal/metrics"
	"github.com/open-ug/conveyor/internal/models"
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

	// Check if TLS is enabled
	tlsEnabled := viper.GetBool("auth.tls.enabled")
	
	// Run server in a goroutine
	go func() {
		var err error
		if tlsEnabled {
			// Create TLS configuration
			conveyorDataDir := viper.GetString("api.data")
			certDir := filepath.Join(conveyorDataDir, "certs")
			
			certManager := auth.NewCertificateManager(certDir)
			tlsConfig := auth.NewTLSConfig(certManager)
			serverTLSConfig, err := tlsConfig.CreateServerTLSConfig()
			if err != nil {
				color.Red("Error creating TLS configuration: %v", err)
				return
			}
			
			// Create TLS listener
			ln, err := tls.Listen("tcp", ":"+port, serverTLSConfig)
			if err != nil {
				color.Red("Error creating TLS listener: %v", err)
				return
			}
			
			color.Green("Starting HTTPS server on port %s with TLS authentication", port)
			err = appCtx.App.Listener(ln)
		} else {
			color.Yellow("Warning: Starting HTTP server without TLS (auth.tls.enabled=false)")
			err = appCtx.App.Listen(":" + port)
		}
		
		if err != nil {
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

type APIServerContext struct {
	App         *fiber.App
	NatsContext *utils.NatsContext
	ETCD        *utils.EtcdClient
	LogModel    *models.LogModel
	BadgerDB    *badger.DB
	AuthManager *AuthManager // Add authentication manager
}

// AuthManager holds authentication components
type AuthManager struct {
	CertificateManager *auth.CertificateManager
	JWTManager         *auth.JWTManager
	TLSConfig          *auth.TLSConfig
}

func (c *APIServerContext) ShutDown() {
	c.App.Shutdown()
	c.NatsContext.Shutdown()
	c.ETCD.ServerStop()
	if c.BadgerDB != nil {
		c.BadgerDB.Close()
	}
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

	// Initialize authentication if enabled
	var authManager *AuthManager
	authEnabled := viper.GetBool("auth.enabled")
	if authEnabled {
		conveyorDataDir := viper.GetString("api.data")
		certDir := filepath.Join(conveyorDataDir, "certs")
		
		certManager := auth.NewCertificateManager(certDir)
		jwtManager := auth.NewJWTManager(certManager)
		tlsConfig := auth.NewTLSConfig(certManager)
		
		authManager = &AuthManager{
			CertificateManager: certManager,
			JWTManager:         jwtManager,
			TLSConfig:          tlsConfig,
		}
		
		// Add authentication middleware to protected routes
		jwtRequired := viper.GetBool("auth.jwt.required")
		if jwtRequired {
			// Apply authentication middleware to all routes except public ones
			authMiddleware := auth.NewAuthMiddleware(jwtManager, certManager)
			authMiddleware.SetSkipAuthRoutes([]string{"/health", "/metrics", "/swagger", "/docs", "/"})
			app.Use(authMiddleware.Handler())
		} else {
			// Use optional authentication
			app.Use(auth.OptionalAuth(jwtManager, certManager))
		}
		
		color.Green("Authentication enabled with TLS certificates and JWT tokens")
	} else {
		color.Yellow("Warning: Authentication is disabled (auth.enabled=false)")
	}

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

	// Initialize BadgerDB
	conveyorDataDir := viper.GetString("api.data")
	badgerOpts := badger.DefaultOptions(conveyorDataDir + "/badger")
	badgerDB, err := badger.Open(badgerOpts)
	if err != nil {
		color.Red("Error opening BadgerDB: %v", err)
		return APIServerContext{}, err
	}

	// Register routes
	logModel := &models.LogModel{DB: badgerDB}
	routes.LogRoutes(app, &handlers.LogHandler{Model: logModel}, natsContext.NatsCon)

	routes.DriverRoutes(app, etcd.Client, natsContext.NatsCon)
	routes.ResourceRoutes(app, etcd.Client, natsContext)
	routes.PipelineRoutes(app, etcd.Client, natsContext)

	return APIServerContext{
		NatsContext: natsContext,
		App:         app,
		ETCD:        etcd,
		LogModel:    logModel,
		BadgerDB:    badgerDB,
		AuthManager: authManager,
	}, nil
}
