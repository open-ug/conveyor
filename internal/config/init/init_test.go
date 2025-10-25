package init

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateDefaultConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "conveyor.yml")
	certDir := filepath.Join(tempDir, "certs")

	// Create options
	opts := &Options{
		APIPort:     9090,
		NatsPort:    4223,
		AuthEnabled: false,
	}

	// Generate config
	err := generateDefaultConfig(opts, tempDir, configPath, certDir)
	if err != nil {
		t.Fatalf("Failed to generate config: %v", err)
	}

	// Check if file was created
	if !fileExists(configPath) {
		t.Fatal("Config file was not created")
	}

	// Read and verify contents
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	configStr := string(content)

	// Check for expected values
	expectedContents := []string{
		"port: 9090",
		"port: 4223",
		"auth_enabled: false",
		"ca: " + filepath.Join(certDir, "ca.pem"),
		"key: " + filepath.Join(certDir, "server.key"),
		"cert: " + filepath.Join(certDir, "server.crt"),
	}

	for _, expected := range expectedContents {
		if !contains(configStr, expected) {
			t.Errorf("Config file missing expected content: %s", expected)
		}
	}
}

func TestCreateDirectories(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	certDir := filepath.Join(tempDir, "certs")

	err := createDirectories(configDir, certDir)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Check if directories were created
	if !dirExists(configDir) {
		t.Fatal("Config directory was not created")
	}

	if !dirExists(certDir) {
		t.Fatal("Cert directory was not created")
	}

	// Check permissions on cert directory (should be more restrictive)
	certInfo, err := os.Stat(certDir)
	if err != nil {
		t.Fatalf("Failed to stat cert directory: %v", err)
	}

	if certInfo.Mode().Perm() != 0700 {
		t.Errorf("Expected cert directory permissions 0700, got %o", certInfo.Mode().Perm())
	}
}

func TestGetSystemDirectories(t *testing.T) {
	configDir, certDir, err := getSystemDirectories()
	if err != nil {
		t.Fatalf("Failed to get system directories: %v", err)
	}

	if configDir == "" {
		t.Error("Config directory should not be empty")
	}

	if certDir == "" {
		t.Error("Cert directory should not be empty")
	}

	// Cert directory should be a subdirectory of config directory
	if !contains(certDir, configDir) && !contains(certDir, "certs") {
		t.Error("Cert directory should be related to config directory or contain 'certs'")
	}
}

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()

	// Test non-existent file
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	if fileExists(nonExistentFile) {
		t.Error("fileExists should return false for non-existent file")
	}

	// Create a file and test
	existentFile := filepath.Join(tempDir, "existent.txt")
	err := os.WriteFile(existentFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !fileExists(existentFile) {
		t.Error("fileExists should return true for existing file")
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				len(s) > len(substr) && s[1:len(substr)+1] == substr ||
				findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
