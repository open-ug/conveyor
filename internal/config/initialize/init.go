package initialize

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/open-ug/conveyor/internal/utils"
)

// Options holds configuration for the init command
type Options struct {
	// Configuration options
	ConfigFile string
	Force      bool

	// API configuration
	AuthEnabled bool
	APIPort     int
	NatsPort    int

	// Certificate options
	CAFile         string
	PrivateKeyFile string
	CertFile       string
}

// Run executes the init command with the given options
func Run(opts *Options) (string, error) {
	fmt.Println("Initializing Conveyor CI...")

	// Determine base directories
	configDir, certDir, err := getSystemDirectories()
	if err != nil {
		return "", fmt.Errorf("failed to determine system directories: %w", err)
	}

	// Create directories with proper permissions
	if err := createDirectories(configDir, certDir); err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	// Handle certificate generation or copying
	if err := handleCertificates(opts, certDir); err != nil {
		return "", fmt.Errorf("failed to handle certificates: %w", err)
	}

	// Generate configuration file
	configPath := filepath.Join(configDir, "conveyor.yml")
	if err := generateConfig(opts, configDir, configPath, certDir); err != nil {
		return "", fmt.Errorf("failed to generate configuration: %w", err)
	}

	printSuccessMessage(configDir, certDir)
	return configPath, nil
}

// getSystemDirectories determines the appropriate config and cert directories
// based on whether we're running as root or regular user
func getSystemDirectories() (configDir, certDir string, err error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", "", err
	}

	if currentUser.Uid == "0" {
		// Running as root - use system directories
		configDir = "/var/lib/conveyor"
		certDir = "/var/lib/conveyor/certs"
	} else {
		// Running as regular user - use XDG directories
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			xdgConfig = filepath.Join(currentUser.HomeDir, ".local", "share")
		}
		configDir = filepath.Join(xdgConfig, "conveyor")
		certDir = filepath.Join(configDir, "certs")
	}

	if utils.IsTestMode() {
		// In test mode, use temporary directories
		tempDir := os.TempDir()
		configDir = filepath.Join(tempDir, "conveyor_test_config")
		certDir = filepath.Join(configDir, "certs")
	}

	return configDir, certDir, nil
}

func printSuccessMessage(configDir, certDir string) {
	fmt.Println("\n✔ Conveyor CI initialization completed successfully!")
	fmt.Println("\nGenerated files:")
	fmt.Printf("  • Configuration: %s/conveyor.yml\n", configDir)
	fmt.Printf("  • CA Certificate: %s/ca.pem\n", certDir)
	fmt.Printf("  • Server Certificate: %s/server.crt\n", certDir)
	fmt.Printf("  • Server Private Key: %s/server.key\n", certDir)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review and customize the configuration in conveyor.yml")
	fmt.Println("  2. Start the Conveyor API server: conveyor up")
	fmt.Println("  3. Visit the documentation: https://conveyor.open.ug")
}
