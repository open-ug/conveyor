use crate::runtime::client::Client;
use crate::runtime::logger::DriverLogger;
use crate::runtime::types::{Driver, DriverMessage, DriverResultEvent, PipelineEvent, Resource};
use async_nats;
use async_nats::jetstream::consumer;
use async_nats::jetstream::consumer::{AckPolicy, DeliverPolicy};
use tokio_stream::StreamExt;

pub struct DriverManager<'a> {
    pub client: &'a Client,
    pub driver: Box<dyn Driver>,
}

impl<'a> DriverManager<'a> {
    pub async fn run(&self) {
        // connect to NATS
        let nats_client = async_nats::connect(&self.client.nats_endpoint)
            .await
            .expect("Failed to connect to NATS");

        let copy_client = nats_client.clone();

        let jetstream = async_nats::jetstream::new(nats_client);

        // subject list for receiving messages for this driver
        let mut filter_subjects = vec![];

        for resource in self.driver.resources() {
            filter_subjects.push(format!("resources.{}", resource));
            filter_subjects.push(format!(
                "drivers.{}.resources.{}",
                self.driver.name(),
                resource
            ));
        }

        let js_consumer = jetstream
            .create_consumer_on_stream(
                consumer::pull::Config {
                    durable_name: Some(self.driver.name()),
                    filter_subjects: filter_subjects,
                    ack_policy: AckPolicy::Explicit,
                    deliver_policy: DeliverPolicy::All,
                    ..Default::default()
                },
                "messages",
            )
            .await
            .expect("Failed to create consumer");

        let mut messages = js_consumer
            .messages()
            .await
            .expect("Failed to subscribe to messages");

        while let Some(message) = messages.next().await {
            match message {
                Ok(msg) => {
                    let driver_message: DriverMessage =
                        serde_json::from_slice(&msg.payload).expect("Failed to parse message");

                    let resource: Resource = serde_json::from_value(driver_message.payload.clone())
                        .expect("Failed to parse resource from message payload");

                    let logger = DriverLogger {
                        run_id: driver_message.run_id.clone(),
                        driver: self.driver.name(),
                        labels: std::collections::HashMap::new(),
                        nats_connection: &copy_client,
                        jetstream_connection: &jetstream,
                    };

                    let result = self.driver.reconcile(
                        driver_message.payload,
                        driver_message.event,
                        driver_message.run_id,
                        &logger,
                    );

                    let result_event = DriverResultEvent {
                        success: result.success,
                        message: result.message,
                        data: result.data,
                        driver: self.driver.name(),
                    };

                    result_event.publish(resource, jetstream.clone()).await;

                    // Acknowledge the message after processing
                    msg.ack().await.expect("Failed to acknowledge message");
                }
                Err(e) => {
                    eprintln!("Error receiving message: {:?}", e);
                }
            }
        }
    }
}

impl DriverResultEvent {
    pub async fn publish(&self, resource: Resource, js: async_nats::jetstream::Context) {
        let event = PipelineEvent {
            event: "update".to_string(),
            run_id: resource.id.clone(),
            resource: serde_json::to_value(&resource).expect("Failed to serialize resource"),
            driver_result: self.clone(),
        };

        let payload = serde_json::to_vec(&event).expect("Failed to serialize pipeline event");

        // Publish to a subject that includes the driver name and resource ID for easy filtering
        let subject = format!("events.{}.{}", self.driver, resource.id);
        js.publish(subject, payload.into())
            .await
            .expect("Failed to publish pipeline event");
    }
}
