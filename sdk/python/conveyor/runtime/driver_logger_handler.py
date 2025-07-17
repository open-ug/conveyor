import logging
import asyncio
import sys

class DriverLoggerHandler(logging.Handler):
    def __init__(self, driver_logger):
        super().__init__()
        self.driver_logger = driver_logger
    
    def emit(self, record):
        labels = {
            "level": record.levelname,
            "module": getattr(record, 'module', 'unknown')
        }
        
        message = record.getMessage()
        
        asyncio.create_task(self.driver_logger.log(labels, message))



class PrintHandler:
    def __init__(self, driver_logger):
        self.driver_logger = driver_logger
        self.original_stdout = sys.stdout
    
    def write(self, text):
        if text.strip():
            labels = {"level": "INFO", "source": "print"}
            asyncio.create_task(self.driver_logger.log(labels, text.strip()))
        return len(text)
    
    def flush(self):
        pass