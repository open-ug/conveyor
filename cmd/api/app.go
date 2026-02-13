/*
Copyright © 2024 - Present Conveyor CI Contributors
*/

// @title Conveyor CI API
// @version 0.5.0
// @description Conveyor is a lightweight, distributed CI/CD engine built for platform developers who demand simplicity without compromise.

// @contact.name Conveyor Support
// @contact.url https://conveyor.open.ug/
// @contact.email info@open.ug

// @license.name Apache 2.0
// @license.url https://opensource.org/license/apache-2-0

// @host localhost:8080
// @BasePath /
// @schemes http https

package api

import (
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	_ "github.com/open-ug/conveyor/internal/swagger"
	"github.com/open-ug/conveyor/pkg/server"
	"github.com/open-ug/conveyor/pkg/types"
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
	config := types.ServerConfig{}

	// transform port from string to int and set in config. Dont fetch from viper.

	intport, err := strconv.Atoi(port)
	if err != nil {
		color.Red("Invalid port number: %v", err)
		return
	}
	config.API.Port = intport
	config.API.AuthEnabled = viper.GetBool("api.auth_enabled")
	config.API.Data = viper.GetString("api.data")

	config.NATS.Port = viper.GetInt("nats.port")
	config.TLS.CA = viper.GetString("tls.ca")
	config.TLS.Key = viper.GetString("tls.key")
	config.TLS.Cert = viper.GetString("tls.cert")

	apiContext, err := server.Setup(&config)
	if err != nil {
		panic(err)
	}

	apiContext.App.Get("/health", healthHandler)

	apiContext.Start()

}
