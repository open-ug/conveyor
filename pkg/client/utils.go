package client

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// doRequest performs an HTTP request. If auth is enabled, it signs a JWT and sets Authorization header.
func (c *Client) doRequest(ctx context.Context, method, path string, body, dest any) error {
	if (method == http.MethodPost || method == http.MethodPut) && body == nil {
		return fmt.Errorf("doRequest: body cannot be nil for POST/PUT requests")
	}
	if dest == nil {
		return fmt.Errorf("doRequest: destination cannot be nil")
	}

	req := c.HTTPClient.R().SetContext(ctx)

	if body != nil {
		jsonMessage, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("doRequest: failed to marshal request body: %w", err)
		}
		req.SetBody(jsonMessage)
	}

	if c.authEnabled {
		tokenString, err := c.createSignedJWT()
		if err != nil {
			return fmt.Errorf("doRequest: failed to create signed jwt: %w", err)
		}
		req.SetHeader("Authorization", "Bearer "+tokenString)
	}

	resp, err := req.Execute(method, strings.TrimRight(c.HTTPClient.BaseURL, "/")+path)
	if err != nil {
		return fmt.Errorf("doRequest: failed to execute %s request: %w", method, err)
	}

	if resp.IsError() {
		return fmt.Errorf("doRequest: server returned %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	if err := json.Unmarshal(resp.Body(), dest); err != nil {
		return fmt.Errorf("doRequest: failed to unmarshal response body: %w", err)
	}

	return nil
}

// ------------------ certificate & JWT helpers ------------------

func parseCertAndKey(certPEM, keyPEM []byte) (*x509.Certificate, crypto.PrivateKey, error) {
	// parse certificate
	var certBlock *pem.Block
	rest := certPEM
	for {
		certBlock, rest = pem.Decode(rest)
		if certBlock == nil {
			break
		}
		if certBlock.Type == "CERTIFICATE" {
			break
		}
	}
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return nil, nil, errors.New("no certificate PEM block found")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse certificate: %w", err)
	}

	// parse private key
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, errors.New("no private key PEM block found")
	}

	var parsedKey any
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		parsedKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		parsedKey, err = x509.ParseECPrivateKey(keyBlock.Bytes)
	case "PRIVATE KEY":
		parsedKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	default:
		parsedKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("parse private key: %w", err)
	}

	switch k := parsedKey.(type) {
	case *rsa.PrivateKey:
		return cert, k, nil
	case *ecdsa.PrivateKey:
		return cert, k, nil
	default:
		return nil, nil, errors.New("unsupported private key type")
	}
}

// createSignedJWT builds and signs a short-lived JWT bound to the client's certificate
func (c *Client) createSignedJWT() (string, error) {
	if !c.authEnabled {
		return "", errors.New("auth not enabled")
	}
	if c.cert == nil || c.privateKey == nil {
		return "", errors.New("certificate or private key not configured")
	}

	now := time.Now().UTC()
	ttl := c.tokenTTL
	if ttl <= 0 {
		ttl = 2 * time.Minute
	}

	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}

	sum := sha256.Sum256(c.cert.Raw)
	thumb := base64.RawURLEncoding.EncodeToString(sum[:])
	claims["cnf"] = map[string]string{"x5t#S256": thumb}

	var signingMethod jwt.SigningMethod
	switch k := c.privateKey.(type) {
	case *rsa.PrivateKey:
		signingMethod = jwt.SigningMethodRS256
	case *ecdsa.PrivateKey:
		switch k.Curve.Params().BitSize {
		case 256:
			signingMethod = jwt.SigningMethodES256
		case 384:
			signingMethod = jwt.SigningMethodES384
		case 521:
			signingMethod = jwt.SigningMethodES512
		default:
			return "", errors.New("unsupported ecdsa curve")
		}
	default:
		return "", errors.New("unsupported private key type")
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	leafB64 := base64.StdEncoding.EncodeToString(c.cert.Raw)
	token.Header["x5c"] = []string{leafB64}

	var signed string
	var err error
	switch k := c.privateKey.(type) {
	case *rsa.PrivateKey:
		signed, err = token.SignedString(k)
	case *ecdsa.PrivateKey:
		signed, err = token.SignedString(k)
	default:
		return "", errors.New("unsupported key type for SignedString")
	}
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}
