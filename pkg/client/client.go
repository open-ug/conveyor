package client

import (
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/open-ug/conveyor/internal/auth"
	config "github.com/open-ug/conveyor/internal/config"
	"github.com/spf13/viper"
)

type Client struct {
	HTTPClient         *resty.Client
	AuthManager        *ClientAuthManager
	jwtToken           string
}

// ClientAuthManager handles client-side authentication
type ClientAuthManager struct {
	CertificateManager *auth.CertificateManager
	JWTManager         *auth.JWTManager
	TLSConfig          *auth.TLSConfig
	ClientID           string
}

func NewClient() *Client {
	return NewClientWithAuth("conveyor-client")
}

func NewClientWithAuth(clientID string) *Client {
	// Initialize the config
	config.InitConfig()

	// Initialize HTTP client
	client := resty.New()
	client.SetBaseURL(viper.GetString("api.host"))
	client.SetHeader("Content-Type", "application/json")

	// Initialize authentication if enabled
	var authManager *ClientAuthManager
	var jwtToken string
	
	authEnabled := viper.GetBool("auth.enabled")
	if authEnabled {
		conveyorDataDir := viper.GetString("api.data")
		certDir := filepath.Join(conveyorDataDir, "certs")
		
		certManager := auth.NewCertificateManager(certDir)
		jwtManager := auth.NewJWTManager(certManager)
		tlsConfig := auth.NewTLSConfig(certManager)
		
		authManager = &ClientAuthManager{
			CertificateManager: certManager,
			JWTManager:         jwtManager,
			TLSConfig:          tlsConfig,
			ClientID:           clientID,
		}
		
		// Configure TLS if enabled
		if viper.GetBool("auth.tls.enabled") {
			clientTLSConfig, err := tlsConfig.CreateClientTLSConfig()
			if err == nil {
				client.SetTLSClientConfig(clientTLSConfig)
			}
		}
		
		// Generate JWT token if required
		if viper.GetBool("auth.jwt.required") {
			token, err := jwtManager.GenerateJWT(clientID)
			if err == nil {
				jwtToken = token
				client.SetAuthToken(token)
			}
		}
	}

	// Create a new client instance
	return &Client{
		HTTPClient:  client,
		AuthManager: authManager,
		jwtToken:    jwtToken,
	}
}

func (c *Client) GetAPIURL() string {
	// Return the API URL
	return viper.GetString("api.host")
}

// RefreshJWTToken generates a new JWT token
func (c *Client) RefreshJWTToken() error {
	if c.AuthManager == nil {
		return fmt.Errorf("authentication not configured")
	}
	
	token, err := c.AuthManager.JWTManager.GenerateJWT(c.AuthManager.ClientID)
	if err != nil {
		return fmt.Errorf("failed to generate JWT token: %w", err)
	}
	
	c.jwtToken = token
	c.HTTPClient.SetAuthToken(token)
	return nil
}

// GenerateJWTWithPrivateKey generates a JWT token using a specific private key
func (c *Client) GenerateJWTWithPrivateKey(privateKey *rsa.PrivateKey) error {
	if c.AuthManager == nil {
		return fmt.Errorf("authentication not configured")
	}
	
	token, err := c.AuthManager.JWTManager.GenerateJWTWithPrivateKey(c.AuthManager.ClientID, privateKey)
	if err != nil {
		return fmt.Errorf("failed to generate JWT token with private key: %w", err)
	}
	
	c.jwtToken = token
	c.HTTPClient.SetAuthToken(token)
	return nil
}

// GetCurrentJWTToken returns the current JWT token
func (c *Client) GetCurrentJWTToken() string {
	return c.jwtToken
}

// SetClientID sets the client ID for JWT generation
func (c *Client) SetClientID(clientID string) {
	if c.AuthManager != nil {
		c.AuthManager.ClientID = clientID
	}
}

// ConfigureTLS configures the client to use custom TLS settings
func (c *Client) ConfigureTLS(tlsConfig *tls.Config) {
	c.HTTPClient.SetTLSClientConfig(tlsConfig)
}
