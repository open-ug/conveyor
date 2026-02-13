package config

import (
	"fmt"

	"github.com/open-ug/conveyor/pkg/types"
)

type ServerConfigWrapper struct {
	*types.ServerConfig
}

func (c *ServerConfigWrapper) GenerateTestYamlConfig() string {
	return fmt.Sprintf(`# Conveyor CI Configuration
# Generated for testing

api:
		port: %d
		auth_enabled: %t
		data: %s
nats:
		port: %d
tls:
		ca: %s
		key: %s
		cert: %s
`,
		c.API.Port,
		c.API.AuthEnabled,
		c.API.Data,
		c.NATS.Port,
		c.TLS.CA,
		c.TLS.Key,
		c.TLS.Cert,
	)
}

func (c *ServerConfigWrapper) GetTestConfig() *types.ServerConfig {
	return c.ServerConfig
}
