<h1 align="center" style="border-bottom: none; height: 200px;">
    <a style="height: 200px; max-width: 200px;" href="https://conveyor.open.ug" target="_blank"><img alt="Conveyor CI Logo" src="https://conveyor.open.ug/logos/logo.svg"
    style="height: 200px; max-width: 200px;"></a>
</h1>

<div align="center">

![Docker Pulls](https://img.shields.io/docker/pulls/openug/conveyor.svg?maxAge=604800)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-ug/conveyor)](https://goreportcard.com/report/github.com/open-ug/conveyor)
![License](https://img.shields.io/github/license/open-ug/conveyor.svg)
![GitHub release](https://img.shields.io/github/v/release/open-ug/conveyor)
![GitHub issues](https://img.shields.io/github/issues/open-ug/conveyor)
![GitHub pull requests](https://img.shields.io/github/issues-pr/open-ug/conveyor)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor?ref=badge_shield)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor?ref=badge_shield&issueType=security)
[![Maintainability](https://qlty.sh/badges/229750f3-9423-4ea6-8528-8e0f8cf854b5/maintainability.svg)](https://qlty.sh/gh/open-ug/projects/conveyor)
[![Code Coverage](https://qlty.sh/badges/229750f3-9423-4ea6-8528-8e0f8cf854b5/test_coverage.svg)](https://qlty.sh/gh/open-ug/projects/conveyor)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10999/badge)](https://www.bestpractices.dev/projects/10999)

</div>

---

**Conveyor CI** is an open-source **lightweight engine for building distributed CI/CD systems with ease.**.

Instead of building your own CI/CD system from scratch, Conveyor CI gives you a modular toolkit, SDKs, APIs, and drivers that handle the hard parts: execution, events, scaling, observability, and more.

## Key features

- **Built-in Observability**: Metrics, tracing, and logging integrated out-of-the-box.
- **Authentication & Security**: TLS certificates and JWT token-based authentication for secure API access.
- **Lightweight & Modular**: Core engine with extensible driver architecture.
- **Event-Driven**: Real-time pipeline execution with comprehensive event streaming.
- **Scalable**: Distributed architecture designed for cloud-native environments.
- **Developer-Friendly**: Rich APIs, SDKs, and documentation for easy integration.
- **Declarative CI/CD Workflows**: Define pipelines using a flexible, extensible API.
- **Real-time Event System**: Event-driven architecture powers responsive, dynamic execution.
- **Effortless Horizontal Scaling**: Built-in runtime scales drivers across distributed systems with zero extra code.
- **Live Log Management**: Stream and store logs in real time from every running task.


## Installation

Conveyor CI is distributed as an OCI container and available on Docker Hub. It depends on `etcd`, `loki`, and `nats`, so a standard Docker Compose configuration is provided.

> **Helm charts coming soon**

To Install it on docker compose you can head over to the Releases page and download `compose.yml` and `loki.yml` or on a linux system you can download them using `curl`.

```sh
curl -s https://api.github.com/repos/open-ug/conveyor/releases/latest | grep browser_download_url | grep compose.yml | cut -d '"' -f 4 | xargs curl -L -o compose.yml

curl -s https://api.github.com/repos/open-ug/conveyor/releases/latest | grep browser_download_url | grep loki.yml | cut -d '"' -f 4 | xargs curl -L -o loki.yml
```

Next start the containers using docker compose.

```sh
docker compose up

# OR

docker compose up -d
```

The Conveyor API Server will be reachable on [http://localhost:8080](http://localhost:8080)

**Note**: For production deployments, configure authentication using TLS certificates and JWT tokens. See the [Authentication Documentation](docs/authentication.md) for detailed setup instructions.

## Authentication

Conveyor CI supports secure authentication using TLS certificates and JWT tokens. This ensures that only trusted clients can connect and interact with the API server.

**Key features**:
- TLS mutual authentication with client certificates
- JWT token-based authorization  
- Certificate Authority (CA) based trust model
- Configurable authentication levels (required, optional, or disabled)

For detailed setup instructions, see [docs/authentication.md](docs/authentication.md).

## More information

Visit the [official documentation](https://conveyor.open.ug). for architecture, SDK usage, and driver development.

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)

## Project Governance

Refer to [Governance Document](https://conveyor.open.ug/docs/contributing/governance)

## License

Apache License 2.0, see [LICENSE](./LICENSE). Copyright © 2024 - Present, Conveyor CI Contributors


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor?ref=badge_large)
