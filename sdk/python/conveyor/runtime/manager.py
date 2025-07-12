from typing import List
from nats.js.api import ConsumerConfig, AckPolicy
import logging
import json
from conveyor.runtime.exceptions import DriverRuntimeException
from driver_logger import DriverLogger,LokiClient



class DriverManager:

    """
    Manage JetStream ,Loki subscription and dispatch events to a Conveyor driver.
    """

    def __init__(self, driver, events: List[str], nats_connection,loki_client:LokiClient):
        self.driver = driver
        self.events = events
        self.nats_connection = nats_connection
        self.loki_client = loki_client
        self.validate()
    


    def validate(self):
        if not self.driver:
            raise DriverRuntimeException("Invalid driver instance provided.")

        if not self.events:
            raise DriverRuntimeException("Events must be a non-empty list.")

        for event in self.events:
            if not isinstance(event, str) or not event.strip():
                raise DriverRuntimeException("Each event must be a non-empty string.")
        
        if not self.nats_connection:
            raise DriverRuntimeException("NATS connection is required for the Driver Manager to start.")

        if not self.loki_client:
            raise DriverRuntimeException("Loki client is required for logging.")

    async def run(self):
        jetstream =  self.nats_connection.jetstream()
        filter_subjects = [f"resources.{r}" for r in self.driver.resources]
        

        cfg = ConsumerConfig(
            durable_name=self.driver.name,
            ack_policy=AckPolicy.EXPLICIT,
            max_ack_pending=1,
            filter_subjects=filter_subjects,
        )

        #create or update consumer on the stream
        self.sub = await jetstream.subscribe(
        subject="resources.*",      # a broad wildcard
        durable=self.driver.name,   # binds to / creates the consumer
        stream="messages",          # stream name   
        manual_ack=True,            # explicit ack mode
        config=cfg,                 # ConsumerConfig 
        cb=self.handle_message,     # async msg callback
        )

        logging.info("Consumer %s created or bound with filter subjects: %s", self.driver.name, filter_subjects)
        logging.info("Driver Manager is running for driver: %s", self.driver.name)




    async def handle_message(self, msg):

        """
         Async callback that handles incoming messages from NATS, processes them accordingly
        :param msg: The message received from NATS.
        """
        try:
            message = json.loads(msg.data)        
            event = message.get("event")

            if not event:  
                logging.warning("Message missing event field")
                raise Exception("Message missing event field for driver %s." % self.driver.name)

            logging.debug("Received message for event: %s", event)

            # quick filter, skip irrelevant events
            if "*" not in self.events and event not in self.events:
                await msg.ack()         # not an event we care about, ack and return
                return

            logger_context = {
                "event": event,
                "id": message.get("id"),
                "run_id": message.get("run_id"),
            }

            logging.debug("Processing event: %s with context: %s", event, logger_context)

            # Create driver logger with context
            driver_logger = DriverLogger(
                self.driver.name, logger_context, self.nats_connection, self.loki_client
            )

            await self.driver.reconcile(
                message.get("payload"),
                event,
                message.get("run_id"),
                driver_logger,
            )
            await msg.ack()                          # success path
            logging.debug("Acknowledged message for event: %s", event)

        except DriverRuntimeException as e:
            logging.error("Driver error: %s", e)
            await msg.nak()                          # redeliver later

        except json.JSONDecodeError as e:
            logging.error("Bad JSON: %s", e)
            await msg.term()                         # drop bad message

        except Exception as e:
            logging.exception("Unhandled error: %s", e)
            await msg.term()                         # at the moment, safest fallback

    async def shutdown(self):
        """ 
        Gracefully shuts down the Driver Manager, unsubscribing from NATS and closing connections.
        """
        logging.info("Shutting down Driver Manager for driver: %s", self.driver.name)
        # Remove interest in subscription.
        await self.sub.unsubscribe()
        # Terminate connection to NATS.
        await self.nats_connection.drain()
        await self.loki_client.close()
        logging.info("Driver Manager for driver %s has been shut down.", self.driver.name)
