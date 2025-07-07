import httpx
from httpx import HTTPStatusError, RequestError
import time
import logging
import json
from nats.aio.client import Client as NATS
from conveyor.runtime.exceptions import InvalidRunIDException


class LokiClient:
    def __init__(self, url: str = "http://localhost:3100"):
        self.url = url
        self.client = httpx.AsyncClient()  # connection reuse

    async def push_log(self, labels: dict, message: str) -> None:
        """
        Sends the logs to Loki via async HTTP requests for storage or log aggregation.

        :param labels: key-value pairs that describe extra metadata about the log.
        :param message: The log message to be sent to Loki.
        :raises HTTPError: If the response from Loki indicates an error (4xx or 5xx).
        :raises RequestError: If the request to Loki fails due to network issues.
        :return: None
        """
        log_entry = {
            "streams": [
                {"stream": labels, "values": [[str(int(time.time() * 1e9)), message]]}
            ]
        }

        url = f"{self.url}/loki/api/v1/push"

        try:
            response = await self.client.post(url, json=log_entry, timeout=None)
            response.raise_for_status()
            logging.info("Log sent to Loki: %s | Labels: %s", message, labels)
            logging.debug(
                "Loki response status code: %s, response: %s, labels: %s",
                response.status_code,
                response.text,
                labels,
            )

        except HTTPStatusError as e:
            status_code = getattr(e.response, "status_code", "N/A")
            reason = getattr(e.response, "reason_phrase", "N/A")
            text = getattr(e.response, "text", "N/A")
            logging.error(
                "HTTPError %s: %s | Reason: %s | Response: %s | Labels: %s",
                status_code,
                e,
                reason,
                text,
                labels,
            )
            raise

        except RequestError as e:
            logging.error("Failed to push log to Loki: %s | Labels: %s", e, labels)
            raise

    # TODO: Close the LokiClient (httpx.AsyncClient) when the app is shutting down
    # This prevents socket leaks and ensures graceful shutdown.
    # Never close after every log — keep the client open for the app lifetime.
    async def close(self):
        await self.client.aclose()


class DriverLogger:

    def __init__(
        self, driver_name: str, labels: dict, nats_conn: NATS, loki_client: LokiClient
    ):
        """
        Initializes the DriverLogger with the driver name, labels, NATS connection, and Loki client.
        """
        self.driver_name = driver_name
        self.labels = labels
        self.nats_conn = nats_conn
        self.loki_client = loki_client

    async def log(self, labels: dict, message: str) -> None:
        """
        Merges the logger's labels with the provided labels and logs the message to the Loki client and NATS.
        :param labels: Additional labels to merge with the logger's labels.
        :param message: The log message to be sent.
        :raises TypeError: If the labels are not serializable eg a set , function .
        :raises UnicodeEncodeError: If the message contains non-ASCII characters.
        """

        merged_labels = {"driver": self.driver_name}
        merged_labels.update(self.labels)
        merged_labels.update(labels)

        run_id = merged_labels.get("run_id", "unknown_run")

        if not run_id or not str(run_id).strip():
            logging.error("Missing run_id from labels: %s", merged_labels)
            raise InvalidRunIDException("Missing run_id from labels.")

        await self.loki_client.push_log(merged_labels, message)

        timestamp = [str(int(time.time())), message]

        logging.debug("encoding data to bytes for NATS : %s", timestamp)
        try:
            payload = json.dumps(timestamp).encode()
        except (TypeError, UnicodeEncodeError) as e:
            logging.error(
                "Failed to encode log payload for NATS: %s | Data: %s", e, timestamp
            )
            raise

        logging.info("Publishing log to NATS: %s ", timestamp)
        await self.nats_conn.publish(
            f"driver:{self.driver_name}:logs:{run_id}", payload
        )
