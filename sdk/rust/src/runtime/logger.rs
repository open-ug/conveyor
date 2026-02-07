use async_nats;
use async_nats::jetstream;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

pub struct DriverLogger<'a> {
    pub run_id: String,
    pub driver: String,
    pub labels: HashMap<String, String>,
    pub nats_connection: &'a async_nats::Client,
    pub jetstream_connection: &'a jetstream::Context,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Log {
    pub run_id: String,
    pub driver: String,
    pub pipeline: Option<String>,
    pub timestamp: String, // RFC 3339 format
    pub message: String,
}

impl<'a> DriverLogger<'a> {
    pub async fn log(&self, message: String) {
        let log_entry = Log {
            run_id: self.run_id.clone(),
            pipeline: None,
            driver: self.driver.clone(),
            timestamp: chrono::Utc::now().to_rfc3339(),
            message,
        };

        let payload = serde_json::to_vec(&log_entry).expect("Failed to serialize log entry");

        self.jetstream_connection
            .publish(format!("logs.{}", self.run_id), payload.clone().into())
            .await
            .expect("Failed to publish log message");

        if let Err(e) = self
            .nats_connection
            .publish(
                format!("live.logs.{}.{}", self.run_id, self.driver),
                payload.into(),
            )
            .await
        {
            eprintln!("Failed to publish log message: {:?}", e);
        }
    }
}
