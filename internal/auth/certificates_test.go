package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestCertificates(t *testing.T) (string, func()) {
	// Create temporary directory for test certificates
	tempDir, err := os.MkdirTemp("", "conveyor-test-certs")
	require.NoError(t, err)

	// Generate CA private key
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create CA certificate template
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Conveyor-CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create self-signed CA certificate
	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	require.NoError(t, err)

	// Save CA certificate
	caCertFile := filepath.Join(tempDir, "ca.pem")
	caCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertDER,
	})
	require.NoError(t, os.WriteFile(caCertFile, caCertPEM, 0644))

	// Save CA private key
	caKeyFile := filepath.Join(tempDir, "ca.key")
	caKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
	})
	require.NoError(t, os.WriteFile(caKeyFile, caKeyPEM, 0600))

	// Generate server private key
	serverPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create server certificate template
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: "Conveyor-API-Server",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
	}

	// Create server certificate signed by CA
	serverCertDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caTemplate, &serverPrivateKey.PublicKey, caPrivateKey)
	require.NoError(t, err)

	// Save server certificate
	serverCertFile := filepath.Join(tempDir, "server.pem")
	serverCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertDER,
	})
	require.NoError(t, os.WriteFile(serverCertFile, serverCertPEM, 0644))

	// Save server private key
	serverKeyFile := filepath.Join(tempDir, "server.key")
	serverKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivateKey),
	})
	require.NoError(t, os.WriteFile(serverKeyFile, serverKeyPEM, 0600))

	// Generate client private key
	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create client certificate template
	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName: "Conveyor-Client",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	}

	// Create client certificate signed by CA
	clientCertDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caTemplate, &clientPrivateKey.PublicKey, caPrivateKey)
	require.NoError(t, err)

	// Save client certificate
	clientCertFile := filepath.Join(tempDir, "client.pem")
	clientCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCertDER,
	})
	require.NoError(t, os.WriteFile(clientCertFile, clientCertPEM, 0644))

	// Save client private key
	clientKeyFile := filepath.Join(tempDir, "client.key")
	clientKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(clientPrivateKey),
	})
	require.NoError(t, os.WriteFile(clientKeyFile, clientKeyPEM, 0600))

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestCertificateManager_LoadServerCertificate(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	cert, err := cm.LoadServerCertificate()

	assert.NoError(t, err)
	assert.NotNil(t, cert.Certificate)
	assert.NotNil(t, cert.PrivateKey)
}

func TestCertificateManager_LoadCACertificate(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	caCert, err := cm.LoadCACertificate()

	assert.NoError(t, err)
	assert.NotNil(t, caCert)
	assert.Equal(t, "Conveyor-CA", caCert.Subject.CommonName)
}

func TestCertificateManager_LoadClientCertificate(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	cert, err := cm.LoadClientCertificate()

	assert.NoError(t, err)
	assert.NotNil(t, cert.Certificate)
	assert.NotNil(t, cert.PrivateKey)
}

func TestCertificateManager_CreateCertPool(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	certPool, err := cm.CreateCertPool()

	assert.NoError(t, err)
	assert.NotNil(t, certPool)
}

func TestCertificateManager_ValidateClientCertificate(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	
	// Load the client certificate
	clientCert, err := cm.LoadClientCertificate()
	require.NoError(t, err)
	
	// Parse the certificate
	cert, err := x509.ParseCertificate(clientCert.Certificate[0])
	require.NoError(t, err)

	// Validate should succeed
	err = cm.ValidateClientCertificate(cert)
	assert.NoError(t, err)
}

func TestCertificateManager_MissingFiles(t *testing.T) {
	// Test with non-existent directory
	cm := NewCertificateManager("/non/existent/path")
	
	_, err := cm.LoadServerCertificate()
	assert.Error(t, err)
	
	_, err = cm.LoadCACertificate()
	assert.Error(t, err)
	
	_, err = cm.LoadClientCertificate()
	assert.Error(t, err)
	
	_, err = cm.CreateCertPool()
	assert.Error(t, err)
}

func TestExtractClientIDFromCertificate(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	clientCert, err := cm.LoadClientCertificate()
	require.NoError(t, err)
	
	cert, err := x509.ParseCertificate(clientCert.Certificate[0])
	require.NoError(t, err)

	clientID := ExtractClientIDFromCertificate(cert)
	assert.Equal(t, "Conveyor-Client", clientID)
}

func TestHashCertificate(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	clientCert, err := cm.LoadClientCertificate()
	require.NoError(t, err)
	
	cert, err := x509.ParseCertificate(clientCert.Certificate[0])
	require.NoError(t, err)

	hash := HashCertificate(cert)
	assert.NotEmpty(t, hash)
	assert.Greater(t, len(hash), 10) // Should be a base64 encoded hash
}