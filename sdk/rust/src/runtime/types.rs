use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct DriverResult {
    pub success: bool,
    pub message: Option<String>,
    pub data: Option<serde_json::Value>, // Placeholder for any data returned by the driver
}

#[derive(Serialize, Deserialize, Debug)]
struct DriverMessage {
    pub event: String,
    pub payload: serde_json::Value, // The actual content of the message, which can be any JSON value
    pub id: Option<String>,         // Optional ID for correlating requests and responses

    // Unique identifier for the driver run, useful for tracking and logging
    pub run_id: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct DriverResultEvent {
    pub success: bool,
    pub message: Option<String>,
    pub data: Option<serde_json::Value>, // Placeholder for any data returned by the driver
    pub driver: String,                  // Identifier for the driver that produced this result
}

pub struct DriverManager {
    // Fields for managing drivers, such as a registry of available drivers,
    // configuration settings, and any necessary state for driver lifecycle management.
}

pub struct DriverLogger {
    // Fields for managing driver logs, such as log storage, formatting options,
    // and any necessary state for log lifecycle management.
}

pub struct Log {
    pub runid: String,
    pub driver: String,
    pub pipeline: Option<String>,
    pub timestamp: String, // RFC 3339 format
    pub message: String,
}
