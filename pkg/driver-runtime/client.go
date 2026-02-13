package driverruntime

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the API client.
type Client struct {
	HTTPClient *resty.Client
	APIURL     string
	NatsURL    string

	// FIX: Store the entire certificate chain
	certs      []*x509.Certificate
	privateKey crypto.PrivateKey
	Options    ConfigOptions
}

type ConfigOptions struct {
	// AuthEnabled indicates if the client should use certificate-bound JWT auth
	// If true, cert and key must be provided.
	// If false, the client will not use auth.
	AuthEnabled bool
	// Cert is the PEM-encoded certificate chain (leaf first).
	// This is required if AuthEnabled is true.
	// It should be a byte slice containing the PEM data.
	Cert []byte
	// Key is the PEM-encoded private key.
	// This is required if AuthEnabled is true.
	// It should be a byte slice containing the PEM data.
	Key []byte

	// RootCA is an PEM-encoded CA certificate to trust when connecting to the API.
	// This is used in NATS connections to verify the server's TLS certificate.
	RootCA []byte
	// TokenTTL is the duration for which short-lived JWTs are valid.
	// If not set, defaults to 2 minutes.
	TokenTTL time.Duration
}

// NewClient initializes the API client.
// It takes the API endpoint, NATS endpoint, and configuration options.
func NewClient(
	ApiEndpoint string,
	NatsEndpoint string,
	options ConfigOptions,
) (*Client, error) {
	client := resty.New()
	client.SetBaseURL(ApiEndpoint)
	client.SetHeader("Content-Type", "application/json")

	c := &Client{
		HTTPClient: client,
		APIURL:     ApiEndpoint,
		NatsURL:    NatsEndpoint,
		Options:    options,
	}

	// If auth is enabled, load cert + key
	if options.AuthEnabled {
		if options.Cert == nil || options.Key == nil {
			return nil, errors.New("auth enabled but cert or key is empty")
		}

		// FIX: Use the corrected parse function
		certs, key, err := parseCertAndKey(options.Cert, options.Key)
		if err != nil {
			return nil, fmt.Errorf("parse cert/key: %w", err)
		}
		// FIX: Store the entire chain
		c.certs = certs
		c.privateKey = key

		// default token TTL if not provided
		if c.Options.TokenTTL <= 0 {
			c.Options.TokenTTL = 2 * time.Minute
		}
	}

	return c, nil
}

// GetAPIURL returns the API URL.
func (c *Client) GetAPIURL() string {
	return c.APIURL
}
