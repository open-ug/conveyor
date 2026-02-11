import json
import logging
from datetime import datetime, timezone

from nats.aio.client import Client as NATS

# Fields that labels must never override.
_RESERVED_FIELDS = frozenset({
    "run_id", "runid", "driver", "timestamp", "message", "pipeline",
})


class DriverLogger:
    """Structured logger that dual-publishes to JetStream and Core NATS."""

    def __init__(self, *args, **kwargs):
        # -----------------------------------------------------------
        # Backward-compatible constructor
        #   New: DriverLogger(run_id: str, driver_name: str, nats_conn)
        #   Old: DriverLogger(driver_name: str, labels: dict, nats_conn, loki_client)
        # -----------------------------------------------------------
        if args and len(args) >= 2 and isinstance(args[1], dict):
            # Old signature: (driver_name, labels, nats_conn, loki_client)
            driver_name = args[0]
            labels = args[1] if args[1] else {}
            nats_conn = args[2] if len(args) > 2 else kwargs.get("nats_conn")
            # Best-effort derive run_id from labels
            self.run_id = str(labels.get("run_id", "") or "")
            self.driver_name = driver_name
            self._nc = nats_conn
            # loki_client (args[3]) is accepted but ignored
        else:
            # New signature: (run_id, driver_name, nats_conn)
            self.run_id = args[0] if len(args) > 0 else kwargs.get("run_id", "")
            self.driver_name = args[1] if len(args) > 1 else kwargs.get("driver_name", "")
            self._nc = args[2] if len(args) > 2 else kwargs.get("nats_conn")

        self._js = self._nc.jetstream() if self._nc else None

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
            "run_id": self.run_id,
            "driver": self.driver_name,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "message": message,
        }
        if pipeline is not None:
            entry["pipeline"] = pipeline

        # Filter labels so reserved fields cannot be overwritten.
        for key, value in labels.items():
            if key not in _RESERVED_FIELDS:
                entry[key] = value

        # Best-effort JSON serialization – never crash async tasks.
        try:
            payload = json.dumps(entry, default=str).encode()
        except Exception:
            return

        # JetStream publish (best-effort)
        try:
            if self._js:
                await self._js.publish(f"logs.{self.run_id}", payload)
        except Exception:
            logging.debug("JetStream publish failed for run %s", self.run_id, exc_info=True)

        # Core NATS publish (best-effort)
        try:
            if self._nc:
                await self._nc.publish(f"live.logs.{self.run_id}.{self.driver_name}", payload)
        except Exception:
            logging.debug("Core NATS publish failed for run %s", self.run_id, exc_info=True)
