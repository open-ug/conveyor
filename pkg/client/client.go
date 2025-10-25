package client

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the API client.
type Client struct {
	HTTPClient *resty.Client
	APIURL     string

	authEnabled bool
	// FIX: Store the entire certificate chain
	certs      []*x509.Certificate
	privateKey crypto.PrivateKey
	tokenTTL   time.Duration
}

// NewClient initializes the API client.
// - apiURL: base API URL
// - authEnabled: true if client should use certificate-bound JWT auth
// - certPath: path to the client's certificate chain PEM file (leaf first)
// - keyPath: path to the client's private key PEM file
// - tokenTTL: duration of short-lived JWTs (optional, pass 0 for default 2m)
func NewClient(apiURL string, authEnabled bool, certPath, keyPath string, tokenTTL time.Duration) (*Client, error) {
	client := resty.New()
	client.SetBaseURL(apiURL)
	client.SetHeader("Content-Type", "application/json")

	c := &Client{
		HTTPClient:  client,
		APIURL:      apiURL,
		authEnabled: authEnabled,
		tokenTTL:    tokenTTL,
	}

	// If auth is enabled, load cert + key
	if authEnabled {
		if certPath == "" || keyPath == "" {
			return nil, errors.New("auth enabled but certPath or keyPath is empty")
		}
		certPEM, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("read cert file %s: %w", certPath, err)
		}
		keyPEM, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("read key file %s: %w", keyPath, err)
		}

		// FIX: Use the corrected parse function
		certs, key, err := parseCertAndKey(certPEM, keyPEM)
		if err != nil {
			return nil, fmt.Errorf("parse cert/key: %w", err)
		}
		// FIX: Store the entire chain
		c.certs = certs
		c.privateKey = key

		// default token TTL if not provided
		if c.tokenTTL <= 0 {
			c.tokenTTL = 2 * time.Minute
		}
	}

	return c, nil
}

// GetAPIURL returns the API URL.
func (c *Client) GetAPIURL() string {
	return c.APIURL
}
