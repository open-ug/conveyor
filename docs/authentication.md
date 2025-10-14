# Authentication Documentation

## Overview

Conveyor CI API Server implements a comprehensive authentication system using TLS certificates and JWT tokens to ensure that only trusted clients can connect and interact with the system. This document provides detailed information about how to set up and use the authentication system.

## Architecture

The authentication system consists of several key components:

1. **Certificate Management**: Handles loading and validation of TLS certificates
2. **JWT Token Management**: Generates and validates JWT tokens
3. **TLS Configuration**: Sets up secure TLS connections
4. **Authentication Middleware**: Protects API routes with authentication

## Authentication Flow

1. **Certificate Setup**: Both server and clients have certificates signed by the same Certificate Authority (CA)
2. **TLS Connection**: Clients establish TLS connections using their certificates
3. **JWT Generation**: Clients generate JWT tokens signed with their private keys
4. **API Requests**: Clients include JWT tokens in the `Authorization` header
5. **Token Validation**: Server validates JWT tokens using the CA certificate

## Certificate Generation

### 1. Certificate Authority (CA) Setup

First, generate a Certificate Authority private key and self-signed certificate:

```bash
# Generate CA private key
openssl genrsa -out ca.key 4096

# Generate self-signed CA certificate
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 \
  -out ca.pem -subj "/CN=Conveyor-CA"
```

### 2. API Server Certificates

Generate certificates for the API server:

```bash
# Generate API server private key
openssl genrsa -out server.key 2048

# Generate API server Certificate Signing Request (CSR)
openssl req -new -key server.key -out server.csr -subj "/CN=Conveyor-API-Server"

# Sign API server certificate with the CA
openssl x509 -req -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial \
  -out server.pem -days 365 -sha256

# Clean up CSR file
rm server.csr
```

### 3. Client Certificates

Generate certificates for clients:

```bash
# Generate client private key
openssl genrsa -out client.key 2048

# Generate client Certificate Signing Request (CSR)
openssl req -new -key client.key -out client.csr -subj "/CN=Conveyor-Client"

# Sign client certificate with the CA
openssl x509 -req -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial \
  -out client.pem -days 365 -sha256

# Clean up CSR file
rm client.csr
```

### 4. Certificate File Summary

After generation, you should have the following files:

| File | Description | Security Level |
|------|-------------|----------------|
| `ca.pem` | Certificate Authority certificate | Public |
| `ca.key` | Certificate Authority private key | **HIGHLY SENSITIVE** |
| `server.pem` | API server certificate | Public |
| `server.key` | API server private key | **SENSITIVE** |
| `client.pem` | Client certificate | Public |
| `client.key` | Client private key | **SENSITIVE** |

⚠️ **Security Note**: Keep private keys (`*.key` files) secure and never share them. The CA private key is especially critical and should be stored securely.

## Certificate Storage

Certificates are stored in different locations based on user context:

- **Root user**: `/var/lib/conveyor/certs/`
- **Non-root user**: `$HOME/.local/share/conveyor/certs/`

The directory structure should be:

```
certs/
├── ca.pem          # Certificate Authority certificate
├── ca.key          # Certificate Authority private key (secure)
├── server.pem      # API server certificate
├── server.key      # API server private key (secure)
├── client.pem      # Client certificate
└── client.key      # Client private key (secure)
```

## Configuration

The authentication system can be configured using environment variables or configuration files:

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CONVEYOR_AUTH_ENABLED` | Enable/disable authentication | `true` |
| `CONVEYOR_TLS_ENABLED` | Enable/disable TLS | `true` |
| `CONVEYOR_JWT_REQUIRED` | Require JWT tokens | `true` |
| `CONVEYOR_DATA_DIR` | Data directory path | Auto-detected |

### Configuration Examples

#### Full Authentication (Production)
```bash
export CONVEYOR_AUTH_ENABLED=true
export CONVEYOR_TLS_ENABLED=true
export CONVEYOR_JWT_REQUIRED=true
```

#### Development Mode (No Authentication)
```bash
export CONVEYOR_AUTH_ENABLED=false
export CONVEYOR_TLS_ENABLED=false
export CONVEYOR_JWT_REQUIRED=false
```

#### TLS Only (No JWT)
```bash
export CONVEYOR_AUTH_ENABLED=true
export CONVEYOR_TLS_ENABLED=true
export CONVEYOR_JWT_REQUIRED=false
```

## API Usage

### Server Side

When starting the API server with authentication enabled:

1. Server loads certificates from the configured directory
2. TLS listener is created with client certificate verification
3. Authentication middleware is applied to protected routes
4. JWT tokens are validated for each request

### Client Side

#### Creating an Authenticated Client

```go
import "github.com/open-ug/conveyor/pkg/client"

// Create client with default authentication
client := client.NewClient()

// Or create client with specific client ID
client := client.NewClientWithAuth("my-client-id")
```

#### Making Authenticated Requests

```go
// The client automatically includes JWT tokens in requests
response, err := client.HTTPClient.R().Get("/protected-endpoint")
```

#### Manual JWT Token Management

```go
// Refresh JWT token
err := client.RefreshJWTToken()

// Use custom private key for JWT generation
err := client.GenerateJWTWithPrivateKey(privateKey)

// Get current token
token := client.GetCurrentJWTToken()
```

### HTTP API Usage

When making direct HTTP requests, include the JWT token in the Authorization header:

```bash
# Generate JWT token (using client with certificates)
TOKEN=$(conveyor-client generate-jwt --client-id "my-client")

# Make authenticated request
curl -H "Authorization: Bearer $TOKEN" \
     --cert client.pem \
     --key client.key \
     --cacert ca.pem \
     https://api.conveyor.example.com/protected-endpoint
```

## Public Endpoints

The following endpoints are publicly accessible and do not require authentication:

- `/health` - Health check endpoint
- `/metrics` - Prometheus metrics
- `/swagger/*` - API documentation
- `/docs/*` - Additional documentation
- `/` - Root endpoint with basic information

## Security Considerations

### Certificate Security

1. **Private Key Protection**: Store private keys in secure locations with appropriate file permissions (600)
2. **CA Key Security**: The CA private key should be stored offline and used only for signing certificates
3. **Certificate Rotation**: Regularly rotate certificates before expiration
4. **Certificate Revocation**: Implement certificate revocation if a private key is compromised

### JWT Token Security

1. **Token Expiration**: JWT tokens expire after 24 hours by default
2. **Signature Verification**: All tokens are verified using the CA certificate
3. **Token Transmission**: Tokens are transmitted over TLS-encrypted connections only

### Network Security

1. **TLS Configuration**: Uses modern TLS versions (1.2+) with strong cipher suites
2. **Certificate Validation**: Both server and client certificates are validated
3. **Mutual TLS**: Server requires and verifies client certificates

## Troubleshooting

### Common Issues

#### 1. "Client certificate validation failed"

**Cause**: Client certificate is not signed by the expected CA or is invalid.

**Solution**:
- Verify client certificate is signed by the same CA as server certificate
- Check certificate expiration dates
- Ensure certificate files are readable

#### 2. "JWT token validation failed"

**Cause**: JWT token is invalid, expired, or signed with wrong key.

**Solution**:
- Generate a new JWT token
- Verify client certificate and private key are correctly matched
- Check system clock synchronization

#### 3. "TLS handshake failed"

**Cause**: TLS configuration mismatch or certificate issues.

**Solution**:
- Verify server certificate matches expected server name
- Check CA certificate is properly installed
- Ensure TLS versions are compatible

#### 4. "Missing certificates"

**Cause**: Certificate files not found in expected location.

**Solution**:
- Verify certificate directory path
- Check file permissions
- Generate certificates if missing

### Debug Commands

```bash
# Verify certificate
openssl x509 -in client.pem -text -noout

# Check certificate against CA
openssl verify -CAfile ca.pem client.pem

# Test TLS connection
openssl s_client -connect api.conveyor.example.com:8080 \
                 -cert client.pem -key client.key -CAfile ca.pem

# Decode JWT token (header and payload only)
echo "JWT_TOKEN" | cut -d. -f1-2 | sed 's/\./\n/' | base64 -d
```

## Integration Examples

### Go Client Integration

```go
package main

import (
    "log"
    "github.com/open-ug/conveyor/pkg/client"
)

func main() {
    // Create authenticated client
    client := client.NewClientWithAuth("integration-client")
    
    // Make request to protected endpoint
    resp, err := client.HTTPClient.R().
        SetResult(&map[string]interface{}{}).
        Get("/api/resources")
    
    if err != nil {
        log.Fatal("Request failed:", err)
    }
    
    log.Printf("Response: %+v", resp.Result())
}
```

### Shell Script Integration

```bash
#!/bin/bash

# Configuration
CLIENT_ID="shell-client"
API_BASE="https://api.conveyor.example.com"
CERT_DIR="/path/to/certs"

# Generate JWT token (pseudo-code, actual implementation needed)
TOKEN=$(generate_jwt_token "$CLIENT_ID" "$CERT_DIR/client.key")

# Make authenticated API call
curl -s \
  -H "Authorization: Bearer $TOKEN" \
  --cert "$CERT_DIR/client.pem" \
  --key "$CERT_DIR/client.key" \
  --cacert "$CERT_DIR/ca.pem" \
  "$API_BASE/api/pipelines" | jq .
```

## Monitoring and Logging

### Authentication Metrics

The system exposes Prometheus metrics for monitoring authentication:

- `auth_requests_total` - Total authentication attempts
- `auth_success_total` - Successful authentications
- `auth_failures_total` - Failed authentications
- `jwt_validations_total` - JWT token validations
- `tls_handshakes_total` - TLS handshake attempts

### Log Messages

Authentication events are logged with appropriate levels:

- **INFO**: Successful authentications
- **WARN**: Authentication attempts without proper credentials
- **ERROR**: Authentication failures and certificate validation errors

## Migration Guide

### Upgrading from Unauthenticated Version

1. **Generate Certificates**: Create CA and certificates as described above
2. **Update Configuration**: Enable authentication in configuration
3. **Update Clients**: Modify clients to use authenticated client SDK
4. **Test Authentication**: Verify authentication works before deploying to production
5. **Deploy**: Deploy server and client updates together

### Gradual Migration

To enable gradual migration, you can:

1. Start with `CONVEYOR_JWT_REQUIRED=false` to make JWT optional
2. Update all clients to support authentication
3. Enable `CONVEYOR_JWT_REQUIRED=true` once all clients are updated

## Best Practices

1. **Use Strong Passwords**: When prompted during certificate generation
2. **Automate Certificate Management**: Use tools like cert-manager for Kubernetes
3. **Monitor Certificate Expiration**: Set up alerts for certificate expiration
4. **Implement Certificate Rotation**: Regularly rotate certificates
5. **Use Hardware Security Modules**: For high-security environments, consider HSMs
6. **Implement Audit Logging**: Log all authentication events for security audit
7. **Test Disaster Recovery**: Have procedures for certificate recovery

## Support and Resources

- **GitHub Issues**: Report bugs and request features
- **Documentation**: Additional documentation at [docs.conveyor.open.ug](https://docs.conveyor.open.ug)
- **Community**: Join discussions in the community forums
- **Security**: Report security issues privately to security@conveyor.open.ug