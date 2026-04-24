use crate::runtime::manager::DriverManager;
use crate::runtime::types::Driver;
use crate::runtime::types::{Resource, ResourceCreateAPIResponse, ResourceDefinition};
use crate::runtime::utils::do_request;
use reqwest::Client as HttpClient;
use reqwest::header;

#[derive(Debug, Clone)]
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

    pub fn auth_enabled(&self) -> bool {
        self.config.auth_enabled()
    }

    pub fn new_driver_manager<'a>(&'a self, driver: Box<dyn Driver>) -> DriverManager<'a> {
        DriverManager {
            client: self,
            driver,
        }
    }
}

impl Client {
    pub async fn create_resource(
        &self,
        resource: &Resource,
    ) -> anyhow::Result<ResourceCreateAPIResponse> {
        let url = format!("{}/resources", self.api_endpoint);
        let payload = serde_json::to_string(resource)?;
        let response_text = do_request(&url, reqwest::Method::POST, Some(&payload)).await?;
        let response: ResourceCreateAPIResponse = serde_json::from_str(&response_text)?;
        Ok(response)
    }

    pub async fn get_resource(
        &self,
        resource_id: &str,
        resource_definition: &str,
    ) -> anyhow::Result<Resource> {
        let url = format!(
            "{}/resources/{}/{}",
            self.api_endpoint, resource_definition, resource_id
        );
        let response_text = do_request(&url, reqwest::Method::GET, None).await?;
        let resource: Resource = serde_json::from_str(&response_text)?;
        Ok(resource)
    }

    pub async fn update_resource(
        &self,
        resource_id: &str,
        data: &Resource,
    ) -> anyhow::Result<Resource> {
        let url = format!("{}/resources/{}", self.api_endpoint, resource_id);
        let payload = serde_json::to_string(data)?;
        let response_text = do_request(&url, reqwest::Method::PUT, Some(&payload)).await?;
        let resource: Resource = serde_json::from_str(&response_text)?;
        Ok(resource)
    }

    pub async fn delete_resource(&self, resource_id: &str) -> anyhow::Result<()> {
        let url = format!("{}/resources/{}", self.api_endpoint, resource_id);
        do_request(&url, reqwest::Method::DELETE, None).await?;
        Ok(())
    }

    pub async fn create_resource_definition(
        &self,
        definition: &ResourceDefinition,
    ) -> anyhow::Result<ResourceDefinition> {
        let url = format!("{}/resource-definitions", self.api_endpoint);
        let payload = serde_json::to_string(definition)?;
        let response_text = do_request(&url, reqwest::Method::POST, Some(&payload)).await?;
        let resource_definition: ResourceDefinition = serde_json::from_str(&response_text)?;
        Ok(resource_definition)
    }

    pub async fn get_resource_definition(
        &self,
        definition_id: &str,
    ) -> anyhow::Result<ResourceDefinition> {
        let url = format!(
            "{}/resource-definitions/{}",
            self.api_endpoint, definition_id
        );
        let response_text = do_request(&url, reqwest::Method::GET, None).await?;
        let resource_definition: ResourceDefinition = serde_json::from_str(&response_text)?;
        Ok(resource_definition)
    }

    pub async fn update_resource_definition(
        &self,
        definition_id: &str,
        data: &ResourceDefinition,
    ) -> anyhow::Result<ResourceDefinition> {
        let url = format!(
            "{}/resource-definitions/{}",
            self.api_endpoint, definition_id
        );
        let payload = serde_json::to_string(data)?;
        let response_text = do_request(&url, reqwest::Method::PUT, Some(&payload)).await?;
        let resource_definition: ResourceDefinition = serde_json::from_str(&response_text)?;
        Ok(resource_definition)
    }
}
