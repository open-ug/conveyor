package auth

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEndToEndAuthentication tests the complete authentication flow
func TestEndToEndAuthentication(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	// Initialize authentication components
	certManager := NewCertificateManager(certDir)
	jwtManager := NewJWTManager(certManager)
	tlsConfig := NewTLSConfig(certManager)

	// Test 1: Certificate loading and validation
	t.Run("Certificate Management", func(t *testing.T) {
		// Load server certificate
		serverCert, err := certManager.LoadServerCertificate()
		assert.NoError(t, err)
		assert.NotEmpty(t, serverCert.Certificate)

		// Load CA certificate
		caCert, err := certManager.LoadCACertificate()
		assert.NoError(t, err)
		assert.Equal(t, "Conveyor-CA", caCert.Subject.CommonName)

		// Load client certificate
		clientCert, err := certManager.LoadClientCertificate()
		assert.NoError(t, err)
		assert.NotEmpty(t, clientCert.Certificate)

		// Create certificate pool
		certPool, err := certManager.CreateCertPool()
		assert.NoError(t, err)
		assert.NotNil(t, certPool)
	})

	// Test 2: TLS configuration
	t.Run("TLS Configuration", func(t *testing.T) {
		// Create server TLS config
		serverTLSConfig, err := tlsConfig.CreateServerTLSConfig()
		assert.NoError(t, err)
		assert.NotNil(t, serverTLSConfig)
		assert.NotEmpty(t, serverTLSConfig.Certificates)

		// Create client TLS config
		clientTLSConfig, err := tlsConfig.CreateClientTLSConfig()
		assert.NoError(t, err)
		assert.NotNil(t, clientTLSConfig)
		assert.NotEmpty(t, clientTLSConfig.Certificates)
	})

	// Test 3: JWT token generation and validation
	t.Run("JWT Token Management", func(t *testing.T) {
		clientID := "integration-test-client"

		// Generate JWT token
		token, err := jwtManager.GenerateJWT(clientID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate JWT token
		claims, err := jwtManager.ValidateJWT(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, clientID, claims.ClientID)
		assert.Equal(t, "conveyor-client", claims.Issuer)
	})

	// Test 4: Complete API authentication flow
	t.Run("API Authentication Flow", func(t *testing.T) {
		// Create Fiber app with authentication middleware
		app := fiber.New()
		
		authMiddleware := NewAuthMiddleware(jwtManager, certManager)
		app.Use(authMiddleware.Handler())

		// Add protected route
		app.Get("/protected", func(c *fiber.Ctx) error {
			clientID := GetClientID(c)
			return c.JSON(fiber.Map{
				"message":   "access granted",
				"client_id": clientID,
			})
		})

		// Generate JWT token for test
		clientID := "api-test-client"
		token, err := jwtManager.GenerateJWT(clientID)
		require.NoError(t, err)

		// Test with valid JWT token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Test without JWT token (should fail)
		reqNoAuth := httptest.NewRequest("GET", "/protected", nil)
		respNoAuth, err := app.Test(reqNoAuth)
		assert.NoError(t, err)
		assert.Equal(t, 401, respNoAuth.StatusCode)
	})

	// Test 5: Public endpoints bypass authentication
	t.Run("Public Endpoints", func(t *testing.T) {
		app := fiber.New()
		
		authMiddleware := NewAuthMiddleware(jwtManager, certManager)
		app.Use(authMiddleware.Handler())

		// Add health check (should be public)
		app.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "healthy"})
		})

		// Test health endpoint without authentication
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

// TestTLSHandshakeSimulation simulates a TLS handshake scenario
func TestTLSHandshakeSimulation(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	certManager := NewCertificateManager(certDir)
	tlsConfig := NewTLSConfig(certManager)

	// Test server TLS configuration
	serverTLSConfig, err := tlsConfig.CreateServerTLSConfig()
	require.NoError(t, err)

	// Test client TLS configuration
	clientTLSConfig, err := tlsConfig.CreateClientTLSConfig()
	require.NoError(t, err)

	// Verify that both configurations are properly set up
	assert.NotEmpty(t, serverTLSConfig.Certificates)
	assert.NotNil(t, serverTLSConfig.ClientCAs)
	assert.Equal(t, tls.RequireAndVerifyClientCert, serverTLSConfig.ClientAuth)

	assert.NotEmpty(t, clientTLSConfig.Certificates)
	assert.NotNil(t, clientTLSConfig.RootCAs)
	assert.Equal(t, "Conveyor-API-Server", clientTLSConfig.ServerName)
}

// TestClientServerInteraction tests the interaction between client and server components
func TestClientServerInteraction(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	// Setup server-side authentication
	serverCertManager := NewCertificateManager(certDir)
	serverJWTManager := NewJWTManager(serverCertManager)
	
	// Setup client-side authentication (using same certificates for testing)
	clientCertManager := NewCertificateManager(certDir)
	clientJWTManager := NewJWTManager(clientCertManager)

	// Test scenario: Client generates JWT, server validates it
	t.Run("Client-Server JWT Exchange", func(t *testing.T) {
		clientID := "test-client"

		// Client generates JWT token
		clientToken, err := clientJWTManager.GenerateJWT(clientID)
		require.NoError(t, err)

		// Server validates the JWT token
		claims, err := serverJWTManager.ValidateJWT(clientToken)
		assert.NoError(t, err)
		assert.Equal(t, clientID, claims.ClientID)
	})

	// Test cross-validation to ensure security
	t.Run("Cross-Validation Security", func(t *testing.T) {
		// Generate token with one key
		token1, err := clientJWTManager.GenerateJWT("client1")
		require.NoError(t, err)

		// It should validate with the same certificate authority
		claims, err := serverJWTManager.ValidateJWT(token1)
		assert.NoError(t, err)
		assert.Equal(t, "client1", claims.ClientID)
	})
}

// TestAuthenticationConfiguration tests various authentication configurations
func TestAuthenticationConfiguration(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	certManager := NewCertificateManager(certDir)
	jwtManager := NewJWTManager(certManager)

	// Test with different middleware configurations
	t.Run("Required Authentication", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequireAuth(jwtManager, certManager))

		app.Get("/secure", func(c *fiber.Ctx) error {
			return c.SendString("secure endpoint")
		})

		// Should fail without token
		req := httptest.NewRequest("GET", "/secure", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("Optional Authentication", func(t *testing.T) {
		app := fiber.New()
		app.Use(OptionalAuth(jwtManager, certManager))

		app.Get("/flexible", func(c *fiber.Ctx) error {
			isAuth := IsAuthenticated(c)
			return c.JSON(fiber.Map{"authenticated": isAuth})
		})

		// Should succeed without token
		req := httptest.NewRequest("GET", "/flexible", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}