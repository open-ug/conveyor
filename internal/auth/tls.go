package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// TLSConfig contains TLS configuration for the server
type TLSConfig struct {
	certificateManager *CertificateManager
}

// NewTLSConfig creates a new TLS configuration
func NewTLSConfig(certManager *CertificateManager) *TLSConfig {
	return &TLSConfig{
		certificateManager: certManager,
	}
}

// CreateServerTLSConfig creates a TLS configuration for the server
func (tc *TLSConfig) CreateServerTLSConfig() (*tls.Config, error) {
	// Load server certificate
	serverCert, err := tc.certificateManager.LoadServerCertificate()
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	// Create CA cert pool for client verification
	caCertPool, err := tc.certificateManager.CreateCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to create CA cert pool: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true,
		// Custom verification function for additional client certificate validation
		VerifyPeerCertificate: tc.verifyClientCertificate,
	}

	return tlsConfig, nil
}

// CreateClientTLSConfig creates a TLS configuration for clients
func (tc *TLSConfig) CreateClientTLSConfig() (*tls.Config, error) {
	// Load client certificate
	clientCert, err := tc.certificateManager.LoadClientCertificate()
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// Create CA cert pool for server verification
	caCertPool, err := tc.certificateManager.CreateCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to create CA cert pool: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		// Server name for certificate verification (should match server certificate)
		ServerName: "Conveyor-API-Server",
	}

	return tlsConfig, nil
}

// verifyClientCertificate provides additional validation for client certificates
func (tc *TLSConfig) verifyClientCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	// Basic validation is already done by TLS layer
	// This function can be used for additional custom validation if needed
	
	if len(rawCerts) == 0 {
		return fmt.Errorf("no client certificate provided")
	}

	// Parse the client certificate
	clientCert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return fmt.Errorf("failed to parse client certificate: %w", err)
	}

	// Additional validation using our certificate manager
	err = tc.certificateManager.ValidateClientCertificate(clientCert)
	if err != nil {
		return fmt.Errorf("client certificate validation failed: %w", err)
	}

	return nil
}

// GetCertificateInfo returns information about a certificate
func GetCertificateInfo(cert *x509.Certificate) map[string]interface{} {
	return map[string]interface{}{
		"subject":      cert.Subject.String(),
		"issuer":       cert.Issuer.String(),
		"serial":       cert.SerialNumber.String(),
		"not_before":   cert.NotBefore,
		"not_after":    cert.NotAfter,
		"dns_names":    cert.DNSNames,
		"ip_addresses": cert.IPAddresses,
		"key_usage":    cert.KeyUsage,
		"ext_key_usage": cert.ExtKeyUsage,
	}
}