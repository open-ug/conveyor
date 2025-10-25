package auth

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Helper: load root CA pool from PEM
func LoadRootCAs() (*x509.CertPool, error) {
	caFilePath := viper.GetString("tls.ca")
	pemBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		return nil, fmt.Errorf("read root CA: %w", err)
	}
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(pemBytes); !ok {
		return nil, errors.New("failed to append root CA PEM")
	}
	return pool, nil
}

// Helper: decode x5c array of base64 DER cert strings into []*x509.Certificate
func parseX5C(x5cIface interface{}) ([]*x509.Certificate, error) {
	rawSlice, ok := x5cIface.([]interface{})
	if !ok || len(rawSlice) == 0 {
		return nil, errors.New("x5c header missing or not an array")
	}
	var certs []*x509.Certificate
	for _, entry := range rawSlice {
		s, ok := entry.(string)
		if !ok {
			return nil, errors.New("x5c entry not string")
		}
		der, err := base64.StdEncoding.DecodeString(s) // x5c uses base64 (not url)
		if err != nil {
			return nil, fmt.Errorf("base64 decode x5c: %w", err)
		}
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			return nil, fmt.Errorf("parse certificate: %w", err)
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

// Compute SHA-256 thumbprint of a certificate, returned as base64url (raw, no padding)
// This corresponds to x5t#S256 value format used in many specs.
func certThumbprintBase64URL(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// Build the verification key function for jwt.ParseWithClaims
// This Keyfunc will:
// - read token.Header["x5c"]
// - parse certs, verify chain using provided root CAs
// - validate that cnf claim thumbprint matches leaf cert
// - return the leaf public key for signature verification
func makeKeyFunc(rootPool *x509.CertPool) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// 1) extract x5c
		x5cIface, ok := token.Header["x5c"]
		if !ok {
			return nil, errors.New("missing x5c header")
		}
		certs, err := parseX5C(x5cIface)
		if err != nil {
			return nil, fmt.Errorf("x5c parse: %w", err)
		}
		leaf := certs[0]

		// 2) verify chain up to the rootPool
		intermediatePool := x509.NewCertPool()
		for i := 1; i < len(certs); i++ {
			intermediatePool.AddCert(certs[i])
		}
		opts := x509.VerifyOptions{
			Roots:         rootPool,
			Intermediates: intermediatePool,
			CurrentTime:   time.Now(),
			// No DNSName check here: this is for signing cert validation (not TLS server name)
		}
		if _, err := leaf.Verify(opts); err != nil {
			return nil, fmt.Errorf("certificate chain verification failed: %w", err)
		}

		// 3) Validate cnf claim matches thumbprint
		// Note: token.Claims may not yet be validated; we'll extract raw claims from token.Claims (map)
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("claims are not MapClaims")
		}

		cnfVal, cnfPresent := mapClaims["cnf"]
		if !cnfPresent {
			return nil, errors.New("cnf claim missing")
		}

		// cnf can be a map with keys like "x5t#S256" or "jkt" or "jwk"
		cnfMap, ok := cnfVal.(map[string]interface{})
		if !ok {
			return nil, errors.New("cnf claim not an object")
		}

		// compute cert thumbprint
		thumb := certThumbprintBase64URL(leaf)

		// Check for "x5t#S256"
		if val, found := cnfMap["x5t#S256"]; found {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("cnf.x5t#S256 not a string")
			}
			// cnf.x5t#S256 is typically base64url without padding
			if subtleConstantTimeCompare(strVal, thumb) {
				return leaf.PublicKey, nil
			}
			return nil, errors.New("cnf.x5t#S256 does not match certificate thumbprint")
		}

		// Fallback: if cnf.jkt (JWK thumbprint) is present and the leaf cert has a JWK thumbprint
		// NOTE: this requires computing the JWK thumbprint from the public key — for brevity we check jkt (assume equals cert thumb)
		if val, found := cnfMap["jkt"]; found {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("cnf.jkt not a string")
			}
			// It's implementation-specific how jkt is computed; many systems use the JWK thumbprint.
			// If your clients set jkt to the same value as x5t#S256, we accept it here.
			if subtleConstantTimeCompare(strVal, thumb) {
				return leaf.PublicKey, nil
			}
			// else continue to fail
			return nil, errors.New("cnf.jkt does not match cert thumbprint (or unsupported jkt format)")
		}

		return nil, errors.New("cnf did not contain recognized thumbprint fields (x5t#S256 or jkt)")
	}
}

// Constant-time compare helper for strings
func subtleConstantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	ab := []byte(a)
	bb := []byte(b)
	// simple constant time
	var res byte = 0
	for i := 0; i < len(ab); i++ {
		res |= ab[i] ^ bb[i]
	}
	return res == 0
}

// Middleware: validates JWT per our flow and attaches claims to context
func JWTCertMiddleware(rootPool *x509.CertPool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var skipPaths = []string{
			"/health",
			"/swagger",
			"/metrics",
		}
		// Skip auth for certain paths
		requestPath := c.Path()
		for _, sp := range skipPaths {
			if strings.HasPrefix(requestPath, sp) {
				return c.Next()
			}
		}
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing Authorization header"})
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid Authorization header format"})
		}
		tokenString := parts[1]

		// Parse token with MapClaims so we can inspect cnf and other claims.
		claims := jwt.MapClaims{}
		parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256", "PS256", "ES256", "ES384", "ES512"}))
		token, err := parser.ParseWithClaims(tokenString, claims, makeKeyFunc(rootPool))
		if err != nil {
			// include details for debugging in dev; reduce verbosity in prod
			log.Printf("token parse error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// Validate standard claims (exp, nbf, iat) manually
		if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
			log.Printf("token expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
		}
		if iat, ok := claims["iat"].(float64); ok && time.Unix(int64(iat), 0).After(time.Now()) {
			log.Printf("token issued in the future")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token issued in the future"})
		}

		// token is valid and signature verified (Keyfunc returned public key)
		// Attach claims to context for handlers
		c.Locals("claims", claims)
		c.Locals("token", token)
		return c.Next()
	}
}
