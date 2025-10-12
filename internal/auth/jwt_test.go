package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateJWT(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	token, err := jm.GenerateJWT("test-client")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Verify it's a valid JWT format (3 parts separated by dots)
	parts := jwt.NewParser().Parse(token, nil)
	assert.NotNil(t, parts)
}

func TestJWTManager_ValidateJWT(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	// Generate a token
	clientID := "test-client"
	token, err := jm.GenerateJWT(clientID)
	require.NoError(t, err)

	// Validate the token
	claims, err := jm.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, clientID, claims.ClientID)
	assert.Equal(t, "conveyor-client", claims.Issuer)
	assert.Equal(t, clientID, claims.Subject)
}

func TestJWTManager_GenerateJWTWithPrivateKey(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	// Generate a new private key for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	clientID := "test-client-with-key"
	token, err := jm.GenerateJWTWithPrivateKey(clientID, privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Parse the token to verify structure
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	
	if claims, ok := parsedToken.Claims.(*Claims); ok {
		assert.Equal(t, clientID, claims.ClientID)
		assert.Equal(t, "conveyor-client", claims.Issuer)
		assert.Equal(t, clientID, claims.Subject)
	}
}

func TestJWTManager_ValidateInvalidJWT(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	// Test with invalid token
	_, err := jm.ValidateJWT("invalid.jwt.token")
	assert.Error(t, err)

	// Test with empty token
	_, err = jm.ValidateJWT("")
	assert.Error(t, err)

	// Test with malformed token
	_, err = jm.ValidateJWT("not-a-jwt")
	assert.Error(t, err)
}

func TestJWTManager_ValidateExpiredJWT(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)

	// Load client certificate to get private key
	clientCert, err := cm.LoadClientCertificate()
	require.NoError(t, err)

	privateKey, ok := clientCert.PrivateKey.(*rsa.PrivateKey)
	require.True(t, ok)

	// Create an expired token
	now := time.Now()
	claims := Claims{
		ClientID: "test-client",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(now.Add(-time.Hour * 2)),
			NotBefore: jwt.NewNumericDate(now.Add(-time.Hour * 2)),
			Issuer:    "conveyor-client",
			Subject:   "test-client",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	jm := NewJWTManager(cm)
	_, err = jm.ValidateJWT(tokenString)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWTManager_MissingCertificates(t *testing.T) {
	// Test with non-existent certificate directory
	cm := NewCertificateManager("/non/existent/path")
	jm := NewJWTManager(cm)

	_, err := jm.GenerateJWT("test-client")
	assert.Error(t, err)

	_, err = jm.ValidateJWT("some.jwt.token")
	assert.Error(t, err)
}

func TestJWTClaims(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	clientID := "detailed-test-client"
	token, err := jm.GenerateJWT(clientID)
	require.NoError(t, err)

	claims, err := jm.ValidateJWT(token)
	require.NoError(t, err)

	// Verify all claim fields
	assert.Equal(t, clientID, claims.ClientID)
	assert.Equal(t, "conveyor-client", claims.Issuer)
	assert.Equal(t, clientID, claims.Subject)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
	
	// Verify token is not expired
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	
	// Verify token was issued in the past (within last minute)
	assert.True(t, claims.IssuedAt.Before(time.Now()))
	assert.True(t, claims.IssuedAt.After(time.Now().Add(-time.Minute)))
}