use crate::runtime::manager::DriverManager;
use reqwest::Client as HttpClient;
use reqwest::header;

#[derive(Debug)]
pub struct Client {
    /// The Conveyor CI API Server endpoint
    pub api_endpoint: String,
    /// The NATS server endpoint for receiving real-time updates
    pub nats_endpoint: String,

    /// Configuration for TLS settings, including mTLS credentials and trusted CA certificates
    pub config: ClientConfig,

    /// HTTP client for making API requests
    pub http_client: HttpClient,
}

/// Configuration for the Conveyor CI client, including TLS settings.
#[derive(Clone, Debug)]
pub struct ClientConfig {
    /// The PEM-encoded certificate chain (leaf first).
    /// If provided, this implicitly enables mTLS (mutual TLS).
    pub cert_pem: Option<Vec<u8>>,

    /// The PEM-encoded private key corresponding to the certificate.
    /// Required if `cert_pem` is provided.
    pub key_pem: Option<Vec<u8>>,

    /// The PEM-encoded Root CA certificate to trust.
    /// Used to verify the server's identity.
    pub root_ca_pem: Option<Vec<u8>>,
}

impl ClientConfig {
    /// Helper to determine if we should attempt authentication
    pub fn auth_enabled(&self) -> bool {
        self.cert_pem.is_some() && self.key_pem.is_some()
    }
}

impl Client {
    /// Creates a new Conveyor CI client with the specified API and NATS endpoints, and configuration.
    pub fn new(api_endpoint: &str, nats_endpoint: &str, config: ClientConfig) -> Self {
        let mut headers = header::HeaderMap::new();
        headers.insert(
            header::CONTENT_TYPE,
            header::HeaderValue::from_static("application/json"),
        );

        let http_client = HttpClient::builder()
            .default_headers(headers)
            .build()
            .expect("Failed to build HTTP client");

        Client {
            api_endpoint: api_endpoint.to_string(),
            nats_endpoint: nats_endpoint.to_string(),
            config,
            http_client,
        }
    }

    pub fn new_driver_manager(&self) -> DriverManager {
        DriverManager {
            // Initialize the driver manager with necessary state and configuration
        }
    }
}
