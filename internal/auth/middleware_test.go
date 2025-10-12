package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_SkipAuthRoutes(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)
	middleware := NewAuthMiddleware(jm, cm)

	app := fiber.New()
	app.Use(middleware.Handler())
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("healthy")
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)
	middleware := NewAuthMiddleware(jm, cm)

	app := fiber.New()
	app.Use(middleware.Handler())
	
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("protected content")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthMiddleware_InvalidAuthorizationHeader(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)
	middleware := NewAuthMiddleware(jm, cm)

	app := fiber.New()
	app.Use(middleware.Handler())
	
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("protected content")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Invalid token-format")
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthMiddleware_EmptyBearerToken(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)
	middleware := NewAuthMiddleware(jm, cm)

	app := fiber.New()
	app.Use(middleware.Handler())
	
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("protected content")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthMiddleware_ValidJWTToken(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)
	middleware := NewAuthMiddleware(jm, cm)

	// Generate a valid JWT token
	token, err := jm.GenerateJWT("test-client")
	require.NoError(t, err)

	app := fiber.New()
	app.Use(middleware.Handler())
	
	app.Get("/protected", func(c *fiber.Ctx) error {
		// Check if authentication context is set
		clientID := GetClientID(c)
		assert.Equal(t, "test-client", clientID)
		
		claims := GetJWTClaims(c)
		assert.NotNil(t, claims)
		assert.Equal(t, "test-client", claims.ClientID)
		
		authenticated := IsAuthenticated(c)
		assert.True(t, authenticated)
		
		return c.SendString("protected content")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthMiddleware_SetSkipAuthRoutes(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)
	middleware := NewAuthMiddleware(jm, cm)

	// Set custom skip routes
	customSkipRoutes := []string{"/public", "/open"}
	middleware.SetSkipAuthRoutes(customSkipRoutes)

	app := fiber.New()
	app.Use(middleware.Handler())
	
	app.Get("/public", func(c *fiber.Ctx) error {
		return c.SendString("public content")
	})
	
	app.Get("/open", func(c *fiber.Ctx) error {
		return c.SendString("open content")
	})

	// Test that custom skip routes work
	req1 := httptest.NewRequest("GET", "/public", nil)
	resp1, err := app.Test(req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)

	req2 := httptest.NewRequest("GET", "/open", nil)
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp2.StatusCode)
}

func TestRequireAuth(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	app := fiber.New()
	app.Use(RequireAuth(jm, cm))
	
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("protected content")
	})

	// Test without authentication
	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestOptionalAuth_WithoutAuth(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	app := fiber.New()
	app.Use(OptionalAuth(jm, cm))
	
	app.Get("/optional", func(c *fiber.Ctx) error {
		authenticated := IsAuthenticated(c)
		assert.False(t, authenticated)
		
		clientID := GetClientID(c)
		assert.Empty(t, clientID)
		
		return c.SendString("optional content")
	})

	req := httptest.NewRequest("GET", "/optional", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestOptionalAuth_WithAuth(t *testing.T) {
	certDir, cleanup := setupTestCertificates(t)
	defer cleanup()

	cm := NewCertificateManager(certDir)
	jm := NewJWTManager(cm)

	// Generate a valid JWT token
	token, err := jm.GenerateJWT("test-client")
	require.NoError(t, err)

	app := fiber.New()
	app.Use(OptionalAuth(jm, cm))
	
	app.Get("/optional", func(c *fiber.Ctx) error {
		authenticated := IsAuthenticated(c)
		assert.True(t, authenticated)
		
		clientID := GetClientID(c)
		assert.Equal(t, "test-client", clientID)
		
		return c.SendString("optional content")
	})

	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthContextHelpers(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		// Test with no auth context
		clientID := GetClientID(c)
		assert.Empty(t, clientID)
		
		claims := GetJWTClaims(c)
		assert.Nil(t, claims)
		
		cert := GetClientCertificate(c)
		assert.Nil(t, cert)
		
		authenticated := IsAuthenticated(c)
		assert.False(t, authenticated)
		
		// Set some context manually for testing
		c.Locals("client_id", "manual-client")
		c.Locals("authenticated", true)
		
		clientID2 := GetClientID(c)
		assert.Equal(t, "manual-client", clientID2)
		
		authenticated2 := IsAuthenticated(c)
		assert.True(t, authenticated2)
		
		return c.SendString("test")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}