/*
Copyright © 2024 Conveyor CI Contributors
*/
package cli

import (
	initpkg "github.com/open-ug/conveyor/internal/init"
	"github.com/spf13/cobra"
)

var initOptions = &initpkg.Options{}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Conveyor CI system with configuration and certificates",
	Long: `Initialize Conveyor CI system for the first time by generating:

- conveyor.yml configuration file
- TLS/SSL certificate files (CA, server certificate, private key)
- Required directories with proper permissions

Examples:
  conveyor init                                          # Initialize with defaults
  conveyor init --config /path/to/conveyor.yml          # Use custom config file
  conveyor init --force                                  # Overwrite existing files
  conveyor init --auth-enabled=false                    # Disable authentication
  conveyor init --api-port 9090 --nats-port 4223       # Custom ports
  conveyor init --ca /path/ca.pem --private-key /path/key.pem --crt /path/cert.pem
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initpkg.Run(initOptions)
	},
}

func init() {
	// Configuration flags
	initCmd.Flags().StringVar(&initOptions.ConfigFile, "config", "", "Custom config file path to copy")
	initCmd.Flags().BoolVar(&initOptions.Force, "force", false, "Overwrite existing files")

	// API configuration flags
	initCmd.Flags().BoolVar(&initOptions.AuthEnabled, "auth-enabled", true, "Enable authentication")
	initCmd.Flags().IntVar(&initOptions.APIPort, "api-port", 8080, "API server port")
	initCmd.Flags().IntVar(&initOptions.NatsPort, "nats-port", 4222, "NATS server port")

	// Certificate flags
	initCmd.Flags().StringVar(&initOptions.CAFile, "ca", "", "Path to existing CA certificate")
	initCmd.Flags().StringVar(&initOptions.PrivateKeyFile, "private-key", "", "Path to existing private key")
	initCmd.Flags().StringVar(&initOptions.CertFile, "crt", "", "Path to existing server certificate")
}