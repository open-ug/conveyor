import logging
import asyncio


class DriverLoggerHandler(logging.Handler):
    """Logging handler that forwards stdlib log records to a DriverLogger."""

    def __init__(self, driver_logger):
        super().__init__()
        self.driver_logger = driver_logger

    def emit(self, record):
        message = record.getMessage()
        labels = {
            "level": record.levelname,
            "module": getattr(record, "module", "unknown"),
        }
        try:
            asyncio.create_task(self.driver_logger.log(message, **labels))
        except RuntimeError:
            pass  # no running event loop – best-effort