---
sidebar_position: 1
---

# Installation & Setup

:::danger

Documentation still in early stages

:::


Conveyor CI is designed to be installed easily. It involved a set of programs and software libraries. In this tutorial we shall take you through a guide to setup Conveyor CI for your development Environment.

## Prerequisites

Conveyor CI requires a few prerequisites and dependecies to be able to run smoothly on your local development environment.

- [Docker](https://docs.docker.com/engine/install/). you can install it on the official website.
- The [Go](https://go.dev/doc/install) language Compiler.

Once you have these installed you can run `go version` and `docker version` to verify their installation.

## Downloading Conveyor Tools

Next we shall install the conveyor dependency tools.

We shall begin by downloading the docker compose files and config files for dependecy containers. In a new directory run the following commands to do this.

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

This will start containers that for the Conveyor CI API Server, ETCD Data Store, NATS and Grafana Loki.

## Setting up Devlopment Project

Next we shall setup the development projects by downloading the software libraries. In this tutorial we shall use a Go project and use the Golang SDKs. Within a new directory, create your Go Project.

```bash
go mod init example-project
```

Then install the Go SDKs

```bash
go get conveyor.open.ug

# OR
go get github.com/open-ug/conveyor
```

Congratulations, you now have Conveyor CI installed on your platform.
