package nginxdriver

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
)

func CreateNginxConfig(app craneTypes.Application) error {
	fmt.Println("NGINX_D: Creating App")
	ports := app.Spec.Ports

	for i := 0; i < len(ports); i++ {
		port := ports[i]
		fmt.Println("NGINX_D: Generating App")
		configContent := generateNginxConfigFromPort(port)
		filename := app.Name + "-" + string(rune(port.External)) + ".conf"
		fmt.Println("NGINX_D: Writing")
		err := writeNginxConfig(filename, configContent)
		fmt.Println("NGINX_D: Wrote file")
		if err != nil {
			return fmt.Errorf("failed to write nginx config: %w", err)
		}
	}

	return nil
}

func DeleteNginxConfig(app craneTypes.Application) error {
	ports := app.Spec.Ports

	for i := 0; i < len(ports); i++ {
		port := ports[i]
		filename := app.Name + "-" + strconv.Itoa(port.External) + ".conf"
		err := removeNginxConfig(filename)
		if err != nil {
			return fmt.Errorf("failed to remove nginx config: %w", err)
		}
	}

	return nil
}

func UpdateNginxConfig(app craneTypes.Application) error {
	// Delete the old config
	err := DeleteNginxConfig(app)
	if err != nil {
		return fmt.Errorf("failed to delete nginx config: %w", err)
	}

	// Create the new config
	err = CreateNginxConfig(app)
	if err != nil {
		return fmt.Errorf("failed to create nginx config: %w", err)
	}

	return nil
}

func generateNginxConfigFromPort(port craneTypes.ApplicationPortMap) string {
	fmt.Println("NGINX_D: Generating Port Spec")
	// Create the server block
	serverBlock := "server {\n"
	serverBlock += "    listen " + port.Domain + ";\n"
	serverBlock += "    server_name " + port.Domain + ";\n"

	// Create the location block
	locationBlock := "    location / {\n"
	locationBlock += "        proxy_pass http://localhost:" + strconv.Itoa(port.External) + ";\n"
	locationBlock += "        proxy_set_header Upgrade $http_upgrade;\n"
	locationBlock += "        proxy_set_header Connection 'upgrade';\n"
	locationBlock += "        proxy_set_header Host $host;\n"
	locationBlock += "        proxy_cache_bypass $http_upgrade;\n"
	locationBlock += "    }\n"

	// Close the location block
	locationBlock += "}"

	// Close the server block
	serverBlock += locationBlock

	return serverBlock

}

// A function that writes the nginx config to /etc/crane/conf/nginx/sites-available/<app-name> and creates a symlink to /etc/crane/conf/nginx/sites-enabled/<app-name>
func writeNginxConfig(appName, configContent string) error {
	// Define directories
	sitesAvailableDir := "/etc/crane/conf/nginx/sites-available"
	sitesEnabledDir := "/etc/crane/conf/nginx/sites-enabled"
	// Ensure directories exist
	if err := os.MkdirAll(sitesAvailableDir, 0755); err != nil {
		return fmt.Errorf("failed to create sites-available directory: %w", err)
	}
	if err := os.MkdirAll(sitesEnabledDir, 0755); err != nil {
		return fmt.Errorf("failed to create sites-enabled directory: %w", err)
	}

	// Define file paths
	availablePath := filepath.Join(sitesAvailableDir, appName)
	enabledPath := filepath.Join(sitesEnabledDir, appName)

	// Write the config file to sites-available
	err := os.WriteFile(availablePath, []byte(configContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write nginx config file: %w", err)
	}

	// Create the symlink in sites-enabled
	err = os.Symlink(availablePath, enabledPath)
	if err != nil {
		// If the symlink creation fails, clean up the written config
		os.Remove(availablePath)
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// A function that removes the nginx config from /etc/crane/conf/nginx/sites-available/<app-name> and removes the symlink from /etc/crane/conf/nginx/sites-enabled/<app-name>
func removeNginxConfig(appName string) error {
	// Define directories
	sitesAvailableDir := "/etc/crane/conf/nginx/sites-available"
	sitesEnabledDir := "/etc/crane/conf/nginx/sites-enabled"

	// Define file paths
	availablePath := filepath.Join(sitesAvailableDir, appName)
	enabledPath := filepath.Join(sitesEnabledDir, appName)

	// Remove the config file from sites-available
	err := os.Remove(availablePath)
	if err != nil {
		return fmt.Errorf("failed to remove nginx config file: %w", err)
	}

	// Remove the symlink from sites-enabled
	err = os.Remove(enabledPath)
	if err != nil {
		return fmt.Errorf("failed to remove symlink: %w", err)
	}

	return nil
}
