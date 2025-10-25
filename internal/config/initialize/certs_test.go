package initialize

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCertificates(t *testing.T) {
	tempDir := t.TempDir()

	err := generateCertificates(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to generate certificates: %v", err)
	}

	// Check if all certificate files were created
	expectedFiles := []string{
		"ca.key",
		"ca.pem",
		"server.key",
		"server.crt",
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(tempDir, filename)
		if !fileExists(path) {
			t.Errorf("Expected certificate file %s was not created", filename)
		}
	}

	// Verify CA certificate
	caCertPath := filepath.Join(tempDir, "ca.pem")
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		t.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		t.Fatal("Failed to decode CA certificate PEM")
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse CA certificate: %v", err)
	}

	// Verify CA certificate properties
	if caCert.Subject.CommonName != "Conveyor CA" {
		t.Errorf("Expected CA CommonName 'Conveyor CA', got '%s'", caCert.Subject.CommonName)
	}

	if !caCert.IsCA {
		t.Error("CA certificate should have IsCA=true")
	}

	// Verify server certificate
	serverCertPath := filepath.Join(tempDir, "server.crt")
	serverCertPEM, err := os.ReadFile(serverCertPath)
	if err != nil {
		t.Fatalf("Failed to read server certificate: %v", err)
	}

	serverCertBlock, _ := pem.Decode(serverCertPEM)
	if serverCertBlock == nil {
		t.Fatal("Failed to decode server certificate PEM")
	}

	serverCert, err := x509.ParseCertificate(serverCertBlock.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse server certificate: %v", err)
	}

	// Verify server certificate properties
	if serverCert.Subject.CommonName != "localhost" {
		t.Errorf("Expected server CommonName 'localhost', got '%s'", serverCert.Subject.CommonName)
	}

	if serverCert.IsCA {
		t.Error("Server certificate should have IsCA=false")
	}

	// Verify server certificate was signed by CA
	err = serverCert.CheckSignatureFrom(caCert)
	if err != nil {
		t.Errorf("Server certificate was not properly signed by CA: %v", err)
	}

	// Check file permissions
	testFilePermissions(t, filepath.Join(tempDir, "ca.key"), 0600)
	testFilePermissions(t, filepath.Join(tempDir, "server.key"), 0600)
	testFilePermissions(t, filepath.Join(tempDir, "ca.pem"), 0644)
	testFilePermissions(t, filepath.Join(tempDir, "server.crt"), 0644)
}

func TestGenerateCertificatesExistingFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create a dummy certificate file
	dummyCertPath := filepath.Join(tempDir, "server.crt")
	err := os.WriteFile(dummyCertPath, []byte("dummy cert"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy certificate: %v", err)
	}

	// Try to generate certificates without force flag - should fail
	err = generateCertificates(tempDir, false)
	if err == nil {
		t.Error("Expected error when certificates already exist without --force flag")
	}

	// Try with force flag - should succeed
	err = generateCertificates(tempDir, true)
	if err != nil {
		t.Errorf("Failed to generate certificates with force flag: %v", err)
	}
}

func TestCopyCertificates(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")

	// Create source and destination directories
	err := os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	// Create dummy certificate files
	caContent := "dummy CA certificate"
	keyContent := "dummy private key"
	certContent := "dummy server certificate"

	caFile := filepath.Join(srcDir, "ca.pem")
	keyFile := filepath.Join(srcDir, "server.key")
	certFile := filepath.Join(srcDir, "server.crt")

	err = os.WriteFile(caFile, []byte(caContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy CA file: %v", err)
	}

	err = os.WriteFile(keyFile, []byte(keyContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create dummy key file: %v", err)
	}

	err = os.WriteFile(certFile, []byte(certContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy cert file: %v", err)
	}

	// Set up options with source file paths
	opts := &Options{
		CAFile:         caFile,
		PrivateKeyFile: keyFile,
		CertFile:       certFile,
	}

	// Copy certificates
	err = copyCertificates(opts, destDir)
	if err != nil {
		t.Fatalf("Failed to copy certificates: %v", err)
	}

	// Verify files were copied
	destFiles := map[string]string{
		filepath.Join(destDir, "ca.pem"):     caContent,
		filepath.Join(destDir, "server.key"): keyContent,
		filepath.Join(destDir, "server.crt"): certContent,
	}

	for destPath, expectedContent := range destFiles {
		if !fileExists(destPath) {
			t.Errorf("Destination file %s was not created", destPath)
			continue
		}

		actualContent, err := os.ReadFile(destPath)
		if err != nil {
			t.Errorf("Failed to read destination file %s: %v", destPath, err)
			continue
		}

		if string(actualContent) != expectedContent {
			t.Errorf("Content mismatch for %s. Expected: %s, Got: %s",
				destPath, expectedContent, string(actualContent))
		}
	}

	// Check permissions (private keys should be 0600)
	testFilePermissions(t, filepath.Join(destDir, "server.key"), 0600)
}

func testFilePermissions(t *testing.T, path string, expectedPerm os.FileMode) {
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("Failed to stat file %s: %v", path, err)
		return
	}

	actualPerm := info.Mode().Perm()
	if actualPerm != expectedPerm {
		t.Errorf("File %s has wrong permissions. Expected: %o, Got: %o",
			path, expectedPerm, actualPerm)
	}
}
