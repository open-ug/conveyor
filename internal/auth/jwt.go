package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	certificateManager *CertificateManager
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(certManager *CertificateManager) *JWTManager {
	return &JWTManager{
		certificateManager: certManager,
	}
}

// Claims represents the JWT claims structure
type Claims struct {
	ClientID string `json:"client_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token using the client's private key
func (jm *JWTManager) GenerateJWT(clientID string) (string, error) {
	// Load client certificate to get the private key
	clientCert, err := jm.certificateManager.LoadClientCertificate()
	if err != nil {
		return "", fmt.Errorf("failed to load client certificate: %w", err)
	}

	// Extract private key from certificate
	privateKey, ok := clientCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("client certificate does not contain an RSA private key")
	}

	// Create claims
	now := time.Now()
	claims := Claims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "conveyor-client",
			Subject:   clientID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// Sign token with private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token using the CA certificate
func (jm *JWTManager) ValidateJWT(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key ID from the token header (if present)
		// For now, we'll validate against any client certificate signed by our CA
		return jm.getValidationKey(token)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid JWT token")
}

// getValidationKey gets the public key for JWT validation
func (jm *JWTManager) getValidationKey(token *jwt.Token) (interface{}, error) {
	// Load client certificate to get the public key
	// In a more sophisticated implementation, we might store multiple
	// client certificates and select based on key ID in the token header
	clientCert, err := jm.certificateManager.LoadClientCertificate()
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// Parse the certificate to get the public key
	if len(clientCert.Certificate) == 0 {
		return nil, fmt.Errorf("no certificate data found")
	}

	cert, err := x509.ParseCertificate(clientCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Validate the certificate against our CA
	err = jm.certificateManager.ValidateClientCertificate(cert)
	if err != nil {
		return nil, fmt.Errorf("client certificate validation failed: %w", err)
	}

	// Return the public key for JWT validation
	publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate does not contain an RSA public key")
	}

	return publicKey, nil
}

// GenerateJWTWithPrivateKey generates a JWT using a provided private key
// This is useful for clients that have their private key separately
func (jm *JWTManager) GenerateJWTWithPrivateKey(clientID string, privateKey *rsa.PrivateKey) (string, error) {
	// Create claims
	now := time.Now()
	claims := Claims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "conveyor-client",
			Subject:   clientID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// Sign token with private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

// ExtractClientIDFromCertificate extracts client ID from certificate subject
func ExtractClientIDFromCertificate(cert *x509.Certificate) string {
	if cert.Subject.CommonName != "" {
		return cert.Subject.CommonName
	}
	return "unknown-client"
}

// HashCertificate creates a SHA256 hash of the certificate for identification
func HashCertificate(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	return base64.URLEncoding.EncodeToString(hash[:])
}