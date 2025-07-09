from abc import ABC, abstractmethod
from exceptions import DriverException

class Driver(ABC):

    def __init__(self, name, resources: list[str]):
        self.name = name
        self.resources = resources
        self.validate()


    @abstractmethod
    async def reconcile(self, payload: str, event: str, run_id: str, driver_logger):
        """
        Notes
        *****
        Declared async so that driver implementations can await
        non-blocking I/O (NATS publishing, Loki logging) without stalling the event loop

        param payload: The payload to be processed by the driver
        param event: The event type that triggered the driver
        param run_id: The unique identifier for the run
        param driver_logger: The logger instance for the driver

        """
        pass
    

    def validate(self) -> None:
        if not self.name or not str(self.name).strip():
            raise DriverException("Driver name must be a non-empty string.")

        if not self.resources or not isinstance(self.resources, list):
            raise DriverException("Resources must be a non-empty list.")
        
        for resource in self.resources:
            if not isinstance(resource, str) or not resource.strip():
                raise DriverException("Each resource must be a non-empty string.")