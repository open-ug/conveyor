package auth

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTLSConfig_CreateServerTLSConfig(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	tc := NewTLSConfig(cm)

	tlsConfig, err := tc.CreateServerTLSConfig()
	assert.NoError(t, err)
	assert.NotNil(t, tlsConfig)

	// Verify TLS configuration
	assert.NotEmpty(t, tlsConfig.Certificates)
	assert.NotNil(t, tlsConfig.ClientCAs)
	assert.Equal(t, 2, int(tlsConfig.ClientAuth)) // RequireAndVerifyClientCert = 2
	assert.GreaterOrEqual(t, int(tlsConfig.MinVersion), int(0x0303)) // TLS 1.2 = 0x0303
	assert.NotEmpty(t, tlsConfig.CipherSuites)
	assert.True(t, tlsConfig.PreferServerCipherSuites)
	assert.NotNil(t, tlsConfig.VerifyPeerCertificate)
}

func TestTLSConfig_CreateClientTLSConfig(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	tc := NewTLSConfig(cm)

	tlsConfig, err := tc.CreateClientTLSConfig()
	assert.NoError(t, err)
	assert.NotNil(t, tlsConfig)

	// Verify TLS configuration
	assert.NotEmpty(t, tlsConfig.Certificates)
	assert.NotNil(t, tlsConfig.RootCAs)
	assert.GreaterOrEqual(t, int(tlsConfig.MinVersion), int(0x0303)) // TLS 1.2 = 0x0303
	assert.NotEmpty(t, tlsConfig.CipherSuites)
	assert.Equal(t, "Conveyor-API-Server", tlsConfig.ServerName)
}

func TestTLSConfig_MissingCertificates(t *testing.T) {
	cm := NewCertificateManager("/non/existent/path")
	tc := NewTLSConfig(cm)

	_, err := tc.CreateServerTLSConfig()
	assert.Error(t, err)

	_, err = tc.CreateClientTLSConfig()
	assert.Error(t, err)
}

func TestGetCertificateInfo(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	
	// Load client certificate for testing
	clientCert, err := cm.LoadClientCertificate()
	require.NoError(t, err)
	
	// Parse the certificate
	cert, err := x509.ParseCertificate(clientCert.Certificate[0])
	require.NoError(t, err)

	info := GetCertificateInfo(cert)
	
	assert.NotNil(t, info)
	assert.Contains(t, info, "subject")
	assert.Contains(t, info, "issuer")
	assert.Contains(t, info, "serial")
	assert.Contains(t, info, "not_before")
	assert.Contains(t, info, "not_after")
	assert.Contains(t, info, "dns_names")
	assert.Contains(t, info, "ip_addresses")
	assert.Contains(t, info, "key_usage")
	assert.Contains(t, info, "ext_key_usage")
	
	// Verify some expected values
	subject, ok := info["subject"].(string)
	assert.True(t, ok)
	assert.Contains(t, subject, "Conveyor-Client")
}

func TestTLSConfig_VerifyPeerCertificate(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	tc := NewTLSConfig(cm)

	// Load client certificate to get raw bytes
	clientCert, err := cm.LoadClientCertificate()
	require.NoError(t, err)

	rawCerts := [][]byte{clientCert.Certificate[0]}
	
	// Test with valid certificate (this is internal function, but we can test the concept)
	err = tc.verifyClientCertificate(rawCerts, nil)
	assert.NoError(t, err)

	// Test with empty certificate list
	err = tc.verifyClientCertificate([][]byte{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no client certificate provided")

	// Test with invalid certificate data
	invalidRawCerts := [][]byte{[]byte("invalid-certificate-data")}
	err = tc.verifyClientCertificate(invalidRawCerts, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse client certificate")
}