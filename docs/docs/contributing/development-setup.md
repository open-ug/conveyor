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

First you need to start the dependency containers. Within the `config` directory, there is a `compose.dev.yml` file. Start is using Docker Compose.

```sh
docker compose -f config/compose.dev.yml up -d
```

## Running the Application

Once you have the containers running. You can run the application.

First, add a `.env` file in the root of your project.

```env
CONVEYOR_SERVER_HOST=http://localhost:8080
NATS_URL=nats://localhost:4222
ETCD_ENDPOINT=localhost:2379
LOKI_ENDPOINT=http://localhost:3100
```

You can then start the api server using the `go run` command.

```sh
go run main.go api-server
```
