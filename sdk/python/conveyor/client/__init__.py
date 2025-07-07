"""
The client Library is a component of the driver runtime that contains funtions that interact with the Conveyor CI API Server. 
these are normal HTTP Requests to the API Server. It also has the ability to fetch Conveyor CI API Server metadata like the 
API Host and Port.

    See: https://conveyor.open.ug/blog/contributing-to-the-conveyor-ci-driver-runtime#client-library
"""

import logging
import requests
from requests.exceptions import RequestException


class Client:
    """
    The client Library is a component of the driver runtime that contains funtions that interact with the Conveyor CI API Server. these are normal HTTP Requests to the API Server. It also has the ability to fetch Conveyor CI API Server metadata like the API Host and Port.

    See: https://conveyor.open.ug/blog/contributing-to-the-conveyor-ci-driver-runtime#client-library
    """
    api_host: str
    api_port: int

    def __init__(self, api_host: str = "localhost", api_port: int = 8080):
        self.api_host = api_host
        self.api_port = api_port

    def get_api_url(self) -> str:
        """
        Returns the full API URL for the Conveyor CI API Server.
        """
        return f"http://{self.api_host}:{self.api_port}"

    def create_resource_definition(self, resource_definition: dict) -> requests.Response:
        """
        Creates a resource definition on the Conveyor CI API Server.

        :param resource_definition: The resource definition, which defines the schema of a resource based on the OpenAPI Specification.
            It determines how the resource will be defined, what properties it will have, and the validation schema for the resource.
        :return: The response from the Conveyor API Server.
        """

        url = f"{self.get_api_url()}/resource-definitions/"
        try:
            response = requests.post(
                url, json=resource_definition, timeout=None)
            response.raise_for_status()
            return response
        except RequestException as e:
            logging.error("Request failed: %s", e)
            raise

    def create_resource(self, resource: dict) -> requests.Response:
        """
        Creates a resource on the Conveyor CI API Server.

        :param resource: The resource,defines the automation pipeline (what to run, when, and how)
        :return: The response from the Conveyor API Server.
        """

        url = f"{self.get_api_url()}/resources/"
        try:
            response = requests.post(url, json=resource, timeout=None)
            response.raise_for_status()
            return response
        except RequestException as e:
            logging.error("Request failed: %s", e)
            raise

    def create_or_update_resource_definition(self, resource_definition: dict) -> requests.Response:
        """
        Creates or updates a resource definition on the Conveyor CI API Server.

        :param resource_definition: The resource definition, which defines the schema of a resource based on the OpenAPI Specification.
            It determines how the resource will be defined, what properties it will have
        :return: The response from the Conveyor API Server.
        """

        url = f"{self.get_api_url()}/resource-definitions/apply/"
        try:
            response = requests.post(
                url, json=resource_definition, timeout=None)
            response.raise_for_status()
            return response
        except RequestException as e:
            logging.error("Request failed: %s", e)
            raise
