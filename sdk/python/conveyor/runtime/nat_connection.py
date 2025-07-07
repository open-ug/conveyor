import nats
from typing import List, Union
import logging


#reusable wrapper for nats connection
async def  connect_to_nats(
    servers: Union[str, List[str]] = "nats://localhost:4222",
    **options
):
    """
    Reusable connection function. Usage:
    nc = await connect_to_nats(["url1", "url2"], **options)
    or
    nc = await connect_to_nats("nats://localhost:4222", **options)
    for more you can refer to the nats documentation. 
    """
    logging.info("Connecting to NATS server(s): %s", servers)
    if isinstance(servers, str):
        servers = [servers]
    return await nats.connect(servers=servers, **options)