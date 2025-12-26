---
sidebar_position: 3
---

# Security

Conveyor CI is designed to be secure by default. It uses system Security frameworks like AppAmmor for Mandatory Access Control and has a Stateless Cryptographic based Authentication system.

## Authentication

Conveyor CI API server uses a Stateless Cryptographic Proof-of-Possession (PoP) Authentication System based on JWTs and X.509 certificates and the NATS event broker users the [TLS Authentication](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/tls_mutual_auth) supported by NATS. The same certificates and private key are used in Authenticating both systems. By default authentication is disabled but can be enabled by setting the `api.auth_enabled` config value to `true`.

### Client Certificates

To authenticate with Conveyor CI, a client must possess a cryptographic identity consisting of a private key (client.key), a signed client certificate (client.crt), and the trusted Certificate Authority (ca.pem). These credentials are issued by the Conveyor CI operator. The recommended flow is:

1. the client generates a private key locally and creates a certificate signing request (CSR)
2. the CSR is submitted to the Conveyor CI administrator
3. Conveyor CI administrator signs the request using its internal Certificate Authority and returns a signed `client.crt` along with the trusted `ca.pem`

The private key never leaves the client machine and must be stored securely. These same files are used both for API authentication and for [NATS mutual TLS (mTLS)](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/tls_mutual_auth) when connecting to Conveyor CI services.

### High-Level Overview of API Server Authentication

1. Clients authenticate using JWTs signed with a private key
2. The corresponding X.509 certificate chain is embedded in the JWT header
3. The server verifies:
   - The certificate chain against a trusted root CA
   - That the token is cryptographically bound to the certificate
   - That the JWT signature is valid
4. No sessions or tokens are stored server-side

> **Note:** Development SDK libraries often handle this implementation for you so you under the hood
