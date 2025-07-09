class DriverRuntimeException(Exception):
    """
    Base class for all exceptions raised by the Conveyor CI Driver Runtime.
    """

    def __init__(self, message: str):
        super().__init__(message)
        self.message = message


class InvalidRunIDException(DriverRuntimeException):
    """Raised when the run_id is missing from the labels ."""

    def __init__(self, message: str):
        super().__init__(message)
        self.message = message


class DriverException(DriverRuntimeException):
    """Raised with issues to do with driver."""

    def __init__(self, message: str):
        super().__init__(message)
        self.message = message
