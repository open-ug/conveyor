import json
import logging
from datetime import datetime, timezone

from nats.aio.client import Client as NATS


class DriverLogger:
    """Structured logger that dual-publishes to JetStream and Core NATS."""

    def __init__(self, run_id: str, driver_name: str, nats_conn: NATS):
        self.run_id = run_id
        self.driver_name = driver_name
        self._nc = nats_conn
        self._js = nats_conn.jetstream()

    async def log(self, message: str, *, pipeline: str | None = None, **labels) -> None:
        """
        Emit a structured JSON log entry.

        The entry is published best-effort to both:
          - JetStream subject ``logs.{run_id}``
          - Core NATS subject ``live.logs.{run_id}.{driver_name}``

        If one publish fails the other still proceeds.

        :param message:  Human-readable log message.
        :param pipeline: Optional pipeline identifier.
        :param labels:   Arbitrary key-value pairs flattened into the log entry.
        """
        entry: dict = {
            "runid": self.run_id,
            "driver": self.driver_name,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "message": message,
        }
        if pipeline is not None:
            entry["pipeline"] = pipeline
        entry.update(labels)

        payload = json.dumps(entry).encode()

        # JetStream publish (best-effort)
        try:
            await self._js.publish(f"logs.{self.run_id}", payload)
        except Exception:
            logging.debug("JetStream publish failed for run %s", self.run_id, exc_info=True)

        # Core NATS publish (best-effort)
        try:
            await self._nc.publish(f"live.logs.{self.run_id}.{self.driver_name}", payload)
        except Exception:
            logging.debug("Core NATS publish failed for run %s", self.run_id, exc_info=True)
