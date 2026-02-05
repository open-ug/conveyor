package driverruntime

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

// Performs an HTTP request. If auth is enabled, it signs a JWT and sets Authorization header.
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

	if c.Options.AuthEnabled {
		tokenString, err := c.createSignedJWT()
		if err != nil {
			return fmt.Errorf("doRequest: failed to create signed jwt: %w", err)
		}
		// fmt.Printf("Generated JWT: %s\n", tokenString) // Optional: for debugging
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

// parseCertAndKey parses the full certificate chain and the private key.
// It validates that the private key matches the *first* certificate in the certPEM.
func parseCertAndKey(certPEM, keyPEM []byte) ([]*x509.Certificate, crypto.PrivateKey, error) {
	// --- 1. Parse Private Key ---
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, errors.New("no private key PEM block found")
	}

	var parsedKey any
	var err error
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		parsedKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		parsedKey, err = x509.ParseECPrivateKey(keyBlock.Bytes)
	case "PRIVATE KEY":
		parsedKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	default:
		parsedKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes) // Assume PKCS8
	}
	if err != nil {
		return nil, nil, fmt.Errorf("parse private key: %w", err)
	}

	var privKey crypto.PrivateKey
	var pubKey crypto.PublicKey

	switch k := parsedKey.(type) {
	case *rsa.PrivateKey:
		privKey = k
		pubKey = &k.PublicKey
	case *ecdsa.PrivateKey:
		privKey = k
		pubKey = &k.PublicKey
	default:
		return nil, nil, errors.New("unsupported private key type")
	}

	// --- 2. Parse ALL Certificates ---
	var certs []*x509.Certificate
	rest := certPEM
	for {
		var certBlock *pem.Block
		certBlock, rest = pem.Decode(rest)
		if certBlock == nil {
			break
		}
		if certBlock.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(certBlock.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("parse certificate in chain: %w", err)
			}
			certs = append(certs, cert)
		}
	}
	if len(certs) == 0 {
		return nil, nil, errors.New("no certificate PEM block found in cert file")
	}

	// --- 3. Validate Leaf Cert (first cert) matches Private Key ---
	leaf := certs[0]
	match := false
	switch pub := pubKey.(type) {
	case *rsa.PublicKey:
		if certPub, ok := leaf.PublicKey.(*rsa.PublicKey); ok {
			match = pub.N.Cmp(certPub.N) == 0 && pub.E == certPub.E
		}
	case *ecdsa.PublicKey:
		if certPub, ok := leaf.PublicKey.(*ecdsa.PublicKey); ok {
			match = pub.X.Cmp(certPub.X) == 0 && pub.Y.Cmp(certPub.Y) == 0 && pub.Curve == certPub.Curve
		}
	}

	if !match {
		return nil, nil, errors.New("private key does not match first certificate in PEM chain")
	}

	// Return the full chain and the private key
	return certs, privKey, nil
}

// builds and signs a short-lived JWT
func (c *Client) createSignedJWT() (string, error) {
	if !c.Options.AuthEnabled {
		return "", errors.New("auth not enabled")
	}
	// FIX: Check the certs slice
	if len(c.certs) == 0 || c.privateKey == nil {
		return "", errors.New("certificate chain or private key not configured")
	}
	// The leaf is the first certificate in our chain
	leaf := c.certs[0]

	now := time.Now().UTC()
	ttl := c.Options.TokenTTL
	if ttl <= 0 {
		ttl = 2 * time.Minute // Should match default in NewClient
	}

	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}

	// FIX: Use the leaf (c.certs[0]) for the thumbprint
	sum := sha256.Sum256(leaf.Raw)
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

	// FIX: Include the *entire* certificate chain in x5c
	x5cChain := make([]string, len(c.certs))
	for i, cert := range c.certs {
		// x5c uses standard base64, not URL encoding
		x5cChain[i] = base64.StdEncoding.EncodeToString(cert.Raw)
	}
	token.Header["x5c"] = x5cChain

	// Sign the token
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
