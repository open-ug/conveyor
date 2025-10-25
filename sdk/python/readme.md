<h1 align="center" style="border-bottom: none; height: 200px;">
    <a style="height: 200px; max-width: 200px;" href="https://conveyor.open.ug" target="_blank"><img alt="Prometheus" src="https://conveyor.open.ug/img/logo.png"
    style="height: 200px; max-width: 200px;"></a>

</h1>

<div align="center">


![Docker Pulls](https://img.shields.io/docker/pulls/openug/conveyor.svg?maxAge=604800)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-ug/conveyor)](https://goreportcard.com/report/github.com/open-ug/conveyor)

</div>

---

The Python Implementation of the Conveyor CI Driver Runtime

## Installation

```sh
pip install conveyor.py
```

## Usage

```py
from conveyor.client import Client

cl = Client()

response = cl.create_or_update_resource_definition(resource_definition)

```

## Development Setup

Check out [this Guide](https://conveyor.open.ug/blog/contributing-to-the-conveyor-ci-driver-runtime#environment-setup) for setting up your development environment.

## License

Apache License 2.0, see [LICENSE](https://github.com/open-ug/conveyor/blob/main/sdk/python/LICENSE).
