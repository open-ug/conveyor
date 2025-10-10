package init

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// handleCertificates manages certificate generation or copying based on options
func handleCertificates(opts *Options, certDir string) error {
	// Check if user provided existing certificates
	if opts.CAFile != "" || opts.PrivateKeyFile != "" || opts.CertFile != "" {
		return copyCertificates(opts, certDir)
	}

	// Generate new certificates
	return generateCertificates(certDir, opts.Force)
}

// copyCertificates copies user-provided certificates to the cert directory
func copyCertificates(opts *Options, certDir string) error {
	certFiles := map[string]string{
		opts.CAFile:         filepath.Join(certDir, "ca.pem"),
		opts.PrivateKeyFile: filepath.Join(certDir, "server.key"),
		opts.CertFile:       filepath.Join(certDir, "server.crt"),
	}

	for srcPath, destPath := range certFiles {
		if srcPath == "" {
			continue
		}

		if err := copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %w", srcPath, destPath, err)
		}
	}

	fmt.Println("📋 Copied existing certificates to cert directory")
	return nil
}

// generateCertificates creates a new CA and server certificate pair
func generateCertificates(certDir string, force bool) error {
	caKeyPath := filepath.Join(certDir, "ca.key")
	caCertPath := filepath.Join(certDir, "ca.pem")
	serverKeyPath := filepath.Join(certDir, "server.key")
	serverCertPath := filepath.Join(certDir, "server.crt")

	// Check if certificates already exist
	if !force {
		if fileExists(caCertPath) || fileExists(serverCertPath) {
			return fmt.Errorf("certificates already exist in %s (use --force to overwrite)", certDir)
		}
	}

	fmt.Println("🔐 Generating CA certificate...")

	// Generate CA private key
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate CA private key: %w", err)
	}

	// Create CA certificate template
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Conveyor CA",
			Organization: []string{"Conveyor CI"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create CA certificate
	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	// Save CA private key
	if err := savePrivateKey(caKey, caKeyPath); err != nil {
		return fmt.Errorf("failed to save CA private key: %w", err)
	}

	// Save CA certificate
	if err := saveCertificate(caCertDER, caCertPath); err != nil {
		return fmt.Errorf("failed to save CA certificate: %w", err)
	}

	fmt.Println("🔐 Generating server certificate...")

	// Generate server private key
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate server private key: %w", err)
	}

	// Create server certificate template
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames: []string{
			"localhost",
		},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
	}

	// Parse CA certificate for signing
	caCert, err := x509.ParseCertificate(caCertDER)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Create server certificate
	serverCertDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("failed to create server certificate: %w", err)
	}

	// Save server private key
	if err := savePrivateKey(serverKey, serverKeyPath); err != nil {
		return fmt.Errorf("failed to save server private key: %w", err)
	}

	// Save server certificate
	if err := saveCertificate(serverCertDER, serverCertPath); err != nil {
		return fmt.Errorf("failed to save server certificate: %w", err)
	}

	fmt.Println("✅ Successfully generated TLS certificates")
	return nil
}

// savePrivateKey saves an RSA private key in PKCS#8 PEM format
func savePrivateKey(key *rsa.PrivateKey, path string) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	keyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, keyPEM)
}

// saveCertificate saves an X.509 certificate in PEM format
func saveCertificate(certDER []byte, path string) error {
	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, certPEM)
}

// copyFile copies a file from src to dst with appropriate permissions
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Determine appropriate permissions based on file extension
	var perm os.FileMode = 0644
	if filepath.Ext(src) == ".key" || filepath.Base(src) == "server.key" {
		perm = 0600
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}