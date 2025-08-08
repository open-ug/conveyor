---
sidebar_position: 1
---

# Contributing

Conveyor CI is an open-source project licenced under the [Apache License 2.0](https://github.com/open-ug/conveyor/blob/main/LICENSE). We welcome anyone who would be interested in contributing to urunc. As a first step, please take a look at the following document. The current document provides a high level overview of Conveyor CI's code structure, along with a few guidelines regarding contributions to the project.


## Code Organisation

Conveyor CI is written in Go and it exposes an Application Programming Interface whose client libraries can be written in any language. We structure the codebase and other files as follows:

- `/`: The root directory contains the non-code files, such as the licence, readme and conveyor's entry point `main.go`.
- `/docs`: This directory contains a [Docusaurous Site](https://docusaurus.io/) for all the documentation related to Conveyor CI.
- `/cmd/cli`: Contains the entrypoint for the command line and definitions for all CLI commands.
- `/cmd/api`: Contains the entrypoint for the API Server.
- `/internal/`: Contains majority of the Conveyor CI codebase.
