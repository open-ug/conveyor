---
sidebar_position: 2
---

# Development Environment Setup

This guide helps you set up a development environment for developing Conveyor CI.

## Prerequisites

Before you begin development, ensure you have the following software on your dev machine.

- [Go](https://go.dev/)
- A Linux environment. This is mainly because Conveyor CI is a Linux only application. If you are contributing to the documentation site or an SDK you can ignore this.


## Running the API Server

Start the API Server by running the following commands.

```sh
sudo go run main.go init

make start
```

## Contributing to the documentation site

The documentation site is located in the `docs/` directiory and is a [docusaurus](https://docusaurus.io/) application. so ensure to `cd docs` and `npm install` to install the dependencies.

You can start the docs site by running

```sh
make docs
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
