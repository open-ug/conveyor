package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// CertificateManager handles loading and managing TLS certificates
type CertificateManager struct {
	CertDir string
}

// NewCertificateManager creates a new certificate manager
func NewCertificateManager(certDir string) *CertificateManager {
	return &CertificateManager{
		CertDir: certDir,
	}
}

// LoadServerCertificate loads the server certificate and private key
func (cm *CertificateManager) LoadServerCertificate() (tls.Certificate, error) {
	certFile := filepath.Join(cm.CertDir, "server.pem")
	keyFile := filepath.Join(cm.CertDir, "server.key")
	
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load server certificate: %w", err)
	}
	
	return cert, nil
}

// LoadCACertificate loads the Certificate Authority certificate
func (cm *CertificateManager) LoadCACertificate() (*x509.Certificate, error) {
	caFile := filepath.Join(cm.CertDir, "ca.pem")
	
	caCertPEM, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	
	caCert, err := x509.ParseCertificate(caCertPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}
	
	return caCert, nil
}

// LoadClientCertificate loads a client certificate and private key
func (cm *CertificateManager) LoadClientCertificate() (tls.Certificate, error) {
	certFile := filepath.Join(cm.CertDir, "client.pem")
	keyFile := filepath.Join(cm.CertDir, "client.key")
	
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load client certificate: %w", err)
	}
	
	return cert, nil
}

// CreateCertPool creates a certificate pool with the CA certificate
func (cm *CertificateManager) CreateCertPool() (*x509.CertPool, error) {
	caFile := filepath.Join(cm.CertDir, "ca.pem")
	
	caCertPEM, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCertPEM) {
		return nil, fmt.Errorf("failed to append CA certificate to pool")
	}
	
	return caCertPool, nil
}

// ValidateClientCertificate validates a client certificate against the CA
func (cm *CertificateManager) ValidateClientCertificate(clientCert *x509.Certificate) error {
	caCertPool, err := cm.CreateCertPool()
	if err != nil {
		return fmt.Errorf("failed to create CA cert pool: %w", err)
	}
	
	opts := x509.VerifyOptions{
		Roots: caCertPool,
	}
	
	_, err = clientCert.Verify(opts)
	if err != nil {
		return fmt.Errorf("client certificate validation failed: %w", err)
	}
	
	return nil
}