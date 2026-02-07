use crate::runtime::logger::DriverLogger;
use serde::{Deserialize, Serialize};
#[derive(Serialize, Deserialize, Debug)]
pub struct DriverResult {
    pub success: bool,
    pub message: Option<String>,
    pub data: Option<serde_json::Value>, // Placeholder for any data returned by the driver
}

#[derive(Serialize, Deserialize, Debug)]
pub struct DriverMessage {
    pub event: String,
    pub payload: serde_json::Value, // The actual content of the message, which can be any JSON value
    pub id: Option<String>,         // Optional ID for correlating requests and responses

    // Unique identifier for the driver run, useful for tracking and logging
    pub run_id: String,
}

#[derive(Clone, Serialize, Deserialize, Debug)]
pub struct DriverResultEvent {
    pub success: bool,
    pub message: Option<String>,
    pub data: Option<serde_json::Value>, // Placeholder for any data returned by the driver
    pub driver: String,                  // Identifier for the driver that produced this result
}

pub trait Driver {
    fn name(&self) -> String;

    fn resources(&self) -> Vec<String>;

    fn reconcile(
        &self,
        payload: serde_json::Value,
        event: String,
        run_id: String,
        logger: &DriverLogger,
    ) -> DriverResult;
}

#[derive(Serialize, Deserialize, Debug)]
pub struct PipelineEvent {
    pub event: String, // e.g "create", "update", "delete"
    pub run_id: String,
    pub resource: serde_json::Value, // Placeholder for the actual resource structure
    pub driver_result: DriverResultEvent,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Resource {
    pub id: String,
    pub name: String,
    pub pipeline: Option<String>,
    pub resource: String,
    pub metadata: std::collections::HashMap<String, String>,
    pub spec: serde_json::Value, // Placeholder for the actual spec structure
}
