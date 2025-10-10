package init

import (
	"fmt"
	"os"
)

// createDirectories creates the necessary directories with proper permissions
func createDirectories(configDir, certDir string) error {
	directories := []struct {
		path string
		perm os.FileMode
		desc string
	}{
		{configDir, 0755, "configuration directory"},
		{certDir, 0700, "certificate directory"},
	}

	for _, dir := range directories {
		if err := createDirectoryIfNotExists(dir.path, dir.perm, dir.desc); err != nil {
			return err
		}
	}

	return nil
}

// createDirectoryIfNotExists creates a directory if it doesn't exist
func createDirectoryIfNotExists(path string, perm os.FileMode, desc string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, perm); err != nil {
			return fmt.Errorf("failed to create %s at %s: %w", desc, path, err)
		}
		fmt.Printf("📁 Created %s: %s\n", desc, path)
	} else if err != nil {
		return fmt.Errorf("failed to check %s at %s: %w", desc, path, err)
	} else {
		fmt.Printf("📁 Using existing %s: %s\n", desc, path)
	}

	return nil
}