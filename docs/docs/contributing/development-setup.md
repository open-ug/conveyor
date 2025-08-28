---
sidebar_position: 2
---

# Development Environment Setup

This guide helps you set up a development environment for developing Conveyor CI.

## Prerequisites

Before you begin development, ensure you have the following software on your dev machine.

- Docker Engine
- Golang

## Starting Dependency Services

First you need to start the Loki container. Within the `config` directory, there is a `compose.dev.yml` file. Start is using Docker Compose.

```sh
docker compose -f config/compose.dev.yml up -d
```

## Running the Application

Once you have the containers running. You can run the application.

First, add a `.env` file in the root of your project.

```env
CONVEYOR_SERVER_HOST=http://localhost:8080
LOKI_ENDPOINT=http://localhost:3100
```

You can then start the api server using the `go run` command.

```sh
make start
```

## API Documentation

Conveyor CI provides comprehensive API documentation using Swagger/OpenAPI. The interactive API documentation is automatically generated from code comments and is available when the server is running.

### Accessing Swagger UI

Once the Conveyor API Server is running, you can access the interactive API documentation at:

**Swagger UI**: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

### Swagger Development Setup

For local development and API documentation generation:

```bash
# Install swag CLI tool
make install-swag

# Generate Swagger documentation
make swagger-init
```

## Testing

Conveyor CI is a Go Application and uses a combination of [stretchr/testify](https://github.com/stretchr/testify) and the go testing package in the standard library.

When you make a change, it is required to write a test for your change as we aim for 100% test coverage. To run your tests, first start all the dependency containers mentioned above then run:

```sh
make test
```

And ensure all tests are passing.
